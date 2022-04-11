package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kabojnk/latdict-api/cache"
	"github.com/kabojnk/latdict-api/db"
	"github.com/kabojnk/latdict-api/entry_filter"
	"github.com/kabojnk/latdict-api/query_filter"
	"github.com/kabojnk/latdict-api/types"
	"net/http"
	"strconv"
)

const DEFAULT_PAGE_SIZE = 20
const MAX_PAGE_SIZE = 200

// @BasePath /entries

// PingExample godoc
// @Summary Get entries
// @Schemes
// @Description Gets Latin dictionary entries
// @Tags get
// @Accept json
// @Produce json
// @Success 200 {object} types.EntriesResponse
// @Router /entries [get]
func getEntries(c *gin.Context) {

	client := db.DBClient{}
	client.Init()
	filter := query_filter.QueryFilter{}
	queryString := c.Request.URL.Query()
	filter.InitWithQueryString(queryString)
	pagination := types.Pagination{}
	pagination.PageNum, _ = strconv.Atoi(queryString.Get("page"))
	pagination.PageSize, _ = strconv.Atoi(queryString.Get("pageSize"))
	if pagination.PageSize == 0 {
		pagination.PageSize = DEFAULT_PAGE_SIZE
	} else if pagination.PageSize > MAX_PAGE_SIZE {
		pagination.PageSize = MAX_PAGE_SIZE
	}
	searchCache := cache.SearchCache{}
	var response types.EntriesResponse
	if searchCache.IsEnabled() {
		searchCache.Open()
		response, err := searchCache.GetCache(filter.Language, filter.QueryText, filter.NeedsExactMatch, pagination.PageNum, pagination.PageSize)
		if err != nil {
			fmt.Printf("Unable to find cache. Getting value from DB...")
			response = getEntriesFromDB(client, pagination, filter)
			searchCache.SaveCache(filter.Language, filter.QueryText, filter.NeedsExactMatch, pagination.PageNum, pagination.PageSize, response)
		} else {
			fmt.Printf("Found cached response and returning it.")
		}
		searchCache.Close()
	} else {
		response = getEntriesFromDB(client, pagination, filter)
	}
	client.Close()
	c.IndentedJSON(http.StatusOK, response)
}

func getEntriesFromDB(client db.DBClient, pagination types.Pagination, filter query_filter.QueryFilter) types.EntriesResponse {
	entries, totalEntries := client.GetEntries(pagination, filter)
	pagination.TotalNumPages = totalEntries
	// @TODO: See if this is more performant than doing a single loop with single SQL queries, or something where we
	//        can get by with one less loop without violating some of Go's mutation limitations.
	if filter.IncludeSenses {
		var entryIDs []int
		for _, entry := range entries {
			entryIDs = append(entryIDs, entry.ID)
		}
		sensesMap := client.GetSensesForEntryIDs(entryIDs)
		for i, entry := range entries {
			entries[i].Senses = sensesMap[entry.ID]
		}
	}
	response := types.EntriesResponse{
		Pagination: pagination,
		Items:      entries,
	}
	return response
}

// PingExample godoc
// @Summary Get entry by UUID
// @Schemes
// @Description Gets Latin dictionary entry by its UUID
// @Tags get
// @Accept json
// @Produce json
// @Success 200 {object} types.Entry
// @Router /entries/:entryUUID [get]
func getEntry(c *gin.Context) {
	var entryURI types.EntryURI
	if err := c.ShouldBindUri(&entryURI); err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}
	client := db.DBClient{}
	client.Init()
	filter := entry_filter.EntryFilter{}
	queryString := c.Request.URL.Query()
	filter.InitWithQueryString(queryString)

	dbEntry := client.GetEntryByUUID(entryURI.EntryUUID)
	orthography := dbEntry.Orthography.String
	entry := types.Entry{
		ID:               dbEntry.ID,
		UUID:             dbEntry.UUID,
		Lemma:            dbEntry.Lemma,
		CommonalityScore: dbEntry.CommonalityScore,
		Speech:           dbEntry.Speech,
		Orthography:      orthography,
	}
	fmt.Printf("EntryID: %d\n", dbEntry.ID)
	if filter.IncludeSenses {
		senses := client.GetSenseForEntryID(dbEntry.ID)
		entry.Senses = senses
	}
	if filter.IncludeGrammarInfo {
		dbGrammarValues := client.GetGrammarValuesForEntryIDs(dbEntry.ID)
		var apiGrammarValues []types.APIGrammarValues
		for _, dbGrammarValue := range dbGrammarValues {
			apiGrammarValues = append(apiGrammarValues, types.APIGrammarValues{
				ID:         dbGrammarValue.ID,
				EntryID:    dbGrammarValue.EntryID,
				UUID:       dbGrammarValue.UUID,
				GrammarKey: dbGrammarValue.GrammarKey,
				Value:      dbGrammarValue.Value.String,
			})
		}
		entry.GrammarValues = apiGrammarValues
	}
	if filter.IncludeAdditionalInfo {
		dbAdditionalInfo := client.GetAdditionalInfoForEntryID(dbEntry.ID)
		apiAdditionalInfo := types.APIAdditionalInfo{
			ID:        dbAdditionalInfo.ID,
			UUID:      dbAdditionalInfo.UUID,
			EntryID:   dbAdditionalInfo.EntryID,
			Age:       dbAdditionalInfo.Age.String,
			Context:   dbAdditionalInfo.Context.String,
			Frequency: dbAdditionalInfo.Frequency.String,
			Geography: dbAdditionalInfo.Geography.String,
			Source:    dbAdditionalInfo.Source.String,
		}
		entry.AdditionalInfo = apiAdditionalInfo
	}

	c.IndentedJSON(http.StatusOK, entry)
}
