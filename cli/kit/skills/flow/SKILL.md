---
name: flow
description: "End-to-end feature delivery pipeline: ticket → specify → implement → verify → deliver → record, synced with the project's issue tracker. Load when asked to run the flow, take a feature from idea to PR, start or resume an end-to-end change, or report progress to the tracker."
---

# Feature flow

One pipeline from idea to merged PR, with the issue tracker as the single
record of what is happening and why. The stages and the tracker are project
configuration — read both from `agent-kit.yaml` before starting:

```yaml
ticketTracker: github        # none | github | jira | linear
featureFlow:                 # ordered stages; only these, only in this order
  - capture
  - specify
  - implement
  - verify
  - deliver
  - record
```

Default shape:

```
idea → ticket → spec proposal → implement → verify → PR → tracker report
```

## Stages

### capture — ticket

Create the tracking ticket and record the acceptance criteria on it.

| Tracker | How |
|---|---|
| `github` | `gh issue create --title "<type>(scope): short summary" --label "<type>"` with a body holding **Context** (what and why) and **Acceptance Criteria** (checklist) |
| `jira` / `linear` | Create in the tracker's app or CLI; record the key (`PROJ-123`) and paste the same Context + Acceptance Criteria sections |
| `none` | Skip this stage; carry context and acceptance criteria into the change folder in specify |

If the user provides an existing ticket (`#12`, `PROJ-123`), read it instead of
creating one: extract context and acceptance criteria, derive a kebab-case
change name if the ticket does not name one.

**Schema decision:** defect/fix work → `bugfix` shape; everything else →
default feature shape. The choice drives branch naming and the PR type later.

### specify — proposal

Follow the `spec-workflow` skill to create the change proposal and its
artifacts. The ticket's acceptance criteria are the source of truth: they
become the scenarios in the spec. Do not invent requirements the ticket does
not have — gaps go back to the ticket, not into the spec silently.

When the proposal is ready, post a summary to the ticket (capabilities, key
decisions, task count, status: ready for implementation) and **pause for the
user to review the artifacts** before implementing.

### implement

1. Create the branch: `feature/<change-name>` — or `fix/<change-name>` for the
   bugfix shape.
2. Follow the `tdd` skill (or `bugfix` for defects). The first task is always
   the failing test.
3. Work through the tasks in order; mark checkboxes as you complete them.
4. If the change touches a user-visible flow the project's e2e suite covers,
   follow the `e2e` skill: extend the journey test before implementing —
   same RED→GREEN rule.

### verify

Follow the `verify` skill. Fill the change's `verification.md` with actual
pasted output — RED output (before), GREEN output (after), what ran, and any
gaps. Then run the project's full validation command. Prose summaries instead
of output do not count.

### deliver

Follow the `pr` skill. Run full validation green first. The PR body states the
user-visible change and links the ticket — `Closes #<number>` on GitHub; an
explicit link in the body for other trackers.

### record — tracker report

Comment on the ticket with the completion report: PR link, branch, what was
built, verification summary (build/test/CI), files changed, status. After the
PR merges, archive the change per the project's spec workflow.

## Stage boundaries update the ticket

The ticket is updated at exactly these points: proposal ready (after specify),
PR opened (after deliver), completion report (record). Skipping a boundary
update breaks the only record other people read.

## Resuming a flow

Given a change name or ticket reference:

1. Inspect the change artifacts: none → resume at specify; tasks incomplete →
   implement; verification missing → verify; no PR → deliver; PR open → report
   status.
2. Find the linked ticket by searching the tracker for the change name.
3. Continue from the first incomplete stage in `featureFlow`.

## Guardrails

1. **Ticket first when a tracker is configured.** It is the tracking record;
   do not start implementing without one unless the user explicitly declines.
2. **Acceptance criteria from the ticket are the source of truth** for the
   spec's scenarios.
3. **Never skip verification.** `verification.md` holds real pasted output.
4. **PR references the ticket.** `Closes #<number>` on GitHub; link in body
   elsewhere.
5. **Scope grew beyond the ticket? Update the ticket first**, then implement.
6. **Stages run in `featureFlow` order.** Skip stages the project dropped;
   never invent extra stages.

## Red flags

- Implementing before any acceptance criteria or spec exists.
- "I'll update the ticket later" — boundary updates are the flow, not
  bookkeeping.
- `verification.md` with claims but no pasted output.
- PR with no ticket link.
- A new ticket created for work that already has one (search first).
