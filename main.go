package main

import (
	"fmt"
	"log"
	"os"

	"webextractor/internal/app"
	"webextractor/internal/cli"
	"webextractor/internal/types"
)

func main() {
	flags, err := cli.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ %s\n", err)
		os.Exit(1)
	}

	config := types.NewExtractionConfig(flags.URL.String(), flags.Sel, flags.Out.String(), flags.Timeout)
	
	application := app.New(config)
	if err := application.Run(); err != nil {
		log.Fatalf("❌ %v", err)
	}
}