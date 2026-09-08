---
name: security
description: "Surface-aware security review — the surfaces specific to this application and what enforces each. Load this before reviewing or writing anything touching authentication, authorization — permissions, roles, access control, who may see or do what — sessions, cookies, uploads, shared resources, free text bound for a model, environment variables, dependencies, or CI/CD workflows. Also load it whenever a security review is asked for, alongside the generic /security-review. Covers response schemas as security controls, prompt injection layering, capability-based sharing, the permission model, and which of the project's automated checks already cover a change mechanically."
---

# Security conventions

**Run `/security-review` too.** It is the generic pass — injection, authentication, common web weaknesses — and it is better at that than this file. This skill carries only what a generic reviewer cannot know: the surfaces that are specific to *this* application, and where the rule for each already lives.

The standing reference is the project's security documentation; the enforcement strategy is in the project's architectural decisions. **Trust those over this file** if they ever disagree, and fix the disagreement.

## What is already mechanical

Your automated guard checks run on push and in CI — structural invariants that catch common security mistakes mechanically. Before reporting a finding, check whether the build already blocks it — a finding the build catches is noise, not a review.

Common guards to check for (add your project's specific guard names here):

| Change | Typical automated checks |
|---|---|
| Anything committed | Secret-literal scanning on push; secret scanning in CI |
| An environment variable read | Env hygiene (declared in `.env.example`), client env exposure (browser bundle leaks) |
| An API route | Response schema enforcement, logging discipline — the latter matters because a logger's redact list is what keeps a cookie or a hash out of a log line, and `console` bypasses it |
| A service module | Service purity — framework core as a value, not a handle type |
| A cookie | Cookie hardening — `httpOnly`, `sameSite`, `secure`, `path` |
| A response schema | Sensitive field exposure checks on response shapes |
| Any code | Dangerous sink detection — `eval`, `new Function`, `child_process`, raw query execution, `dangerouslySetInnerHTML` |
| A CI/CD workflow | Workflow permission checks, plus linters in CI |
| A dependency | Dependency audit at the appropriate severity level in CI |

**What no check can see** is the whole point of a review: authorisation logic, whether a session is actually invalidated, whether uploaded content is what it claims to be, and everything below.

## 1. Shared-hostname cookie hazards

When multiple applications share a hostname (different path prefixes, ports, or subdomains), cookies are not reliably isolated. This is a common source of security issues.

- Use a distinct cookie name and a `Path` scoped to your application's prefix. Both values should come from centralized configuration.
- **`Path` scoping decides where a session exists at all.** Session middleware typically constructs a session only for URLs under the cookie's path. Anywhere else, the session object may be empty or missing methods — truthy but non-functional.
- **Never treat another application's session cookie as an authentication signal.** A signed cookie store scoped to the whole domain arrives regardless of who the user is to your application.
- `SameSite=Lax` alone is not sufficient CSRF protection when a sibling application shares the site.

Rule and reasoning: document your cookie isolation strategy in your security documentation and architectural decisions.

## 2. The response schema is a security control

When your API framework serialises responses against a schema, a field with no place in the schema cannot leave the process. Secure password handling typically depends on two independent filters: a safe projection at the data layer *and* response schema enforcement.

A route with no response schema is not a style problem — it is a security gap. A password digest field appearing in a response schema defeats both filters at once. Automated guards catch the mechanical half; **what they cannot catch is a schema that is present but too generous** — a user response shape handing every field of the row to whoever asks.

## 3. External identity sync

When user identity and roles are mastered in an external system and synced to your application, the sync handler carries a critical surface.

The consequence people forget: **the sync handler must invalidate sessions** — and it should do so **in the same transaction as the identity update**, not alongside it. Outside the transaction, a crash between the two leaves a revoked user holding a working session, which is the exact window server-side sessions were chosen to close.

Things to check on sync endpoints:

| Concern | What goes wrong |
|---|---|
| Signature verification | **Verify against the raw body.** Parsing then re-serializing (e.g. `JSON.parse` → `stringify`) reorders keys — verifying against the re-serialised object verifies nothing |
| Idempotency | Redelivery is the normal case, so both an event ID guard and a version guard are needed, and both must answer 2xx on a no-op |
| Credential channel | If the sync carries password digests, anyone able to forge a request sets the credential a user signs in with. Rotate the signing secret on the same schedule as your session secret, and treat a leak as an authentication incident |
| Sensitive data in responses/logs | Password digests must never reach a response (enforce via safe select lists and route schemas) or a log (redact the field name; do not log the payload whole) |
| GDPR erasure | Erasure is the one sync event where retaining the previous value is the defect. If you keep the old email in an audit record, you have retained what was legally erased. Everywhere else, keeping the prior value is good practice — which is exactly why this gets missed |

## 4. Free text reaching a model, and uploads reaching a third party

If your application accepts natural language destined for a prompt, or uploads files (especially photographs of third parties' property) to external providers:

- **Prompt injection** should be resisted structurally: user text enters as a delimited user message, never concatenated into a system prompt, and the assistant's power is bounded by the tools it is given rather than by input filtering. Document the layering in your AI pipeline documentation.
- **Data handling** — what may be sent to a provider, its retention terms, and whether it trains on submissions — must be documented and bounded. Do not widen what is sent without answering these questions.

## 5. Share links are capability tokens, not credentials

If your application shares generated content or resources via links with no login required, those links are capabilities: their own table, their own expiry, their own revocation, and no relationship to the session machinery.

Ask of any share feature: can the link be guessed, does it expire, can it be revoked, and does it expose more than the single resource it was issued for.

## 6. Being signed in is not being allowed

Authentication is who gets *in*. Authorization is what they may do once they are, and the two fail differently: a broken session shows a sign-in form, a broken permission check shows one tenant's data to another tenant's administrator.

What to look for in a review — all four of which have a wrong version that reads as correct:

- **A role string used as a capability.** `role === 'Admin'` anywhere outside the policy module is a permission check that will be missed the day the capability moves. The answer is `can(subject, Permission.X)`.
- **A scoped permission asked without its resource.** It must be `false` — unanswerable is not permitted, and a default of `true` there is the most dangerous line this model could hold.
- **A tenant comparison reached before the internal-admin check.** Internal staff may hold no tenant, so order decides whether they are refused everything; and the obvious repair — treating a null tenant as a wildcard — hands the wildcard to every unassigned account.
- **A client gate with no server gate behind it.** Hiding a control is a courtesy. If the endpoint does not re-read the role from the database on the same request, the button was the only thing stopping anyone.

The client's role-based routing is a *display* decision — what is shown — and must not grow into a second authority about what is reachable.

## When you find something

Say plainly what an attacker gets and how. If it is real, fix it or write it down — report through the project's security reporting path, and an unresolved question goes to the project's open questions rather than staying in a review comment.

If the finding is a rule that could be checked mechanically, that is a new guard check: write the rule in your conventions documentation first, then the guard implementation. Your architectural decisions document sets out which layer owns what.
