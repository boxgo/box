package main

import (
	"log"

	"github.com/boxgo/box/cmd/rest"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "box",
	Short: "Box toolchains",
}

func main() {
	rest.Register(root)

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}
