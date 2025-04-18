package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisOps struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Db       string `json:"db"`
}

func RedisInstance(redisConfig map[string]interface{}) *redis.Client {

	redisMap := redisConfig["redis"].(map[string]interface{})

	var password string
	password, ok := redisMap["password"].(string)
	if !ok {
		log.Printf("password not found %s\n", password)
	}
	log.Println("Setting up redis ")
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%v", redisMap["host"], redisMap["port"]),
		Password:     password,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  5 * time.Second,
		MinIdleConns: 50,
		MaxIdleConns: 200,
	})

	if pong, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		log.Printf("Redis Ping Response: %s", pong)
	}

	return client

}
