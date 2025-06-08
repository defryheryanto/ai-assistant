package user

import "context"

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetUserByWhatsAppJID(ctx context.Context, jid string) (*User, error) {
	return s.repository.FindUserByWhatsAppJID(ctx, jid)
}
