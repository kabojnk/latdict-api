package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kabojnk/latdict-api/query_filter"
	_ "github.com/lib/pq"
	"time"
)

type DBClient struct {
	Config DBConfig
	DB     *sqlx.DB
}

func (client *DBClient) Init() {
	client.Config = DBConfig{}
	client.Config.Init()
}

// Open Opens a Postgres connection
func (client *DBClient) Open() {
	shouldEnableSSL := "disable"
	if !client.Config.SSLModeDisable {
		shouldEnableSSL = "enable"
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", client.Config.Host, client.Config.Port, client.Config.User, client.Config.Password, client.Config.DBName, shouldEnableSSL)
	pqDb, pqErr := sqlx.Open("postgres", dsn)
	if pqErr != nil {
		panic(pqErr)
	}
	client.DB = pqDb
	client.DB.SetConnMaxLifetime(time.Minute * 3)
	client.DB.SetMaxOpenConns(10)
	client.DB.SetMaxIdleConns(10)
}

// Close closes a Postgres connection
func (client *DBClient) Close() {
	if client.DB != nil {
		err := client.DB.Close()
		if err != nil {
			panic(err)
		}
	}
}

// GetEntries Gets a list of entries from the lexicon.
func (client *DBClient) GetEntries(pagination Pagination, filter query_filter.QueryFilter) ([]Entry, int) {
	client.Open()
	totalResults := 0
	rows, err := client.queryBasedOnFilterAndPagination(filter, pagination)
	if err != nil {
		panic(err)
	}
	var entries []Entry
	defer rows.Close()
	for rows.Next() {
		searchKeyResult := DBSearchKeyResult{}
		err := rows.StructScan(&searchKeyResult)
		if err != nil {
			panic(err)
		}
		// Query entry
		dbEntry := client.GetEntryByID(searchKeyResult.EntryID)
		if totalResults == 0 {
			totalResults = searchKeyResult.TotalResults
		}
		orthography := dbEntry.Orthography.String
		entries = append(entries, Entry{
			UUID:             dbEntry.UUID,
			Lemma:            dbEntry.Lemma,
			CommonalityScore: dbEntry.CommonalityScore,
			Speech:           dbEntry.Speech,
			Orthography:      orthography,
		})
	}
	client.Close()
	return entries, totalResults
}

// GetEntryByID Gets an entry by its PK
func (client *DBClient) GetEntryByID(entryId int) DBEntry {
	dbEntry := DBEntry{}
	row := client.DB.QueryRowx(`select * 
		from entries 
		WHERE id = $1`, entryId)
	err := row.StructScan(&dbEntry)
	if err != nil {
		panic(err)
	}
	return dbEntry
}

// Performs a SQL query based on filter and pagination data and returns the result (or errors)
func (client *DBClient) queryBasedOnFilterAndPagination(filter query_filter.QueryFilter, pagination Pagination) (*sqlx.Rows, error) {

	table := "latin"
	if filter.Language == "english" {
		table = "english"
	}
	query := fmt.Sprintf(`select distinct sk.entry_id 
			from search_keys_%s 
			where sk.search_key = $1
			order by commonality_score desc LIMIT $2 OFFSET $3
			limit $4 
			offset $5`, table)
	if filter.NeedsExactMatch {
		return client.DB.Queryx(query,
			filter.QueryText,
			pagination.PageSize,
			pagination.PageNum*pagination.PageSize)
	}
	// Important note, it's cheaper for us to do our full, starting, and partial matching in our select statements, and
	// using the widest net (WHERE LIKE %query%) as the sole WHERE statement.
	query = fmt.Sprintf(`select distinct sk.entry_id,
		length(sk.search_key) as search_key_length,
		count(sk.entry_id) OVER() as total_results,
		commonality_score,
		case when sk.search_key = $1 then 3 else 0 end as full_match, 
		case when sk.search_key like $2 then 2 else 0 end as starting_match, 
		case when sk.search_key like $3 then 1 else 0 end as partial_match 
		from search_keys_%s sk
		where sk.search_key LIKE $3
		order by full_match desc, starting_match desc, partial_match desc, sk.commonality_score desc, search_key_length ASC
		limit $4
		offset $5`, table)
	return client.DB.Queryx(query,
		filter.QueryText,
		fmt.Sprintf("%%%s", filter.QueryText),
		fmt.Sprintf("%%%s%%", filter.QueryText),
		pagination.PageSize,
		pagination.PageNum*pagination.PageSize)
}
