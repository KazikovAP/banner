package main

import (
	"banner/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}
