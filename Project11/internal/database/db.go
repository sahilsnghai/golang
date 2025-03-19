package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(3 * time.Minute)

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

	log.Println("Database connected successfully")
	return &services.Storage{Db: db, Parameters: parameters, Db2: nil, Client: nil, Ctx: nil}, nil
}
