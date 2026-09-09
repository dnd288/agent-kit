/**
 * The BDD test instance — where a step definition gets its page objects and its memory.
 *
 * ── Why steps take fixtures and never construct anything ──────────────────────────────────
 *
 * A step definition is a sentence somebody in the business can read. The moment it contains
 * `new SignInPage(page)` it is code again, and the second step that needs the same page object
 * has to construct a second one. Fixtures make the page objects the vocabulary of a step and keep
 * the wiring in one file — which is also what lets a scenario be re-read as a stack of sentences
 * in the report.
 *
 * ── The `world` fixture, and why it is not a global ───────────────────────────────────────
 *
 * Given/When/Then are separate functions. "the user has created a project" has to hand the
 * project id to "the user removes the item", and the only honest channel between them is a
 * per-scenario object. Test-scoped, so it is fresh for every scenario and two scenarios can
 * never see each other's subject — which matters here more than usual, because there is no
 * per-worker database isolation (`lib/suite.ts` says why) and cross-talk would look like a
 * product defect.
 *
 * ── Why `test` comes from `lib/suite.js` and not from `@playwright/test` ──────────────────
 *
 * `suite.ts` carries the fixture that gives every test its own `x-forwarded-for`. Sign-in may
 * be capped per address (a credential-stuffing guard) and this suite signs in many times.
 * Without that fixture every one of them lands in one bucket and the run may die partway through
 * with the API answering 429. Building the BDD instance on top of `base` rather than on
 * `suite`'s `test` would silently drop it.
 */
import type { Page } from '@playwright/test';
import { mergeTests } from '@playwright/test';
import { test as bdd } from 'playwright-bdd';
import { test as suite } from '../lib/suite.js';
import { SignInPage } from '../pom/sign-in.page.js';

/**
 * The two test instances, merged rather than nested.
 *
 * `mergeTests` is the API for this and the shortcut is not: spreading one test object into
 * another's `extend` looks like it works and silently drops the fixtures, because a Playwright
 * test object is not a plain fixtures record. What would be lost here is the `address` fixture —
 * an `auto` one, so nothing declares a dependency on it and nothing would fail loudly. The suite
 * would simply start answering 429 partway through, and the failures would read as the API being
 * broken.
 */
const base = mergeTests(bdd, suite);

/**
 * What one scenario knows about itself.
 *
 * Everything here is written by a Given or a When and read by a Then. Nothing is seeded: a
 * scenario that needs a resource creates one, because a listing paginates and a seeded row
 * borrowed by name can fall off page one once another scenario adds rows in front of it.
 *
 * Add fields as your project's flows grow. The starter set covers the patterns most projects
 * need immediately; the convention is that every field is optional and every read goes through
 * a guard.
 */
export interface ScenarioWorld {
  /** The project or primary resource this scenario is working on. */
  projectId?: string;
  /** A secondary resource id — a child, a related record, or a created artefact. */
  resourceId?: string;
  /** A share link this scenario minted. */
  shareUrl?: string;
  /** The last response a navigation produced — status codes are assertions here. */
  lastStatus?: number;
  /** A page belonging to somebody who is not the scenario's subject — a stranger, a colleague. */
  openedPage?: Page;
  /** The name a scenario typed, so an assertion can look for it afterwards. */
  createdName?: string;
  /** A value read BEFORE a mutation, so the Then can compare before and after. */
  valueBefore?: string;
  /** A credential a scenario created or received — an address, a token, a link. */
  credential?: string;
}

/**
 * Guard: fail with a readable message rather than `undefined is not a string`.
 *
 * Every read of `world` goes through a guard so a scenario missing a Given fails with
 * "this scenario has no project — a step that acts on one must be preceded by a Given that
 * creates or finds one" rather than a TypeError deep in a page object.
 */
export function requireId(id: string | undefined, noun = 'resource'): string {
  if (!id) {
    throw new Error(
      `This scenario has no ${noun} — a step that acts on one must be preceded by a ` +
        `Given that creates or finds one.`,
    );
  }
  return id;
}

/**
 * Add page object fixtures here as the project grows. The pattern:
 *
 * ```ts
 * import { DashboardPage } from '../pom/dashboard.page.js';
 * // then in the extend block:
 * dashboard: async ({ page }, use) => { await use(new DashboardPage(page)); },
 * ```
 *
 * A fixture named `console` shadows the global inside every step that destructures it, so
 * prefer `adminConsole` or similar when the POM wraps a console-like surface.
 */
export const test = base.extend<{
  signIn: SignInPage;
  world: ScenarioWorld;
}>({
  signIn: async ({ page }, use) => {
    await use(new SignInPage(page));
  },

  // biome-ignore lint/correctness/noEmptyPattern: Playwright requires the destructuring pattern.
  world: async ({}, use) => {
    await use({});
  },
});

export { expect } from '../lib/suite.js';
