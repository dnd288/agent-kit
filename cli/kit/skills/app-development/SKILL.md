---
name: app-development
description: "Application-level conventions for Next.js App Router work. Load this when adding or changing a page, layout, loading or error boundary, route group, metadata, middleware, a server action, or a fetch to the API. Also load it for basePath, assetPrefix, redirects, Link hrefs, next/image, next/font, Server versus Client Components, streaming and Suspense, prefetching and the router cache, soft vs hard navigation, parallel and intercepting routes, forwarding the session cookie from a Server Component, the proxy hop to the API, the presigned-upload hop, or deciding where a piece of state should live (server, URL, query, context, store or props). Load it too for access control in a route — permissions, `can`, the `usePolicy` and `useCan` hooks, gating a control, and what a route renders for a refusal — and for failure handling: which of the four boundary levels a failure belongs to (widget, page content, under auth, global), why a layout's throw escapes its own `error.tsx`, and why a Server Component cannot call a `\"use client\"` module's exports. If the question is about the full data-flow model or where a store sits, load state-management instead."
---

# App Router conventions

Next.js App Router, React, Tailwind v4, TypeScript `strict`. The styling decision is
your project's architectural decisions; the visual language is `design`.

## 1. Server Components by default

`"use client"` only where interactivity requires it, and **kept at the leaves**. Tailwind is
build-time CSS with no runtime, so it works in Server Components without the client-boundary problems
runtime CSS-in-JS has here — that is part of why it was chosen.

The boundary lands on the component file, not the route. Never add `"use client"` to a `layout.tsx`
or `page.tsx` to make a child work; push it down. `component-development` covers which component files get
the directive and which must not.

**A `"use client"` module's exports are client references, not functions.** A Server Component that
imports a plain helper from one — a row mapper, a parser, a constant table — gets "Attempted to call
`x()` from the server" at *runtime*. Nothing static sees it: the directive is legal, the import is
legal, the types agree, and `your validation command` is green, because vitest has no client/server boundary and
the function is just a function there.

So a module a Server Component imports from carries **no directive at all** — it compiles into
whichever graph imports it, which is what lets one file hold both the mapper a page calls and the
column renderers a client table uses (`<domain>/columns.tsx` is the pattern). Where a file genuinely
needs the directive, split the server-safe half out. Only e2e can observe this — see `tdd` for why.

## 2. The URL prefix is configuration

When an application is served alongside another application **on the same hostname**, split by URL
prefix, the prefix must be read from configuration everywhere, with no exceptions:

- `basePath` and `assetPrefix` in `next.config` come from the environment
- every `href` and every asset path is relative, so `basePath` can do its job
- the API base path is configuration too
- the session cookie's `Path` is scoped to the prefix (§4)

**The failure mode is why this matters more than it sounds.** A hardcoded prefix works perfectly in
local development, where there is no prefix, and breaks only once the app is served behind the real
path — by which point it is in every route, link and image. `your guard checks` rejects literals for exactly
this reason; if the guard fails, fix the code rather than the guard.

## 3. Route groups

Separate route groups for audiences with different layouts and different authentication requirements.
Do not try to serve both from one layout with a conditional;
the auth requirement differs, and that is a routing concern rather than a rendering one.

## 4. Talking to the API — the Next server is the only client of it

The client and the API **never import each other**. Request and response types come from
`your contract package` — zod schemas, platform-neutral, the wire agreement. The client reaches it through `your API package` (typings + typed services), never by naming `your contract package` directly. `your guard checks` enforces both.

**The browser never calls the API directly** (your project's architectural decisions).
Traffic goes browser → Next → API, so there is no CORS to configure, no API path or error
taxonomy in the bundle, and one place to attach credentials.

When fetching from a Server Component, **forward the session cookie explicitly.** A server-side
fetch carries no browser cookies by default, so an authenticated request made without forwarding
looks anonymous to the API and returns a 401 that is easy to misread as a session bug.

**Node's `fetch` cannot resolve a relative URL.** "Empty base means same-origin" is true of a browser
and false of everything running on the server; the API origin is required and absolute, even in
production where the browser sees one hostname.

### Failures: boundaries, classification, and the serialization trap

Read your project's error boundary documentation before adding a boundary; the documentation describes
**the four levels** — widget, page content, under auth, global — and what each one cannot catch.
The level decides how much of the screen a failure takes with it, and the commonest mistake is
reaching for the route when the subtree was the right size.

Four traps to watch for:

- **A layout's throw escapes its own `error.tsx`.** A boundary renders *inside* the layout of its
  segment, so the layout that failed has none. That is why every gate classifies rather than throws:
  `getSession` answers `null` for a 401 and **throws** for every other status, and an unclassified
  throw in a layout skips the segment's boundary entirely. Adding
  a new gated segment means adding its classification, not just its boundary.
- **`redirect()` in a layout half-renders a client navigation.** It is followed correctly on a fresh
  document request and leaves a soft navigation partial. If a redirect must fire on exactly one
  route beneath a layout, it was never the layout's to make; it belongs to that route's `page.tsx`.
- **A new segment needs `loading.tsx` and `error.tsx`, and the page needs neither to look fine.**
  Both are absences, so nothing in the build notices. Make it a checklist item.
- **An error state that cannot say why hides a defect in your own code.** A `catch` that renders a
  failure for a cause it did not classify logs the error itself, or a bug in this file becomes
  indistinguishable from an outage.

The three rules that decide where a failure is *classified*:

- **A client `error.tsx` cannot classify a server-render failure.** Next serializes the error
  across the server boundary — the boundary receives a generic `Error` with a `digest`, never the
  class. So `instanceof ApiUnavailableError` inside `error.tsx` would silently do nothing in
  production. Classify where the typed error is visible (the server fetch site), and let
  `error.tsx` be the generic safety net it can truthfully be.
- **The API being down is not signed out.** The portal
  layout should render an unavailable state (retry = `router.refresh()`) instead of redirecting to
  sign-in. A down API must never render as an empty wall or a sign-in prompt.
- **Across a Server Action, classification travels as a marker**
  (e.g. `{ unavailable: true }`) because a Server Action cannot throw across the wire
  and have the boundary classify it.

## 5. `next/image` — explicit sizes

Always declare `width` and `height` on `<Image>`, or use `fill` with a positioned parent. `sizes`
is mandatory for responsive images and must match the grid — an `all 100vw` wastes bandwidth on
desktop; a wrong `sizes` loads a blurred-up 640 px image at 1200 px. This is a layout dimension
choice, not styling.

For external images, add the hostname to `remotePatterns` in `next.config`.

## 6. Fonts — one mechanism

Pick one mechanism and use it everywhere:
- `next/font` for Google Fonts or local files — built-in subsetting and self-hosting
- Tracked `@font-face` declarations in your CSS — explicit control

Do not mix. Two mechanisms means two sets of fallback stacks, two flash-of-unstyled-text behaviours,
and a visual difference that shows only in production builds.

## 7. Component conventions in a route context

A route `page.tsx` or `layout.tsx` is a Server Component by default; a `"use client"` client component
renders inside it. The boundary is the component file, not the route.

### What changes at each layer

- **Server Components** can `async`, can read `cookies()`, `headers()`, search params; cannot use
  hooks. They are the data-fetching and access-control layer.
- **Client Components** can use hooks, handle events, manage local state; cannot read server-only
  APIs. They are the interaction layer.
- **Shared modules** (no directive) compile into whichever graph imports them.

### Composites are presentational

Reusable UI composites in the component library receive data through props and emit events through
callbacks. They do not fetch, mutate, read global application state, or know about routing. A
composite that needs a `fetch` or a route push has left its layer — split it.

## 7e. Integrations — every outward boundary behind the Next server

The Next server is the single client of the API (§4), so integrations stack up here:

- **The API hop** — a Server Component or Action calls the API endpoint with the session cookie
  forwarded; the browser never reaches the API directly.
- **The upload hop** — large files go **direct to object storage** with a presigned
  PUT the proxy mints (§4); the object is never streamed through the app.
- **Background jobs** — accepting work enqueues a message; the worker
  processes it asynchronously. The client polls status, never holds a queue handle.
- **Images and fonts** — `next/image` with explicit sizes (§5); fonts via one mechanism (§6).
- **Third-party** — external integrations sit behind the proxy boundary
  (integration tests prove the seams — `test`).

Every integration that crosses the wire has a contract (or schema in the integration package for
server-to-server), tested at its seam. A new external hop without a schema behind the proxy is a
feature that has not been specified yet.

## 7f. Access control — asking the question from a route

Read the security documentation for your project's authorization model first,
and load `security` for a review. Your contract package defines the policy vocabulary;
your app's access module is how a component asks.

What is specific to a route, and therefore lives here:

- **Server Components compute; Client Components ask a hook.** `can(subjectOf(session), Permission.X)`
  in a `page.tsx` (`<domain>/page.tsx` is the pattern), `usePolicy()` / `useCan(permission,
  resource)` below a `"use client"` boundary. Both reach the policy through your API package — naming
  the contract package from the client fails the boundary check.
- **Pass the decision, not the subject.** A screen receives `canCreate={…}`, not a session it can
  interrogate: UI composites are presentational, and a policy call inside one is the same
  violation as a `fetch` (§7). It also removes the appearance a pixel baseline could not predict.
- **A refused read is a routing outcome, and each status has a different one.** 401 → the sign-in
  redirect carrying where they were going; 403/404 → `notFound()`; anything else → the segment's
  boundary. A page that renders an empty list for a refusal is telling the reader there is nothing
  there, which is the same defect class as rendering an outage as "you have no items" (§4).
- **Say at each gate that it is a courtesy.** The comment is load-bearing: without it the next
  reader assumes the client check or the server check is redundant, and removes one.

## 8. Business logic does not live here

Not in a route handler, not in a server action. It belongs in services as plain functions
with no framework types in their signatures — see `api-contract`. A route handler that grew a
transaction and three branches is a service that ended up in the wrong workspace.

## 9. Responsive is a functional requirement

Treat responsive design as testable behaviour, not polish.

Use `dvh`, never `vh` — mobile browser chrome makes `vh` wrong. The viewport matrix and the
tablet-landscape breakpoint are documented in `design`; that is the one copy.
