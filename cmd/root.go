package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "woodhouse",
	Short: "Your own personal Sir Arthur Henry Woodhouse, ready to help you with mundane tasks.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
