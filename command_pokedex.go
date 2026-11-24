package main

import (
	"fmt"
)

func commandPokedex(cfg *config, args []string) error {
	fmt.Println("Your Pokedex:")

	for _, name := range cfg.pokedexNames() {
		fmt.Printf(" - %s\n", name)
	}

	return nil
}
