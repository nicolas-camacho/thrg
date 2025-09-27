package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CheckAdminExists(ctx context.Context) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).Model(&User{}).Where("role = ?", RoleAdmin).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check admin existence: %w", result.Error)
	}
	return count > 0, nil
}

func (r *Repository) CreateAdmin(ctx context.Context, username, password string) error {
	admin := User{
		Username: username,
		Role:     RoleAdmin,
	}

	if err := admin.SetPassword(password); err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	result := r.db.WithContext(ctx).Create(&admin)
	if result.Error != nil {
		return fmt.Errorf("failed to create admin user: %w", result.Error)
	}

	return nil
}

func (r *Repository) Authenticate(ctx context.Context, username, password string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}

		return nil, fmt.Errorf("failed to query user: %w", result.Error)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
