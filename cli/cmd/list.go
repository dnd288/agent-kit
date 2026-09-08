package cmd

import (
	"fmt"
	"io/fs"
	"strings"

	kit "github.com/agent-kit/agent-kit-cli/kit"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available skills",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

type skillInfo struct {
	name        string
	description string
	category    string
}

func runList(cmd *cobra.Command, args []string) error {
	skills, err := discoverSkills()
	if err != nil {
		return err
	}

	fmt.Println("Available skills:")
	fmt.Println()

	categories := map[string][]skillInfo{}
	order := []string{"Core Methodology", "Development", "Testing", "Operations"}

	for _, s := range skills {
		categories[s.category] = append(categories[s.category], s)
	}

	for _, cat := range order {
		items, ok := categories[cat]
		if !ok {
			continue
		}
		fmt.Printf("  %s:\n", cat)
		for _, s := range items {
			fmt.Printf("    %-25s  %s\n", s.name, s.description)
		}
		fmt.Println()
	}

	return nil
}

func discoverSkills() ([]skillInfo, error) {
	var skills []skillInfo

	entries, err := fs.ReadDir(kit.Content, "skills")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		desc := extractDescription(kit.Content, "skills/"+name+"/SKILL.md")
		cat := categorizeSkill(name)

		skills = append(skills, skillInfo{
			name:        name,
			description: desc,
			category:    cat,
		})
	}

	return skills, nil
}

func categorizeSkill(name string) string {
	switch name {
	case "review", "verify", "tdd", "bugfix", "refactor", "spec-workflow", "pr":
		return "Core Methodology"
	case "component-development", "ui-development", "app-development", "api-contract", "state-management", "design":
		return "Development"
	case "e2e", "test":
		return "Testing"
	case "security", "local-dev":
		return "Operations"
	default:
		return "Other"
	}
}

func extractDescription(fsys fs.FS, path string) string {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return "(no description)"
	}

	content := string(data)

	// Look for YAML frontmatter description
	if strings.HasPrefix(content, "---") {
		end := strings.Index(content[3:], "---")
		if end > 0 {
			frontmatter := content[3 : end+3]
			for _, line := range strings.Split(frontmatter, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "description:") {
					desc := strings.TrimPrefix(line, "description:")
					desc = strings.TrimSpace(desc)
					desc = strings.Trim(desc, `"'`)
					// Truncate long descriptions
					if len(desc) > 80 {
						desc = desc[:77] + "..."
					}
					return desc
				}
			}
		}
	}

	return "(no description)"
}
