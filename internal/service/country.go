package service

import (
	"context"
	"fmt"
	"strings"

	"country-search-api/internal/cache"
	"country-search-api/internal/client"
	"country-search-api/internal/models"
)

// CountryService defines the interface for the country service operations
type CountryService interface {
	SearchCountry(ctx context.Context, name string) (*models.Country, error)
}

type countryService struct {
	cache  cache.Cache
	client client.Client
}

// NewCountryService creates a new instance of CountryService
func NewCountryService(cache cache.Cache, client client.Client) CountryService {
	return &countryService{
		cache:  cache,
		client: client,
	}
}

// SearchCountry searches for a country first in the cache, and if not found, via the HTTP client
func (service *countryService) SearchCountry(context context.Context, name string) (*models.Country, error) {
	// Normalize the name to handle cache keys consistently (e.g., "  CaNaDa " -> "canada")
	cacheKey := strings.ToLower(strings.TrimSpace(name))

	// 1. Check cache: Ask the memory cache if we already have this country
	if val, found := service.cache.Get(cacheKey); found {
		fmt.Printf("Cache HIT for key: %s\n", cacheKey)
		
		// Attempt to safely cast the cached value to our Country struct
		if country, ok := val.(*models.Country); ok {
			return country, nil
		}
		fmt.Printf("Cache cast error for key: %s, continuing to fetch\n", cacheKey)
	}

	fmt.Printf("Cache MISS for key: %s\n", cacheKey)

	// 2. Fetch from client: Since it wasn't in the cache, ask our REST client to get it from the Internet
	country, err := service.client.FetchCountry(context, name)
	if err != nil {
		return nil, fmt.Errorf("service failed to fetch country: %w", err)
	}

	// 3. Store in cache: Save the new country in the cache so we don't have to fetch it again later
	service.cache.Set(cacheKey, country)
	fmt.Printf("Stored in Cache for key: %s\n", cacheKey)

	// 4. Return result: Give the country back to the handler
	return country, nil
}
