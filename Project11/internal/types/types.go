package types

import (
	"sync"

	"github.com/sahilsnghai/golang/Project11/internal/services"
	"github.com/valyala/fasthttp"
)

type APIError struct {
	Error string `json:"error"`
}

// type APIFunc func(http.ResponseWriter, *http.Request) error
type APIFunc func(*fasthttp.RequestCtx)

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
	GetMetadata(map[string]interface{}) (map[string]interface{}, error)
	Migration(map[string]interface{}) (map[string]interface{}, error)
	Migrate(map[string]interface{}, *sync.WaitGroup) (map[string]interface{}, error)
	UpdateReqCtx(*fasthttp.RequestCtx)
}

type MysqlStorage struct {
	services.Storage
}
