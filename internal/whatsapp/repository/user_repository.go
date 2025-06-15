package repository

import (
	"context"
	"database/sql"

	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByJID(ctx context.Context, jid string) (*whatsapp.User, error) {
	var res whatsapp.User
	err := r.db.QueryRowContext(ctx, queryFindUserByJID, jid).Scan(
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

func (r *UserRepository) Insert(ctx context.Context, data *whatsapp.User) (int64, error) {
	var id int64
	_, err := r.db.ExecContext(
		ctx,
		queryInsertUser,
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
