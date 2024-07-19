package api

import (
	"log"
	"net/http"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/core"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
)

type App struct {
	Engine *gin.Engine
}

func NewApp() *App {
	return &App{
		Engine: gin.Default(),
	}
}

func (app *App) InitRouter() {
	log.Println("--> Hello from Core service")

	app.Engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	app.Engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Core service is running!")
	})

	authMiddleware := auth.NewMiddleware(coredb.NewUsersRepo(db.GetDBManager().DB))

	app.Engine.Use(authMiddleware.CheckRequestBody)

	app.Engine.Use(authMiddleware.CheckAuth)

	app.Engine.POST("/graphql", func(c *gin.Context) {
		postData, ok := c.Request.Context().Value(auth.PostDataKey).(utils.GraphqlQueryData)
		if !ok {
			c.String(http.StatusBadRequest, "invalid post data")
			return
		}

		result := graphql.Do(graphql.Params{
			Context:        c.Request.Context(),
			Schema:         core.InitSchema(),
			RequestString:  postData.Query,
			VariableValues: postData.Variables,
			OperationName:  postData.Operation,
		})

		// Log errors
		if len(result.Errors) > 0 {
			log.Printf("errors: %v", result.Errors)
		}

		c.JSON(http.StatusOK, result)
	})
}
