# Compound engineering principles

Compound engineering is the practice of making each change lower the cost of the next one. A codebase compounds when growth mostly adds data, records, small modules, and well-owned behaviours. It decays when growth adds branches, flags, duplicated flows, and exceptions that only one person remembers.

These principles are deliberately written as review standards rather than mechanical rules. They ask whether the change made future work cheaper.

## 1. Leave it easier to change than you found it

The standard is not whether the code is clever, modern, or locally tidy. The standard is whether the next person can make the next adjacent change faster and with less risk.

A change leaves code easier to change when it:

- Removes a hidden assumption.
- Names a concept that was previously implicit.
- Moves behaviour to the layer that owns it.
- Adds a test at the level where the behaviour can break.
- Deletes obsolete branches, comments, flags, or docs.
- Narrows an interface to the values callers actually need.

A change makes code harder to change when it:

- Adds a second way to do the same thing.
- Solves one case by hiding a condition in a distant layer.
- Requires future callers to know historical accidents.
- Leaves documentation describing the old world.
- Adds configuration whose valid combinations nobody can enumerate.

### Examples

A form needs a new field. The compounding change adds the field to the schema, service, UI, validation, tests, and docs in the places that already own fields. The decaying change adds a special form variant because this one customer is different.

A report needs another grouping. The compounding change models grouping as data and updates one renderer. The decaying change copies the report and changes three labels.

An endpoint has repeated authorization checks. The compounding change extracts the policy into a named service and replaces the repeated checks. The decaying change adds a helper while leaving the old checks scattered.

### What would reverse this

This rule would be softened only if the cost of making the nearby future cheaper consistently exceeded the cost of localized duplication. Evidence would need to show repeated cases where cleanup blocked urgent delivery and the duplicated code stayed isolated, obvious, and cheap to remove.

## 2. A new abstraction earns its place by deleting

An abstraction is not free. It adds a name, an interface, a file, a concept, and a place to debug. It earns that cost by removing more complexity than it adds.

Good abstractions replace repeated policy or structure. They make invalid states harder to express. They shrink callers. They let the code say the product concept once.

Weak abstractions sit beside the old code. They add a wrapper, but callers still need to know the wrapped details. They centralize syntax while leaving policy duplicated. They make tests mock the abstraction instead of proving the behaviour.

Before adding an abstraction, ask:

- What code disappears in this change?
- Which future change becomes smaller?
- Which decisions move behind the interface?
- Which invalid use becomes impossible or obvious?
- What would make this abstraction obsolete?

### Examples

Three routes shape the same error body. A useful abstraction makes every route throw a typed error and deletes local response shaping. A weak abstraction adds `formatError()` while every route still chooses status codes and body fields itself.

Several components repeat the same loading, empty, and error states. A useful abstraction creates a state component and replaces all repeated markup. A weak abstraction exports constants while the branches remain copied.

Two payment providers need integration. Two may be a coincidence. When a third arrives, a provider interface earns its place if it replaces the three custom flows. It does not earn its place if each provider still has a bespoke path plus an adapter.

### What would reverse this

This rule would be reversed for domains where abstraction must exist before duplication for safety or compliance: cryptography, permission policy, money movement, migrations, or legal record keeping. Even then, the abstraction must record the risk that justified it before deletion was possible.

## 3. Two is a coincidence; three is a missing abstraction

The first implementation teaches the shape of the problem. The second tests whether the shape repeats. The third is evidence that the codebase is missing a concept.

Do not abstract on the first case unless the domain already has a stable boundary. Do not ignore the third case unless the similarities are superficial. At three, either create the abstraction or write down why the three cases must remain separate.

The abstraction must replace all three. Leaving three copies plus a shared helper creates four places to understand.

### Examples

One import wizard can be custom. Two import wizards may still differ enough to stay separate. A third import wizard probably means there is an import workflow: parse, validate, preview, commit, report. The abstraction should own that sequence and let each importer provide data-specific parts.

One notification channel can be direct. Two channels can be separate while the behaviour is still settling. Three channels usually mean there is a delivery interface with shared retry, templates, audit logging, and failure reporting.

Three tables with copied pagination logic mean pagination is a component, hook, or query model. The fix is not a shared `PAGE_SIZE` constant while each table still implements navigation differently.

### What would reverse this

This rule would be weakened when the three cases are likely to diverge faster than they converge, or when abstraction would freeze an unstable product question. The exception must be explicit and dated: what evidence would show convergence, and when the team will look again.

## 4. Delete the special case; do not flag it

Flags and special cases are sometimes necessary for rollout, experiments, compliance, or emergency mitigation. They become dangerous when they outlive the reason they were added.

A special case multiplies the states maintainers must reason about. A flag creates at least two products: on and off. Two flags create four combinations. Many flags make the code impossible to read in one pass.

Prefer changes that make the general rule true. If one customer, route, locale, or provider needs different behaviour, ask whether the product has discovered a real dimension. If it has, model that dimension as data. If it has not, keep the exception small, named, owned, and temporary.

### Examples

A customer needs a different approval limit. The compounding change stores approval limits per account or plan. The decaying change adds `if (customerId === ...)`.

A migration needs old and new behaviour for one deploy. The healthy flag has an owner, removal condition, and deletion task. The unhealthy flag remains for years and every new feature must support both paths.

A region requires different tax handling. That is probably a jurisdiction model, not a region flag in checkout code.

### What would reverse this

This rule would be relaxed in systems where runtime switching is the product: experimentation platforms, plugin hosts, compatibility layers, or feature management tools. In those systems, flags are first-class data with ownership, expiry, audit, and test coverage. The rule still applies to unowned branches hidden in product code.

## 5. An exception is written down and bounded, or it becomes the rule

Unwritten exceptions become folklore. Folklore becomes architecture. A future maintainer sees the exception, copies it, and reasonably concludes it is allowed.

Every exception needs:

- The rule it breaks.
- The reason it exists.
- The owner or team responsible for removing it.
- The scope where it is allowed.
- The date or condition for review.
- The condition that would remove it.

The record can live in an ADR, a spec, a code comment next to the exception, or a tracked issue. The right home depends on the breadth of the exception. A cross-system exception belongs in docs. A one-line temporary workaround can be a comment with an issue link.

### Examples

A route handler contains business logic during an incident. The exception comment names the incident, states that only this route may do it, links the follow-up service extraction, and gives the removal condition.

A package imports across a forbidden boundary for a migration. The ADR records the boundary break, why the migration needs it, which files may do it, and the release after which the import must disappear.

A raw colour appears in a component because the design token does not exist yet. The comment links the design-system issue and says the value is temporary until the token is added.

### What would reverse this

This rule would be reversed only in throwaway code that will not be maintained: experiments, spikes, or prototypes with a deletion date. Once code enters the maintained product path, exceptions need bounds.

## Enforced by review, not by guard

Static guards are excellent for facts:

- This file imports a forbidden package.
- This component contains a raw hex colour.
- This route lacks a response schema.
- This migration edits an existing migration file.
- This dependency is disallowed.

Compound engineering is a judgement:

- Did this abstraction remove more than it added?
- Is this duplicated code waiting for a third case, or already a missing concept?
- Is this flag a bounded rollout mechanism or a permanent fork?
- Did this exception deserve its cost?

A guard that tries to answer those questions will either miss the point or produce false positives. False positives are corrosive: once people stop trusting a guard, they disable it or route around it, and its real findings disappear too.

Use guards to protect crisp invariants. Use review to protect design judgement. The review should name the symptom, explain the cost, and ask for a smaller shape, deletion, or a written exception.

## Applying the principles

When reviewing or authoring a change, ask:

1. What concept did this change reveal?
2. Where does that concept belong?
3. What code became unnecessary?
4. What branch, flag, or exception did we avoid?
5. What future change is now cheaper?
6. What remains exceptional, and where is it bounded?

A change does not need to solve the whole codebase. It does need to avoid making the next honest change more expensive.
