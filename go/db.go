package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetDB(batch bool) (*sqlx.DB, error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = GetEnv("MYSQL_HOSTNAME", "127.0.0.1") + ":" + GetEnv("MYSQL_PORT", "3306")
	mysqlConfig.User = GetEnv("MYSQL_USER", "isucon")
	mysqlConfig.Passwd = GetEnv("MYSQL_PASS", "isucon")
	mysqlConfig.DBName = GetEnv("MYSQL_DATABASE", "isucholar")
	mysqlConfig.Params = map[string]string{
		"time_zone": "'+00:00'",
	}
	mysqlConfig.ParseTime = true
	mysqlConfig.MultiStatements = batch
	mysqlConfig.InterpolateParams = true

	if batch {
		return sqlx.Open("mysql", mysqlConfig.FormatDSN())
	}

	db, err := sql.Open(tracedDriver("mysql"), mysqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}

	initDB(db)

	return sqlx.NewDb(db, "mysql"), nil
}

func initDB(db *sql.DB) {
	waitDB(db)
	go pollDB(db)

	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxIdleConns(128)
	db.SetMaxOpenConns(128)
}

func waitDB(db *sql.DB) {
	for {
		err := db.Ping()
		if err == nil {
			return
		}

		log.Printf("Failed to ping DB: %s", err)
		log.Println("Retrying...")
		time.Sleep(time.Second)
	}
}

func pollDB(db *sql.DB) {
	for {
		err := db.Ping()
		if err != nil {
			log.Printf("Failed to ping DB: %s", err)
		}

		time.Sleep(time.Second)
	}
}
