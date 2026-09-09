/**
 * The cheapest claims in the suite: steps for `features/cross-cutting/health.feature`.
 *
 * Every other scenario silently depends on these — the front door, the sign-in screen, and both
 * halves of the application answering for themselves. They are worth their own file because their
 * failures are indistinguishable from everybody else's at a glance: a run where every scenario
 * fails at "signed in" is either a broken product or an application that never started, and these
 * say which within seconds.
 */
import { expect } from '@playwright/test';
import { createBdd } from 'playwright-bdd';
import { test } from '../fixtures/index.js';

const { Given, When, Then } = createBdd(test);

/* ── The front door ────────────────────────────────────────────────────────────────────────── */

Given('nobody is signed in', async () => {
  /* A fresh browser context has no session. This step exists so the Gherkin reads correctly
   * and because the sentence may eventually do cleanup — clearing cookies, for example — when
   * scenarios start sharing contexts. */
});

When('they open the application root', async ({ page }) => {
  await page.goto('/');
});

When('they open the sign-in screen', async ({ page }) => {
  await page.goto('sign-in');
});

Then('they are asked to sign in, rather than shown an error', async ({ page }) => {
  /* The claim is the form, not the URL. A redirect to `/sign-in` that renders an error page
   * would pass a URL check and fail this one. */
  await expect(page.getByLabel('Email')).toBeVisible();
  await expect(page.getByLabel('Password')).toBeVisible();
});

Then('it offers an email, a password and a way in', async ({ page }) => {
  await expect(page.getByLabel('Email')).toBeVisible();
  await expect(page.getByLabel('Password')).toBeVisible();
  await expect(page.getByRole('button', { name: /sign in|log in/i })).toBeVisible();
});

/* ── Both halves answer for themselves ─────────────────────────────────────────────────────── */

When('the API\'s health endpoint is asked', async ({ request, world }) => {
  /* Replace with your project's actual API health route. The health endpoint should be outside
   * any application prefix so it answers regardless of configuration. */
  const response = await request.get('your API health endpoint URL');
  world.lastStatus = response.status();
});

When('the client\'s health endpoint is asked', async ({ request, world }) => {
  /* Replace with your project's actual client health route. */
  const response = await request.get('your client health endpoint URL');
  world.lastStatus = response.status();
});

Then('it answers that it is ok', async ({ world }) => {
  expect(world.lastStatus, 'the health endpoint did not answer 200').toBe(200);
});
