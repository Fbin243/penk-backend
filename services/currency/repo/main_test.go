package repo_test

import (
	"os"
	"testing"

	"tenkhours/pkg/db"
	"tenkhours/services/currency/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb   *mongo.Database
	fishRepo *repo.FishRepo
)

func TestMain(m *testing.M) {
	testdb = db.InitDBManagerFromURL("mongodb://localhost:27017", "test").DB
	fishRepo = repo.NewFishRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
