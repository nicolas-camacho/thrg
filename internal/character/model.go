package character

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Character struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`

	CurrentStoryID *uuid.UUID `gorm:"type:uuid"`
	CurrentActID   *uuid.UUID `gorm:"type:uuid"`

	Misfortune float64 `gorm:"not null;default:0"`
	Locura     float64 `gorm:"not null;default:0"`
	Panico     float64 `gorm:"not null;default:0"`
	Ansiedad   float64 `gorm:"not null;default:0"`
	Brillantes float64 `gorm:"not null;default:0"`
}
