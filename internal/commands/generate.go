package commands

import (
	"fmt"

	"github.com/isacikgoz/mattermost-suite-utilities/internal/generator"
	"github.com/isacikgoz/mattermost-suite-utilities/internal/parser"
	"github.com/spf13/cobra"
)

func generateCmdF(cmd *cobra.Command, args []string) error {
	cmd.Printf("parsing %q\n", args[0])
	st, err := parser.ParseDirectory(args[0])

	if err != nil {
		return err
	}

	err = generator.Render(st, "client")
	if err != nil {
		return err
	}

	fmt.Printf("OUTPUT %s\n", st.Name)

	return nil
}
