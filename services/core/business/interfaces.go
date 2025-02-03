package business

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/base"
	"tenkhours/services/core/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

// Business
type IProfileBusiness interface {
	GetProfile(ctx context.Context) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, input entity.ProfileInput) (*entity.Profile, error)
	DeleteProfile(ctx context.Context) (*entity.Profile, error)
	IntrospectProfile(ctx context.Context, firebaseProfile auth.FirebaseProfile) (*entity.Profile, error)
	CheckPermission(ctx context.Context, profileID, characterID string, metricID *string) error
	BuyItem(ctx context.Context, profileID, characterID, metricID *string, item entity.ItemType, amount int32) error
}

type ICharacterBusiness interface {
	GetCharacterByID(ctx context.Context, id string) (*entity.Character, error)
	GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error)
	UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error)
	DeleteCharacter(ctx context.Context, id string) (*entity.Character, error)
	UpdateTimeInCharacter(ctx context.Context, characterID string, metricID *string, time int32) error
}

type IGoalBusiness interface {
	GetGoals(ctx context.Context, characterID string, status *entity.GoalStatusFilter) ([]entity.Goal, error)
	UpsertGoal(ctx context.Context, characterID string, input entity.GoalInput) (*entity.Goal, error)
}

type ITemplateBusiness interface {
	GetTemplates(ctx context.Context) ([]entity.Template, error)
	GetTemplateCategory(ctx context.Context, id string) (*entity.TemplateCategory, error)
}

// Repository
type IProfileRepo interface {
	base.IBaseRepo[entity.Profile]
	GetProfileByFirebaseUID(ctx context.Context, firebaseUID string) (*entity.Profile, error)
	DeleteProfileByFirebaseUID(ctx context.Context, firebaseUID string) error
}

type ICharacterRepo interface {
	base.IBaseRepo[entity.Character]
	GetCharactersByProfileID(ctx context.Context, profileID string) ([]entity.Character, error)
	CountCharactersByProfileID(ctx context.Context, profileID string) (int64, error)
	GetAllCharacters(ctx context.Context) ([]entity.Character, error)
	DeleteCharacter(ctx context.Context, id string) (*entity.Character, error)
	DeleteCharactersByProfileID(ctx context.Context, profileID string) error
}

type IGoalRepo interface {
	base.IBaseRepo[entity.Goal]
	GetGoalsByCharacterID(ctx context.Context, characterID string, status *entity.GoalStatusFilter) ([]entity.Goal, error)
	// TODO: @Fbin243 refactor the return type
	UpdateOneMetricInGoals(ctx context.Context, metric entity.CustomMetric, goalIDs []string) (*mongo.UpdateResult, error)
	RemoveOnePropertyFromGoals(ctx context.Context, metricID, propertyID string, goalIDs []string) (*mongo.UpdateResult, error)
	UpdateStatusOfGoals(ctx context.Context, goalIDs []string, status entity.GoalFinishStatus) (*mongo.UpdateResult, error)
}

type ITemplateRepo interface {
	base.IBaseRepo[entity.Template]
}

type ITemplateCategoryRepo interface {
	base.IBaseRepo[entity.TemplateCategory]
}

type ICache interface {
	DeleteProfileData(ctx context.Context, profile *entity.Profile) error
}

// RPCs
type ICurrencyClient interface {
	CreateFish(ctx context.Context, profileID string) error
}

type IAnalyticClient interface {
	DeleteCapturedRecords(ctx context.Context, profileID string) error
}
