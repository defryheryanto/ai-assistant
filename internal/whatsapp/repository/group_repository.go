package repository

import (
	"context"
	"database/sql"

	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{
		db: db,
	}
}

func (r *GroupRepository) FindByJID(ctx context.Context, jid string) (*whatsapp.Group, error) {
	var res whatsapp.Group
	err := r.db.QueryRowContext(ctx, queryFindGroupByJID, jid).Scan(
		&res.ID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.GroupJID,
		&res.IsActive,
		&res.RegisteredBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &res, nil
}

func (r *GroupRepository) Insert(ctx context.Context, data *whatsapp.Group) (int64, error) {
	var id int64
	_, err := r.db.ExecContext(
		ctx,
		queryInsertGroup,
		data.GroupJID,
		data.IsActive,
		data.RegisteredBy,
	)
	if err != nil {
		return 0, err
	}

	return id, nil
}
