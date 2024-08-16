package api

import (
	"encoding/json"
	"net/http"

	"github.com/marcodd23/gopernet/internal/services"
)

// JSONResponse defines a structure for consistent API responses.
type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// GetBurrowsHandler returns the list of burrows.
func GetBurrowsHandler(service *services.DefaultBurrowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		burrows := service.GetAllBurrows()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONResponse{
			Status: "success",
			Data:   burrows,
		})
	}
}

// RentBurrowHandler allows a burrow to be rented.
func RentBurrowHandler(service *services.DefaultBurrowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := service.RentBurrow(request.Name)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JSONResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		// Respond with success and the name of the rented burrow
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JSONResponse{
			Status:  "success",
			Message: "Burrow rented successfully",
			Data:    map[string]string{"name": request.Name},
		})
	}
}

// GenerateReportHandler generates a report of the current state of the burrows.
func GenerateReportHandler(service *services.DefaultBurrowService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Generate the report
		report, err := service.GenerateReport()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSONResponse{
				Status:  "error",
				Message: "Failed to generate report",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JSONResponse{
			Status:  "success",
			Message: "Report generated successfully",
			Data:    report,
		})
	}
}
