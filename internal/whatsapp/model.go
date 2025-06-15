package whatsapp

import "time"

type Group struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	GroupJID     string    `json:"group_jid"`
	IsActive     bool      `json:"is_active"`
	RegisteredBy string    `json:"registered_by"`
}
