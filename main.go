package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type word struct {
	ID string `json:"id"`
	Entry string `json:"entry"`
}

var words = []word {
	{
		ID: "ab-1",
		Entry: "ab, a",
	},
	{
		ID: "ad-1",
		Entry: "ad",
	},
	{
		ID: "ex-1",
		Entry: "ex, e",
	},
}

func getWords(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, words)
}

func main() {
	router := gin.Default()
	router.GET("/words", getWords)

	router.Run("localhost:8080")
}


