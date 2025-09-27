package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	RoleAdmin  = "admin"
	RolePlayer = "player"
)

type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"default:player"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}
