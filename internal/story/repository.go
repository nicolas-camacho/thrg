package story

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

func (r *Repository) SaveStory(ctx context.Context, story *Story) error {
	if err := r.db.WithContext(ctx).Save(story).Error; err != nil {
		return fmt.Errorf("error saving story: %w", err)
	}
	return nil
}

func (r *Repository) GetStoryByID(ctx context.Context, storyID uuid.UUID) (*Story, error) {
	var story Story
	if err := r.db.WithContext(ctx).Preload("Acts.Options.Consequences").First(&story, "id = ?", storyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Retorna nil si no se encuentra.
		}
		return nil, fmt.Errorf("error getting story by ID:%w", err)
	}
	return &story, nil
}

func (r *Repository) GetStoryActByID(ctx context.Context, actID uuid.UUID) (*Act, error) {
	var act Act
	if err := r.db.WithContext(ctx).Preload("Options.Consequences").First(&act, "id = ?", actID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting story act by ID:%w", err)
	}
	return &act, nil
}

func (r *Repository) GetAllStories(ctx context.Context) ([]Story, error) {
	var stories []Story
	if err := r.db.WithContext(ctx).Find(&stories).Error; err != nil {
		return nil, fmt.Errorf("error getting all stories:%w", err)
	}
	return stories, nil
}

func (r *Repository) LoadStoriesFromData(ctx context.Context, storiesData []StoryData) error {
	for _, storyData := range storiesData {
		// Crea la historia principal
		story := Story{
			HolderName:          storyData.HolderName,
			Title:               storyData.Title,
			Description:         storyData.Description,
			MisfortuneThreshold: storyData.MisfortuneThreshold,
		}

		// Mapea y enlaza los actos
		var acts []*Act
		for _, actData := range storyData.Acts {
			act := &Act{
				StoryID: story.ID, // Asignado después de guardar la historia
				Order:   actData.Order,
				Text:    actData.Text,
			}
			acts = append(acts, act)

			// Mapea y enlaza las opciones
			var options []*Option
			for _, optionData := range actData.Options {
				option := &Option{
					ActID: act.ID, // Asignado después de guardar el acto
					Text:  optionData.Text,
				}
				options = append(options, option)

				// Mapea y enlaza las consecuencias
				var consequences []Consequence
				for _, consequenceData := range optionData.Consequences {
					consequence := Consequence{
						OptionID: option.ID,
						Type:     ConsequenceType(consequenceData.Type),
						Value:    consequenceData.Value,
					}
					consequences = append(consequences, consequence)
				}
				option.Consequences = consequences
				options = append(options, option)
			}

			act.Options = make([]Option, len(options))
			for i, optPtr := range options {
				act.Options[i] = *optPtr
			}
			acts = append(acts, act)
		}
		story.Acts = make([]Act, len(acts))
		for i, actPtr := range acts {
			story.Acts[i] = *actPtr
		}

		// GORM no puede predecir IDs en cascada sin una operación de guardado,
		// por lo que el mapeo NextAct se hace después.
		if err := r.db.WithContext(ctx).Save(&story).Error; err != nil {
			return fmt.Errorf("error saving story %s: %w", story.Title, err)
		}
	}

	// Segundo bucle para enlazar NextAct ahora que los UUIDs existen
	for _, storyData := range storiesData {
		var story Story
		r.db.WithContext(ctx).First(&story, "holder_name = ?", storyData.HolderName)

		for _, actData := range storyData.Acts {
			var act Act
			r.db.WithContext(ctx).First(&act, "story_id = ? AND \"order\" = ?", story.ID, actData.Order)

			for _, optionData := range actData.Options {
				if optionData.NextActOrder != nil {
					var nextAct Act
					if r.db.WithContext(ctx).First(&nextAct, "story_id = ? AND \"order\" = ?", story.ID, *optionData.NextActOrder).Error == nil {
						var option Option
						r.db.WithContext(ctx).First(&option, "act_id = ? AND text = ?", act.ID, optionData.Text)
						option.NextAct = &nextAct.ID
						r.db.WithContext(ctx).Save(&option)
					}
				}
			}
		}
	}

	return nil
}
