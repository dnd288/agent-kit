/**
 * The sign-in screen.
 *
 * ── Why page objects at all ───────────────────────────────────────────────────────────────
 *
 * A step definition is written once and read by everybody: `When the user signs in` appears in
 * a dozen features and the person reading the report is not going to open a locator. The page
 * object is where the locator's REASONING lives, so the step stays a sentence about the business
 * and the "why this query and not that one" stays next to the query.
 *
 * The rule: role, then label, then text the component owns. A structural selector is the last
 * rung and carries a comment saying why no role exists.
 */
import type { Page } from '@playwright/test';
import type { Account } from '../lib/accounts.js';
import { expect } from '../lib/suite.js';

export class SignInPage {
  constructor(private readonly page: Page) {}

  async open(path = 'sign-in'): Promise<void> {
    await this.page.goto(path);
  }

  /**
   * Fill the form and wait for the Server Action's POST, not merely for the click.
   *
   * The action is what carries the API's `Set-Cookie` onto a response the browser receives, so
   * returning at the click leaves the next line racing an in-flight submission — a `reload()`
   * aborts it, and the test then reports a missing session for a sign-in that never finished.
   * Assertions that poll (`toBeVisible`) hide the race; anything reading state once, like
   * `context.cookies()`, does not.
   *
   * Replace the button name and the response URL pattern with your project's actual sign-in
   * surface.
   */
  async submit(account: Account): Promise<void> {
    await this.page.getByLabel('Email').fill(account.email);
    await this.page.getByLabel('Password').fill(account.password);

    await Promise.all([
      this.page.waitForResponse(
        (response) =>
          response.request().method() === 'POST' && response.url().includes('sign-in'),
      ),
      /* Replace the button name with whatever your sign-in screen calls its submit control.
       * The role query is deliberate: a `<button>` that is really a `<div>` would fail this
       * check, which is the point. */
      this.page.getByRole('button', { name: /sign in|log in/i }).click(),
    ]);
  }

  /**
   * Sign in from wherever the browser currently is, and wait until the authenticated surface
   * admits them.
   *
   * Replace the waitForURL pattern with the route your project sends a newly signed-in user to.
   */
  async signIn(account: Account): Promise<void> {
    await this.open();
    await this.submit(account);
    /* Replace with the URL pattern that proves the authenticated shell loaded. */
    await this.page.waitForURL(/\/dashboard|\/projects|\/home/);
    await expect(this.page.getByRole('button', { name: /sign in|log in/i })).toBeHidden();
  }

  /**
   * The form's own failure message.
   *
   * SCOPED INSIDE THE FORM, and that is the whole point of this getter existing. Some frameworks
   * inject a route announcer with `role="alert"` into every page, so a bare alert role matches
   * it: `toBeVisible()` would pass on pages showing no message at all, and `toBeHidden()` could
   * never pass. Scoping inside the form is what makes the assertion honest.
   */
  refusal() {
    return this.page.getByRole('form').getByRole('alert');
  }
}
