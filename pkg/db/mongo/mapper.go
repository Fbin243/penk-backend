package mongodb

import (
	"log"

	"tenkhours/pkg/db/base"
	"tenkhours/pkg/types"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
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

func ToPaginationPineline(paginate *types.Pagination) []bson.M {
	if paginate == nil {
		return nil
	}

	pipeline := []bson.M{}
	if paginate.Limit != nil {
		pipeline = append(pipeline, bson.M{"$limit": *paginate.Limit})
	}
	if paginate.Offset != nil {
		pipeline = append(pipeline, bson.M{"$skip": *paginate.Offset})
	}
	return pipeline
}
