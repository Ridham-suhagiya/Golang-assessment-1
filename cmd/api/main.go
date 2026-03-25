package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"country-search-api/internal/cache"
	"country-search-api/internal/client"
	"country-search-api/internal/handler"
	"country-search-api/internal/service"
)

func main() {
	// 1. Initialize the components for our app
	fmt.Println("Initializing components...")

	// Create an in-memory cache to store country data for faster repeated access
	inMemoryCache := cache.NewMemoryCache()

	// Set up the HTTP client that will be used to make outgoing requests
	httpClient := &http.Client{Timeout: 10 * time.Second}
	// Create the REST Countries API client
	restClient := client.NewRESTCountriesClient(httpClient)

	// Create the service layer to handle the business logic (cache + external API)
	countryService := service.NewCountryService(inMemoryCache, restClient)

	// Create the handler to connect HTTP input to our service
	countryHandler := handler.NewCountryHandler(countryService)

	// 2. Set up the HTTP router
	mux := http.NewServeMux()
	
	// Register the /api/countries/search endpoint to our handler's Search function
	mux.HandleFunc("/api/countries/search", countryHandler.Search)

	// 3. Determine the port for the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if none is provided via environment variables
	}

	// 4. Start the backend HTTP server
	fmt.Printf("Server is starting on port %s...\n", port)
	
	// ListenAndServe blocks forever, keeping our server alive and listening for incoming requests
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Could not listen on %s: %v\n", port, err)
		os.Exit(1)
	}
}
