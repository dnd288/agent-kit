# ADR format and conventions

An Architecture Decision Record (ADR) is a short document that records one technical decision, the context that made it necessary, the consequences of choosing it, and the evidence that would cause the team to change its mind.

ADRs are for decisions future maintainers would otherwise need to rediscover or re-litigate. They are not status reports, design dumps, or meeting minutes.

## When to write an ADR

Write an ADR when a decision is:

- Hard to reverse.
- Cross-cutting across packages, services, teams, or deployment surfaces.
- Non-obvious from the final code.
- A deliberate rejection of a common approach.
- A constraint that future work must obey.
- A bounded exception to a repository rule.

Do not write an ADR for every implementation detail. If the code explains the choice and the blast radius is local, a comment or test name may be enough.

## File naming

Store ADRs in `docs/adr/` unless the project uses another documented location.

Use this file name format:

```text
NNNN-short-title.md
```

Examples:

```text
0001-use-postgres-for-primary-storage.md
0002-serve-the-api-behind-the-web-app.md
0003-keep-generated-clients-checked-in.md
```

Rules:

- Use a four-digit sequence number.
- Use lowercase kebab case for the title.
- Do not renumber existing ADRs.
- If a decision is replaced, add a new ADR and update the old one's status.

## Required template

```md
# NNNN. Short decision title

## Status

Proposed | Accepted | Superseded by [NNNN](NNNN-new-title.md) | Deprecated

## Context

Describe the forces that made the decision necessary.

Include:

- The problem.
- Relevant constraints.
- Current behaviour.
- Options considered.
- Why the decision matters now.

Keep this factual. Link to specs, issues, incidents, benchmarks, or previous ADRs when they matter.

## Decision

State the decision directly.

Include:

- What will be true after this decision.
- Where the rule applies.
- Which option was chosen.
- Which important alternatives were rejected.

A reader should be able to quote this section as the rule.

## Consequences

Describe the results of the decision.

Include positive and negative consequences:

- What becomes simpler.
- What becomes harder.
- What new obligations exist.
- What migration or follow-up is required.
- What risks remain.

## What would reverse this

State the exit condition.

Include the evidence, product change, operational change, scale threshold, cost threshold, vendor change, or failure mode that would make this decision wrong.
```

## The reversal section is mandatory

Every decision carries its exit condition. "What would reverse this" prevents decisions from becoming permanent because nobody remembers why they were made.

Good reversal conditions are observable:

- "If p95 response time exceeds 300 ms for this endpoint under production traffic after caching, reconsider the storage choice."
- "If two more clients need offline writes, replace the server-only mutation model with a sync protocol."
- "If the vendor removes the pricing tier this depends on, reopen the build-versus-buy decision."
- "If this exception remains after the migration completes, delete it or write a replacement ADR."

Weak reversal conditions are vague:

- "If requirements change."
- "If this becomes a problem."
- "If we find a better way."

## Conventions

### One decision per record

An ADR should record one decision. If the team chooses a database, a queue, and an auth model in one meeting, write three ADRs when each decision has independent consequences.

### Link from specs

Specifications should link ADRs when the ADR explains a constraint that shapes the solution. ADRs should link specs when a product requirement caused the decision.

### Link from code sparingly

Use code comments to link ADRs only where the decision would otherwise look surprising. Do not decorate every implementation of the decision with a link.

### Update when reversed

When a decision changes:

1. Write a new ADR for the new decision.
2. Change the old ADR status to `Superseded by [NNNN](NNNN-new-title.md)`.
3. Update specs and docs that linked the old rule.
4. Remove or migrate code comments that cite the old rule.

Do not edit history to make the old ADR look correct. The old decision is useful context.

### Record rejected alternatives

Rejected options help future maintainers avoid repeating the same investigation. Keep the detail proportional. Name the option, why it was plausible, and why it lost.

### Keep ADRs short enough to read

Prefer two to four pages over a thesis. Link supporting benchmarks, incident writeups, diagrams, or vendor docs instead of embedding all evidence.

## Review checklist

Before accepting an ADR, check:

- [ ] The title states the decision, not the topic.
- [ ] The status is current.
- [ ] The context explains why the decision is needed now.
- [ ] The decision is direct and scoped.
- [ ] Consequences include trade-offs, not only benefits.
- [ ] "What would reverse this" is concrete and observable.
- [ ] Related specs, issues, or ADRs are linked.
- [ ] The decision is not duplicating another active ADR.
