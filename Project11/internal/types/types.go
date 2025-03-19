package types

import (
	"context"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/sahilsnghai/golang/Project11/internal/services"
)

var ctx = context.Background()

type APIError struct {
	Error string `json:"error"`
}

type APIFunc func(http.ResponseWriter, *http.Request) error

type DatabaseConfig struct {
	Ccplatform struct {
		Driver   string
		Host     string
		Name     string
		Password string
		Port     int
		Schema   string
		UserName string
	} `json:"ccplatform"`
}
type Config struct {
	Database   DatabaseConfig         `json:"dataSources"`
	Parameters map[string]interface{} `json:"parameters"`
}

type GenericAPIRequest struct {
	Data map[string]interface{} `json:"data"`
}

type Storage interface {
	GetMetadata(map[string]interface{}, *redis.Client) (map[string]interface{}, error)
	Migration(map[string]interface{}, *redis.Client) (map[string]interface{}, error)
	Migrate(map[string]interface{}) (map[string]interface{}, error)
}

type MysqlStorage struct {
	services.Storage
}
