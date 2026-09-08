// Package skills is the single source of truth for which skills the kit ships,
// how they are categorised, whether they are opt-in, and their one-line
// summaries. Everything that needs to know about the set of skills reads from
// here:
//
//   - `list`   — categorisation, display order, and the optional marker
//   - `add`/`init` (scaffold.TransformContent) — cross-reference prefixing
//   - `init`   — which skills are optional (opt-in) and their summaries
//   - scaffold — the AGENTS.md skills table (summaries)
//   - `validate` — checks this registry against the embedded skills/ directory
//
// Adding a skill is one edit: add its directory under kit/skills/<name>/ and a
// line to All below.
package skills

// Skill is one methodology skill shipped in the kit.
type Skill struct {
	Name string
	// Category groups the skill in `list` output.
	Category string
	// Optional skills are not installed by default; `init` asks before
	// installing them and they can be added later with `agent-kit add`.
	Optional bool
	// Meta skills are used from the kit source tree and are never scaffolded into
	// target projects.
	Meta bool
	// Summary is a one-line human description for the AGENTS.md table and the
	// interactive prompts.
	Summary string
}

// All is the canonical registry of shipped skills.
var All = []Skill{
	// Core Methodology
	{Name: "review", Category: "Core Methodology", Summary: "Five-axis code review gate"},
	{Name: "verify", Category: "Core Methodology", Summary: "Adversarial claim verification against specifications"},
	{Name: "tdd", Category: "Core Methodology", Summary: "RED→GREEN discipline with mode-to-claim mapping"},
	{Name: "bugfix", Category: "Core Methodology", Summary: "inspect→debug→impact→fix→verify loop"},
	{Name: "refactor", Category: "Core Methodology", Summary: "MEASURE→PIN→MOVE→PROVE→RECORD loop"},
	{Name: "simplicity", Category: "Core Methodology", Summary: "Simplest-sufficient-shape decision — the rung ladder"},
	{Name: "spec-workflow", Category: "Core Methodology", Summary: "Typed changes with end-to-end-first ordering"},
	{Name: "pr", Category: "Core Methodology", Summary: "Five PR types with checklists and conventions"},
	{Name: "optimization-loop", Category: "Core Methodology", Optional: true, Summary: "BASELINE→CHANGE→MEASURE→KEEP-OR-REVERT performance loop"},
	{Name: "problem-solving", Category: "Core Methodology", Optional: true, Summary: "Impasse techniques and the dev-note handed to the user"},

	// Development
	{Name: "component-development", Category: "Development", Summary: "Primitive component development patterns"},
	{Name: "ui-development", Category: "Development", Summary: "Feature-folder composite UI patterns"},
	{Name: "app-development", Category: "Development", Summary: "Application routing, data flow, access control"},
	{Name: "api-contract", Category: "Development", Summary: "Schema-first API design and service purity"},
	{Name: "state-management", Category: "Development", Summary: "State placement decision ladder"},
	{Name: "design", Category: "Development", Summary: "Design system methodology"},

	// Testing
	{Name: "e2e", Category: "Testing", Summary: "End-to-end test authoring and healing"},
	{Name: "test", Category: "Testing", Summary: "Test mode taxonomy and operational runbook"},
	{Name: "scenario-explorer", Category: "Testing", Optional: true, Summary: "Enumerate the case space before writing tests"},

	// Operations
	{Name: "security", Category: "Operations", Summary: "Surface-aware security review"},
	{Name: "local-dev", Category: "Operations", Summary: "Local development environment setup"},

	// Setup
	{Name: "agent-kit-init", Category: "Setup", Optional: true, Meta: true, Summary: "Initialize Agent Kit via any coding agent"},
}

// CategoryOrder is the order categories are displayed by `list`.
var CategoryOrder = []string{"Core Methodology", "Development", "Testing", "Operations", "Setup"}

// find returns the registered skill and whether it exists.
func find(name string) (Skill, bool) {
	for _, s := range All {
		if s.Name == name {
			return s, true
		}
	}
	return Skill{}, false
}

// Category returns the category for a skill name, or "Other" if it is not
// registered.
func Category(name string) string {
	if s, ok := find(name); ok {
		return s.Category
	}
	return "Other"
}

// Summary returns the one-line summary for a skill name, or "" if it is not
// registered.
func Summary(name string) string {
	if s, ok := find(name); ok {
		return s.Summary
	}
	return ""
}

// IsOptional reports whether a skill is opt-in.
func IsOptional(name string) bool {
	s, ok := find(name)
	return ok && s.Optional
}

// IsMeta reports whether a skill is used from the kit and not installed into projects.
func IsMeta(name string) bool {
	s, ok := find(name)
	return ok && s.Meta
}

// Names returns every installable (non-meta) skill name, in registry order.
// Used by scaffold.TransformContent to prefix cross-skill references.
func Names() []string {
	var names []string
	for _, s := range All {
		if !s.Meta {
			names = append(names, s.Name)
		}
	}
	return names
}

// Optional returns the opt-in skills that are installable (not meta), in registry order.
func Optional() []Skill {
	var out []Skill
	for _, s := range All {
		if s.Optional && !s.Meta {
			out = append(out, s)
		}
	}
	return out
}
