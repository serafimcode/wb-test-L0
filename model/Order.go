package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID        string `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	Info      json.RawMessage `gorm:"type:json" json:"info"`
}

func (m *Order) BeforeCreate(*gorm.DB) (err error) {
	m.CreatedAt = time.Now().UTC()
	return nil
}
