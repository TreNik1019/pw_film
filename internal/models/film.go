package models

import "time"

type Film struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Version           int       `json:"version"`
	Titel             string    `json:"titel"`
	Art               string    `json:"art"` // z.B. BlueRay, 4K
	Erscheinungsdatum time.Time `json:"erscheinungsdatum" gorm:"type:date"`
	Genre             string    `json:"genre"`
	Rating            float64   `json:"rating"`
	Verfuegbar        bool      `json:"verfuegbar"`
	Preis             float64   `json:"preis"`
	Schlagwoerter     string    `json:"schlagwoerter"` // Wird als JSON-String in der CSV gehalten (z.B. "["AUFREGEND","SPANNEND"]")
	Erzeugt           time.Time `json:"erzeugt"`
	Aktualisiert      time.Time `json:"aktualisiert"`
}

// TableName überschreibt den Standard-Plural von GORM
func (Film) TableName() string {
	return "film"
}
