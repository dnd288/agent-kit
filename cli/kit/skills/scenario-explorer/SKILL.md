---
name: scenario-explorer
description: "Enumerating the input and state space of a behaviour before writing code or tests — turning a vague requirement into the explicit list of cases that must then be handled and covered. Load this when starting a feature or change and deciding what it must handle; when a requirement is stated only as its happy path; when a bug reveals a case nobody considered; when asked for edge cases, boundaries, or 'what could go wrong'; before writing tests, so the test list is complete rather than whatever came to mind. Owns the decomposition lenses (happy path, boundaries, empty, error, concurrency, adversarial, state transitions), the disposition of each case (tested, handled, or bounded out of scope), and the rule that the case list is anchored in BOTH the spec and the pull request. Routes to tdd to turn each scenario into a RED test, to e2e for journey-level cases, to security for the adversarial lens, to spec-workflow for where the case list lives, and to review's correctness axis, which reads that list as the claim set."
---

# Exploring the scenario space

**You cannot test, handle, or review what you have not enumerated.** The list of cases *is* the
claim set — a case nobody wrote down is a case nothing covers, and it ships as a defect wearing the
disguise of "we didn't think of that". This skill is the step before `tdd`: it discovers *what*
claims exist, so the test list is the space, not whatever came to mind first.

The requirement almost always arrives as its happy path. The work here is turning "let a user rename
a project" into the dozen cases hiding inside it.

## The decomposition lenses

Walk every lens; most requirements have cases in more than one. Name the cases, do not just nod at
the category.

| Lens | What it asks | Typical cases |
|---|---|---|
| **Happy path** | The stated behaviour, working | the obvious one — usually the only one given |
| **Boundaries** | The edges of each input range | 0, 1, many, max, off-by-one, first/last, exactly-at-limit |
| **Empty / missing** | Absent or null input | empty string, no rows, unset field, missing param |
| **Error / failure** | A step that can fail | downstream error, timeout, malformed input, partial write |
| **Concurrency / ordering** | More than one at once, or out of order | two writers, retry, double-submit, idempotency, stale read |
| **Adversarial** | A hostile or unauthorised caller | injection, oversized input, another user's resource — route `security` |
| **State transitions** | Each state crossed with each event | edit while archived, cancel after commit, act on a deleted thing |

**Partition; do not enumerate every value.** Three positive integers test the same path — one
representative plus the boundaries is the coverage, not all three. The skill is choosing the cases
that exercise *different* behaviour, and stopping when the next case would exercise the same one.

## Every scenario gets a disposition

A case list with no decision attached to each case is a worry list. Each scenario is one of:

- **Tested** — a test observes it. `tdd` turns it into a RED test first; `e2e` owns the journey-level
  ones; `test` owns which mode can see it.
- **Handled in code** — the behaviour exists and a test above proves it.
- **Out of scope, on purpose** — written down and bounded, never silently dropped. "Concurrent
  rename by two admins: not handled this change, single-admin assumption, revisit if teams ship"
  is a decision (golden rule 5); an unmentioned gap is a bug waiting.

## Anchor the list in two places

The case list is not a scratch note — it is a durable artefact, and it lives in **both**:

1. **The spec, in the openspec change-folder.** The scenarios are the claim set the change is built
   against and reviewed against — `review`'s correctness axis reads exactly this list, and a
   scenario with no covering test is an unverified claim. `spec-workflow` owns the change folder and
   where the scenarios sit in it.
2. **The pull request checklist.** A short "scenarios covered / deliberately out of scope" block in
   the PR body so a reviewer sees the whole space and what was decided about each edge — not just the
   diff. `pr` owns the PR shape; this is what fills its scenario section.

The two are the same list at different lifetimes: the spec is where it is agreed, the PR is where it
is shown to have been honoured.

## Red flags

- A scenario on the list with neither a test nor an "out of scope" note — the one case guaranteed to
  break.
- A requirement whose only case is the happy path, accepted as complete.
- A bug in production that maps to a lens you never walked — add the lens to your defaults, not just
  the one case.
- "We'll handle that edge later" with no record of the edge — later has no list to read.

## Before you finish

- [ ] Every lens was walked, not just the ones that came to mind
- [ ] Cases are partitioned, not every value listed
- [ ] Each scenario has a disposition: tested, handled, or bounded out of scope
- [ ] The list is in the change-folder spec and in the PR checklist
- [ ] Adversarial cases were routed through `security`, not guessed at

## Read next

- [`tdd`](../tdd/SKILL.md) — turning each scenario into a RED test before the code.
- [`e2e`](../e2e/SKILL.md) — the journey-level scenarios across the stack.
- [`spec-workflow`](../spec-workflow/SKILL.md) — where the case list lives in the change folder.
- [`review`](../review/SKILL.md) — the correctness axis that reads the scenario list as the claim set.
- [`security`](../security/SKILL.md) — the adversarial lens, done properly.
- [`simplicity`](../simplicity/SKILL.md) — the non-lazy boundary: handling every case the caller will
  actually pass.
