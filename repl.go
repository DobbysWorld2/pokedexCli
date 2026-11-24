package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// config holds state shared by the REPL and its commands.
type config struct {
	// commands is an ordered slice so that `help` prints a stable list.
	commands []cliCommand

	// registry maps a command name to its command definition.
	registry map[string]cliCommand

	// apiClient talks to the PokeAPI.
	apiClient *apiClient

	// pokedex holds the Pokemon the user has caught, keyed by name.
	pokedex map[string]Pokemon

	// rng generates the random numbers used for catch chances.
	rng *rand.Rand

	// Pagination URLs for the location-area endpoint.
	nextLocationsURL     string
	previousLocationsURL string
}

// pokedexNames returns the names of all caught Pokemon, sorted alphabetically.
func (cfg *config) pokedexNames() []string {
	names := make([]string, 0, len(cfg.pokedex))

	for name := range cfg.pokedex {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		userInput := scanner.Text()

		cleanedInput := cleanInput(userInput)
		if len(cleanedInput) == 0 {
			continue
		}

		commandName := cleanedInput[0]
		args := cleanedInput[1:]

		cmd, ok := cfg.registry[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(cfg, args); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
