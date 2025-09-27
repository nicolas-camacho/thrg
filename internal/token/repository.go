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

	newToken := RegistrationToken{
		Value:       tokenValue,
		CreatedByID: adminID,
	}

	result := r.db.WithContext(ctx).Create(&newToken)
	if result.Error != nil {
		return "", fmt.Errorf("error creating new token: %w", result.Error)
	}

	return tokenValue, nil
}
