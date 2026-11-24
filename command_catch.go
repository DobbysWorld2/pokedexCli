package main

import (
	"fmt"
)

const pokemonBaseURL = "https://pokeapi.co/api/v2/pokemon/"

// catchChance returns the probability of catching a Pokemon with the given
// base experience.
//
// Higher base experience means a rarer, harder Pokemon, so the chance falls
// as base experience rises. A Pokemon with no base experience is always
// catchable.
func catchChance(baseExperience int) float64 {
	if baseExperience <= 0 {
		return 1.0
	}

	return 1.0 / (1.0 + float64(baseExperience)/50.0)
}

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: catch <pokemon>")
	}

	name := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := cfg.apiClient.getPokemon(name)
	if err != nil {
		return err
	}

	baseExperience := 0
	if pokemon.BaseExperience != nil {
		baseExperience = *pokemon.BaseExperience
	}

	if cfg.rng.Float64() > catchChance(baseExperience) {
		fmt.Printf("%s escaped!\n", name)
		return nil
	}

	fmt.Printf("%s was caught!\n", name)
	fmt.Println("You may now inspect it with the inspect command.")

	// Add the Pokemon to the user's Pokedex.
	cfg.pokedex[pokemon.Name] = pokemon

	return nil
}
