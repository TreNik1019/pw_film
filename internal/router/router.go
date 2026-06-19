package router

import (
	"pw_film/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter registers middleware and maps handlers to endpoints
func SetupRouter(filmHandler *handlers.FilmHandler) *gin.Engine {
	r := gin.Default()

	// Setup basic CORS middleware
	r.Use(corsMiddleware())

	// Versioned route group
	v1 := r.Group("/api/v1")
	{
		v1.POST("/films", filmHandler.CreateFilm)
		v1.GET("/films", filmHandler.GetFilms)
		v1.GET("/films/:id", filmHandler.GetFilmByID)
	}

	return r
}

// CORS middleware allowing generic development integration
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
