package main

import "fmt"

const locationAreaURL = "https://pokeapi.co/api/v2/location-area/"

func commandMap(cfg *config, args []string) error {
	url := locationAreaURL
	if cfg.nextLocationsURL != "" {
		url = cfg.nextLocationsURL
	}

	data, err := cfg.apiClient.getLocationAreas(url)
	if err != nil {
		return err
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

	// Store the pagination URLs for the next page and for `mapb`.
	cfg.nextLocationsURL = ""
	if data.Next != nil {
		cfg.nextLocationsURL = *data.Next
	}

	cfg.previousLocationsURL = ""
	if data.Previous != nil {
		cfg.previousLocationsURL = *data.Previous
	}

	return nil
}
