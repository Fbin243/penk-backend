package business

import (
	"context"
	"time"

	"tenkhours/pkg/db/base"
	"tenkhours/services/core/entity"

	rdb "tenkhours/pkg/db/redis"
)

// Business
type IPermissionBusiness interface {
	CheckOwnCharacter(ctx context.Context, profileID, characterID string) error
	CheckOwnEntities(ctx context.Context, characterID string, entities []PermissionEntity) error
}

type IProfileBusiness interface {
	GetProfile(ctx context.Context) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, input entity.ProfileInput) (*entity.Profile, error)
	DeleteProfile(ctx context.Context) (*entity.Profile, error)
	IntrospectUser(ctx context.Context, token, userID, deviceID string) (*rdb.AuthSession, error)
}

type ICharacterBusiness interface {
	GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error)
	UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error)
	DeleteCharacter(ctx context.Context, id string) (*entity.Character, error)
}

type IGoalBusiness interface {
	base.IBaseBusiness[entity.Goal, entity.GoalInput, entity.GoalFilter, entity.GoalOrderBy]
}

type IMetricBusiness interface {
	base.IBaseBusiness[entity.Metric, entity.MetricInput, entity.MetricFilter, entity.MetricOrderBy]
}

type ICategoryBusiness interface {
	base.IBaseBusiness[entity.Category, entity.CategoryInput, entity.CategoryFilter, entity.CategoryOrderBy]
}

type IHabitBusiness interface {
	base.IBaseBusiness[entity.Habit, entity.HabitInput, entity.HabitFilter, entity.HabitOrderBy]
	CountHabitLog(ctx context.Context, filter *entity.HabitLogFilter) (int, error)
	GetHabitLogs(ctx context.Context, filter *entity.HabitLogFilter, orderBy *entity.HabitLogOrderBy, limit, offset *int) ([]entity.HabitLog, error)
	UpsertHabitLog(ctx context.Context, input *entity.HabitLogInput) (*entity.HabitLog, error)
}

type ITimeTrackingBusiness interface {
	UpsertTimeTracking(ctx context.Context, input *entity.TimeTrackingInput) (*entity.TimeTracking, error)
}

type ITaskBusiness interface {
	base.IBaseBusiness[entity.Task, entity.TaskInput, entity.TaskFilter, entity.TaskOrderBy]
	CountTaskSession(ctx context.Context, filter *entity.TaskSessionFilter) (int, error)
	GetTaskSessions(ctx context.Context, filter *entity.TaskSessionFilter, orderBy *entity.TaskSessionOrderBy, limit, offset *int) ([]entity.TaskSession, error)
	UpsertTaskSession(ctx context.Context, input *entity.TaskSessionInput) (*entity.TaskSession, error)
	UpsertTaskSessions(ctx context.Context, inputs []entity.TaskSessionInput) ([]entity.TaskSession, error)
	DeleteTaskSession(ctx context.Context, id string) (*entity.TaskSession, error)
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
	CountCharactersByProfileID(ctx context.Context, profileID string) (int, error)
	GetAllCharacters(ctx context.Context) ([]entity.Character, error)
	DeleteCharacter(ctx context.Context, id string) (*entity.Character, error)
	DeleteCharactersByProfileID(ctx context.Context, profileID string) error
	Exist(ctx context.Context, profileID, characterID string) error
}

type ICategoryRepo interface {
	base.IBaseRepo[entity.Category]
	Find(ctx context.Context, p entity.CategoryPipeline) ([]entity.Category, error)
	CountByFilter(ctx context.Context, filter *entity.CategoryFilter) (int, error)
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, categoryID string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
}

type IMetricRepo interface {
	base.IBaseRepo[entity.Metric]
	Find(ctx context.Context, p entity.MetricPipeline) ([]entity.Metric, error)
	CountByFilter(ctx context.Context, filter *entity.MetricFilter) (int, error)
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	CountByCategoryID(ctx context.Context, categoryID string) (int, error)
	CountUnassigned(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, metricID string) error
	UnassignCategory(ctx context.Context, categoryID string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
}

type IGoalRepo interface {
	base.IBaseRepo[entity.Goal]
	Find(ctx context.Context, p entity.GoalPipeline) ([]entity.Goal, error)
	CountByFilter(ctx context.Context, filter *entity.GoalFilter) (int, error)
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
	Exist(ctx context.Context, characterID, goalID string) error
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
}

type ICache interface {
	GetAuthSession(ctx context.Context, firebaseUID string) (*rdb.AuthSession, error)
	SetAuthSession(ctx context.Context, profile *entity.Profile, session *rdb.AuthSession) error
	DeleteAuthSession(ctx context.Context, firebaseUID string) error
	DeleteProfileData(ctx context.Context, profile *entity.Profile) error
}

type IHabitRepo interface {
	base.IBaseRepo[entity.Habit]
	Find(ctx context.Context, pipeline entity.HabitPipeline) ([]entity.Habit, error)
	CountByFilter(ctx context.Context, filter *entity.HabitFilter) (int, error)
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	CountByCategoryID(ctx context.Context, categoryID string) (int, error)
	CountUnassigned(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, habitID string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
	UnassignCategory(ctx context.Context, categoryID string) error
}

type IHabitLogRepo interface {
	base.IBaseRepo[entity.HabitLog]
	Find(ctx context.Context, pipeline entity.HabitLogPipeline) ([]entity.HabitLog, error)
	CountByFilter(ctx context.Context, filter *entity.HabitLogFilter) (int, error)
	UpsertByTimestamp(ctx context.Context, timestamp time.Time, habit *entity.HabitLog) error
	DeleteByHabitID(ctx context.Context, habitID string) error
	DeleteByHabitIDs(ctx context.Context, habitIDs []string) error
	CountByHabitID(ctx context.Context, habitID string) (int, error)
}

type ITimeTrackingRepo interface {
	base.IBaseRepo[entity.TimeTracking]
	FindByCategoryID(ctx context.Context, categoryID string) ([]entity.TimeTracking, error)
	FindByCharacterID(ctx context.Context, characterID string) ([]entity.TimeTracking, error)
	GetTotalTimeByCategoryID(ctx context.Context, categoryID string) (int, error)
	GetTotalTimeOfUnassigned(ctx context.Context, characterID string) (int, error)
	GetTotalTimeByCharacterID(ctx context.Context, characterID string) (int, error)
	UnassignCategory(ctx context.Context, categoryID string) error
	UnassignReference(ctx context.Context, referenceID string, referenceType entity.EntityType) error
	UpdateCategoryByReferenceID(ctx context.Context, referenceID string, categoryID *string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
	FindByReferenceIDAndTimestamp(ctx context.Context, refID string, timestamp time.Time) (*entity.TimeTracking, error)
	FindByReferenceID(ctx context.Context, referenceID string) ([]entity.TimeTracking, error)
	DeleteByIDs(ctx context.Context, ids []string) error
}

type ITaskRepo interface {
	base.IBaseRepo[entity.Task]
	Find(ctx context.Context, pipeline entity.TaskPipeline) ([]entity.Task, error)
	CountByFilter(ctx context.Context, filter *entity.TaskFilter) (int, error)
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	CountByCategoryID(ctx context.Context, categoryID string) (int, error)
	CountUnassigned(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, taskID string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
	UnassignCategory(ctx context.Context, categoryID string) error
}

type ITaskSessionRepo interface {
	base.IBaseRepo[entity.TaskSession]
	Find(ctx context.Context, pipeline entity.TaskSessionPipeline) ([]entity.TaskSession, error)
	CountByFilter(ctx context.Context, filter *entity.TaskSessionFilter) (int, error)
	DeleteByTaskID(ctx context.Context, taskID string) error
	DeleteByTaskIDs(ctx context.Context, taskIDs []string) error
	CountByTaskID(ctx context.Context, taskID string) (int, error)
}

// RPCs
type ICurrencyClient interface {
	CreateFish(ctx context.Context, profileID string) error
	DeleteFish(ctx context.Context, profileID string) error
}

type INotificationClient interface {
	SendNotification(ctx context.Context, req *entity.SendNotiReq) (bool, error)
}

type IRewardRepo interface {
	DeleteReward(ctx context.Context, profileID string) error
}
