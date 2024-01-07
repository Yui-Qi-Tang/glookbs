package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1")
	},
}
