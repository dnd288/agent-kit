# Healing a red flow

`SKILL.md` holds the six verdicts. This is how to reach one, what separates them, and three
worked cases from this repository.

**Reach a verdict before editing anything.** A red flow that gets an edit before a diagnosis
becomes a green flow that proves less than it did, and nobody can tell afterwards.

## 0. Read the prompt the run already wrote for you

Every failure carries an AI-fix attachment (`aiFix.promptAttachment`): the feature, the failing
step, that step's own source, and the page's ARIA snapshot. It is in the terminal output and in
`reports/business-flows.html` against the step that failed.

Read it first. It is usually enough to reach a verdict, and it is always enough to rule out
"the element moved" — the snapshot shows what was actually on the screen.

## 1. Separate the environment from the product, before anything else

If **every** scenario in a file fails identically, or the whole run reports one error and zero
tests, it is the environment and `acme-test` owns it. Three shapes, all seen here:

- **A generation failure, before any browser starts.** `bddgen` runs first, and
  `missingSteps: 'fail-on-gen'` means a sentence with no step behind it fails there: the run reports
  the feature, the line and the sentence, and zero tests. That is not a red flow — it is an
  unfinished one, and the fix is a step definition, never a rewritten sentence.
- **A module-scope throw in a step file.** Playwright imports every generated spec before running
  any, so a throw while loading takes the whole run down and reports one error about a variable
  nothing else needs. The withdrawn `user-sync` spec did this for `RAILS_SYNC_HMAC_SECRET` — empty
  by default locally — so the suite could not run at all and the reason looked like a external system
  problem. Its claims are `known-issues/user-mirror.feature` now, and the variable is read at the
  step that uses it.
- **The wrong database.** The long-lived `pnpm --filter api dev` in this workspace is usually
  pointed at the e2e database, while the repo `.env` says the dev database. Seeding the wrong one
  produces a full run of failures that all look like product defects. Read it off the process:
  `ps eww $(pgrep -f 'filter api dev' | tail -1) | tr ' ' '\n' | grep ^DATABASE_URL=`.
- **The provider refusing us.** A run with no credit left sits on the generating screen exactly
  as a healthy one does, and the scenario spends its whole budget before reporting a timeout. The
  room's `generationContext` carries the reason —
  `{"analyseError":{"code":"provider_no_credits"}}` — and `staging.feature` asserts the screen's
  own message early so the failure names the environment rather than the product.

## 2. Then work the five remaining verdicts

| Verdict | How to confirm it | What to do |
|---|---|---|
| **Stale fixture** | The app works by hand, and the value the scenario names is absent from the seed | Reseed. If the scenario NAMED a seeded row, that is the real fault — ask the API instead |
| **Stale locator** | The ARIA snapshot shows the element under a different query | Fix the page object, and comment why the new query is the right one |
| **Stale specification** | Nothing moved, no element changed: a decision record overtook the journey | Rewrite the scenario against the current requirement, and say in the feature what changed |
| **A race** | `--repeat-each=5` fails some runs, not all — or fails only after a reload | Wait for the thing that actually happens. Never `waitForTimeout` |
| **A defect** | Everything above is excluded | Hand off to `acme-bugfix`. Do not diagnose it here |

The command that separates a race from a defect:

```
cd e2e && npx bddgen && npx playwright test --repeat-each=5 -g "<scenario name>"
```

`-g` matches the generated title, which is the `Scenario:` sentence — so the scenario is addressed
by the words the business wrote. There is one project and one worker locally, so neither needs a
flag; adding `--workers=4` to "speed up a repeat" reintroduces exactly the contention the repeat
is trying to rule out.

Five green is a race that has been fixed. Five red is not a race.

## Worked case 1 — a stale locator wearing four disguises

Six scenarios failed across the products and moodboards consoles, each with a different message.
All four were the same class:

- `getByText('What is missing from the catalogue')` found nothing — the sentence is the coverage
  strip's `aria-label`, so it is a name a screen reader announces and never text on the page.
  `getByRole('region', { name: … })` is the query.
- `getByRole('dialog')` found nothing — `AdminConfirmDialog` is destructive and carries
  `role="alertdialog"`, which `getByRole('dialog')` does not match.
- `getByRole('link', { name: 'your web framework page' })` found nothing — the pager renders a `<button>`.
- `selectOption` threw "Element is not a `<select>`" — the facets are Base UI Selects, so they are
  a `<button role="combobox">` that is opened, and an option that is clicked.

**The tell for all four**: the element is plainly on the screen in the ARIA snapshot, under a
different role or a different name. That is a stale locator every time, and it is a page-object
fix with a comment saying what the element actually is.

## Worked case 2 — a race that only happened after a reload

`When they upload a photograph of the room` failed with "the photo never became ready", but only
in the one scenario that reloads the page first. The trace showed the presign Server Action
running and **no PUT to S3 following** — and then, on closer reading, no presign either: the panel
was in its empty drop state, showing no preview, no progress and no error.

`setInputFiles` dispatches a `change` event. React hears it only once the page has hydrated, and
`page.reload()` resolves at `load`, which is earlier. The event reached no listener.

Two things made the fix correct rather than approximate. It is a **retry**, not a sleep — a fixed
wait is a guess that grows every time somebody watches it flake. And the input is **cleared
before each retry**, because setting the same path onto an input that already holds it fires no
change event, so the first version of the fix span three times over the same silence and failed
identically.

## Worked case 3 — a stale specification that cost seven minutes a run

`await page.waitForURL(/\?step=RESULT/i, { timeout: 120_000 })` timed out on a generation that had
succeeded. Nothing had moved and no element had changed: the wizard stopped landing on
`?step=result` when success started taking the agent to the project's own record at `/details`.

The tell is the shape of the failure — a wait that spends its entire budget while the application
does the right thing. It had also been invisible for weeks, because the scenario skipped on every
machine that could not see the generation stack's environment variables, and a skip nobody asked
for reads as a pass.

Two fixes, not one: the wait was corrected, and `playwright.config.ts` now loads the repo `.env`
so the runner and the servers agree about what is configured.

## What healing may never do

`SKILL.md` has the full list. The three that get reached for under time pressure:

- Moving a failing scenario to `@known-issue`. That directory is for behaviour that has NOT BEEN
  BUILT, with the blocker named in a sentence. A scenario that fails because the application is
  wrong is a defect.
- Weakening a Gherkin sentence until it describes what the product does rather than what it must
  do. This is worse here than in a plain spec, because the business reads these.
- Deleting an assertion rather than the whole scenario. A scenario missing its point still looks
  like coverage — and now it looks like coverage in a report somebody outside engineering is
  reading.
