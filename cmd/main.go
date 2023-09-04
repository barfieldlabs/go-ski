package main

import (
	"context"
	"log"

	"github.com/barfieldlabs/go-ski/core"
	"github.com/chromedp/cdproto/target"
)

func main() {
	// Initialize core procedures
	proc := core.NewProcedures()

	// Create a context
	ctx := context.Background()

	// Initialize or fetch the initial targets
	var initialTargets []*target.Info

	// Perform actions
	err := proc.Execute(ctx, initialTargets)
	if err != nil {
		log.Fatalf("Failed to perform actions: %v", err)
	}

	log.Println("Successfully completed web scraping.")
}
