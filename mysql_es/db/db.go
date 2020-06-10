package db

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DbEngine *gorm.DB
	EsClient *elasticsearch.Client
)

func GetDb() {
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/test?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("gorm open failed, err:", err)
	}
	DbEngine = db
}

//初始化es
func GetEs() {
	config := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		fmt.Println("elasticsearch init failed, err:", err)
	}
	res, err := client.Info()
	if err != nil {
		panic(err)
	}
	if res.IsError() {
		panic(res.String())
	}
	EsClient = client
}

func init() {
	go GetDb()
	go GetEs()
}
