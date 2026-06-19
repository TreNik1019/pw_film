package main

import (
	"log"

	"pw_film/internal/config"
	"pw_film/internal/database"
	"pw_film/internal/handlers"
	"pw_film/internal/repository"
	"pw_film/internal/router"
)

func main() {
	log.Println("Initializing pw_film server...")

	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Setup database connection and run auto-migration
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// 3. Initialize repository
	filmRepo := repository.NewFilmRepository(db)

	// 4. Initialize handlers
	filmHandler := handlers.NewFilmHandler(filmRepo)

	// 5. Setup Router
	r := router.SetupRouter(filmHandler)

	// 6. Start server
	log.Printf("Starting HTTP server on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
