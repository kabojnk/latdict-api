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
	dbConfig.Host = os.Getenv("DB_HOST")
	dbConfig.Port = os.Getenv("DB_PORT")
	dbConfig.User = os.Getenv("DB_USER")
	dbConfig.Password = os.Getenv("DB_PASSWORD")
	dbConfig.DBName = os.Getenv("DB_NAME")
	dbConfig.SSLModeDisable, _ = strconv.ParseBool(os.Getenv("DB_SSL_MODE_DISABLE"))
}
