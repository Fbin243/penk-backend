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
	IntrospectToken(ctx context.Context, token string) (*rdb.AuthSession, error)
	CheckPermission(ctx context.Context, profileID, characterID, categoryID *string) error
}

type ICharacterBusiness interface {
	GetCharacterByID(ctx context.Context, id string) (*entity.Character, error)
	GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error)
	UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error)
	DeleteCharacter(ctx context.Context, id string) (*entity.Character, error)
}

type IGoalBusiness interface {
	GetGoals(ctx context.Context, characterID string, status *entity.GoalStatusFilter) ([]entity.Goal, error)
	UpsertGoal(ctx context.Context, input entity.GoalInput) (*entity.Goal, error)
	DeleteGoal(ctx context.Context, id string) (*entity.Goal, error)
}

type ITemplateBusiness interface {
	GetTemplates(ctx context.Context) ([]entity.Template, error)
	GetTemplateCategory(ctx context.Context, id string) (*entity.TemplateTopic, error)
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
	ValidateCharacter(ctx context.Context, profileID, characterID string) error
}

type IGoalRepo interface {
	base.IBaseRepo[entity.Goal]
	GetGoalsByCharacterID(ctx context.Context, characterID string, status *entity.GoalStatusFilter) ([]entity.Goal, error)
	ValidateGoal(ctx context.Context, profileID, goalID string) error
	UpdateStatusOfGoals(ctx context.Context, goalIDs []string, status entity.GoalFinishStatus) error
}

type ITemplateRepo interface {
	base.IBaseRepo[entity.Template]
}

type ITemplateCategoryRepo interface {
	base.IBaseRepo[entity.TemplateTopic]
}

type ICache interface {
	GetAuthSession(ctx context.Context, firebaseUID string) (*rdb.AuthSession, error)
	SetAuthSession(ctx context.Context, profile *entity.Profile, session *rdb.AuthSession) error
	DeleteProfileData(ctx context.Context, profile *entity.Profile) error
}

// RPCs
type ICurrencyClient interface {
	CreateFish(ctx context.Context, profileID string) error
}

type IAnalyticClient interface {
	DeleteCapturedRecords(ctx context.Context, profileID string) error
}
