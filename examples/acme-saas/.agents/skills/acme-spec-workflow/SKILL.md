---
name: spec-workflow
description: "The specification-driven loop, end to end and END-TO-END FIRST — pick the change type, then map the business logic to flows, specify, design, tasks, apply, verify, archive. Load this when starting or continuing any user-visible work: a new feature, a screen, a cross-workspace change, an openspec change folder, a proposal.md, a flows.md, an interface.md, a capability spec, build.json, verification.md, or the decision artifact. Also load it when deciding whether work is one change or several, what schema a change should use (ui, fe, be, integration, bugfix or cross-cutting), what a task needs to be hand-off ready for a subagent, when a change is blocked by an open question, or when archiving. Owns the type-selection rule, the end-to-end-first order, the gates within a typed change, the build.json task manifest, the subagent handover contract, and the per-type surface table. Routes to your specification tool folder/README.md and config.yaml for mechanics, docs/product/prd.md for FR- identifiers, docs/product/open-questions.md for what gates a change, your project's architectural decisions for the workflow and your project's architectural decisions for why changes are split by layer."
---

# The specification workflow

Anything user-visible, or anything spanning more than one workspace, gets a specification
before it gets code — `AGENTS.md` says so, `CONTRIBUTING.md` step 2 says so, and since
your project's architectural decisions the tooling exists to make it
true. This skill is the working procedure. The decision record is the why; `your specification tool folder/README.md`
is the map; `your specification tool folder/config.yaml` carries the injected context and per-artifact rules. This
file owns the sequence and what goes wrong in it.

Copy changes, and defects whose cause is obvious from the symptom, skip all of this and go
straight to a pull request. **A defect whose cause is not obvious gets a `acme-bugfix` change** —
see `acme-bugfix`.

## First: pick the type

A change has a **type**, chosen at creation and pinned into `changes/<id>/.openspec.yaml`.
It decides which artifacts the change owns and what proves it.

| Type | Owns | Proven by | Its flow, usually |
|---|---|---|---|
| `ui` | your component library path, your UI package path — props-only components and screens | stories · quality | PENDING — no route, no data |
| `fe` | `your app path/src/app` — routes, actions, DTO mapping, state placement | unit, jsdom | the first real RED |
| `be` | your API service path, your data layer path, your shared contract package path | contract · service · route tests, one migration | PENDING — no screen |
| `integration` | the wiring of the above, and the app build | `your e2e path/features/**/*.feature` with steps behind them, one Chromium project | the pending ones, turned green |
| `acme-bugfix` | a defect at any layer | inspect · debug · fix · verify | `features/defects/<symptom>.feature` |
| `cross-cutting` | genuinely cross-cutting work — infrastructure, deployment, documentation consistency | whatever the change names | green here, or not user-visible |

**Every type carries `flows.md` and `verification.md`.** The end-to-end map and its pasted proof
are the two artifacts no type is exempt from — the rest of the graph is what differs.

```
openspec new change <id> --schema be
```

**Pass `--schema` every time.** Omitted, it silently falls back to `config.yaml`'s default and
the change is pinned to it — `your specification tool folder/README.md` records this biting once already, when a
broken config left a change on the packaged schema and nothing failed.

**A feature is a sequence, not one change.** `ui` and `be` run in parallel; `fe` wires them;
`integration` proves the result end to end:

```
ui  ─┐
     ├─→ fe ─→ integration
be  ─┘
```

**Why `ui` and `be` are parallel rather than sequential** is the same argument that used to say
"UI first". `your UI package` composites are props-only — no `next/*`, no `your API package`, no `your shared contract package`, no `fetch` —
so a screen's prop types *are* the specification of what the API must return. Drafting them from
the design means the contract falls out of the design instead of the design bending to a guessed
schema, and the pixel and visual layers exist before data can perturb them, which is the only
reason they mean anything. Neither change waits on the other: the `ui` change states the contract
shapes its caller will map in, the `be` change states which consumer it is satisfying, and both
write that down in `interface.md`.

## End-to-end first, within every type

**The order inside a change is flows → tests → code → pasted proof.** `flows.md` is written
immediately after the proposal and before `interface.md`, the specs and the tasks; task 1.1 is
the flow and it is a barrier; `verification.md` closes the change with the run.

The argument is one sentence: **a flow written after the code it checks proves that the code
passes it, not that it would have caught the code being wrong.** Everything else here follows
from that.

**A layer that cannot turn its flow green still writes it.** A `ui` change has no route and a
`be` change has no screen, so the person in the scenario cannot reach the thing yet. The
scenarios go to `your e2e path/features/known-issues/`, tagged `@known-issue` with the blocker named in a
sentence — and they are still checked, because **a sentence with no step behind it fails
`bddgen`** before a browser opens. The `fe` change is where the RED becomes real; the
`integration` change moves the file to the directory `your e2e path/AGENTS.md` places it in and deletes the
tag. Nothing is deferred: what moves is the tag, not the writing.

```
ui  ─┐  flow PENDING ─┐
     ├─→ fe ─────────→ ├─→ integration: tag off, flow GREEN
be  ─┘  flow PENDING ─┘
```

**`known-issues/` is not a parking space for red tests.** A scenario goes there when the thing
it asserts is not built and the blocker fits in a sentence. A scenario that fails because the
application is *wrong* is a defect, and belongs to `acme-bugfix`.

**That `interface.md` is the seam.** `your shared contract package` and the state-placement rule are what the
three types meet at, and each states its own side: `ui` its props, `be` its schemas, `fe` the
mapping between them and where the state lands. Three changes that each name the others cannot
drift apart silently — which is the failure a single spanning change hides until the end.

**Split when it is genuinely one movement.** A change that touches one workspace and can be
verified on its own is the right size. If a proposal's "What Changes" names files in two layers,
that is the signal to split, not to widen.

## Source

Adapted from `spec-driven-development`, `interview-me` and `incremental-implementation` in
[`addyosmani/agent-skills`](https://github.com/addyosmani/agent-skills) (MIT). The generic variant is rewritten to route to your repository's own documents and specification artifact structure.

## The loop, and who owns what

| Phase | Artifact | Command | Rule owner |
|---|---|---|---|
| **Propose** | `your change tracking folder/<id>/proposal.md` | `openspec new change <id> --schema <type>` | `docs/product/prd.md` mints the `FR-*` it cites |
| **Flows** | `changes/<id>/flows.md` | — | `acme-e2e` — the business logic, as Gherkin, before anything else |
| **Interface** | `changes/<id>/interface.md` | — | the seam: contracts, props, state placement |
| **Specify** | `changes/<id>/specs/<capability>/spec.md` | `/opsx:propose` | `#### Scenario:` WHEN/THEN, with a `*Flow:*` line |
| **Design** | `changes/<id>/design.md` | — | only when the how is not obvious |
| **Migration** | `changes/<id>/migration.md` *(`be` only)* | — | one per PR, or an explicit "no" |
| **Tasks** | `changes/<id>/tasks.md` **+ `build.json`** | — | the flow first, then the schema's own ordering |
| **Apply** | code | `/opsx:apply`, or the `feature workflow` workflow | one subagent per task |
| **Verify** | `changes/<id>/verification.md` | `your validation command`, then `cd e2e && your e2e test command <slice>` | `acme-verify`, `acme-review` |
| **Archive** | the living spec **and a numbered ADR** | `openspec archive <id>` | the typed schemas |

`openspec list`, `openspec status --change <id>`, `openspec view` — say where a change is at
any point. The `/opsx:*` commands and `openspec-*` skills under `.claude/` are **generated** by
OpenSpec and never hand-edited; the CI `spec` job proves they are in step.

## The gates, inside one typed change

Gates are barriers: later work genuinely composes earlier work. Within a gate, work fans out
concurrently; between gates, nothing starts until the gate before it has drained.

| Gate | What it covers | Fan-out |
|---|---|---|
| **G1 Specify** | proposal · **flows** · interface · specs · design · tasks + `build.json` | serial — one author, one thread |
| **G2 Decompose** | *(`ui` only)* screen manifest per screen · assets | 1 per screen; 1 per asset |
| **G3 Flow** | the Gherkin and its steps, written and **run**, RED recorded | 1 per feature file; **a barrier — nothing in G4 starts until it drains** |
| **G4 Build** | the change's own layer, in the order its schema's `tasks` instruction gives | 1 per element; **barrels and the migration serial and single** |
| **G5 Prove** | the surfaces of this type · the flow re-run and GREEN pasted into `verification.md` · verify · review | 1 per spec; then serial |
| **G6 Archive** | the living spec **and a numbered ADR** | serial |

**G1 names what the change removes.** The proposal states what it deletes, reuses or unifies — or records that it removes nothing, and why that is right. This is the [golden rules](../../../AGENTS.md#golden-rules) gate at the only moment it is cheap: once the tasks are written, "what could this replace" is a rewrite rather than a sentence. It is one line in `proposal.md`, not a new artifact.

**These run per change, not per feature.** A `ui` change and a `be` change each walk all five,
concurrently and independently. That is the point: a change that spans four layers can only be
verified once every layer is finished, so it is verified last or not at all — which is exactly
what happened to `create-project-creation-and-listing`, whose one unverifiable task ("run the
three layers together") was the one that mattered.

G2 exists only for `ui`, because only a screen has a Figma frame to decompose. **G3 is the
barrier that makes the whole thing mean anything** — the flow is written and run before the layer
is built, and its outcome (a real RED, or a `bddgen`-checked skip on a `@known-issue` file) is
recorded before anybody writes a component. G4's ordering is the schema's own: `ui` runs tokens →
primitives → widgets → screens → barrels, `be` runs contracts → migration → services → plugins →
routes, `fe` runs mapping → page → action → leaf. The old single ordering that ran all of those in
one sequence described a change spanning all of them, and is kept only in the `cross-cutting` schema for the
cross-cutting case.

**Two waves are single-writer by design:**

- **Barrels** — `your component library path/src/index.ts`, `your UI package path/src/index.ts`,
  `your shared contract package path/src/schemas/index.ts`. Every new item appends to one, so parallel
  agents conflict here and nowhere else. One serial agent replays their reports into them,
  then `typecheck` proves it.
- **The your ORM migration** — one per pull request, backward compatible with the previous
  release. Two agents writing migrations produce two directories and a conflict the
  `migration-required` guard cannot untangle.

**`blocked` is a real outcome at every wave, not just the Figma ones.** A task gated by an
open question is kept in the list, marked with the question number — deleting it loses the
fact that it was scoped. The screen manifest's `blocked[]` already records this shape: the
element, the owner, the question, and an `unblocks` recipe. Carry it forward. A fully blocked
gate stops the run; it does not build half a feature and call it done.

## `build.json` — the task manifest

`tasks.md` stays OpenSpec's tracked checklist; the apply phase parses its `- [ ]` checkboxes
and nothing else. The detail a subagent needs cannot live on a checkbox line, so it goes in a
sibling `your change tracking folder/<id>/build.json` that `tasks.md` indexes by id — the same split as
`your resources package path/figma/screens/sign-in.json`, which carries eighteen elements each with a
verdict, a target and a paragraph of build detail.

```json
{
  "id": "3.2",
  "wave": "primitives",
  "title": "Add the `social` variant to Button",
  "skill": "component-development",
  "target": "your component library path/src/components/button/",
  "implements": ["FR-2.1-01"],
  "scenario": "sign-in › Agent signs in with email and password",
  "source": "your resources package path/figma/screens/sign-in.json#continue-with-google",
  "surfaces": ["component", "unit", "stories", "quality"],
  "conflicts": ["your component library path/src/index.ts"],
  "verify": "pnpm --filter your component library test && pnpm --filter your component library test:quality",
  "blockedBy": "OQ15 — OAuth providers not agreed"
}
```

**The subagent handover contract.** A task is well formed when its entry alone is enough to
build from — no conversation history, no sibling task's output. Six fields carry that:

- `skill` — what to load. Every builder loads its layer's skill and follows its sequence.
- `target` — where it lands.
- `scenario` — the spec sentence it must satisfy. If the change's `specs/**` has no scenario
  covering it, the task is the wrong shape; the scenario was written first.
- `source` — the design or contract it reads instead of asking. A builder never re-derives a
  fact that is already tracked in the manifest or the contract.
- `surfaces` — what must ship with it, from the table below.
- `acme-verify` — the command that proves it. Run by the verify gate, not by the builder.

`conflicts` names files the agent must **not** touch — the barrels and the migration, always,
plus anything a sibling task owns.

**`wave` is a closed vocabulary, and the workflow dispatches on the exact string.** A wave name
outside its type's list is a task the `feature workflow` workflow will not build — it says so
rather than skipping quietly, and refuses the run outright if *no* task matches, because
building nothing and reporting success is the worst shape a failure can take.

| Type | Waves, in order |
|---|---|
| `ui` | **`flows`** · `assets` · `primitives` · `widgets` · `screens` |
| `fe` | **`flows`** · `app-routes` · `server-actions` · `client-state` |
| `be` | **`flows`** · `contracts` · `migration` · `services` · `plugins` · `routes` |
| `integration` | `flows` |
| `acme-bugfix` | `acme-test` · `fix` — in that order, and the order is the discipline |
| `cross-cutting` | all of the above, in the cross-cutting sequence |

**`flows` is first for every type that has one, and it is a barrier rather than a wave that
happens to be first** — the same shape as `acme-bugfix`'s `acme-test` before `fix`, and for the same
reason. A `flows` task's `acme-verify` is the selection it runs: `cd e2e && your e2e test command <split>/<name>`.

Several `build.json` files written before this vocabulary settled use names from an older one
(`ui`, `acme-design`, `api`, `worker`, `infra`). They are not built by the workflow, and were not
before either — the difference is that it now says so.

## The surface table — rich surfaces, made concrete

A task is not done when the file compiles; it is done when its layer's set exists. The table
is the checklist; the skills cited for each layer own the detail.

**Read the rows your type owns, and only those** — plus the Flow row, which **every** type ships.
A `ui` change ships Flow and the two component rows; `be` ships Flow and Contract through Route;
`fe` ships Flow and App route through Slice; `integration` ships Flow alone. Checking a `ui`
change against the Route row is how a review produces findings nobody can act on.

| Layer | Type | Location | Must ship |
|---|---|---|---|
| Primitive | `ui` | `your component library path/src/components/<n>/` | `<n>.tsx` (cva table co-located) · `<n>.test.tsx` · `<n>.stories.tsx` enumerating every variant — *a control the a11y run never toggles is a variant never checked* · `<n>.quality.test.tsx` |
| Widget / screen | `ui` | `your UI package path/src/<feature>/<n>/` | the four above **plus** a `PseudoLocale` story. The pixel and visual layers are retired (your project's architectural decisions, your project's architectural decisions), so `<n>.quality.test.tsx` is the whole rendering surface — `component-has-surfaces` now asks both packages for it |
| Contract | `be` | `your shared contract package path/src/schemas/<n>.schema.ts` | `<N>Schema` · inferred `<N>` type · a case in `core.test.ts` · the barrel line (serial) |
| Migration | `be` | `your data layer path/prisma/migrations/<ts>_<n>/` | exactly one, backward compatible |
| Service | `be` | `your API service path/src/services/<n>.ts` | no framework types in its signature, enforced by `service-purity` · `<n>.test.ts` |
| Plugin | `be` | `your API service path/src/plugins/<n>.ts` | one concern · wrapped in `fastify-plugin` so decorators and hooks escape encapsulation · registered in `buildApp()` before its dependants · `<n>.test.ts` through `app.inject()` |
| Route | `be` | `your API service path/src/routes/<n>.ts` | handler · **response schema** · failures thrown as `AppError` rather than shaped by hand · an `app.inject()` test via `buildApp()` |
| App route | `fe` | `your app path/src/app/<group>/<n>/page.tsx` | page · DTO→props mapping · forwarded session cookie · links passed into slots · **state placed by your project's architectural decisions's ordered rule**, and a mutation that revalidates what it changed |
| Server Action | `fe` | `your app path/src/app/<group>/<n>/actions.ts` | `'use server'` — **async functions and nothing else**, or it throws at runtime from a build that typechecked · a unit test · revalidation |
| Client leaf | `fe` | `your app path/src/app/<group>/<n>/<n>-client.tsx` | `'use client'` at the leaf · `useActionState` for pending and error · a unit test |
| Slice | `fe` | `your app path/src/app/<group>/<n>/<n>.store.ts` | a plain function, unit tested with no React · scoped to the feature, never global |
| Flow | **every type** | `your e2e path/features/<split>/<n>.feature` + `your e2e path/steps/<domain>.steps.ts` | Gherkin the business reads, one `Scenario:` per `#### Scenario:` · `acme-test`/`expect` from `../fixtures/index.js` · a sentence with no step fails `bddgen` · written and RUN before the rows below it · PENDING for `ui` and `be`, under `known-issues/` with the blocker named · `acme-e2e` owns the layers below it |

One naming rule belongs here because it is not obvious: catalogue keys are `emailLabel`,
never a bare `password:` — a `password`-named key assigned a string literal trips the
`secret-literals` guard.

## Three invariants that make the fan-out work

Each with the failure it prevents. They come from the screen workflow, measured there, and
they transfer:

1. **Discovered facts are interpolated, never hardcoded.** A probe in the run's first step
   establishes what is actually available — Figma tiers, chromium, the manifest — and every
   later prompt embeds the result. A sentence baked into a prompt stays true-by-assertion
   long after the facts move.
2. **Workers write tests; they do not run suites.** Two self-verifying agents once launched
   vitest worker pools simultaneously and spent most of their tool calls in `pgrep` wait
   loops: ~10 minutes of judgement, ~50 minutes of waiting. Verification is serial, once,
   with the machine to itself.
3. **Read the filesystem, not a document's summary of it.** If `AGENTS.md` or a skill states
   a different component count from what is on disk, the disk is right and the drift is
   worth reporting — a stale inventory is how a component gets built twice.

## Common rationalizations

| Rationalization | Reality |
|---|---|
| "The screen is small, I'll skip the spec" | The spec is what makes the scenario-test link possible. Without it, every check verifies against the code's own history. |
| "The open question will answer itself by the time we build" | It has an owner and a severity for a reason. Name it in the proposal or you are assuming an answer the client has not given. |
| "The manifest already documents the screen" | The manifest is *what to build*; the spec is *what the system does*. They answer different questions, and the guard knows the difference. |
| "One agent can do the whole feature" | One agent can; it just cannot do it in parallel, and the fan-out exists to make the cost of a screen bearable. Serial is the fallback, not the plan. |
| "This change is small, I will fold the route into the UI change" | Then nothing can verify the UI change on its own, and the route arrives with no specification of its own. Two small changes review faster than one medium one. |
| "The integration change is just wiring, it needs no spec" | Its scenarios are the only ones that cross a layer boundary. Every other test in the repository passes with the seam broken. |
| "The flow belongs in the integration change, not here" | Then this layer is built against an imagined path and nobody finds out for three changes. Write it now, tag it `@known-issue`, and the integration change removes a tag instead of inventing a journey. |
| "The flow can't pass yet, so writing it is busywork" | `bddgen` refuses a sentence with no step behind it, so a pending flow already proves its own vocabulary exists. And the day the gap closes, the proof is written rather than negotiated. |
| "The tests pass, so it works" | The layer's tests pass. Whether a person can do the thing is a different claim, and `verification.md` is where it is answered with a pasted run rather than an inference. |
| "A blocked task is a failed task" | A blocked task with the question named is a correctly scoped task. Deleting it loses the fact that it was scoped. |

## Red flags

- A proposal with no `FR-`/`NFR-` identifier and no open-question number. It has not read the PRD.
- A scenario written with three hashes — it validates clean and asserts nothing
  (`openspec validate --all --strict` is the backstop, but the author is the first line).
- Two tasks in one change both touching a barrel or the migration.
- A task whose `acme-verify` command is missing. It cannot be handed over.
- A gate skipped because "it will be fine" — gates exist because later layers compose earlier ones.
- Archiving a change with no decision artifact when the change made a hard-to-reverse choice.
- A proposal whose "What Changes" names files in two layers. Split it; do not widen the type.
- An `integration` change that names no antecedents. It is a feature wearing the wrong type.
- A change whose `flows.md` is missing, or written after `tasks.md`. The order is the mechanism.
- A `flows.md` naming a selector, a URL, a status code or an id. That is a step, a page object or a spec — `your e2e path/AGENTS.md`'s five layers say which.
- A `verification.md` whose GREEN is pasted and whose RED is not. A test that was never seen failing is a test nobody has evidence works.
- A skip reported as a pass — a `@known-issue` scenario counted as coverage, or a `@generation` scenario that skipped because the stack was absent.
- A `ui` or `be` change claiming an end-to-end proof. It has no route or no screen; the claim is the sibling's.
- A `acme-bugfix` whose `diagnosis.md` carries no observation. That is a hypothesis, and a fix for a hypothesis is how the defect comes back under a different symptom.

## Verification

- [ ] The change was created with an explicit `--schema`, and `.openspec.yaml` names the type intended
- [ ] The change folder exists and `openspec status --change <id>` shows the artifacts in order
- [ ] `interface.md` names the sibling changes on the other side of the seam — a `ui` change with no consumer named, or a `be` change with no consumer named, is a guess about a screen nobody has drawn
- [ ] Every proposal cites the PRD identifiers it implements and the open questions that gate it
- [ ] `flows.md` exists, was written before `interface.md` and the specs, and names where each feature file lives among the seven directories
- [ ] Every flow this change cannot turn green is written anyway, under `known-issues/`, tagged `@known-issue` with the blocker in a sentence
- [ ] Every requirement has at least one `#### Scenario:` written so a failing test names it, and a `*Flow:*` line that either names the scenario or says the claim is layer-internal
- [ ] `verification.md` pastes the RED and the GREEN, names the selection that ran, and says what did not run
- [ ] `tasks.md` is ordered within this type's own layer; barrels and the migration are single-writer
- [ ] Every task in `build.json` passes the handover contract: skill, target, scenario, source, surfaces, verify
- [ ] Every blocked task names its question number and stays in the list
- [ ] `openspec validate --all --strict` passes on the change
- [ ] The archive step produced both the living spec and the numbered decision record
