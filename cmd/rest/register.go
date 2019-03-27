package rest

import (
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	root.AddCommand(newInitCommand())
}
