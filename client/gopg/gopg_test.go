package gopg_test

import (
	"fmt"

	gopg2 "github.com/boxgo/box/v2/client/gopg"
)

type (
	TestUserinfo struct {
		UserName string
		Age      int
	}
)

func Example() {
	pg := gopg2.StdConfig("default").Build()

	ts := &TestUserinfo{
		UserName: "box",
		Age:      18,
	}
	if err := pg.Model(ts).CreateTable(&gopg2.CreateTableOptions{
		IfNotExists: true,
		Temp:        false,
	}); err != nil {
		panic(err)
	}

	if _, err := pg.Model(ts).Insert(); err != nil {
		panic(err)
	}

	fmt.Println(pg.Model(ts).DropTable(&gopg2.DropTableOptions{}))
	// output: <nil>
}
