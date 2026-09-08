package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "agent-kit",
	Short: "Agent-driven development kit",
	Long: `Agent Kit scaffolds structured agent-assisted development into your project.

It installs process skills, instruction templates, specification schemas,
and tooling that encode battle-tested development practices:

  - Five-axis code review
  - Adversarial verification
  - RED→GREEN test-driven development
  - Structured bug fixing and refactoring
  - Specification-driven change management
  - And more

Usage:
  agent-kit init          Initialize agent-kit in the current project
  agent-kit add <skill>   Add a single skill to an existing project
  agent-kit list          List all available skills`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("agent-kit v%s\n", version)
	},
}
