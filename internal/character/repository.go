package character

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateCharacter(ctx context.Context, userID uuid.UUID) (*Character, error) {
	newCharacter := Character{
		UserID: userID,
	}

	if err := r.db.WithContext(ctx).Create(&newCharacter).Error; err != nil {
		return nil, fmt.Errorf("error creating character: %w", err)
	}

	return &newCharacter, nil
}

func (r *Repository) GetCharacterByID(ctx context.Context, characterID uuid.UUID) (*Character, error) {
	var character Character
	if err := r.db.WithContext(ctx).First(&character, "id = ?", characterID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting character by ID: %w", err)
	}
	return &character, nil
}

func (r *Repository) GetCharacterByUserID(ctx context.Context, userID uuid.UUID) (*Character, error) {
	var character Character
	if err := r.db.WithContext(ctx).First(&character, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting character by user ID: %w", err)
	}
	return &character, nil
}

func (r *Repository) UpdateCharacter(ctx context.Context, character *Character) error {
	if err := r.db.WithContext(ctx).Save(character).Error; err != nil {
		return fmt.Errorf("error updating character: %w", err)
	}
	return nil
}
