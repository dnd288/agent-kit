package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/agent-kit/agent-kit-cli/internal/config"
	"github.com/agent-kit/agent-kit-cli/internal/prompt"
	"github.com/agent-kit/agent-kit-cli/internal/scaffold"
	"github.com/agent-kit/agent-kit-cli/internal/skills"
	kit "github.com/agent-kit/agent-kit-cli/kit"
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
	start := time.Now()
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
		cfg.TicketTracker = detectTicketTracker(cwd)
		cfg.FeatureFlow = defaultFeatureFlow()
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

	if err := writeInitLog(cwd, cfg, time.Since(start)); err != nil {
		return fmt.Errorf("write init log: %w", err)
	}

	fmt.Println("\nAgent Kit installed.")
	fmt.Printf("Init log: %s\n", config.InitLogFile)
	fmt.Println("Next steps:")
	fmt.Println("  1. Read AGENTS.md and fill in project-specific boundaries")
	fmt.Println("  2. Review CONTRIBUTING.md and adjust commands")
	fmt.Println("  3. Run your validation command")

	return nil
}

// writeInitLog records the completed setup at the project root for later
// benchmarking.
func writeInitLog(targetDir string, cfg *config.ProjectConfig, elapsed time.Duration) error {
	fileCount, err := countSkillFiles(targetDir)
	if err != nil {
		return err
	}

	optional := 0
	for _, name := range cfg.InstalledSkills {
		if skills.IsOptional(name) {
			optional++
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Agent Kit init log\n\n")
	fmt.Fprintf(&b, "- Date: %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(&b, "- Duration: %s\n", elapsed.Round(time.Millisecond))
	fmt.Fprintf(&b, "- Method: agent-kit CLI\n")
	fmt.Fprintf(&b, "- Project: %s\n", cfg.ProjectName)
	fmt.Fprintf(&b, "- Prefix: %s\n", cfg.Prefix)
	fmt.Fprintf(&b, "- Package manager: %s\n", cfg.PackageManager)
	fmt.Fprintf(&b, "- Stack: %s\n", strings.Join(cfg.Stack, ", "))
	fmt.Fprintf(&b, "- Ticket tracker: %s\n", cfg.TicketTracker)
	fmt.Fprintf(&b, "- Modules: %s\n", strings.Join(enabledModules(cfg), ", "))
	fmt.Fprintf(&b, "- Skills installed: %d (%d optional)\n", len(cfg.InstalledSkills), optional)
	fmt.Fprintf(&b, "- Skill files: %d under .agents/skills/\n", fileCount)

	if len(cfg.InstalledSkills) > 0 {
		b.WriteString("\n## Installed skills\n\n| Skill | Category |\n|---|---|\n")
		for _, name := range cfg.InstalledSkills {
			fmt.Fprintf(&b, "| %s-%s | %s |\n", cfg.Prefix, name, skills.Category(name))
		}
	}

	return os.WriteFile(filepath.Join(targetDir, config.InitLogFile), []byte(b.String()), 0644)
}

// countSkillFiles counts files under <targetDir>/.agents/skills.
func countSkillFiles(targetDir string) (int, error) {
	count := 0
	err := filepath.WalkDir(filepath.Join(targetDir, ".agents", "skills"),
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				count++
			}
			return nil
		})
	return count, err
}

// enabledModules lists the modules enabled in the config, in install order.
func enabledModules(cfg *config.ProjectConfig) []string {
	modules := []string{"skills", "templates"}
	if cfg.IncludeClaude {
		modules = append(modules, "claude")
	}
	if cfg.IncludeOpenSpec {
		modules = append(modules, "openspec")
	}
	if cfg.IncludeHooks {
		modules = append(modules, "hooks")
	}
	if cfg.IncludeCI {
		modules = append(modules, "ci")
	}
	return modules
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

	fmt.Printf("\nTicket tracker detected: %s\n", detectTicketTracker(cwd))
	cfg.TicketTracker, err = prompt.AskSelect("Ticket tracker", trackerChoices(detectTicketTracker(cwd)))
	if err != nil {
		return nil, err
	}

	fmt.Println("\nFeature flow stages: capture → specify → implement → verify → deliver → record")
	standard, err := prompt.AskConfirm("Use the standard feature flow?", true)
	if err != nil {
		return nil, err
	}
	if standard {
		cfg.FeatureFlow = defaultFeatureFlow()
	} else {
		custom, err := prompt.AskString("Custom stages (comma-separated)", strings.Join(defaultFeatureFlow(), ", "))
		if err != nil {
			return nil, err
		}
		cfg.FeatureFlow = parseFeatureFlow(custom)
	}

	optional := skills.Optional()
	if len(optional) > 0 {
		fmt.Println("\nOptional skills (not installed by default):")
		for _, s := range optional {
			fmt.Printf("  - %s — %s\n", s.Name, s.Summary)
		}
		for _, s := range optional {
			yes, confirmErr := prompt.AskConfirm(fmt.Sprintf("Install optional skill %s?", s.Name), false)
			if confirmErr != nil {
				return nil, confirmErr
			}
			if yes {
				cfg.OptionalSkills = append(cfg.OptionalSkills, s.Name)
			}
		}
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

// detectTicketTracker infers the ticket tracker from the git remote and the
// gh CLI. Only GitHub is auto-detectable; other trackers are asked.
func detectTicketTracker(cwd string) string {
	out, err := exec.Command("git", "-C", cwd, "remote", "get-url", "origin").Output()
	if err == nil && strings.Contains(string(out), "github.com") {
		if _, err := exec.LookPath("gh"); err == nil {
			return "github"
		}
	}
	return "none"
}

// defaultFeatureFlow is the standard six-stage pipeline the flow skill ships with.
func defaultFeatureFlow() []string {
	return []string{"capture", "specify", "implement", "verify", "deliver", "record"}
}

// trackerChoices returns the tracker options with the detected one first.
func trackerChoices(detected string) []string {
	all := []string{"github", "jira", "linear", "none"}
	choices := []string{}
	for _, c := range all {
		if c == detected {
			choices = append(choices, c)
		}
	}
	for _, c := range all {
		if c != detected {
			choices = append(choices, c)
		}
	}
	return choices
}

// parseFeatureFlow splits and normalizes a comma-separated stage list.
func parseFeatureFlow(s string) []string {
	var stages []string
	for _, part := range strings.Split(s, ",") {
		stage := strings.ToLower(strings.TrimSpace(part))
		if stage != "" {
			stages = append(stages, stage)
		}
	}
	if len(stages) == 0 {
		return defaultFeatureFlow()
	}
	return stages
}
