---
name: review
description: "The review gate before merge — five axes, evidence over opinion. Load this when reviewing a change: your own, another developer's, or a subagent's, before a pull request is opened or when one is under review. Covers what your guard checks can and cannot see, the five axes in order of weight, the approval standard, and where each axis routes for the real rule. Load alongside security for any change touching auth, uploads, sharing or model input."
---

# Reviewing a project change

Every change gets reviewed before merge. Not a sprint review, not a skim — a five-axis pass
with the repository's own gates as the baseline. your guard checks have already run by the time a
change reaches you; your job starts where it stops.

## Source

Adapted from `code-review-and-quality` in
[`addyosmani/agent-skills`](https://github.com/addyosmani/agent-skills) (MIT), rewritten to
route to this repository's documents and to what its guard checks actually cover.

## What the guard already did, so you do not repeat it

Your guard checks — thirty-seven checks in two groups — have run on
push (your static validation) and in CI. It catches hex colours, wrong-package imports, a route with no
response schema, a `process.env` read absent from `.env.example`, a cookie missing
`httpOnly`, a workflow interpolating context into `run:`, a component past its complexity budget,
and a memo that can never hit. your engineering conventions
owns the structural rules and your security documentation the security ones; the guard is
their mechanical arm. **Do not re-verify what the guard verifies** — do verify that what the
guard *cannot* see is right, because that is the entire job:

- authorisation logic — the guard sees a session endpoint, not whether one user's project
  returns to another user
- whether a session is genuinely invalidated on sync-deactivation
- an over-broad response schema — declared and enforced are both checked; too generous is not
- whether a business rule lives in a service rather than a route handler
- whether the change satisfies the specification it was written against

## The five axes, in order of weight

### 1. Correctness — does it do what the spec says?

The change's your change tracking folder specs scenarios are the claim set. Each scenario is a
test case; a scenario with no test covering it is an unverified claim. `acme-verify` does the
adversarial pass against the spec; this axis is the confirming read — does the code, not the
test, satisfy the requirement?

### 2. Architecture — does it belong where it is?

The boundaries in your agent instructions are load-bearing and mostly guard-checked; the judgement cases
are not.

**Read your change configuration first and check the cases the change's type can actually commit.** A
`ui` change cannot put a your server framework type in a service signature, and telling its author to look
for one produces a finding nobody can act on — which is how a review starts being skimmed.

| Type | The judgement cases that apply |
|---|---|
| `ui` | a composite that wants data — it belongs in your app path, passed down as props; do not relax `ui-stays-presentational`. A store reached in from outside is the same violation as a `fetch`; a screen that rebuilds the shell instead of composing one; **a screen that has stopped being one component** — see below |
| `be` | a service signature with a your server framework type in it — wrong, and tests will fight you; business logic in a route handler — services are where the tests live; a response schema that names more than it needs; a new endpoint that reimplements a service instead of calling it |
| `fe` | a DTO from your API contract package leaking into your UI package prop types — the route maps at the boundary; state placed against your project architectural decisions' order; a Server Action file exporting anything but async functions; a container copied to a second route instead of lifted |
| `integration` | a defect patched at the seam instead of in the layer that owns it |
| `acme-bugfix` | a guard or a test edited to stop a failure; a fix with no RED output recorded |

**A change that trips a case outside its own row is usually mis-typed rather than wrong** — say
that, rather than asking for the code to move.

**Compound engineering** (your project golden rules) is this axis's other half, and the one with no guard behind it at all. The question is not whether the code is good but whether it left the next change cheaper. Five symptoms, each reproducible by pointing at the files:

- **a block copied rather than extracted.** The second copy is the moment to extract, not the fifth — the wizard shell reached nine files one reasonable commit at a time
- **a fallback constant standing in for a required prop.** `DEFAULT_USER` and its kin turn a caller's bug into a plausible-looking render; a required prop fails loudly instead
- **a third copy of a shape**, where the abstraction that would replace all three is sitting one file away
- **a flag added where a case should have been deleted.** Ask what removes the flag, and when — no answer means it is permanent
- **a new concept the feature did not need** — a layer, a wrapper, a config surface. Business growth is meant to be data growth

**A change whose whole purpose is this axis is a refactor**, and `acme-refactor` is the loop it should have followed — so the questions to ask of one are its: was the behaviour pinned before it moved, are structure and interface in separate commits, and did the measurement it was started for actually move.

**The question is "what did this delete", not "is this perfect".** A change that removes nothing is not automatically wrong: a genuinely new capability is genuinely new code. It is wrong when what it added already existed under another name.

**State placement** (your project architectural decisions, with the
identity-slice exception in your project architectural decisions) is the
newest judgement call and the guard cannot see it. Ask the rule in order — identity → database →
URL → query → context → slice → props — and look for the three failures it exists to prevent:

- **a store or a query cache holding a copy of a database row.** The answer was rule 1, and the copy
  is a bug waiting to disagree with its source (the session store is the one exception, and only
  for the identity the server resolved)
- **a Server Action that mutates and does not revalidate.** The stale render that follows reads
  exactly like a caching bug and is the most common defect in this shape
- **something in `useState` that should be in the URL** — a wizard step, a filter, an open dialog. If
  a reload or a shared link should reproduce it and does not, that is the finding
- **a session store that fetched or mutated the session itself** — the store is a projection; a
  `fetch` in the provider or a second session source is your project architectural decisions' reversal condition in the wild

A store or TanStack Query import reached in from OUTSIDE your UI package path or your component library path is the same
violation as a `fetch`: it gives a component an appearance its props do not describe, which is what
the pixel baseline then cannot assert. A screen's OWN per-mount view store, seeded from props and
holding no server data, is not that violation (your project architectural decisions) —
check the three constraints rather than the word.

**Component complexity** (your project architectural decisions) is
mechanically covered by `composite-complexity` — 24 props, 6 stateful hooks, 800 code lines,
complexity 45 — so the axis's job is the part a count cannot do:

- **a ledger entry added instead of a split.** A new line in `composite-complexity`'s `RECORDED` is
  an exemption wearing a debt record's clothes. The ledger is for what was already there
- **a split that drew the wrong seam.** Six widgets is not automatically better than one screen.
  Ask what any single future change now touches: if it touches four of the six, the seam is wrong
- **a screen store holding a server answer** — a product, a project, a generation. That is your project architectural decisions'
  stated reversal condition, and the value's home was rule 1
- **a widget reused across screens that reads one screen's context.** It has quietly become that
  screen's widget, and the second screen will discover it
- **a mock module used as a prop default.** `DEFAULT_RESULT_PRODUCTS` and `createMockResultState()`
  are imported by a shipped component today; the non-empty fallback is what put an "add product"
  control on the share recipient's read-only view

### 3. Readability — will this make sense at handover?

This repository is client-bound; the person reading this in six months may not have been here.
A comment explaining *why* (the load-bearing ones are in `theme.css`'s focus-visible rule and
the visual-testing ledger's predicates) is worth more than a comment restating what. A reviewer should be
able to say what a file does from its header and a skim.

### 4. Security — is the surface as documented?

your security documentation is the surface-by-surface map; `acme-security` is the
skill-length version. Read it before reviewing anything touching auth, sessions, cookies,
uploads, share links, model input, env vars, dependencies or workflows. The mechanical half
is guarded; the authorisation and over-broad-schema half is yours.

### 5. Performance — is anything unbounded?

The obvious ones: an N+1 your data layer query, a component that re-renders a whole screen because a
prop identity changed (`expectRenderCount` exists for exactly this), an image without
dimensions or `srcset` on the LCP element. Not micro-optimisation — unbounded things.

## The approval standard

Approve a change when it definitely improves overall code health, even if it is not perfect.
Do not block a change because it is not exactly how you would have written it. Block for:
the change does not satisfy its specification, it violates a stated rule, or it will cost a
future reader real time. The fix must be concrete — "this fails when X" beats "this feels
wrong". If the specification itself is the problem, say so; the spec is reviewable like
anything else.

**Green CI is not the whole gate on a UI change.** CI runs the jsdom modes only; stories,
quality, pixel, visual and e2e run on demand (your project architectural decisions,
amended). So for anything touching your component library path, your UI package path or your app path, ask where
those were run — the your CI labels for browser, e2e, or deep checks on the pull request, or the author's word that
your deep validation was green. A UI change whose deep modes nobody ran has been checked by
nothing that can see a layout.

**An e2e slice is evidence; it is not the suite.** A local run is expected to be narrowed
(`e2e/README.md`, *Running a slice*), so ask *which* slice and whether it covers the change — and
take the e2e CI label as the answer for the rest. "Started and killed" is **unrun**: a run
that reported no verdict is not weak evidence, it is none.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "The guard passed, it must be fine" | The guard checks declared and enforced. Authorisation and over-broad schemas are precisely what it cannot see. |
| "The tests pass, ship it" | Passing tests are evidence, not proof. Did a human read the diff? |
| "CI is green on this UI change" | CI no longer runs the browser modes. Green means jsdom passed; nothing has measured the layout unless someone ran the deep modes. |
| "I wrote it, so I'm the wrong reviewer" | A fresh-context read catches what the author's accumulated context hid — that is the point of reviewing before merge, not a reason to skip it. |
| "This is a small change" | Small is when review is cheapest. The skip-the-fan-out threshold in `your ship workflow` is two files and fifty lines with auth, uploads, sharing, model input and config untouched — anything above that gets reviewed. |
| "It's not exactly how I'd write it" | That is not a block. Correctness against the spec is; taste is not. |
| "We'll clean it up later" | Later is a person with less context and no reason to look. The change that created the duplication is the cheapest place to remove it, and the only place anyone is looking. |
| "It's only one more flag" | Every flag was only one more. Ask what removes it and when; no answer means permanent, and permanent flags are how a file stops being readable in one pass. |
| "Copying is faster than extracting" | For this change, yes. The bill arrives on the third copy, and it is paid by whoever has to change all three — usually not the author. |
| "The abstraction would be premature" | Premature abstraction guesses at a shape it has not seen. Extracting the third identical copy is not a guess; that rule exists to tell the two apart. |

## Red flags

- A review that only lists what the guard would have caught anyway
- An approval with no reference to the change's specification
- A blocked finding without a concrete failure scenario
- A security finding on an auth-touching change that did not open your security documentation
- "LGTM" on a change that grew a transaction and three branches inside a route handler

## Verification

- [ ] The change's spec scenarios were read, not just the diff
- [ ] Each axis was answered with evidence: the spec, the guard output, the security doc
- [ ] Every blocked finding names the failure it prevents
- [ ] The guard's blind spots were the focus, not its findings
- [ ] The verdict is one of: approve, approve with nits, or block-with-concrete-failure
