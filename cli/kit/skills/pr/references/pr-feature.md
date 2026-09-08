# Feature PR — `pr-feature`

**When:** implementing a new capability, or substantially changing an existing one.

**Rule:** a spec PR must be merged first unless the change is small enough to review the same day without one. `spec-workflow` owns the spec-first path.

## The rules

| | |
|---|---|
| Branch | `feature/<issue>-short-description` |
| Title | `feat(scope):` — no `—` suffix; all context goes in the body |
| Issue link | `Closes #N` — this finishes the issue |
| Proof | `your validation command`; `your deep validation` if UI changed |
| Deep modes | Run `your deep validation` if your UI package, component library, or app path changed. For any other workspace, `your validation command` is sufficient. |
| Labels | Every `layer:*` the change touches, plus the epic and priority from the issue |

## Template body fields

```
Closes #<N>
```

## Red flags

- **Mixed scope** — a PR that addresses more than one issue, or more than one unrelated change. One issue = one PR.
- A feature PR whose change folder was created the same day — spec was back-filled
- your UI package or component library changed and no `your CI labels for extended checks` label applied
- Touching two layers — it is two PRs
