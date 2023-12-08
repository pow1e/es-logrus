package es_logrus

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"time"
)

var LogrusObj *logrus.Logger

const FileNameLayOut = "2006-01-02"

// InitLog 初始化日志
func InitLog() {
	if LogrusObj != nil {
		file, err := setOutputFile()
		if err != nil {
			panic(err)
		}
		LogrusObj.SetOutput(file)
		return
	}
	// 初始化
	logger := logrus.New()
	file, err := setOutputFile()
	if err != nil {
		panic(err)
	}
	logger.SetOutput(file)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogrusObj = logger
}

// setOutputFile 按照日期进行分割，如果不存在则创建，存在就写入日至中
func setOutputFile() (*os.File, error) {
	now := time.Now()
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 设置日志文件位置
	logFilePath := filepath.Join(dir, "/logs/")

	// 使用os.Stat 判断文件或文件夹是否存在
	// 如果返回 err是nil则表明文件或文件夹存在
	// 否则使用os.IsNotExists() 判断 ==》 如果为true则表明文件或文件夹不存在
	_, err = os.Stat(logFilePath)
	if os.IsNotExist(err) {
		// true 文件夹或文件不存在
		if err = os.MkdirAll(logFilePath, 0777); err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	// 文件夹存在
	logFileName := now.Format(FileNameLayOut) + ".log"
	fileName := filepath.Join(logFilePath, logFileName)

	// err不为空 则需要使用os.IsNotExists()判断文件是否存在
	if _, err = os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			// 创建文件
			if _, err = os.Create(fileName); err != nil {
				log.Fatal(err)
				return nil, err
			}
		}
	}

	// 写入文件
	// O_WRONLY 以只写的方式
	// O_APPEND 追加
	// ModeAppend 只能写 且只能写到末尾
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	return file, err
}
