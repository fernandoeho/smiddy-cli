package cmd

import (
	"fmt"
	"os"

	"github.com/fernandoeho/smiddy/internal/ai"
	"github.com/fernandoeho/smiddy/internal/ai/claude"
	"github.com/fernandoeho/smiddy/internal/fs"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the current sprint",
	Run: func(cmd *cobra.Command, args []string) {
		if !fs.FileExists(".smiddy") {
			fmt.Fprintln(os.Stderr, "No .smiddy directory found. Run `smiddy init` first.")
			os.Exit(1)
		}

		n, err := fs.NextSprintNumber(".smiddy/sprints")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if n == 1 {
			fmt.Fprintln(os.Stderr, "No sprints found. Run `smiddy new` first.")
			os.Exit(1)
		}

		current := n - 1
		sprintDir := fmt.Sprintf(".smiddy/sprints/%d", current)

		client, err := claude.NewClient()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := ai.RunSprint(client, sprintDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
