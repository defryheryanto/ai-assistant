package repository

import (
	"context"
	"database/sql"

	"github.com/defryheryanto/ai-assistant/internal/user"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindUserByWhatsAppJID(ctx context.Context, jid string) (*user.User, error) {
	var res user.User
	err := r.db.QueryRowContext(ctx, queryFindUserByWhatsAppJID, jid).Scan(
		&res.ID,
		&res.Name,
		&res.WhatsAppJID,
		&res.Role,
		&res.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &res, nil
}
