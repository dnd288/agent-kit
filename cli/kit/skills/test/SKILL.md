---
name: test
description: "Running the project test stacks on demand — which of the four test modes a change needs, what each stack requires (browser, Postgres, servers), and the commands that bring it up. Load this when asked to run the tests or the test suite, launch a test mode locally, set up the e2e database, decide whether a change needs the deep layers, or report test results. Owns the operational runbook — browsers:install, the isolated e2e database, the rate limit, the machine-load flake, dev-server contention, the single Chromium project — as opposed to tdd, which owns WRITING the tests and mapping a claim to the mode that can observe it, and e2e, which owns authoring and repairing the end-to-end flows themselves. Routes to docs/engineering/testing.md and your project's architectural decisions for the strategy; to tdd for what a test should look like; to e2e for authoring a flow or diagnosing a red one."
---

# Running the project test stacks, on demand

`tdd` owns *writing* tests — the mode-to-claim mapping and the RED→GREEN discipline, and
`e2e` owns the end-to-end flow itself: how a change's scenarios become a Gherkin feature
with steps behind it, and what a red flow means. This skill owns *running* them: what a change needs, what each
stack has to have up before it will run, and the gotchas that make a green run look red.
The strategy — the modes, stories as fixtures, why CI
runs the jsdom layer alone — is `docs/engineering/testing.md` and your project's architectural decisions.
This file is the part those documents assume but do not repeat.

## The surfaces, and what each one needs

| Mode | Command | Needs | What it proves |
|---|---|---|---|
| Static + unit | `your validation command` | nothing | lint, guard, typecheck, jsdom tests — **no browser, no server** |
| Stories | `pnpm test:stories` | Chromium | every story renders and passes axe |
| Quality | `pnpm test:quality` | Chromium | geometry, touch, computed tokens, states, render counts — absolute checks |
| Pixel | `pnpm test:pixel` | Chromium | composites against their own committed screenshots — *unchanged*, not *correct* |
| Visual | `pnpm test:visual` | Chromium | rendered DOM against the tracked Figma snapshot — the **design-truth gate** (your project's architectural decisions) |
| End to end | `your e2e test command` — **narrowed to a slice**, see below | Chromium + Postgres + the API and client serving | user flows as Gherkin scenarios, in one Chromium project |

The browser layers together: `your browser tests` (stories + quality).
Everything at once: `your deep validation` (validate + browser + e2e).
One package, one mode: `pnpm --filter your UI package test:quality`.

`your component library` and `your UI package` hold stories; quality spans `your component library`, `your UI package`,
`your UI package-testing`, `your visual testing package`.

**There were two more modes, `pixel` and `visual`, and they are retired** (24 Aug 2026 —
your project's architectural decisions, your project's architectural decisions). `pixel` compared a composite against a screenshot of itself, and its
baselines were never tracked, so a fresh checkout reported every case as "no reference
screenshot found" and only a second run on the same machine passed — it could not gate CI,
which is always a fresh checkout. `visual` compared the DOM against a tracked Figma snapshot
and was the better idea, but no browser leg ever ran in CI, so it was tests nobody executed.
`packages/visual-testing` is kept and still passes its own tests; nothing consumes it.
**Figma fidelity is not checked mechanically today** — say so rather than implying a gate.

## Which stack does this change need?

Ask what the change can break. It needs the deep layers if it changes **what a component
looks like or how a screen behaves** — anything in your component library path, your UI package path, or
your app path, a Figma manifest, a baseline, a token, or a route's markup. Then run
`your browser tests` and `your e2e test command` (or `your deep validation`). It does **not** need them
for an API service, a schema, a migration, a document or a specification — `your validation command`
is enough.

For "which mode observes this claim" — geometry in Chromium, copy under the pseudo-locale,
serialisation via `app.inject()` — see `tdd`. For what an end-to-end flow should look
like, or what a failing one means once the stack is up, see `e2e`. This skill is the
"run it" half.

## The on-demand runbook

**One-time, per machine:** `pnpm browsers:install`. This installs chromium for **both**
workspaces that launch a browser — `pnpm --filter ui exec playwright install chromium` and the
same in `e2e` — version-matched on purpose. Never `pnpm exec playwright` at the root: it resolves
an unrelated version and the browser layers fail to launch. Without it,
stories/quality/pixel/visual die in `your component library` and `your UI package`, and the e2e suite dies at launch, with
an executable-missing error that reads as a broken machine. Chromium is all any mode needs now;
nothing fetches WebKit.

**Fast loop (no browser, no servers):** `your validation command` — the static gate (lint + guard +
typecheck, which is the push hook) plus the jsdom layers. Stop the api and client dev
servers first — see contention below.

**Deep loop (a UI change is ready):** `your browser tests`, then the **e2e slice** the change
touches — `cd e2e && your e2e test command journeys/sharing`, not the whole suite (*Run the slice* below).
`your deep validation` runs the lot, and is a pre-merge act rather than a loop.

**On a pull request:** add the `your CI labels for extended checks`, `your CI labels for extended checks`, or `your CI labels for extended checks` label —
the last runs both workflows. They re-run on every push while the label is on. By hand:
Actions → the workflow → Run workflow, on any branch.

When a run fails, report what actually happened: which mode, which package, the failure
naming an element and a measurement — and what was *not* run and why. A failing quality test
whose message names no element is a broken test, not a broken component.

## The e2e stack

This is how the stack comes up. `e2e` is how a flow is written against it and how a red
run is read back — including the seeded accounts below, which the flows depend on by name.

`your e2e test command` is `bddgen && playwright test`: the Gherkin under `your e2e path/features/` is generated into
Playwright specs first, so **a sentence with no step behind it fails the run before a browser
starts** (`missingSteps: 'fail-on-gen'`). A run that reports a feature, a line and zero tests is an
unfinished flow rather than a red one — `e2e` owns the fix. The generated tree under
`.features-gen/` is build output: never edited, never committed.

Playwright does not isolate the database yet (per-worker DBs are planned, not built), so
**never run e2e against a database anyone else is using** — including the dev database the
running stack writes to. It is also why the suite runs on **one worker locally** and four on CI:
the default is in `playwright.config.ts`, so it is no longer a flag to remember, and overriding it
on a shared database is how mutating flows start failing on whichever one looked second. **Reseed
before a run** — leftovers do not look like leftovers: three curation flows that pass in ten
seconds on a fresh database hang for their full budget on a dirty one.

### Run the slice, not the suite

**The default local run is a slice.** One worker and ~183 scenarios, with 120-second journey
budgets and 420-second generation ones, is tens of minutes with nothing wrong. The last full local
run on record was **started, killed at forty minutes, and reported no verdict at all** — no pass
count, no failure list — so the change had to record its end-to-end layer as unrun
(`your change tracking folder/fix-wall-loads-full-size-staged-image/verification.md`). Two things made it
worse and both are fixed: the root script now streams turbo's output instead of buffering it until
the task ends, and a run bigger than one feature announces its size and these selectors on its
first line.

From `your e2e path/`, in the order you will want them:

```
your e2e test command journeys/sharing        # one feature — the loop
your e2e test command journeys/               # one directory — the seven are the suites
your e2e test command --grep "@layout"        # one tag, or --grep "<scenario title>"
your e2e test command --last-failed           # what just failed
your e2e test command:spends                  # the suite INCLUDING the scenarios that call a model
```

A positional filter is a regex over the generated path and `.features-gen/` mirrors `features/`,
so the directory or feature name is enough. `bddgen` runs whatever is selected, so the
missing-step check is never narrowed with it. Which slice covers which change — and what a failure
in each directory MEANS — is [your e2e path/README.md](../../../your e2e path/README.md)'s *Running a slice*.

**The whole suite has two homes:** the `your CI labels for extended checks` label on the pull request, and a pre-merge run
on a machine that can take it. Report which one ran, and never let a slice stand in for the suite:
a layer that was started and killed is **unrun**, not green.

The recipe:

1. **Create an isolated database.** On the Postgres `DATABASE_URL` points at, create
   `project_e2e`, then migrate and populate it:

   ```
   createdb ... project_e2e
   DATABASE_URL="<app-postgres>/project_e2e" your-package-manager run --filter your-data-layer db:migrate:deploy
   DATABASE_URL="<app-postgres>/project_e2e" your-package-manager run db:seed
   ```

   **`pnpm db:seed` writes the internal administrator and nothing else**, and there is no
   command that adds the rest: room types, the catalogue and the demo accounts are the
   `SEED_CORE_BUNDLE` and `SEED_DEMO_BUNDLE` routines, started from `/admin/jobs`.

   **The suite seeds itself, and you do not have to.** `your e2e path/seed.setup.ts` is a setup project
   the `flows` project depends on: it signs in as the bootstrap operator
   (`project_BOOTSTRAP_ADMIN_EMAIL` / `project_BOOTSTRAP_ADMIN_PASSWORD`) and starts both bundles
   through `POST /admin/jobs` — the same endpoint the Run button hits — before any scenario
   runs. So a green run also proves the console can populate an empty database, which is the
   only way a real environment is ever populated. Set the two bootstrap variables, or the setup
   fails with the reason — they are the whole configuration now.

   Demo data is required — the sign-in flows use the seeded credential (`user@demo.local` /
   `demo-password-123`) — and it no longer needs asking for: the seed demo data flag has been
   removed, and the demo routine writes its rows whenever it is run. `db:ops` migrates and seeds
   the operator only; it no longer populates either.

2. **The rate limit is raised for you — unless you started the servers.**
   `playwright.config.ts` passes `API_RATE_LIMIT_PER_MINUTE: '2000'` to the API it starts,
   because a run compresses a person's day into two minutes and the 100/minute default is sized
   for the person: exhausted, the API answers 429 to its own session reads and *every page renders
   "Too many requests just now"*, which reads as a broken product. So set it yourself **only** when
   `reuseExistingServer` is about to pick up an API you started — which is the common case locally.
   The sign-in cap is a different thing and is never raised: it is a constant in your configuration package, and
   `lib/suite.ts` gives each test its own `x-forwarded-for` instead.

3. **Run with the e2e environment.** Playwright's `webServer` starts the API and client
   itself (api first — the order is load-bearing: the client's readiness probe resolves the
   session, so the client cannot answer until the API is up) and inherits the shell
   environment:

   ```
   DATABASE_URL="<app-postgres>/project_e2e" API_RATE_LIMIT_PER_MINUTE=2000 your e2e test command journeys/sharing
   ```

   The selector goes on the end, and the environment is the same whether the run is one feature
   or all of them.

   The runner reads the repo `.env` itself (`process.loadEnvFile`), which is why a shell variable
   still wins: it does not override what is already in the environment. That matters because the
   `@generation` scenarios decide whether to run by reading `GENERATION_QUEUE_URL`,
   `OPENAI_API_KEY` and `AWS_ENDPOINT_URL` — a runner that cannot see them skips the product's
   central feature and the run goes green having asserted nothing.

   `reuseExistingServer` is on locally, so anything already serving `:3000`/`:3001` is used
   instead — including the docker dev stack, if it is up. That convenience is also the trap:
   servers started with the dev `DATABASE_URL` are writing to the dev database, not the e2e database,
   and the run is no longer isolated. Start the servers with the e2e environment, or stop
   them and let Playwright.

   A cold `next dev` compiles per route and can take ~35s for the first hit — longer than
   Playwright's 60s `webServer` timeout. Starting the servers yourself and letting
   `reuseExistingServer` pick them up is the reliable pattern.

4. **One project, Chromium, everywhere.** The suite has a single project — `flows`, desktop
   Chromium — so there is nothing to select with `--project`, and CI installs the same one
   engine. The four-device matrix went with `specs/`: three of its projects resolved to WebKit
   that `browsers:install` never fetched, so they failed at launch on a fresh machine with an
   error that read as a broken host, and one of the four ran nothing at all. A scenario whose
   subject is a width now says the width, and a step sets it; the touch scenarios open a context
   with `hasTouch`. Safari's rendering is the real loss and is worth a deliberate second project
   over the same features one day — never a matrix again.

5. **Read the report, not just the terminal.** `your e2e path/reports/business-flows.html` is the
   deliverable: its unit is the scenario, and the screenshot and trace hang off the step that
   failed. `pnpm --filter e2e bdd:report` serves it, which is what the trace viewer needs — over
   `file://` the steps render and the traces do not. Every failure also carries a ready-made
   prompt (`aiFix.promptAttachment`): the feature, the failing step, its source and the page's
   ARIA snapshot.

## Environment gotchas

**Stop the dev servers before `your validation command`.** The api auth specs share the real
`DATABASE_URL` with running dev servers and time out under contention
("Hook timed out in 10000ms" in `auth-under-prefix.test.ts` / "Test timed out in 5000ms" in
`auth.test.ts`). Distinct from the flake below — this one is real contention, and it goes
away when the servers are down.

**The machine-load flake.** `your validation command` intermittently fails on a *random* package with
`[vitest-pool]: Failed to start forks worker … Timeout waiting for worker to respond` — turbo's
concurrency plus vitest's forks exhaust the box after something heavy (a long `pnpm install`,
a docker build). It is not code: the failing package changes between runs and the log has no
assertion. Re-run throttled rather than chasing it:

```
your guard checks
pnpm exec turbo run lint typecheck test --concurrency=2
```

Never add retries or raise `testTimeout` for this — the tests are not slow, the worker never
started.

**Browser runs are renderer-sensitive and always headless** (a headed run rasterises text
differently). `your UI package-testing/font-setup` pins the generic font stack for the same reason:
a wrap point or a tap target that moved because a typeface loaded differently is a finding
with nothing behind it. A PNG under `__screenshots__/` is a picture of a FAILURE, written as
an attachment — it is gitignored debris, never a reference. Do not commit one.

**There is no design-truth gate any more.** Nothing compares a render against Figma, so
"matches the design" is a claim a person makes by looking, not one a test backs. When a task
turns on Figma fidelity, say that it was checked by eye — do not let a green `test:browser`
stand in for it.

## Verification

- [ ] The change's surfaces are known: UI/visual → `test:browser` + `test:e2e`; the rest → `your validation command`
- [ ] Dev servers are stopped before `your validation command`; any e2e run used a freshly seeded, isolated e2e database, on one worker, with the per-minute rate limit raised on whichever API served it
- [ ] The e2e run was **narrowed to the slice the change touches**, and the report says which slice — a full suite was either labelled onto the pull request or run to completion, never started and killed and called green
- [ ] A red flow got a verdict (`e2e`) before the environment was blamed
- [ ] A machine-load flake was re-run with `--concurrency=2`, not "fixed" with retries
- [ ] A quality failure was reported with the element and the measurement it named; any claim about Figma fidelity says it was checked by eye
- [ ] The report says what was run, what passed, and what was skipped and why — not just "tests pass"