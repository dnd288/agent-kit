# Verification

What actually happened, not what was expected to happen. A fix with no RED is a change that
might have done nothing.

## The failing test

<!-- Its path and name, and the mode that can observe the claim — cross-cutting-tdd owns the mapping:
     jsdom cannot see geometry, a unit test cannot see a route. Written BEFORE the fix.

     If the defect was user-visible, this is the scenario from flows.md, in
     your end-to-end feature directorydefects/<symptom>.feature, and the command that selects it:
     `your end-to-end test command defects/<symptom>`. -->

## RED — what it asserted before the fix

<!-- The actual failure output. Paste it. "It failed" is not evidence that it failed for the
     right reason, and a test that passes before the fix is testing something else. -->

## GREEN — after the fix

<!-- The same test passing, pasted, plus whatever else was run. -->

## What ran, and what did not

<!-- By name. A run that was started and killed is UNRUN. -->

## What is still not proven

<!-- Honestly. If the fix is verified in jsdom but the symptom was reported in a browser,
     say so. -->
