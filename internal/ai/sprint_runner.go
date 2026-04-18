package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/fernandoeho/smiddy/internal/ai/agents"
	"github.com/fernandoeho/smiddy/internal/ai/claude"
	"github.com/fernandoeho/smiddy/internal/fs"
	"github.com/fernandoeho/smiddy/internal/ui"
)

const maxIterations = 5

// SprintState represents the current phase of the sprint loop
type SprintState int

const (
	StateAnalyze SprintState = iota
	StateImplement
	StateReview
	StateComplete
	StateFailed
)

func (s SprintState) String() string {
	switch s {
	case StateAnalyze:
		return "analyzing"
	case StateImplement:
		return "implementing"
	case StateReview:
		return "reviewing"
	case StateComplete:
		return "complete"
	case StateFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// RunSprint executes the full sprint loop for the given sprint directory
func RunSprint(client *claude.Client, sprintDir string) error {
	ctx := context.Background()

	specsPath := fmt.Sprintf("%s/specs.md", sprintDir)
	statusPath := fmt.Sprintf("%s/status.md", sprintDir)

	if !fs.FileExists(specsPath) {
		return fmt.Errorf("specs.md not found at %s — run `smiddy new` first", specsPath)
	}

	specs, err := fs.ReadFile(specsPath)
	if err != nil {
		return err
	}

	goalsPath := ".smiddy/project-goals.md"
	goals := ""
	if fs.FileExists(goalsPath) {
		goals, _ = fs.ReadFile(goalsPath)
	}

	mapPath := ".smiddy/project-map.md"
	projectMap := ""
	if fs.FileExists(mapPath) {
		projectMap, _ = fs.ReadFile(mapPath)
	}

	state := StateAnalyze
	iteration := 0
	var architectPlan string
	var implementationReport string

	ui.Bold("Starting sprint...\n\n")

	for state != StateComplete && state != StateFailed {
		switch state {

		case StateAnalyze:
			ui.Info("[%s] Architect analyzing specs...\n", StateAnalyze)

			userPrompt := fmt.Sprintf(
				"## Project Goals\n%s\n\n## Project Map\n%s\n\n## Sprint Specs\n%s",
				goals, projectMap, specs,
			)

			plan, err := client.Complete(ctx, agents.ArchitectSystem, []claude.Message{
				{Role: "user", Content: userPrompt},
			})
			if err != nil {
				return fmt.Errorf("architect analysis failed: %w", err)
			}

			architectPlan = plan
			fmt.Println(architectPlan)

			reportPath := fmt.Sprintf("%s/architect-plan.md", sprintDir)
			_ = fs.WriteFile(reportPath, architectPlan)

			state = StateImplement

		case StateImplement:
			iteration++
			ui.Info("\n[%s] Iteration %d/%d — implementation agent running...\n", StateImplement, iteration, maxIterations)

			userPrompt := fmt.Sprintf(
				"## Sprint Specs\n%s\n\n## Architect Plan\n%s",
				specs, architectPlan,
			)

			report, err := client.Complete(ctx, agents.GoHorseSystem, []claude.Message{
				{Role: "user", Content: userPrompt},
			})
			if err != nil {
				return fmt.Errorf("implementation agent failed: %w", err)
			}

			implementationReport = report
			fmt.Println(implementationReport)

			reportPath := fmt.Sprintf("%s/implementation-%d.md", sprintDir, iteration)
			_ = fs.WriteFile(reportPath, implementationReport)

			state = StateReview

		case StateReview:
			ui.Info("\n[%s] Architect reviewing iteration %d...\n", StateReview, iteration)

			userPrompt := fmt.Sprintf(
				"## Original Specs\n%s\n\n## Architect Plan\n%s\n\n## Implementation Report\n%s",
				specs, architectPlan, implementationReport,
			)

			review, err := client.Complete(ctx, agents.ArchitectReviewSystem, []claude.Message{
				{Role: "user", Content: userPrompt},
			})
			if err != nil {
				return fmt.Errorf("architect review failed: %w", err)
			}

			fmt.Println(review)

			reviewPath := fmt.Sprintf("%s/review-%d.md", sprintDir, iteration)
			_ = fs.WriteFile(reviewPath, review)

			if strings.Contains(review, "COMPLETE") {
				state = StateComplete
			} else if iteration >= maxIterations {
				state = StateFailed
			} else {
				architectPlan = review
				state = StateImplement
			}
		}

		writeStatus(statusPath, state, iteration)
	}

	if state == StateComplete {
		ui.Success("\nSprint complete after %d iteration(s).\n", iteration)
	} else {
		ui.Warn("\nSprint reached max iterations (%d) without completing. Review reports in %s.\n", maxIterations, sprintDir)
	}

	return nil
}

func writeStatus(path string, state SprintState, iteration int) {
	content := fmt.Sprintf("# Sprint Status\n\nState: %s\nIteration: %d\n", state, iteration)
	_ = fs.WriteFile(path, content)
}
