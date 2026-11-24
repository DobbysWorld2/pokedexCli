package main

import (
	"fmt"
	"net/url"
)

func commandExplore(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: explore <location_area>")
	}

	areaName := args[0]

	fmt.Printf("Exploring %s...\n", areaName)

	areaURL := locationAreaURL + url.PathEscape(areaName)

	data, err := cfg.apiClient.getLocationArea(areaURL)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}
