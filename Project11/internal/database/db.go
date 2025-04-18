package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sahilsnghai/golang/Project11/internal/cache"
	"github.com/sahilsnghai/golang/Project11/internal/services"
	"github.com/sahilsnghai/golang/Project11/internal/types"
)

func NewDatabaseStore(database types.DatabaseConfig, parameters map[string]interface{}) (*services.Storage, error) {
	cc := database.Ccplatform
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cc.UserName,
		cc.Password,
		cc.Host,
		cc.Port,
		cc.Name,
	)

	log.Println("Connection string created successfully, connecting to DB...")

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db.SetMaxIdleConns(40)
	db.SetMaxOpenConns(1000)
	db.SetConnMaxLifetime(10 * time.Minute)

	// db2, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to database: %v", err)
	// }

	// sqlDB, err := db2.DB()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to retrieve raw DB from GORM: %v", err)
	// }

	// sqlDB.SetMaxIdleConns(10)
	// sqlDB.SetMaxOpenConns(100)
	// sqlDB.SetConnMaxLifetime(time.Minute * 5)

	client := cache.RedisInstance(parameters)

	log.Println("Database connected successfully")

	return &services.Storage{
		Db:         db,
		Parameters: parameters,
		Client:     client,
		Ctx:        context.Background(),
		Headers:    nil,
	}, nil

}
