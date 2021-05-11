package main

import (
	"github.com/isacikgoz/mattermost-suite-utilities/internal/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:               "clientgen",
		Short:             "Mattermost suite client generator for Go",
		Long:              `This CLI tool provides Mattermost suite client code genration for Go.`,
		DisableAutoGenTag: true,
	}

	rootCmd.AddCommand(commands.GetCommands()...)

	rootCmd.Execute()
}
