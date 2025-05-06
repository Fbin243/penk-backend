package business

import (
	"context"
	"fmt"
	"tenkhours/pkg/auth"
	"tenkhours/services/notification/entity"

	"tenkhours/pkg/db/base"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/pkg/errors"
)

type ReminderBusiness struct {
	ReminderRepo IReminderRepo
}

func NewReminderBusiness(reminderRepo IReminderRepo) *ReminderBusiness {
	return &ReminderBusiness{
		ReminderRepo: reminderRepo,
	}
}

func (biz *ReminderBusiness) CreateReminder(ctx context.Context, reminder entity.ReminderInput) (*entity.Reminder, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	newReminder := &entity.Reminder{
		BaseEntity:   &base.BaseEntity{},
		ProfileID:    authSession.ProfileID,
		Title:        reminder.Title,
		Type:         reminder.Type,
		RemindTime:   reminder.RemindTime,
		Recurrence:   reminder.Recurrence,
		LinkedItemID: reminder.LinkedItemID,
	}

	createdReminder, err := biz.ReminderRepo.CreateReminder(ctx, newReminder)
	if err != nil {
		return nil, fmt.Errorf("failed to create reminder: %w", err)
	}

	return createdReminder, nil
}

func (biz *ReminderBusiness) GetRemindersByProfileID(ctx context.Context) ([]*entity.Reminder, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	reminders, err := biz.ReminderRepo.GetRemindersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminders - err: %v", err)
	}

	return reminders, nil
}

func (biz *ReminderBusiness) GetReminderByID(ctx context.Context, reminderID string) (*entity.Reminder, error) {
	reminder, err := biz.ReminderRepo.GetReminderByID(ctx, reminderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder - err: %v", err)
	}

	return reminder, nil
}

func (biz *ReminderBusiness) UpdateReminder(ctx context.Context, input entity.ReminderInput) (*entity.Reminder, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	reminder := &entity.Reminder{
		ProfileID:    authSession.ProfileID,
		Title:        input.Title,
		Type:         input.Type,
		RemindTime:   input.RemindTime,
		Recurrence:   input.Recurrence,
		LinkedItemID: input.LinkedItemID,
	}

	updatedReminder, err := biz.ReminderRepo.UpdateReminder(ctx, reminder)
	if err != nil {
		return nil, fmt.Errorf("failed to update reminder: %w", err)
	}

	return updatedReminder, nil
}

func (biz *ReminderBusiness) DeleteReminder(ctx context.Context, reminderID string) (bool, error) {
	deleted, err := biz.ReminderRepo.DeleteReminder(ctx, reminderID)
	if err != nil {
		return false, fmt.Errorf("failed to delete reminder - err: %v", err)
	}

	return deleted, nil
}
