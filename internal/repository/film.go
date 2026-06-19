package repository

import (
	"pw_film/internal/models"

	"gorm.io/gorm"
)

// FilmRepository defines database contract for managing movies
type FilmRepository interface {
	Create(film *models.Film) error
	GetAll() ([]models.Film, error)
	GetByID(id uint) (*models.Film, error)
}

type gormFilmRepository struct {
	db *gorm.DB
}

// NewFilmRepository constructs a GORM film repository
func NewFilmRepository(db *gorm.DB) FilmRepository {
	return &gormFilmRepository{db: db}
}

// Create adds a new film to the database
func (r *gormFilmRepository) Create(film *models.Film) error {
	return r.db.Create(film).Error
}

// GetAll returns all films from the database
func (r *gormFilmRepository) GetAll() ([]models.Film, error) {
	var films []models.Film
	err := r.db.Find(&films).Error
	return films, err
}

// GetByID finds a film by its primary key
func (r *gormFilmRepository) GetByID(id uint) (*models.Film, error) {
	var film models.Film
	if err := r.db.First(&film, id).Error; err != nil {
		return nil, err
	}
	return &film, nil
}
