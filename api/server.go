package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/serafimcode/wb-test-L0/dbadapter"
	"github.com/serafimcode/wb-test-L0/memorycache"
	"github.com/serafimcode/wb-test-L0/repository"
)

var localCache *memorycache.Cache

func InitServer(cache *memorycache.Cache) {
	r := gin.Default()
	localCache = cache

	r.GET("/orders/:id", getOrderById)
	r.Run()
}

func getOrderById(c *gin.Context) {
	id := c.Param("id")

	orderInfo, isExist := localCache.Get(id)

	if !isExist {
		repo := repository.OrderRepository{Db: dbadapter.GetDb()}
		order := *repo.GetById(id)

		localCache.Set(id, order.Info, 0)
		orderInfo = order.Info
	}

	c.JSON(http.StatusOK, orderInfo)
}
