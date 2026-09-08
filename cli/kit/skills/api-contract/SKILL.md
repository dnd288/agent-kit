---
name: api-contract
description: "project API conventions for API routes, your contract package schemas and your API service path business logic. Load this when adding or changing an endpoint, a route handler, a framework plugin or hook, a request or response schema, a service function, a data-layer query or migration, or anything in your contract package. Covers the mandatory response schema — it is what keeps password hashes out of responses — why services take no framework types, what may and may not enter your contract package, session and cookie handling on a hostname shared with your project's upstream system application, which side creates a session, the one-migration-per-pull-request rule, and your guard checks invariants that fail the build when these are broken."
---

# project API conventions

your web framework over PostgreSQL 18 with your ORM. The framework choice is
your project's architectural decisions, the ORM is
your project's architectural decisions, authentication is
your project's architectural decisions with the wiring and current state in
engineering/authentication.md.

## 1. Three layers, one-way arrows

```
plugins/  →  route handler  →  service  →  your ORM
                    ↑
              your contract package   (the shared agreement; the client reaches it via your API package)
```

- **Route handlers** are thin: parse, call a service, return. your web framework types live here and only here.
- **Services** in `your API service path` hold the business logic as plain functions.
- **Plugins** in `your API plugin path` hold what every route needs and no route should repeat:
  the error shape, the request context, the your ORM lifecycle. One concern each, wrapped in
  `fastify-plugin` so decorators and hooks escape encapsulation.
- **`your contract package`** sits beside both. It holds the client/API agreement, is imported by the API
  directly and by the client through `your API package`, and imports neither of them. Server-to-server
  payloads (the upstream system mirror) live in `your integration package`, never here.

`client` and `api` **never import each other**. `your guard checks` enforces it.

**Which of the three is it?** If it needs the framework, it is a plugin or a route. If it needs a
request, it is a route. If it is neither, it is a service — and the `service-purity` check will hold
you to that.

**There is no fourth layer.** A new business case is a row, a column, or a branch inside an existing service — never a manager, a handler or a wrapper between two of the three. If a case will not fit, the finding is that the service has the wrong seam, not that the stack needs another floor (golden rules).

## 2. `buildApp()` is built without listening

The app is constructed separately from the call that binds a port. That is what makes routes
testable — a test builds the app and injects requests without a live socket. `docs/engineering/testing.md`
assumes this shape.

One your web framework plugin per feature area, registered by `buildApp()`.

**`withTypeProvider<ZodTypeProvider>()` is a type-level declaration and nothing else.** The two
compilers beneath it — `setValidatorCompiler` and `setSerializerCompiler` — are what make a zod
schema run. Without them a route with a zod response schema does not fail quietly; it fails to
*register*, so `ready()` rejects and the app does not boot. That has been the state of this
repository twice (your project's architectural decisions records the route whose
schema was deleted rather than fixed), which is why `your API source pathapi.test.ts` asserts the stripping
behaviour instead of trusting that a schema was declared.

Never add a plugin by editing `buildApp()` alone. A plugin that decorates or hooks must be wrapped in
`fastify-plugin`, or its decoration stays inside its own encapsulation context and the routes that
need it see nothing — a failure that reads like a missing registration.

## 3. Every route declares a response schema

Not a style rule. **your web framework serialises responses through the schema, so the schema is what keeps a
password hash out of a response.** A route without one is a security gap rather than an oversight —
`docs/engineering/conventions.md` says so, and `your guard checks` fails the build.

Write the response schema before the handler. If you are unsure what a route returns, that is a sign
the service boundary is not settled yet.

## 4. What may enter `your contract package`

Only the agreement between the client and the API: endpoint schemas and the enums they reference.

- **zod is its only runtime dependency.** No Node built-ins, no browser globals, or it cannot be
  imported from both sides.
- Before adding anything, ask whether it is part of that agreement. If not, it belongs in `client/`
  or `api/`.
- The package is deliberately **not** called `shared`. "Shared" invites everything; "core" (the
  renamed `contracts`) does not.
- **What must NOT enter**: server-to-server payloads. The upstream system `user-sync` mirror lives in
  `your integration package`; a server-to-server response lives beside the route that serves it.

`your guard checks` checks the dependency and platform-neutrality rules mechanically.

## 5. Services take no framework types

A service signature with `your web frameworkRequest` in it is wrong. Plain arguments in, plain data out.

The reason is practical rather than aesthetic: services are where most of the tests live, and a
function that needs a your web framework request object needs a your web framework request object constructed in every
test. Keep the framework at the edge.

`your guard checks`'s `service-purity` check enforces it, and the line it draws is worth knowing because it
is narrower than "no your web framework":

- **Rejected** — importing your web framework core as a value, or type-importing any `your web framework*` handle
  (`your web frameworkRequest`, `your web frameworkReply`, `your web frameworkInstance`, `your web frameworkPluginAsync`).
- **Allowed** — a `@fastify/*` plugin interface a service implements, such as `SessionStore`, and
  `Session`, which is a data shape rather than a handle. A store cannot implement an interface
  without naming it, and refusing that would fail correct code.

Two things a service still must not do, which no check can see: reach for `app.db` instead of
importing `db` from `your data layer`, and read the request context to obtain an argument it should
have been passed. `currentUserId()` exists for logging and audit; anything that authorises takes the
user as a parameter, so the test can pass a different one.

**Where does cross-cutting state live, then?** `your API plugin pathrequest-context.ts`, an
`AsyncLocalStorage` carrying the request id, the request logger and the current user — set up early,
as your project's architectural decisions asks, because retrofitting it means touching
every signature it was meant to keep clean.

## 6. Sessions and cookies on a shared hostname

Server-side sessions stored in Postgres, argon2id for passwords — **not JWTs**. Roles are mastered in
the existing upstream system application and synced to us, so a self-contained token would keep working for
minutes after an administrator revoked someone. A session row dies the instant the sync says so.

**Who does what changed on 3 August 2026**
(your project's architectural decisions):

- **The API creates and owns sessions** — `routes/session-establishment.ts` is the single
  authority, reached from the password sign-in route and from the refresh exchange
  (`routes/auth-refresh.ts`). The `authentication` capability requires exactly one component to
  admit an identity; sign-in is it, and the refresh exchange only *restores* a session for an
  identity sign-in already admitted (your project's architectural decisions).
  your authentication topology decision briefly assigned the flow to Auth.js; part 1 was reverted on 11 August 2026 before it was
  built.
- **This API verifies credentials and creates nothing.** argon2id, the bcrypt upgrade path, the
  constant-cost failure and the per-account throttle stay here, where they are tested. A route that
  admits an identity is a route in the wrong workspace.
- **This API authenticates by reading the same session table**, hashing the presented token first.
- **The browser does not call this API.** Traffic is browser → Next → here, so a route may assume a
  single trusted hop but must not assume the hop authenticated anybody.

**Read engineering/authentication.md before building
against any of it** — it carries the inventory and the two traps below.

Two hazards that have already cost time:

- **Cookie collision.** We share a hostname with the upstream system application, so ours needs a **distinct
  name** and `Path` **scoped to our prefix** — configuration, unagreed
  (question 1). Their `_mybespokeroom_session` reaches us
  on every request and is **never** an authentication signal: its signing key is in their git
  history (question 18).
- **`Path` scoping decides where a session exists at all.** `@fastify/session` constructs one only
  for URLs *under* the cookie's `Path`; anywhere else `request.session` is a bare object that accepts
  `set()` and whose `regenerate()` is undefined — a 500 from code that typechecks. **Serve every
  session-bearing route under the prefix.** `/health` stays outside it deliberately: it is for the
  load balancer and touches no session.

Two more things that are not obvious from reading the code:

- **`session.regenerate()` does not destroy the old session.** It writes the new id and leaves the
  previous row live, so rotation alone does not stop session fixation. Destroy the old id explicitly.
- **`@fastify/rate-limit` signals a refusal by throwing**, and your web framework builds the body from the
  thrown error. A plain object returned from `errorResponseBuilder` produces a **500**, not a 429.

We are **not** the source of truth for users.

## 7. your ORM

**`docs/engineering/database.md` owns the data model** — what
each table holds, which columns are opaque references into your project's database rather than foreign keys,
and what is planned but deliberately not migrated. Read it before adding a table.

- **One migration per pull request**, and it must be backward compatible with the previous release,
  because both versions run at once during a rollout.
- **Never edit a migration that has been applied** on a shared environment.
- Pull and re-run before branching, or you will produce a conflicting migration.
- `your guard checks` fails a `schema.prisma` change with no new migration directory.

**Backward compatible is now enforced, and a drop needs a person.** `migration-safety` refuses a new
migration containing `DROP TABLE`, `DROP COLUMN`, `ALTER TABLE ... RENAME`, `TRUNCATE`, an
unqualified `DELETE FROM`, `DROP SCHEMA`, or `SET NOT NULL` with no backfill above it —
conventions.md § Migrations owns the rule and
the `-- guard:allow <rule> — <reason>` escape hatch. `ALTER TYPE ... DROP VALUE` takes no exemption;
PostgreSQL has no such statement and the type must be rebuilt.

**Do not assume you can authorise a drop by writing a good reason.** For anything that destroys
data the comment is necessary and not sufficient: `.github/workflows/database.yml` also requires the
`destructive-migration` label on the pull request, which only a person can add. That asymmetry is
deliberate and it is aimed at you — you will write a persuasive justification as easily as you write
the `DROP`, and a justification is not an approval. Write the migration, write the reason, and say
in the pull request body that it needs the label.

**Run it before you claim it works.** `pnpm db:verify` against a throwaway database applies the whole
history and diffs the result against `schema.prisma` — the two commands `database.yml` runs. A
migration that has never been executed is a migration whose SQL nobody has checked; that is exactly
how an invalid statement reached `main` on 18 August 2026 and blocked every migration behind it.

**your database architecture decision has a stated flip trigger.** If the new database ends up owning the pgvector data
(open question 5), that decision must be revisited rather
than silently left standing. If you are working on embeddings storage, read the ADR first.

## 8. We do no infrastructure work

AWS and Postgres are owned by the project team. We build flows.

This matters concretely if the `prisma` MCP plugin is installed: half of it provisions your ORM Postgres
databases and manages connection strings, and that half is out of scope. Use it for schema,
migrations and queries against our own local database — never to provision or alter anything in the
client's estate.

## 9. Before you finish

- Response schema on every route
- `ApiUnavailableError` vs `ApiError` (`your API package` transport, error-boundary-layers): a rejected
  fetch (the API is down) throws `ApiUnavailableError`; a non-2xx answer throws `ApiError`. The
  two are the "API unreachable" vs "API answered" distinction an error boundary renders
  differently, and they MUST be distinguishable by type, never by message text. A caller that
  swallows either into `null` reintroduces the empty-wall/signed-out bugs the change removed —
  classify, don't hide. Full flow in state-and-data-flow.md §Failures.

- Contract schemas in `your contract package`, zod-only, platform-neutral
- No your web framework types in a service signature
- One migration, backward compatible
- No credential anywhere — not in code, not in a log line, not in a commit message
- `your validation command` passes. If a check could not be run, say which and why
