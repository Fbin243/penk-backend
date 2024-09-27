package coredb

import (
	"os"
	"testing"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb         *mongo.Database
	profilesRepo   *ProfilesRepo
	charactersRepo *CharactersRepo
)

func TestMain(m *testing.M) {
	testdb = db.InitDBManagerFromURL("mongodb://localhost:27017", "test").DB
	profilesRepo = NewProfilesRepo(testdb)
	charactersRepo = NewCharactersRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
