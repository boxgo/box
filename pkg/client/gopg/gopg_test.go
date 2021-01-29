package gopg_test

import (
	"testing"

	"github.com/boxgo/box/pkg/client/gopg"
)

type (
	TestUserinfo struct {
		UserName string
		Age      int
	}
)

func TestModel(t *testing.T) {
	pg := (&gopg.Config{URI: "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"}).Build()

	ts := &TestUserinfo{
		UserName: "box",
		Age:      18,
	}
	if err := pg.Model(ts).CreateTable(&gopg.CreateTableOptions{
		IfNotExists: true,
		Temp:        false,
	}); err != nil {
		t.Fatalf("create table error %#v", err)
	} else {
		t.Log("create table success")
	}

	if result, err := pg.Model(ts).Insert(); err != nil {
		t.Fatalf("insert error %#v", err)
	} else {
		t.Logf("insert success %#v", result)
	}
}
