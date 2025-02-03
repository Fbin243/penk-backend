package mongorepo_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"tenkhours/services/core/business"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"

	mongorepo "tenkhours/services/core/repo/mongo"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb               *mongo.Database
	profileRepo          business.IProfileRepo
	characterRepo        business.ICharacterRepo
	goalRepo             business.IGoalRepo
	templateRepo         business.ITemplateRepo
	templateCategoryRepo business.ITemplateCategoryRepo
)

func TestMain(m *testing.M) {
	if godotenv.Load(filepath.Join(utils.GetRoot(), ".env.test")) != nil {
		log.Fatal("Error loading .env.test" + " file")
	}

	testdb = mongodb.InitDBManagerFromEnv().DB

	profileRepo = mongorepo.NewProfileRepo(testdb)
	characterRepo = mongorepo.NewCharacterRepo(testdb)
	goalRepo = mongorepo.NewGoalRepo(testdb)
	templateRepo = mongorepo.NewTemplateRepo(testdb)
	templateCategoryRepo = mongorepo.NewTemplateCategoryRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
