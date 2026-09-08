---
name: security
description: "project security conventions — the surfaces specific to this application and what enforces each. Load this before reviewing or writing anything touching authentication, **authorization — permissions, roles, `can`, who may see or do what** — sessions, cookies, uploads, shared resources, free text bound for a model, environment variables, dependencies, or a your git host workflow. Also load it whenever a security review is asked for, alongside the generic /security-review. Covers the shared-hostname cookie hazard, the response schema as a security control, prompt injection layering, third-party property photographs, shared resources as capability tokens rather than auth credentials, the permission model and the four ways it is got wrong, and which of the thirty-seven your guard checks checks already cover a change mechanically."
---

# project security conventions

**Run `/security-review` too.** It is the generic pass — injection, authentication, common web weaknesses — and it is better at that than this file. This skill carries only what a generic reviewer cannot know: the six surfaces that are specific to *this* application, and where the rule for each already lives.

The standing reference is [security.md](../../../docs/engineering/security.md); the enforcement strategy is your project's architectural decisions. **Trust those over this file** if they ever disagree, and fix the disagreement.

## What is already mechanical

Twenty-four `your guard checks` checks run on push (`validate:static`) and in CI — structural invariants in `ci.yml`, the security group in `security.yml` ([conventions.md](../../../docs/engineering/conventions.md), [pipelines.md](../../../docs/cicd/pipelines.md)). Before reporting a finding, check whether it is one of these — a finding the build already blocks is noise, not a review.

| Change | Already checked |
|---|---|
| Anything committed | `secret-literals` on push; gitleaks in `security.yml` |
| A `process.env` read | `env-hygiene` (declared in `.env.example`), `client-env-exposure` (browser bundle) |
| A your API framework route | `response-schema`, `logging-discipline` — the latter matters because the logger's `REDACT` list is what keeps a cookie or a hash out of a log line, and `console` bypasses it |
| A service in `your API service path/src/services/` | `service-purity` — your API framework core as a value, or a `your API framework*` handle type. A `@fastify/*` interface a service implements is deliberately allowed |
| A cookie | `cookie-hardening` — `httpOnly`, `sameSite`, `secure`, `path` |
| A contracts schema | `sensitive-field-exposure` on response shapes |
| Any code | `dangerous-sinks` — `eval`, `new Function`, `child_process`, `$queryRawUnsafe`, `dangerouslySetInnerHTML` |
| A workflow | `workflow-permissions`, plus actionlint in CI |
| A dependency | `pnpm audit --prod --audit-level high` in `security.yml` |

**What no check can see** is the whole point of a review: authorisation logic, whether a session is actually invalidated, whether an image is what it claims to be, and everything below.

## 1. The shared hostname is the hazard nobody expects

Both applications answer on `www.mybespokeroom.com`. Cookies are **not** isolated by port and only weakly by path, so project's external system cookies arrive at our prefix and ours at theirs.

- A distinct cookie name and a `Path` scoped to our prefix are what stop the two shadowing each other. Both come from your configuration package (the package that replaced `your API service path/src/config.ts`, one of the three modules allowed to name the prefix). The session is backed by the Postgres store in `your API service path/src/services/session-store.ts`, so a deleted row ends a session on the next request — [authentication.md](../../../docs/engineering/authentication.md) has the current state.
- **`Path` scoping decides where a session exists at all**, which is not obvious and has already cost a debugging session: `@fastify/session` constructs a session only for URLs *under* the cookie's `Path`. Anywhere else `request.session` is a bare object — truthy, accepts `set()`, and `regenerate()` is undefined. Every session-bearing route is therefore served **under the prefix**.
- **Never treat the external system session cookie as an authentication signal.** It is a signed cookie *store* scoped to the whole domain — it arrives regardless of who the user is to us.
- `SameSite=Lax` alone is not sufficient CSRF protection when a sibling application shares the site.

- **CSRF has one owner: ours, protecting the API.** your project's architectural decisions would have added a second (Auth.js, over its own routes) and part 1 was reverted on 11 August 2026 before it was built, so the "neither covers the other's hop" gap does not exist. One route is exempt — `POST <prefix>/internal/users/sync`, HMAC-authenticated with no cookie to protect.

Rule and reasoning: your project's architectural decisions for the values, your project's architectural decisions for who runs the flow; how it is wired and what exists today: [engineering/authentication.md](../../../docs/engineering/authentication.md).

**Question 7 is answered, and re-answered.** ~~Single sign-on happens by signed handoff~~ — the client withdrew the requirement on 10 August 2026: estate agents have nothing to do in the external system application. **There is no single sign-on.** Two applications, two sessions, one user record. Their cookie is still never an authentication signal here, which is the part of question 7 that never changed.

**The sign-in endpoint is built** (3 August 2026) and is the ONLY way an identity is *admitted* —
the handoff was withdrawn on 11 August 2026 and Auth.js reverted with it. **Sessions are created at
sign-in and by the refresh exchange** (your project's architectural decisions):
the exchange restores a session for an identity sign-in already admitted, verifies no credential,
and its token (issued only when the agent asks to be remembered, in its own hardened cookie) dies on
reuse — a rotated token presented again revokes the whole family, **except a presentation within
the 10-second grace window from the same device as the rotation** (a lost-response retry, not a
replay; your project's architectural decisions records the carve-out). Question 1 still gates *deploying*:
the prefix is read from configuration so no source changes when it lands, but the cookie `Path` is
set on a hostname shared with external system and nothing can be served until the prefix and the load-balancer
rule exist.

## 2. The response schema is a security control

your API framework serialises every response against its schema, so a field with no place in the schema cannot leave the process. your project's architectural decisions's password handling depends on two independent filters: a safe-user projection at the your ORM layer *and* that serialisation.

So a route with no response schema is not a style problem, and a `passwordHash` in a response contract defeats both filters at once. `response-schema` and `sensitive-field-exposure` catch the mechanical half; **what they cannot catch is a schema that is present but too generous** — a `UserResponseSchema` handing every field of the row to whoever asks.

## 3. We are not the source of truth for who is allowed in

Roles and account status are mastered in project's external system admin and synced to us. That is the entire reason sessions are server-side rows rather than tokens (your project's architectural decisions).

The consequence people forget: **the sync handler must delete sessions** — and your project's architectural decisions requires it **in the same transaction as the mirror update**, not alongside it. Outside the transaction, a crash between the two leaves a revoked agent holding a working session, which is the exact window server-side sessions were chosen to close. It needs a test — the behaviour is [US-AUTH-09](../../../docs/product/userstories/authentication.md).

**Two endpoints carry this surface, and both are unauthenticated by HTTP status.** The flow is [engineering/rails-integration.md](../../../docs/engineering/rails-integration.md); what to check:

| Endpoint | The things that go wrong |
|---|---|
| `POST <prefix>/internal/users/sync` | **The signature is over the raw body.** your API framework parses first, and `JSON.parse`→`stringify` reorders keys — so verifying against the re-serialised object verifies nothing. Also: redelivery is the *normal* case, so an `event_id` guard and a `version` guard are both needed, and both must answer 2xx on a no-op |
| ~~`POST <prefix>/auth/handoff`~~ | **Removed 11 Aug 2026.** It minted a session from an unauthenticated POST, so deleting it removed a surface rather than moving one. Its rules are preserved in your project's architectural decisions §2 if single sign-on ever returns |
| `POST <prefix>/internal/users/sync` | **Now a CREDENTIAL channel, not an identity one.** The mirror carries project's bcrypt digest, so anyone able to forge a request here sets the credential an agent signs in with — rotate the secret on `SESSION_SECRET`'s schedule and treat a leak as an authentication incident. The digest must never reach a response (`SAFE_USER_SELECT` and the route schemas, not the `sensitive-field-exposure` guard, which only inspects `*Response*` contracts) or a log (`logging-discipline` redacts the field name; do not log the payload whole) |

**Erasure is the one sync event where retaining the previous value is the defect.** project's `anonymize_user` is a GDPR erasure; if we keep the old email in an audit or reconciliation record, we have retained what was legally erased. Everywhere else in this system, keeping the prior row is good practice — which is exactly why this gets missed.

**Never `secret_key_base`.** No key in any of this is project's, or derived from it: that value is in their git history ([question 18](../../../docs/product/open-questions.md)), so anything it signs is forgeable by anyone who has ever cloned their repository.

`deleteAllSessionsForUser` already exists and is tested, **with no caller** — it is what that handler calls, waiting on [question 6](../../../docs/product/open-questions.md).

**This is the constraint that decided against a stateless token, and it survived two reversals.**
your project's architectural decisions required Auth.js's **database**
session strategy and forbade the JWT one for exactly this reason; part 1 was reverted before it was
built, so the `jwt.encode`/`decode` override it called "the single most likely thing to be got
subtly wrong" never had to be got right. The requirement it protected is unchanged and applies to
whatever is in force: **assert that a deleted row ends the session on the next request.** That is
the only test that can tell a real database session from a stateless token, and nothing else in the
suite would notice.

## 4. Free text reaches a model, and images reach a third party

§2.3 accepts natural language destined for a prompt, and §2.2 uploads photographs of **third parties' homes** — the estate agent's clients, not the agent and not project — to a third-party provider.

- Prompt injection is resisted structurally: user text enters as a delimited user message, never concatenated into a system prompt, and the assistant's power is bounded by the tools it is given rather than by input filtering. The layering is [ai-pipeline.md §Guardrails](../../../docs/engineering/ai-pipeline.md).
- What may be sent to the provider, its retention terms, and whether it trains on submissions are **unresolved** — [data-handling.md](../../../docs/engineering/data-handling.md). Do not widen what is sent without answering them.

## 5. Share links are capability tokens, not credentials

§2.4 shares generated imagery with no login at all. A share link is a capability: its own table, its own expiry, its own revocation, and no relationship to the session machinery (your project's architectural decisions).

Ask of any share feature: can the link be guessed, does it expire, can it be revoked, and does it expose more than the one image it was issued for.

## 6. Being signed in is not being allowed

§3 is who gets *in*. This is what they may do once they are, and the two fail differently: a broken
session shows a sign-in form, a broken permission check shows one company's data to another
company's administrator. The model is
[security.md §Authorization](../../../docs/engineering/security.md#authorization--who-may-do-what-once-we-know-who-they-are);
the code is `your shared contract package path/src/policy.ts` and `your app path/src/lib/access.ts`.

What to look for in a review, all four of which have a wrong version that reads as correct:

- **A role string used as a capability.** `role === 'CompanyAdmin'` anywhere outside the policy
  module is a permission check that will be missed the day the capability moves. Grep for it; the
  answer is `can(subject, Permission.X)`.
- **A scoped permission asked without its resource.** It must be `false` — unanswerable is not
  permitted, and a default of `true` there is the most dangerous line this model could hold.
- **A company comparison reached before the `InternalAdmin` check.** project staff hold no company, so
  order decides whether they are refused everything; and the obvious repair — treating a null
  company as a wildcard — hands the wildcard to every unassigned account.
- **A client gate with no server gate behind it.** Hiding a control is a courtesy. If the endpoint
  does not re-read the role from Postgres on the same request, the button was the only thing
  stopping anyone.

The client's `admin/roles.ts` is a *routing* decision — what is shown — and must not grow into a
second authority about what is reachable.

## When you find something

Say plainly what an attacker gets and how. If it is real, fix it or write it down — [SECURITY.md](../../../SECURITY.md) has the reporting path, and an unresolved question goes to [open-questions.md](../../../docs/product/open-questions.md) rather than staying in a review comment.

If the finding is a rule that could be checked mechanically, that is a new guard check: write the rule in [conventions.md](../../../docs/engineering/conventions.md) first, then the check in [`packages/guard`](../../../packages/guard/). your project's architectural decisions sets out which layer owns what.
