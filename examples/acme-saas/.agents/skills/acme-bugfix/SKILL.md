---
name: bugfix
description: "The defect loop — inspect, debug, IMPACT, fix, verify. Load this when fixing any bug, investigating a failure, or deciding whether a defect needs a change folder at all. Also load it when a test is failing and the cause is not yet known, when a fix is ready and needs proving, or when writing report.md, diagnosis.md, flows.md or verification.md in a `acme-bugfix` change. Owns the threshold that decides between a straight pull request and a bugfix change, the rule that the diagnosis carries its evidence, the impact pass that keeps a fix backward compatible, the run set a fix owes — the defect's own feature, the cross-cutting health slice, and whatever the impact named — and the RED-before-GREEN discipline. Routes to tdd for which test mode can observe which claim, test for bringing the stacks up, e2e for the feature file, your contribution guide for the pull-request rules, and your project architectural decisions for why the type exists."
---

# Fixing a defect

Five steps, in order, and the order is the whole procedure: **inspect → debug → impact → fix →
verify.** Each one produces something the next reads, which is what stops the loop collapsing into
"change something and see if the symptom goes away".

The last two exist for the same reason as the first three. **A fix is a change to working software**
— everything around the defect is behaviour somebody is relying on right now, and a fix that repairs
one path and breaks two is a worse outcome than the defect. Impact is where that is looked for, and
the run set is where it is proven.

## First: does this need a change folder at all?

**No, if the cause is obvious from the symptom.** A typo, a copy change, an off-by-one with a
stack trace pointing straight at it — go straight to a pull request. `your contribution guide` says so
and it stays true.

**Yes, if the cause is not obvious.** That is precisely when the inspect and debug steps have
something to record, and when skipping them produces a fix aimed at a guess.

```
your change-tool new change fix-<short-name> --schema bugfix
```

The test is not severity. A trivial-looking defect with a non-obvious cause needs the loop; a
serious one with an obvious cause does not.

**And it is not a defect at all if behaviour is correct and only the shape is wrong** — that is
`acme-refactor`, whose defining constraint is the opposite of this loop's: behaviour does not change.
When a fix and a restructure are both wanted, they are two commits, because a behavioural change
buried in a large move is a change nobody reviewed.

## Inspect — `report.md`

What is wrong, **without naming a cause**.

The symptom in the words of whoever saw it, the reproduction, and the blast radius. Data being
wrong, data being lost and data being *displayed* wrongly are three different severities and
should not be written as one.

**Leave the cause out, even when you think you know it.** A report that opens with a cause has
skipped the next step, and the fix will aim at the guess. The commonest way a defect survives
its own fix is that the first plausible explanation was never tested against a second one.

If it reproduces intermittently, say how often. "Sometimes" is not a reproduction, and a fix
verified against an intermittent failure that happened not to occur is not verified.

The artifacts are `report.md` → `diagnosis.md` → `flows.md` → `tasks.md` → `verification.md`, and
the schema in your bugfix change schema is what asks for each.

## Debug — `diagnosis.md`

The root cause, **and the observation that proves it**.

The evidence is the point of the artifact. A log line, a failing assertion, a query plan, a
network trace, a diff between a working and a broken state — something that distinguishes this
cause from the others that would produce the same symptom. **A diagnosis with no evidence is a
hypothesis**, and shipping a fix for a hypothesis is how a defect comes back wearing a different
symptom six weeks later.

Record what you ruled out, not just what you concluded. The alternatives are what a later reader
needs when the fix turns out to be wrong.

Then answer: **which test should have failed and did not?** That question is usually the one that
names the test to write next, and it is the difference between fixing this defect and fixing its
class.

**If the diagnosis shows the specification was wrong rather than the code, stop.** A requirement
is changing, so this is a `ui`, `fe` or `be` change. Say so in `diagnosis.md`, open the right
change, and let this one close — a bugfix that quietly rewrites a requirement is a specification
change nobody reviewed.

## Impact — the rest of `diagnosis.md`

**The blast radius of the FIX, not of the defect.** `report.md` asked who the defect hurts; this
asks who the *repair* could. They are different sets, and the second one is the one nobody looks
for.

Ask, of the thing about to change: **who else calls it, and who is relying on it behaving as it
does today?** Grep for the callers rather than remembering them. `acme-refactor`'s smell catalogue
is the same question asked of shape; this is it asked of behaviour.

Then answer the backward-compatibility question explicitly. In this repository it has a small
number of concrete shapes:

| What the fix touches | What is already relying on it |
|---|---|
| A response schema in your API contract package | your server framework drops every field the schema does not name, so **removing one is a silent 200 with the data gone**. A client on the previous release still asks for it — during a deploy both releases are live |
| A your data layer migration | One per pull request, backward compatible, never an edit to an applied one. Say what runs against a database still at the previous migration |
| A share link | A **capability token**, not a credential. Change how one is issued or read and every link already sent stops working — `acme-security` owns the surface |
| The session cookie's name, path or flags | Everybody signed in is signed out, and the hostname is shared with project's the external system application |
| A contract enum, or a state name | The wire agreement. Both sides parse it, and the ones already persisted do not migrate themselves |
| A shared **e2e step** in `steps/<domain>.steps.ts` | A step is a sentence, and **every feature that says that sentence changes meaning with it**. Grep `features/` for the sentence before editing the step |
| A seeded row or a catalogue shape | The console flows and the sharing flows mutate one shared catalogue; `e2e/your agent instructions`'s five rules say why a name is not a key |
| A your component library path primitive | Every composite that renders it, in both packages |

**"Nothing else calls it" is a fine answer — once it has been looked for.** Write it down either
way: this section is what chooses the run set below, and a fix whose impact was never mapped is a
fix verified against the symptom and nothing else.

**If the safe fix is bigger than this change, say so rather than shrinking the fix to fit.** A
conditional wrapped around the symptom to avoid touching a caller is the guarded-symptom failure
the [golden rules](../../../your agent instructions#golden-rules) name, and the follow-up nobody opens is how it
becomes permanent.

## Fix — `tasks.md`

**The first task is always the failing test. Not the fix.**

`acme-tdd` maps the claim to the mode that can observe it — geometry needs Chromium, copy needs
the pseudo-locale, serialisation needs `app.inject()`, a cross-layer claim needs e2e. A test
written in the wrong mode passes for the wrong reason and proves nothing.

**A defect proven at the e2e layer is a feature in `e2e/features/defects/`** — planned in
`flows.md` and written here, one file per fixed defect, named for the SYMPTOM, in the words
somebody would report it in. Where the claim belongs to a faster mode, `flows.md` says which and
why in a line, and no browser scenario is written: a suite that grows a journey for every unit
claim gets slow and vague, and `acme-tdd` owns the mapping. Those are the only
scenarios allowed to look narrow, because a defect is narrow; and `defects/` is where a red run
tells the reader "the defect came back" rather than "the product does not work". Never
`known-issues/`: that directory is for behaviour nobody has built, and a defect is behaviour that
is wrong. `acme-e2e` owns the file, the layers under it and the placement question.

Then the fix. **Fix the code, not the test and not the guard.** A guard failure means the code
broke an invariant. A failing pixel baseline means the render changed — open the PNG before
deciding which. A guard that produces a false positive is a higher-priority bug than one that
misses a case, but it is a *separate change*.

**Remove the cause; do not add a branch around the symptom.** A fix that adds a conditional, a flag or a null check where the diagnosis named a wrong seam leaves the defect in place and adds a reader's obstacle on top — and the next occurrence gets a second branch. If the honest fix is too large for this change, say so in `diagnosis.md` and open the follow-up; a guarded symptom with no follow-up is permanent ([golden rules](../../../your agent instructions#golden-rules)).

Keep it small. A bugfix that grows a feature is two changes, and the feature half has no
specification.

## Verify — `verification.md`

**RED before GREEN, and paste the RED.**

"It failed" is not evidence it failed for the right reason. A test that would also fail without
the defect present is testing something else, and it will keep passing after a regression
reintroduces the bug. The failure output is what shows the test was pointed at the right thing.

Then the same test passing — **and the run set the impact pass named.**

### The run set a fix owes

Two rows are owed by every fix, whatever it touched. The rest are chosen by the impact table
above, and **choosing them is the deliverable**: a slice list that does not follow from the impact
section is a guess about coverage.

| Owed | Run | Why this one |
|---|---|---|
| **always** | your validation command | lint, guards, typecheck, jsdom. The common path CI runs |
| **always** | `cd e2e && pnpm test:e2e cross-cutting/` | **The health slice.** The application answers, a visitor with no session lands at the front door, sign-in creates a session that survives a reload, the cookie stays scoped to our prefix. 26 scenarios, the fastest useful slice in the suite — and if these are red, no other result in the run means anything |
| **always** | `pnpm test:e2e defects/<symptom>` | The feature written in the Fix step. This is the regression, and it is the one that must go RED before it goes GREEN |
| a component's appearance | your browser tests, plus the journeys that render it | jsdom cannot see geometry (`acme-tdd`) |
| a route, a service or a response schema | the `journeys/` that call it | Every layer's own tests pass with the seam broken — that is the whole argument for the layer existing |
| a permission, a role, a session or a cookie | `pnpm test:e2e roles/` | 11 scenarios over four roles: what each may and may not reach. Cheap, and the failure it catches is somebody seeing a surface that is not theirs |
| the staff console | `pnpm test:e2e internal/` | Staff curation. No estate agent is affected, which is exactly why nobody notices |
| a migration | the slice above, against a database migrated from the **previous** release | Both releases are live during a deploy |
| the assistant or the generation pipeline | `pnpm test:e2e:chat`; `@generation` needs LocalStack, a worker and a key | Both are excluded from the default run, so the default run says nothing about them |

**Reseed before the run and select a slice, not the suite.** Leftovers do not look like leftovers —
three curation flows that pass in 9.7s on a fresh database hung for their full 270s budget on a
dirty one. The whole suite is CI's job behind the your CI label for e2e checks label; `e2e/README.md`'s *Running a
slice* is the reference and `acme-test` owns bringing the stacks up.

**A run that was started and killed is UNRUN.** One change had to record its end-to-end layer that
way after forty minutes produced no pass count and no verdict. Write that, not "e2e passed".

**Say what is still not proven.** A fix verified in jsdom for a symptom reported in a browser has
a gap, and the gap belongs in writing rather than in someone's memory.

## Red flags

- A `report.md` that names a cause. The debug step has been skipped.
- A `diagnosis.md` with no observation. It is a hypothesis wearing a conclusion's clothes.
- A fix committed before the test. Nothing then proves the test would have caught it.
- A `verification.md` with no RED output.
- A guard or a test edited to make the failure stop. That is not a fix.
- A bugfix that adds a requirement. Requirements are `ui`, `fe` and `be` work.
- A `diagnosis.md` with no "what else depends on this". The fix has been scoped by the symptom alone.
- A fix that removes or renames a response-schema field, an enum value or a cookie attribute with no sentence about the release still running beside it.
- A `verification.md` whose only run is the defect's own feature. That proves the symptom is gone, not that the fix is safe.
- A shared step definition edited without grepping `features/` for the sentence. One edit, every feature that says it.
- "e2e passed" with no selection named. Which slice ran is the coverage claim.

## Verification

- [ ] The change exists on the `acme-bugfix` schema, or the cause was genuinely obvious and it is a plain pull request
- [ ] `report.md` describes the symptom and names no cause
- [ ] `diagnosis.md` carries the observation that proves the cause, and what it ruled out
- [ ] `diagnosis.md` names the test that should have failed
- [ ] `diagnosis.md` names what else depends on the thing being changed, and answers the backward-compatibility question — "nothing else calls it" counts, once it has been looked for
- [ ] The failing test was written first and its RED output is in `verification.md`
- [ ] your validation command green
- [ ] `cross-cutting/` green, and the defect's own `defects/<symptom>` feature RED then GREEN
- [ ] Every slice the impact section named was run, or is listed as not run with the reason
- [ ] `verification.md` names the selections that ran — not "e2e passed" — and states what is still not proven
