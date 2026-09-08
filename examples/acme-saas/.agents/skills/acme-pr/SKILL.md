---
name: pr
description: "Opening an project pull request — choose the type that matches your work: pr-spec, pr-feature, pr-bugfix, pr-chore, pr-infra. Load when asked to open, raise or submit a pull request, write a PR title or body, name a branch, or link a PR to its issue and the project board. Routes to review for the five review axes and verify for the adversarial pass."
---

# Opening a pull request

Every PR goes through the same three steps in the same order: decide the type, apply the rules, do the post-open steps.

## Step 1 — Which type?

| | Spec | Feature | Bugfix | Chore | Infra |
|---|---|---|---|---|---|
| Carries | `your change tracking folder/<id>/` only | implementation | implementation | implementation | implementation |
| Title | `docs(spec): <id>` | `feat(scope):` | `fix(scope):` | `chore(scope):` | `chore(infra):` |
| Issue link | `Part of #N` | `Closes #N` | `Closes #N` | none | none |
| Branch prefix | `feature/` | `feature/` | `fix/` | `docs/` | `docs/` |
| Proof | `openspec validate --strict` | `your validation command` | `your validation command` | `your validation command` | `your validation command` |
| Needs issue | yes | yes | no if obvious, yes if not | no | no |
| Spec first | required | required unless obvious | N/A | N/A | N/A |

**The separator is whether the work is new or existing.** A new feature or enhancement is a `feature`. A defect is a `acme-bugfix`. Everything else — tooling, deps, CI, docs — is a `chore` or `infra`. `infra` is for AWS and pipeline work owned by your infrastructure package.

A bugfix with an unclear cause gets a `acme-bugfix` change folder first (`acme-bugfix` owns this). A feature touching two layers is two PRs.

## Step 2 — Apply the rules for that type

Read the reference file for your PR type:

| Type | Reference |
|---|---|
| Spec | `references/pr-spec.md` |
| Feature | `references/pr-feature.md` |
| Bugfix | `references/pr-bugfix.md` |
| Chore | `references/pr-chore.md` |
| Infra | `references/pr-infra.md` |

## Step 3 — Shared mechanics

These apply to every PR type.

### The branch

Branches are `feature/`, `fix/`, or `docs/` from `main`. The prefix carries the issue number when one exists: `feature/184-transactional-email`. A branch without a number is only valid for chores and infra — work with no issue.

Never open a PR from `main`. Never `--force` a push to a branch that has been reviewed.

### The commit message

Conventional Commits. **Scope is a workspace from the table in `conventions.md`**, not the feature name.

```
fix(ai): handle credit_balance_exhausted

The SDK error carries code credit_balance_exhausted (with type
insufficient_quota); the classifier listed only the type, so the live
no-credit refusal fell through to the transient path and redelivered.
Cover both codes, keep the message fingerprint, and pin the real shape
in the worker test.
```

Problem first, then the fix, then what proves it. The body is the most durable documentation in the repository.

### The body

Use `.github/pull_request_template.md` — fill the checklist honestly. A ticked box that is false is worse than an unticked box with a reason. Put the essay in `## What and why`.

`## What and why` says what the change **deleted or unified**, or that it removed nothing and why — the [golden rules](../../../AGENTS.md#golden-rules) in the one place a reviewer reads before the diff. Applies to all five types: a `chore` that removes a duplicate is as much a result as a `feature` that adds a screen.

### The labels

`ci.yml`'s `on_demand` job turns on modes based on labels: `your CI labels for extended checks`, `your CI labels for extended checks`, `your CI labels for extended checks`, `performance-checks`. `on: pull_request` includes `labeled`, so adding one re-runs CI without a new push.

`acme-test` owns the runbook. CI is currently disabled (`CI` workflow is `disabled_manually`) — a label is a note to a human, not an automatic gate.

### After it is open

1. **Comment the PR URL onto the issue** — `gh issue comment <N> --body "PR: <url>"`
2. **Move the issue to `status:in-review`**, remove `status:in-progress`
3. **Add the PR to the board** — `gh project item-add 1 --owner MBL-Innovations-LTD --url <pr-url>`
4. **Request a reviewer by hand** — `.github/CODEOWNERS` is inert
5. Run `acme-review` and `acme-verify`

### The Development panel

your git host auto-links a branch to its issue only when the PR body contains `Closes`, `Fixes` or `Resolves`. The `feature/` and `fix/` branch prefixes are configured as autolink patterns in the repo settings — matching branches auto-appear in the Development panel without any keyword.

For spec PRs, `Part of #N` does not populate the panel. Add `Branch: feature/<issue>-short-description` as a workaround — the branch name is readable from the issue sidebar.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "The spec and the code in one PR is fewer steps" | It is, and the spec then gets read as an explanation of the diff rather than a contract |
| "CI is off, so the label does not matter" | The label is how a reviewer knows which modes you claim to have run |
| "`## Summary` / `## Test plan` says the same things" | It drops the migration, credential and security-surface items, which nobody remembers unaided |
| "I will link the issue when I merge" | The issue is how the reviewer finds context. Linking at merge is linking after it mattered |
| "The scope is obvious from the files" | Not in `git log --oneline`, which is the only place a scope is ever read |
| "It is one commit, the body can be the title" | Then the reasoning exists nowhere |

## Red flags (all types)

- **Mixed scope** — a PR that addresses more than one issue, or more than one unrelated change, even within the same epic. One issue = one PR. If the title needs a list of things, it is several PRs.
- A PR with no issue link and no board item — nothing on the board moves when it merges
- Two your ORM migration directories in one diff — one per PR, and the guard cannot untangle a conflict
- A ticked checklist item for a check that was never run
- A force-push after a review has been left
- A green PR on a UI change with no browser mode run anywhere

## Verification (all types)

- [ ] Type is decided and title matches the shape above
- [ ] Branch off `main` with correct prefix; issue number present when one exists
- [ ] Commit scope is a workspace from `conventions.md`
- [ ] Commit body says why, not just what
- [ ] Body follows `.github/pull_request_template.md`, checklist filled honestly
- [ ] Correct issue link keyword (`Closes` / `Part of` / none)
- [ ] `Branch:` line on spec PRs
- [ ] `your validation command` green; deep modes run or explicitly named as not run
- [ ] Right CI label applied if UI changed
- [ ] Issue commented with PR URL and moved to `status:in-review`
- [ ] PR added to project 1, reviewer requested by hand
