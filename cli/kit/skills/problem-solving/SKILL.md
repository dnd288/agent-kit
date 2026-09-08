---
name: problem-solving
description: "Breaking an impasse — the techniques to apply when you are stuck, and the dev note you hand the user once you are through. Load this when repeated attempts at the same approach keep failing the same way; when a bug will not reproduce or resists diagnosis; when a design has no obvious path; when you notice you are thrashing — changing things and hoping. Owns the recognise-you-are-stuck signal, six techniques to regain traction (minimal repro, restate, bisect, one variable, invert, escalate), and the dev-note format — problem, solution, affects, implements, result — that is produced on the way out and handed to the user for review before the change is finalised. This is meta and cross-cutting: it applies inside any loop, whether the stuck thing is a bugfix, a refactor, or a design. Routes to bugfix for the defect loop itself, and to simplicity/refactor when the impasse is a design one."
---

# Breaking an impasse

**Recognising that you are stuck is the skill; the techniques are the easy part.** The signal is
concrete: two or three attempts at the *same shape* of approach, each failing the *same way*. That is
not progress needing one more push — it is a wrong assumption you cannot see, and the next attempt of
the same shape will fail identically. Stop adding attempts. Change technique.

This skill is meta. It does not own the bugfix loop (`bugfix` does) or the refactor loop (`refactor`
does) — it owns what you reach for *inside* any of them when forward motion stops.

## The six techniques

Reach for them roughly in this order; each is cheap and each attacks the impasse from a different
side.

1. **Minimal repro.** Strip the situation to the smallest thing that still fails — delete inputs,
   collapse config, remove layers. The stripping itself is diagnostic: the case that makes the
   failure vanish names the cause.
2. **Restate / rubber-duck.** Write the problem in plain words: what you expected, what you saw, and
   why you believed they would match. The gap between the last two is almost always the wrong
   assumption, and writing it down is what makes it visible.
3. **Bisect.** Find the last known-good point — a commit, a smaller input, a prior config — and halve
   the distance to the break. `git bisect` is the literal form; the mindset applies to inputs and
   data too.
4. **Change one variable.** Stop changing three things and hoping. Isolate one, observe, then the
   next. Speed comes from attribution, not from breadth.
5. **Invert.** Ask the problem backwards: "what would make this fail on purpose?" or "what must be
   true for the *current wrong* output to happen?" The answer is the assumption that is actually
   false.
6. **Escalate / ask.** After real elimination, a bounded ask — here is the repro, here is what I have
   ruled out, here is the specific decision I need — is not giving up. It is the right move, and it is
   only cheap because the five techniques above made the question precise.

**Three failed attempts of the same shape is itself a red flag.** If you are reaching for a fourth,
you have skipped this skill — pick a *different* technique, not a bigger version of the last one.

## The dev note — the way out, and the thing the user reviews

When a technique breaks the impasse, the deliverable is not just the fix — it is a short note that
makes the fix reviewable. Produce it in this exact shape:

```
Problem:    the symptom, and the wrong assumption that hid the cause
Solution:   what actually fixes it, and why the assumption was wrong
Affects:    the files, behaviours, and other callers this touches
Implements: the concrete change made
Result:     the observed outcome after the fix, with evidence (test output, the number, the repro now passing)
```

**Then hand the note to the user and let them review it before the change is finalised or
committed.** The note is the review artefact — it is faster for a human to check a five-line
"problem → result" than to reconstruct the reasoning from a diff, and it is where a wrong assumption
you *still* hold gets caught. Do not close the loop by committing; close it by getting the note
reviewed.

## Red flags

- A fourth attempt of the same shape after three failed the same way.
- A "fix" whose **Result** line has no evidence — you have a hypothesis, not a result.
- Reaching for **escalate** first, before any elimination — the ask will be too vague to answer.
- Committing the change before the user has seen the dev note.

## Before you finish

- [ ] The impasse was broken by a named technique, not by a fourth attempt of the same shape
- [ ] The dev note is complete — all five fields, **Result** carrying real evidence
- [ ] The note was handed to the user and reviewed before the change was finalised
- [ ] If the stuck thing was a defect, `bugfix` owned the loop and this owned the technique

## Read next

- [`bugfix`](../bugfix/SKILL.md) — the defect loop this plugs into when the impasse is a bug.
- [`refactor`](../refactor/SKILL.md) / [`simplicity`](../simplicity/SKILL.md) — when the impasse is a
  design one.
- [`verify`](../verify/SKILL.md) — the adversarial pass that attacks the fix's claim, not confirms it.
