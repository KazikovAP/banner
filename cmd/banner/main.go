package main

import (
	"banner/internal/app"
	"banner/internal/lib/logger"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start server %v", logger.Err(err))
	}
}
