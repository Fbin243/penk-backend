package business

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"github.com/teambition/rrule-go"
)

type CreateReminderInput struct {
	CharacterID   string
	Name          string
	RemindTimeStr string
	RRule         string
	ReferenceID   *string
	ReferenceType *entity.EntityType
}

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
		BaseEntity: &base.BaseEntity{
			ID: mongodb.GenObjectID(),
		},
		CharacterID: authSession.CurrentCharacterID,
	}

	r, err := rrule.StrToRRule(input.RRule)
	if err != nil {
		return nil, err
	}

	if input.ID == nil {
		count, err := b.reminderRepo.CountByCharacterID(ctx, authSession.CurrentCharacterID)
		if err != nil {
			return nil, err
		}
		if count >= utils.LimitedReminderNumber {
			return nil, errors.ErrLimitReminder
		}

		remindTime, err := calculateRemindTime(r.GetDTStart(), input.RemindTimeStr)
		if err != nil {
			return nil, err
		}

		reminder.RemindTime = remindTime
		if utils.IsToday(*remindTime) {
			err = b.reminderCache.UpsertReminder(ctx, reminder)
			if err != nil {
				return nil, err
			}
		}
	} else {
		reminder, err = b.reminderRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}

		if input.RemindTimeStr != reminder.RemindTimeStr {
			if b.reminderCache.Exist(ctx, reminder) == nil {
				reminder.RemindTime, err = calculateRemindTime(r.GetDTStart(), input.RemindTimeStr)
				if err != nil {
					return nil, err
				}

				// This reminder is already in cache, so we need to update it with new remindTime
				err = b.reminderCache.UpsertReminder(ctx, reminder)
				if err != nil {
					return nil, err
				}
			}
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

	err = b.reminderCache.DeleteReminder(ctx, reminderID)
	if err != nil {
		return nil, err
	}

	return b.reminderRepo.FindAndDeleteByID(ctx, reminderID)
}

func calculateRemindTime(date time.Time, remindTimeStr string) (*time.Time, error) {
	// Parse time from RemindTimeStr (format: hh:mm)
	timeParts := strings.Split(remindTimeStr, ":")
	if len(timeParts) != 2 {
		return nil, fmt.Errorf("invalid time format, expected hh:mm")
	}
	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid hour: %v", err)
	}
	minute, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minute: %v", err)
	}

	return lo.ToPtr(time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, time.UTC)), nil
}
