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

	listenAddr := ":8080"

	s := server.NewAPIServer(listenAddr, store, *config)
	// fmt.Printf("\nconfig parameters %+v \n", config)
	// "metadata", "aDD", "aDF", "inactive",
	s.Run()
	// req := interface{}{"domainId": 27074, "Keys": []string{"aDD", "aDF", "inactive", "shortcuts"}, UserId: 10910}
	// req := map[string]interface{}{
	// 	"domainId": float64(27074),
	// 	"keys":     []interface{}{"metadata", "aDD"},
	// }
	// redis := RedisInstance(s.config.Parameters)
	// defer redis.Close()
	// metadata, err := s.store.GetMetadata(req, redis)
	// if err != nil {
	// 	log.Fatalf("unable to connect to database: %s", err.Error())
	// 	return
	// }

	// log.Println(redis.Ping(ctx))
	// fmt.Printf("metadata : %+v\n", metadata)
}
