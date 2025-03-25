package business

import (
	"context"

	"tenkhours/pkg/db/base"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"
)

// Business
type IProfileBusiness interface {
	GetProfile(ctx context.Context) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, input entity.ProfileInput) (*entity.Profile, error)
	DeleteProfile(ctx context.Context) (*entity.Profile, error)
	IntrospectToken(ctx context.Context, token, deviceID string) (*rdb.AuthSession, error)
	CheckPermission(ctx context.Context, profileID, characterID, categoryID *string) error
}

type ICharacterBusiness interface {
	GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error)
	UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error)
	DeleteCharacter(ctx context.Context, id string) (*entity.Character, error)
}

type IGoalBusiness interface {
	GetGoals(ctx context.Context, characterID string, status *entity.GoalStatus) ([]entity.Goal, error)
	UpsertGoal(ctx context.Context, input entity.GoalInput) (*entity.Goal, error)
	DeleteGoal(ctx context.Context, id string) (*entity.Goal, error)
}

type IMetricBusiness interface {
	GetMetrics(ctx context.Context, characterID string) ([]entity.Metric, error)
	UpsertMetric(ctx context.Context, input entity.MetricInput) (*entity.Metric, error)
	DeleteMetric(ctx context.Context, id string) (*entity.Metric, error)
}

type ICategoryBusiness interface {
	GetCategories(ctx context.Context, characterID string) ([]entity.Category, error)
	UpsertCategory(ctx context.Context, input entity.CategoryInput) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id string) (*entity.Category, error)
}

// Repository
type IProfileRepo interface {
	base.IBaseRepo[entity.Profile]
	ProfileExists(ctx context.Context, firebaseUID string) (bool, error)
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
	Exist(ctx context.Context, characterID, categoryID string) error
}

type ICategoryRepo interface {
	base.IBaseRepo[entity.Category]
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, categoryID string) error
	FindByCharacterID(ctx context.Context, characterID string) ([]entity.Category, error)
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
}

type IMetricRepo interface {
	base.IBaseRepo[entity.Metric]
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	CountByCategoryID(ctx context.Context, categoryID string) (int, error)
	CountUnassigned(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, categoryID string) error
	FindByCharacterID(ctx context.Context, characterID string) ([]entity.Metric, error)
	UnassignCategory(ctx context.Context, categoryID string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
}

type IGoalRepo interface {
	base.IBaseRepo[entity.Goal]
	GetGoalsByCharacterID(ctx context.Context, characterID string, status *entity.GoalStatus) ([]entity.Goal, error)
	ValidateGoal(ctx context.Context, profileID, goalID string) error
	UpdateStatusOfGoals(ctx context.Context, goalIDs []string, status entity.GoalStatus) error
	SyncGoalStatus(ctx context.Context, characterID string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
}

type ICache interface {
	GetAuthSession(ctx context.Context, firebaseUID string) (*rdb.AuthSession, error)
	SetAuthSession(ctx context.Context, profile *entity.Profile, session *rdb.AuthSession) error
	DeleteProfileData(ctx context.Context, profile *entity.Profile) error
}

// RPCs
type ICurrencyClient interface {
	CreateFish(ctx context.Context, profileID string) error
	DeleteFish(ctx context.Context, profileID string) error
}

// TODO: Allow Core service fetch data from Timetracking repo
type ITimeTrackingRepo interface {
	GetTotalTimeByCategoryID(ctx context.Context, categoryID string) (int, error)
	GetTotalTimeOfUnassigned(ctx context.Context, characterID string) (int, error)
	GetTotalTimeByCharacterID(ctx context.Context, characterID string) (int, error)
	UnassignCategory(ctx context.Context, categoryID string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
}
