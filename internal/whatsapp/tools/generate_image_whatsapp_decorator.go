package tools

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	"github.com/tmc/langchaingo/llms"
	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// ImageWhatsAppDecorator sends generated image URLs to WhatsApp.
type ImageWhatsAppDecorator struct {
	base   tools.Tool
	client *whatsmeow.Client
}

// NewImageWhatsAppDecorator wraps a GenerateImageTool with WhatsApp delivery.
func NewImageWhatsAppDecorator(base tools.Tool, client *whatsmeow.Client) *ImageWhatsAppDecorator {
	return &ImageWhatsAppDecorator{base: base, client: client}
}

func (t *ImageWhatsAppDecorator) SystemPrompt() string  { return t.base.SystemPrompt() }
func (t *ImageWhatsAppDecorator) Definition() llms.Tool { return t.base.Definition() }

func (t *ImageWhatsAppDecorator) Execute(ctx context.Context, call llms.ToolCall) (*llms.MessageContent, error) {
	resp, err := t.base.Execute(ctx, call)
	if err != nil {
		return resp, err
	}

	if t.client != nil {
		waCtx := contextgroup.GetWhatsAppContext(ctx)
		if waCtx != nil && resp != nil && len(resp.Parts) > 0 {
			if r, ok := resp.Parts[0].(llms.ToolCallResponse); ok {
				go t.sendImage(ctx, waCtx.CurrentChatJID, r.Content)
			}
		}
	}

	return resp, nil
}

func (t *ImageWhatsAppDecorator) sendImage(ctx context.Context, jidStr, url string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	upload, err := t.client.Upload(ctx, data, whatsmeow.MediaImage)
	if err != nil {
		return
	}

	parts := strings.Split(jidStr, "@")
	jid := types.JID{User: parts[0], Server: parts[1]}
	_, _ = t.client.SendMessage(ctx, jid, &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			Url:        proto.String(upload.URL),
			DirectPath: proto.String(upload.DirectPath),
			MediaKey:   upload.MediaKey,
			FileLength: proto.Uint64(uint64(len(data))),
			Mimetype:   proto.String(resp.Header.Get("Content-Type")),
		},
	})
}
