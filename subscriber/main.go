package main

import (
	"log"
	"time"

	env "github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/serafimcode/wb-test-L0/api"
	"github.com/serafimcode/wb-test-L0/dbadapter"
	"github.com/serafimcode/wb-test-L0/memorycache"
	"github.com/serafimcode/wb-test-L0/model"
	"github.com/serafimcode/wb-test-L0/stanclient"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatal(err)
	}

	dbadapter.InitDb()

	cache := memorycache.New(5*time.Minute, 10*time.Minute)

	populateCache(dbadapter.GetDb(), cache)

	stanService := stanclient.InitStan()
	defer stanService.Connection.Close()

	api.InitServer(cache)
}

func populateCache(db *gorm.DB, cache *memorycache.Cache) {
	var orders []model.Order
	db.Where("created_at >= ? ", time.Now().Add(-time.Hour)).Find(&orders)

	for _, order := range orders {
		cache.Set(order.ID, order.Info, 0)
	}

	cache.Log()
}
