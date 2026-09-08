package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	kit "github.com/agent-kit/agent-kit-cli/kit"
	"github.com/agent-kit/agent-kit-cli/internal/config"
	"github.com/agent-kit/agent-kit-cli/internal/scaffold"
	"github.com/agent-kit/agent-kit-cli/internal/skills"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <skill>",
	Short: "Add a single skill to an existing project",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	skillName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("could not read agent-kit.yaml; run `agent-kit init` first: %w", err)
	}

	if cfg.HasSkill(skillName) {
		fmt.Printf("Skill %s is already installed.\n", skillName)
		return nil
	}

	if !skillExists(skillName) {
		return fmt.Errorf("unknown skill %q; run `agent-kit list` to see available skills", skillName)
	}

	if skills.IsMeta(skillName) {
		return fmt.Errorf("skill %q is a meta-skill and cannot be installed into projects; use it from the kit source tree", skillName)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := installSingleSkill(cfg, cwd, skillName); err != nil {
		return err
	}

	cfg.AddSkill(skillName)
	if err := cfg.Save(cwd); err != nil {
		return err
	}

	fmt.Printf("Installed skill %s-%s.\n", cfg.Prefix, skillName)
	return nil
}

func skillExists(skillName string) bool {
	_, err := fs.Stat(kit.Content, "skills/"+skillName+"/SKILL.md")
	return err == nil
}

func installSingleSkill(cfg *config.ProjectConfig, targetDir, skillName string) error {
	skillsDir := filepath.Join(targetDir, ".agents", "skills")
	claudeSkillsDir := filepath.Join(targetDir, ".claude", "skills")
	prefixedName := cfg.Prefix + "-" + skillName

	if err := os.MkdirAll(filepath.Join(skillsDir, prefixedName), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(claudeSkillsDir, 0755); err != nil {
		return err
	}

	root := "skills/" + skillName
	err := fs.WalkDir(kit.Content, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(path, root)
		relPath = strings.TrimPrefix(relPath, "/")
		destPath := filepath.Join(skillsDir, prefixedName, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		content, err := fs.ReadFile(kit.Content, path)
		if err != nil {
			return err
		}

		transformed := scaffold.TransformContent(string(content), cfg)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(destPath, []byte(transformed), 0644)
	})
	if err != nil {
		return err
	}

	linkPath := filepath.Join(claudeSkillsDir, prefixedName)
	targetPath, err := filepath.Rel(claudeSkillsDir, filepath.Join(skillsDir, prefixedName))
	if err != nil {
		return err
	}
	os.Remove(linkPath)
	return os.Symlink(targetPath, linkPath)
}
