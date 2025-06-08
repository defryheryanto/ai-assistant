package calendar

import "context"

//go:generate mockgen -source service.go -package mock -destination mock/mock.go

type Service interface {
	CreateEvent(ctx context.Context, params CreateEventParams) (string, error)
}
