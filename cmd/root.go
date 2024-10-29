package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mobitec",
	Short: "Control a mobitec flipdot destination board",
}

func init() {}

func Execute() error {
	return rootCmd.Execute()
}
