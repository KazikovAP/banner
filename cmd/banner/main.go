package main

import (
	"banner/internal/app"
	logerr "banner/internal/lib/logger/logerr"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start server %v", logerr.Err(err))
	}
}
