# Context safety

How to keep an agent's working context small enough that it actually follows every instruction.

## The problem

A language model has a fixed context window. As the prompt grows, three things degrade:

1. **Instruction adherence** — rules in the middle of a long prompt are followed less reliably than rules at the start or end.
2. **Precision** — the model summarises rather than quoting, rounds rather than counting, and substitutes "close enough" for "exact".
3. **Cost and latency** — every token in the prompt is billed and processed, whether or not it contributes to the answer.

Context safety is the discipline of sending an agent the minimum it needs — and nothing else.

---

## 1. The minimum context set

Every agent prompt has exactly four required sections:

| Section | Contains | Does NOT contain |
|---|---|---|
| **Identity** | What the agent is, its role, its quality standard | The history of how the project got here |
| **Scope** | The files it may read and write, as globs | Every file in the repo |
| **Contract** | The types, schemas, or interfaces it must conform to | The implementation of those types elsewhere |
| **Deliverable** | What "done" looks like — files, commands, evidence | A tutorial on how to write the code |

Everything else is noise. If the agent needs a reference document, link to it by path — the agent can read it. Do not paste it into the prompt.

---

## 2. What to cut

### Full file contents you could link

```
# Bad — 200 lines of schema pasted into prompt
Here is the full schema:
```typescript
export const ProjectSchema = z.object({ ... })
```

# Good — 3 lines, agent reads if needed
The project schema is at `packages/core/src/schemas/project.ts`.
Your output must conform to `ProjectResponseSchema` exported from the same file.
```

### History and rationale

The agent does not need to know *why* a decision was made. It needs to know *what* the decision is. Link to the ADR if it needs the reasoning.

### Other agents' deliverables

If agent B needs the output of agent A, pass the **file paths** agent A wrote — not the file contents. Agent B reads them.

### Unused skill text

If the agent is building a backend service, do not include the UI development skill. Skills are loaded on demand, not bundled into every prompt.

---

## 3. The 4K rule of thumb

A well-scoped agent prompt is under 4,000 tokens. If your prompt exceeds this:

1. **Check for pasted file contents.** Replace with paths.
2. **Check for background context.** Move to a linked document.
3. **Check for instructions the agent won't use.** Remove them.
4. **Check for repeated instructions.** Say it once.

This is a guideline, not a hard limit. Some agents (reviewers, security auditors) legitimately need more context. But if a *builder* agent needs more than 4K tokens of instruction, the task is probably too large for one agent.

---

## 4. Context pollution patterns

| Pattern | Symptom | Fix |
|---|---|---|
| **Kitchen-sink prompt** | Agent ignores rules in the middle | Cut to the minimum context set |
| **Copy-paste escalation** | Each iteration pastes the previous output back in | Pass file paths, not contents |
| **Defensive over-specification** | Prompt lists every edge case because the agent missed one | Fix the one case; don't add 20 rules to prevent it |
| **Conversation accumulation** | Multi-turn agent accumulates all prior turns | Start a fresh agent for each phase; pass only the handover |
| **Skill stacking** | Agent receives 5 skills when it needs 1 | Load only the skill that matches the task |

---

## 5. When to start a fresh agent

Start a new agent (not a follow-up message to the same one) when:

- The task changes character (building → reviewing → testing).
- The accumulated conversation exceeds ~30 turns.
- The agent starts making mistakes it did not make earlier in the conversation.
- You need the agent to follow a different skill's instructions.

A fresh agent with a clean, minimal prompt outperforms a tired agent with a long, noisy context.

---

## 6. Security note

Context safety is also a security boundary. An agent that receives the full codebase can leak secrets it was never meant to see. Scope the readable paths, and never paste credentials, API keys, or customer data into an agent prompt — not even to "show an example".

---

## Related

- [Orchestration protocol](orchestration-protocol.md) — the fan-out contract that uses minimal context
- [Failure journal format](failure-journal-format.md) — how to record what went wrong without bloating context
- [skill-format.md](skill-format.md) — how skills are structured to load on demand
