# e2e — workspace rules

The layer-specific rules for this workspace are owned by the documents and skills the
root [AGENTS.md](../AGENTS.md) routes to. This file adds nothing they already state;
it exists so the workspace itself is self-describing for an agent that starts here.

- Stack and commands: this workspace's `package.json` and the root [README](../README.md)
- Writing or repairing a flow: `e2e` (.agents/skills); running one: `test`
- Tests: `docs/engineering/testing.md`
- What this workspace contains today: [README.md](./README.md)

## The layer a change belongs in

Flows here are **BDD**: Gherkin in `features/`, run by playwright-bdd. The generated Playwright
specs under `.features-gen/` are build output — never edited, never committed.

The whole design is one rule: **a reader of a feature file must never need to open anything
else.** Ask in order, first yes wins:

1. Is it a sentence about what the product must do, in the words the business uses? → `features/`
2. Is it the MEANING of such a sentence, in code — including its assertions? → `steps/`
3. Is it how you REACH something on a screen — a locator, a wait, the reasoning behind either? → `pom/`
4. Is it needed by more than one step? → `fixtures/`
5. Is it needed by more than one layer — an identity, the API origin, a helper? → `lib/`

A step that constructs a page object, a page object that asserts a business outcome, or a feature
file naming a selector, a URL, a status code or an id, is in the wrong layer.

## Where a new feature goes

`features/` is split into seven directories by **what a failure means**, not by the screen it
touched — that is the question somebody reading a red run is actually asking.
[README.md](./README.md) lists what is in each one today. The placement is one question asked in
order, and the first yes wins:

1. Does a real person walk it end to end, across several steps? → `journeys/`
2. Is the claim about what a ROLE may or may not reach? → `roles/`
3. Is it a surface only staff ever see? → `internal/`
4. Must it hold on every screen — health, sign-in, a layout treatment? → `cross-cutting/`
5. Is it a defect that has been fixed? → `defects/`, one feature per defect, named for the symptom
6. Is the behaviour not built yet? → `known-issues/`, tagged `@known-issue` with the blocker named
7. Otherwise → `capabilities/`

A feature that answers yes twice is usually two features. One that answers yes to none is usually
asserting an implementation detail, and `tdd` places that one layer down.

## The five rules the split rests on

- **A scenario creates what it asserts on.** There is no per-worker database isolation —
  `lib/suite.ts` says the fixture is planned and not built — so a scenario that borrows a seeded
  row is a scenario that fails after somebody else's run. Where borrowing is unavoidable, the step
  says why at the point it borrows.
- **Never name a seeded row; ASK the API, and page.** The wall paginates and the suite's own
  creation scenarios push rows in front of any named card. Find rows through the API or by
  position, never by a literal name from the seed.
- **Reseed before a run.** The same missing isolation, from the running side: flows that mutate
  shared data leave leftovers that do not look like leftovers — flows that pass on a freshly
  seeded database may hang on a dirty one.
- **`known-issues/` is not a parking space for red tests.** A scenario goes there only when the
  thing it asserts has not been built and the blocker can be named in a sentence. A scenario that
  fails because the application is *wrong* is a defect, and belongs to `bugfix`.
- **A name is not a key.** Seeded data may carry duplicate names, so a step that addresses a row
  by name may address several. Create the subject, or address it by position, and say which in
  the step.

## The costs are declared in the feature file

A tag is not decoration: `steps/hooks.ts` is where each one becomes behaviour, so the cost of a
scenario is visible to the person reading it. `@spends` calls an expensive service and is billed;
`@generation` needs external infrastructure and an extended budget; `@money` marks a scenario
asserting something must NOT spend. The full table is [README.md](./README.md)'s.

There is no `@desktop-only`: the suite runs in ONE Playwright project (`flows`, desktop Chromium),
so a scenario cannot be repeated across a device matrix in the first place. A guarantee in the
shape of the run beats one every expensive scenario has to remember to opt into — and a claim that
IS about a width says so in its own words, with a step that sets it.

## A run is a slice

**Select the slice your change touches; the suite, whole, is CI's job.** Many scenarios on one
worker — a journey may take 120 seconds and a `@generation` scenario 420 — is tens of minutes with
nothing wrong. A run that produces no verdict has spent the time and proved nothing.

The seven directories are the suites, and Playwright's own selectors do the rest —
`your e2e test command journeys/`, `journeys/sharing`, `--grep "@layout"`, `--last-failed`.
The default run excludes `@spends`; `your e2e test command:spends` is the one that bills.
[README.md](./README.md)'s *Running a slice* is the table.

Two things this does not weaken. `bddgen` runs whatever is selected, so a sentence with no step
behind it still fails before a browser starts. And a narrowed run is not a green suite: say which
slice ran, and put the appropriate CI label on the pull request for the rest.

When the selection is bigger than any single feature, `run-scope.reporter.ts` says so on the run's
first line. It is a notice rather than a gate — a full run is exactly right before a merge, and a
flag to silence a warning is a flag that accumulates.

## The import depth

Steps and page objects sit one directory below the e2e root, features two below:

```ts
import { expect, test } from '../fixtures/index.js';   // never '@playwright/test'
import { API_ORIGIN } from '../lib/api.js';
```

`fixtures/index.ts` merges playwright-bdd's test instance with `lib/suite.ts` using `mergeTests`.
Spreading one into the other instead looks like it works and silently drops `suite`'s `auto`
fixtures — the run then starts answering 429 partway through, and the failures read as the API
being broken.

A wrong depth fails at COLLECTION with `Cannot find module`, which takes the whole run down rather
than one file — the same shape as a module that throws while loading.
