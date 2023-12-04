package repository

import (
	"log"

	"gorm.io/gorm"

	"github.com/serafimcode/wb-test-L0/model"
)

type OrderRepository struct {
	Db *gorm.DB
}

func (repository *OrderRepository) Create(order *model.Order) {
	res := repository.Db.Create(&order)
	if res.Error != nil {
		log.Println(res.Error)
	}
}

func (repository *OrderRepository) GetById(id string) *model.Order {
	var order model.Order
	res := repository.Db.Find(&order, "id = ?", id)

	if res.Error != nil {
		log.Println(res.Error)
	}

	return &order
}
