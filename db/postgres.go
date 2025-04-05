package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	
	psqlInfo := "host=localhost port=5432 user=postgres dbname=filesharing sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open DB: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	fmt.Println("Hardcoded Connected to PostgreSQL!")
}
