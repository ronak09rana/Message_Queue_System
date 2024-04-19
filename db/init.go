package db

import (
	"database/sql"
	"fmt"
	"log"
	"message_queue_system/domain"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Client *sql.DB

func Init(dbCred domain.DBCred) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbCred.Username, dbCred.Password, dbCred.Hostname, dbCred.DBName))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Printf("Error: %v, unable to ping database", pingErr)
		return err
	}
	log.Printf("Connected to DB %s successfully\n", dbCred.DBName)
	Client = db
	return nil
}
