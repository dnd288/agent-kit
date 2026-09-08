---
name: refactor
description: "Changing the shape of project code without changing what it does — the loop that makes that safe, at any layer. Load this when asked to refactor, clean up, tidy, restructure, simplify, split, extract, decompose, modularise or 'improve the quality of' existing code; when a component, service, route, hook, store or module has grown too big, too complex, too coupled or untestable; when your guard checks reports a composite over its complexity budget or a ledger entry to pay down; when a screen has too many props or too many useState calls; when logic sits in the wrong layer — a fetch in a composite, business logic in a route handler, a store holding a database row, a DTO leaking into your UI package; when duplication has reached a third copy; when asked to make something testable, to separate concerns, to fix the data flow, or to bring a file up to standard. Owns the MEASURE → PIN → MOVE → PROVE → RECORD loop, the rule that behaviour does not change, the characterisation-harness requirement, the phase discipline that keeps structure, behaviour and interface in separate commits, and the smell catalogue from your database to the pixel. Routes to ui-development §9 for splitting a screen, state-management for where a value lives, api-contract for the service boundary, tdd for which mode can observe a claim, and review for the gate."
---

# Refactoring project code

**The defining constraint: behaviour does not change.** Everything below exists to make that
constraint checkable rather than hoped for. A change that improves the shape *and* fixes something
is two changes, and the reason to separate them is not tidiness — it is that a behavioural change
buried in a thousand-line move is a change nobody reviewed.

That constraint is also what makes this the highest-leverage work in the repository. your agent instructions's
[golden rules](../../../your agent instructions#golden-rules) set the standard — *the codebase gets simpler as the
product gets larger*, and *the test is not "is this good code" but "is the next change cheaper than
this one was"*. Refactoring is the only activity whose entire purpose is that test.

## The loop

**MEASURE → PIN → MOVE → PROVE → RECORD.** Each step produces what the next one reads, which is
what stops a refactor collapsing into "move things around until it feels better".

Skipping **PIN** is the one failure that matters. Every other shortcut costs time; that one silently
ships a behavioural change wearing a refactor's commit message.

## First: is this a refactor at all?

Three things get called refactoring and only one is:

| What you are doing | It is | Skill |
|---|---|---|
| Same behaviour, different shape | a **refactor** | this one |
| Different behaviour, because the old one was wrong | a **bugfix** | `bugfix` |
| Different behaviour, because we want more of it | a **feature** | `spec-workflow` |

**Mixed is the normal case, and the answer is to split it.** Restructuring a screen and, while you
are in there, correcting a gate that was wrong: two commits, the refactor first, and the fix second
where it is a two-line diff a reviewer can actually see. Doing the fix first is also fine. Doing
both at once is not, because the harness in PIN cannot hold behaviour still while behaviour is
deliberately moving.

**If the refactor reveals a defect, stop and write it down** before deciding which order to take
them in. A defect found mid-move is the most valuable thing a refactor produces, and it is also the
easiest to accidentally "fix" as part of the move and never record.

## Does it need a change folder?

The same threshold `bugfix` uses, for the same reason:

**No** — go straight to a pull request when the move is local and mechanical: extracting a function,
renaming, deleting a dead branch, splitting one file into two, replacing a duplicated block with a
call.

**Yes** — when the refactor changes an interface other workspaces name, moves a value between layers
of your project architectural decisions's ladder, or takes more than
one commit to land safely. That is exactly when there is something to specify and someone to warn.

```
your change-tool new change refactor-<short-name> --schema <layer>
```

**There is no `refactor` schema, deliberately.** A change has a type and it is one layer
(your project architectural decisions) — `ui` for your component library path and
your UI package path, `fe` for your app path, `be` for your API app path and the data model, `cross-cutting` when it
genuinely crosses. A refactor that cannot pick one is a refactor that is doing too much at once, and
that is a finding rather than a reason for a sixth type.

---

## 1. MEASURE — find the bad thing; do not guess

**A refactor aimed at a feeling is a rewrite with better manners.** Start from a signal that
existed before you had an opinion.

| Signal | How to get it | What it means |
|---|---|---|
| A composite past its budget | your guard checks → `composite-complexity` | 24 props / 6 stateful hooks / 800 code lines / complexity 45. Over any of them, the component is several components sharing a file (your project architectural decisions) |
| A standing debt entry | `RECORDED` in `packages/guard/src/checks/composite-complexity.ts` | The repository's own list of what it already knows is wrong, with numbers. Paying one down is always legitimate work |
| A hook that cannot work | your guard checks → `component-performance` | A memo with no dependency array, or a dependency built inside its own array |
| A boundary already crossed | your guard checks → `ui-stays-presentational`, `service-purity`, `cross-app-imports`, `client-api-boundary`, `core-neutrality`, `vanilla-store-factories` | The build is red; the refactor is not optional |
| A missing surface | your guard checks → `component-has-surfaces`, `component-has-stories` (warnings) | A component with no quality test has nothing measuring its real rendering |
| The test that is too big | `wc -c` on the file's test | A test file that is the largest in its package is telling you every claim needs the whole thing rendered |
| A waterfall or a duplicate read | read the page's `await`s | Sequential independent server reads, or two callers of an uncached reader in one render (state-and-data-flow.md) |
| The third copy | the change you are about to make touches three files identically | *Two is a coincidence; three is a missing abstraction* — and the abstraction **replaces** all three |

`references/smells.md` is the full catalogue, layer by layer from your database to the pixel, each smell
with the rule that owns it and — where there is one — the instance already in this repository.

**Write down what you measured before you touch anything.** The number is what PROVE compares
against, and "it feels much cleaner" is not a number.

## 2. PIN — behaviour cannot be held still by hoping

**You may not refactor what nothing observes.** Before the first move, there is a test that fails if
behaviour changes. One of three cases:

**a. A harness already exists.** Check what it imports. A test that names only the module's public
surface and drives everything else through it covers any internal move you are about to make — and
`result-screen.test.tsx` is the worked example: 42 KB, importing exactly two symbols, which is why
the whole of phase 1 was verifiable without editing a line of it.

**b. A harness exists but reaches inside.** A test that imports an internal helper, or asserts on a
private shape, pins the *implementation* rather than the behaviour and will fail on a correct move.
Rewrite those assertions to go through the public surface **first, as their own commit**, and watch
them still pass. Only then start moving.

**c. There is no harness.** Write one, and write it against **current behaviour including its
bugs** — a characterisation test records what the code does, not what it should do. If something it
records is wrong, that is a defect to raise separately (see the mixed case above); encoding the
correct behaviour now means the test goes green on code that has never run.

`tdd` owns which of the four modes can observe which claim, and the answer decides where the
harness lives: a pure function is a unit test, geometry and computed tokens are Chromium
(your browser tests), a service is a unit test, a route is a route test, and a journey across the
whole stack is `e2e/`.

**The harness is read-only for the duration.** Editing it during a move is how a safety net becomes
a rubber stamp. If a move genuinely requires the harness to change, that move is not a refactor.

## 3. MOVE — structure, behaviour and interface never share a commit

Three kinds of change, in this order, and each finishes before the next starts:

1. **Structure** — code moves; nothing that calls it changes. Extract a function, split a file, pull
   a widget out, lift a service out of a handler. Re-export moved symbols from where they used to
   live so the public surface is byte-identical and the harness stays valid.
2. **Behaviour** — nothing here, by definition. If you have some, it belongs to a different change.
3. **Interface** — the signature, the props, the response schema, the exported names. The breaking
   one. Callers update in the same commit, and this is the phase that needs the change folder, the
   review and the e2e slice.

**Run the harness after every move and commit each one.** A red run then names the move that did it,
which is worth far more than the tidiness of one large commit. The commit message says what moved
and why the next change is cheaper — that message is the most durable documentation the repository
has ([conventions](../../../your engineering conventions)).

**Delete rather than deprecate.** A prop still declared and no longer read typechecks and does
nothing, which is precisely how `result-screen`'s browse surface broke: four props passed into
silence with a green build behind them. [Golden rule 2](../../../your agent instructions#golden-rules) is the
general form — *a new abstraction earns its place by deleting*; if the move adds a layer and removes
nothing, say what it will remove and when, or do not add it.

**Splitting a screen has its own procedure**, phase by phase, in
[`ui-development` §9](../ui-development/SKILL.md) — *Refactoring one that is already over
budget*. Do not restate it here; it owns the widget seam, the screen store and the prop groups.

## 4. PROVE — the same standard as any other change

Green is not the claim. The claim is **"behaviour did not change, and here is what would have caught
it if it had"**.

- your validation command — lint, guards, typecheck, jsdom tests. The common path, and the floor.
- your browser tests — if anything about appearance or geometry could have moved. It is not in CI
  by default (your project architectural decisions, amended), so a your UI package path
  refactor that skips it has verified nothing about the layer whose whole subject is appearance.
- The e2e slice the code is on — `cd e2e && pnpm test:e2e journeys/<name>`, or the your CI label for e2e checks
  label. Required for any interface-phase change and any move across a layer boundary.
- **The measurement from step 1, taken again.** A refactor with no movement in the number it was
  started for did not do the thing it was for.

**Say what you did not run and why.** `review`'s standard applies unchanged, and a refactor is
the change type where "the tests pass" is most likely to be true and least likely to be sufficient.

For anything user-visible, `verify`'s adversarial pass — CLAIM → EXTRACT → DOUBT → RECONCILE →
STOP — attacks the claim rather than confirming it.

## 5. RECORD — leave the map matching the ground

A refactor that does not update what described the old shape has moved the problem into the
documentation:

- **The ledger.** `composite-complexity` fails when a recorded number improves and its entry has not
  followed, and names the number to write. Delete the entry when the last measure is inside budget.
  **Adding an entry is never how a refactor is finished.**
- **The documents that described the old shape.** your agent instructions's boundaries, the workspace
  your agent instructions, `docs/engineering/*`. `conventions.md` is explicit: *changing a rule means checking
  the skills that cite it, in the same pull request.*
- **A decision record**, when the refactor establishes a rule rather than applying one — a new
  boundary, a new place a kind of value lives, a pattern the next screen should copy. Seven sections,
  and *What would change our mind* is mandatory. your project architectural decisions
  is the recent example, and it exists because splitting one screen decided how every later screen
  is shaped.
- **Nothing, if it applied an existing rule.** Most refactors are this. Do not write a record for
  obeying one.

---

## When NOT to refactor

Each of these is a real answer, and saying it is better work than a plausible-looking move:

- **The code is ugly but stable, and nothing is about to touch it.** The test is whether the *next*
  change is cheaper — with no next change, the benefit is zero and the risk is not.
- **You cannot pin the behaviour and cannot afford to.** Then the honest sequence is: write the
  harness as its own piece of work, land it, and refactor later. Not: move carefully and hope.
- **The requirement is about to change.** Restructuring around a shape that is being replaced next
  sprint spends the budget twice.
- **You are one commit into a feature.** Finish it, then refactor. A half-built feature has no stable
  behaviour to hold still.
- **The "duplication" is a coincidence.** Two things that look alike and change for different reasons
  are not duplication, and merging them creates a shared thing with two masters — which is worse than
  the copies, and harder to undo.
- **It would be a rewrite.** Rewriting is not refactoring: it discards the behaviour nobody wrote
  down, which is most of what old code is worth.

## Before you finish

- [ ] Behaviour is unchanged, and a test would have failed if it were not
- [ ] The harness was not edited during the structural phases
- [ ] Structure, behaviour and interface are in separate commits
- [ ] Moved symbols are re-exported, or every caller is updated in the same commit
- [ ] The measurement from step 1, taken again and stated
- [ ] your validation command; your browser tests if appearance could have moved; the e2e slice if a
      boundary or an interface did
- [ ] The ledger followed the improvement down, and no entry was added
- [ ] Every document and skill that described the old shape is updated in this change
- [ ] The commit message says why the next change is now cheaper

## Read next

- `references/smells.md` — the catalogue, layer by layer, with the repository's own instances.
- [`ui-development` §9](../ui-development/SKILL.md) — splitting an over-budget screen.
- [`state-management`](../state-management/SKILL.md) — where a value lives, and the
  anti-pattern catalogue the guard scans for.
- [`api-contract`](../api-contract/SKILL.md) — the service boundary, response schemas,
  migrations.
- [`tdd`](../tdd/SKILL.md) / [`test`](../test/SKILL.md) — which mode observes which
  claim, and how to run it.
- [`review`](review) — the gate, and the architecture axis that judges whether
  the seam was the right one.
- [your agent instructions golden rules](../../../your agent instructions#golden-rules) — the standard the whole loop serves.
