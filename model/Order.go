package model

type Order struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Info []byte `gorm:"type:json" json:"info"`
}
