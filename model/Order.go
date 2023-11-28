package model

type Order struct {
	ID   string                 `gorm:"primaryKey" json:"id"`
	Info map[string]interface{} `gorm:"type:json" json:"info"`
}
