package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"

	"github.com/samber/lo"
)

type CharacterBusiness struct {
	CharacterRepo ICharacterRepo
	ProfileRepo   IProfileRepo
	GoalRepo      IGoalRepo
}

func NewCharacterBusiness(characterRepo ICharacterRepo, profileRepo IProfileRepo, goalRepo IGoalRepo) *CharacterBusiness {
	return &CharacterBusiness{
		CharacterRepo: characterRepo,
		ProfileRepo:   profileRepo,
		GoalRepo:      goalRepo,
	}
}

func (biz *CharacterBusiness) GetCharacterByID(ctx context.Context, id string) (*entity.Character, error) {
	character, err := biz.CharacterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (biz *CharacterBusiness) GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	characters, err := biz.CharacterRepo.GetCharactersByProfileID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (biz *CharacterBusiness) UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	profile, err := biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	var character *entity.Character
	if input.ID == nil {
		// Insert new character
		charactersCount, err := biz.CharacterRepo.CountCharactersByProfileID(ctx, authSession.ProfileID)
		if err != nil {
			return nil, err
		}

		if charactersCount >= int64(utils.LimitedCharacterNumber) {
			return nil, errors.NewGQLError(errors.ErrCodeLimitCharacter, nil)
		}

		character = &entity.Character{
			BaseEntity: &base.BaseEntity{},
			ProfileID:  profile.ID,
			Categories: []entity.Category{},
			Metrics:    []entity.Metric{},
		}
	} else {
		// Update existing character
		character, err = biz.CharacterRepo.FindByID(ctx, *input.ID)
		if err != nil {
			return nil, err
		}

		if ok, _ := auth.CheckPermission(profile, character, "write"); !ok {
			return nil, errors.ErrPermissionDenied
		}
	}

	character.Name = input.Name
	character.Gender = input.Gender
	character.Tags = input.Tags
	if input.Tags != nil {
		character.Tags = input.Tags
	}

	if input.Categories != nil {
		err = biz.upsertCategoriesInCharacter(ctx, character, input.Categories)
		if err != nil {
			return nil, err
		}
	}

	if input.Metrics != nil {
		err = biz.upsertMetrics(ctx, upsertMetricsInput{
			character:    character,
			metricInputs: input.Metrics,
		})
		if err != nil {
			return nil, err
		}
	}

	// Update the character's goals
	if input.Metrics != nil ||
		input.Categories != nil {
		// Create a map of metrics for easy access
		metricsMap := map[string]entity.Metric{}
		for _, category := range character.Categories {
			for _, metric := range category.Metrics {
				metricsMap[metric.ID] = metric
			}
		}
		for _, metric := range character.Metrics {
			metricsMap[metric.ID] = metric
		}

		// Find all unfinished and unexpired goals of the character
		goals, err := biz.GoalRepo.GetGoalsByCharacterID(ctx, character.ID, &entity.GoalStatusFilter{
			FinishStatus: lo.ToPtr(entity.GoalFinishStatusUnfinished),
			ExpireStatus: lo.ToPtr(entity.GoalExpireStatusUnexpired),
		})
		if err != nil {
			return nil, err
		}

		finishedGoalIDs := lo.Map(lo.Filter(goals, func(goal entity.Goal, _ int) bool {
			goal.UpdateStatus(metricsMap)
			return goal.Status == entity.GoalFinishStatusFinished
		}), func(goal entity.Goal, _ int) string {
			return goal.ID
		})

		err = biz.GoalRepo.UpdateStatusOfGoals(ctx, finishedGoalIDs, entity.GoalFinishStatusFinished)
		if err != nil {
			return nil, err
		}
	}

	if input.ID == nil {
		character, err = biz.CharacterRepo.InsertOne(ctx, character)
		if err != nil {
			return nil, err
		}

		// Update the current character of the profile with the new character
		profile.CurrentCharacterID = lo.ToPtr(character.ID)
		_, err = biz.ProfileRepo.UpdateByID(ctx, profile.ID, profile)
		if err != nil {
			return nil, err
		}
	} else {
		character, err = biz.CharacterRepo.UpdateByID(ctx, *input.ID, character)
		if err != nil {
			return nil, err
		}
	}

	return character, nil
}

func (biz *CharacterBusiness) upsertCategoriesInCharacter(ctx context.Context, character *entity.Character, categoryInputs []entity.CategoryInput) error {
	if len(categoryInputs) > int(utils.LimitedCategoryNumber) {
		return errors.NewGQLError(errors.ErrCodeLimitCategory, nil)
	}

	categories := []entity.Category{}
	for _, categoryInput := range categoryInputs {
		category := entity.Category{
			Name: categoryInput.Name,
		}

		if categoryInput.ID == nil {
			category.ID = mongodb.GenObjectID()
		} else {
			category.ID = *categoryInput.ID
		}

		if categoryInput.Description != nil {
			category.Description = *categoryInput.Description
		}

		if categoryInput.Style != nil {
			category.Style = entity.CategoryStyle{
				Color: categoryInput.Style.Color,
				Icon:  categoryInput.Style.Icon,
			}
		}

		if categoryInput.Metrics != nil {
			err := biz.upsertMetrics(ctx, upsertMetricsInput{
				category:     &category,
				metricInputs: categoryInput.Metrics,
			})
			if err != nil {
				return err
			}
		}

		categories = append(categories, category)
	}

	character.Categories = categories

	return nil
}

type upsertMetricsInput struct {
	character    *entity.Character
	category     *entity.Category
	metricInputs []entity.MetricInput
}

func (biz *CharacterBusiness) upsertMetrics(_ context.Context, input upsertMetricsInput) error {
	if len(input.metricInputs) > int(utils.LimitedMetricNumber) {
		return errors.NewGQLError(errors.ErrCodeLimitMetric, nil)
	}

	metrics := []entity.Metric{}
	for _, metricInput := range input.metricInputs {
		metric := entity.Metric{
			Name:  metricInput.Name,
			Value: metricInput.Value,
			Unit:  metricInput.Unit,
		}

		if metricInput.ID == nil {
			metric.ID = mongodb.GenObjectID()
		} else {
			metric.ID = *metricInput.ID
		}

		metrics = append(metrics, metric)
	}

	if input.character != nil {
		input.character.Metrics = metrics
	}
	if input.category != nil {
		input.category.Metrics = metrics
	}

	return nil
}

func (biz *CharacterBusiness) DeleteCharacter(ctx context.Context, id string) (*entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.ErrUnauthorized
	}

	character, err := biz.CharacterRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	profile := &entity.Profile{
		BaseEntity: &base.BaseEntity{
			ID: authSession.ProfileID,
		},
	}

	if ok, _ := auth.CheckPermission(profile, character, "write"); !ok {
		return nil, errors.ErrPermissionDenied
	}

	deletedCharacter, err := biz.CharacterRepo.DeleteCharacter(ctx, id)
	if err != nil {
		return nil, err
	}

	return deletedCharacter, nil
}
