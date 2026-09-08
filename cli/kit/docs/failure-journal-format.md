# Failure journal format

How to record what went wrong, what fixed it, and what to watch for next time — without the record becoming noise.

## Why keep one

Failures repeat. The same misconfiguration, the same off-by-one, the same misread assumption comes back in a different file, a different sprint, a different person's hands. A failure journal turns a solved problem into a pattern the team can search.

The journal is not a post-mortem. Post-mortems are events; the journal is a growing reference. A post-mortem asks "what happened to the system"; a journal entry asks "what did I learn that someone else will hit".

---

## 1. Entry format

Each entry is a markdown file. One entry per failure. The filename is a date-slug: `2024-03-15-prisma-migration-ordering.md`.

```markdown
# Prisma migration ordering

**Date:** 2024-03-15
**Author:** (who wrote the entry)
**Severity:** broke CI for 2 hours

## Problem

The symptom, and the wrong assumption that hid the cause.

> Two migrations created on the same day ran in the wrong order because
> Prisma sorts by filename and both started with the same timestamp.
> We assumed timestamps were unique; they are unique per second, but
> two developers ran `prisma migrate dev` within the same second.

## Solution

What actually fixes it, and why the assumption was wrong.

> Prisma migration names are `<timestamp>_<slug>`. The timestamp is
> `YYYYMMDDHHmmss` — second precision. Two migrations in the same
> second get the same prefix and sort by slug, which is alphabetical
> and unrelated to intent.
>
> Fix: always check `prisma/migrations/` for an existing migration
> with today's timestamp before creating a new one. If one exists,
> wait one second or rename.

## Affects

The files, behaviours, and other callers this touches.

> - `prisma/migrations/` — ordering of all future migrations
> - CI pipeline — migration step assumes sorted order is correct order
> - Local dev — `prisma migrate dev` replays in the same order

## Implements

The concrete change made (PR link, commit, config change).

> PR #142 — added a pre-migration check to the dev script.
> Also added a CI step that verifies migration timestamp uniqueness.

## Result

The observed outcome after the fix, with evidence.

> CI green after the reorder. The pre-migration check caught a
> duplicate timestamp in PR #157 two weeks later.

## Tags

`prisma`, `migration`, `ci`, `race-condition`
```

---

## 2. Where entries live

```
docs/failures/
├── 2024-03-15-prisma-migration-ordering.md
├── 2024-03-22-session-cookie-path.md
├── 2024-04-01-embed-fs-dotfiles.md
└── ...
```

The directory is flat. Entries are sorted by date in the filename. Subdirectories add structure nobody navigates.

---

## 3. When to write an entry

Write one when:

- A bug took more than 30 minutes to diagnose.
- The root cause was a wrong assumption, not a typo.
- The fix required understanding something that is not in the documentation.
- Someone else on the team will hit the same thing.

Do not write one for:

- A typo or a missing import.
- A bug caught by an existing test (the test *is* the journal entry).
- A failure you cannot explain yet — wait until you can.

---

## 4. The dev-note shorthand

Not every failure earns a full journal entry. The `problem-solving` skill produces a **dev note** during debugging:

```
Problem:    the symptom, and the wrong assumption that hid the cause
Solution:   what actually fixes it, and why the assumption was wrong
Affects:    the files, behaviours, and other callers this touches
Implements: the concrete change made
Result:     the observed outcome after the fix, with evidence
```

A dev note is the minimum viable record. If the failure is worth keeping, promote the dev note to a journal entry by adding the date, author, severity, and tags. If it is not, the dev note lives in the commit message or PR description and does its job there.

---

## 5. Using the journal

### Search

```bash
grep -rl "prisma" docs/failures/
grep -rl "race-condition" docs/failures/
```

Tags at the bottom of each entry make `grep` the only index you need.

### Onboarding

New team members read `docs/failures/` in date order. Each entry is a lesson that took someone else hours to learn and takes five minutes to read.

### Review

When reviewing a change that touches a known failure area, search the journal. If there is a relevant entry, check whether the change respects the lesson or re-introduces the problem.

---

## 6. Anti-patterns

| Pattern | Why it fails | Fix |
|---|---|---|
| Writing the entry before the fix | The "solution" is a guess, not a fact | Write after the fix is verified |
| Entries without "Result" | No evidence the fix worked | Add the evidence, or mark the entry as unverified |
| Entries that blame people | The journal becomes a shame list | Blame the system, not the person |
| Entries for every bug | The journal becomes noise | Apply the "when to write" filter |
| Never reading the journal | The lessons are learned twice | Search before debugging; read on onboarding |

---

## Related

- [Orchestration protocol](orchestration-protocol.md) — the failure protocol agents follow during fan-out
- [Context safety](context-safety.md) — how dev notes stay small enough to be useful
- [philosophy.md](philosophy.md) — the compound engineering principles: leave it easier to change
