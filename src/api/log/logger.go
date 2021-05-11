package log

import (
	"github.com/PTLam25/git-hub-api-microservice/src/api/config"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Log *logrus.Logger
)

func init() {
	// получаем log level из конфига и создаем объект logrud Level
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}

	// создаем логгер
	Log = &logrus.Logger{
		Level: level,
		Out:   os.Stdout, // устанавливаем, что логируем в консоль
	}

	if config.IsProduction() {
		// для прода формат логов будет JSON, чтобы можно было поместить в БД Elastic Search и его индексировать
		Log.Formatter = &logrus.JSONFormatter{}
	} else {
		// для разработки будет удобнее в текст, чтобы просто читать
		Log.Formatter = &logrus.TextFormatter{}
	}
}
