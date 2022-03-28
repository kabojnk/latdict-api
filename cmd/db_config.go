package main

import (
	"os"
	"strconv"
)

type DBConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	DBName         string
	SSLModeDisable bool
}

func (dbConfig *DBConfig) Init() {
	dbConfig.Host = os.Getenv("POSTGRES_HOST")
	dbConfig.Port = os.Getenv("POSTGRES_PORT")
	dbConfig.User = os.Getenv("POSTGRES_USER")
	dbConfig.Password = os.Getenv("POSTGRES_PASSWORD")
	dbConfig.DBName = os.Getenv("POSTGRES_NAME")
	dbConfig.SSLModeDisable, _ = strconv.ParseBool(os.Getenv("POSTGRES_SSL_MODE_DISABLE"))
}
