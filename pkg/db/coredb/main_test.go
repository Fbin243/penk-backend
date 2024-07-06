package coredb

import (
	"os"
	"testing"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testdb         *mongo.Database
	usersRepo      *UsersRepo
	charactersRepo *CharactersRepo
)

func TestMain(m *testing.M) {
	testdb = db.InitDBFromURL("mongodb://localhost:27017", "test")
	usersRepo = NewUsersRepo(testdb)
	charactersRepo = NewCharactersRepo(testdb)

	code := m.Run()

	os.Exit(code)
}
