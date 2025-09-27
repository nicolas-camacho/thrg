package token

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

const TokenLength = 16

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func generateRandomToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (r *Repository) CreateNewToken(ctx context.Context, adminID uint) (string, error) {
	var tokenValue string

	for {
		tokenValue = generateRandomToken(TokenLength)

		var existing RegistrationToken
		if err := r.db.WithContext(ctx).Where("value = ?", tokenValue).First(&existing).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return "", fmt.Errorf("error checking existing token: %w", err)
			}
			break
		}
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	newToken := RegistrationToken{
		Value:       tokenValue,
		CreatedByID: adminID,
		ExpiresAt:   &expirationTime,
	}

	result := r.db.WithContext(ctx).Create(&newToken)
	if result.Error != nil {
		return "", fmt.Errorf("error creating new token: %w", result.Error)
	}

	return tokenValue, nil
}

func (r *Repository) GetAllTokens(ctx context.Context) ([]RegistrationToken, error) {
	var tokens []RegistrationToken

	result := r.db.WithContext(ctx).Order("created_at desc").Find(&tokens)
	if result.Error != nil {
		return nil, fmt.Errorf("error retrieving tokens: %w", result.Error)
	}
	return tokens, nil
}
