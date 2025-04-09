package mongodb

import (
	"tenkhours/pkg/db/base"

	"github.com/jinzhu/copier"
)

func ToMongoEntity[M base.IBaseEntity, N any](entity *M) *N {
	mongoEntity := new(N)
	copier.Copy(mongoEntity, entity)
	return mongoEntity
}

func ToEntity[N any, M base.IBaseEntity](mongoEntity *N) *M {
	entity := new(M)
	copier.Copy(entity, mongoEntity)
	return entity
}
