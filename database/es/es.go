// Package es
//
// @author: xwc1125
// @date: 2020/8/14 0014
package es

import (
	"context"
	"errors"
	"fmt"
	log2 "log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

type ES struct {
	es *elastic.Client
}

func NewES(hosts []string) (*ES, error) {
	if hosts == nil || len(hosts) == 0 {
		return nil, errors.New("es hosts is empty")
	}
	client, err := elastic.NewClient(
		elastic.SetURL(hosts...), // 设置host
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second), // 心跳检测时间
		elastic.SetMaxRetries(5),
		elastic.SetErrorLog(log2.New(os.Stderr, "ELASTIC ", log2.LstdFlags)), // 错误日志
		// elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),          // 正常日志
	)
	if err != nil {
		return nil, err
	}
	info, code, err := client.Ping(hosts[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esVersion, err := client.ElasticsearchVersion(hosts[0])
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elasticsearch version %s\n", esVersion)

	return &ES{
		es: client,
	}, nil
}

func (es *ES) Client() *elastic.Client {
	return es.es
}

// 批量创建时
func (es *ES) BulkIndexRequest(index string) *elastic.BulkIndexRequest {
	return elastic.NewBulkIndexRequest().Index(index)
}
