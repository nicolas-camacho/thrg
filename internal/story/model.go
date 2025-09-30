package story

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoryBase struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Story struct {
	StoryBase
	HolderName          string `gorm:"uniqueIndex;not null"`
	Title               string `gorm:"not null"`
	Description         string
	MisfortuneThreshold float64 `gorm:"not null"`
	Acts                []Act
}

type Act struct {
	StoryBase
	StoryID uuid.UUID `gorm:"type:uuid;not null"`
	Order   int       `gorm:"not null"`
	Text    string    `gorm:"not null"`
	Options []Option  `gorm:"foreignKey:ActID"`
}

type Option struct {
	StoryBase
	ActID        uuid.UUID     `gorm:"type:uuid;not null"`
	Text         string        `gorm:"not null"`
	NextAct      *uuid.UUID    `gorm:"type:uuid"`
	Consequences []Consequence `gorm:"foreignKey:OptionID"`
}

type ConsequenceType string

const (
	TypeLocura     ConsequenceType = "locura"
	TypePanico     ConsequenceType = "panico"
	TypeAnsiedad   ConsequenceType = "ansiedad"
	TypeBrillantes ConsequenceType = "brillantes"
	TypeMisfortune ConsequenceType = "desgracia"
)

type Consequence struct {
	gorm.Model
	OptionID uuid.UUID       `gorm:"type:uuid;not null"`
	Type     ConsequenceType `gorm:"not null"`
	Value    float64         `gorm:"not null"`
}
