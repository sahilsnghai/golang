package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sahilsnghai/golang/Project11/cmd/handlers"
	"github.com/sahilsnghai/golang/Project11/internal/types"
	"github.com/sahilsnghai/golang/Project11/internal/utils"
)

type APIServer struct {
	handlers.APIServer
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/get-metadata", utils.MiddleWare(s.dispatchMethod(
		map[string]types.APIFunc{
			"POST": s.HandleGetMetadata,
		},
	)))

	router.HandleFunc("/data-migration", utils.MiddleWare(s.dispatchMethod(
		map[string]types.APIFunc{
			"POST": s.HandleMigration,
		},
	)))

	router.HandleFunc("/health", utils.MiddleWare(s.dispatchMethod(
		map[string]types.APIFunc{
			"GET": s.HandleHealthCheck,
		},
	)))

	log.Printf("Starting Server at port: %s\n", s.ListenAddr)
	http.ListenAndServe(s.ListenAddr, router)
}

func NewAPIServer(listenAddr string, store types.Storage, config types.Config) *APIServer {
	return &APIServer{handlers.APIServer{ListenAddr: listenAddr, Store: store, Config: config}}
}

func (s *APIServer) dispatchMethod(methods map[string]types.APIFunc) types.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("Redirecting to methods")
		if handler, ok := methods[r.Method]; ok {
			return handler(w, r)
		}
		return utils.WriteJson(w, http.StatusMethodNotAllowed, true, fmt.Sprint("Method Not Allowed: ", r.Method))
	}
}
