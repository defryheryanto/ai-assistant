package calendar

import "context"

//go:generate mockgen -source service.go -package mock -destination mock/mock.go

// Service handles interactions with a calendar provider.
type Service interface {
	CreateEvent(ctx context.Context, params CreateEventParams) (string, error)
}
