package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/sahilsnghai/golang/Project11/internal/types"
	"github.com/sahilsnghai/golang/Project11/internal/utils"
	"github.com/valyala/fasthttp"
)

type APIServer struct {
	ListenAddr string
	Store      types.Storage
	Config     types.Config
	Ctx        *fasthttp.RequestCtx
}

func (s *APIServer) HandleGetMetadata(ctx *fasthttp.RequestCtx) {
	start := time.Now()
	log.Println("Fetching metadata from redis")
	getMetadataReq := new(types.GenericAPIRequest)

	if err := json.Unmarshal(ctx.Request.Body(), getMetadataReq); err != nil {
		utils.WriteJson(ctx, fasthttp.StatusBadRequest, true, map[string]interface{}{
			"error":   true,
			"message": fmt.Sprintf("Invalid request body: %v", err),
		})
		return

	}

	userVo, ok := ctx.UserValue("userVo").(map[string]interface{})
	if !ok {
		utils.WriteJson(ctx, fasthttp.StatusInternalServerError, true, map[string]interface{}{
			"error":   true,
			"message": "userVo not found or invalid in context",
		})
		return

	}

	getMetadataReq.Data["userId"] = userVo["id"]

	metadata, err := s.Store.GetMetadata(getMetadataReq.Data)
	if err != nil {
		utils.WriteJson(ctx, fasthttp.StatusInternalServerError, true, map[string]interface{}{
			"error":   true,
			"message": fmt.Sprintf("Failed to fetch metadata: %v", err),
		})
		return
	}

	log.Printf("total time taken by complete process in HandleGetMetadata is %v ms", time.Since(start).Milliseconds())

	utils.WriteJson(ctx, fasthttp.StatusOK, false, metadata)
}

func (s *APIServer) HandleMigration(ctx *fasthttp.RequestCtx) {
	log.Println("Publishing data from database to redis")

	migrateReq := new(types.GenericAPIRequest)
	if err := json.Unmarshal(ctx.Request.Body(), migrateReq); err != nil {
		utils.WriteJson(ctx, fasthttp.StatusBadRequest, true, map[string]interface{}{
			"error":   true,
			"message": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	userVo, ok := ctx.UserValue("userVo").(map[string]interface{})
	if !ok {
		utils.WriteJson(ctx, fasthttp.StatusInternalServerError, true, map[string]interface{}{
			"error":   true,
			"message": "userVo not found or invalid in context",
		})
		return
	}

	if domainId, ok := migrateReq.Data["domainId"]; !ok || domainId == nil {
		utils.WriteJson(ctx, fasthttp.StatusBadRequest, true, map[string]interface{}{
			"error":   true,
			"message": "domainId is missing or invalid",
		})
		return
	}

	migrateReq.Data["userId"] = userVo["id"]
	migrateReq.Data["orgId"] = userVo["organizationId"]
	migrateReq.Data["email"] = userVo["email"]
	migrateReq.Data["askme_filter_limits"] = s.Config.Parameters["askme_filter_limits"]

	// migrateReq.Data["HEADERS"] = map[string]string{
	// 	"Authorization": string(ctx.Request.Header.Peek("Authorization")),
	// 	"Content-Type":  string(ctx.Request.Header.Peek("Content-Type")),
	// 	"Version":       string(ctx.Request.Header.Peek("Version")),
	// }

	_metadata, err := s.Store.Migration(migrateReq.Data)
	if err != nil {
		log.Printf("Error during migration: %s\n", err.Error())
		utils.WriteJson(ctx, fasthttp.StatusInternalServerError, true, map[string]interface{}{
			"error":   true,
			"message": "Unable to publish metadata",
		})
		return
	}

	if toShow, ok := migrateReq.Data["show-metadata"].(bool); ok && toShow {
		utils.WriteJson(ctx, fasthttp.StatusOK, false, _metadata)
	}

	utils.WriteJson(ctx, fasthttp.StatusOK, false, map[string]interface{}{
		"status": map[string]interface{}{"success": true, "error": nil},
		"sample": 0,
	})

}

func (s *APIServer) HandleHealthCheck(ctx *fasthttp.RequestCtx) {
	utils.WriteJson(ctx, fasthttp.StatusOK, false, "health check")

}
