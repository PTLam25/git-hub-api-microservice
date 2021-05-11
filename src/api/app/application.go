package app

import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/log"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartApp() {
	log.Info("about to map the urls", "step:1", "status:pending")
	// выводим логику запуска в отдельную функцию, чтобы можно было тестировать запуск приложения
	mapUrls()
	log.Info("the urls successfully mapped", "step:2", "status:success")

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
