package repo_test

import (
	"os"
	"testing"

	"tenkhours/pkg/db"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb         *mongo.Database
	profilesRepo   *repo.ProfilesRepo
	charactersRepo *repo.CharactersRepo
)

func TestMain(m *testing.M) {
	testdb = db.InitDBManagerFromURL("mongodb://localhost:27017", "test").DB
	profilesRepo = repo.NewProfilesRepo(testdb)
	charactersRepo = repo.NewCharactersRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
