package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getEntries(c *gin.Context) {
	client := DBClient{}
	client.Init()
	filter := QueryFilter{}
	pagination := Pagination{}
	entries := client.GetEntries(pagination, filter)
	c.IndentedJSON(http.StatusOK, entries)
}
