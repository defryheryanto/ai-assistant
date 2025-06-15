package whatsapp

import (
	"context"
	"fmt"
)

type userService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) GetByJID(ctx context.Context, jid string) (*User, error) {
	return s.repository.FindByJID(ctx, jid)
}

func (s *userService) Create(ctx context.Context, params CreateUserParams) (int64, error) {
	role := params.Role
	if role == "" || (role != UserRoleUser && role != UserRoleAdmin) {
		role = UserRoleUser
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
