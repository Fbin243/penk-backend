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
		return nil, errors.Unauthorized()
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
		return nil, errors.Unauthorized()
	}

	profile, err := biz.ProfileRepo.FindByID(ctx, authSession.ProfileID)
	if err != nil {
		return nil, err
	}

	if input.ID == nil {
		// Insert new character
		charactersCount, err := biz.CharacterRepo.CountCharactersByProfileID(ctx, authSession.ProfileID)
		if err != nil {
			return nil, err
		}

		if charactersCount >= int64(utils.LimitedCharacterNumber) {
			return nil, errors.NewGQLError(errors.ErrCodeLimitCharacter, nil)
		}

		character := &entity.Character{
			BaseEntity: &base.BaseEntity{},
			Name:       input.Name,
			Gender:     input.Gender,
			ProfileID:  profile.ID,
			Categories: []entity.Category{},
			Metrics:    []entity.Metric{},
		}

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
			err = biz.upsertMetricsInCharacter(ctx, character, input.Metrics)
			if err != nil {
				return nil, err
			}
		}

		createdCharacter, err := biz.CharacterRepo.InsertOne(ctx, character)
		if err != nil {
			return nil, err
		}

		// Update the current character of the profile
		profile.CurrentCharacterID = character.ID
		_, err = biz.ProfileRepo.UpdateByID(ctx, profile.ID, profile)
		if err != nil {
			return nil, err
		}

		return createdCharacter, nil
	}

	// Update existing character
	character, err := biz.CharacterRepo.FindByID(ctx, *input.ID)
	if err != nil {
		return nil, err
	}

	if ok, _ := auth.CheckPermission(profile, character, "write"); !ok {
		return nil, errors.PermissionDenied()
	}

	character.Name = input.Name

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
		err = biz.upsertMetricsInCharacter(ctx, character, input.Metrics)
		if err != nil {
			return nil, err
		}
	}

	return biz.CharacterRepo.UpdateByID(ctx, *input.ID, character)
}

func (biz *CharacterBusiness) upsertCategoriesInCharacter(ctx context.Context, character *entity.Character, categoryInputs []entity.CategoryInput) error {
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
				return errors.BadRequest()
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

func (biz *CharacterBusiness) upsertMetricsInCharacter(ctx context.Context, character *entity.Character, metricInputs []entity.MetricInput) error {
	metrics := make([]entity.Metric, 0)
	metricsMap := make(map[string]entity.Metric)
	for _, metric := range character.Metrics {
		metricsMap[metric.ID] = metric
	}

	if len(metricInputs) > int(utils.LimitedMetricNumber) {
		return errors.NewGQLError(errors.ErrCodeLimitMetric, nil)
	}

	for _, metricInput := range metricInputs {
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
			if _, ok := metricsMap[*metricInput.ID]; !ok {
				return errors.BadRequest()
			}

			oldMetric := metricsMap[*metricInput.ID]
			oldMetric.CategoryID = metricInput.CategoryID
			oldMetric.Name = metricInput.Name
			oldMetric.Unit = metricInput.Unit
			oldMetric.Value = metricInput.Value

			metrics = append(metrics, oldMetric)
		}
	}

	character.Metrics = metrics

	return nil
}

func (biz *CharacterBusiness) DeleteCharacter(ctx context.Context, id string) (*entity.Character, error) {
	authSession, ok := ctx.Value(auth.AuthSessionKey).(rdb.AuthSession)
	if !ok {
		return nil, errors.Unauthorized()
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
		return nil, errors.PermissionDenied()
	}

	deletedCharacter, err := biz.CharacterRepo.DeleteCharacter(ctx, id)
	if err != nil {
		return nil, err
	}

	return deletedCharacter, nil
}

// TODO: Fbin do goals later
// func (biz *CharacterBusiness) checkGoalsFinished(ctx context.Context, goals []entity.Goal, newCategories, oldCategories []entity.Category) error {
// 	newCategoriesMap := CategoryMap{}
// 	for _, category := range newCategories {
// 		newCategoriesMap[category.ID] = MetricMap{}
// 		for _, metric := range category.Metrics {
// 			newCategoriesMap[category.ID][metric.ID] = metric
// 		}
// 	}

// 	oldCategoriesMap := CategoryMap{}
// 	for _, category := range oldCategories {
// 		oldCategoriesMap[category.ID] = MetricMap{}
// 		for _, metric := range category.Metrics {
// 			oldCategoriesMap[category.ID][metric.ID] = metric
// 		}
// 	}

// 	finishedGoalIDs := make([]string, 0)
// 	// TODO: Migrate logic check goals
// 	for _, goal := range goals {
// 		for _, category := range goal.Target {
// 			newCategory, ok := newCategoriesMap[category.ID]
// 			if !ok {
// 				return errors.PermissionDenied()
// 			}

// 			oldCategory, ok := oldCategoriesMap[category.ID]
// 			if !ok {
// 				return errors.PermissionDenied()
// 			}

// 			for _, metric := range category.Metrics {
// 				newMetric, ok := newCategory[metric.ID]
// 				if !ok {
// 					return errors.PermissionDenied()
// 				}

// 				oldMetric, ok := oldCategory[metric.ID]
// 				if !ok {
// 					return errors.PermissionDenied()
// 				}

// 				if (newMetric.Value <= metric.Value && metric.Value <= oldMetric.Value) ||
// 					(oldMetric.Value <= metric.Value && metric.Value <= newMetric.Value) {
// 					goal.Status = entity.GoalFinishStatusFinished
// 					finishedGoalIDs = append(finishedGoalIDs, goal.ID)
// 				}
// 			}
// 		}
// 	}

// 	// Update the status of the finished goals
// 	if len(finishedGoalIDs) > 0 {
// 		_, err := biz.GoalRepo.UpdateStatusOfGoals(ctx, finishedGoalIDs, entity.GoalFinishStatusFinished)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
