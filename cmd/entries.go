package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kabojnk/latdict-api/query_filter"
	"net/http"
	"strconv"
)

const DEFAULT_PAGE_SIZE = 20
const MAX_PAGE_SIZE = 200

func getEntries(c *gin.Context) {

	client := DBClient{}
	client.Init()
	filter := query_filter.QueryFilter{}
	queryString := c.Request.URL.Query()
	filter.InitWithQueryString(queryString)
	pagination := Pagination{}
	pagination.PageNum, _ = strconv.Atoi(queryString.Get("page"))
	pagination.PageSize, _ = strconv.Atoi(queryString.Get("pageSize"))
	if pagination.PageSize == 0 {
		pagination.PageSize = DEFAULT_PAGE_SIZE
	} else if pagination.PageSize > MAX_PAGE_SIZE {
		pagination.PageSize = MAX_PAGE_SIZE
	}
	entries, totalEntries := client.GetEntries(pagination, filter)
	pagination.TotalNumPages = totalEntries
	response := EntriesResponse{
		Pagination: pagination,
		Items:      entries,
	}
	c.IndentedJSON(http.StatusOK, response)
}
