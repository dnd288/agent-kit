# Flows

The end-to-end features this change is proven by. Written BEFORE the interface, the specs
and the tasks — the business logic first, then the sentences that prove it, then the code.

## The business logic, mapped

<!-- One row per behaviour this change serves. The right-hand column is the whole point:
     if a failure here means nothing to the business, the claim is a layer-internal detail
     and belongs one layer down (cross-cutting-tdd places it). -->

| Business behaviour | Feature file | What a failure there means |
|---|---|---|
| <what the product must do> | `your end-to-end feature directory<split>/<name>.feature` | <the consequence, in the business's words> |

## `your end-to-end feature directory<split>/<name>.feature`

**Purpose**

<!-- The paragraph that becomes the `Feature:` preamble. Why this file exists, and what a
     red run here tells whoever is reading it. your end-to-end testing instructions's ordered question chooses the
     directory; first yes wins. -->

**Scenarios**

<!-- Gherkin, in the words the business uses. No selector, no URL, no status code, no id —
     a feature file naming one of those is in the wrong layer. -->

```gherkin
Scenario: <name>
  Given <the state somebody is in>
  When <what they do>
  Then <what they get>
```

**Tags**: <!-- @journey · @roles · @internal · @generation · @spends · @money · @known-issue.
     A tag is a cost declared in the scenario. your end-to-end testing documentation's table is what each one buys. -->

**Status when this change lands**: <!-- GREEN, or PENDING — and then the sibling change that
     turns it green, named. A pending flow is written NOW, in known-issues/, tagged
     @known-issue with the blocker in a sentence. It is still checked: a sentence with no
     step behind it fails `bddgen` before a browser starts. -->

**Steps**: <!-- Which of your end-to-end steps directory<domain>.steps.ts already says these sentences, and which
     are owed. A sentence reused is a sentence already proven to work. -->

## What this change cannot prove end to end

<!-- Required. Say it plainly, and name what pins it instead at the layer below — a quality
     test, a service test, an app.inject() test. A change that claims to prove everything
     end to end has not read its own boundary. -->
