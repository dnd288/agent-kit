package cmd

import (
	"fmt"
	"io/fs"

	"github.com/agent-kit/agent-kit-cli/internal/skills"
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
	discovered, err := discoverSkills()
	if err != nil {
		return err
	}

	fmt.Println("Available skills:")
	fmt.Println()

	categories := map[string][]skillInfo{}
	order := skills.CategoryOrder

	for _, s := range discovered {
		categories[s.category] = append(categories[s.category], s)
	}

	for _, cat := range order {
		items, ok := categories[cat]
		if !ok {
			continue
		}
		fmt.Printf("  %s:\n", cat)
		for _, s := range items {
			name := s.name
			if skills.IsMeta(name) {
				name += " (meta)"
			} else if skills.IsOptional(name) {
				name += " (optional)"
			}
			fmt.Printf("    %-25s  %s\n", name, s.description)
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
	return skills.Category(name)
}

func extractDescription(fsys fs.FS, path string) string {
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return "(no description)"
	}

	fm, ok := parseFrontmatter(string(data))
	if !ok {
		return "(no description)"
	}

	desc := frontmatterField(fm, "description")
	if desc == "" {
		return "(no description)"
	}
	if len(desc) > 80 {
		desc = desc[:77] + "..."
	}
	return desc
}
