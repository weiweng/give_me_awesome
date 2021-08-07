package handler

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"give_me_awesome/es"
	"give_me_awesome/logs"
)

func More(c *gin.Context) {
	ctx := context.Background()
	index := getValue("book", c, "")
	hitId := getValue("id", c, "")

	more := make([]template.HTML, 0)
	if hitId != "" {
		list, err := es.GetContent(ctx, index, getIds(hitId))
		if err != nil {
			logs.Errorf("es query error[%v]", err)
		} else {
			for _, s := range list {
				more = append(more, template.HTML(s.Data))
			}
		}
	}
	// 往w里写入内容，就会在浏览器里输出
	c.HTML(http.StatusOK, "more.tmpl", gin.H{
		"Data": more,
	})
}

func getIds(hitId string) []string {
	id, err := strconv.Atoi(hitId)
	if err != nil {
		return []string{hitId}
	}
	ret := make([]string, 0, 4)
	id -= 1
	if id <= 0 {
		id = 0
	}
	ret = append(ret, strconv.Itoa(id))
	for i := 0; i <= 2; i++ {
		ret = append(ret, strconv.Itoa(id+i))
	}
	return ret
}
