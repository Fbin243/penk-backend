package business

import (
	"context"

	"tenkhours/services/core/entity"
)

type PermissionEntity struct {
	ID   string
	Type entity.EntityType
}

func (biz *PermissionBusiness) CheckOwnCharacter(ctx context.Context, profileID, characterID string) error {
	err := biz.CharacterRepo.Exist(ctx, profileID, characterID)
	if err != nil {
		return err
	}
	return nil
}

func (biz *PermissionBusiness) CheckOwnEntities(ctx context.Context, characterID string, entities []PermissionEntity) error {
	for _, entityItem := range entities {
		switch entityItem.Type {
		case entity.EntityTypeCategory:
			err := biz.CategoryRepo.Exist(ctx, characterID, entityItem.ID)
			if err != nil {
				return err
			}
		case entity.EntityTypeMetric:
			err := biz.MetricRepo.Exist(ctx, characterID, entityItem.ID)
			if err != nil {
				return err
			}
		case entity.EntityTypeGoal:
			err := biz.GoalRepo.Exist(ctx, characterID, entityItem.ID)
			if err != nil {
				return err
			}
		case entity.EntityTypeHabit:
			err := biz.HabitRepo.Exist(ctx, characterID, entityItem.ID)
			if err != nil {
				return err
			}
		case entity.EntityTypeTask:
			err := biz.TaskRepo.Exist(ctx, characterID, entityItem.ID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
