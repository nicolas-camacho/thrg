package token

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type TokenDTO struct {
	Value          string    `json:"Value"`
	IsUsed         bool      `json:"IsUsed"`
	UsedByUsername string    `json:"UsedByUsername"`
	CreatedAt      time.Time `json:"CreatedAt"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateNewToken(ctx context.Context, adminID uuid.UUID) (string, error) {
	tokenUUID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating UUID: %w", err)
	}
	tokenValue := tokenUUID.String()

	newToken := RegistrationToken{
		Value:       tokenValue,
		CreatedByID: adminID,
	}

	result := r.db.WithContext(ctx).Create(&newToken)
	if result.Error != nil {
		return "", fmt.Errorf("error creating token: %w", result.Error)
	}
	return tokenValue, nil
}

func (r *Repository) ValidateAndUseToken(ctx context.Context, tokenValue string, usedByID uuid.UUID) (*RegistrationToken, error) {
	var token RegistrationToken

	result := r.db.WithContext(ctx).Where("value = ?", tokenValue).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("token not found")
		}
		return nil, fmt.Errorf("error retrieving token: %w", result.Error)
	}

	if token.IsUsed {
		return nil, fmt.Errorf("token has already been used")
	}

	token.IsUsed = true
	token.UsedByID = &usedByID

	result = r.db.WithContext(ctx).Save(&token)
	if result.Error != nil {
		return nil, fmt.Errorf("error updating token as used: %w", result.Error)
	}

	return &token, nil
}

func (r *Repository) GetAllTokens(ctx context.Context) ([]RegistrationToken, error) {
	var tokens []RegistrationToken

	result := r.db.WithContext(ctx).Order("created_at desc").Find(&tokens)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving tokens: %w", result.Error)
	}
	return tokens, nil
}
