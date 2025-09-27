package token

import "gorm.io/gorm"

type RegistrationToken struct {
	gorm.Model
	Value       string `gorm:"uniqueIndex;not null"`
	IsUsed      bool   `gorm:"default:false"`
	UsedByID    *uint
	ExpiresAt   *int64 `gorm:"autoCreateTime:milli"`
	CreatedByID uint
}
