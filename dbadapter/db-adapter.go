package dbadapter

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/serafimcode/wb-test-L0/model"
)

var dbInstance *gorm.DB

func InitDb() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&model.Order{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")
	dbInstance = db
}

func GetDb() *gorm.DB {
	if dbInstance == nil {
		panic("Db was not initialized")
	}
	return dbInstance
}
