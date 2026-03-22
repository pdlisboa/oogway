package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(chatsCommands)
}

var rootCmd = &cobra.Command{
	Use: "agent",
}

func Execute() error {
	return rootCmd.Execute()
}
