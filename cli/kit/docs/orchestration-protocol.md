# Orchestration protocol

How to coordinate multiple agents working on the same codebase without stepping on each other.

## The problem

When a feature workflow fans out — a UI agent, a backend agent, and a test agent all running in parallel — three failure modes appear:

1. **File collision** — two agents edit the same file at the same time, and one write overwrites the other.
2. **Semantic collision** — two agents create compatible files that make incompatible assumptions (different type shapes, different naming, different error handling).
3. **Context blowout** — an agent receives a prompt so large it loses the instruction that matters.

The protocol below prevents all three.

---

## 1. Fan-out contract

A coordinating agent (the "orchestrator") assigns work to sub-agents. Each assignment carries:

| Field | Required | Purpose |
|---|---|---|
| `scope` | yes | The directories and files the agent may write. Exclusive — no two agents share a path. |
| `contract` | yes | The types, schemas, or interfaces the agent must conform to. Defined before fan-out, not during. |
| `deliverable` | yes | What "done" looks like — the files that must exist and the command that must pass. |
| `context` | yes | The minimum the agent needs to do its job — no more. See [Context safety](context-safety.md). |

### Scope rules

- A scope is a list of glob patterns (`apps/api/src/services/project/**`).
- Two scopes must not overlap. If they must (e.g. a shared type file), the type file is created **before fan-out** and is read-only to all agents.
- An agent that needs to touch a file outside its scope **stops and reports** rather than writing.

### Contract rules

- The contract is a concrete artifact: a TypeScript interface, a zod schema, an OpenAPI fragment, a function signature.
- It is written by the orchestrator (or a dedicated "contract agent") **before** any implementation agent starts.
- An implementation agent may not change the contract. If the contract is wrong, the agent stops and reports.

---

## 2. Handover format

When an agent finishes, it returns a structured result:

```
FILES_WRITTEN:
  - path/to/file.ts
  - path/to/file.test.ts

FILES_READ (not written):
  - path/to/contract.ts

COMMANDS_RUN:
  - pnpm typecheck (passed)
  - pnpm test --filter @scope/package (passed)

ISSUES:
  - none | list of unresolved items
```

The orchestrator collects these before starting the next phase. A phase does not start until every agent in the previous phase has reported.

---

## 3. Phase ordering

The default phase sequence for a feature:

```
Phase 0: Contract
  → Define shared types, schemas, interfaces
  → These become read-only inputs to all later phases

Phase 1: Backend + UI (parallel)
  → Backend agent: services, routes, migrations
  → UI agent: components, screens, stories
  → Scopes are disjoint by definition (different packages)

Phase 2: Frontend wiring
  → Connects UI to backend through the application layer
  → Reads Phase 0 contracts, Phase 1 deliverables

Phase 3: Integration + E2E
  → Proves the whole thing works end to end
  → May read everything, writes only test files
```

Phases are not mandatory. A small change may have one phase. The discipline is: **define scope and contract before fan-out, collect results before the next phase.**

---

## 4. Failure protocol

When an agent hits an obstacle:

1. **Stop writing.** Do not guess past the obstacle.
2. **Report what you know.** The symptom, the file, the line, and what you tried.
3. **Report what you need.** The decision, the missing information, or the contract change.
4. **Return control.** The orchestrator decides the next step — retry, reassign, or escalate to the user.

An agent that guesses past an obstacle creates a semantic collision. The cost of stopping is one round trip. The cost of guessing is debugging an assumption that was never stated.

---

## 5. Anti-patterns

| Pattern | Why it fails | Fix |
|---|---|---|
| Shared mutable file | Two agents race on the same file | Move to Phase 0 contract or serialize access |
| Implicit contract | Agents infer types from each other's code | Write the contract explicitly before fan-out |
| Full-repo context | Agent receives the entire codebase | Send only scope + contract + deliverable |
| Phase skip | Integration starts before backend reports | Gate phases on handover reports |
| Silent override | Agent changes something outside scope without reporting | Scope check in the orchestrator |

---

## Related

- [Context safety](context-safety.md) — how to keep agent prompts small enough to work
- [Failure journal format](failure-journal-format.md) — how to record what went wrong and what fixed it
- [philosophy.md](philosophy.md) — the compound engineering principles that inform this protocol
