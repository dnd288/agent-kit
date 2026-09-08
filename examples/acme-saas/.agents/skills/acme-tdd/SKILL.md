---
name: tdd
description: "Test-first development for project tasks, and which of the four test modes applies to which claim. Load this when implementing any logic, fixing any bug, or changing any behaviour — a component, a service, a route, a migration — when a task from your task plan says a scenario needs coverage, or when asked what a test should look like for a given layer. Owns the RED→GREEN cycle, the <name>.<kind>.tsx naming that five globs key on, why a story is a fixture and not a test, and why geometry claims belong in Chromium, not jsdom. Routes to your testing documentation and your project architectural decisions for the strategy."
---

# Test-first development for project tasks

*Later* is the load-bearing word. A task is not done when it works; it is done when the claim
its scenario makes is defended by a test that fails without it. The strategy — four modes,
stories as fixtures, pixel baselines on composites only — is your project architectural decisions
and your testing documentation. This skill owns the discipline and the mode mapping.

## Source

Adapted from `test-driven-development` in
[`addyosmani/agent-skills`](https://github.com/addyosmani/agent-skills) (MIT), rewritten
around this repository's four test modes and its naming contract.

## The naming contract first, because everything keys on it

`<name>.<kind>.tsx` — never `index.tsx`. Five config globs and guard regexes key on the
suffix, and three of them would **silently weaken** rather than error if the file were named
otherwise (your project architectural decisions). The kinds:

| Suffix | Mode | What it asserts |
|---|---|---|
| `<name>.test.tsx` | unit (jsdom) | structure, logic, the i18n contract — no layout, no cascade |
| `<name>.stories.tsx` | stories (Chromium) | renders and passes axe; **a fixture, not a test** |
| `<name>.quality.test.tsx` | quality (Chromium) | geometry, computed tokens, states, render counts — absolute checks, valid at any viewport |

`*.quality.test.tsx` holds absolute assertions that name the element and the measurement; a
failure there is a defect, not a drift. That absoluteness is why it is the last browser layer
standing: two comparative ones, `pixel` and `visual`, are retired (your project architectural decisions, your project architectural decisions),
because a check that needs a prior render to compare against cannot answer on a fresh clone.
**Nothing verifies design source fidelity mechanically now** — when a claim rests on it, check by eye
and say so. A component with no stories is never rendered in a browser and axe never sees it
— the `component-has-stories` guard enforces that. A story is what the browser renders; it is
not evidence of behaviour.

## Which mode answers which claim

| Claim | Mode that can observe it | Why the others cannot |
|---|---|---|
| A control is >=44px, visible, unobscured | `quality` | jsdom measures zero; a story only renders |
| Every computed colour is a theme token | `quality` (`expectTokenisedStyles`) | the class-name guard cannot see inline styles or Base UI defaults |
| Copy survives pseudo-localisation | `unit` under `renderPseudo` (`expectNoHardcodedCopy`) | only the pseudo-locale brackets catalogue strings |
| The layout matches the design source frame | `visual` (the design-truth gate) | `pixel` only proves *unchanged*; the visual layer compares against the tracked design source snapshot (your project architectural decisions) |
| Focus is visible, states differ non-colourfully | `quality` (`expectFocusVisible`, `expectNonColourCue`) | `:focus-visible` needs real keyboard modality |
| A route strips what its schema does not name | `api` unit via `app.inject()` | nothing else exercises serialisation |
| A service returns plain data | `api` unit | the test constructs the plain arguments |
| A session is revocable — a deleted row ends it | `api` unit against the store, **and** `acme-e2e` | the only test that distinguishes a database session from a JWT wearing its name (your project architectural decisions) |
| A refusal reveals nothing — unknown address vs wrong password | `api` unit comparing **status and bytes**, plus that the expensive path ran | a body-only assertion misses the timing oracle |
| A cookie keeps `httpOnly`/`SameSite`/`Path` across the proxy hop | `api` unit on the parse, **and** `acme-e2e` on the browser's jar | each side is correct alone; only the browser proves the hop |
| A route works under the configured URL prefix | `api` unit with the prefix **stubbed before import** | `config` is read at module load, so a normal test runs at `/` and passes vacuously |
| A mutation invalidates what it changed | `acme-e2e` | a unit test sees the action return, not the next render |
| A wizard step survives a reload or a shared link | `acme-e2e`, navigating **straight to the URL** | if it is not in the URL there is nothing to navigate to — which is the finding |
| A slice's transition is correct | `unit`, calling the slice as a plain function | no React needed; if it needs React it is not a slice |

The mapping rule: **assert the property where it can be measured, and name that property in
the test.** `assertHasLayout` throws in jsdom rather than passing vacuously — a test that
cannot measure should fail, not pass.

**An `acme-e2e` claim is written as a sentence, not as a file.** The flows are Gherkin under
`e2e/features/`, so the claim goes in a `Scenario:` and its meaning goes in a step. RED there is a
scenario whose steps exist and whose assertion fails; a sentence with NO step behind it fails
`bddgen` before a browser starts, which is an unfinished test rather than a red one and proves
nothing either way. `acme-e2e` owns the translation and the layer it lands in.

## The RED→GREEN cycle

1. **Write the failing test first** — the scenario's claim, in the mode that can observe it.
   It fails for the right reason: the behaviour does not exist. This is the doubt step of
   `acme-verify` made concrete; a failing test is a disproof attempt.
2. **Implement the minimum** that turns it green — in the layer the surface table names, with
   the layer's skill loaded (`acme-component-development`, `acme-ui-development`,
   `acme-api-contract`).
3. **Watch it pass for the right reason.** A test that passes because `assertHasLayout`
   threw, or because the component is never rendered, has not passed — it has not run.
4. **Ship the surfaces the table demands**, not just the one test. A primitive arrives with
   unit + stories + quality; a composite adds pseudo + pixel.

The stories are written to enumerate **variants and states**, not to look pretty in
Storybook: *a control the a11y run never toggles is a variant never checked.* Enumerate
every variant in a story rather than leaving them to controls.

## The first-run pixel baseline

The first run of a pixel test **writes the baseline and fails on purpose**. Open the PNG
before keeping it — an unreviewed baseline records whatever was on screen, and the first
baselines written in this repository recorded entirely unstyled components. Never delete or
regenerate a baseline to make a test pass; that is how a wrong render becomes the reference.
The **visual layer is the design-truth gate** (your project architectural decisions): it compares the rendered DOM
against the tracked design source snapshot and can fail because the code disagrees with **the
design** rather than with its own history. A failure names the node, the delta and the
fix — report it verbatim, and never edit a snapshot to make it pass; re-run
`your design sync command` only when the design actually moved.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "I'll write the test after the code works" | After is the load-bearing word. There is no after. |
| "The story is my test" | A story renders and runs axe. It asserts nothing about behaviour. |
| "It's just a unit, jsdom is fine for layout" | jsdom has no layout. `assertHasLayout` throws there — the harness is telling you the mode is wrong. |
| "The baseline already exists" | Then it asserts *unchanged*, not *correct*. The design source diff is the correctness check. |
| "The guard requires stories, so stories are the requirement" | Stories are the fixture the browser modes run against. The requirement is the behaviour. |

## Red flags

- A test file that cannot fail (asserts `true`, or measures zero and passes)
- A baseline regenerated to make a test green without being opened
- A claim about geometry defended only by a jsdom test
- A component with stories but no quality test
- A failing test whose failure message names no element and no measurement

## Verification

- [ ] Every scenario's claim maps to a test that fails without the implementation
- [ ] The mode matches the claim — geometry in Chromium, copy under the pseudo-locale, serialisation via `app.inject()`
- [ ] Stories enumerate every variant and state the design specifies
- [ ] Every new pixel baseline was opened and looked at
- [ ] your validation command is green — it is jsdom only and needs no browser
- [ ] **If the change alters what a component looks like or how a screen behaves**, the deep
      layers were run too: your browser tests and the e2e slice it touches (`cd e2e && pnpm
      test:e2e journeys/sharing` — the whole suite is CI's, `e2e/README.md` § *Running a slice*).
      They are not in CI by default (your project architectural decisions,
      amended), so nothing else will run them for you — the your CI label for browser checks, your CI label for e2e checks,
      or your CI label for deep checks label on a pull request is the CI route
