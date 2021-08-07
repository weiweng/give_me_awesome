package logs

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func Init() {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/tmp/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	Logger = logrus.New()

	//设置输出
	Logger.Out = src

	//设置日志级别
	Logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	Logger.Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	Logger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	Logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	Logger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	Logger.Panicf(format, args...)
}
