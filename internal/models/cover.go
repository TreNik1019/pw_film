package models

type Cover struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Info        string `json:"info"`        // z.B. "Cover 1"
	ContentType string `json:"contentType"` // z.B. "img/png"
	FlmID       uint   `json:"flm_id" gorm:"column:flm_id"`
}

func (Cover) TableName() string {
	return "cover"
}
