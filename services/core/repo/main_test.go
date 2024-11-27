package repo_test

import (
	"os"
	"testing"

	"tenkhours/pkg/db"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb         *mongo.Database
	profilesRepo   *repo.ProfilesRepo
	charactersRepo *repo.CharactersRepo
	goalsRepo      *repo.GoalsRepo
)

func TestMain(m *testing.M) {
	// Local version
	// testdb = db.InitDBManagerFromURL("mongodb://localhost:27017", "test").DB

	// Remote version
	utils.ReadEnvFile()
	testdb = db.InitDBManagerFromEnv("test").DB

	profilesRepo = repo.NewProfilesRepo(testdb, db.GetRedisClient())
	charactersRepo = repo.NewCharactersRepo(testdb)
	goalsRepo = repo.NewGoalsRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
