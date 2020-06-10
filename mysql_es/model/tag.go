package model

import (
	"bytes"
	"encoding/json"
)

type TagBody struct {
	Name string `json:"name"`
}

type Tag struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type SearchTagReqBody struct {
	Keyword string `json:"keyword"`
}

func (t *Tag) MarshalToJson() string {
	datas, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(datas)
}

type O map[string]interface{}

func (o *O) MarshalTOJsonBytes() *bytes.Buffer {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(o); err != nil {
		panic(err)
	}
	return &buf
}
