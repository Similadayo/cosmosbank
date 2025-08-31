package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, Cosmos Bank!")

	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Run() error {
	// TODO: Wire CosmosBank app + CLI commands
	return nil
}
