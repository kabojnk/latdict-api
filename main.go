package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	docs "github.com/kabojnk/latdict-api/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.BasePath = "/"

	router := gin.Default()
	router.GET("/entries", getEntries)
	router.GET("/entries/:entryUUID", getEntry)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.Run("localhost:8000")
}
