package api

import (
	postHandler "belajar-api-golang/api/handler/post"
	userHandler "belajar-api-golang/api/handler/user"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, db *sql.DB) {
	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello, World!",
			})
		})

		//user handler
		api.GET("/user/all", userHandler.GetAll(db))
		api.GET("/user/:id", userHandler.Get(db))
		api.POST("/user", userHandler.Create(db))
		api.PUT("/user/:id", userHandler.Update(db))
		api.DELETE("/user/:id", userHandler.Delete(db))

		//post handler
		api.GET("/post/all", postHandler.GetAll(db))
		api.GET("/post/:id", postHandler.Get(db))
		api.POST("/post", postHandler.Create(db))
		api.PUT("/post/:id", postHandler.Update(db))
		api.DELETE("/post/:id", postHandler.Delete(db))
	}
}
