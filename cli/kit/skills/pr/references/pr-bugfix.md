# Bugfix PR — `pr-bugfix`

**When:** fixing a defect. Two paths depending on how obvious the cause is.

**Obvious cause** (can explain the symptom in one sentence): goes straight to a PR.  
**Unclear cause** (diagnosis needed): open a `bugfix` change folder first — `bugfix` owns this path.

## The rules

| | |
|---|---|
| Branch | `fix/<issue>-short-description` |
| Title | `fix(scope):` — no `—` suffix; all context goes in the body |
| Issue link | `Closes #N` |
| Proof | `your validation command`; reproduction recorded in the commit body |
| Labels | Every `layer:*` the fix touches, plus the epic and priority from the issue |

## Template body fields

```
Closes #<N>
```

## Red flags

- A fix with no RED output (test or observation) proving the defect existed
- A ticked checklist item for a check that was never run
- A guard edited to make a failure go away rather than fixing the code
