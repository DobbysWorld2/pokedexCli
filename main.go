package main

import (
	"math/rand"
	"time"
)

func main() {
	cfg := config{
		apiClient: newAPIClient(),
		pokedex:   make(map[string]Pokemon),
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		commands: []cliCommand{
			{
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			{
				name:        "map",
				description: "Display 20 location areas in the Pokemon world. Subsequent calls display the next 20.",
				callback:    commandMap,
			},
			{
				name:        "mapb",
				description: "Display the previous 20 location areas in the Pokemon world.",
				callback:    commandMapBack,
			},
			{
				name:        "explore",
				description: "Explore a location area and list the Pokemon found there.",
				callback:    commandExplore,
			},
			{
				name:        "catch",
				description: "Attempt to catch a Pokemon.",
				callback:    commandCatch,
			},
			{
				name:        "pokedex",
				description: "List all the Pokemon you've caught.",
				callback:    commandPokedex,
			},
			{
				name:        "inspect",
				description: "View details of a Pokemon you've caught.",
				callback:    commandInspect,
			},

			{
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
						

		},
	}

	cfg.registry = make(map[string]cliCommand, len(cfg.commands))
	for _, cmd := range cfg.commands {
		cfg.registry[cmd.name] = cmd
	}

	startRepl(&cfg)
}
