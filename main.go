package main

import (
	"belajar-api-golang/api"
	"belajar-api-golang/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	api.SetupRouter(r, db)
	r.Run(":8080")
}
