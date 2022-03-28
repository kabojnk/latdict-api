package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kabojnk/latdict-api/query_filter"
	"net/http"
)

func getEntries(c *gin.Context) {

	client := DBClient{}
	client.Init()
	filter := query_filter.QueryFilter{}
	filter.InitWithQueryString(c.Request.URL.Query())
	pagination := Pagination{}
	entries := client.GetEntries(pagination, filter)
	c.IndentedJSON(http.StatusOK, entries)
}
