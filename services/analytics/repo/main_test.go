package repo_test

import (
	"context"
	"os"
	"testing"

	"tenkhours/pkg/db"
	"tenkhours/services/analytics/repo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testdb        *mongo.Database
	snapshotsRepo *repo.SnapshotsRepo
)

func TestMain(m *testing.M) {
	testdb = db.InitDBManagerFromURL("mongodb://localhost:27017", "test").DB
	testdb.CreateCollection(context.Background(), db.SnapshotsCollection,
		options.CreateCollection().
			SetTimeSeriesOptions(
				options.TimeSeries().
					SetTimeField("timestamp").
					SetMetaField("metadata"),
			),
	)

	snapshotsRepo = repo.NewSnapshotsRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
