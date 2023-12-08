package es_logrus

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	entry := logger.WithFields(logrus.Fields{
		"作者": "cz",
	})
	entry.Info("nihao")
}

func TestEs(t *testing.T) {
	InitEs()
	hook, err := NewElasticHook(EsClient, "localhost", logrus.DebugLevel, "my_index")
	if err != nil {
		t.Fatal(err)
	}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(hook)
	logger.Error("你好")
}
