package token

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TokenModelBase struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RegistrationToken struct {
	TokenModelBase
	Value       string `gorm:"uniqueIndex;not null"`
	IsUsed      bool   `gorm:"default:false"`
	UsedByID    *uuid.UUID
	ExpiresAt   *time.Time
	CreatedByID uuid.UUID
}
