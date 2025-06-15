package contextgroup

import "context"

type WhatsAppContext struct {
	CurrentChatJID string
	SenderJID      string
}

type whatsappKeyType string

var whatsappKey whatsappKeyType = "whatsapp_key"

func SetWhatsAppContext(ctx context.Context, data *WhatsAppContext) context.Context {
	return context.WithValue(ctx, whatsappKey, data)
}

func GetWhatsAppContext(ctx context.Context) *WhatsAppContext {
	currentSession, ok := ctx.Value(whatsappKey).(*WhatsAppContext)
	if !ok || currentSession == nil {
		return nil
	}

	return currentSession
}
