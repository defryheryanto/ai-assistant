package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/defryheryanto/ai-assistant/config"
	"github.com/defryheryanto/ai-assistant/internal/app"
	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var errSenderNotAuthenticated = fmt.Errorf("sender not authenticated")
var errGroupNotAuthenticated = fmt.Errorf("group not authenticated")
var errUserNotMentioned = fmt.Errorf("ignoring message, user not mentioned")

type EventHandler struct {
	client       *whatsmeow.Client
	toolRegistry tools.Registry
	services     *app.Services
}

func (h *EventHandler) Handle(ctx context.Context) whatsmeow.EventHandler {
	return func(evt any) {
		switch v := evt.(type) {
		case *events.Message:
			var err error
			if !v.Info.IsGroup {
				ctx, err = h.authenticateSender(ctx, v, true)
				if err != nil && err != errSenderNotAuthenticated {
					log.Printf("error authenticating sender: %v\n", err)
					return
				}
				if err == errSenderNotAuthenticated && config.IsUserWhitelistEnabled {
					return
				}
			} else {
				ctx, err = h.authenticateGroup(ctx, v)
				if err != nil && err != errSenderNotAuthenticated && err != errGroupNotAuthenticated && err != errUserNotMentioned {
					log.Printf("error authenticating group: %v\n", err)
					return
				}
				if err == errUserNotMentioned || err == errGroupNotAuthenticated {
					return
				}

			}

			textMessage := ""
			switch {
			case getMessage(v) != "":
				if !hasActualText(v) {
					return
				}
				textMessage = getMessage(v)
			case v.Message.GetAudioMessage() != nil:
				audioMessage := v.Message.GetAudioMessage()

				audioFileName := fmt.Sprintf("%s/transcriptions/%s.wav", config.TempFolderPath, uuid.New().String())
				f, _ := os.Create(audioFileName)
				err := h.client.DownloadToFile(ctx, audioMessage, f)
				if err != nil {
					log.Printf("error downloading audio: %v\n", err)
					return
				}
				f.Close()
				defer os.Remove(audioFileName)

				ff, err := os.Open(audioFileName)
				if err != nil {
					log.Printf("error opening audio: %v\n", err)
					return
				}
				defer ff.Close()

				res, err := h.services.OpenAIClient.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
					Model: openai.AudioModelWhisper1,
					File:  ff,
				})
				if err != nil {
					log.Printf("error transcripting audio: %v\n", err)
					return
				}

				textMessage = res.Text
			}

			if textMessage == "" {
				return
			}

			respCtx := &contextgroup.ResponseContext{}
			ctx = contextgroup.SetResponseContext(ctx, respCtx)

			resp, err := h.toolRegistry.Execute(ctx, v.Info.Chat.ToNonAD().String(), textMessage)
			if err != nil {
				log.Printf("error executing tool: %v\n", err)
				return
			}

			if respCtx.MediaSent {
				return
			}

			_, err = h.client.SendMessage(ctx, v.Info.Chat, &waE2E.Message{
				Conversation: proto.String(resp),
			})
			if err != nil {
				log.Printf("error sending response message: %v\n", err)
				return
			}
		}
	}
}

func getMessage(evt *events.Message) string {
	if evt.Message.GetConversation() != "" {
		return evt.Message.GetConversation()
	}
	if evt.Message.GetExtendedTextMessage() != nil && evt.Message.GetExtendedTextMessage().Text != nil {
		return evt.Message.GetExtendedTextMessage().GetText()
	}
	return ""
}

func hasActualText(evt *events.Message) bool {
	text := getMessage(evt)
	if text == "" {
		return false
	}
	ext := evt.Message.GetExtendedTextMessage()
	if ext != nil && ext.GetContextInfo() != nil {
		for _, jid := range ext.GetContextInfo().GetMentionedJID() {
			if idx := strings.Index(jid, "@"); idx > 0 {
				number := jid[:idx]
				text = strings.ReplaceAll(text, "@"+number, "")
			}
		}
	}
	return strings.TrimSpace(text) != ""
}
func (h *EventHandler) authenticateSender(ctx context.Context, evt *events.Message, setWhatsAppContext bool) (context.Context, error) {
	senderJID := evt.Info.Sender.ToNonAD().String()
	usr, err := h.services.UserService.GetByJID(ctx, senderJID)
	if err != nil {
		return ctx, err
	}
	if usr == nil {
		return ctx, errSenderNotAuthenticated
	}

	ctx = contextgroup.SetUserContext(ctx, &contextgroup.UserContext{
		ID:          usr.ID,
		Name:        usr.Name,
		WhatsAppJID: usr.WhatsAppJID,
		Role:        string(usr.Role),
		Email:       usr.Email,
	})
	if setWhatsAppContext {
		chatJID := evt.Info.Chat.ToNonAD().String()
		ctx = contextgroup.SetWhatsAppContext(ctx, &contextgroup.WhatsAppContext{
			CurrentChatJID: chatJID,
			SenderJID:      senderJID,
		})
	}

	return ctx, nil
}

func (h *EventHandler) authenticateGroup(ctx context.Context, evt *events.Message) (context.Context, error) {
	if !evt.Info.IsGroup {
		return ctx, nil
	}

	// Process only mentioned message in group
	mentionedJIDs := evt.Message.GetExtendedTextMessage().GetContextInfo().GetMentionedJID()
	currentLoggedInUserJID := h.client.Store.ID.ToNonAD().String()
	if !slices.Contains(mentionedJIDs, currentLoggedInUserJID) {
		return ctx, errUserNotMentioned
	}

	var err error
	isSenderAuthenticated := true
	ctx, err = h.authenticateSender(ctx, evt, false)
	if err != nil && err != errSenderNotAuthenticated {
		return ctx, err
	}
	if err == errSenderNotAuthenticated {
		isSenderAuthenticated = false
	}

	chatJID := evt.Info.Chat.ToNonAD().String()
	if config.IsWhatsAppGroupWhitelistEnabled || !isSenderAuthenticated {
		whatsappGroup, err := h.services.WhatsAppGroupService.GetByJID(ctx, chatJID)
		if err != nil {
			return ctx, err
		}
		if whatsappGroup == nil {
			return ctx, errGroupNotAuthenticated
		}
	}

	senderJID := evt.Info.Sender.ToNonAD().String()
	ctx = contextgroup.SetWhatsAppContext(ctx, &contextgroup.WhatsAppContext{
		CurrentChatJID: chatJID,
		SenderJID:      senderJID,
	})

	return ctx, nil
}
