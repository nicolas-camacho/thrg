package token

import (
	"time"

	"gorm.io/gorm"
)

type RegistrationToken struct {
	gorm.Model
	Value       string `gorm:"uniqueIndex;not null"`
	IsUsed      bool   `gorm:"default:false"`
	UsedByID    *uint
	ExpiresAt   *time.Time
	CreatedByID uint
}
