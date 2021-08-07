package model

import (
	"html/template"
	"time"
)

type Book struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	Index     string    `json:"index"`
	CreatedAt time.Time `json:"created_at"`
}

type Content struct {
	Index  string `json:"index"`
	Offset int64  `json:"offset"`
	ID     string `json:"id"`
	Data   string `json:"data"`
	Tag    string `json:"tag"`
}

type DataInfo struct {
	Content template.HTML
	More    template.HTML
}

type QueryInfo struct {
	Content string
	Id      string
}
