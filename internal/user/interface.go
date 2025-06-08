package user

import "context"

//go:generate mockgen -source interface.go -package mock -destination mock/mock.go

type Service interface {
	GetUserByWhatsAppJID(ctx context.Context, jid string) (*User, error)
}

type Repository interface {
	FindUserByWhatsAppJID(ctx context.Context, jid string) (*User, error)
}
