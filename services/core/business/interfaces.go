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
	IntrospectToken(ctx context.Context, token, deviceID string) (*rdb.AuthSession, error)
}

type ICharacterBusiness interface {
	GetCharactersByProfileID(ctx context.Context) ([]entity.Character, error)
	UpsertCharacter(ctx context.Context, input entity.CharacterInput) (*entity.Character, error)
	DeleteCharacter(ctx context.Context, id string) (*entity.Character, error)
}

type IGoalBusiness interface {
	GetGoals(ctx context.Context, status *entity.GoalStatus) ([]entity.Goal, error)
	UpsertGoal(ctx context.Context, input *entity.GoalInput) (*entity.Goal, error)
	DeleteGoal(ctx context.Context, id string) (*entity.Goal, error)
}

type IMetricBusiness interface {
	GetMetrics(ctx context.Context) ([]entity.Metric, error)
	UpsertMetric(ctx context.Context, input *entity.MetricInput) (*entity.Metric, error)
	DeleteMetric(ctx context.Context, id string) (*entity.Metric, error)
}

type ICategoryBusiness interface {
	GetCategories(ctx context.Context) ([]entity.Category, error)
	UpsertCategory(ctx context.Context, input entity.CategoryInput) (*entity.Category, error)
	DeleteCategory(ctx context.Context, id string) (*entity.Category, error)
}

type IHabitBusiness interface {
	GetHabits(ctx context.Context) ([]entity.Habit, error)
	UpsertHabit(ctx context.Context, input *entity.HabitInput) (*entity.Habit, error)
	DeleteHabit(ctx context.Context, id string) (*entity.Habit, error)
	GetHabitLogs(ctx context.Context, filter *entity.HabitLogFilter, sort *entity.HabitLogOrderBy, limit, offset *int) ([]entity.HabitLog, error)
	UpsertHabitLog(ctx context.Context, input *entity.HabitLogInput) (*entity.HabitLog, error)
}

type ITimeTrackingBusiness interface {
	UpsertTimeTracking(ctx context.Context, input *entity.TimeTrackingInput) (*entity.TimeTracking, error)
}

type ITaskBusiness interface {
	GetTasks(ctx context.Context, filter *entity.TaskFilter) ([]entity.Task, error)
	UpsertTask(ctx context.Context, input *entity.TaskInput) (*entity.Task, error)
	DeleteTask(ctx context.Context, id string) (*entity.Task, error)
	GetTaskSessions(ctx context.Context, filter *entity.TaskSessionFilter) ([]entity.TaskSession, error)
	UpsertTaskSession(ctx context.Context, input *entity.TaskSessionInput) (*entity.TaskSession, error)
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
	Exist(ctx context.Context, characterID, metricID string) error
	FindByCharacterID(ctx context.Context, characterID string) ([]entity.Metric, error)
	UnassignCategory(ctx context.Context, categoryID string) error
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
}

type IGoalRepo interface {
	base.IBaseRepo[entity.Goal]
	GetGoalsByCharacterID(ctx context.Context, characterID string) ([]entity.Goal, error)
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, goalID string) error
}

type ICache interface {
	GetAuthSession(ctx context.Context, firebaseUID string) (*rdb.AuthSession, error)
	SetAuthSession(ctx context.Context, profile *entity.Profile, session *rdb.AuthSession) error
	DeleteAuthSession(ctx context.Context, firebaseUID string) error
	DeleteProfileData(ctx context.Context, profile *entity.Profile) error
}

type IHabitRepo interface {
	base.IBaseRepo[entity.Habit]
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	CountByCategoryID(ctx context.Context, categoryID string) (int, error)
	CountUnassigned(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, habitID string) error
	FindByCharacterID(ctx context.Context, characterID string) ([]entity.Habit, error)
	FindByCharacterIDs(ctx context.Context, characterIDs []string) ([]entity.Habit, error)
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
	UnassignCategory(ctx context.Context, categoryID string) error
}

type IHabitLogRepo interface {
	base.IBaseRepo[entity.HabitLog]
	Find(ctx context.Context, pineline entity.HabitLogPineline) ([]entity.HabitLog, error)
	UpsertByTimestamp(ctx context.Context, timestamp time.Time, habit *entity.HabitLog) error
	DeleteByHabitID(ctx context.Context, habitID string) error
	DeleteByHabitIDs(ctx context.Context, habitIDs []string) error
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
	CountByCharacterID(ctx context.Context, characterID string) (int, error)
	CountByCategoryID(ctx context.Context, categoryID string) (int, error)
	CountUnassigned(ctx context.Context, characterID string) (int, error)
	Exist(ctx context.Context, characterID, taskID string) error
	Find(ctx context.Context, pineline entity.TaskPineline) ([]entity.Task, error)
	DeleteByCharacterID(ctx context.Context, characterID string) error
	DeleteByCharacterIDs(ctx context.Context, characterIDs []string) error
	UnassignCategory(ctx context.Context, categoryID string) error
}

type ITaskSessionRepo interface {
	base.IBaseRepo[entity.TaskSession]
	Find(ctx context.Context, pineline entity.TaskSessionPineline) ([]entity.TaskSession, error)
	DeleteByTaskID(ctx context.Context, taskID string) error
	DeleteByTaskIDs(ctx context.Context, taskIDs []string) error
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
