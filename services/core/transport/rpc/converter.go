package rpc

import (
	"fmt"
	"time"

	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
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
}
