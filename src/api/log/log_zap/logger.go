package log_zap

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"}, // записывать в консоль
		Encoding:    "json",             // формат JSON
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		EncoderConfig: zapcore.EncoderConfig{ // настройка как будет выглядеть объект лога
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,    // отображения времени
			EncodeLevel:  zapcore.LowercaseLevelEncoder, //  слова с маленькой буквой
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	Log, err = logConfig.Build()
	if err != nil {
		panic(err)
	}

}

func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

func Debug(message string, tags ...zap.Field) {
	// debug log interface
	Log.Debug(message, tags...)
	// закрыть связь с консолем по вывода лога в консоль
	Log.Sync()
}

func Info(message string, tags ...zap.Field) {
	Log.Info(message, tags...)
	// закрыть связь с консолем по вывода лога в консоль
	Log.Sync()
}

func Error(message string, err error, tags ...zap.Field) {
	// error log interface
	message = fmt.Sprintf("%s - ERROR - %v", message, err)
	Log.Error(message, tags...)
	// закрыть связь с консолем по вывода лога в консоль
	Log.Sync()
}
