package es

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/olivere/elastic/v7"

	"give_me_awesome/logs"
	"give_me_awesome/model"
)

const (
	bookIndex    = "book" // 书籍的固定索引
	noIndexError = "elastic: Error 404 (Not Found): no such index"
)

func DelBook(ctx context.Context, index string) error {
	_, err := client.DeleteIndex(index).Do(ctx) // 执行请求，需要传入一个上下文对象
	return err
}

func AddBook(ctx context.Context, book *model.Book) error {
	_, err := client.Index().
		Index(bookIndex). // 设置索引名称
		Id(book.ID).      // 设置文档id
		BodyJson(book).   // 指定前面声明struct对象
		Do(ctx)           // 执行请求，需要传入一个上下文对象
	return err
}

// GetDataByIndexAndID 根据id查询文档
func GetAllBook(ctx context.Context) ([]model.Book, error) {
	searchResult, err := client.Search().
		Index(bookIndex). // 设置索引名
		Pretty(true).     // 查询结果返回可读性较好的JSON格式
		Do(ctx)           // 执行请求
	if err != nil {
		if strings.Contains(err.Error(), noIndexError) {
			return []model.Book{}, nil
		}
		return nil, err
	}
	logs.Infof("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())
	ret := make([]model.Book, 0, searchResult.TotalHits())
	if searchResult.TotalHits() > 0 {
		var c model.Book
		for _, item := range searchResult.Each(reflect.TypeOf(c)) {
			// 转换成Article对象
			if t, ok := item.(model.Book); ok {
				ret = append(ret, t)
			}
		}
		return ret, nil
	}
	return ret, nil
}

func IsExists(ctx context.Context, index string) bool {
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		logs.Errorf("isExists error[%v]", err)
		return false
	}
	return exists
}

// BatchAddContent 批量新增内容
func BatchAddContent(ctx context.Context, data []*model.Content) error {
	var err error
	cnt := 0
	for _, item := range data {
		_, err = client.Index().
			Index(item.Index). // 设置索引名称
			Id(item.ID).       // 设置文档id
			BodyJson(item).    // 指定前面声明struct对象
			Do(ctx)            // 执行请求，需要传入一个上下文对象
		if err != nil {
			err = fmt.Errorf("[%w]", err)
		}
		cnt++
		if cnt%1000 == 0 {
			logs.Errorf("======= write %v ", item.ID)
		}
	}
	return err
}

// GetDataByIndexAndID 根据id查询文档
func GetDataByIndexAndID(ctx context.Context, index string, id string) (*model.Content, error) {
	get1, err := client.Get().
		Index(index). // 指定索引名
		Id(id).       // 设置文档id
		Do(ctx)       // 执行请求
	if err != nil {
		return nil, err
	}
	data, _ := get1.Source.MarshalJSON()
	ret := model.Content{}
	err = json.Unmarshal(data, &ret)
	return &ret, err
}

// DeleteByIndexAndID 根据id删除一条数据
func DeleteByIndexAndID(ctx context.Context, index string, id string) error {
	_, err := client.Delete().
		Index(index).
		Id(id). // 文档id
		Do(ctx)
	return err
}

// PhraseQuery 模糊搜索
func GetContent(ctx context.Context, index string, ids []string) ([]model.Content, error) {
	query := elastic.NewIdsQuery().Ids(ids...)
	searchResult, err := client.Search(index).
		Query(query). // 设置查询条件
		Size(len(ids)).
		Pretty(true). // 查询结果返回可读性较好的JSON格式
		Do(ctx)       // 执行请求

	if err != nil {
		return nil, err
	}
	logs.Infof("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())
	ret := make([]model.Content, 0, 10)
	if searchResult.TotalHits() > 0 {
		var c model.Content
		for _, item := range searchResult.Each(reflect.TypeOf(c)) {
			// 转换成Article对象
			if t, ok := item.(model.Content); ok {
				ret = append(ret, t)
			}
		}
		return ret, nil
	}
	return ret, nil
}

// PhraseQuery 模糊搜索
func PhraseQuery(ctx context.Context, index string, query []string, pageNumber int) ([]*model.Content, error) {

	boolQuery := elastic.NewBoolQuery().Must()

	for _, item := range query {
		termQuery := elastic.NewMatchPhraseQuery("data", item)
		boolQuery.Must(termQuery)
	}

	searchResult, err := client.Search().
		Index(index).     // 设置索引名
		Query(boolQuery). // 设置查询条件
		// Sort("Created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(10 * (pageNumber - 1)). // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(10).                    // 设置分页参数 - 每页大小
		Pretty(true).                // 查询结果返回可读性较好的JSON格式
		Do(ctx)                      // 执行请求

	if err != nil {
		return nil, err
	}
	logs.Infof("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())
	ret := make([]*model.Content, 0, 10)
	if searchResult.TotalHits() > 0 {
		var c model.Content
		for _, item := range searchResult.Each(reflect.TypeOf(c)) {
			// 转换成Article对象
			if t, ok := item.(model.Content); ok {
				ret = append(ret, &t)
			}
		}
		return ret, nil
	}
	return ret, nil
}

// PhraseQueryHighlight 模糊搜索高亮
func PhraseQueryHighlight(ctx context.Context, index string, query []string, pageNumber int) ([]model.QueryInfo, error) {

	boolQuery := elastic.NewBoolQuery().Must()

	for _, item := range query {
		termQuery := elastic.NewMatchPhraseQuery("data", item)
		boolQuery.Must(termQuery)
	}

	// 定义highlight
	highlight := elastic.NewHighlight()
	// 指定需要高亮的字段
	highlight = highlight.Fields(elastic.NewHighlighterField("data"))
	// 指定高亮的返回逻辑 <span style='color: red;'>...msg...</span>
	highlight = highlight.PreTags("<span style='color: blue;'>").PostTags("</span>")

	searchResult, err := client.Search().
		Index(index). // 设置索引名
		Highlight(highlight).
		Query(boolQuery). // 设置查询条件
		// Sort("Created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(10 * (pageNumber - 1)). // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(10).                    // 设置分页参数 - 每页大小
		Pretty(true).                // 查询结果返回可读性较好的JSON格式
		Do(ctx)                      // 执行请求

	if err != nil {
		return nil, err
	}
	logs.Infof("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())
	ret := make([]model.QueryInfo, 0, 10)
	if searchResult.TotalHits() > 0 {
		for _, item := range searchResult.Hits.Hits {
			if item == nil {
				continue
			}
			con := ""
			for _, v := range item.Highlight {
				con += strings.Join(v, "")
			}
			ret = append(ret,
				model.QueryInfo{
					Content: con,
					Id:      item.Id,
				})
		}
		return ret, nil
	}
	return ret, nil
}
