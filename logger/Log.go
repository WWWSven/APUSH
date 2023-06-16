package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
)

var logger = logrus.New()

func init() {
	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	f, err := os.OpenFile("push.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("日志文件创建失败! %+v", err)
		return
	}
	writer := io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = writer
	logger.Out = writer
	logger.Level = logrus.InfoLevel
	logger.Info("日志配置成功")
}

func GetLogger() *logrus.Logger {
	return logger
}

func Info(format string, args ...interface{}) {
	if len(args) == 0 {
		logger.Info(format)
		return
	}
	logger.Infof(format, args)
}
