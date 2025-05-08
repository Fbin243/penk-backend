package redisrepo

import (
	"time"

	"tenkhours/pkg/db/base"
	core_entity "tenkhours/services/core/entity"

	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
)

func ToRedisZ(reminder *core_entity.Reminder) *redis.Z {
	var score float64
	if reminder.RemindTime != nil {
		score = float64(reminder.RemindTime.Unix())
	}

	return &redis.Z{
		Member: reminder.ID,
		Score:  score,
	}
}

func ToReminder(z *redis.Z) *core_entity.Reminder {
	var remindTime *time.Time
	if z.Score != 0 {
		remindTime = lo.ToPtr(time.Unix(int64(z.Score), 0))
	}

	return &core_entity.Reminder{
		BaseEntity: &base.BaseEntity{
			ID: z.Member.(string),
		},
		RemindTime: remindTime,
	}
}
