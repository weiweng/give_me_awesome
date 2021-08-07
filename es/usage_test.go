package es

import (
	"context"
	"give_me_awesome/logs"
	"testing"
)

func TestGetAllBook(t *testing.T) {
	logs.Init()
	Init()
	bookList, err := GetAllBook(context.Background())
	t.Logf("%v %v", bookList, err)
	t.FailNow()
}

func TestDelBook(t *testing.T) {
	logs.Init()
	Init()
	err := DelBook(context.Background(), "book")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log("del ok")
	t.FailNow()
}

func TestGetContent(t *testing.T) {
	logs.Init()
	Init()
	data, err := GetContent(context.TODO(), "91274646efe920929aca3019542e7a15", []string{"196"})
	t.Logf("err[%v] \n", err)
	t.Logf("data[%v] \n", data)
	t.Fail()
}
