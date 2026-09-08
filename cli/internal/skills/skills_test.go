package skills

import "testing"

func TestAgentKitInitRegistered(t *testing.T) {
	for _, s := range All {
		if s.Name == "agent-kit-init" {
			return
		}
	}
	t.Fatal("agent-kit-init not found in All registry")
}

func TestAgentKitInitCategory(t *testing.T) {
	got := Category("agent-kit-init")
	if got != "Setup" {
		t.Fatalf("Category(\"agent-kit-init\") = %q, want \"Setup\"", got)
	}
}

func TestAgentKitInitIsOptional(t *testing.T) {
	if !IsOptional("agent-kit-init") {
		t.Fatal("IsOptional(\"agent-kit-init\") = false, want true")
	}
}

func TestAgentKitInitIsMeta(t *testing.T) {
	if !IsMeta("agent-kit-init") {
		t.Fatal("IsMeta(\"agent-kit-init\") = false, want true")
	}
}

func TestMetaSkillsExcludedFromOptional(t *testing.T) {
	for _, s := range Optional() {
		if s.Meta {
			t.Errorf("Optional() returned meta skill %q", s.Name)
		}
	}
}

func TestMetaSkillsExcludedFromNames(t *testing.T) {
	names := Names()
	for _, n := range names {
		if IsMeta(n) {
			t.Errorf("Names() returned meta skill %q", n)
		}
	}
}

func TestAgentKitInitSummary(t *testing.T) {
	got := Summary("agent-kit-init")
	if got == "" {
		t.Fatal("Summary(\"agent-kit-init\") is empty")
	}
}

func TestSetupCategoryInOrder(t *testing.T) {
	for _, cat := range CategoryOrder {
		if cat == "Setup" {
			return
		}
	}
	t.Fatal("\"Setup\" not found in CategoryOrder")
}

func TestAllSkillsHaveCategory(t *testing.T) {
	valid := make(map[string]bool, len(CategoryOrder))
	for _, cat := range CategoryOrder {
		valid[cat] = true
	}
	for _, s := range All {
		if !valid[s.Category] {
			t.Errorf("skill %q has category %q which is not in CategoryOrder", s.Name, s.Category)
		}
	}
}

func TestAllSkillsHaveSummary(t *testing.T) {
	for _, s := range All {
		if s.Summary == "" {
			t.Errorf("skill %q has empty Summary", s.Name)
		}
	}
}

func TestNoDuplicateSkillNames(t *testing.T) {
	seen := make(map[string]bool, len(All))
	for _, s := range All {
		if seen[s.Name] {
			t.Errorf("duplicate skill name: %q", s.Name)
		}
		seen[s.Name] = true
	}
}
