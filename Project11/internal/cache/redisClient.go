package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisOps struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Db       string `json:"db"`
}

// Function to create a Redis client
func RedisInstance(redisConfig map[string]interface{}) *redis.Client {

	redisMap := redisConfig["redis"].(map[string]interface{})

	var password string
	password, ok := redisMap["password"].(string)
	if !ok {
		log.Printf("password not found %s\n", password)
	}
	// log.Printf("%s:%v\n", redisMap["host"], redisMap["port"])

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", redisMap["host"], redisMap["port"]),
		Password: password,
	})

	log.Println(client.Ping(context.Background()))

	return client

}
