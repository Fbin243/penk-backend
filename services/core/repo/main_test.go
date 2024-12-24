package repo_test

import (
	"context"
	"os"
	"testing"

	"tenkhours/pkg/db"
	"tenkhours/pkg/utils"
	"tenkhours/services/core/repo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testdb                 *mongo.Database
	profilesRepo           *repo.ProfilesRepo
	charactersRepo         *repo.CharactersRepo
	goalsRepo              *repo.GoalsRepo
	templatesRepo          *repo.TemplatesRepo
	templateCategoriesRepo *repo.TemplateCategoriesRepo
	snapshotsRepo          *repo.SnapshotsRepo
)

func TestMain(m *testing.M) {
	// Local version
	// testdb = db.InitDBManagerFromURL("mongodb://localhost:27017", "test").DB

	// Remote version
	utils.ReadEnvFile()
	testdb = db.InitDBManagerFromEnv("TenK-Hours-Dev").DB

	profilesRepo = repo.NewProfilesRepo(testdb)
	charactersRepo = repo.NewCharactersRepo(testdb)
	goalsRepo = repo.NewGoalsRepo(testdb)
	templatesRepo = repo.NewTemplatesRepo(testdb)
	templateCategoriesRepo = repo.NewTemplateCategoriesRepo(testdb)
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
