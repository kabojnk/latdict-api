package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kabojnk/latdict-api/query_filter"
	"github.com/kabojnk/latdict-api/types"
	"github.com/lib/pq"
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
func (client *DBClient) GetEntries(pagination types.Pagination, filter query_filter.QueryFilter) ([]types.Entry, int) {
	client.Open()
	totalResults := 0
	rows, err := client.queryBasedOnFilterAndPagination(filter, pagination)
	if err != nil {
		panic(err)
	}
	var entries []types.Entry
	defer rows.Close()
	for rows.Next() {
		searchKeyResult := types.DBSearchKeyResult{}
		err := rows.StructScan(&searchKeyResult)
		if err != nil {
			panic(err)
		}
		// Query entry -- don't get any extra data. Senses will be appended in a separate group query
		dbEntry := client.GetEntryByID(searchKeyResult.EntryID)
		if totalResults == 0 {
			totalResults = searchKeyResult.TotalResults
		}
		orthography := dbEntry.Orthography.String
		entries = append(entries, types.Entry{
			ID:               dbEntry.ID,
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
func (client *DBClient) GetEntryByID(entryId int) types.DBEntry {
	dbEntry := types.DBEntry{}
	row := client.DB.QueryRowx(`select * 
		from entries 
		WHERE id = $1`, entryId)
	err := row.StructScan(&dbEntry)
	if err != nil {
		panic(err)
	}
	return dbEntry
}

func (client *DBClient) GetEntryByUUID(entryUUID string) types.DBEntry {
	client.Open()
	dbEntry := types.DBEntry{}
	row := client.DB.QueryRowx(`select * 
		from entries 
		WHERE uuid = $1`, entryUUID)
	err := row.StructScan(&dbEntry)
	if err != nil {
		panic(err)
	}
	client.Close()
	return dbEntry
}

func (client *DBClient) GetSenseForEntryID(entryID int) []types.Sense {
	var senses []types.Sense
	client.Open()
	err := client.DB.Select(&senses, "select * from senses where entry_id = $1", entryID)
	if err != nil {
		panic(err)
	}
	client.Close()
	return senses
}

func (client *DBClient) GetSensesForEntryIDs(entryIDs []int) map[int][]types.Sense {
	sensesMap := make(map[int][]types.Sense)
	fmt.Printf("Entry IDs: %v\n", entryIDs)
	client.Open()
	rows, err := client.DB.Queryx("select * from senses where entry_id = any($1) order by senses.order ASC", pq.Array(entryIDs))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		sense := types.Sense{}
		err := rows.StructScan(&sense)
		if err != nil {
			panic(err)
		}
		sensesMap[sense.EntryID] = append(sensesMap[sense.EntryID], sense)
	}
	client.Close()
	return sensesMap
}

func (client *DBClient) GetAdditionalInfoForEntryID(entryID int) types.DBAdditionalInfo {
	var additionalInfo types.DBAdditionalInfo
	client.Open()
	row := client.DB.QueryRowx("select * from entry_additional_info where entry_id = $1", entryID)
	err := row.StructScan(&additionalInfo)
	if err != nil {
		panic(err)
	}
	client.Close()
	return additionalInfo
}

func (client *DBClient) GetGrammarValuesForEntryIDs(entryID int) []types.DBGrammarValues {
	var grammarValues []types.DBGrammarValues
	client.Open()
	err := client.DB.Select(&grammarValues, "select * from entry_grammar_values where entry_id = $1", entryID)
	if err != nil {
		panic(err)
	}
	client.Close()
	return grammarValues
}

// Performs a SQL query based on filter and pagination data and returns the result (or errors)
func (client *DBClient) queryBasedOnFilterAndPagination(filter query_filter.QueryFilter, pagination types.Pagination) (*sqlx.Rows, error) {

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
