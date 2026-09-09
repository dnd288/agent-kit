# e2e

What this workspace contains today, and nothing more — this file is the inventory, not
the documentation. The rules live in the documents and skills the workspace
[AGENTS.md](./AGENTS.md) routes to.

- **e2e** — [playwright-bdd](https://github.com/vitalets/playwright-bdd) over Playwright. A
  `.feature` file is the source; the Playwright spec under `.features-gen/` is build output and is
  neither edited nor committed
- Scripts: `your e2e test command` (generates, then runs **everything that does not call an
  expensive service**), `your e2e test command:spends` (the whole suite, billed scenarios
  included), `bdd:gen`, `bdd:report` (serves the business report so its traces open), `typecheck`

## Running a slice

**A local run is a slice. The suite, whole, is CI's job.** The suite runs on **one worker**
locally, because there is no per-worker database isolation yet. A journey is allowed 120 seconds
and a `@generation` scenario 420, so a full local run is tens of minutes before anything is wrong
with it.

Select from `your e2e path/`. Anything under `your e2e test command` takes Playwright's own
selectors:

| To run | Command |
|---|---|
| One directory — the seven under `features/` are the suites | `your e2e test command journeys/` |
| One feature | `your e2e test command journeys/first-project` |
| One scenario, by title | `your e2e test command --grep "creates"` |
| One tag | `your e2e test command --grep "@layout"` |
| What failed last time | `your e2e test command --last-failed` |
| A suspected race, five times over | `your e2e test command --repeat-each=5 --grep "<scenario>"` |
| The billed scenarios too — the whole suite | `your e2e test command:spends` |

A positional filter is a regular expression over the **generated** path, and `.features-gen/`
mirrors `features/` — so a directory or a feature name is enough, and no extension is needed.
`bddgen` runs first whatever you select, so a sentence with no step behind it still fails before a
browser starts: narrowing the run never narrows *that* check.

**Which slice covers a change** is the same question as which failure would mean what — the
directory table below is the answer.

The run says so itself. When the selection is larger than any single feature,
`run-scope.reporter.ts` prints the count and these selectors before the first browser
opens — the last moment Ctrl-C is still cheap. It never fires on CI, and it refuses nothing.

## The layout

| Directory | Holds | Read by |
|---|---|---|
| `features/` | Gherkin. What the product must do, in the words the business uses | Anyone. This is the point |
| `steps/` | The Given/When/Then implementations, grouped by domain, plus the tag hooks | Engineers |
| `pom/` | Page objects. Where a locator's REASONING lives, so a step stays a sentence | Engineers |
| `fixtures/` | The BDD test instance: the page objects and the per-scenario `world` | Engineers |
| `lib/` | What more than one layer needs — identities, the API origin, helpers | Engineers |
| `reports/` | The run's output: the business flows document, and the same run as Cucumber Messages | Anyone |
| `.features-gen/` | Generated. Never edited, never committed | Nobody |

**A sentence in a feature file with no step behind it fails the GENERATION**, before anything
runs (`missingSteps: 'fail-on-gen'`). That is the property the whole layout rests on: a
business-readable suite that can quietly skip a scenario is a suite that starts describing
behaviour nobody implemented, which is worse than having no suite.

## Two projects

`features/` is the whole suite. There is one project in the default run — `flows`, desktop
Chromium — and a `seed` project that populates the database before any scenario runs.

**One worker locally, four on CI.** There is no per-worker database isolation (`lib/suite.ts` says
why) and this suite mutates what it reads. Four workers on one seeded database are four copies of
every mutating flow racing each other, and the failures land on whichever scenario looked second.
CI keeps four because its database is its own, and `retries: 2` covers the machine being slow.

## The tags, and what they cost

A tag in a feature file is not decoration: `steps/hooks.ts` is where each one becomes behaviour,
so the cost of a scenario is visible in the scenario rather than buried in code.

| Tag | Means |
|---|---|
| `@journey`, `@roles`, `@internal`, `@cross-cutting` | 120s budget — several routes, each compiling on first hit |
| `@capability` | 120s budget, same as a journey. A `capabilities/` feature asserts one capability rather than walking a path, but still signs in and opens a screen |
| `@generation` | Needs external infrastructure and an extended budget; 420s. Skipped, with the reason, when the stack is absent |
| `@spends` | Calls an expensive external service. **Excluded by default** — `your e2e test command:spends` is what opts in |
| `@money` | Asserts that something must NOT spend. Costs nothing to make and is the only guard against a regression that bills silently |
| `@optional-infra` | Excluded by default; a named script runs it when the infrastructure is configured |
| `@known-issue` | Skipped: the behaviour is not shipped. The feature names what must ship first |

**`your e2e test command` excludes `@spends`.** The default run costs nothing; `your e2e test command:spends` is the whole suite and the only thing that bills.

## How the features are split

Seven directories under `features/`, split by **what a failure means** rather than by which screen
it touched — the question somebody reading a red run is actually asking.

| Directory | Holds | A failure here means |
|---|---|---|
| `journeys/` | Multi-step paths a real person walks end to end | The product does not work. These are the money paths |
| `roles/` | One file per role, asserting what that role may and may not reach | Somebody sees a surface that is not theirs, or is refused one that is |
| `internal/` | Staff-only console surfaces | Staff cannot operate. No end user is affected |
| `capabilities/` | One capability, asserted on its own, where no journey needs it | A capability regressed without breaking a path |
| `cross-cutting/` | What must hold on every screen: health, sign-in, layout | Everything is affected — read this first when a run is broadly red |
| `defects/` | One feature per fixed defect, named for the symptom | The defect came back. These are the only scenarios that may look narrow |
| `known-issues/` | Written against behaviour that does not exist yet, skipped with the blocker named | Nothing — they are skipped. They exist so the day the gap closes, the proof is already written |

Which directory a NEW feature goes in, and the rules the split rests on, are
[AGENTS.md](./AGENTS.md)'s. This file says only what is here today.

### What is in each one today

<!-- Fill in as features are added. Example row:
| `journeys/first-project` | A user creates a project, finds it on the wall, opens it |
-->

| Feature | The business claim it defends |
|---|---|
| `cross-cutting/health` | The front door, the sign-in screen, and both halves of the application answering for themselves |

`lib/` carries what more than one layer needs:

| Module | Holds | Why it is not inlined |
|---|---|---|
| `suite.ts` | `test` and `expect`, and the per-test `x-forwarded-for` fixture | Every spec imports the first two from here; the fixture is what keeps the sign-in cap from failing the run |
| `accounts.ts` | The seeded identities and their credentials | One module rather than a copy per spec |

**A scenario never names a seeded row.** Ask the API which rows exist, because the wall
paginates and the suite's own creation tests push rows in front of any named card.

## The reports

`reports/business-flows.html` is the deliverable, and its unit is the SCENARIO: the Given/When/Then
somebody wrote, in order, with the screenshot and the trace hung off the step that failed. It is
the answer to "what does a green run actually prove" — the steps are sentences about the business
rather than about locators.

`reports/messages.ndjson` is the same run as Cucumber Messages: machine-readable, diffable between
runs, and the input any other tool would want. The Playwright HTML report stays alongside them
because it is the only one carrying the trace viewer.

**Every failure carries an AI-fix attachment** (`aiFix.promptAttachment`): the feature, the failing
step, that step's own source and the page's ARIA snapshot, assembled into a prompt. A red run is
something an agent can act on rather than something it has to reconstruct.
