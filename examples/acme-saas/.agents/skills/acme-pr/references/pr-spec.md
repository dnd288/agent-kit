# Spec PR — `pr-spec`

**When:** new feature or enhancement, big enough that another team builds against the contract.

**Rule:** the spec ships before anyone writes code. A spec reviewed after the code is a spec written to match the code.

## The rules

| | |
|---|---|
| Branch | `feature/<issue>-short-description` |
| Title | `docs(spec): <id>` — no `—` suffix; all context goes in the body |
| Issue link | `Part of #N` — not `Closes`, because agreeing a contract is not finishing the work |
| Branch line | `Branch: feature/<issue>-short-description` on its own line |
| Proof | `pnpm exec openspec validate <id> --strict` |
| Labels | `layer:docs` |

## Template body fields

```
Part of #<N>
Branch: feature/<issue>-<desc>
```

## Red flags

- `Closes #N` on a spec PR — the issue is not done when its contract is agreed
- No `Branch:` line — the Development panel will be empty
- A spec PR whose change folder was created the same day — the spec was back-filled
