package es

import (
	"fmt"

	"github.com/olivere/elastic/v7"
)

var client *elastic.Client

func Init() {
	var err error
	client, err = elastic.NewClient(
		// 设置基于http base auth验证的账号和密码
		// elastic.SetBasicAuth("user", "secret"),
		// 设置ES服务地址，支持多个地址
		elastic.SetSniff(false),
		elastic.SetURL("http://localhost:9300"))
	if err != nil {
		fmt.Printf("连接失败: %v\n", err)
	} else {
		fmt.Println("连接成功")
	}
}
