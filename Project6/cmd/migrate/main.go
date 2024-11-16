package main

import (
	"fmt"
	"log"
	"os"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sahilsnghai/Project6/config"
	"github.com/sahilsnghai/Project6/database"
)

func main() {
	db, err := database.NewMySQLStore(mysqlCfg.Config{
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

	driver, err := mysqlMigrate.WithInstance(db, &mysqlMigrate.Config{})

	defer db.Close()

	if err != nil {
		fmt.Println("error found after drive instace")
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	if cmd := os.Args[(len(os.Args) - 1)]; cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	} else if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}

	}

}
