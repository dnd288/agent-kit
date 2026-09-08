package cmd

import (
	"io/fs"
	"testing"

	kit "github.com/agent-kit/agent-kit-cli/kit"
)

func TestAgentKitInitSkillExists(t *testing.T) {
	_, err := fs.ReadFile(kit.Content, "skills/agent-kit-init/SKILL.md")
	if err != nil {
		t.Fatalf("skills/agent-kit-init/SKILL.md not found in embedded FS: %v", err)
	}
}

func TestAgentKitInitFrontmatterValid(t *testing.T) {
	data, err := fs.ReadFile(kit.Content, "skills/agent-kit-init/SKILL.md")
	if err != nil {
		t.Fatalf("read SKILL.md: %v", err)
	}

	fm, ok := parseFrontmatter(string(data))
	if !ok {
		t.Fatal("SKILL.md has no valid YAML frontmatter")
	}

	name := frontmatterField(fm, "name")
	if name != "agent-kit-init" {
		t.Fatalf("frontmatter name = %q, want \"agent-kit-init\"", name)
	}

	desc := frontmatterField(fm, "description")
	if desc == "" {
		t.Fatal("frontmatter description is empty")
	}
}
