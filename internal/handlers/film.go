package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"pw_film/internal/models"
	"pw_film/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FilmHandler mediates between the HTTP requests and the Repository layer
type FilmHandler struct {
	repo repository.FilmRepository
}

// NewFilmHandler constructs a new FilmHandler
func NewFilmHandler(repo repository.FilmRepository) *FilmHandler {
	return &FilmHandler{repo: repo}
}

// CreateFilm handles POST /api/v1/films
func (h *FilmHandler) CreateFilm(c *gin.Context) {
	var input models.Film
	// Gin automatically validates structural requirements using the model validation tags
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create film: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GetFilms handles GET /api/v1/films
func (h *FilmHandler) GetFilms(c *gin.Context) {
	films, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch films: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, films)
}

// GetFilmByID handles GET /api/v1/films/:id
func (h *FilmHandler) GetFilmByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid film ID format"})
		return
	}

	film, err := h.repo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Film not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve film: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, film)
}
