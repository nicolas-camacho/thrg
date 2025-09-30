package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nicolas-camacho/thrg/internal/core"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepositoryLookup interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*core.UserLookupModel, error)
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func NewUserRepositoryLookup(repo *Repository) UserRepositoryLookup {
	return repo
}

func (r *Repository) CheckAdminExists(ctx context.Context) (bool, error) {
	var count int64

	result := r.db.WithContext(ctx).Model(&User{}).Where("role = ?", RoleAdmin).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check admin existence: %w", result.Error)
	}
	return count > 0, nil
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

func (r *Repository) CreateUser(ctx context.Context, username, password, role string) (*User, error) {
	newUser := User{
		Username: username,
		Role:     role,
	}

	if err := newUser.SetPassword(password); err != nil {
		return nil, fmt.Errorf("failed to set password: %w", err)
	}

	result := r.db.WithContext(ctx).Create(&newUser)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("username already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", result.Error)
	}
	return &newUser, nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&User{}, userID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", userID)
	}
	return nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (*core.UserLookupModel, error) {
	var user User

	result := r.db.WithContext(ctx).First(&user, "id = ?", userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", result.Error)
	}
	return &core.UserLookupModel{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

func (r *Repository) GetAllPlayers(ctx context.Context) ([]User, error) {
	var users []User

	result := r.db.WithContext(ctx).Where("role = ?", RolePlayer).Order("created_at DESC").Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all players: %w", result.Error)
	}
	return users, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).First(&user, "username = ?", username)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", result.Error)
	}
	return &user, nil
}
