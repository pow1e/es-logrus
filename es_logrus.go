package es_logrus

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
	"time"
)

type ElasticHook struct {
	Client        *elasticsearch.Client // es的客户端
	Host          string                // es的host
	Index         IndexNameFunc         // 获取索引的名字
	Level         []logrus.Level        // 日志级别
	Ctx           context.Context       // 上下文
	CtxCancelFunc context.CancelFunc    // 上下文的函数
	FireFunc      FireFunc              // 执行hook的方法
}

// message 发送es的信息
type message struct {
	Host      string
	TimeStamp string `json:"@timestamp"`
	Message   string
	Data      logrus.Fields // 自定义的fileds
	Level     string
}

// createMessage 创建信息
func createMessage(entry *logrus.Entry, hook *ElasticHook) *message {
	//获取当前日至实例的等级
	level := entry.Level.String()

	// 如果Data中有logrus.ErrorKey 则将其的value转换成string(如果value的err)
	if value, ok := entry.Data[logrus.ErrorKey]; ok && value != nil {
		// 如果获取的值是err 则转换成string
		if err, ok := value.(error); ok {
			entry.Data[logrus.ErrorKey] = err.Error()
		}
	}

	return &message{
		Host:      hook.Host,
		TimeStamp: entry.Time.UTC().Format(time.RFC3339Nano),
		Message:   entry.Message,
		Data:      entry.Data,
		Level:     strings.ToUpper(level),
	}
}

// IndexNameFunc 获取index name
type IndexNameFunc func() string

// FireFunc 需要实现hook执行逻辑的函数
type FireFunc func(*logrus.Entry, *ElasticHook) error

// NewElasticHook 创建一个es hook对象
func NewElasticHook(client *elasticsearch.Client, host string, level logrus.Level, index string) (*ElasticHook, error) {
	return newElasticHookWithFunc(client, host, level, func() string { return index })
}

// newElasticHookWithFunc 创建一个es hook对象，通过IndexNameFunc这个方式
func newElasticHookWithFunc(client *elasticsearch.Client, host string, level logrus.Level, indexFunc IndexNameFunc) (*ElasticHook, error) {
	var levels []logrus.Level
	for _, l := range logrus.AllLevels {
		// 判断传入的level是在哪个等级上面
		if level >= l {
			levels = append(levels, l)
		}
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &ElasticHook{
		Client:        client,
		Host:          host,
		Index:         indexFunc,
		Level:         levels,
		Ctx:           ctx,
		CtxCancelFunc: cancelFunc,
		FireFunc:      syncFireFunc,
	}, nil
}

// syncFireFunc 异步发送 实现hook函数的具体逻辑
func syncFireFunc(entry *logrus.Entry, hook *ElasticHook) error {
	msg := createMessage(entry, hook)
	data, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	// 操作es
	request := esapi.IndexRequest{
		Index:   hook.Index(),
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	// 发送请求
	resp, err := request.Do(hook.Ctx, hook.Client)
	if err != nil {
		return err
	}

	// 解析响应
	// todo 完善解析json的操作
	log.Println(resp.String())

	return nil
}

// Levels 返回一个日至级别的切片，如果存在则会触发hook 即会调用Fire
func (hook *ElasticHook) Levels() []logrus.Level {
	return hook.Level
}

// Fire hook执行内容(需要三Levels中存在的日至等级)
func (hook *ElasticHook) Fire(entry *logrus.Entry) error {
	return hook.FireFunc(entry, hook)
}
