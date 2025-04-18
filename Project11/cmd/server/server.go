package server

import (
	"fmt"
	"log"

	"github.com/sahilsnghai/golang/Project11/cmd/handlers"
	"github.com/sahilsnghai/golang/Project11/internal/types"
	"github.com/sahilsnghai/golang/Project11/internal/utils"
	"github.com/valyala/fasthttp"
)

type APIServer struct {
	handlers.APIServer
}

func (s *APIServer) Run() {
	routes := map[string]map[string]func(ctx *fasthttp.RequestCtx){
		"/get-metadata": {
			"POST": utils.MiddleWare(s.HandleGetMetadata),
		},
		"/data-migration": {
			"POST": utils.MiddleWare(s.HandleMigration),
		},
		"/health": {
			"GET": utils.MiddleWare(s.HandleHealthCheck),
		},
	}

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		method := string(ctx.Method())

		if handlers, exists := routes[string(ctx.Path())]; exists {
			if __handler, methodExists := handlers[method]; methodExists {
				s.Store.UpdateReqCtx(ctx)
				__handler(ctx)

			} else {
				utils.WriteJson(ctx, fasthttp.StatusMethodNotAllowed, true, map[string]string{"message": fmt.Sprintf("Method %s Not Allowed", method)})
			}
		} else {
			utils.WriteJson(ctx, fasthttp.StatusNotFound, true, "Route Not Found")
		}
	}

	log.Printf("Starting Server at port: %s\n", s.ListenAddr)
	if err := fasthttp.ListenAndServe(s.ListenAddr, requestHandler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

func NewAPIServer(listenAddr string, store types.Storage, config types.Config) *APIServer {
	return &APIServer{
		handlers.APIServer{
			ListenAddr: listenAddr,
			Store:      store,
			Config:     config,
			Ctx:        nil,
		}}
}
