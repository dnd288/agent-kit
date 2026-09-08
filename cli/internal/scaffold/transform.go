package scaffold

import (
	"strings"

	"github.com/agent-kit/agent-kit-cli/internal/config"
)

// TransformContent replaces generic placeholders in kit content
// with project-specific values from the config.
func TransformContent(content string, cfg *config.ProjectConfig) string {
	// Replace project name references
	content = strings.ReplaceAll(content, "your-project-name", cfg.ProjectName)
	content = strings.ReplaceAll(content, "YOUR_PROJECT", strings.ToUpper(strings.ReplaceAll(cfg.Prefix, "-", "_")))

	// Replace package manager commands
	pm := cfg.PackageManager
	if pm == "" {
		pm = "npm"
	}

	pmRun := pm + " run"
	if pm == "pnpm" || pm == "bun" {
		pmRun = pm
	}

	content = strings.ReplaceAll(content, "your-package-manager run", pmRun)
	content = strings.ReplaceAll(content, "your-package-manager", pm)

	// Replace skill name prefixes in cross-references
	// Skills reference each other by name (e.g. "the `review` skill")
	// When installed, they get prefixed, so update the references
	coreSkills := []string{
		"review", "verify", "tdd", "bugfix", "refactor",
		"spec-workflow", "pr", "e2e", "test", "security",
		"component-development", "ui-development", "app-development",
		"api-contract", "state-management", "local-dev", "design",
	}

	for _, skill := range coreSkills {
		// Replace backtick-quoted skill references: `review` → `acme-review`
		content = strings.ReplaceAll(content,
			"`"+skill+"`",
			"`"+cfg.Prefix+"-"+skill+"`",
		)
		// Replace skill references in descriptions
		content = strings.ReplaceAll(content,
			"the "+skill+" skill",
			"the "+cfg.Prefix+"-"+skill+" skill",
		)
	}

	return content
}
