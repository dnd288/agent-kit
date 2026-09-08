---
name: optimization-loop
description: "Making code measurably faster, smaller, or cheaper without guessing — the loop where every kept change is justified by a number and correctness never moves. Load this when asked to optimise, speed up, reduce latency/memory/bundle size/cost, fix a performance regression, meet an SLO or budget that is being missed, remove an N+1, warm a cold path, or when a profile points at a hot spot. Owns the BASELINE → CHANGE → MEASURE → KEEP-OR-REVERT loop, the rule that a change with no measured improvement is reverted, the discipline of one variable at a time, and profiling before guessing. Routes to simplicity for whether to optimise at all (premature optimisation is the default mistake), to refactor when the hot path must be reshaped before it can be tuned, to tdd/test for pinning correctness while runtime behaviour changes, and to review's performance axis for the gate. This is the opposite of refactor: here runtime behaviour changes on purpose and the observable answer must not."
---

# The optimization loop

**The defining constraint: every change you keep is justified by a number, and the observable answer
never changes.** A faster wrong answer is not an optimisation, it is a regression that happened to
run quickly. Everything below exists to keep those two facts checkable.

This is the mirror image of `refactor`. There, the shape changes and behaviour is pinned still.
Here, runtime behaviour — time, memory, bytes, cost — changes on purpose, and *correctness* is the
thing pinned still. The loop is what stops "optimising" from meaning "changing things until it feels
snappier".

## First: should you optimise at all?

Most of the time the answer is no, and saying so is the highest-value move in this skill.

- **Is there a number that is actually missed?** An SLO, a budget, a benchmark, a user-visible
  stall. No target means no optimisation — `simplicity` owns the case where the honest answer is
  "this is fast enough and the tuning would only add complexity". Premature optimisation is the
  default mistake, not a rare one.
- **Do you know where the time goes?** If not, you are about to optimise a guess. Profile first
  (next section). The bottleneck is almost never where it feels like it is.
- **Is the shape wrong?** If the hot path is tangled enough that you cannot change one variable
  cleanly, the honest sequence is `refactor` first (behaviour pinned), then optimise. Not both at
  once — you will not be able to attribute the delta.

## The loop

**BASELINE → CHANGE → MEASURE → KEEP-OR-REVERT.** Each step feeds the next; skipping BASELINE or
MEASURE turns the whole thing back into guessing.

### 1. BASELINE — measure before you touch anything

Run `your benchmark command` (or the profiler) under conditions you can reproduce, and **write the
number down** along with the target you are trying to hit. This is what MEASURE compares against, and
"it felt faster" is not a number. Record the noise too — run it more than once so you know what a
real improvement has to beat.

**Profile before you guess.** The hot path is measured, not assumed. A day spent optimising the 2%
is a day the 80% did not move.

### 2. CHANGE — one variable at a time

One optimisation per iteration. Two changes at once and MEASURE cannot tell you which one paid — so
a kept change may be carrying a silent regression alongside a real win.

**Pin correctness before it moves.** The test suite is green before the change and green after it —
`your validation command`, plus the specific tests that observe the behaviour you are speeding up
(`tdd`/`test` own which mode observes which claim). If you cannot prove the answer is unchanged, you
are not optimising, you are rewriting.

### 3. MEASURE — the same benchmark, the same way

Re-run the identical benchmark under identical conditions and compare to BASELINE. Different input,
different machine, or a warm cache where BASELINE was cold makes the comparison meaningless.

### 4. KEEP-OR-REVERT — the number decides, not the effort

Keep the change **only if all three hold**: the target metric improved beyond the noise you
measured, correctness held, and nothing else you care about regressed (a latency win that doubles
memory is a trade to make on purpose, not by accident). Otherwise **revert it** — a change that
bought no measured improvement bought only complexity, and keeping it because it "should be faster"
is how a codebase fills with cargo-cult tuning.

**Record the reverts.** A one-line note of what was tried and the number it did (or did not) move
saves the next person from re-trying it. `problem-solving`'s dev note is a good shape for this when
the attempt was substantial.

## When NOT to optimise

- **No target.** Fast enough is done. See `simplicity`.
- **The bottleneck is elsewhere.** You profiled the wrong layer. Measure, then move.
- **The gain costs more than it returns.** An unreadable micro-optimisation for a 1% win on a path
  nobody waits on is a net loss — the bill is paid by every future reader.
- **The shape is about to change.** Tuning code that a pending feature will replace spends the
  budget twice.

## Before you finish

- [ ] A target number existed before the work started
- [ ] BASELINE and after numbers are both stated, from the same benchmark run the same way
- [ ] Correctness is green — the observable answer did not change
- [ ] One variable per kept change, so each delta is attributable
- [ ] Every change that did not beat the noise was reverted, not kept "just in case"
- [ ] Any regression traded away (memory for speed, etc.) was a deliberate, stated trade

## Read next

- [`simplicity`](../simplicity/SKILL.md) — whether to optimise at all, and the cost of structure the
  problem did not ask for.
- [`refactor`](../refactor/SKILL.md) — reshaping the hot path safely before tuning it.
- [`tdd`](../tdd/SKILL.md) / [`test`](../test/SKILL.md) — which mode pins the correctness you must
  hold still.
- [`review`](../review/SKILL.md) — the performance axis at the merge gate.
- [`problem-solving`](../problem-solving/SKILL.md) — the dev-note shape for recording a substantial attempt.
