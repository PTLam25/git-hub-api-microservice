package app

import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/log/log_logrus"
	"github.com/PTLam25/git-hub-api-microservice/src/api/log/log_zap"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartApp() {
	log_logrus.Info("about to map the urls", "step:1", "status:pending")
	log_zap.Info("about to map the urls", log_zap.Field("step", 1), log_zap.Field("status", "pending"))
	// выводим логику запуска в отдельную функцию, чтобы можно было тестировать запуск приложения
	mapUrls()
	log_logrus.Info("the urls successfully mapped", "step:2", "status:success")
	log_zap.Info("about to map the urls", log_zap.Field("step", 2), log_zap.Field("status", "success"))

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
