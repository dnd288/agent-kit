# Chore PR — `pr-chore`

**When:** tooling, dependencies, CI, docs, or any work with no user-facing behaviour change.

**Rule:** no spec, no deep modes, no issue required. One workspace only — if it touches two workspaces it is probably not a chore.

## The rules

| | |
|---|---|
| Branch | `docs/<short-description>` — no issue number, work without one is allowed here |
| Title | `chore(scope):` — no `—` suffix; all context goes in the body. Drop `docs(scope):` — `chore(scope):` covers everything non-feature. |
| Issue link | None |
| Proof | `your validation command` |
| Labels | `layer:docs` if pure docs, otherwise the workspace `layer:*` it touches |

## Template body fields

None required.

## Red flags

- **Mixed scope** — a PR that addresses more than one issue, or more than one unrelated change. One issue = one PR.
- A `deps` scope change with no `pnpm audit` output in the body
- A chore touching two workspaces — it may be two changes
- No validation run
