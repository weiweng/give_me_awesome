package handler

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"give_me_awesome/es"
	"give_me_awesome/model"
)

func WriteToES(file string, index string) error {
	contentList, err := readBook(file, index)
	if err != nil {
		return err
	}
	exists := es.IsExists(context.TODO(), index)
	if exists {
		return nil
	}
	if len(contentList) > 0 {
		fmt.Printf("sum num is %v \n", len(contentList))
		return es.BatchAddContent(context.TODO(), contentList)
	}
	return nil
}

func readBook(filePath string, index string) ([]*model.Content, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("文件打开失败[%w]", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	contentList := make([]*model.Content, 0, 1024)
	id := 0
	off := 0
	for {
		p, err := reader.ReadString('\n')
		// io.EOF 表示文件的末尾
		if err == io.EOF {
			break
		}
		p = strings.TrimSpace(p)
		contentList = append(contentList, &model.Content{
			Index:  index,
			Offset: int64(off + utf8.RuneCountInString(p)),
			ID:     strconv.Itoa(id),
			Data:   p,
			Tag:    "",
		})
		id++
	}
	return contentList, nil
}
