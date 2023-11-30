package repository

import (
	"gorm.io/gorm"

	"github.com/serafimcode/wb-test-L0/model"
)

type OrderRepository struct {
	Db *gorm.DB
}

func (repository *OrderRepository) Create(order model.Order) {
	repository.Db.Create(order)
}

func (repository *OrderRepository) GetById(id string) *model.Order {
	var order model.Order
	repository.Db.Find(&order, "id = ?", id)
	return &order
}
