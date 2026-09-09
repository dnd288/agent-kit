/**
 * The suite seeds itself the way an operator seeds a deployed environment: through the API.
 *
 * ── Why this is not a seed script ─────────────────────────────────────────────────────────
 *
 * A seed script is a path only CI takes and nobody exercises in production. If the suite
 * seeds through the same admin endpoint an operator's click hits, a green run proves
 * something a seed script never could: that the application can actually populate an empty
 * database, which is the only way a real environment is ever populated.
 *
 * ── Why a setup PROJECT rather than globalSetup ───────────────────────────────────────────
 *
 * It needs the API to be listening. `webServer` entries come up before any project runs, so
 * a project the flows depend on is guaranteed to find one; `globalSetup`'s ordering relative
 * to `webServer` has changed across Playwright versions and is not worth depending on.
 *
 * It also means a failure here reports as a failed setup rather than as every scenario
 * failing at "invalid credentials".
 */

import { expect, test as setup } from '@playwright/test';
import { OPERATOR } from './lib/accounts.js';
import { API_ORIGIN } from './lib/api.js';

/* Long, and deliberately so: seed routines that ingest reference data can take tens of
 * seconds on a warm laptop. A CI runner under load is slower, and a timeout here reads as
 * "the application cannot seed" rather than as "the runner was busy". */
setup.setTimeout(10 * 60 * 1000);

/** The seed endpoint answers as soon as the row exists; the work runs behind it. */
const POLL_INTERVAL_MS = 500;

interface SeedResult {
  id: string;
  status: 'RUNNING' | 'SUCCEEDED' | 'FAILED';
  total: number;
  processed: number;
  succeeded: number;
  failed: number;
  message: string | null;
}

setup('the operator seeds the database from the admin API', async ({ request }) => {
  /* ── Sign in as the bootstrap administrator ──────────────────────────────────────────
   *
   * Replace the sign-in call below with your project's actual authentication flow.
   * The bootstrap admin is the only identity that exists before seeding — every other
   * account is written by the seed routine itself. */
  const signedIn = await request.post(`${API_ORIGIN}/auth/sign-in`, {
    data: { email: OPERATOR.email, password: OPERATOR.password },
  });

  const howToFixSignIn = [
    `could not sign in as the bootstrap operator (${OPERATOR.email}).`,
    'Set BOOTSTRAP_ADMIN_EMAIL and BOOTSTRAP_ADMIN_PASSWORD in your environment',
    'and run your database seed command — it writes the only account that exists',
    'before the seed routines run.',
    await signedIn.text(),
  ].join(' ');
  expect(signedIn.status(), howToFixSignIn).toBe(200);

  /* ── Run seed routines ───────────────────────────────────────────────────────────────
   *
   * Replace the endpoint and routine names below with your project's actual seed API.
   * Order may matter: demo accounts that belong to a company need the company written
   * first. */
  for (const routine of ['SEED_CORE', 'SEED_DEMO']) {
    const started = await request.post(`${API_ORIGIN}/admin/seed`, {
      data: { routine },
    });

    /* 409 is "already running", which on a reused stack means a previous run is still
     * going. Fall through to the poll. */
    if (started.status() !== 409) {
      expect(started.status(), `${routine} would not start: ${await started.text()}`).toBe(200);
    }

    const result = await settle(routine);

    expect(
      result.status,
      `${routine} finished ${result.status} — ${result.processed}/${result.total} processed, ` +
        `${result.failed} failed. ${result.message ?? ''}`,
    ).toBe('SUCCEEDED');

    /* A routine that succeeded having written nothing is the failure this catches. Without
     * it, every spec that signs in as a demo account fails one by one with "invalid
     * credentials" instead of here, once. */
    expect(
      result.succeeded,
      `${routine} succeeded without writing anything. The demo accounts the suite signs in as were not written.`,
    ).toBeGreaterThan(0);

    process.stdout.write(`[seed] ${routine}: ${result.succeeded}/${result.total} written\n`);
  }

  /** Poll until the routine leaves RUNNING. */
  async function settle(routine: string): Promise<SeedResult> {
    for (;;) {
      const response = await request.get(`${API_ORIGIN}/admin/seed/${routine}`);
      expect(response.status(), await response.text()).toBe(200);
      const result = (await response.json()) as SeedResult;
      if (result.status !== 'RUNNING') return result;
      await new Promise((resolve) => setTimeout(resolve, POLL_INTERVAL_MS));
    }
  }
});
