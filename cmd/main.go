package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.GET("/entries", getEntries)
	router.Run("localhost:8000")
}
