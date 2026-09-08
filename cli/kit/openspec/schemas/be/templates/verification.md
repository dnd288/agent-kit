# Verification

What actually happened, not what was expected to happen. This artifact is the difference
between a feature that is proven and a feature that is assumed.

## The flow

<!-- The feature file and scenarios from flows.md, and the command that selects them:
     `your end-to-end test command <split>/<name>`. A positional filter is a regular expression
     over the generated path, so a directory or a feature name is enough. -->

## RED — before the change

<!-- Pasted output. "It failed" is not evidence that it failed for the right reason: a
     scenario that would fail with this change already in place is asserting something else.

     For a PENDING flow the red is `bddgen` accepting every sentence and the @known-issue tag
     skipping the run. Say that in those words. A skip reported as a pass is the failure this
     artifact exists to prevent. -->

## GREEN — after the change

<!-- The same selection passing, pasted, with the count. Plus `the project validation command`, and this
     layer's own tests. -->

## What ran, and what did not

<!-- By name. The suite whole is CI's job, behind the `e2e-checks` label; a local run is a
     slice, and which slice is a claim about coverage.

     **A run that was started and killed is UNRUN.** One change had to record its end-to-end
     layer as unrun after forty minutes produced no pass count, no failure list and no
     verdict. Saying so is the difference between a gap and a lie. -->

## What is still not proven

<!-- Honestly. A claim verified in jsdom for a symptom reported in a browser is a claim with
     a gap, and the gap belongs in writing rather than in somebody's memory. -->
