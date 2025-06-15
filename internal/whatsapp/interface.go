package whatsapp

import "context"

//go:generate mockgen -source interface.go -package mock -destination mock/mock.go

type GroupRepository interface {
	FindByJID(ctx context.Context, jid string) (*Group, error)
	Insert(ctx context.Context, data *Group) (int64, error)
}

type GroupService interface {
	GetByJID(ctx context.Context, jid string) (*Group, error)
	Create(ctx context.Context, params CreateGroupParams) (int64, error)
}

type UserRepository interface {
	FindByJID(ctx context.Context, jid string) (*User, error)
	Insert(ctx context.Context, data *User) (int64, error)
}

type UserService interface {
	GetByJID(ctx context.Context, jid string) (*User, error)
	Create(ctx context.Context, params CreateUserParams) (int64, error)
}
