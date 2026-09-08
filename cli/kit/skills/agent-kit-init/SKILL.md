---
name: agent-kit-init
description: "Initialize Agent Kit in a project — scan the codebase to detect stack, package manager, and structure; confirm with the user; then scaffold skills, templates, hooks, CI, and specifications. Use this skill when setting up agent-driven development practices in a new or existing repository."
---

# Initialize Agent Kit

This skill turns a repository into an agent-kit project. It detects the project's
stack, proposes a configuration, confirms with the user, then scaffolds skills,
templates, and tooling from the bundled kit.

The procedure is deterministic: two runs with the same config produce identical
output. The agent replaces interactive prompts — it reads the codebase, proposes,
and confirms instead of running a terminal questionnaire.

## When to use this skill

- A repository has no `.agents/skills/` directory and the user wants agent-driven
  development practices.
- The user asks to "set up agent-kit", "initialize agent-kit", or "add development
  skills to this project".
- The user wants structured code review, TDD, spec-workflow, security review, or
  similar methodology installed.

Do not use this skill when:

- Agent-kit is already installed (`agent-kit.yaml` exists at the project root).
  Use the `add` procedure to install additional skills instead.
- The user only wants a single skill added to an existing agent-kit project.

## Read first

Before starting, read these references shipped alongside this skill:

- `references/config-schema.md` — every configuration field, its type, default,
  and how to detect it from the codebase.
- `references/skill-registry.md` — the full skill catalog with categories,
  stack requirements, and optional flags.
- `references/transform-rules.md` — placeholder replacement rules applied to
  every file copied from the kit.

The kit assets directory — containing all skills, templates, and tooling to be
scaffolded — is located at `../../` relative to this skill's directory. That is
the root of the kit tree with the following layout:

```
../../                          (kit root)
├── skills/                     # Methodology skills to copy
├── templates/                  # Root instruction templates
├── claude/                     # Claude Code agents, commands, workflows
├── openspec/                   # Specification schemas
├── hooks/                      # Git hook templates
├── ci/                         # CI workflow templates
└── docs/                       # Methodology documentation
```

All `kit/` paths referenced in this skill resolve from that root. For example,
`kit/skills/review/` means `../../skills/review/` relative to this skill's
directory.

## Procedure

### Phase 1 — Detect

Scan the project root to infer configuration. Do not ask the user yet.

**Project name:**
- Read the `name` field from `package.json`, the module path from `go.mod`,
  the package name from `Cargo.toml` or `pyproject.toml`.
- Fallback: use the directory name.

**Prefix:**
- Derive from the project name. Single word → use it (truncated to 8 chars).
  Multiple words → use initials. Lowercase, hyphens only.
- Examples: `my-app` → `my-app`, `acme-web-platform` → `awp`.

**Package manager:**
- `pnpm-lock.yaml` or `pnpm-workspace.yaml` → `pnpm`
- `bun.lockb` or `bun.lock` → `bun`
- `yarn.lock` → `yarn`
- `package-lock.json` → `npm`
- No lockfile → `npm` (default)

**Stack (scan dependencies from package.json, go.mod, Cargo.toml, etc.):**
- `react` or `react-dom` → React
- `next` → Next.js
- `vue` → Vue
- `svelte` or `@sveltejs/kit` → Svelte
- `fastify` → Fastify
- `express` → Express
- `@nestjs/core` → NestJS
- `prisma` or `@prisma/client` → Prisma
- `@playwright/test` or `playwright` → Playwright
- `cypress` → Cypress
- For non-JS projects, detect from equivalent dependency files.

**Monorepo:**
- `pnpm-workspace.yaml`, `lerna.json`, `nx.json`, or `workspaces` field in
  `package.json` → monorepo.

**Existing setup:**
- Check if `.agents/`, `.claude/`, `AGENTS.md`, `agent-kit.yaml` already exist.
  If `agent-kit.yaml` exists, stop and tell the user agent-kit is already
  initialized.

### Phase 2 — Confirm

Present the detected configuration to the user as a structured summary:

```
Detected configuration:
  Project name:    my-awesome-app
  Prefix:          maa
  Package manager: pnpm
  Stack:           React, Next.js, Playwright
  Monorepo:        yes

Modules to install:
  ✓ Skills (15 core + 2 stack-matched)
  ✓ Root templates (AGENTS.md, CLAUDE.md, CONTRIBUTING.md, CONTEXT.md)
  ? OpenSpec schemas          — structured specification workflow
  ? Claude Code integration   — agents, commands, workflows for Claude Code
  ? Git hooks                 — pre-commit, commit-msg, pre-push via .husky
  ? CI templates              — GitHub Actions workflows
  ? Optional skills:
      optimization-loop       — performance optimization loop
      problem-solving         — impasse techniques
      scenario-explorer       — enumerate the case space before writing tests
```

Ask the user to:
1. Confirm or change the detected values (project name, prefix, package manager).
2. Decide on each optional module (OpenSpec, Claude Code, hooks, CI).
3. Decide on each optional skill.

If the user says "use defaults" or "yes to all", enable everything.

### Phase 3 — Scaffold

Execute these steps in order. For every file copied from the kit, apply the
transform rules from `references/transform-rules.md`.

The kit root is at `../../` relative to this skill's directory. All source paths
below are relative to that kit root.

#### 3.1 — Install skills

For each skill in the registry (see `references/skill-registry.md`):

1. **Check if it should be installed:**
   - Core skills → always install.
   - Stack-specific skills → install only if the detected stack matches.
   - Optional skills → install only if the user opted in.

2. **Copy the skill directory** from the kit root:
   ```
   ../../skills/<skill-name>/  →  <project>/.agents/skills/<prefix>-<skill-name>/
   ```
   Apply transforms to every `.md` file.

3. **Create the Claude Code symlink** (if Claude Code integration is enabled):
   ```
   <project>/.claude/skills/<prefix>-<skill-name>  →  ../../.agents/skills/<prefix>-<skill-name>
   ```

Record each installed skill name for the config file.

#### 3.2 — Write root templates

Copy and transform these files from the kit root to the project root:

| Source (relative to kit root) | Destination (project root) |
|---|---|
| `../../templates/AGENTS.md` | `AGENTS.md` |
| `../../templates/CLAUDE.md` | `CLAUDE.md` |
| `../../templates/CONTRIBUTING.md` | `CONTRIBUTING.md` |
| `../../templates/CONTEXT.md` | `CONTEXT.md` |

In `AGENTS.md`, replace the `<!-- SKILLS_TABLE -->` marker with a generated
markdown table:

```markdown
| Skill | Owns |
|---|---|
| `<prefix>-review` | Five-axis code review gate |
| `<prefix>-verify` | Adversarial claim verification against specifications |
| ...
```

Use the `Summary` field from the skill registry for each installed skill.

#### 3.3 — Claude Code integration (if opted in)

Copy and transform from the kit root:

```
../../claude/agents/*.md       →  <project>/.claude/agents/*.md
../../claude/commands/*.md     →  <project>/.claude/commands/*.md
../../claude/workflows/*.js    →  <project>/.claude/workflows/*.js
```

Create directories as needed.

#### 3.4 — OpenSpec schemas (if opted in)

Copy and transform from the kit root:

```
../../openspec/                →  <project>/openspec/
```

Preserve the full directory tree (config.yaml, schemas/, templates/).

#### 3.5 — Git hooks (if opted in)

Copy and transform from the kit root:

```
../../hooks/pre-commit   →  <project>/.husky/pre-commit     (chmod 755)
../../hooks/commit-msg   →  <project>/.husky/commit-msg     (chmod 755)
../../hooks/pre-push     →  <project>/.husky/pre-push       (chmod 755)
```

#### 3.6 — CI templates (if opted in)

Copy and transform from the kit root:

```
../../ci/*.yml           →  <project>/.github/workflows/*.yml
```

#### 3.7 — Write config file

Write `agent-kit.yaml` at the project root with the final configuration:

```yaml
projectName: my-awesome-app
prefix: maa
stack:
  - React
  - Next.js
  - Playwright
packageManager: pnpm
monorepo: true
includeOpenSpec: true
includeClaude: true
includeHooks: true
includeCI: true
optionalSkills: []
installedSkills:
  - review
  - verify
  - tdd
  - bugfix
  - ...
```

### Phase 4 — Report

After scaffolding, report:

1. **Summary** of what was installed — skill count, modules enabled, files created.
2. **Next steps** the user should take:
   - Read `AGENTS.md` and fill in project-specific boundaries.
   - Review `CONTRIBUTING.md` and replace placeholder commands with real ones.
   - If OpenSpec was installed, customise `openspec/config.yaml` with project context.
   - Run the project's validation command to confirm nothing is broken.
   - Commit the new files.

## Rules

1. **Never overwrite existing files without asking.** If `AGENTS.md`, `CONTRIBUTING.md`,
   or any target file already exists, ask the user whether to overwrite, merge, or skip.

2. **Transforms are exact string replacement.** Do not invent new placeholders.
   Apply only the replacements documented in `references/transform-rules.md`.

3. **Symlinks use relative paths.** The `.claude/skills/<name>` symlink must use a
   relative path to `.agents/skills/<name>`, not an absolute path.

4. **Skill filtering is deterministic.** Follow the stack-matching rules exactly
   as documented in `references/skill-registry.md`. Do not guess which skills
   are relevant — use the lookup table.

5. **Do not modify kit source files.** The kit directory is read-only reference
   material. Copy, transform, and write — never edit in place.

6. **Config file is YAML.** Use `agent-kit.yaml`, not JSON or TOML. Field names
   use camelCase to match the existing schema.

## Red flags

- **`agent-kit.yaml` already exists** — the project is already initialized.
  Stop and tell the user. Suggest `agent-kit add <skill>` for adding skills.
- **`.agents/skills/` exists with content** — partial installation. Ask whether to
  continue and overwrite, or abort.
- **No package.json / go.mod / equivalent** — cannot detect stack. Ask the user
  to provide stack information manually.
- **Git is not initialized** — hooks cannot be installed. Warn the user and skip
  the hooks module.
- **Existing CLAUDE.md has custom content** — do not overwrite. Ask the user
  whether to append `@AGENTS.md` to the existing file or skip.

## Verification

After scaffolding:

- [ ] `agent-kit.yaml` exists at the project root and is valid YAML.
- [ ] Every installed skill has a directory under `.agents/skills/<prefix>-<name>/`
      with a `SKILL.md` file.
- [ ] If Claude Code integration is enabled, every skill has a symlink under
      `.claude/skills/<prefix>-<name>` pointing to the correct `.agents/` directory.
- [ ] `AGENTS.md` exists and contains the skills table (not the `<!-- SKILLS_TABLE -->`
      placeholder).
- [ ] All placeholder strings (`your-project-name`, `your-package-manager`,
      `YOUR_PROJECT`) have been replaced in every generated file.
- [ ] Cross-skill references use the prefixed names (e.g. `` `<prefix>-review` ``
      not `` `review` ``).
- [ ] No absolute paths in symlinks.
- [ ] Hook files (if installed) are executable (chmod 755).

Run a quick spot-check: open one installed skill's SKILL.md and verify that
`your-project-name` does not appear anywhere in the file.

## Report

When done, output:

```
Agent Kit initialized.

  Project:          <name>
  Prefix:           <prefix>
  Skills installed: <count> (<core count> core + <stack count> stack-matched + <optional count> optional)
  Modules:          skills, templates[, openspec][, claude][, hooks][, ci]
  Config:           agent-kit.yaml

Next steps:
  1. Read AGENTS.md and fill in project-specific boundaries
  2. Review CONTRIBUTING.md and adjust commands
  3. Run your validation command
  4. Commit the new files
```
