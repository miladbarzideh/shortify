package main

import (
	"github.com/miladbarzideh/shortify/internal/app"
	"github.com/miladbarzideh/shortify/internal/infra/logger"
)

func main() {
	log := logger.InitLogger()
	server := app.NewServer(log)
	if server.Run() != nil {
		log.Fatal("failed to start the app")
	}
}
