import { defineConfig, devices } from '@playwright/test';
import { cucumberReporter, defineBddConfig } from 'playwright-bdd';

/*
 * THE RUNNER READS THE SAME `.env` THE SERVERS DO.
 *
 * Without this the Playwright process has whatever the shell exported and nothing else, while
 * the API and the client — started by `webServer` below, or already running — have the repo's
 * `.env`. The two disagreeing is not a tidiness problem; it silently changes what the suite
 * PROVES. A scenario that decides whether to run by reading an environment variable will SKIP
 * on a machine where that variable is only in `.env`, and a skip nobody asked for is worse than
 * a failure.
 *
 * `process.loadEnvFile` (Node 20.12+) rather than dotenv: no dependency, and it does not
 * override variables already in the environment — so an explicit shell export still wins.
 */
try {
  process.loadEnvFile(new URL('../.env', import.meta.url).pathname);
} catch {
  /* No `.env` — CI, or a machine configured entirely through the shell. Every consumer of these
   * variables already has a documented default or a skip, so this is not an error. */
}

/*
 * ── THE BDD LAYER ─────────────────────────────────────────────────────────────────────────
 *
 * `features/` holds Gherkin, `steps/` holds the code that carries it out, and playwright-bdd
 * generates a real Playwright spec per feature into `.features-gen/`. The generated files are
 * build output — never edited, never committed — and the thing a human reads is the `.feature`.
 *
 * `featuresRoot` is what makes the generated tree mirror the feature tree, so a failure in
 * `features/journeys/first-project.feature` reports under `journeys/` and the directory split
 * keeps meaning what `AGENTS.md` says it means.
 *
 * `missingSteps: 'fail-on-gen'` is deliberate and is the whole safety property of this layout: a
 * sentence in a feature file with no step behind it FAILS THE GENERATION, loudly, before anything
 * runs. The alternative — a scenario that quietly skips — is how a business-readable suite starts
 * describing behaviour nobody implemented, which is worse than having no suite.
 *
 * `aiFix.promptAttachment` attaches a ready-made prompt to every failure, carrying the feature,
 * the failing step, the step's own source and the page snapshot. That is what makes a red run
 * something an agent can act on rather than something it has to reconstruct.
 */
const bddTestDir = defineBddConfig({
  features: 'features/**/*.feature',
  featuresRoot: 'features',
  steps: ['steps/**/*.ts', 'fixtures/*.ts'],
  outputDir: '.features-gen',
  missingSteps: 'fail-on-gen',
  quotes: 'single',
  aiFix: { promptAttachment: true },
});

/* ── PORTS ──────────────────────────────────────────────────────────────────────────────────
 *
 * Replace these with your project's actual ports. They are variables rather than literals so
 * that a second worktree or a CI runner can shift them without editing this file.
 */
const clientPort = process.env.E2E_CLIENT_PORT ?? '3000';
const apiPort = process.env.E2E_API_PORT ?? '3001';
const baseURL = `http://localhost:${clientPort}/`;

export default defineConfig({
  testDir: bddTestDir,
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  /*
   * ONE WORKER LOCALLY, and the reason is the database rather than the machine.
   *
   * There is no per-worker database isolation, and this suite mutates what it reads. Four
   * workers on one seeded database are four copies of every mutating flow racing each other,
   * and the failures land on whichever scenario looked second. CI keeps four because its
   * database is its own, and `retries: 2` covers the machine being slow.
   */
  workers: process.env.CI ? 4 : 1,
  timeout: 15_000,
  /*
   * THREE REPORTS, AND ONLY ONE OF THEM IS FOR ENGINEERS.
   *
   * `cucumberReporter('html')` is the deliverable: a document whose unit is the SCENARIO and
   * whose lines are the Given/When/Then somebody in the business wrote, with the screenshot and
   * the trace hung off the step that failed. It is the answer to "what does a green run actually
   * prove" — you read the steps, and they are sentences about the business rather than locators.
   *
   * `message` is the same run as newline-delimited Cucumber Messages: machine-readable, diffable
   * between runs, and the input any other tool would want.
   *
   * The Playwright HTML report stays because it is the only one that carries the trace viewer,
   * and `list` stays because a terminal run has to say something while it works.
   */
  reporter: [
    /*
     * `externalAttachments` writes each screenshot and trace as a file next to the report
     * rather than base64 inside it, which is what keeps the HTML openable and what lets the
     * embedded trace viewer work.
     */
    cucumberReporter('html', {
      outputFile: 'reports/business-flows.html',
      externalAttachments: true,
    }),
    cucumberReporter('message', { outputFile: 'reports/messages.ndjson' }),
    ['html', { outputFolder: 'playwright-report', open: 'never' }],
    /*
     * The fourth reporter reports on the RUN rather than on the tests: how much of the suite
     * was selected, and — when that is more than any one feature — how to select less. The
     * convention is that a local run is a slice and CI runs the suite.
     */
    ['./run-scope.reporter.ts'],
    ...(process.env.CI ? ([['github'], ['list']] as const) : ([['list']] as const)),
  ],

  use: {
    baseURL,
    /*
     * WHAT THE REPORT CAN SHOW IS DECIDED HERE, NOT IN THE REPORTER.
     *
     * The Cucumber report attaches screenshots, videos and traces that PLAYWRIGHT recorded;
     * it records none itself. `retain-on-failure` records always and deletes on success: the
     * cost is paid by failing runs, which are the runs somebody is about to investigate.
     */
    screenshot: 'only-on-failure',
    trace: 'retain-on-failure',
    video: 'off',

    /*
     * ── HOW LONG A SINGLE STUCK ACTION MAY COST ──────────────────────────────────────────
     *
     * Unset, an action inherits the TEST timeout. A journey is allowed 120 seconds, so one
     * locator that matches nothing hangs until the whole budget is gone and then reports
     * "Test timeout exceeded", which says nothing about the screen. Playwright's own message
     * for a timed-out action names the locator and prints the call log; the test-level one
     * does not, and it arrives two minutes later.
     *
     * TWENTY SECONDS is chosen against the slowest legitimate action rather than the fastest.
     * Anything genuinely slower than this is a WAIT, and a wait is written as one —
     * `expect.poll` and `toBeVisible({ timeout })` set their own ceilings.
     *
     * Navigation is separate and larger: the first hit to a route in dev mode pays for a
     * compile, and that can genuinely take most of a minute.
     */
    actionTimeout: 20_000,
    navigationTimeout: 60_000,
  },

  /*
   * ── TWO PROJECTS ────────────────────────────────────────────────────────────────────────
   *
   * `seed` signs in as the bootstrap operator and runs the seed routines through the admin
   * API before any scenario runs. A dependency rather than `globalSetup`: this needs the API
   * listening, `webServer` entries come up before any project runs, and `globalSetup`'s
   * ordering against them has moved between Playwright versions.
   *
   * `flows` is the one project that runs scenarios — desktop Chromium. One project cannot
   * accidentally inherit a `testDir` pointing at nothing, which is a mistake that silently
   * produces a pass.
   */
  projects: [
    {
      name: 'seed',
      testDir: '.',
      testMatch: /seed\.setup\.ts$/,
    },
    {
      name: 'flows',
      use: { ...devices['Desktop Chrome'] },
      dependencies: ['seed'],
    },
  ],

  /* THE API COMES FIRST, AND THE ORDER IS LOAD-BEARING.
   *
   * Playwright brings these up in array order, waiting for each one's `url` to answer before
   * starting the next. If the client depends on the API (e.g. resolving a session on every
   * render), putting the client first deadlocks: it cannot answer its readiness probe until
   * the API is up, and the API is never started because Playwright is still waiting for the
   * client.
   *
   * Replace the commands and URLs below with your project's actual dev server commands.
   */
  webServer: [
    {
      // Replace with your API dev server command
      command: 'your-package-manager run dev:api',
      env: {
        /* Raise the per-minute rate limit for a run. A suite compresses a person's day into
         * two minutes; the production limit is sized for the person and exhausting it makes
         * every page render a rate-limit error. */
        API_RATE_LIMIT_PER_MINUTE: '2000',
        API_PORT: apiPort,
      },
      url: `http://localhost:${apiPort}/health`,
      reuseExistingServer: !process.env.CI,
    },
    {
      // Replace with your client dev server command
      command: 'your-package-manager run dev:client',
      env: {
        PORT: clientPort,
      },
      url: baseURL,
      reuseExistingServer: !process.env.CI,
    },
  ],
});
