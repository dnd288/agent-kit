---
name: verify
description: "The adversarial verification pass — prove a change satisfies its own specification, or prove it does not. Load this before marking any user-visible work ready, when a change's scenarios look green but nobody has challenged them, when a subagent's build needs checking against the change folder, or when asked to audit a change against its openspec spec. Owns the CLAIM → EXTRACT → DOUBT → RECONCILE → STOP cycle, which claims to extract from the change's specs, what counts as evidence (your testing documentation), and the three-round bound. Complements review: that reviews the change; this attacks it."
---

# Verifying a change against its specification

Passing tests are evidence, not proof. The change's own tests were written by the same
context that wrote the code — they share its blind spots. This skill is the discipline of
materialising a fresh-context reviewer, biased to **disprove**, before any user-visible work
stands. It is not a verdict on the finished artifact (that is `acme-review`); it is the
in-flight posture that catches wrong directions while course-correction is cheap.

## Source

Adapted from `doubt-driven-development` in
[`addyosmani/agent-skills`](https://github.com/addyosmani/agent-skills) (MIT). The project
variant takes its claim set from the change folder rather than from the author's head.

## The claim set: the specification

In this repository the claims are not what you think you built — they are the change's
your change tracking folder specs scenarios, written before the code, each one a
`#### Scenario:` with a WHEN and a THEN. A scenario is a claim that the system behaves a
certain way. The spec is the contract; the change is the artifact.

A change with no change folder cannot be verified by this skill, only reviewed — which is
the distinction the `spec-required` guard check draws, and the reason it warns.

## The cycle

```
- [ ] CLAIM — a scenario, named as a claim
- [ ] EXTRACT — the smallest unit that satisfies it, reasoning stripped
- [ ] DOUBT — a fresh-context reviewer, told to find failure
- [ ] RECONCILE — every finding classified against the artifact
- [ ] STOP — trivial findings, three cycles, or an explicit "ship it"
```

### CLAIM

For each scenario in the change's spec, write the claim compactly:

```
CLAIM: "Sign in shows the floating-card layout below the large breakpoint and the 50/50 split above it"
WHY THIS MATTERS: the responsive switch is the non-obvious part of this screen; a
                  regression reads as a design decision.
```

If a claim cannot be written compactly, the scenario is vague — that is a finding about the
spec, and it goes back to the spec, not around it.

### EXTRACT

The reviewer gets the artifact and the contract, not the journey. A scenario maps to the code
that implements it, the test that covers it, and the pixel baseline that records it. Strip
your reasoning — if you hand over conclusions you get back validation of your conclusions.
The unit must fit in one read; a scenario that needs a whole screen to verify is a scenario
that should have been decomposed.

### DOUBT

The reviewer's prompt **must be adversarial** — framing decides the answer:

```
Adversarial review. Find what is wrong with this artifact. Assume the author is
overconfident. Look for: unstated assumptions, edge cases, hidden coupling, ways the
contract could be violated, existing conventions this might break, failure under
unexpected input. Do NOT validate. Do NOT summarize. Find issues, or state explicitly
that you cannot find any after thorough examination.

CONTRACT: <the scenario, verbatim>
ARTIFACT: <the smallest unit that satisfies it>
```

**Pass the CONTRACT and the ARTIFACT only. Never pass the CLAIM.** Handing the reviewer your
conclusion biases it toward agreement. The reviewer must independently decide whether the
artifact satisfies the contract. In Claude Code, a fresh-context reviewer is the mechanism;
on any other agent, the same prompt to a fresh session.

**What counts as evidence.** your testing documentation owns the answer, and it is stricter
than it looks: a jsdom unit test cannot measure layout (`assertHasLayout` throws there for
exactly this reason); the browser modes exist to be used. A scenario about geometry is
verified by a `*.quality.test.tsx`. A scenario about **design fidelity has no mechanical verifier** — the visual gate is retired (your project architectural decisions) — so it is verified by eye and the report must say that rather than pointing at a green browser run. A scenario about copy
is verified under the pseudo-locale, because that is the only render where unbracketed text
on screen means hardcoded copy. Map the scenario to the mode that can actually observe it —
the five modes (unit, stories, quality, pixel, visual) are in `acme-tdd`.

### RECONCILE

The reviewer's output is data, not verdict. Classify each finding in precedence order:

1. **Contract misread** — flagged because the scenario was unclear. Fix the spec first, re-loop.
2. **Valid and actionable** — a real gap. Fix it, re-loop.
3. **Valid trade-off** — real but the fix costs more than the acceptance. Document it so the user sees it.
4. **Noise** — correct under context the reviewer lacked. Note it, and ask whether the scenario
   should have carried that context.

A fresh reviewer can be wrong because it lacks context. Do not defer just because it is fresh.

### STOP

Stop when the next round returns only trivial or already-considered findings, after three
cycles, or on an explicit "ship it". Three unresolved cycles is information about the
artifact — surface it, do not grind a fourth round. If three cycles is "obviously
insufficient", the artifact is too big: return to EXTRACT and decompose, do not lift the bound.

## Where this sits in the pipeline

The the feature workflow's verify gate is this skill bounded to three rounds, with the
observer schema-constrained and forbidden from fixing — the fixer is a separate pass that
only sees the failures. The same shape applies by hand: **observe and fix are different
agents, or different turns, never the same one.** A verifier that fixes its own findings has
stopped verifying.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "The tests pass" | They were written by the same context as the code. A scenario nobody challenged is a scenario nobody verified. |
| "Re-verifying is expensive" | Debugging a wrong screen in production is more expensive. The check is bounded; the bug is not. |
| "The reviewer will nitpick" | Only if unscoped. Constrain to "issues that would make this fail under the contract." |
| "I'll verify at the end with review" | Review is a verdict on the finished artifact. Verification catches the wrong direction while the fix is cheap. |
| "The reviewer disagreed, so I was wrong" | Disagreement is information, not verdict. Re-read the artifact, classify, then decide. |
| "The scenario is basically covered by the story" | A story is a fixture for a11y and render; it asserts almost nothing. The mode must match the claim. |

## Red flags

- Verifying a claim you wrote the answer to (passing the CLAIM to the reviewer)
- A reviewer prompt that asks "is this good?" instead of "find issues"
- Zero findings classified actionable across two rounds — you are validating, not doubting
- A geometry claim verified by a jsdom test
- Looping past three rounds without escalating
- Fixing a finding in the same pass that found it

## Verification

- [ ] Every scenario in the change's spec was named as a claim
- [ ] Each claim got a fresh-context adversarial review with ARTIFACT + CONTRACT only
- [ ] Findings were classified against the artifact, not rubber-stamped
- [ ] Evidence matched the claim: the mode that can actually observe it
- [ ] A stop condition was met — trivial findings, three cycles, or explicit ship
