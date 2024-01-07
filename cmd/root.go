package cmd

import (
	"github.com/spf13/cobra"
)

// Execute executes cmd
func Execute() error {
	var service = &cobra.Command{
		Use:   "glookbs",
		Short: "a task rest api application",
	}
	service.AddCommand(
		runserver(),
		version,
	)
	return service.Execute()
}
