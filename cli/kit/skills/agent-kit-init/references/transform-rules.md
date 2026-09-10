# Transform rules

When copying files from the kit to the project, apply these exact string
replacements to every `.md`, `.yaml`, `.yml`, and `.js` file. Binary files and
other extensions are copied without transformation.

## Replacements

Apply in this order. Each replacement is global (replace all occurrences).

### 1. Project name

| Find | Replace with |
|---|---|
| `your-project-name` | `<projectName>` from config |

### 2. Uppercase project prefix

| Find | Replace with |
|---|---|
| `YOUR_PROJECT` | Uppercase version of `<prefix>`, with hyphens converted to underscores |

Example: prefix `acme-web` → `ACME_WEB`

### 3. Package manager run command

| Find | Replace with |
|---|---|
| `your-package-manager run` | `<pmRun>` |

Where `<pmRun>` is derived from the package manager:

| Package manager | `<pmRun>` value |
|---|---|
| `pnpm` | `pnpm` |
| `bun` | `bun` |
| `npm` | `npm run` |
| `yarn` | `yarn run` |

**Important:** this replacement must run before the generic package manager
replacement below, because `your-package-manager run` contains `your-package-manager`
as a substring.

### 4. Package manager name

| Find | Replace with |
|---|---|
| `your-package-manager` | `<packageManager>` from config |

### 5. Cross-skill references (backtick-quoted)

For every skill name in the registry, replace backtick-quoted references:

| Find | Replace with |
|---|---|
| `` `<skill-name>` `` | `` `<prefix>-<skill-name>` `` |

Example with prefix `acme`:
- `` `review` `` → `` `acme-review` ``
- `` `tdd` `` → `` `acme-tdd` ``
- `` `spec-workflow` `` → `` `acme-spec-workflow` ``

### 6. Cross-skill references (prose)

For every skill name in the registry, replace prose references:

| Find | Replace with |
|---|---|
| `the <skill-name> skill` | `the <prefix>-<skill-name> skill` |

Example with prefix `acme`:
- `the review skill` → `the acme-review skill`
- `the tdd skill` → `the acme-tdd skill`

### 7. Cross-skill relative paths

In the kit, sibling skills sit at `../<skill-name>/`. Installed directories
are prefixed, so rewrite the path segment to keep links resolving:

| Find | Replace with |
|---|---|
| `../<skill-name>/` | `../<prefix>-<skill-name>/` |

Examples with prefix `acme`:
- `../refactor/references/smells.md` → `../acme-refactor/references/smells.md`
- `../tdd/SKILL.md` → `../acme-tdd/SKILL.md`
- `../../state-management/references/data-flow.md` →
  `../../acme-state-management/references/data-flow.md`

The name must be bounded by slashes on both sides (`../` before, `/` after);
never rewrite inside longer names such as `../pre-analysis/`.

## Skill names to transform

Apply rules 5, 6, and 7 for each of these skill names:

```
review, verify, tdd, bugfix, refactor, simplicity, spec-workflow, pr,
optimization-loop, problem-solving, component-development, ui-development,
app-development, api-contract, state-management, design, e2e, test,
scenario-explorer, security, local-dev
```

## SKILLS_TABLE marker

In `AGENTS.md` only, replace the literal string `<!-- SKILLS_TABLE -->` with a
generated markdown table:

```markdown
| Skill | Owns |
|---|---|
| `<prefix>-<skill-name>` | <summary from skill registry> |
```

Include only installed skills, in registry order (Core Methodology → Development
→ Testing → Operations, then by order within each category).

## Files that are NOT transformed

- Binary files (images, compiled assets)
- Files outside the kit (the skill's own SKILL.md and references are not
  transformed — they are instructions for the agent, not project content)

## Idempotency note

Transforms are not idempotent. Running them twice would replace
`<prefix>-review` with `<prefix>-<prefix>-review`. The agent must track which
files have been transformed and never re-transform a file that has already been
processed.
