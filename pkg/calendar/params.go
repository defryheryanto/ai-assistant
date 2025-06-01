package calendar

type CreateEventParams struct {
	Summary     string   `json:"summary"`
	Description string   `json:"Description"`
	Location    string   `json:"Location"`
	Start       string   `json:"Start"`
	End         string   `json:"End"`
	Attendees   []string `json:"Attendees"`
}
