---
name: api-contract
description: "API conventions for routes, contract schemas and service business logic. Load this when adding or changing an endpoint, a route handler, a framework plugin or hook, a request or response schema, a service function, a data-layer query or migration, or anything in the contract package. Covers the mandatory response schema — it is what keeps password hashes out of responses — why services take no framework types, what may and may not enter the contract package, session and cookie handling, which side creates a session, the one-migration-per-pull-request rule, and guard invariants that fail the build when these are broken."
---

# API conventions

Your web framework over PostgreSQL with your ORM. The framework choice is
your project's architectural decisions, the ORM is
your project's architectural decisions, authentication is
documented in your engineering/authentication documentation.

## 1. Three layers, one-way arrows

```
plugins/  →  route handler  →  service  →  your ORM
                    ↑
              your contract package   (the shared agreement; the client reaches it via your API package)
```

- **Route handlers** are thin: parse, call a service, return. Framework types live here and only here.
- **Services** in `your API service path` hold the business logic as plain functions.
- **Plugins** in `your API plugin path` hold what every route needs and no route should repeat:
  the error shape, the request context, the ORM lifecycle. One concern each, wrapped in
  the framework's plugin system so decorators and hooks escape encapsulation.
- **`your contract package`** sits beside both. It holds the client/API agreement, is imported by the API
  directly and by the client through `your API package`, and imports neither of them. Server-to-server
  payloads live in `your integration package`, never here.

`client` and `api` **never import each other**. `your guard checks` enforces it.

**Which of the three is it?** If it needs the framework, it is a plugin or a route. If it needs a
request, it is a route. If it is neither, it is a service — and the `service-purity` check will hold
you to that.

**There is no fourth layer.** A new business case is a row, a column, or a branch inside an existing service — never a manager, a handler or a wrapper between two of the three. If a case will not fit, the finding is that the service has the wrong seam, not that the stack needs another floor (golden rules).

## 2. `buildApp()` is built without listening

The app is constructed separately from the call that binds a port. That is what makes routes
testable — a test builds the app and injects requests without a live socket. Your testing documentation
assumes this shape.

One framework plugin per feature area, registered by `buildApp()`.

**Type providers are a type-level declaration and nothing else.** The compilers beneath them —
validator and serializer — are what make a schema run. Without them a route with a response schema
does not fail quietly; it fails to *register*, so `ready()` rejects and the app does not boot.

Never add a plugin by editing `buildApp()` alone. A plugin that decorates or hooks must be wrapped in
the framework's plugin wrapper, or its decoration stays inside its own encapsulation context and the
routes that need it see nothing — a failure that reads like a missing registration.

## 3. Every route declares a response schema

Not a style rule. **The framework serialises responses through the schema, so the schema is what keeps a
password hash out of a response.** A route without one is a security gap rather than an oversight —
your engineering conventions say so, and `your guard checks` fails the build.

Write the response schema before the handler. If you are unsure what a route returns, that is a sign
the service boundary is not settled yet.

## 4. What may enter `your contract package`

Only the agreement between the client and the API: endpoint schemas and the enums they reference.

- **zod is its only runtime dependency.** No Node built-ins, no browser globals, or it cannot be
  imported from both sides.
- Before adding anything, ask whether it is part of that agreement. If not, it belongs in `client/`
  or `api/`.
- The package is deliberately **not** called `shared`. "Shared" invites everything; a specific name
  enforces scope.
- **What must NOT enter**: server-to-server payloads. Those live in
  `your integration package`; a server-to-server response lives beside the route that serves it.

`your guard checks` checks the dependency and platform-neutrality rules mechanically.

## 5. Services take no framework types

A service signature with framework request types in it is wrong. Plain arguments in, plain data out.

The reason is practical rather than aesthetic: services are where most of the tests live, and a
function that needs a framework request object needs one constructed in every
test. Keep the framework at the edge.

`your guard checks`'s `service-purity` check enforces it, and the line it draws is worth knowing because it
is narrower than "no framework":

- **Rejected** — importing the framework as a value, or type-importing any request/reply/instance handles.
- **Allowed** — a plugin interface a service implements, such as `SessionStore`, and
  `Session`, which is a data shape rather than a handle.

Two things a service still must not do, which no check can see: reach for `app.db` instead of
importing `db` from `your data layer`, and read the request context to obtain an argument it should
have been passed. `currentUserId()` exists for logging and audit; anything that authorises takes the
user as a parameter, so the test can pass a different one.

**Where does cross-cutting state live, then?** In a request-context plugin using
`AsyncLocalStorage` carrying the request id, the request logger and the current user — set up early.

## 6. Sessions and cookies

Server-side sessions stored in the database, strong password hashing — **not JWTs** unless
your architecture specifically chose them. Roles mastered in an external identity provider are synced
rather than embedded in tokens, so a revocation takes effect immediately rather than waiting for
token expiry.

**Authentication topology:**

- **The API creates and owns sessions** — a single session-establishment route is the
  authority, reached from the password sign-in route and from the refresh exchange.
  The `authentication` capability requires exactly one component to admit an identity.
- **The API verifies credentials and creates sessions.** Strong hashing, constant-cost
  failure and per-account throttle stay here, where they are tested.
- **The browser does not call the API directly.** Traffic is browser → Next → API, so a route may assume a
  single trusted hop but must not assume the hop authenticated anybody.

**Read your authentication documentation before building
against any of it** — it carries the inventory and the traps below.

Two hazards:

- **Cookie collision.** When sharing a hostname with another application, use a **distinct
  cookie name** and `Path` **scoped to the application's prefix**. Another application's session cookie
  is **never** an authentication signal.
- **`Path` scoping decides where a session exists at all.** The session middleware constructs a session only
  for URLs *under* the cookie's `Path`; anywhere else `request.session` is a bare object whose
  `regenerate()` is undefined — a 500 from code that typechecks. **Serve every
  session-bearing route under the prefix.** `/health` stays outside it deliberately: it is for the
  load balancer and touches no session.

Two more things that are not obvious from reading the code:

- **`session.regenerate()` does not destroy the old session.** It writes the new id and leaves the
  previous row live, so rotation alone does not stop session fixation. Destroy the old id explicitly.
- **Rate-limit middleware signals a refusal by throwing**, and the framework builds the body from the
  thrown error. A plain object returned from the error response builder produces a **500**, not a 429.

## 7. ORM and migrations

**Your database documentation owns the data model** — read it before adding a table.

- **One migration per pull request**, and it must be backward compatible with the previous release,
  because both versions run at once during a rollout.
- **Never edit a migration that has been applied** on a shared environment.
- Pull and re-run before branching, or you will produce a conflicting migration.
- `your guard checks` fails a schema change with no new migration directory.

**Backward compatible is now enforced, and a drop needs a person.** The migration safety check refuses a new
migration containing `DROP TABLE`, `DROP COLUMN`, `ALTER TABLE ... RENAME`, `TRUNCATE`, an
unqualified `DELETE FROM`, `DROP SCHEMA`, or `SET NOT NULL` with no backfill above it.
Your conventions document owns the rule and the escape-hatch syntax.

**Do not assume you can authorise a drop by writing a good reason.** For anything that destroys
data the comment is necessary and not sufficient: a CI workflow may also require a
label on the pull request, which only a person can add. Write the migration, write the reason, and say
in the pull request body that it needs the label.

**Run it before you claim it works.** Verify the migration against a throwaway database that applies the whole
history and diffs the result against the schema. A migration that has never been executed is a
migration whose SQL nobody has checked.

## 8. Infrastructure boundaries

Know what you own and what you do not. If the ORM plugin can provision databases and manage
connection strings, that capability may be out of scope. Use it for schema, migrations and queries
against your own local database — never to provision or alter anything in another team's estate.

## 9. Before you finish

- Response schema on every route
- Error type distinction: a rejected fetch (the API is down) vs a non-2xx answer (the API responded
  with an error). The two are the "API unreachable" vs "API answered" distinction an error boundary
  renders differently, and they MUST be distinguishable by type, never by message text.
- Contract schemas in `your contract package`, zod-only, platform-neutral
- No framework types in a service signature
- One migration, backward compatible
- No credential anywhere — not in code, not in a log line, not in a commit message
- `your validation command` passes. If a check could not be run, say which and why
