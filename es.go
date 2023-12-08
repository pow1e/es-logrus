package es_logrus

import (
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"log"
)

var EsClient *elasticsearch.Client

// InitEs 初始化es
func InitEs() {
	esConn := fmt.Sprintf("http://%s:%s", "192.168.61.129", "9200")
	cfg := elasticsearch.Config{
		Addresses: []string{esConn},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Panic(err)
	}
	EsClient = client
}
