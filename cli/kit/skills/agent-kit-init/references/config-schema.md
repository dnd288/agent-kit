# Configuration schema

This reference documents every field in `agent-kit.yaml`, its type, default value,
and how to auto-detect it from the codebase.

## Fields

### `projectName`
- **Type:** string
- **Required:** yes
- **Default:** directory name of the project root
- **Detection:**
  - `package.json` → `name` field (strip `@scope/` prefix if present)
  - `go.mod` → last segment of the module path
  - `Cargo.toml` → `[package].name`
  - `pyproject.toml` → `[project].name` or `[tool.poetry].name`
  - `mix.exs` → `:app` value
  - Fallback → basename of the current working directory

### `prefix`
- **Type:** string (lowercase, hyphens only, max 8 characters)
- **Required:** yes
- **Default:** derived from `projectName`
- **Derivation rules:**
  - Single word → use it, truncated to 8 chars. `dashboard` → `dashboar`
  - Two+ words → use first letter of each word. `acme-web-platform` → `awp`
  - Strip non-alphanumeric characters, lowercase, replace underscores with hyphens
  - If result is empty → `proj`
- **Purpose:** prefixes every installed skill name to avoid collisions with other
  agent-kit installations. `review` → `<prefix>-review`.

### `stack`
- **Type:** string array
- **Required:** no (empty means no stack-specific skills)
- **Default:** auto-detected from dependencies
- **Detection (Node.js — read `package.json` `dependencies` + `devDependencies`):**

  | Dependency key | Stack value |
  |---|---|
  | `react` or `react-dom` | `React` |
  | `next` | `Next.js` |
  | `vue` | `Vue` |
  | `svelte` or `@sveltejs/kit` | `Svelte` |
  | `@angular/core` | `Angular` |
  | `fastify` | `Fastify` |
  | `express` | `Express` |
  | `@nestjs/core` | `NestJS` |
  | `hono` | `Hono` |
  | `prisma` or `@prisma/client` | `Prisma` |
  | `@playwright/test` or `playwright` | `Playwright` |
  | `cypress` | `Cypress` |

- **Detection (Go, Rust, Python, Elixir):** report as "unable to auto-detect
  stack" and ask the user to provide it manually. The stack values above are
  the canonical set; non-JS projects can use these if they apply or leave stack
  empty.

### `packageManager`
- **Type:** string enum: `pnpm`, `npm`, `yarn`, `bun`
- **Required:** no (default: `npm`)
- **Default:** auto-detected from lockfiles
- **Detection (check in this order — first match wins):**

  | File exists | Package manager |
  |---|---|
  | `pnpm-lock.yaml` or `pnpm-workspace.yaml` | `pnpm` |
  | `bun.lockb` or `bun.lock` | `bun` |
  | `yarn.lock` | `yarn` |
  | `package-lock.json` | `npm` |
  | None of the above | `npm` |

### `monorepo`
- **Type:** boolean
- **Default:** `false`
- **Detection:**
  - `pnpm-workspace.yaml` exists → `true`
  - `lerna.json` exists → `true`
  - `nx.json` exists → `true`
  - `package.json` has `workspaces` field → `true`
  - `Cargo.toml` has `[workspace]` section → `true`

### `ticketTracker`
- **Type:** string enum: `none`, `github`, `jira`, `linear`
- **Default:** `none`
- **Detection:**
  - `git remote get-url origin` contains `github.com` and the `gh` CLI is on
    `PATH` → `github`
  - Otherwise → `none` (ask the user; they know if the team uses Jira or
    Linear)
- **Purpose:** the `flow` skill reads this to decide whether feature work
  starts with a ticket and where progress reports go.

### `featureFlow`
- **Type:** string array (ordered stage names)
- **Default:** `["capture", "specify", "implement", "verify", "deliver", "record"]`
- **Valid values:** any ordered subset or custom naming of the six canonical
  stages — `capture` (ticket), `specify` (proposal), `implement`, `verify`,
  `deliver` (PR), `record` (tracker report)
- **Detection:** not detectable from code — ask the user what the project's
  end-to-end flow for one feature is, then record their answer here
- **Purpose:** the `flow` skill runs exactly these stages in this order.
  Projects without a tracker commonly drop `capture`/`record` or fold them
  into `specify`.

### `includeOpenSpec`
- **Type:** boolean
- **Default:** `true`
- **Purpose:** install OpenSpec specification schemas under `openspec/`.

### `includeClaude`
- **Type:** boolean
- **Default:** `true`
- **Purpose:** install Claude Code integration — agents, commands, workflows,
  and `.claude/skills/` symlinks.

### `includeHooks`
- **Type:** boolean
- **Default:** `true`
- **Prerequisite:** git must be initialized (`.git/` exists).
- **Purpose:** install git hooks under `.husky/`.

### `includeCI`
- **Type:** boolean
- **Default:** `true`
- **Purpose:** install GitHub Actions CI templates under `.github/workflows/`.

### `optionalSkills`
- **Type:** string array
- **Default:** `[]`
- **Valid values:** `optimization-loop`, `problem-solving`, `scenario-explorer`
- **Purpose:** tracks which optional skills the user chose to install.

### `installedSkills`
- **Type:** string array
- **Default:** populated during scaffolding
- **Purpose:** records every skill that was installed (by unprefixed name).
  Used by `agent-kit add` to avoid reinstalling existing skills.
