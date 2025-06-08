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

func (r *Repository) Insert(ctx context.Context, data *user.User) (int64, error) {
	var id int64
	_, err := r.db.ExecContext(
		ctx,
		queryInsert,
		data.Name,
		data.WhatsAppJID,
		data.Role,
		data.Email,
	)
	if err != nil {
		return 0, err
	}

	return id, nil
}
