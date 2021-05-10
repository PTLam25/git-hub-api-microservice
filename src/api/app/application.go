package app

import (
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartApp() {
	// выводим логику запуска в отдельную функцию, чтобы можно было тестировать запуск приложения
	mapUrls()

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
