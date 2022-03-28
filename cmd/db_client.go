package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

const DEFAULT_PAGE_SIZE = 20

type DBClient struct {
	Config DBConfig
	DB     *sqlx.DB
}

func (client *DBClient) Init() {
	client.Config = DBConfig{}
	client.Config.Init()
}

func (client *DBClient) Open() {
	shouldEnableSSL := "disable"
	if !client.Config.SSLModeDisable {
		shouldEnableSSL = "enable"
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", client.Config.Host, client.Config.Port, client.Config.User, client.Config.Password, client.Config.DBName, shouldEnableSSL)
	fmt.Println(dsn)
	pqDb, pqErr := sqlx.Open("postgres", dsn)
	if pqErr != nil {
		panic(pqErr)
	}
	client.DB = pqDb
	client.DB.SetConnMaxLifetime(time.Minute * 3)
	client.DB.SetMaxOpenConns(10)
	client.DB.SetMaxIdleConns(10)
}

func (client *DBClient) Close() {
	if client.DB != nil {
		err := client.DB.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (client *DBClient) GetEntries(pagination Pagination, filter QueryFilter) []APIEntry {
	client.Open()
	pageSize := pagination.PageSize
	if pageSize == 0 {
		pageSize = DEFAULT_PAGE_SIZE
	}
	rows, err := client.DB.Queryx("SELECT * FROM entries LIMIT $1 OFFSET $2", pageSize, pagination.PageNum*pageSize)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var entries []APIEntry
	for rows.Next() {
		dbEntry := DBEntry{}
		err := rows.StructScan(&dbEntry)
		if err != nil {
			panic(err)
		}
		orthography := dbEntry.Orthography.String
		entries = append(entries, APIEntry{
			UUID:             dbEntry.UUID,
			Lemma:            dbEntry.Lemma,
			CommonalityScore: dbEntry.CommonalityScore,
			Speech:           dbEntry.Speech,
			Orthography:      orthography,
		})
	}
	client.Close()
	return entries
}
