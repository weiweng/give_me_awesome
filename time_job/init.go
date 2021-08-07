package time_job

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"path"
	"strconv"
	"time"

	"give_me_awesome/es"
	"give_me_awesome/handler"
	"give_me_awesome/helper"
	"give_me_awesome/logs"
	"give_me_awesome/model"
)

func Init() {
	go func() {
		du := 300 * time.Second
		ticker := time.NewTicker(du)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				logs.Infof("time_job do begin.............. \n")
				do()
			}
		}
	}()
}

func do() {
	defer func() {
		if err := recover(); err != nil {
			logs.Infof("time_job panic.............. \n")
		}
		logs.Infof("time_job do end.............. \n")
	}()
	bookList, err := es.GetAllBook(context.TODO())
	if err != nil {
		logs.Infof("get all book from es err[%v] \n", err)
		return
	}
	logs.Infof("online bookList[%v] \n", bookList)

	onlineBook := make(map[string]bool, len(bookList))
	for _, book := range bookList {
		onlineBook[book.Name] = false
	}

	bookPath, books, err := helper.GetBookList()
	if err != nil {
		logs.Infof("get all book from disk err[%v] \n", err)
		return
	}
	logs.Infof("book dir bookList[%v] \n", books)

	var target string
	for _, b := range books {
		if _, ok := onlineBook[b]; !ok {
			target = b
			break
		}
	}
	if target != "" {
		// 将target写入es
		book := path.Join(bookPath, target)
		index := getIndex(target, bookList)
		if err := handler.WriteToES(book, index); err != nil {
			logs.Infof("write book detail to es err[%v] \n", err)
			return
		}
		if err := es.AddBook(context.TODO(), &model.Book{
			ID:        strconv.Itoa(len(bookList) + 1),
			Name:      target,
			Index:     index,
			CreatedAt: time.Now(),
		}); err != nil {
			logs.Infof("write book info to es err[%v] \n", err)
			return
		}
	}
}

// getIndex 将bookName md5的后5位作为index，在onlineBook里去重，重复的话继续操作
func getIndex(bookName string, onlineBook []model.Book) string {
	h := md5.New()
	h.Write([]byte(bookName))
	bookMD5 := hex.EncodeToString(h.Sum(nil))
	for _, b := range onlineBook {
		if b.Index == bookMD5 {
			return getIndex(bookMD5, onlineBook)
		}
	}
	return bookMD5
}
