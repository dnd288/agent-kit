/**
 * The shared test surface for E2E suites.
 *
 * Import from here instead of `@playwright/test` directly. The indirection is what lets a
 * fixture reach every suite at once: a database-backed fixture — seed per worker, truncate
 * between tests — will be added here and every suite picks it up without changing a line. A
 * suite that imports `@playwright/test` directly today would be the one left out tomorrow.
 *
 * ── EVERY TEST ARRIVES FROM ITS OWN ADDRESS ──────────────────────────────────────────────
 *
 * This is the one fixture that exists today, and it is here because the alternative is a
 * suite that cannot go green.
 *
 * Sign-in may be capped per address — a credential-stuffing guard, N attempts per M minutes.
 * The suite signs in many times. Left alone, every request arrives from the loopback address,
 * they all land in ONE bucket, and the run dies partway through with the API answering 429 to
 * its own session reads — every page then renders an error and the failures read as application
 * defects.
 *
 * Giving each test its own address is not a bypass of that cap; it is the topology the cap
 * was designed against. In production every user has their own address, and a suite in which
 * hundreds of different people share one address models nothing real. With this fixture the
 * limiter is genuinely exercised — each synthetic user gets N attempts, and a spec that signs
 * in twice still spends two of them.
 */
import type { APIRequest, Browser } from '@playwright/test';
import { test as base, request } from '@playwright/test';

/** The options `request.newContext()` takes, tied to the factory `apiContext()` forwards to. */
type APIRequestOptions = NonNullable<Parameters<APIRequest['newContext']>[0]>;

export { expect, type Page, type Request } from '@playwright/test';

/**
 * A stable address for one test.
 *
 * Carrier-grade NAT space (100.64.0.0/10, RFC 6598) on purpose: it is routable-looking and,
 * unlike 10/8 and 192.168/16, it is NOT one of the ranges a reverse proxy typically treats as
 * an internal hop — so the API reads it as a client rather than discarding it. TEST-NET
 * (203.0.113.0/24) would be the more obvious choice and is too small: it holds 256 addresses
 * and a large suite can run more tests than that.
 *
 * Derived from the test's own id rather than a counter, so it is the same on a retry — a
 * retry that landed in a fresh bucket would hide a test that genuinely exhausted its own.
 */
export function addressForTest(testId: string): string {
  let hash = 0;
  for (const character of testId) hash = (hash * 31 + character.charCodeAt(0)) >>> 0;

  /* Two octets of room — 65,536 buckets. A collision costs nothing anyway: two tests sharing a
   * bucket still have N attempts between them, and they sign in once. */
  return `100.64.${(hash >>> 8) & 0xff}.${hash & 0xff}`;
}

/**
 * The address of the test currently running in THIS worker.
 *
 * A module-level value rather than a fixture parameter because the two things that need it —
 * the `browser` proxy below and `apiContext()` — are reachable from places a test-scoped
 * fixture is not: `browser` is worker-scoped, and `apiContext()` is called from module-level
 * helpers that take no fixtures. Safe because Playwright runs a worker's tests one at a time,
 * and the auto fixture below sets it before each one.
 */
let currentAddress = '100.64.0.0';

/** The header every context this test opens must carry. */
function forwardedFrom(address: string): Record<string, string> {
  return { 'x-forwarded-for': address };
}

/**
 * An API request context aimed straight at the API, arriving from this test's address.
 *
 * Some specs talk to the API directly — the half of a flow that has no reachable surface yet
 * — and they do it with `request.newContext()` imported from `@playwright/test`, which no
 * fixture can reach. Without this they arrive from the loopback address: the one exhausted
 * bucket the whole fixture exists to escape, and they sign in.
 */
export async function apiContext(options: APIRequestOptions = {}) {
  return request.newContext({
    ...options,
    extraHTTPHeaders: { ...forwardedFrom(currentAddress), ...options.extraHTTPHeaders },
  });
}

export const test = base.extend<{ address: string }>({
  /* The default `page` and its context.
   *
   * The empty pattern is required, not sloppy: Playwright READS the destructuring to discover a
   * fixture's dependencies and refuses a named parameter outright — "First argument must use the
   * object destructuring pattern" — so `{}` is how this API spells "depends on nothing". */
  // biome-ignore lint/correctness/noEmptyPattern: Playwright requires the destructuring pattern.
  extraHTTPHeaders: async ({}, use, testInfo) => {
    await use(forwardedFrom(addressForTest(testInfo.testId)));
  },

  /* Publishes this test's address for the two consumers a fixture parameter cannot reach.
   * `auto` so no spec has to remember to depend on it. */
  address: [
    // biome-ignore lint/correctness/noEmptyPattern: Playwright requires the destructuring pattern.
    async ({}, use, testInfo) => {
      currentAddress = addressForTest(testInfo.testId);
      await use(currentAddress);
    },
    { auto: true },
  ],

  /*
   * THE CONTEXTS A SPEC OPENS ITSELF.
   *
   * `extraHTTPHeaders` reaches the `page` fixture's context and nothing else, and this suite
   * may open contexts by hand: a stranger, a second identity, a returning user. Every one of
   * those falls back to the loopback address — the single exhausted bucket this fixture exists
   * to escape — while the default page is correctly isolated.
   *
   * Wrapping the factory rather than editing the call sites, because forgetting one is
   * invisible: the context still works, the test still passes, and the suite only fails once
   * enough of them accumulate.
   */
  browser: async ({ browser }, use) => {
    await use(
      new Proxy(browser, {
        get(target, property, receiver) {
          if (property !== 'newContext') return Reflect.get(target, property, receiver);

          return (options: Parameters<Browser['newContext']>[0] = {}) =>
            target.newContext({
              ...options,
              extraHTTPHeaders: {
                ...forwardedFrom(currentAddress),
                ...options.extraHTTPHeaders,
              },
            });
        },
      }),
    );
  },
});
