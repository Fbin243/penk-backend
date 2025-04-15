package mongodb

import (
	"log"

	"tenkhours/pkg/db/base"

	"github.com/jinzhu/copier"
)

func ToMongoEntity[M base.IBaseEntity, N any](entity *M) *N {
	mongoEntity := new(N)
	err := copier.Copy(mongoEntity, entity)
	if err != nil {
		log.Printf("Error copying entity: %v", err)
	}

	return mongoEntity
}

func ToEntity[N any, M base.IBaseEntity](mongoEntity *N) *M {
	entity := new(M)
	err := copier.Copy(entity, mongoEntity)
	if err != nil {
		log.Printf("Error copying entity: %v", err)
	}
	return entity
}
