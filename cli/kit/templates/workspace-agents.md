# Workspace agent instructions

This file is the workspace-level instruction template. Copy it to `AGENTS.md` inside an application, package, service, tool, or documentation workspace, then replace placeholders with rules that apply only here.

## What this workspace is

`<workspace-name>` owns `<responsibility>`.

It contains:

- `<directory or module>` — `<purpose>`
- `<directory or module>` — `<purpose>`
- `<directory or module>` — `<purpose>`

It does not own:

- `<responsibility owned elsewhere>` — see `<path>`
- `<responsibility owned elsewhere>` — see `<path>`

## Boundaries specific to this workspace

- **Inputs:** `<where data/config/events enter>`
- **Outputs:** `<what this workspace exports or serves>`
- **Allowed dependencies:** `<packages, layers, services>`
- **Forbidden dependencies:** `<packages, layers, services>`
- **State ownership:** `<database/cache/store/local state rules>`
- **Generated files:** `<what is generated, command that regenerates it, whether hand edits are forbidden>`
- **External systems:** `<APIs, queues, webhooks, files, browsers, CLIs>`

Keep rules specific. If a rule applies everywhere, move it to the root `AGENTS.md`.

## Key files

- `<path>` — `<why this file matters>`
- `<path>` — `<why this file matters>`
- `<path>` — `<why this file matters>`

Read these before large changes:

- `<path to README or design doc>`
- `<path to ADR or spec>`

## Testing

Use the smallest mode that can observe the claim, then run this workspace's ready checks.

Common commands:

```sh
<workspace install or build command>
<workspace lint command>
<workspace typecheck command>
<workspace unit test command>
<workspace integration or browser test command>
```

Run `<required validation command>` before marking workspace changes ready.

Additional guidance:

- Add or update tests beside changed behaviour.
- Prefer characterization tests before refactors.
- For bug fixes, add the test that fails before the fix.
- For generated outputs, test the source and regenerate the output rather than hand-editing it.
- If a required test cannot run locally, report the command, the blocker, and where it must run.
