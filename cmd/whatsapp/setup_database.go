package main

import (
	"database/sql"

	"github.com/defryheryanto/ai-assistant/config"
	_ "github.com/lib/pq"
)

func setupDatabaseConnection() *sql.DB {
	conn, err := sql.Open("postgres", config.DatabaseConnectionString)
	if err != nil {
		panic(err)
	}

	return conn
}
