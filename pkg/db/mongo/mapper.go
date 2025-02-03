package mongodb

import "github.com/jinzhu/copier"

type Mapper[M any, N any] struct{}

func (m *Mapper[M, N]) ToMongoEntity(entity *M) *N {
	mongoEntity := new(N)
	copier.Copy(mongoEntity, entity)
	return mongoEntity
}

func (m *Mapper[M, N]) ToEntity(mongoEntity *N) *M {
	entity := new(M)
	copier.Copy(entity, mongoEntity)
	return entity
}
