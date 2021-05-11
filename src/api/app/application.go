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
	log.Log.Info("about to map the urls")
	// выводим логику запуска в отдельную функцию, чтобы можно было тестировать запуск приложения
	mapUrls()
	log.Log.Info("the urls successfully mapped")

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
