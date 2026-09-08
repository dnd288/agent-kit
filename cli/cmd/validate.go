package cmd

import (
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"

	"github.com/agent-kit/agent-kit-cli/internal/skills"
	kit "github.com/agent-kit/agent-kit-cli/kit"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the embedded skill kit and CLI wiring",
	Long: `Validate checks the skills shipped inside the CLI:

  - every skill directory has a SKILL.md with a non-empty name and description
  - the frontmatter name matches the directory name
  - every skill is registered in internal/skills, so list categorises it and
    add/init prefix its cross-references instead of silently filing it as "Other"
  - cross-skill links (../<name>/SKILL.md) resolve to a real skill
  - reference links (references/<file>.md) point to files that exist

It exits non-zero if any check fails.`,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

var (
	crossRefRe = regexp.MustCompile(`\.\./([a-z0-9-]+)/SKILL\.md`)
	refFileRe  = regexp.MustCompile(`(?:\.\./[a-z0-9-]+/)?references/[A-Za-z0-9._-]+\.md`)
)

func runValidate(cmd *cobra.Command, args []string) error {
	var problems []string
	report := func(format string, a ...interface{}) {
		problems = append(problems, fmt.Sprintf(format, a...))
	}

	entries, err := fs.ReadDir(kit.Content, "skills")
	if err != nil {
		return fmt.Errorf("reading embedded skills: %w", err)
	}

	// Pass 1: collect skill directories and their SKILL.md content. dirs must be
	// complete before cross-references can be resolved in pass 2.
	dirs := map[string]bool{}
	content := map[string]string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		dirs[name] = true

		data, err := fs.ReadFile(kit.Content, "skills/"+name+"/SKILL.md")
		if err != nil {
			report("%s: directory has no SKILL.md", name)
			continue
		}
		content[name] = string(data)
	}

	// Pass 2: per-skill checks.
	for name, body := range content {
		fm, ok := parseFrontmatter(body)
		if !ok {
			report("%s/SKILL.md: missing or unterminated YAML frontmatter", name)
			continue
		}

		fmName := frontmatterField(fm, "name")
		switch {
		case fmName == "":
			report("%s/SKILL.md: frontmatter has no non-empty `name`", name)
		case fmName != name:
			report("%s/SKILL.md: frontmatter name %q does not match its directory", name, fmName)
		}
		if frontmatterField(fm, "description") == "" {
			report("%s/SKILL.md: frontmatter has no non-empty `description`", name)
		}

		if skills.Category(name) == "Other" {
			report("%s: not registered in internal/skills — `list` would file it under \"Other\" and `add`/`init` would not prefix its cross-references", name)
		}

		for _, m := range crossRefRe.FindAllStringSubmatch(body, -1) {
			ref := m[1]
			if ref != name && !dirs[ref] {
				report("%s/SKILL.md: cross-reference ../%s/SKILL.md does not resolve to a skill", name, ref)
			}
		}

		for _, rel := range refFileRe.FindAllString(body, -1) {
			target := path.Clean(path.Join("skills", name, rel))
			if _, err := fs.ReadFile(kit.Content, target); err != nil {
				report("%s/SKILL.md: reference link %q points to a missing file", name, rel)
			}
		}
	}

	// Pass 3: registry entries that have no directory (the reverse drift).
	for _, s := range skills.All {
		if !dirs[s.Name] {
			report("%s: registered in internal/skills but has no skills/%s directory", s.Name, s.Name)
		}
	}

	sort.Strings(problems)
	if len(problems) > 0 {
		fmt.Printf("✗ %d problem(s) found:\n\n", len(problems))
		for _, p := range problems {
			fmt.Printf("  - %s\n", p)
		}
		return fmt.Errorf("validation failed: %d problem(s)", len(problems))
	}

	fmt.Printf("✓ kit valid: %d skills, all registered, frontmatter and references OK\n", len(content))
	return nil
}
