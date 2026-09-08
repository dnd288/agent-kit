# Skill format

A skill is a portable instruction package for a recurring task. It tells an agent when the task applies, which procedure to follow, which rules matter, what red flags to stop on, and how to verify the result.

Skills keep operational know-how out of chat history. A project should be usable by different coding agents because the durable procedure lives in the repository.

## Agent Skills format

Store skills under `.agents/skills/<skill-name>/SKILL.md`.

Recommended layout:

```text
.agents/
  skills/
    <skill-name>/
      SKILL.md
      references/
        <supporting-doc>.md
```

`SKILL.md` contains YAML frontmatter followed by Markdown instructions.

## YAML frontmatter

Required fields:

```yaml
---
name: project-task-name
description: Use this skill when ...
---
```

### `name`

Use a stable, lowercase, hyphenated name. Prefix with the project or domain when skills may be installed alongside other projects.

Examples:

```yaml
name: app-bugfix
name: app-pr
name: app-api-contract
name: app-ui-development
```

### `description`

The description is trigger text. Write it so an agent can decide when to load the skill.

Good description:

```yaml
description: Use this skill when fixing a defect, investigating a regression, or deciding whether an observed failure needs a bugfix specification.
```

Weak description:

```yaml
description: Bugfix stuff.
```

Include concrete triggers:

- File areas: `apps/api`, `packages/ui`, `infra/`.
- Task types: bugfix, PR, migration, e2e, design implementation.
- Risk surfaces: auth, sessions, uploads, secrets, external input.
- Deliverables: specification, ADR, component, route, workflow.

## Body structure

Use headings that match how the agent should work.

```md
# Skill name

## When to use this skill

State the triggers and non-triggers.

## Read first

List owning documents and key paths.

## Procedure

Give ordered steps.

## Rules

List rules that are specific to this procedure, or route to the owning document.

## Red flags

Name conditions that require stopping, asking, escalating, or changing mode.

## Verification

List checks and evidence required before declaring the task done.

## Report

State what the final response must include.
```

Keep the procedure actionable. Prefer "Run the affected unit test, then the workspace validation command" over "ensure quality".

## Routing principle

**Skills route; they do not restate.**

A skill should point to the document that owns a rule instead of copying the rule into the skill. Duplication creates drift: the doc changes, the skill keeps the old rule, and agents follow the wrong one.

Use the skill for:

- Task order.
- Which documents to read.
- Which checks prove the task.
- Which red flags require escalation.
- Tool-specific gotchas.

Use docs, specs, and ADRs for:

- Product requirements.
- Architecture decisions.
- Security policy.
- Testing strategy.
- Design-system rules.
- Deployment policy.

Good skill text:

```md
Read `docs/engineering/security.md` before changing auth, sessions, uploads, webhooks, secrets, or untrusted input. Follow its surface-specific controls. If the change violates one, write or update an ADR that bounds the exception.
```

Weak skill text:

```md
Authentication uses server-side sessions with these fifteen detailed rules...
```

The weak version will drift when the security document changes.

## Symlink convention

Keep one source of truth for skills.

Recommended convention:

```text
.agents/skills/<skill-name>/SKILL.md   # canonical copy
.claude/skills/<skill-name>            # symlink to canonical skill directory
```

This lets generic agents read `.agents/skills/` while Claude Code can discover the same skill through `.claude/skills/`. Do not maintain two copies.

Example:

```sh
mkdir -p .agents/skills/app-bugfix
ln -s ../../.agents/skills/app-bugfix .claude/skills/app-bugfix
```

Adjust the relative path for your repository layout and operating system.

## References subdirectory

Use `references/` for supporting material that belongs with the procedure but is too detailed for the main skill.

Examples:

```text
.agents/skills/app-e2e/references/locator-conventions.md
.agents/skills/app-pr/references/pr-labels.md
.agents/skills/app-release/references/release-checklist.md
```

Use references for:

- Command matrices.
- Long examples.
- Troubleshooting tables.
- Tool output interpretation.
- Templates that are specific to the skill.

Do not use references to hide core rules. The main `SKILL.md` should still tell the agent what to read and when.

## Citing documents instead of duplicating rules

When a skill depends on a rule, cite the owning document and section.

Examples:

```md
- State placement follows `docs/engineering/state.md#decision-ladder`.
- API response schemas follow `docs/engineering/api-contract.md#responses`.
- ADRs follow `docs/adr/README.md#format`.
- Security-sensitive changes follow `docs/engineering/security.md#surfaces`.
```

If the cited document changes, update the skill only when the procedure changes. If the skill had copied the old rule, both files would need edits.

## Red flags

A red flag tells the agent not to continue blindly.

Good red flags:

- The change touches credentials, tokens, private keys, or production data.
- A migration is not backward compatible.
- The observed bug has no reproduction or log evidence.
- A generated file was edited by hand.
- The task spans more workspaces than the requested change type allows.
- A guard appears wrong and would need to be weakened.
- A test passes only when run after another test.

For each red flag, say what to do: stop and ask, write a specification, run a security skill, create an ADR, or narrow the change.

## Verification

Every skill should define done evidence.

Examples:

```md
## Verification

- Run `<unit test command>` for changed services.
- Run `<component test command>` for changed UI states.
- Run `<validate command>` before marking ready.
- If a check cannot run, report the command and blocker.
```

Do not let skills end with "ensure it works". Name the proof.

## Skill review checklist

Before adding or changing a skill:

- [ ] The `name` is stable and specific.
- [ ] The `description` includes clear trigger text.
- [ ] The procedure is ordered and executable.
- [ ] Durable rules are cited, not duplicated.
- [ ] Red flags say what action to take.
- [ ] Verification names commands or evidence.
- [ ] References are linked from the main skill.
- [ ] The `.claude/skills/` entry points to the canonical `.agents/skills/` copy when using the symlink convention.
