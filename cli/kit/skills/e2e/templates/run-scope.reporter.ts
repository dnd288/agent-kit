import type { FullConfig, Reporter, Suite, TestCase } from '@playwright/test/reporter';

/**
 * WHAT A RUN IS ABOUT TO COST, SAID BEFORE IT COSTS IT.
 *
 * The convention is **run the slice, and let CI run the suite**. A full local run on one
 * worker with 120-second journey budgets is tens of minutes before anything is wrong with it.
 * This reporter makes the convention visible at the only moment it can still be acted on:
 * the first second of the run, while Ctrl-C is still cheap.
 *
 * It is a NOTICE rather than a gate. Nothing here refuses a full run, because a full run is
 * exactly right before a merge and on CI — and a flag to silence a warning is a flag that
 * accumulates. It fires only when the selection is bigger than any single feature and never
 * on CI.
 */

/**
 * Above this many selected scenarios, the run is not a slice.
 *
 * The number is derived rather than chosen: the largest single feature in a typical suite is
 * around 25–30 scenarios once outlines expand. A selection bigger than that therefore cannot
 * be "one feature", and is either a directory or the lot. Both are worth a word; a feature
 * is not.
 */
const LARGEST_FEATURE = 30;

const RULE = '─'.repeat(92);

function formatDuration(ms: number): string {
  const seconds = Math.round(ms / 1000);
  if (seconds < 90) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  return minutes < 60 ? `${minutes}m` : `${Math.floor(minutes / 60)}h ${minutes % 60}m`;
}

/**
 * playwright-bdd carries a scenario's Gherkin tags onto the Playwright test, so `@spends` is
 * readable here. If a future version stops doing that, the count comes back zero and the line
 * is dropped rather than the banner lying about the bill.
 */
function billed(tests: TestCase[]): number {
  return tests.filter((test) => test.tags?.includes('@spends')).length;
}

export default class RunScopeReporter implements Reporter {
  private announced = false;
  private startedAt = 0;
  private selected = 0;

  onBegin(config: FullConfig, suite: Suite) {
    // `--list` prints the tree and runs nothing, so there is no cost to warn about.
    if (process.env.CI || process.argv.includes('--list')) return;

    const tests = suite.allTests();
    this.selected = tests.length;
    if (this.selected <= LARGEST_FEATURE) return;

    this.announced = true;
    this.startedAt = Date.now();
    const workers = config.workers;
    const spends = billed(tests);

    const lines = [
      '',
      RULE,
      ` ${this.selected} scenarios queued on ${workers} worker${workers === 1 ? '' : 's'}. This is the pre-merge run, not the loop.`,
      ...(spends > 0
        ? [` ${spends} of them call an external service (@spends) and may be billed.`]
        : []),
      '',
      ' Run the slice your change touches instead:',
      '   one directory        your e2e test command journeys/',
      '   one feature          your e2e test command journeys/first-project',
      '   one scenario         your e2e test command --grep "creates a project"',
      '   what just failed     your e2e test command --last-failed',
      '',
      ' The whole suite belongs on CI.',
      ' your e2e path/README.md § Running a slice.',
      RULE,
      '',
    ];
    console.log(lines.join('\n'));
  }

  onEnd() {
    if (!this.announced) return;
    console.log(
      `\n${RULE}\n ${this.selected} scenarios, ${formatDuration(Date.now() - this.startedAt)}. A slice would have been minutes — your e2e path/README.md § Running a slice.\n${RULE}\n`,
    );
  }
}
