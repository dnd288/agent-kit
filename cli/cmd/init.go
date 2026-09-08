package cmd

import (
	"fmt"
	"os"
	"strings"

	kit "github.com/agent-kit/agent-kit-cli/kit"
	"github.com/agent-kit/agent-kit-cli/internal/config"
	"github.com/agent-kit/agent-kit-cli/internal/prompt"
	"github.com/agent-kit/agent-kit-cli/internal/scaffold"
	"github.com/spf13/cobra"
)

var (
	nonInteractive bool
	projectName    string
	prefix         string
	packageManager string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize agent-kit in the current project",
	Long: `Initialize agent-kit in the current project.

This installs process skills, instruction templates, Claude Code integration,
OpenSpec schemas, git hooks, and CI templates based on your answers.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&nonInteractive, "yes", false, "Use defaults and skip prompts")
	initCmd.Flags().StringVar(&projectName, "project", "", "Project name")
	initCmd.Flags().StringVar(&prefix, "prefix", "", "Skill prefix")
	initCmd.Flags().StringVar(&packageManager, "package-manager", "", "Package manager (pnpm, npm, yarn, bun)")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg := &config.ProjectConfig{}

	if nonInteractive {
		cfg.ProjectName = defaultProjectName(cwd)
		if projectName != "" {
			cfg.ProjectName = projectName
		}
		cfg.Prefix = defaultPrefix(cfg.ProjectName)
		if prefix != "" {
			cfg.Prefix = prefix
		}
		cfg.PackageManager = "pnpm"
		if packageManager != "" {
			cfg.PackageManager = packageManager
		}
		cfg.Stack = []string{"React", "Next.js", "Playwright"}
		cfg.Monorepo = true
		cfg.IncludeOpenSpec = true
		cfg.IncludeClaude = true
		cfg.IncludeHooks = true
		cfg.IncludeCI = true
	} else {
		cfg, err = askConfig(cwd)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\nInstalling agent-kit into %s\n", cwd)
	fmt.Printf("Project: %s\n", cfg.ProjectName)
	fmt.Printf("Skill prefix: %s\n\n", cfg.Prefix)

	if err := scaffold.Scaffold(kit.Content, cfg, cwd); err != nil {
		return err
	}

	if err := cfg.Save(cwd); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Println("\nAgent Kit installed.")
	fmt.Println("Next steps:")
	fmt.Println("  1. Read AGENTS.md and fill in project-specific boundaries")
	fmt.Println("  2. Review CONTRIBUTING.md and adjust commands")
	fmt.Println("  3. Run your validation command")

	return nil
}

func askConfig(cwd string) (*config.ProjectConfig, error) {
	var err error
	cfg := &config.ProjectConfig{}

	cfg.ProjectName, err = prompt.AskString("Project name", defaultProjectName(cwd))
	if err != nil {
		return nil, err
	}

	cfg.Prefix, err = prompt.AskString("Skill prefix", defaultPrefix(cfg.ProjectName))
	if err != nil {
		return nil, err
	}
	cfg.Prefix = sanitizePrefix(cfg.Prefix)

	cfg.PackageManager, err = prompt.AskSelect("Package manager", []string{"pnpm", "npm", "yarn", "bun"})
	if err != nil {
		return nil, err
	}

	cfg.Stack, err = prompt.AskMultiSelect("Stack", []string{
		"React",
		"Next.js",
		"Vue",
		"Svelte",
		"Fastify",
		"Express",
		"NestJS",
		"Prisma",
		"Playwright",
		"Cypress",
	})
	if err != nil {
		return nil, err
	}

	cfg.Monorepo, err = prompt.AskConfirm("Monorepo?", true)
	if err != nil {
		return nil, err
	}

	cfg.IncludeOpenSpec, err = prompt.AskConfirm("Include OpenSpec schemas?", true)
	if err != nil {
		return nil, err
	}

	cfg.IncludeClaude, err = prompt.AskConfirm("Include Claude Code integration?", true)
	if err != nil {
		return nil, err
	}

	cfg.IncludeHooks, err = prompt.AskConfirm("Include git hooks?", true)
	if err != nil {
		return nil, err
	}

	cfg.IncludeCI, err = prompt.AskConfirm("Include CI templates?", true)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultProjectName(cwd string) string {
	parts := strings.Split(strings.TrimRight(cwd, string(os.PathSeparator)), string(os.PathSeparator))
	if len(parts) == 0 {
		return "my-project"
	}
	return parts[len(parts)-1]
}

func defaultPrefix(name string) string {
	words := strings.FieldsFunc(strings.ToLower(name), func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	if len(words) == 0 {
		return "proj"
	}

	// Use first word if short, initials if many words
	if len(words) == 1 {
		return sanitizePrefix(words[0])
	}

	var initials strings.Builder
	for _, w := range words {
		if w != "" {
			initials.WriteByte(w[0])
		}
	}
	return sanitizePrefix(initials.String())
}

func sanitizePrefix(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}
