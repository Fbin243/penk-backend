package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

func (b *ReminderBusiness) Upsert(ctx context.Context, input *entity.ReminderInput) (*entity.Reminder, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	permEntities := []PermissionEntity{}
	if input.ID != nil {
		permEntities = append(permEntities, PermissionEntity{
			ID:   *input.ID,
			Type: entity.EntityTypeReminder,
		})
	}

	if input.ReferenceID != nil && input.ReferenceType != nil {
		permEntities = append(permEntities, PermissionEntity{
			ID:   *input.ReferenceID,
			Type: *input.ReferenceType,
		})
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, permEntities)
	if err != nil {
		return nil, err
	}

	reminder := &entity.Reminder{
		BaseEntity:  &base.BaseEntity{},
		CharacterID: authSession.CurrentCharacterID,
	}

	if input.ID == nil {
		count, err := b.reminderRepo.CountByCharacterID(ctx, authSession.CurrentCharacterID)
		if err != nil {
			return nil, err
		}
		if count >= utils.LimitedReminderNumber {
			return nil, errors.ErrLimitCategory
		}
	} else {
		reminder, err = b.reminderRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}
	}

	err = copier.Copy(reminder, input)
	if err != nil {
		return nil, err
	}

	if input.ID != nil {
		return b.reminderRepo.FindAndUpdateByID(ctx, *input.ID, reminder)
	}

	return b.reminderRepo.InsertOne(ctx, reminder)
}

func (b *ReminderBusiness) Delete(ctx context.Context, reminderID string) (*entity.Reminder, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	err = b.permBiz.CheckOwnEntities(ctx, authSession.CurrentCharacterID, []PermissionEntity{
		{
			ID:   reminderID,
			Type: entity.EntityTypeReminder,
		},
	})
	if err != nil {
		return nil, err
	}

	return b.reminderRepo.FindAndDeleteByID(ctx, reminderID)
}
