package contextgroup

import "context"

type UserContext struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	WhatsAppJID string `json:"whatsapp_jid"`
	Role        string `json:"role"`
	Email       string `json:"email"`
}

type userKeyType string

var userKey userKeyType = "user_key"

func SetUserContext(ctx context.Context, data *UserContext) context.Context {
	return context.WithValue(ctx, userKey, data)
}

func GetUserContext(ctx context.Context) *UserContext {
	currentSession, ok := ctx.Value(userKey).(*UserContext)
	if !ok || currentSession == nil {
		return nil
	}

	return currentSession
}
