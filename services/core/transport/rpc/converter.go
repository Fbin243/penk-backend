package rpc

import (
	"fmt"
	"time"

	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
	"github.com/samber/lo"
)

var MetricConditionConverter = []copier.TypeConverter{
	{
		SrcType: core.MetricCondition(0),
		DstType: entity.MetricCondition(fmt.Sprint(0)),
		Fn: func(src any) (any, error) {
			return entity.MetricCondition(
				core.MetricCondition_name[int32(src.(core.MetricCondition).Number())],
			), nil
		},
	},
	{
		SrcType: entity.MetricCondition(fmt.Sprint(0)),
		DstType: core.MetricCondition(0),
		Fn: func(src any) (any, error) {
			return core.MetricCondition(
				core.MetricCondition_value[string(src.(entity.MetricCondition))],
			), nil
		},
	},
}

var UnixTimeConverter = []copier.TypeConverter{
	{
		SrcType: int64(0),
		DstType: time.Time{},
		Fn: func(src any) (any, error) {
			return time.Unix(src.(int64), 0), nil
		},
	},
	{
		SrcType: time.Time{},
		DstType: int64(0),
		Fn: func(src any) (any, error) {
			return src.(time.Time).Unix(), nil
		},
	},
	// Handle for time.Time to *int64
	{
		SrcType: time.Time{},
		DstType: lo.ToPtr(int64(0)),
		Fn: func(src any) (any, error) {
			return lo.ToPtr(src.(time.Time).Unix()), nil
		},
	},
	// Handle for *int64 to *time.Time
	{
		SrcType: lo.ToPtr(int64(0)),
		DstType: &time.Time{},
		Fn: func(src any) (any, error) {
			if src.(*int64) == nil {
				return nil, nil
			}
			return lo.ToPtr(time.Unix(*src.(*int64), 0)), nil
		},
	},
	// Handle for int64 to *time.Time
	{
		SrcType: int64(0),
		DstType: &time.Time{},
		Fn: func(src any) (any, error) {
			return time.Unix(src.(int64), 0), nil
		},
	},
}

var EntityTypeConverter = []copier.TypeConverter{
	{
		SrcType: core.EntityType(0),
		DstType: entity.EntityType(fmt.Sprint(0)),
		Fn: func(src any) (any, error) {
			return entity.EntityType(
				core.EntityType_name[int32(src.(core.EntityType).Number())],
			), nil
		},
	},
	{
		SrcType: entity.EntityType(fmt.Sprint(0)),
		DstType: core.EntityType(0),
		Fn: func(src any) (any, error) {
			return core.EntityType(
				core.EntityType_value[string(src.(entity.EntityType))],
			), nil
		},
	},
}

var HabitConverter = []copier.TypeConverter{
	{
		SrcType: core.CompletionType(0),
		DstType: entity.CompletionType(fmt.Sprint(0)),
		Fn: func(src any) (any, error) {
			return entity.CompletionType(
				core.CompletionType_name[int32(src.(core.CompletionType).Number())],
			), nil
		},
	},
	{
		SrcType: entity.CompletionType(fmt.Sprint(0)),
		DstType: core.CompletionType(0),
		Fn: func(src any) (any, error) {
			return core.CompletionType(
				core.CompletionType_value[string(src.(entity.CompletionType))],
			), nil
		},
	},
	{
		SrcType: core.HabitReset(0),
		DstType: entity.HabitReset(fmt.Sprint(0)),
		Fn: func(src any) (any, error) {
			return entity.HabitReset(
				core.HabitReset_name[int32(src.(core.HabitReset).Number())],
			), nil
		},
	},
	{
		SrcType: entity.HabitReset(fmt.Sprint(0)),
		DstType: core.HabitReset(0),
		Fn: func(src any) (any, error) {
			return core.HabitReset(
				core.HabitReset_value[string(src.(entity.HabitReset))],
			), nil
		},
	},
}
