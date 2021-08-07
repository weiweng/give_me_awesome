package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"give_me_awesome/es"
	"give_me_awesome/logs"
	"give_me_awesome/model"
)

func getValue(key string, c *gin.Context, def string) string {
	ret, ok := c.GetPostForm(key)
	if ok {
		return ret
	}
	ret, ok = c.GetQuery(key)
	if ok {
		return ret
	}
	return def
}

func Query(c *gin.Context) {
	ctx := context.Background()
	index := getValue("book", c, "")
	pageStr := getValue("page", c, "1")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	query := getValue("key", c, "")
	queryList := strings.Split(query, ",")
	hitId := getValue("id", c, "")

	allBook, err := es.GetAllBook(ctx)
	if err != nil {
		logs.Errorf("es getAllBook error[%v]", err)
		return
	}

	tmp := make([]model.DataInfo, 0)
	if len(queryList) > 0 {
		list, err := es.PhraseQueryHighlight(ctx, index, queryList, page)
		if err != nil {
			logs.Errorf("es query error[%v]", err)
			return
		}
		for _, s := range list {
			tmp = append(tmp, model.DataInfo{
				Content: template.HTML(s.Content),
				More:    template.HTML(fmt.Sprintf("/v1/query?book=%s&key=%s&page=%d&id=%s", index, query, page, s.Id)),
			})
		}
	}

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

	pre := template.HTML(fmt.Sprintf("/v1/query?book=%s&key=%s&page=%d", index, query, page-1))
	next := template.HTML(fmt.Sprintf("/v1/query?book=%s&key=%s&page=%d", index, query, page+1))
	// 往w里写入内容，就会在浏览器里输出
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"Data":     tmp,
		"MoreData": more,
		"Pre":      pre,
		"Next":     next,
		"Book":     index,
		"Query":    query,
		"BookList": allBook,
	})
}
