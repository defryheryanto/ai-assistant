package app

import (
	"context"
	"database/sql"

	"github.com/defryheryanto/ai-assistant/internal/user"
	userrepository "github.com/defryheryanto/ai-assistant/internal/user/repository"
	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	whatsapprepository "github.com/defryheryanto/ai-assistant/internal/whatsapp/repository"
)

type Repositories struct {
	UserRepository          user.Repository
	WhatsAppGroupRepository whatsapp.GroupRepository
}

func SetupRepository(ctx context.Context, db *sql.DB) *Repositories {
	userRepository := userrepository.New(db)
	whatsappGroupRepository := whatsapprepository.NewGroupRepository(db)

	return &Repositories{
		UserRepository:          userRepository,
		WhatsAppGroupRepository: whatsappGroupRepository,
	}
}
