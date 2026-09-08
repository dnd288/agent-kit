# Agent Kit

A reusable agent-driven development methodology — skills, templates, specifications, and tooling — extracted from production use and ready to adopt in any software project.

## What is this?

Agent Kit is a collection of **process skills**, **instruction templates**, and a **CLI scaffolder** that brings structured agent-assisted development to your codebase. It encodes battle-tested practices for:

- **Code review** — five-axis review beyond what linters catch
- **Verification** — adversarial claim-testing against specifications
- **Test-driven development** — RED→GREEN with mode-to-claim mapping
- **Bug fixing** — structured inspect→debug→impact→fix→verify loop
- **Refactoring** — MEASURE→PIN→MOVE→PROVE→RECORD with characterisation harnesses
- **Specification workflow** — typed changes with end-to-end ordering
- **Pull requests** — typed PRs with checklists and conventions
- **Security** — surface-aware security review
- And more: component development, UI development, API contracts, state management, E2E testing

## Philosophy

**Compound engineering: the codebase gets simpler as the product gets larger.**

1. Leave it easier to change than you found it
2. A new abstraction earns its place by deleting
3. Two is a coincidence; three is a missing abstraction — and the abstraction replaces all three
4. Delete the special case; do not flag it
5. An exception is written down and bounded, or it becomes the rule

See [kit/docs/philosophy.md](kit/docs/philosophy.md) for the full treatment.

## Quick Start

### Install the CLI

```bash
# From source
cd cli && go build -o agent-kit && mv agent-kit /usr/local/bin/

# Or with go install
go install github.com/your-org/agent-kit/cli@latest
```

### Initialize a project

```bash
cd your-project
agent-kit init
```

The CLI will ask:
- **Project name** — used in instructions and agent definitions
- **Project prefix** — short prefix for skill names (e.g. "acme" → `acme-review`)
- **Stack** — which technologies you use (filters relevant skills)
- **Package manager** — pnpm / npm / yarn / bun
- **Monorepo?** — affects path references
- What to include: OpenSpec schemas, Claude Code integration, git hooks, CI templates

### Add a skill later

```bash
agent-kit add security
agent-kit add e2e
```

### List available skills

```bash
agent-kit list
```

## What gets generated

```
your-project/
├── AGENTS.md                    # Golden rules + project boundaries
├── CLAUDE.md                    # Loads AGENTS.md
├── CONTRIBUTING.md              # Development workflow
├── CONTEXT.md                   # Project glossary
├── .agents/skills/              # Tool-neutral skills
│   ├── <prefix>-review/
│   ├── <prefix>-verify/
│   ├── <prefix>-tdd/
│   └── ...
├── .claude/                     # Claude Code integration
│   ├── skills/                  # Symlinks → .agents/skills/
│   ├── agents/                  # Agent definitions
│   └── workflows/               # Multi-agent orchestration
├── openspec/                    # Spec-driven development (optional)
│   ├── config.yaml
│   └── schemas/
├── .husky/                      # Git hooks (optional)
└── .github/workflows/           # CI templates (optional)
```

## Structure

```
agent-kit/
├── kit/                         # The methodology content
│   ├── skills/                  # 15+ generic process skills
│   ├── templates/               # Root instruction templates
│   ├── claude/                  # Claude Code integration
│   ├── openspec/                # Specification schemas
│   ├── hooks/                   # Git hook templates
│   ├── ci/                      # CI workflow templates
│   └── docs/                    # Methodology documentation
├── cli/                         # Go CLI scaffolder
└── examples/                    # Example output
```

## Skills

Skills follow the [Agent Skills format](kit/docs/skill-format.md) — YAML frontmatter with trigger descriptions, plus a markdown body encoding the process. They are tool-neutral: any agent platform that reads `SKILL.md` picks them up.

| Skill | Owns |
|---|---|
| `review` | Five-axis code review gate |
| `verify` | Adversarial CLAIM→EXTRACT→DOUBT→RECONCILE→STOP verification |
| `tdd` | RED→GREEN discipline with mode-to-claim mapping |
| `bugfix` | inspect→debug→impact→fix→verify loop |
| `refactor` | MEASURE→PIN→MOVE→PROVE→RECORD loop |
| `spec-workflow` | Typed changes with end-to-end-first ordering |
| `pr` | Five PR types with checklists |
| `e2e` | End-to-end test authoring and healing |
| `test` | Test mode taxonomy and operational runbook |
| `security` | Surface-aware security review |
| `component-development` | Primitive component patterns |
| `ui-development` | Feature-folder composite patterns |
| `app-development` | Application routing and data flow |
| `api-contract` | Schema-first API design |
| `state-management` | State placement decision ladder |

## License

MIT
