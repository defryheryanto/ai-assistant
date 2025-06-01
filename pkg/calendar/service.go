package calendar

type Service interface {
	CreateEvent(params CreateEventParams) (string, error)
}
