package logger

import "go.uber.org/zap"

func New() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	sugar := logger.Sugar()

	return sugar
}
