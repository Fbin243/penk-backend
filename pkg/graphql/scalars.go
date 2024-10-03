package graphql

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MarshalNullableString(s string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if s == "" {
			w.Write([]byte(`null`))
			return
		}

		w.Write([]byte(`"` + s + `"`))
	})
}

func UnmarshalNullableString(v interface{}) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("NullableString must be a string")
	}

	return s, nil
}

func MarshalObjectID(oid primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if oid.IsZero() {
			w.Write([]byte(`null`))
			return
		}

		w.Write([]byte(`"` + oid.Hex() + `"`))
	})
}

func UnmarshalObjectID(v interface{}) (primitive.ObjectID, error) {
	id, ok := v.(string)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("ObjectID must be a string")
	}

	if id == "" {
		return primitive.NilObjectID, fmt.Errorf("ObjectID must not be empty")
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return oid, nil
}
