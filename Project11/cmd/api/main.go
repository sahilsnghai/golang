package main

import (
	"log"

	"github.com/sahilsnghai/golang/Project11/cmd/server"
	db "github.com/sahilsnghai/golang/Project11/internal/database"
	"github.com/sahilsnghai/golang/Project11/internal/utils"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config, err := utils.GetConfig()
	if err != nil {
		log.Fatalf("unable to read configuration file: %s", err.Error())
		return
	}

	store, err := db.NewDatabaseStore(config.Database, config.Parameters)
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err.Error())
		return
	}
	defer store.Client.Close()

	listenAddr := ":8080"

	s := server.NewAPIServer(listenAddr, store, *config)
	s.Run()
}
