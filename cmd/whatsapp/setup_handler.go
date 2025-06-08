package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/defryheryanto/ai-assistant/config"
	"github.com/defryheryanto/ai-assistant/internal/app"
	"github.com/defryheryanto/ai-assistant/internal/user"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	"github.com/openai/openai-go"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func eventHandler(ctx context.Context, client *whatsmeow.Client, toolRegistry tools.Registry, services *app.Services) whatsmeow.EventHandler {
	return func(evt any) {
		switch v := evt.(type) {
		case *events.Message:
			// Only respond to a personal message
			chatJID := v.Info.MessageSource.Chat.String()
			if !strings.HasSuffix(chatJID, "@s.whatsapp.net") {
				return
			}

			if config.IsUserWhitelistEnabled {
				usr, err := services.UserService.GetUserByWhatsAppJID(ctx, chatJID)
				if err != nil {
					log.Printf("error getting user: %v\n", err)
					return
				}
				if usr == nil {
					return
				}

				ctx = user.SetUserToContext(ctx, usr)
			}

			textMessage := ""
			switch {
			case getMessage(v) != "":
				textMessage = getMessage(v)
			case v.Message.GetAudioMessage() != nil:
				audioMessage := v.Message.GetAudioMessage()

				f, _ := os.Create("./audio.wav")
				err := client.DownloadToFile(ctx, audioMessage, f)
				if err != nil {
					log.Printf("error downloading audio: %v\n", err)
					return
				}
				f.Close()

				ff, err := os.Open("./audio.wav")
				if err != nil {
					log.Printf("error opening audio: %v\n", err)
					return
				}
				defer ff.Close()

				res, err := services.OpenAIClient.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
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

			resp, err := toolRegistry.Execute(ctx, textMessage)
			if err != nil {
				log.Printf("error executing tool: %v\n", err)
				return
			}

			_, err = client.SendMessage(ctx, v.Info.Chat, &waE2E.Message{
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
