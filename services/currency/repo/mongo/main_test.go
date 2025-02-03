package mongorepo_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"tenkhours/services/currency/business"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"

	mongorepo "tenkhours/services/currency/repo/mongo"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb   *mongo.Database
	fishRepo business.IFishRepo
)

func TestMain(m *testing.M) {
	if godotenv.Load(filepath.Join(utils.GetRoot(), ".env.test")) != nil {
		log.Fatal("Error loading .env.test" + " file")
	}

	testdb = mongodb.InitDBManagerFromEnv().DB
	fishRepo = mongorepo.NewFishRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
