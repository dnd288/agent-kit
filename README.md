# Agent Kit

A reusable agent-driven development methodology — skills, templates, specifications, and tooling — extracted from production use and ready to adopt in any software project.

## What is this?

Agent Kit is a collection of **process skills**, **instruction templates**, and an **initialization skill** that brings structured agent-assisted development to your codebase. It works with any coding agent — Claude Code, OMP, Codex, OpenCode, or anything that reads markdown instructions. It encodes battle-tested practices for:

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

See [cli/kit/docs/philosophy.md](cli/kit/docs/philosophy.md) for the full treatment.

## Quick Start

### Option A — via any coding agent (recommended)

No binary installation needed. Any coding agent that can read files and run shell
commands can initialize Agent Kit in your project.

**Step 1 — get the kit into your project:**

```bash
# Clone and copy the kit directory
git clone https://github.com/your-org/agent-kit.git /tmp/agent-kit
cp -r /tmp/agent-kit/cli/kit /path/to/your-project/.agent-kit-source
```

Or add as a git submodule:

```bash
git submodule add https://github.com/your-org/agent-kit.git .agent-kit
```

**Step 2 — ask your agent to initialize:**

Tell your coding agent:

> Read the skill at `.agent-kit-source/skills/agent-kit-init/SKILL.md` and follow
> its procedure to initialize agent-kit in this project.

The agent will:
1. **Detect** your stack, package manager, and monorepo structure from the codebase
2. **Confirm** the configuration with you
3. **Scaffold** skills, templates, hooks, CI, and specifications
4. **Report** what was installed and next steps

**Step 3 — clean up the source (optional):**

After initialization, the kit source is no longer needed — all assets have been
copied and transformed into your project:

```bash
rm -rf .agent-kit-source   # if you copied it
# or remove the submodule if you used that approach
```

#### Agent-specific examples

<details>
<summary><strong>Claude Code</strong></summary>

```
# In Claude Code, just say:
Read .agent-kit-source/skills/agent-kit-init/SKILL.md and initialize agent-kit in this project.
```

Or if you've set up the skill as a slash command:

```
/agent-kit-init
```

</details>

<details>
<summary><strong>Codex / OpenAI Codex CLI</strong></summary>

```bash
codex "Read the skill file at .agent-kit-source/skills/agent-kit-init/SKILL.md and follow its procedure to set up agent-kit in this project."
```

</details>

<details>
<summary><strong>OMP / OpenCode / other agents</strong></summary>

Point your agent to the skill file and ask it to follow the procedure:

```
Please read .agent-kit-source/skills/agent-kit-init/SKILL.md and execute
the initialization procedure for this project.
```

Any agent that can read markdown, scan a directory tree, and write files can
run the init skill.

</details>

### Option B — via the CLI

A Go CLI is also available for scripted or non-agent workflows.

```bash
# From source
cd cli && go build -o agent-kit && mv agent-kit /usr/local/bin/

# Or with go install
go install github.com/your-org/agent-kit/cli@latest
```

```bash
cd your-project
agent-kit init
```

The CLI will interactively ask for project name, prefix, stack, package manager,
and which modules to include.

### Add a skill later

Via an agent — ask it to read the `agent-kit-init` skill's add procedure, or:

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
├── agent-kit.yaml               # Configuration (prefix, stack, installed skills)
├── .agents/skills/              # Tool-neutral skills
│   ├── <prefix>-review/
│   ├── <prefix>-verify/
│   ├── <prefix>-tdd/
│   └── ...
├── .claude/                     # Claude Code integration (optional)
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
├── cli/
│   ├── kit/                     # The methodology content
│   │   ├── skills/              # 21 process skills + init meta-skill
│   │   ├── templates/           # Root instruction templates
│   │   ├── claude/              # Claude Code integration
│   │   ├── openspec/            # Specification schemas
│   │   ├── hooks/               # Git hook templates
│   │   ├── ci/                  # CI workflow templates
│   │   └── docs/                # Methodology documentation
│   ├── cmd/                     # CLI commands (Go)
│   └── internal/                # CLI internals
└── examples/                    # Example output
```

## Skills

Skills follow the [Agent Skills format](cli/kit/docs/skill-format.md) — YAML frontmatter with trigger descriptions, plus a markdown body encoding the process. They are **tool-neutral**: any agent platform that reads `SKILL.md` picks them up.

### Core Methodology

| Skill | Owns |
|---|---|
| `review` | Five-axis code review gate |
| `verify` | Adversarial CLAIM→EXTRACT→DOUBT→RECONCILE→STOP verification |
| `tdd` | RED→GREEN discipline with mode-to-claim mapping |
| `bugfix` | inspect→debug→impact→fix→verify loop |
| `refactor` | MEASURE→PIN→MOVE→PROVE→RECORD loop |
| `simplicity` | Simplest-sufficient-shape decision — the rung ladder |
| `spec-workflow` | Typed changes with end-to-end-first ordering |
| `pr` | Five PR types with checklists |
| `optimization-loop` | BASELINE→CHANGE→MEASURE→KEEP-OR-REVERT performance loop *(optional)* |
| `problem-solving` | Impasse techniques and the dev-note handed to the user *(optional)* |

### Development

| Skill | Owns |
|---|---|
| `component-development` | Primitive component patterns |
| `ui-development` | Feature-folder composite patterns |
| `app-development` | Application routing and data flow |
| `api-contract` | Schema-first API design |
| `state-management` | State placement decision ladder |
| `design` | Design system methodology |

### Testing

| Skill | Owns |
|---|---|
| `e2e` | End-to-end test authoring and healing |
| `test` | Test mode taxonomy and operational runbook |
| `scenario-explorer` | Enumerate the case space before writing tests *(optional)* |

### Operations

| Skill | Owns |
|---|---|
| `security` | Surface-aware security review |
| `local-dev` | Local development environment setup |

### Setup (meta)

| Skill | Owns |
|---|---|
| `agent-kit-init` | Initialize Agent Kit via any coding agent |

> Meta skills are used from the kit source tree and are not scaffolded into
> target projects.

## Using skills after initialization

Once Agent Kit is initialized in your project, the installed skills are available
under `.agents/skills/<prefix>-<skill-name>/SKILL.md`. Tell your coding agent to
load and follow a skill when performing the corresponding task:

```
# Code review
Load the <prefix>-review skill and review the current diff.

# Bug fix
Load the <prefix>-bugfix skill and fix the issue described in #123.

# TDD
Load the <prefix>-tdd skill and implement the feature described in the spec.

# Security review
Load the <prefix>-security skill and audit the changes on this branch.

# Create a PR
Load the <prefix>-pr skill and open a pull request for this branch.
```

For Claude Code, skills in `.claude/skills/` are auto-discovered — the agent
loads them based on the trigger text in the `description` frontmatter field.

For other agents, point them to the `.agents/skills/` directory or the specific
skill file.

## License

MIT
