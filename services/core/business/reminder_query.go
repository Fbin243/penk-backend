package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/types"
	"tenkhours/services/core/entity"
)

func (b *ReminderBusiness) Get(ctx context.Context, filter *entity.ReminderFilter, orderBy *entity.ReminderOrderBy, limit, offset *int) ([]entity.Reminder, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = &entity.ReminderFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	return b.reminderRepo.Find(ctx, entity.ReminderPipeline{
		Filter:  filter,
		OrderBy: orderBy,
		Pagination: &types.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	})
}

func (b *ReminderBusiness) Count(ctx context.Context, filter *entity.ReminderFilter) (int, error) {
	authSession, err := auth.GetAuthSession(ctx)
	if err != nil {
		return 0, err
	}

	if filter == nil {
		filter = &entity.ReminderFilter{}
	}
	filter.CharacterID = &authSession.CurrentCharacterID

	return b.reminderRepo.CountByFilter(ctx, filter)
}
