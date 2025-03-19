package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sahilsnghai/golang/Project11/internal/cache"
	"github.com/sahilsnghai/golang/Project11/internal/types"
	"github.com/sahilsnghai/golang/Project11/internal/utils"
)

type APIServer struct {
	ListenAddr string
	Store      types.Storage
	Config     types.Config
}

func (s *APIServer) HandleGetMetadata(w http.ResponseWriter, r *http.Request) error {

	start := time.Now()
	log.Println("Fetching metadata from redis")
	getMetadataReq := new(types.GenericAPIRequest)
	if err := json.NewDecoder(r.Body).Decode(getMetadataReq); err != nil {
		return err
	}

	userVo, ok := r.Context().Value("userVo").(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid userVo in context", http.StatusInternalServerError)
		return fmt.Errorf("userVo not found or invalid in context")
	}

	getMetadataReq.Data["userId"] = userVo["id"]

	client := cache.RedisInstance(s.Config.Parameters)
	defer client.Close()

	metadata, err := s.Store.GetMetadata(getMetadataReq.Data, client)
	if err != nil {
		return err
	}

	log.Printf("total time taken by complete process in HandleGetMetadata is %v ms", time.Since(start).Milliseconds())

	return utils.WriteJson(w, http.StatusOK, false, metadata)

}

func (s *APIServer) HandleMigration(w http.ResponseWriter, r *http.Request) error {

	log.Println("Publishing data from database to redis")
	migrateReq := new(types.GenericAPIRequest)
	if err := json.NewDecoder(r.Body).Decode(migrateReq); err != nil {
		return err
	}

	userVo, ok := r.Context().Value("userVo").(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid userVo in context", http.StatusInternalServerError)
		return fmt.Errorf("userVo not found or invalid in context")
	}
	migrateReq.Data["userId"] = userVo["id"]
	migrateReq.Data["orgId"] = userVo["organizationId"]
	migrateReq.Data["email"] = userVo["email"]
	migrateReq.Data["askme_filter_limits"] = s.Config.Parameters["askme_filter_limits"]

	client := cache.RedisInstance(s.Config.Parameters)
	defer client.Close()

	_metadata, err := s.Store.Migration(migrateReq.Data, client)
	if err != nil {
		http.Error(w, "not able to publish metadata", http.StatusInternalServerError)
		log.Printf("error while migration %s\n", err.Error())
		return fmt.Errorf("unable to publish")
	}

	if toShow, ok := migrateReq.Data["show-metadata"].(bool); ok && toShow {
		return utils.WriteJson(w, http.StatusOK, false, _metadata)
	}
	return utils.WriteJson(w, http.StatusOK, true, fmt.Sprint("DataMigration called"))

}

func (s *APIServer) HandleHealthCheck(w http.ResponseWriter, r *http.Request) error {
	return utils.WriteJson(w, http.StatusOK, false, fmt.Sprint("health check "))

}
