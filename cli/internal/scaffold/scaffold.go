package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/agent-kit/agent-kit-cli/internal/config"
)

// Scaffold copies kit content into the target project directory,
// transforming file paths and content based on the project config.
func Scaffold(kitFS embed.FS, cfg *config.ProjectConfig, targetDir string) error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Installing skills", func() error { return installSkills(kitFS, cfg, targetDir) }},
		{"Writing root templates", func() error { return installTemplates(kitFS, cfg, targetDir) }},
	}

	if cfg.IncludeClaude {
		steps = append(steps, struct {
			name string
			fn   func() error
		}{"Setting up Claude Code integration", func() error { return installClaude(kitFS, cfg, targetDir) }})
	}

	if cfg.IncludeOpenSpec {
		steps = append(steps, struct {
			name string
			fn   func() error
		}{"Installing OpenSpec schemas", func() error { return installOpenSpec(kitFS, cfg, targetDir) }})
	}

	if cfg.IncludeHooks {
		steps = append(steps, struct {
			name string
			fn   func() error
		}{"Installing git hooks", func() error { return installHooks(kitFS, cfg, targetDir) }})
	}

	if cfg.IncludeCI {
		steps = append(steps, struct {
			name string
			fn   func() error
		}{"Installing CI templates", func() error { return installCI(kitFS, cfg, targetDir) }})
	}

	for _, step := range steps {
		fmt.Printf("  → %s...\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

// installSkills copies all skills or a subset based on stack selection.
func installSkills(kitFS embed.FS, cfg *config.ProjectConfig, targetDir string) error {
	skillsDir := filepath.Join(targetDir, ".agents", "skills")
	claudeSkillsDir := filepath.Join(targetDir, ".claude", "skills")

	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(claudeSkillsDir, 0755); err != nil {
		return err
	}

	// Walk all skills in the kit
	err := fs.WalkDir(kitFS, "skills", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root "skills" directory itself
		if path == "skills" {
			return nil
		}

		// Get the relative path within skills/
		relPath := strings.TrimPrefix(path, "skills/")

		// Extract skill name (first path segment)
		parts := strings.SplitN(relPath, "/", 2)
		skillName := parts[0]

		// Check if this skill should be installed based on stack
		if !shouldInstallSkill(skillName, cfg.Stack) {
			if d.IsDir() && len(parts) == 1 {
				return fs.SkipDir
			}
			return nil
		}

		// Prefixed destination name
		prefixedName := cfg.Prefix + "-" + skillName
		destRel := prefixedName
		if len(parts) > 1 {
			destRel = filepath.Join(prefixedName, parts[1])
		}

		destPath := filepath.Join(skillsDir, destRel)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Read and transform content
		content, err := fs.ReadFile(kitFS, path)
		if err != nil {
			return err
		}

		transformed := TransformContent(string(content), cfg)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		if err := os.WriteFile(destPath, []byte(transformed), 0644); err != nil {
			return err
		}

		// Record installed skill
		cfg.AddSkill(skillName)

		return nil
	})

	if err != nil {
		return err
	}

	// Create symlinks from .claude/skills/ to .agents/skills/
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		linkPath := filepath.Join(claudeSkillsDir, entry.Name())
		targetPath, err := filepath.Rel(claudeSkillsDir, filepath.Join(skillsDir, entry.Name()))
		if err != nil {
			return err
		}
		// Remove existing symlink if present
		os.Remove(linkPath)
		if err := os.Symlink(targetPath, linkPath); err != nil {
			return fmt.Errorf("symlink %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// shouldInstallSkill determines if a skill is relevant for the project's stack.
func shouldInstallSkill(skillName string, stack []string) bool {
	// Core methodology skills are always installed
	coreSkills := map[string]bool{
		"review":        true,
		"verify":        true,
		"tdd":           true,
		"bugfix":        true,
		"refactor":      true,
		"spec-workflow": true,
		"pr":            true,
		"test":          true,
		"security":      true,
	}

	if coreSkills[skillName] {
		return true
	}

	// Stack-specific skills
	stackSkills := map[string][]string{
		"component-development": {"React", "Vue", "Svelte", "Angular"},
		"ui-development":        {"React", "Vue", "Svelte", "Angular"},
		"app-development":       {"Next.js", "Nuxt", "SvelteKit", "Remix"},
		"api-contract":          {"Fastify", "Express", "NestJS", "Hono"},
		"state-management":      {"React", "Vue", "Svelte", "Angular"},
		"e2e":                   {"Playwright", "Cypress"},
		"local-dev":             {},  // always install
		"design":                {},  // always install
	}

	required, exists := stackSkills[skillName]
	if !exists {
		return true // unknown skills are always installed
	}

	if len(required) == 0 {
		return true // no stack requirement
	}

	for _, s := range stack {
		for _, r := range required {
			if strings.EqualFold(s, r) {
				return true
			}
		}
	}

	return false
}

// installTemplates writes root-level instruction files.
func installTemplates(kitFS embed.FS, cfg *config.ProjectConfig, targetDir string) error {
	templates := map[string]string{
		"templates/AGENTS.md":       "AGENTS.md",
		"templates/CLAUDE.md":       "CLAUDE.md",
		"templates/CONTRIBUTING.md": "CONTRIBUTING.md",
		"templates/CONTEXT.md":      "CONTEXT.md",
	}

	for src, dest := range templates {
		content, err := fs.ReadFile(kitFS, src)
		if err != nil {
			return fmt.Errorf("read %s: %w", src, err)
		}

		transformed := TransformContent(string(content), cfg)

		// Generate skills table for AGENTS.md
		if dest == "AGENTS.md" {
			table := generateSkillsTable(cfg)
			transformed = strings.ReplaceAll(transformed, "<!-- SKILLS_TABLE -->", table)
		}

		destPath := filepath.Join(targetDir, dest)
		if err := os.WriteFile(destPath, []byte(transformed), 0644); err != nil {
			return fmt.Errorf("write %s: %w", dest, err)
		}
	}

	return nil
}

// generateSkillsTable builds a markdown table of installed skills.
func generateSkillsTable(cfg *config.ProjectConfig) string {
	skillDescriptions := map[string]string{
		"review":                "Five-axis code review gate",
		"verify":                "Adversarial claim verification against specifications",
		"tdd":                   "RED→GREEN discipline with mode-to-claim mapping",
		"bugfix":                "inspect→debug→impact→fix→verify loop",
		"refactor":              "MEASURE→PIN→MOVE→PROVE→RECORD loop",
		"spec-workflow":         "Typed changes with end-to-end-first ordering",
		"pr":                    "Five PR types with checklists and conventions",
		"e2e":                   "End-to-end test authoring and healing",
		"test":                  "Test mode taxonomy and operational runbook",
		"security":              "Surface-aware security review",
		"component-development": "Primitive component development patterns",
		"ui-development":        "Feature-folder composite UI patterns",
		"app-development":       "Application routing, data flow, access control",
		"api-contract":          "Schema-first API design and service purity",
		"state-management":      "State placement decision ladder",
		"local-dev":             "Local development environment setup",
		"design":                "Design system methodology",
	}

	var sb strings.Builder
	sb.WriteString("| Skill | Owns |\n")
	sb.WriteString("|---|---|\n")

	for _, skill := range cfg.InstalledSkills {
		desc := skillDescriptions[skill]
		if desc == "" {
			desc = skill
		}
		sb.WriteString(fmt.Sprintf("| `%s-%s` | %s |\n", cfg.Prefix, skill, desc))
	}

	return sb.String()
}

// installClaude sets up Claude Code agent definitions and workflows.
func installClaude(kitFS embed.FS, cfg *config.ProjectConfig, targetDir string) error {
	claudeDir := filepath.Join(targetDir, ".claude")

	// Copy agent definitions
	agentFiles := []string{"agents/reviewer.md", "agents/security-auditor.md", "agents/test-engineer.md"}
	for _, f := range agentFiles {
		content, err := fs.ReadFile(kitFS, "claude/"+f)
		if err != nil {
			return err
		}

		transformed := TransformContent(string(content), cfg)
		destPath := filepath.Join(claudeDir, f)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(destPath, []byte(transformed), 0644); err != nil {
			return err
		}
	}

	// Copy workflows
	workflowFiles := []string{"workflows/feature-wf.js"}
	for _, f := range workflowFiles {
		content, err := fs.ReadFile(kitFS, "claude/"+f)
		if err != nil {
			// Workflows are optional
			continue
		}

		transformed := TransformContent(string(content), cfg)
		destPath := filepath.Join(claudeDir, f)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(destPath, []byte(transformed), 0644); err != nil {
			return err
		}
	}

	// Copy commands
	err := fs.WalkDir(kitFS, "claude/commands", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		content, readErr := fs.ReadFile(kitFS, path)
		if readErr != nil {
			return readErr
		}

		relPath := strings.TrimPrefix(path, "claude/")
		destPath := filepath.Join(claudeDir, relPath)

		if mkErr := os.MkdirAll(filepath.Dir(destPath), 0755); mkErr != nil {
			return mkErr
		}
		transformed := TransformContent(string(content), cfg)
		return os.WriteFile(destPath, []byte(transformed), 0644)
	})

	return err
}

// installOpenSpec copies specification schemas.
func installOpenSpec(kitFS embed.FS, cfg *config.ProjectConfig, targetDir string) error {
	openspecDir := filepath.Join(targetDir, "openspec")

	return fs.WalkDir(kitFS, "openspec", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "openspec" {
			return nil
		}

		relPath := strings.TrimPrefix(path, "openspec/")
		destPath := filepath.Join(openspecDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		content, err := fs.ReadFile(kitFS, path)
		if err != nil {
			return err
		}

		transformed := TransformContent(string(content), cfg)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(destPath, []byte(transformed), 0644)
	})
}

// installHooks copies git hook templates.
func installHooks(kitFS embed.FS, cfg *config.ProjectConfig, targetDir string) error {
	huskyDir := filepath.Join(targetDir, ".husky")
	if err := os.MkdirAll(huskyDir, 0755); err != nil {
		return err
	}

	hooks := []string{"commit-msg", "pre-commit", "pre-push"}
	for _, hook := range hooks {
		content, err := fs.ReadFile(kitFS, "hooks/"+hook)
		if err != nil {
			continue // hooks are optional
		}

		transformed := TransformContent(string(content), cfg)
		destPath := filepath.Join(huskyDir, hook)

		if err := os.WriteFile(destPath, []byte(transformed), 0755); err != nil {
			return err
		}
	}

	return nil
}

// installCI copies CI workflow templates.
func installCI(kitFS embed.FS, cfg *config.ProjectConfig, targetDir string) error {
	ciDir := filepath.Join(targetDir, ".github", "workflows")
	if err := os.MkdirAll(ciDir, 0755); err != nil {
		return err
	}

	return fs.WalkDir(kitFS, "ci", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		content, err := fs.ReadFile(kitFS, path)
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, "ci/")
		transformed := TransformContent(string(content), cfg)
		destPath := filepath.Join(ciDir, relPath)

		return os.WriteFile(destPath, []byte(transformed), 0644)
	})
}
