# Skill registry

This reference is the single source of truth for which skills the kit ships, how
they are categorised, which stacks they apply to, and whether they are optional.

## Categories and display order

1. Core Methodology
2. Development
3. Testing
4. Operations

## Full registry

### Core Methodology

| Skill | Summary | Optional | Stack requirement |
|---|---|---|---|
| `review` | Five-axis code review gate | no | — (always) |
| `verify` | Adversarial claim verification against specifications | no | — (always) |
| `tdd` | RED→GREEN discipline with mode-to-claim mapping | no | — (always) |
| `bugfix` | inspect→debug→impact→fix→verify loop | no | — (always) |
| `refactor` | MEASURE→PIN→MOVE→PROVE→RECORD loop | no | — (always) |
| `simplicity` | Simplest-sufficient-shape decision — the rung ladder | no | — (always) |
| `spec-workflow` | Typed changes with end-to-end-first ordering | no | — (always) |
| `pr` | Five PR types with checklists and conventions | no | — (always) |
| `flow` | End-to-end feature pipeline: ticket → specify → implement → verify → deliver → record | no | — (always) |
| `optimization-loop` | BASELINE→CHANGE→MEASURE→KEEP-OR-REVERT performance loop | **yes** | — |
| `problem-solving` | Impasse techniques and the dev-note handed to the user | **yes** | — |

### Development

| Skill | Summary | Optional | Stack requirement |
|---|---|---|---|
| `component-development` | Primitive component development patterns | no | React, Vue, Svelte, Angular |
| `ui-development` | Feature-folder composite UI patterns | no | React, Vue, Svelte, Angular |
| `app-development` | Application routing, data flow, access control | no | Next.js, Nuxt, SvelteKit, Remix |
| `api-contract` | Schema-first API design and service purity | no | Fastify, Express, NestJS, Hono |
| `state-management` | State placement decision ladder | no | React, Vue, Svelte, Angular |
| `design` | Design system methodology | no | — (always) |

### Testing

| Skill | Summary | Optional | Stack requirement |
|---|---|---|---|
| `e2e` | End-to-end test authoring and healing | no | Playwright, Cypress |
| `test` | Test mode taxonomy and operational runbook | no | — (always) |
| `scenario-explorer` | Enumerate the case space before writing tests | **yes** | — |

### Operations

| Skill | Summary | Optional | Stack requirement |
|---|---|---|---|
| `security` | Surface-aware security review | no | — (always) |
| `local-dev` | Local development environment setup | no | — (always) |

## Installation rules

### Core skills (always install)

Skills with no stack requirement and `Optional: no` are always installed regardless
of the detected stack.

### Stack-matched skills

A skill with a stack requirement is installed only when the project's detected
`stack` array includes at least one of its required values (case-insensitive match).

Example: `component-development` requires `React`, `Vue`, `Svelte`, or `Angular`.
If the project's stack is `["React", "Next.js"]`, it matches on `React` and is
installed.

### Optional skills

Skills with `Optional: yes` are installed only when the user explicitly opts in
during the confirm phase. They appear in a separate section of the confirmation
prompt.

### Unknown skills

If a skill directory exists under `kit/skills/` but is not listed in this registry,
it is always installed (fail-open). This prevents new skills from being silently
skipped.

## Skill directory structure

Each skill under `kit/skills/<name>/` has:

```
<name>/
├── SKILL.md              # Required — frontmatter + instructions
└── references/           # Optional — supporting documents
    └── *.md
```

The `SKILL.md` frontmatter must have:
- `name:` — must match the directory name exactly
- `description:` — trigger text for agent discovery
