package mongorepo_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	mongorepo "tenkhours/services/analytic/repo/mongo"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb             *mongo.Database
	capturedRecordRepo *mongorepo.CapturedRecordRepo
)

func TestMain(m *testing.M) {
	if godotenv.Load(filepath.Join(utils.GetRoot(), ".env.test")) != nil {
		log.Fatal("Error loading .env.test" + " file")
	}

	testdb = mongodb.InitDBManagerFromEnv().DB
	capturedRecordRepo = mongorepo.NewCapturedRecordRepo(testdb)

	os.Exit(m.Run())
}
