package repo_test

import (
	"os"
	"testing"

	"tenkhours/pkg/db"
	"tenkhours/services/timetrackings/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb            *mongo.Database
	timeTrackingsRepo *repo.TimeTrackingsRepo
)

func TestMain(m *testing.M) {
	testdb = db.InitDBManagerFromURL("mongodb://localhost:27017", "test").DB
	timeTrackingsRepo = repo.NewTimeTrackingsRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
