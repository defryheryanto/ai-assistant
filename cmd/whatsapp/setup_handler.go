package main

import (
	"context"
	"log"
	"strings"

	"github.com/defryheryanto/ai-assistant/pkg/tools"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func eventHandler(ctx context.Context, client *whatsmeow.Client, toolRegistry tools.Registry) whatsmeow.EventHandler {
	return func(evt any) {
		switch v := evt.(type) {
		case *events.Message:
			// Only respond to a personal message
			chatJID := v.Info.MessageSource.Chat.String()
			if strings.HasSuffix(chatJID, "@s.whatsapp.net") {
				message := getMessage(v)
				resp, err := toolRegistry.Execute(ctx, message)
				if err != nil {
					log.Printf("error executing tool: %v\n", err)
				}

				_, err = client.SendMessage(ctx, v.Info.Chat, &waE2E.Message{
					Conversation: proto.String(resp),
				})
				if err != nil {
					log.Printf("error sending response message: %v\n", err)
				}
			}
		}
	}
}

func getMessage(evt *events.Message) string {
	if evt.Message.GetConversation() != "" {
		return evt.Message.GetConversation()
	}
	if evt.Message.GetExtendedTextMessage().Text != nil {
		return evt.Message.GetExtendedTextMessage().GetText()
	}
	return ""
}
