package user

import (
	"context"
	"fmt"
)

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

func (s *service) Create(ctx context.Context, params CreateUserParams) (int64, error) {
	role := params.Role
	if role == "" || (role != RoleUser && role != RoleAdmin) {
		role = RoleUser
	}

	id, err := s.repository.Insert(ctx, &User{
		Name:        params.Name,
		WhatsAppJID: fmt.Sprintf("%s@s.whatsapp.net", params.Phone),
		Role:        role,
		Email:       params.Email,
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}
