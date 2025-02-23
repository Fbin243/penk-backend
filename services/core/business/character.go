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

type (
	CategoryMap map[string]MetricMap
	MetricMap   map[string]entity.Metric
)

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
	metricsMap := map[string]entity.Metric{}

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

	if input.Vision != nil {
		character.Vision = entity.Vision{
			Name: input.Vision.Name,
		}

		if input.Vision.Description != nil {
			character.Vision.Description = *input.Vision.Description
		}
	}

	if input.Categories != nil {
		err = biz.upsertCategoriesInCharacter(ctx, character, input.Categories)
		if err != nil {
			return nil, err
		}
	}

	if input.Metrics != nil {
		for _, metric := range character.Metrics {
			metricsMap[metric.ID] = metric
		}

		err = biz.upsertMetricsInCharacter(ctx, upsertMetricsInCharacterInput{
			character:    character,
			metricsMap:   metricsMap,
			metricInputs: input.Metrics,
		})
		if err != nil {
			return nil, err
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

func (biz *CharacterBusiness) upsertCategoriesInCharacter(_ context.Context, character *entity.Character, categoryInputs []entity.CategoryInput) error {
	categories := []entity.Category{}
	categoriesMap := map[string]entity.Category{}
	for _, category := range character.Categories {
		categoriesMap[category.ID] = category
	}

	if len(categoryInputs) > int(utils.LimitedCategoryNumber) {
		return errors.NewGQLError(errors.ErrCodeLimitCategory, nil)
	}

	for _, categoryInput := range categoryInputs {
		if categoryInput.ID == nil {
			// Insert new category
			category := entity.Category{
				ID:   mongodb.GenObjectID(),
				Name: categoryInput.Name,
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

			categories = append(categories, category)
		} else {
			// Update old category
			if _, ok := categoriesMap[*categoryInput.ID]; !ok {
				return errors.ErrBadRequest
			}

			oldCategory := categoriesMap[*categoryInput.ID]
			oldCategory.Name = categoryInput.Name

			if categoryInput.Description != nil {
				oldCategory.Description = *categoryInput.Description
			}

			if categoryInput.Style != nil {
				oldCategory.Style = entity.CategoryStyle{
					Color: categoryInput.Style.Color,
					Icon:  categoryInput.Style.Icon,
				}
			}

			categories = append(categories, oldCategory)
		}
	}

	character.Categories = categories

	return nil
}

type upsertMetricsInCharacterInput struct {
	character    *entity.Character
	metricsMap   map[string]entity.Metric
	metricInputs []entity.MetricInput
}

func (biz *CharacterBusiness) upsertMetricsInCharacter(_ context.Context, input upsertMetricsInCharacterInput) error {
	metrics := []entity.Metric{}
	if len(input.metricInputs) > int(utils.LimitedMetricNumber) {
		return errors.NewGQLError(errors.ErrCodeLimitMetric, nil)
	}

	for _, metricInput := range input.metricInputs {
		if metricInput.ID == nil {
			// Insert new metric
			metric := entity.Metric{
				ID:         mongodb.GenObjectID(),
				CategoryID: metricInput.CategoryID,
				Name:       metricInput.Name,
				Value:      metricInput.Value,
				Unit:       metricInput.Unit,
			}

			metrics = append(metrics, metric)
		} else {
			// Update old metric
			if _, ok := input.metricsMap[*metricInput.ID]; !ok {
				return errors.ErrBadRequest
			}

			oldMetric := input.metricsMap[*metricInput.ID]
			oldMetric.CategoryID = metricInput.CategoryID
			oldMetric.Name = metricInput.Name
			oldMetric.Unit = metricInput.Unit
			oldMetric.Value = metricInput.Value

			metrics = append(metrics, oldMetric)
		}
	}

	input.character.Metrics = metrics

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
