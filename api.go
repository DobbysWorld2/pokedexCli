package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"pokedexcli/internal/pokecache"
)

// ---------------------------------------------------------------- response types

// NamedResource is the shape PokeAPI uses for a name plus a URL to fetch more
// detail from.
type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// locationAreas is one page of the location-area list endpoint.
type locationAreas struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []locationArea `json:"results"`
}

// locationArea is a single entry in a page of location areas.
type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// locationAreaDetail is a single location area, including the Pokemon that
// can be encountered there.
type locationAreaDetail struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	GameIndex         int                `json:"game_index"`
	Locations         []NamedResource    `json:"locations"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

// PokemonEncounter is one Pokemon that can be encountered in a location area.
type PokemonEncounter struct {
	Pokemon        NamedResource `json:"pokemon"`
	VersionDetails []struct {
		EncounterDetails []struct {
			Chance          int              `json:"chance"`
			ConditionValues []NamedResource  `json:"condition_values"`
		} `json:"encounter_details"`
		Version NamedResource `json:"version"`
	} `json:"version_details"`
}

// ---------------------------------------------------------------- client

type apiClient struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

func newAPIClient() *apiClient {
	return &apiClient{
		httpClient: http.Client{
			Timeout: 10 * time.Second,
		},
		cache: pokecache.NewCache(5 * time.Minute),
	}
}

// getLocationAreas fetches a page of location areas, serving it from the
// cache when possible.
func (c *apiClient) getLocationAreas(url string) (locationAreas, error) {
	body, err := c.fetch(url)
	if err != nil {
		return locationAreas{}, err
	}

	var data locationAreas
	if err := json.Unmarshal(body, &data); err != nil {
		return locationAreas{}, fmt.Errorf("could not unmarshal response: %w", err)
	}

	return data, nil
}

// getLocationArea fetches the details of a single location area, serving it
// from the cache when possible.
func (c *apiClient) getLocationArea(url string) (locationAreaDetail, error) {
	body, err := c.fetch(url)
	if err != nil {
		return locationAreaDetail{}, err
	}

	var data locationAreaDetail
	if err := json.Unmarshal(body, &data); err != nil {
		return locationAreaDetail{}, fmt.Errorf("could not unmarshal response: %w", err)
	}

	return data, nil
}

// fetch returns the raw response bytes for url, using the cache when the
// entry is still present. Network failures and invalid responses are never
// cached.
func (c *apiClient) fetch(url string) ([]byte, error) {
	// Cache hit: no request needed.
	if val, ok := c.cache.Get(url); ok {
		return val, nil
	}

	// Cache miss: make the request.
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	// Verify the body is valid JSON before caching it, so that a malformed
	// response is never stored.
	if !json.Valid(body) {
		return nil, fmt.Errorf("invalid JSON in response")
	}

	c.cache.Add(url, body)

	return body, nil
}

// Pokemon is a single Pokemon, as returned by the pokemon endpoint.
type Pokemon struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	BaseExperience *int             `json:"base_experience"`
	Height         int              `json:"height"`
	Weight         int              `json:"weight"`
	Types          []NamedResource  `json:"types"`
	Stats          []PokemonStat    `json:"stats"`
}

// PokemonStat is one of a Pokemon's six stats.
type PokemonStat struct {
	BaseStat int            `json:"base_stat"`
	Stat     NamedResource  `json:"stat"`
}

// getPokemon fetches a single Pokemon by name, serving it from the cache
// when possible.
func (c *apiClient) getPokemon(name string) (Pokemon, error) {
	pokemonURL := pokemonBaseURL + url.PathEscape(name)

	body, err := c.fetch(pokemonURL)
	if err != nil {
		return Pokemon{}, err
	}

	var data Pokemon
	if err := json.Unmarshal(body, &data); err != nil {
		return Pokemon{}, fmt.Errorf("could not unmarshal response: %w", err)
	}

	return data, nil
}
