/**
 * What a tag in a feature file actually does.
 *
 * ── Why tags rather than `if` inside a step ───────────────────────────────────────────────
 *
 * A feature file is read by people who will never open a step definition, so the cost of a
 * scenario has to be visible in the scenario. `@spends` says this one calls an expensive
 * external service and is billed; `@generation` says it needs specific infrastructure. What a
 * scenario costs is therefore readable in the scenario.
 *
 * The hooks below are where those words become behaviour. Putting the condition inside a step
 * instead would hide the cost in code and leave the feature file lying about what it does.
 *
 * ── A skip must say why ───────────────────────────────────────────────────────────────────
 *
 * Every skip here carries the reason and the fix. A check that goes quietly amber is one people
 * learn to ignore, and the worst outcome available to this suite is a green run that asserted
 * nothing about the product's central feature — which is exactly what happens when the runner
 * cannot see the required environment variables.
 */
import { createBdd } from 'playwright-bdd';
import { test } from '../fixtures/index.js';

const { Before } = createBdd(test);

/**
 * How long a scenario is allowed to take, and why the number lives here.
 *
 * The suite's default is 15 seconds, sized for a warm route, and it stays that way on purpose: a
 * cross-cutting check that hangs should fail fast and loudly. A JOURNEY is a different animal —
 * it walks several routes, and a route's first hit against the dev server is a compile. A failure
 * at 15s there is the application behaving correctly and slowly rather than incorrectly.
 *
 * Every assertion inside a journey polls, so a generous ceiling costs nothing when things are
 * fast. A GENERATION is minutes of real external work and gets its own, larger, ceiling.
 *
 * `@capability` is the same ceiling for the same reason. A feature under `capabilities/` asserts
 * one capability on its own rather than walking a path, but it still has to REACH the surface —
 * signing in and opening a screen is several routes compiling on first hit.
 */
Before({ tags: '@journey or @internal or @roles or @cross-cutting or @capability' }, async () => {
  test.setTimeout(120_000);
});

Before({ tags: '@generation' }, async () => {
  test.setTimeout(420_000);
});

/**
 * Scenarios that need optional infrastructure skip with a reason when it is absent.
 *
 * Replace this with the actual environment check for your project's generation or processing
 * stack. The pattern: read an env var the infrastructure provides, and skip when it is missing.
 */
// const generationStack = !!process.env.YOUR_GENERATION_STACK_URL;
//
// Before({ tags: '@generation' }, async () => {
//   test.skip(
//     !generationStack,
//     'The generation stack is not configured — this scenario needs the processing ' +
//       'infrastructure. A submitted request with no worker stays queued forever, which reads ' +
//       'as a hang rather than as a missing worker.',
//   );
// });

/**
 * A gap in the product, kept visible.
 *
 * `known-issues` scenarios are written against behaviour that does not exist yet, and they are
 * skipped HERE rather than with Gherkin's own `@skip` tag — because a feature whose scenarios are
 * all tagged `@skip` generates no spec at all, and a gap that appears nowhere in the report is a
 * gap nobody is reminded of. Skipped from a hook, the scenario is generated, listed, and carries
 * its reason.
 *
 * A scenario goes here only when the thing it asserts has NOT BEEN BUILT and the blocker can be
 * named in a sentence. A scenario that fails because the application is WRONG is a defect, and
 * belongs in `features/defects/` once it is fixed — never here. This directory is not a parking
 * space for red tests.
 */

/**
 * The escape hatch, and why it is a variable rather than a hand-edit.
 *
 * A `ui`, `fe` or `be` change has to be able to ASK whether the flow it was blocking now passes —
 * that is the evidence it shipped anything. Without a way to run these, the only options were
 * deleting the tag before the change that owns the deletion (which loses the record of what was
 * pending) or editing the feature file for a run and putting it back (which is a hand-edit nobody
 * reviews and which silently becomes permanent the one time somebody forgets).
 *
 * Off by default, so an ordinary run and CI are unchanged: these scenarios stay skipped and stay
 * listed. The env var is a deliberate act by somebody who wants to know.
 */
const RUN_KNOWN_ISSUES = process.env.YOUR_PROJECT_RUN_KNOWN_ISSUES === '1';

Before({ tags: '@known-issue' }, async () => {
  test.skip(
    !RUN_KNOWN_ISSUES,
    'Written against behaviour that does not exist yet. The feature file names what must ' +
      'ship first; it comes off this tag in the change that ships it, which is what proves ' +
      'the change. Run them anyway with YOUR_PROJECT_RUN_KNOWN_ISSUES=1.',
  );
});
