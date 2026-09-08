---
name: e2e
description: Writing and repairing project's end-to-end flows as BDD — turning a change's scenarios into Gherkin under your e2e path/features/ with steps, page objects and fixtures behind them, and turning a red run into a verdict. Load this when a change needs an end-to-end test, when asked to author or extend a feature or a step definition, when a flow is red or flaky, when the business asks what a green run actually proves, or when a scenario has to be traced from openspec to the file that proves it. Owns the scenario-to-Gherkin translation, the five-layer split (feature, step, page object, fixture, lib) and which of the seven feature directories a flow goes in, the locator and wait conventions and the concrete failure each one defends, the tag vocabulary and what each tag costs, and the six-verdict healing procedure — as opposed to test, which owns bringing the stack up, and tdd, which owns whether a claim belongs in this mode at all. Routes to your e2e path/README.md for the inventory, your specification tool folder/schemas/integration/schema.yaml for the scenario mapping, and bugfix once a defect is confirmed.
---

# Writing and healing project's end-to-end flows

An end-to-end flow is the only test in this repository that can see a **seam** — a failure
living between two layers that each pass their own tests. It is also the only one a person
outside engineering can read. When flows are written in
**Gherkin** and run by [playwright-bdd](https://github.com/vitalets/playwright-bdd), the
document a run produces has scenarios and Given/When/Then in it rather than locators.

Those two properties are the same property. A seam is described in business terms — "a user
submits an order and receives a confirmation" — and the layers underneath are the
implementation detail that keeps failing between themselves.

This skill is the middle of three. `tdd` decides **whether** a claim belongs in this mode;
`test` brings the stack up and owns every environment gotcha. This one owns what happens in
between: how a scenario becomes a feature file, what goes in which of the five layers, and what
a red run means. The strategy is [docs/engineering/testing.md](../../../docs/engineering/testing.md)
and your project's architectural decisions; the inventory of what exists today is
[your e2e path/README.md](../../../your e2e path/README.md).

## What this skill does not own

Read this table before adding a rule here. A rule that belongs to one of these lives there.

| Question | Owner |
|---|---|
| Bringing the stack up — the isolated e2e database, the rate limit, `reuseExistingServer` — and how much of the suite one run selects | `test`; the selectors are your e2e README's *Running a slice* |
| Whether this claim needs `e2e` at all, and the RED→GREEN order | `tdd` |
| The four-viewport matrix the component suites measure against, and why CI runs the jsdom layer alone | [docs/engineering/testing.md](../../../docs/engineering/testing.md) |
| Which of the seven directories under `features/` a new flow goes in, and the five rules the split rests on | [your e2e path/AGENTS.md](../../../your e2e path/AGENTS.md) |
| What each feature contains, the tag vocabulary's current membership, and the three things a shared-resource flow deliberately does not prove | [your e2e path/README.md](../../../your e2e path/README.md) |
| Everything after a defect is confirmed — diagnosis, change folder, fix | `bugfix` |

## The five layers, and what may go in each

The whole design is one rule: **a reader of a feature file must never need to open anything
else.** Everything that would break that promise is pushed down a layer.

| Layer | Holds | The test for "does this belong here" |
|---|---|---|
| `features/<split>/*.feature` | Gherkin. Business language, business outcomes. The directory is chosen by what a failure MEANS — [your e2e path/AGENTS.md](../../../your e2e path/AGENTS.md) asks the question in order | Would somebody outside engineering read this sentence and agree it is what the product must do? |
| `steps/*.steps.ts` | The Given/When/Then implementations, plus assertions | Is this the *meaning* of the sentence, in code? |
| `pom/*.page.ts` | Locators, waits, and the REASONING behind each one | Is this "how you reach the thing", rather than "what must be true of it"? |
| `fixtures/index.ts` | The page objects and the per-scenario `world` | Does more than one step need it? |
| `lib/*.ts` | Identities, the API origin, the anonymous context, the outbox, the viewport matrix | Does more than one LAYER need it? |

Four things follow from that split and none of them are negotiable:

- **A step never constructs a page object.** It takes one as a fixture. The moment a step says
  `new WizardPage(page)` it is code again, and the second step needing the same object has to
  build a second one.
- **A page object never asserts a business outcome.** It may assert its own preconditions —
  "the create did not land on a project URL" — because that is a precondition, not a verdict.
  A page object that quietly asserted the wall's contents would make every scenario calling it
  claim something the feature file never said.
- **A feature file names no selector, no URL, no status code and no id.** `Then they are
  answered as though it does not exist` is the claim; that it is a `404` and not a redirect is
  the step's business, and the step says why.
- **State moves between steps through `world`, never through a module-level variable.** A Given
  writes it, a Then reads it, and every read goes through a guard that fails with "this scenario
  has no project — a step that acts on one must be preceded by a Given that creates or finds one"
  rather than `undefined is not a string`.

## From a change's scenarios to a feature file

A flow comes from a **scenario**, never from a screen — and it comes **first**. Every typed change
writes `flows.md` immediately after its proposal, before the interface, the specs and the tasks,
and task 1.1 is the flow: written, run, and its outcome recorded before the layer is built
(`your specification tool folder/schemas/*/schema.yaml`, your project's architectural decisions). A flow authored after
the code it checks proves that the code passes it, not that it would have caught the code being
wrong.

**A `ui` or `be` change cannot turn its flow green** — no route, no screen — so it writes it into
`known-issues/` tagged `@known-issue` with the blocker named in a sentence. That is not a
deferral: `bddgen` still refuses a sentence with no step behind it, so the vocabulary and its steps
exist at the moment the layer is built, and the `integration` change removes a tag and moves a file
rather than inventing a journey from three changes' worth of assumptions. `known-issues/` remains
closed to red tests: a scenario that fails because the application is *wrong* is a defect.

The `integration` schema
([your specification tool folder/schemas/integration/schema.yaml](../../../your specification tool folder/schemas/integration/schema.yaml))
owns the mapping and its own test of whether a requirement belongs there at all. Within one
file, the correspondence is now almost literal:

| In the change | In the feature |
|---|---|
| `#### Scenario:` | exactly one `Scenario:` |
| `### Requirement:` | a `Rule:` — Gherkin has the keyword, and it is the right one |
| the requirement's `*Test:*` marker | the `.feature` it names, which must exist |
| the requirement's `*Flow:*` marker | the `Scenario:` in `flows.md` it traces to — or a line saying the claim is layer-internal |
| WHEN / THEN | `When` / `Then`, in the requirement's own words where they survive translation |

Three rules the schema does not state:

- **A journey is one scenario across several steps**, not four scenarios that each sign in and
  set themselves up. Use `Background:` for the setup every scenario in a feature shares, and a
  `Scenario Outline:` with an `Examples:` table where the same sentence is genuinely true of
  several values — a width, a reason a sign-in was refused, an off-site destination. The table is
  the enumeration a reader can check; four near-identical scenarios are one they have to diff.
- **Assert the seam, not the unit.** If the assertion would hold with the other process stubbed,
  it belongs one layer down and `tdd` says where.
- **Name the requirement in a Gherkin comment above the scenario** — `US-AUTH-02`, `FR-2.4-03`.
  It is what survives the change folder being archived.

The step-by-step path from a change folder to a finished feature, with an annotated skeleton, is
[references/authoring.md](references/authoring.md).

## The tags, and why they are not decoration

`steps/hooks.ts` is where a tag becomes behaviour. That indirection is the point: **the cost of a
scenario has to be visible in the scenario**, not buried in a condition inside a step.

| Tag | What the hook does | Why the tag rather than an `if` |
|---|---|---|
| `@journey` `@roles` `@internal` `@cross-cutting` | 120s budget | 15s is sized for a warm route; a journey walks several, each compiling on first hit |
| `@external` | Extended budget; skips with a reason when required external infrastructure is absent | A scenario waiting on a worker, webhook, queue, or third-party callback otherwise reads as a product hang |
| `@spends` | **Excluded from `your e2e test command`** by its own `--grep-invert`; `your e2e test command:spends` is what runs it | It calls an expensive external service, and a reader deserves to know |
| `@billing` | Nothing — a label | Marks a scenario asserting something must NOT spend. These cost nothing and are the only guard against a regression that bills silently |
| `@optional-infra` | **Excluded from `your e2e test command`** by the same `--grep-invert`; a named script runs it when the infrastructure is configured | Applied to a whole feature area that depends on infrastructure the default run does not provide; the README owns the reversal condition |
| `@known-issue` | Skips, with the blocker named | The behaviour is not built. See below |
| `@mode:default` | Nothing in `hooks.ts` — it is **playwright-bdd's own** tag, and it makes the feature's scenarios run one at a time without skipping the ones after a failure | A feature whose scenarios take the catalogue out from under each other cannot run beside itself. `serial` would stop the business report at the first red line, which is the one report nobody can act on |
| `@layout` | Nothing — a label | Marks the scenarios whose subject IS a width, so a run can be narrowed to them. The width itself is in the sentence |

**There is no viewport tag, and there is no need for one.** The suite runs in a single
Playwright project — `flows`, desktop Chromium — and there is no device matrix left to repeat it
across; `specs/` is gone and every claim it carried is a feature. A scenario therefore cannot be
repeated across four viewports, so nothing has to opt out of that: the guarantee lives in
`playwright.config.ts`'s project layout rather than in a tag every expensive scenario has to
remember to carry.

**A claim that IS about a width says which width, in its own sentence.** `Given the screen is
390 pixels wide`, `And a stranger opens the link on a tablet held sideways` — a step sets it, and
the failure names the width instead of reporting under a project name that did not. Where the
claim genuinely is the matrix, that is a **`Scenario Outline` with an `Examples` table**:
`cross-cutting/sign-in-layout.feature` lists the four viewports as rows, so the widths are visible
to the reader and only the scenarios that are about width pay for them — which is precisely what a
device project could not do. A shared-resource feature can open a second context with
`hasTouch` for the scenarios whose subject is the device, so a touch claim needs no second
engine either. Geometry against a committed baseline is still the pixel and visual modes' —
`tdd` says which.

**`@known-issue` is not a parking space for red tests.** A scenario goes there only when the
thing it asserts has not been built and the blocker fits in a sentence. A scenario that fails
because the application is *wrong* is a defect and belongs to `bugfix`. And it is skipped
from the HOOK rather than with Gherkin's own `@skip`, because a feature whose scenarios are all
`@skip` generates no spec at all — a gap that appears nowhere in the report is a gap nobody is
reminded of.

## The conventions, and the failure each defends

Every one of these was paid for. The right-hand column is the point: a rule whose failure you
cannot name is a rule that gets relaxed the first time it is inconvenient.

| Rule | The failure it defends |
|---|---|
| `test` and `expect` come from a shared fixtures file that merges playwright-bdd's instance with project-specific test helpers | The merged fixture carries per-test rate-limit bypass (e.g. an IP override header). Without it, sign-in-heavy suites hit the API's rate limit partway through and fail with 429. It is typically an `auto` fixture, so nothing declares a dependency on it and **nothing fails loudly when it is dropped** — spreading a test object instead of merging it loses it silently |
| Navigate **relative**: `goto('sign-in')` | Playwright resolves with the WHATWG URL API, so `goto('/sign-in')` discards the base path entirely and your web framework 404s. `playwright.config.ts` explains the load-bearing trailing slash |
| Take `basePath` from the Playwright config; never write a prefix | The prefix is deployment configuration, so a hardcoded value can pass locally and fail only behind the real path. Keep the allow-list in your URL-prefix guard narrow |
| Locate by **role**, then label, then text the component owns. A structural selector is the last rung and carries a comment saying why no role exists | A test id survives a layout being rebuilt wrongly, and noticing exactly that is the point. The three structural selectors in `pom/` each name why: a visually-hidden file input, a link whose only name is an opaque id, a table cell addressed by position |
| **Scope a role inside its region** — `getByRole('form').getByRole('alert')` | Some frameworks inject route announcers or global live regions into every page. A bare alert role can match the framework element, so visibility assertions pass on pages showing no message or fail for the wrong reason |
| Pair a submit with the response it causes: `Promise.all([waitForResponse, click])` | The Server Action is what carries the API's `Set-Cookie` onto a response the browser receives. Returning at the click races an in-flight submission. Polling assertions hide the race; anything reading state once, like `context.cookies()`, does not |
| Type through rich text editors the way a user does | Many contenteditable editors update internal state only through input transactions. Directly setting text can leave the editor's document empty, so submit actions wait on a request that was never sent |
| **Clear a file input before re-setting it** | Setting the same path onto an input that already holds it fires no `change` event, so a retry is a no-op that spins over the same silence |
| **After a reload, an interaction may need a retry** | `page.reload()` resolves at `load`; React hydrates later. A `change` event fired in between reaches no listener, the panel sits in its empty state, and the failure reads as "the upload did not complete" against a stack that never received a byte |
| Create what you assert on, with a **run-unique name**; never assert a count or an empty table | There is no per-worker database isolation yet. A suite that assumed an empty portal passes once and fails for the rest of the day |
| **Never name a seeded row.** Ask the API for one matching the property you need, and page through | The suite's own creation scenarios can push seeded records behind pagination during a run. A named lookup then fails together across suites, while an API query for the needed property stays tied to the scenario's intent |
| Bind a seeded credential to a name; never inline the literal | `password: '<literal>'` is what the `secret-literals` guard looks for, and teaching it to skip test files would blind it to a real one landing in one |

## Identity, environment, and where a file lands

**A stranger needs a second context.** `lib/anonymous.ts` — `openAnonymously()` and
`expectNoSession()` — exists because reusing the signed-in context proves only that a signed-in
person can open a URL, and would pass against an application with no unauthenticated path at
all. The caller closes the context.

**A missing environment variable is a skip with a reason, never a throw at module load.**
Playwright imports every generated spec before running any of them, so a module-scope throw in a
step file is a COLLECTION failure: it takes the whole run down — every other feature with it — and
reports one error about a variable nothing else needs. Treat optional integrations as
`known-issues/` or skip with a reason, and read the variable at the
moment it is used. A step file is imported by the generator as well as the runner, so the rule
applies to `steps/` and `lib/` alike.

**The runner reads the same `.env` the servers do** (`process.loadEnvFile` in
`playwright.config.ts`), because a runner that cannot see the external service env vars skips every
scenario that would have proved the product's central feature, and the run goes green having
asserted nothing. A skip nobody asked for is worse than a failure.

**A feature file's directory is a contract, not a filing preference.** It is what a reader of a
red run is told first — `journeys/` means the product does not work, `internal/` means an internal
workflow is broken and no public user is affected, `known-issues/` means nothing at all. Putting a
flow in the wrong one misreports the failure before anybody opens it;
[your e2e path/AGENTS.md](../../../your e2e path/AGENTS.md) asks the placement question in order.

## Healing a failing flow

**A verdict is not a diagnosis.** Reach one before editing anything, and record the observation
that produced it. Five verdicts resolve here or route sideways; the sixth is a hand-off, and
performing it yourself does `bugfix`'s work out of order.

| Verdict | The observation that distinguishes it |
|---|---|
| **Stale fixture** | The app works by hand, and the value the scenario names is absent from the seed |
| **Stale locator** | The element is present under a different query — a different role (`alertdialog`, not `dialog`), a different element (`button`, not `link`), an `aria-label` rather than visible text, or a strict-mode violation naming two hits |
| **Stale specification** | Nothing moved and no element changed: the *journey* was overtaken by a decision record. A wait for one route or query state outlived the product landing somewhere else, and spent its whole budget on a run that had already succeeded |
| **A race** | Intermittent, or one-sided: hydration, a `router.replace` that has not landed, a redirect still in flight. `--repeat-each=5` fails some runs, not all |
| **Environment** | Every scenario in the file fails identically, or generation fails outright — `test` owns it |
| **A defect** | Everything else is excluded and the application is genuinely wrong — hand off to `bugfix` |

**Heal on one feature, not on the suite.** `your e2e test command --last-failed` from `your e2e path/` is the first
loop, `your e2e test command journeys/sharing` is the second, and `--repeat-each=5 --grep "<scenario>"` is
what separates a race from a defect. A verdict reached in ninety seconds is reached; a forty-minute
full run between edits is how a red flow stays red for a day. The suite comes back at the end, once
— [your e2e path/README.md](../../../your e2e path/README.md)'s *Running a slice*.

**A red run hands you the prompt.** `aiFix.promptAttachment` attaches the feature, the failing
step, that step's own source and the page's ARIA snapshot to every failure. Read it before
reaching for the trace; it is usually enough to reach a verdict, and the trace is what
distinguishes a race from a defect once you have one.

The triage, the commands that separate a race from a defect, and three worked cases from this
repository are [references/healing.md](references/healing.md).

## What may never be done to make a flow pass

Each of these turns a red test into a green one that proves less than it did:

- Weakening a Gherkin sentence so it describes what the product does rather than what it must do
- Unscoping a role back to a bare one, or dropping to a class selector, to make a locator match
- Deleting an assertion rather than the whole scenario — a scenario missing its point still looks
  like coverage, and now it looks like coverage to the business too
- Moving a failing scenario to `@known-issue` when the application is wrong rather than unbuilt
- Adding `retries` locally, or raising a timeout, to cover something intermittent
- Adding `waitForTimeout` in place of waiting for the thing that actually happens
- Asserting copy the client owns rather than the behaviour the requirement asks for
- Asserting the output of a stubbed port, which asserts the stub

The test for all eight: **would this still fail against an application with the feature
removed?** If not, it has stopped being a flow.

## Verification

- [ ] Every `#### Scenario:` in the change has a `Scenario:`, and every `*Test:*` marker resolves to a feature that exists
- [ ] Every sentence in the feature reads as a business outcome, and names no selector, URL, status code or id
- [ ] The feature is in the directory that matches what its failure MEANS, and answers `your e2e path/AGENTS.md`'s placement question only once
- [ ] `npx bddgen` passes — a sentence with no step behind it fails generation, and that is the check
- [ ] Every locator is a role or a label, scoped inside its own region — or the escape hatch is commented with the reason no role exists
- [ ] Every submit is paired with the response it causes; no bare click after a Server Action
- [ ] Everything asserted on was created by this scenario under a unique name, or found by asking the API; no assertion depends on the table's prior contents
- [ ] Every scenario that costs money carries `@spends`, and every one needing optional external infrastructure carries `@external` or the project's equivalent
- [ ] A feature that depends on optional infrastructure carries its feature-level tag, and the author knows whether the default run covers it
- [ ] A claim about a width says the width in the scenario, and a step sets it — no scenario is waiting for a project to supply one
- [ ] Every non-obvious line names the concrete failure it defends, and ambiguous assertions carry a diagnosis string
- [ ] A red flow got a verdict, with the observation that produced it, before it got an edit
- [ ] No assertion was weakened to reach green
