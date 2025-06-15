package whatsapp

import "context"

type GroupRepository interface {
	FindByJID(ctx context.Context, jid string) (*Group, error)
	Insert(ctx context.Context, data *Group) (int64, error)
}

type GroupService interface {
	GetByJID(ctx context.Context, jid string) (*Group, error)
	Create(ctx context.Context, params CreateGroupParams) (int64, error)
}
