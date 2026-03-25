package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"country-search-api/internal/service"
)

// CountryHandler is the HTTP handler for country endpoints
type CountryHandler struct {
	svc service.CountryService
}

// NewCountryHandler creates a new CountryHandler
func NewCountryHandler(svc service.CountryService) *CountryHandler {
	return &CountryHandler{
		svc: svc,
	}
}

// Search handles the GET /api/countries/search?name={name} endpoint
func (handler *CountryHandler) Search(responseWriter http.ResponseWriter, request *http.Request) {
	// 1. Verify that the incoming HTTP method is exactly what we expect (GET)
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Extract the 'name' query parameter from the URL
	name := request.URL.Query().Get("name")
	if name == "" {
		http.Error(responseWriter, "Query parameter 'name' is required", http.StatusBadRequest)
		return
	}

	// 3. Call the business logic layer to find the country
	country, err := handler.svc.SearchCountry(request.Context(), name)
	if err != nil {
		// If the error message indicates the country doesn't exist, return a 404 Not Found error
		if strings.Contains(err.Error(), "country not found") {
			http.Error(responseWriter, err.Error(), http.StatusNotFound)
			return
		}
		
		// If the error is something unexpected, return a generic 500 Server Error
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Set the response content type to JSON so the client knows how to parse it
	responseWriter.Header().Set("Content-Type", "application/json")
	
	// Send a 200 OK status code indicating success
	responseWriter.WriteHeader(http.StatusOK)

	// 5. Convert the 'country' struct into JSON and write it directly to the response output stream
	if err := json.NewEncoder(responseWriter).Encode(country); err != nil {
		http.Error(responseWriter, "Failed to encode response", http.StatusInternalServerError)
	}
}
