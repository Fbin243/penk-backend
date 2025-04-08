package rpc

import "github.com/jinzhu/copier"

// MapEntityToRPC maps an entity to a gRPC response
func MapEntityToRPC[E, R any](entity *E, converters []copier.TypeConverter) (*R, error) {
	resp := new(R)
	err := copier.CopyWithOption(resp, entity, copier.Option{
		Converters: converters,
	})
	return resp, err
}

// MapRPCInputToEntityInput maps a gRPC request to an entity input
func MapRPCInputToEntityInput[R, E any](req *R, converters []copier.TypeConverter) (*E, error) {
	entity := new(E)
	err := copier.CopyWithOption(entity, req, copier.Option{
		Converters: converters,
	})
	return entity, err
}
