# playwright.config.ts — annotated guide

The scaffold's `playwright.config.ts` is opinionated. This reference explains each setting
block: what it does, why the value was chosen, and what breaks when it is wrong.

## The BDD layer

```ts
const bddTestDir = defineBddConfig({
  features: 'features/**/*.feature',
  featuresRoot: 'features',
  steps: ['steps/**/*.ts', 'fixtures/*.ts'],
  outputDir: '.features-gen',
  missingSteps: 'fail-on-gen',
  quotes: 'single',
  aiFix: { promptAttachment: true },
});
```

| Setting | Why | What breaks when it is wrong |
|---|---|---|
| `featuresRoot: 'features'` | Makes the generated tree mirror the feature tree, so a failure in `features/journeys/sharing.feature` reports under `journeys/` and the directory split keeps meaning what `AGENTS.md` says it means | Without it, generated paths are flat and Playwright's positional filters (`journeys/`) stop working |
| `missingSteps: 'fail-on-gen'` | **The whole safety property.** A sentence with no step behind it FAILS generation, loudly, before anything runs. The alternative — a scenario that quietly skips — is how a business-readable suite starts describing behaviour nobody implemented | Without it, unfinished flows silently pass |
| `aiFix: { promptAttachment: true }` | Attaches a ready-made prompt to every failure — the feature, the failing step, the step source, and the ARIA snapshot. Makes a red run something an agent can act on | Without it, diagnosing a failure requires manually assembling context |
| `quotes: 'single'` | Consistency in generated code | No functional impact |

## Environment loading

```ts
try {
  process.loadEnvFile(new URL('../.env', import.meta.url).pathname);
} catch { /* no .env — CI or shell-configured */ }
```

The runner reads the same `.env` the servers do. Without this, the Playwright process has
whatever the shell exported and nothing else, while the API and client (started by `webServer`)
have the repo `.env`. The two disagreeing silently changes what the suite PROVES — scenarios
that check for external infrastructure environment variables will skip on a machine with a
fully configured stack, and the run goes green having asserted nothing.

`process.loadEnvFile` (Node 20.12+) does not override variables already in the environment, so
a shell export still wins. No dependency on dotenv.

## Workers

```ts
workers: process.env.CI ? 4 : 1,
```

**One locally, four on CI.** There is no per-worker database isolation, and this suite mutates
what it reads. Four workers on one seeded database are four copies of every mutating flow racing
each other; failures land on whichever scenario looked second. CI keeps four because its database
is fresh and its own, and `retries: 2` covers the machine being slow.

What breaks: running multiple workers locally on a shared database produces intermittent failures
that look like product defects but are data races between scenarios.

## Timeouts

```ts
timeout: 15_000,
use: {
  actionTimeout: 20_000,
  navigationTimeout: 60_000,
},
```

| Timeout | Value | Why |
|---|---|---|
| Test | 15s | Sized for a warm route. Journeys and capabilities override to 120s via tag hooks |
| Action | 20s | A stuck locator should fail in 20s, not borrow the scenario's full 120s budget. Without this, one missing element hangs silently until the test timeout, wasting minutes and producing a message about the clock instead of about the screen |
| Navigation | 60s | A cold `next dev` compiles a route on first hit. This is one round trip, but compilation of a large app can genuinely take most of a minute |

The tag hooks in `steps/hooks.ts` override the test timeout per tag:
- `@journey`, `@roles`, `@internal`, `@cross-cutting`, `@capability` → 120s
- `@generation` → 420s

## Reporters

```ts
reporter: [
  cucumberReporter('html', { outputFile: 'reports/business-flows.html', externalAttachments: true }),
  cucumberReporter('message', { outputFile: 'reports/messages.ndjson' }),
  ['html', { outputFolder: '../playwright-report', open: 'never' }],
  ['./run-scope.reporter.ts'],
  ...(process.env.CI ? ([['github'], ['list']] as const) : ([['list']] as const)),
],
```

| Reporter | Audience | What it carries |
|---|---|---|
| Cucumber HTML | **The business.** The deliverable — scenarios as Given/When/Then with screenshots and traces on the failing step | `externalAttachments: true` keeps it openable and lets the trace viewer work |
| Cucumber Messages | Machines | NDJSON — diffable between runs, input for other tools |
| Playwright HTML | Engineers | The only reporter carrying the embedded trace viewer |
| run-scope | The person about to wait | How much was selected, and how to select less |
| github | CI | Annotates the failing step on the pull request |
| list | Terminal | Something to read while the run works |

## Projects

```ts
projects: [
  { name: 'seed', testDir: '.', testMatch: /seed\.setup\.ts$/ },
  { name: 'flows', use: { ...devices['Desktop Chrome'] }, dependencies: ['seed'] },
],
```

**Why `seed` is a project, not `globalSetup`:** it needs the API listening. `webServer` entries
come up before any project runs, so a project the flows depend on is guaranteed to find one;
`globalSetup`'s ordering relative to `webServer` has changed across Playwright versions. A
failure in seed also reports as a failed setup rather than as N scenarios failing with the same
error.

**Why `testDir` is overridden on `seed`:** the default `testDir` is `bddTestDir` — the GENERATED
output of `defineBddConfig`. `seed.setup.ts` is hand-written and lives at the suite root, so a
project inheriting `bddTestDir` would match nothing and report a green having seeded no database.
That silently produces a PASS.

**One project, Chromium, everywhere.** The suite runs in one browser engine. A device matrix that
repeats every scenario at four widths proves the same sentence four times and reports the failure
under a name that does not say which width broke. A claim about a width says so in the scenario,
and a step sets it.

## webServer

```ts
webServer: [
  { command: '...api...', url: 'http://localhost:PORT/health', ... },
  { command: '...client...', url: 'http://localhost:PORT/', ... },
],
```

**API first, and the order is load-bearing.** Playwright brings these up in array order, waiting
for each one's `url` to answer before starting the next. With the client first this deadlocks:
the client's root layout resolves the session on every render, so it cannot answer its readiness
probe until the API is listening — and the API is never started because Playwright is still
waiting for the client. It times out with `ECONNREFUSED` and no output from the API.

**Rate limit:** the per-minute API rate limit should be raised for a run. A scenario walks
through steps in seconds that a real person takes minutes between, so the per-address budget
exhausts and the API starts answering 429 — every page then renders an error and the failures
read as a broken product.

**`reuseExistingServer`:** convenient but a trap. Servers started with the dev `DATABASE_URL`
write to the dev database, not the e2e database, and the run is no longer isolated.

## Screenshots, traces, video

```ts
use: {
  screenshot: 'only-on-failure',
  trace: 'retain-on-failure',
},
```

Both record always and delete on success: the cost is paid by failing runs, which are the runs
somebody is about to investigate. Without these, the business report promises a screenshot on
the failing step and carries none, and `on-first-retry` means that locally (where `retries` is 0)
no trace is ever written either.

Video is off by default. Recording costs wall-clock on every scenario. Enable it when evidence of
a green run is needed, not for debugging.
