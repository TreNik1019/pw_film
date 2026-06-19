package models

type Regisseur struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Geburtsjahr int    `json:"geburtsjahr"`
	FilmID      uint   `json:"film_id" gorm:"column:film_id"`
}

func (Regisseur) TableName() string {
	return "regisseur"
}
