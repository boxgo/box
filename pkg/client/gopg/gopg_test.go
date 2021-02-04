package gopg_test

import (
	"fmt"

	"github.com/boxgo/box/pkg/client/gopg"
)

type (
	TestUserinfo struct {
		UserName string
		Age      int
	}
)

func Example() {
	pg := gopg.StdConfig("default").Build()

	ts := &TestUserinfo{
		UserName: "box",
		Age:      18,
	}
	if err := pg.Model(ts).CreateTable(&gopg.CreateTableOptions{
		IfNotExists: true,
		Temp:        false,
	}); err != nil {
		panic(err)
	}

	if _, err := pg.Model(ts).Insert(); err != nil {
		panic(err)
	}

	fmt.Println(pg.Model(ts).DropTable(&gopg.DropTableOptions{}))
	// output: <nil>
}
