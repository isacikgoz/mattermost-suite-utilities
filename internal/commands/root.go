package commands

import "github.com/spf13/cobra"

func GetCommands() []*cobra.Command {
	commands := make([]*cobra.Command, 0)

	commands = append(commands, &cobra.Command{
		Use:     "generate <path>",
		Short:   "Generate client for Go",
		Example: `  generate ../mattermost-server/app`,
		Args:    cobra.ExactArgs(1),
		RunE:    generateCmdF,
	})

	return commands
}
