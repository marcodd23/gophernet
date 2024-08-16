package api

import (
	"github.com/marcodd23/gopernet/internal/config"
	"net/http"

	"github.com/marcodd23/gopernet/internal/services"
)

func RegisterRoutes(mux *http.ServeMux, service *services.DefaultBurrowService, config *config.ServiceConfig) {
	mux.HandleFunc(config.Rest.Endpoints["get-burrows"].Path, GetBurrowsHandler(service))
	mux.HandleFunc(config.Rest.Endpoints["rent-burrow"].Path, RentBurrowHandler(service))
	mux.HandleFunc(config.Rest.Endpoints["get-report"].Path, GenerateReportHandler(service))
}
