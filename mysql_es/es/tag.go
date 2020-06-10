package es

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go_wyy_micro/mysql_es/db"
	"go_wyy_micro/mysql_es/model"
	"log"
	"strconv"
	"strings"
)

func ReportTagToES(tag *model.Tag) {
	req := esapi.IndexRequest{
		Index:        "test",
		DocumentType: "tag",
		DocumentID:   strconv.Itoa(tag.Id),
		Body:         strings.NewReader(tag.MarshalToJson()),
		Refresh:      "true",
	}
	fmt.Println(111)

	resp, err := req.Do(context.TODO(), db.EsClient)
	if err != nil {
		log.Printf("ESIndexRequestErr: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.IsError() {
		log.Printf("ESIndexRequestErr: %s", resp.String())
	} else {
		log.Printf("ESIndexRequestOk: %s", resp.String())
	}
}
