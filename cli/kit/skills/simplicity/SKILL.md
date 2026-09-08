---
name: simplicity
description: "Choosing the simplest shape that fully solves the problem — the decision you make before you write, not the cleanup you do after. Load this when starting any change and deciding how much structure it needs; when tempted to add an abstraction, a config option, a flag, a wrapper, a layer, an interface or a 'manager'; when a design review asks 'is this over-engineered?' or 'is this doing too much?'; when a change feels heavier than the problem; when reaching for a pattern because it is familiar rather than because the problem asked for it; when deciding whether to build for a future requirement that has not arrived; or when auditing existing code for accidental complexity. Owns the simplicity ladder (take the lowest rung that fully solves the real problem), the non-lazy boundary (simple is the whole problem in a minimal shape, not part of the problem), the `// shortcut:` convention for a deliberately-lower rung, and the accidental-complexity hunt list. Routes to refactor for the loop that removes complexity safely once it exists, to review for the over-engineering axis at the gate, to `../../docs/philosophy.md` for the compound-engineering principles this serves, and to state-management for where a value lives."
---

# Choosing the simplest sufficient shape

**The defining constraint: the lowest rung that fully solves the actual problem.** Not the lowest
rung — that under-builds. Not the rung that would be nice if the problem were bigger — that
over-builds. The one where the whole problem is solved and nothing above it is carried.

This skill is the decision you make **before** you write. `refactor` owns removing complexity that
already exists; `review` owns catching it at the gate. This owns not creating it in the first place,
which is the cheapest place to win. The standard it serves is `../../docs/philosophy.md` — *the
codebase gets simpler as the product gets larger* — and the five principles there are the law; this
skill is how you apply them at authoring time. Do not restate them here; read them.

## Simple is not small, and it is not lazy

Two different mistakes wear the word "simple", and telling them apart is most of the job.

- **Over-building** adds structure the problem did not ask for: an interface with one implementer, a
  config option nobody sets, a flag with no removal condition, a layer that only forwards. It looks
  like foresight. It is usually a guess about a future that has not arrived.
- **Under-building** ships part of the problem and calls the rest out of scope: no empty state, no
  error branch, no authorization check, no handling for the input the real caller will actually
  pass. It looks like simplicity. It is an unfinished change, and the missing part becomes a defect
  someone with less context pays for.

**The test for the non-lazy boundary:** does the code handle every input the real caller will pass —
including the empty, the error, the unauthorized, the concurrent? If yes, it is simple. If no, it is
not simpler than the full version, it is *less than* the full version, and the difference is a bug
you have not written down. Simplicity removes structure the problem does not need; it never removes
the problem.

## The ladder

Structure is a cost paid in names, files, indirection, and places to debug. Each rung costs more
than the one below. **Take the lowest rung that fully solves the real problem, and climb only by
naming what the next rung deletes.**

| Rung | The shape | Earns its place when |
|---|---|---|
| 0 | **Nothing** — the case does not need solving yet | The requirement is imagined, not asked for. Defer or delete it (YAGNI) |
| 1 | **Inline** — a few lines where they are used, no new name | Used once, read in place, unlikely to be reused |
| 2 | **A named local** — a function or constant in the same file | The inline version needs a name to be readable, but only this file uses it |
| 3 | **A shared function or module** | A *second* caller genuinely needs the *same* behaviour — not a similar-looking one |
| 4 | **A parameter or option on the existing thing** | The variation is one axis on something that already exists — extend it, do not fork it |
| 5 | **A new abstraction or interface** | It **replaces** three existing shapes and deletes them (principle 2 and 3). Two is a coincidence |
| 6 | **A new subsystem, config surface, or plugin seam** | It removes a whole *class* of future change, and you can name the changes it removes |

**Climbing is the claim, not the default.** "I might need it later" is rung 0's argument for
staying at rung 0. The move from rung *n* to *n+1* is only justified by what disappears: an
abstraction that sits *beside* the code it was meant to replace (principle 2) has not earned rung 5,
it has added a rung-5 cost for a rung-1 benefit.

**Descending is always legitimate.** If the problem turned out smaller than the shape you reached
for, drop a rung. That is not a refactor to schedule; it is the change getting simpler as you
understand it.

## The `// shortcut:` convention — a deliberate lower rung, written down

Sometimes you knowingly take rung *n* when the problem will need *n+1* soon — the third case has not
arrived, the config dimension is coming next sprint, the extraction is real but out of this change's
scope. That is fine, and it is the *only* honest form of "do the simple thing now", because it is
bounded (principle 5). Silent under-building is the failure; a recorded shortcut is not.

```
// shortcut: <what is simplified, in one line>
// promote to <the higher rung> when <the concrete condition that triggers it>
```

The condition must be observable by someone who is not you: *"when a third provider is added"*,
*"when this option needs a non-default value"* — not *"when we have time"*. A shortcut with no
promotion condition is not a shortcut, it is the permanent shape, and it should be the honest lower
rung with no comment at all.

When the shortcut is more than a one-liner — a whole seam deferred, a boundary knowingly crossed —
it belongs in `refactor`'s debt ledger as a `RECORDED` entry with numbers, not only in a comment.
Route it there; do not let a paragraph of debt live as a comment.

## The accidental-complexity hunt list

When auditing existing code — or your own change before you open it — these are the greppable
signatures of a rung that was climbed without earning it. Each has a fixed answer.

| Signature | How to find it | The move |
|---|---|---|
| An interface with one implementer | The type has exactly one `implements`/one concrete return | Inline it; the abstraction is rung 5 paying for a rung-2 problem |
| A wrapper that only forwards | A function whose body is a single call passing the same arguments through | Delete it; call the wrapped thing |
| A config option nobody sets | Grep the option name → one definition, every caller uses the default | Delete the option; hard-code the default |
| A flag with no removal condition | Grep the flag → no `promote to`/`remove when`/issue link near it | Ask what removes it; no answer means delete the branch, not the flag |
| A parameter always passed one value | Grep the call sites → the argument never varies | Drop the parameter; inline the constant |
| A `-manager` / `-helper` / `-utils` grab-bag | A module whose functions share only a suffix | Move each function to the concept that owns it; delete the bag |
| Two names for one concept | The same shape appears under two identifiers | Pick one; principle 3's "shared helper beside three copies" in miniature |
| A dead branch | A condition whose other side is unreachable given the callers | Delete the branch and the condition |

This list overlaps `refactor`'s MEASURE table and `review`'s over-engineering symptoms on purpose —
same disease, three moments to catch it: before writing (here), while removing (`refactor`), at the
gate (`review`). It is not the full smell catalogue; that is `../refactor/references/smells.md`.

## When simplicity is the wrong goal

Naming these is better work than a plausible-looking minimal shape:

- **Safety, money, permissions, migrations, crypto** — principle 2's stated reversal: here the
  abstraction may need to exist *before* the duplication, because the invalid state must be
  impossible, not merely unlikely. Under-building these is not simplicity, it is a defect class.
- **The domain already has a stable boundary.** Abstracting on the first case is right when the shape
  is not a guess — a well-known protocol, an established seam. Rung 5 on case one is justified by the
  boundary, not by the count.
- **The simplest local shape pushes complexity uphill.** A rung-1 inline that forces every caller to
  know a historical accident (principle 1) is not simpler, it has moved the cost to more places.
  Count the total structure, not the structure in this file.

## Before you finish

- [ ] The shape is the lowest rung that handles every input the real caller will pass
- [ ] Every rung above the minimum is justified by something it deletes, named in the commit
- [ ] Any deliberate lower rung has a `// shortcut:` with an observable promotion condition, or is in the ledger
- [ ] No config option, flag, parameter, or interface with a single non-default user survived the hunt list
- [ ] Nothing safety-, money-, or permission-critical was made "simpler" by dropping a case

## Read next

- `../../docs/philosophy.md` — the five compound-engineering principles this skill applies.
- [`refactor`](../refactor/SKILL.md) — the MEASURE → PIN → MOVE → PROVE → RECORD loop for removing
  complexity that already exists, and the debt ledger a shortcut routes to.
- [`review`](../review/SKILL.md) — the over-engineering axis at the merge gate, and the
  rationalizations table.
- [`state-management`](../state-management/SKILL.md) — the placement ladder, which is this skill's
  logic applied to where a value lives.
