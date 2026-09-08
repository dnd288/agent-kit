# Authoring a flow, from a change folder to a feature file

The step-by-step path. `SKILL.md` holds the conventions themselves; this is the order to apply
them in, worked against flows that are in the repository today —
`your e2e path/features/journeys/` contains the user journey features.

## 1. Read the change for flow-shaped scenarios

Open the change's `specs/` and look for scenarios whose failure would live BETWEEN two layers.
A scenario belongs here when the sentence you would write to describe the failure names two
processes. "The wizard step is persisted" is one layer. "The agent closes the tab and comes back
to the step they were on" is two, and it is the second that becomes a `Scenario:`.

Everything else goes one layer down, and `tdd` says where.

## 2. Place the file

`your e2e path/AGENTS.md` asks the placement question in order. The short version: does a real person walk
it end to end → `journeys/`; is it about what a ROLE may reach → `roles/`; staff-only surface →
`internal/`; every screen → `cross-cutting/`; a fixed defect → `defects/`; not built yet →
`known-issues/`; otherwise → `capabilities/`.

## 3. Write the Gherkin FIRST, and write it for somebody who will never open the code

This is the step people skip, and skipping it is how a suite ends up with sentences like
`When the user clicks the button with text "Generate design"`. Three tests:

- **Would a person outside engineering agree this is what the product must do?** If the sentence
  is about a mechanism rather than an outcome, it is a step-definition comment, not a step.
- **Does it name a selector, a URL, a status code or an id?** Then it is in the wrong layer.
  `Then they are answered as though it does not exist` is the claim; that it is a `404` and not a
  redirect is the step's business, and the step says why.
- **Would the scenario still fail if the feature were removed?** If not, it is not a flow.

An annotated skeleton:

```gherkin
@journey
Feature: Staging a room

  # The prose under Feature: is where the WHY goes — what only a browser can see here, what the
  # scenarios cost, and what they deliberately do not prove. It is read by everybody and costs
  # nothing to run.

  Background:
    # Only what EVERY scenario in the file needs. A Background that sets up three quarters of
    # them is three quarters of a file in the wrong place.
    Given 'a logged-in user' is signed in

  # A `Rule:` is the change's `### Requirement:`. Gherkin has the keyword; use it.
  Rule: A job is started when the user submits the form, and never by a page load

    # A comment above a scenario carries the requirement id — US-PRJ-06, FR-2.2-D1 — and the
    # failure the scenario exists for. It is what survives the change folder being archived.
    @billing @external @spends
    Scenario: Landing on the processing screen does not start a second run
      Given they have completed the prerequisite steps
      When they submit the form and a job is queued
      And the job is being processed
      And they land on the processing screen again
      Then no further job was started
```

## 4. Tag it honestly

The tag is how a reader learns what the scenario costs, so it is part of the sentence rather than
an annotation on it. `@spends` if it calls an expensive external service; `@external` if it needs optional infrastructure; `@billing` if it asserts
something must NOT spend. `@mode:default` is playwright-bdd's own, and belongs on a feature whose
scenarios take the data out from under each other. There is no viewport tag: the suite runs in one
project, so a scenario is never repeated across a matrix — and a claim that IS about a width says
the width in its own sentence, with a step that sets it.

**A `@billing` scenario usually needs `@spends` too**, and that is not a contradiction. Proving
"landing did not start a SECOND run" needs a first run to exist: without one, the pointer is
empty before and after, they match, and the scenario reports that landing spent nothing — which
is true and proves nothing. It would go green against a product that could not process a job at
all. Assert the precondition in the Given.

## 5. Generate, and let the generator tell you what is missing

```
cd e2e && npx bddgen
```

`missingSteps: 'fail-on-gen'` means a sentence with no step behind it fails here, loudly, before
anything runs. That is the check: a business-readable suite that can quietly skip a scenario is
one that starts describing behaviour nobody implemented.

## 6. Write the step, then push everything that is not the meaning downwards

A step is the sentence's MEANING in code, and its assertions. Everything else moves:

- the locator, the wait, and the reasoning behind either → `pom/`
- anything two steps need → `fixtures/`
- anything two layers need → `lib/`

```ts
Then('no further job was started', async ({ result, world }) => {
  const after = await result.currentJobId(requireResourceId(world.resourceId));
  expect(
    after,
    // The diagnosis string is not decoration. A failure here has to explain, to somebody who did
    // not write it, why the number mattered — and what it costs when it is wrong.
    'landing on the processing screen started a NEW job. A refresh, a typed URL or a restored tab is not a request to process, and each one costs real money.',
  ).toBe(world.jobIdBefore);
});
```

Two things a step may never do: construct a page object (take it as a fixture), and read `world`
without a guard. Every read goes through `requireResourceId`-style guards so a scenario missing a
Given fails with "this scenario has no resource" rather than `undefined is not a string`.

## 7. Put the reasoning in the page object

The page object is the only layer where a comment about a locator belongs, and the repository's
standard is that a non-obvious line names the failure it defends:

```ts
/**
 * Upload the file and continue.
 *
 * `setInputFiles` dispatches a `change` event. React only hears it once the page has HYDRATED,
 * and `page.reload()` resolves at `load` — which is earlier. So on a freshly reloaded upload step
 * the file is set, the event fires into a page with no listener, and the panel sits in its empty
 * drop state. The failure reads as "the upload did not complete" against a stack that never
 * received a byte.
 *
 * Retried, and the input is CLEARED before each retry: setting the same path onto an input that
 * already holds it fires no change event, so a bare second attempt is a no-op.
 */
```

## 8. Run it, and read the report rather than the terminal

```
cd e2e && your e2e test command journeys/sharing
```

**Run the feature you just wrote, not the suite.** A full local run is ~183 scenarios on one
worker and tens of minutes, and the one on record was killed at forty having reported nothing;
`your e2e path/README.md`'s *Running a slice* has the selectors — a directory, a feature, `--grep`,
`--last-failed`. `your e2e test command` already excludes the scenarios that call a model; the whole
suite, billed scenarios included, is `your e2e test command:spends` — and it is cheaper, not shorter.

`test:e2e` is `bddgen && playwright test --grep-invert "@spends"`: the generation is part of the
run whatever is selected,
so a sentence with no step behind it fails before a browser starts and narrowing never narrows that
check. There is one project (`flows`) and one worker locally, so neither takes a flag.

Then read the report — `pnpm --filter e2e bdd:report` serves it, which is what the embedded trace
viewer needs; over `file://` the steps render and the traces do not. The unit is the scenario and
the lines are the steps, so this is where you check the thing the terminal cannot tell you: **does
the scenario read as a claim about the business?** If the steps read as a script, the Gherkin needs
rewriting, not the code.

## 9. Check it against the verification list

`SKILL.md`'s final section. The two that catch the most: every sentence reads as a business
outcome and names no selector; and everything asserted on was created by this scenario or found
by asking the API, never named from the seed.
