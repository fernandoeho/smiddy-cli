package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fernandoeho/smiddy/internal/fs"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive project onboarding",
	Run: func(cmd *cobra.Command, args []string) {
		if !fs.FileExists(".smiddy") {
			fmt.Fprintln(os.Stderr, "No .smiddy directory found. Run `smiddy init` first.")
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Smiddy setup — let's define your project.")

		vision := prompt(reader, "What are you building? (one sentence)")
		audience := prompt(reader, "Who is it for?")
		metrics := prompt(reader, "How will you measure success?")
		constraints := prompt(reader, "Any hard constraints? (tech, time, budget — or press Enter to skip)")

		content := fmt.Sprintf(`# Project Goals

## Vision
%s

## Target audience
%s

## Success metrics
%s

## Constraints
%s
`, vision, audience, metrics, constraints)

		if err := fs.WriteFile(".smiddy/project-goals.md", content); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("\nproject-goals.md updated. Run `smiddy new` to create your first sprint.")

	},
}

func prompt(reader *bufio.Reader, question string) string {
	fmt.Printf("  %s\n  > ", question)
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
