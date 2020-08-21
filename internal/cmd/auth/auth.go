package auth

import "github.com/spf13/cobra"

// NewCommand creates new report command
func NewCommand() *cobra.Command {
	var c = &cobra.Command{
		Use:   "auth",
		Short: "Manage Glass Factory authentication credentials.",
	}
	c.AddCommand(NewLoginCommand())
	return c
}
