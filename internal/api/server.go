package api

import (
	"fmt"
	"github.com/marcodd23/gopernet/internal/config"
	"net/http"

	"github.com/marcodd23/gopernet/internal/services"
)

func NewServer(service *services.DefaultBurrowService, config *config.ServiceConfig) *http.Server {
	mux := http.NewServeMux()
	RegisterRoutes(mux, service, config)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Server.Port),
		Handler: mux,
	}
}
