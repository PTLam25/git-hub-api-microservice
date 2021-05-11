package log

import (
	"fmt"
	"github.com/PTLam25/git-hub-api-microservice/src/api/config"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
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

func Debug(message string, tags ...string) {
	// debug log interface
	if Log.Level < logrus.DebugLevel {
		return
	}

	Log.WithFields(parseFields(tags...)).Debug(message)
}

func Info(message string, tags ...string) {
	// info log interface
	if Log.Level < logrus.InfoLevel {
		return
	}
	// WithFields позволяет добавить дополнительные поля к объекту логгинга для сохранения в БД Elastic Search
	// Пример: INFO[0025] the urls successfully mapped                  status=success step=2
	Log.WithFields(parseFields(tags...)).Info(message)
}

func Error(message string, err error, tags ...string) {
	// error log interface
	if Log.Level < logrus.ErrorLevel {
		return
	}

	message = fmt.Sprintf("%s - ERROR - %s", message, err)
	Log.WithFields(parseFields(tags...)).Error(message)
}

func parseFields(tags ...string) logrus.Fields {
	// выделяем память на создания Map (Fields-это Map) с такой фиксированной длинны
	result := make(logrus.Fields, len(tags))

	for _, tag := range tags {
		elements := strings.Split(tag, ":")
		result[strings.TrimSpace(elements[0])] = strings.TrimSpace(elements[1])
	}

	return result
}
