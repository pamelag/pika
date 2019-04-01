package main

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultPort   = "8080"
	defaultDBURL  = "152.28.1.1"
	defaultDBPort = "3306"
	defaultDBUser = "pamela"
	defaultDBPass = "1234"
	defaultDBName = "ny_cab_data"
	retries       = 30
)

var (
	port   = envString("PORT", defaultPort)
	dburl  = envString("DB_HOST", defaultDBURL)
	dbport = envString("DB_PORT", defaultDBPort)
	dbuser = envString("DB_USER", defaultDBUser)
	dbpass = envString("DB_PASS", defaultDBPass)
	dbname = envString("DB_NAME", defaultDBName)
	//dbtype = envString("DB_TYPE", defaultDBType)

	httpAddr = ":" + port
)

func getConnectionString() string {
	var buf bytes.Buffer
	buf.WriteString(dbuser)
	buf.WriteString(":")
	buf.WriteString(dbpass)
	buf.WriteString("@tcp(")
	buf.WriteString(dburl)
	buf.WriteString(":")
	buf.WriteString(dbport)
	buf.WriteString(")/")
	buf.WriteString(dbname)
	return buf.String()
}

func connect() (*sql.DB, error) {
	return sql.Open("mysql", getConnectionString())
}

func getConnection() (*sql.DB, error) {
	var connection *sql.DB
	var err error
	var n uint
	for n < retries {
		connection, err = connect()
		if err != nil {
			log.Println("retrying to connect to db, ", n)
			time.Sleep(10 * time.Second)
			n++
			continue
		}

		break

	}

	return connection, err
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
