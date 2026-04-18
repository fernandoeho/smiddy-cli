package cmd

import (
	"fmt"
	"os"

	"github.com/fernandoeho/smiddy/internal/fs"
	"github.com/fernandoeho/smiddy/internal/templates"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new sprint",
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

		sprintDir := fmt.Sprintf(".smiddy/sprints/%d", n)

		if err := fs.EnsureDir(sprintDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		specsPath := fmt.Sprintf("%s/specs.md", sprintDir)
		content := fmt.Sprintf(templates.SpecsTemplate, n)

		if err := fs.WriteFile(specsPath, content); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Created sprint %d → %s\n", n, specsPath)
		fmt.Println("Edit specs.md, then run `smiddy run` to start the sprint.")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
