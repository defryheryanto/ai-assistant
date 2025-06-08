package app

import (
	"context"
	"database/sql"

	"github.com/defryheryanto/ai-assistant/internal/user"
	userrepository "github.com/defryheryanto/ai-assistant/internal/user/repository"
)

type Repositories struct {
	UserRepository user.Repository
}

func SetupRepository(ctx context.Context, db *sql.DB) *Repositories {
	userRepository := userrepository.New(db)

	return &Repositories{
		UserRepository: userRepository,
	}
}
