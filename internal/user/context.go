package user

import "context"

type userKey string

var key userKey = "user_key"

func SetUserToContext(ctx context.Context, data *User) context.Context {
	return context.WithValue(ctx, key, data)
}

func GetUserFromContext(ctx context.Context) *User {
	currentSession, ok := ctx.Value(key).(*User)
	if !ok || currentSession == nil {
		return nil
	}

	return currentSession
}
