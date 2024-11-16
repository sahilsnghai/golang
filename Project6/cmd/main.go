package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/sahilsnghai/Project6/cmd/api"
	"github.com/sahilsnghai/Project6/config"
	"github.com/sahilsnghai/Project6/database"
)

func main() {

	log.Println(config.EnvsInit)
	db, err := database.NewMySQLStore(mysql.Config{
		User:                 config.EnvsInit.DBUser,
		Passwd:               config.EnvsInit.DBPassword,
		Addr:                 config.EnvsInit.DBAddress,
		DBName:               config.EnvsInit.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}
	dbConnectionCheck(db)
	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func dbConnectionCheck(db *sql.DB) {
	log.Println("Checking database connection")
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB Connect is setup and connected")
}
