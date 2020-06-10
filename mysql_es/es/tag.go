package es

import (
	"context"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
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

func SearchTagsFromEs(keyword string) ([]*model.Tag, error) {
	//构建查询
	query := model.O{
		"query": model.O{
			"match_phrase_prefix": model.O{
				"name":           keyword,
				"max_expansions": 50,
			},
		},
	}

	jsonBuf := query.MarshalTOJsonBytes()

	//发出查询请求
	resp, err := db.EsClient.Search(
		db.EsClient.Search.WithContext(context.TODO()),
		db.EsClient.Search.WithIndex("test"),
		db.EsClient.Search.WithBody(jsonBuf),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, errors.New(resp.Status())
	}

	js, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	hitsjs := js.GetPath("hist", "hits")
	hits, err := hitsjs.Array()
	if err != nil {
		return nil, err
	}

	hitsLen := len(hits)
	if hitsLen == 0 {
		return []*model.Tag{}, nil
	}

	tags := make([]*model.Tag, 0, len(hits))
	for idx := 0; idx < hitsLen; idx++ {
		sourceJS := hitsjs.GetIndex(idx).Get("_source")

		tagID, err := sourceJS.Get("tag_id").Int()
		if err != nil {
			return nil, err
		}

		tagName, err := sourceJS.Get("name").String()
		if err != nil {
			return nil, err
		}

		tagEntity := &model.Tag{Id: tagID, Name: tagName}
		tags = append(tags, tagEntity)
	}

	return tags, nil
}
