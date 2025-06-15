package app

import (
	"context"
	"database/sql"

	"github.com/defryheryanto/ai-assistant/internal/whatsapp"
	whatsapprepository "github.com/defryheryanto/ai-assistant/internal/whatsapp/repository"
)

type Repositories struct {
	UserRepository          whatsapp.UserRepository
	WhatsAppGroupRepository whatsapp.GroupRepository
}

func SetupRepository(ctx context.Context, db *sql.DB) *Repositories {
	userRepository := whatsapprepository.NewUserRepository(db)
	whatsappGroupRepository := whatsapprepository.NewGroupRepository(db)

	return &Repositories{
		UserRepository:          userRepository,
		WhatsAppGroupRepository: whatsappGroupRepository,
	}
}
