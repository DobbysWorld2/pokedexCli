package main

import "fmt"

func commandMapBack(cfg *config, args []string) error {
	if cfg.previousLocationsURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	data, err := cfg.apiClient.getLocationAreas(cfg.previousLocationsURL)
	if err != nil {
		return err
	}

	for _, area := range data.Results {
		fmt.Println(area.Name)
	}

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
