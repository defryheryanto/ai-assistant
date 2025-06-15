package whatsapp

import "context"

type groupService struct {
	repository GroupRepository
}

func NewGroupService(repository GroupRepository) GroupService {
	return &groupService{
		repository: repository,
	}
}

func (s *groupService) GetByJID(ctx context.Context, jid string) (*Group, error) {
	return s.repository.FindByJID(ctx, jid)
}

func (s *groupService) Create(ctx context.Context, params CreateGroupParams) (int64, error) {
	return s.repository.Insert(ctx, &Group{
		GroupJID:     params.GroupJID,
		IsActive:     true,
		RegisteredBy: params.RegisteredBy,
	})
}
