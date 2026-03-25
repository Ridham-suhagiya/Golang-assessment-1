package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"country-search-api/internal/models"
)

// HTTPClient interface allows us to mock the *http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the interface for interacting with the REST Countries API
type Client interface {
	FetchCountry(ctx context.Context, name string) (*models.Country, error)
}

type restCountriesClient struct {
	httpClient HTTPClient
	baseURL    string
}

// NewRESTCountriesClient creates a new Client for the REST Countries API
func NewRESTCountriesClient(httpClient HTTPClient) Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &restCountriesClient{
		httpClient: httpClient,
		baseURL:    "https://restcountries.com/v3.1",
	}
}

// FetchCountry makes a request to the REST Countries API to fetch country details by name
func (client *restCountriesClient) FetchCountry(context context.Context, name string) (*models.Country, error) {
	// Construct the URL: https://restcountries.com/v3.1/name/{name}
	// url.PathEscape ensures spaces or special characters in the name don't break the URL
	endpoint := fmt.Sprintf("%s/name/%s", client.baseURL, url.PathEscape(name))

	// Create a new HTTP GET request
	request, err := http.NewRequestWithContext(context, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Actually send the HTTP request over the internet
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	
	// Ensure the response body is always closed to prevent memory leaks
	defer response.Body.Close()

	// Check if the API returned a successful 200 OK status
	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Prepare a variable to securely hold exactly what we expect the API to return
	var parsedRows []models.RestCountryResponse
	
	// json.NewDecoder reads the incoming data stream and converts it into our Go structs
	if err := json.NewDecoder(response.Body).Decode(&parsedRows); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(parsedRows) == 0 {
		return nil, fmt.Errorf("country not found")
	}

	firstMatch := parsedRows[0]

	country := &models.Country{
		Name:       firstMatch.Name.Common,
		Population: firstMatch.Population,
	}

	if len(firstMatch.Capital) > 0 {
		country.Capital = firstMatch.Capital[0]
	}

	for _, currency := range firstMatch.Currencies {
		country.Currency = currency.Symbol
		// Just take the first currency's symbol
		break
	}

	return country, nil
}
