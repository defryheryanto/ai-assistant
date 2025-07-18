package tools

import (
	"context"
	"io"
	"net/http"

	"github.com/defryheryanto/ai-assistant/internal/contextgroup"
	"github.com/defryheryanto/ai-assistant/pkg/tools"
	"github.com/tmc/langchaingo/llms"
	"go.mau.fi/whatsmeow"
	waE2E "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

// GenerateImageTool sends generated image URLs to WhatsApp.
type GenerateImageTool struct {
	base   tools.Tool
	client *whatsmeow.Client
}

// NewGenerateImageTool wraps another tool and forwards images to WhatsApp when possible.
func NewGenerateImageTool(base tools.Tool, client *whatsmeow.Client) *GenerateImageTool {
	return &GenerateImageTool{base: base, client: client}
}

func (t *GenerateImageTool) SystemPrompt() string  { return t.base.SystemPrompt() }
func (t *GenerateImageTool) Definition() llms.Tool { return t.base.Definition() }

func (t *GenerateImageTool) Execute(ctx context.Context, call llms.ToolCall) (*llms.MessageContent, error) {
	resp, err := t.base.Execute(ctx, call)
	if err != nil {
		return resp, err
	}

	if t.client != nil {
		waCtx := contextgroup.GetWhatsAppContext(ctx)
		if waCtx != nil && resp != nil && len(resp.Parts) > 0 {
			if r, ok := resp.Parts[0].(llms.ToolCallResponse); ok {
				contextgroup.MarkMediaSent(ctx)
				go t.sendImage(ctx, waCtx.CurrentChatJID, r.Content)
			}
		}
	}

	return resp, nil
}

func (t *GenerateImageTool) sendImage(ctx context.Context, jidStr, url string) {
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

	jid, err := types.ParseJID(jidStr)
	if err != nil {
		return
	}
	_, err = t.client.SendMessage(ctx, jid, &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           proto.String(upload.URL),
			DirectPath:    proto.String(upload.DirectPath),
			MediaKey:      upload.MediaKey,
			FileLength:    proto.Uint64(uint64(len(data))),
			Mimetype:      proto.String(resp.Header.Get("Content-Type")),
			FileEncSHA256: upload.FileEncSHA256,
			FileSHA256:    upload.FileSHA256,
		},
	})
	if err != nil {
		return
	}
}
