---
name: app-development
description: "project client conventions for Next.js App Router work. Load this when adding or changing anything under your app route path: a page, layout, loading or error boundary, route group, metadata, middleware/proxy, a server action, or a fetch to your API. Also load it for basePath, assetPrefix, redirects, Link hrefs, next/image, next/font, Server versus Client Components, streaming and Suspense, prefetching and the router cache, soft vs hard navigation, parallel and intercepting routes, forwarding the session cookie from a Server Component, the Next-as-proxy hop to your web framework, the S3 presigned-upload hop, or deciding where a piece of state should live (server, URL, query, context, store or props). Load it too for access control in a route — permissions, `can`, the `useUserPolicy` and `useCan` hooks, gating a control, and what a route renders for a refusal — and for failure handling: which of the four boundary levels a failure belongs to (widget, page content, under auth, global), why a layout's throw escapes its own `error.tsx`, and why a Server Component cannot call a `"use client"` module's exports. If the question is about the full data-flow model or where a store sits, load state-management instead. This application shares a hostname with your project's upstream system application and its URL prefix is configuration that has not been agreed yet. The client never imports the API; your API package is how it talks to it, and the wire agreement lives in your contract package."
---

# project App Router conventions

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
the function is just a function there. It has cost two defects, and in one of them the page's own
`catch` rendered a designed error state over it for as long as the page existed.

So a module a Server Component imports from carries **no directive at all** — it compiles into
whichever graph imports it, which is what lets one file hold both the mapper a page calls and the
column renderers a client table uses (`admin/companies/columns.tsx` is the pattern; `admin/users/columns.tsx`
is what happens when it is not). Where a file genuinely needs the directive, split the server-safe
half out (`members/people-tab.ts`). Only e2e can observe this — see `tdd` for why.

## 2. The URL prefix is configuration — the sharpest hazard here

This application is served alongside your project's upstream system application **on the same hostname**, split by URL
prefix. **The prefix has not been agreed** (open question 1,
blocking). The `/staging/*` in the architecture diagram is illustrative, not decided.

So it is read from configuration, everywhere, with no exceptions:

- `basePath` and `assetPrefix` in `next.config` come from the environment
- every `href` and every asset path is relative, so `basePath` can do its job
- the API base path is configuration too
- the session cookie's `Path` is scoped to the prefix (§4)

**The failure mode is why this matters more than it sounds.** A hardcoded prefix works perfectly in
local development, where there is no prefix, and breaks only once the app is served behind the real
path — by which point it is in every route, link and image. `your guard checks` rejects literals for exactly
this reason; if the guard fails, fix the code rather than the guard.

## 3. Route groups

Two audiences with different layouts and different authentication:

- the **authenticated agent portal** — the §2.2 workflow
- the **unauthenticated share view** — a recipient opening a shared result, who has no account

Separate route groups, separate layouts. Do not try to serve both from one layout with a conditional;
the auth requirement differs, and that is a routing concern rather than a rendering one.

## 4. Talking to the API — the Next server is the only client of it

The client and the API **never import each other**. Request and response types come from
`your contract package` — zod schemas, platform-neutral, the wire agreement. The client reaches it through `your API package` (typings + typed services), never by naming `your contract package` directly. `your guard checks` enforces both.

**The browser never calls your web framework directly** (your project's architectural decisions).
Traffic goes browser → Next → your web framework, so there is no CORS to configure, no API path or error
taxonomy in the bundle, and one place to attach credentials. That single hop is also where the upstream system
API (question 5) and S3 signed-URL minting will live.

One planned exception: a room photograph is uploaded **direct to S3 with a presigned `PUT`** the
proxy issues. Streaming 10 MB through two servers to reach a bucket is cost with no benefit, and the
proxy stays in control by being the only thing that can mint the URL.

When fetching from a Server Component, **forward the session cookie explicitly.** A server-side
fetch carries no browser cookies by default, so an authenticated request made without forwarding
looks anonymous to the API and returns a 401 that is easy to misread as a session bug.

**Node's `fetch` cannot resolve a relative URL.** "Empty base means same-origin" is true of a browser
and false of everything running on the server; the API origin is required and absolute, even in
production where the browser sees one hostname.

### Failures: boundaries, classification, and the serialization trap

The full architecture is state-and-data-flow.md §Failures
(error-boundary-layers), including **the four levels** — widget, page content, under auth, global —
and what each one cannot catch. Read that table before adding a boundary; the level decides how much
of the screen a failure takes with it, and the commonest mistake is reaching for the route when the
subtree was the right size.

Four traps that have each cost a real defect here:

- **A layout's throw escapes its own `error.tsx`.** A boundary renders *inside* the layout of its
  segment, so the layout that failed has none. That is why every gate classifies rather than throws:
  `getSession` answers `null` for a 401 and **throws** for every other status, and an unclassified
  throw in `(portal)/layout.tsx` or `admin/layout.tsx` skips the segment's boundary entirely. Adding
  a new gated segment means adding its classification, not just its boundary.
- **`redirect()` in a layout half-renders a client navigation.** It is followed correctly on a fresh
  document request and leaves a soft navigation partial — the console's landing bug, 15,715 bytes
  against 47,400. If a redirect must fire on exactly one route beneath a layout, it was never the
  layout's to make; it belongs to that route's `page.tsx`.
- **A new segment needs `loading.tsx` and `error.tsx`, and the page needs neither to look fine.**
  Both are absences, so nothing in the build notices. `ADMIN.md` §5.2 makes it a checklist item for
  the console; the same is true everywhere.
- **An error state that cannot say why hides a defect in our own code.** A `catch` that renders a
  failure for a cause it did not classify logs the error itself, or a bug in this file becomes
  indistinguishable from an outage — which is exactly how the global user list stayed broken.

The three rules that decide where a failure is *classified*:

- **A client `error.tsx` cannot classify a server-render failure.** Next serializes the error
  across the server boundary — the boundary receives a generic `Error` with a `digest`, never the
  class. So `instanceof ApiUnavailableError` inside `error.tsx` would silently do nothing in
  production. Classify where the typed error is visible (the server fetch site), and let
  `error.tsx` be the generic safety net it can truthfully be.
- **The API being down is not signed out.** `getSession` rethrows `ApiUnavailableError`; the portal
  layout renders `<UnavailableState />` (retry = `router.refresh()`) instead of redirecting to
  sign-in. A down API must never render as an empty wall or a sign-in prompt.
- **Across a Server Action, classification travels as the `{ unavailable: true }` marker**
  (`lib/api-unavailable.ts`), never as a message to match. An action catches the reachability
  failure and returns the marker; the client narrows with `isApiUnavailable`.

New fetch sites follow the existing pattern: catch `ApiUnavailableError` → render
`<UnavailableState />`; throw everything else to the segment's `error.tsx`. Do not reintroduce a
swallowing variant of `apiTryFetch` — it was removed because it hid the distinction.

Two more rules from the review of the pattern:

- **Classify EVERY server fetch in a page.** The wizard page classifies both its reads (the
  project and the room types) — a connection drop between two calls in one render is still "the
  API is down", and a page with one classified fetch and one unclassified one renders the wrong
  state when the second fails.
- **A state rendered without the shell carries a way off, and the catch sites log.** Retry is
  not an escape: `UnavailableState` and the wizard's `error.tsx` footer-link to `/projects`; the
  classified catch sites write one `console.error` naming the route, so an outage produces the
  same correlation signal the generic boundaries' digest logs do.

### Sessions

**The your web framework API owns the sign-in flow and the session cookie**, against our own `Session` table.
A Server Action makes the Next → your web framework hop and relays the `Set-Cookie`; the browser never calls
the API directly.

Auth.js was chosen for this in your project's architectural decisions and
**part 1 was reverted on 11 August 2026 before it was built** — there is no `next-auth` dependency
and no `[...nextauth]` route. Parts 3 and 4 stand: your web framework verifies credentials, and the Next
server remains the only client of the API.

Two hazards, both of which have already cost time:

- **The cookie collides.** We share a hostname with the upstream system application, so ours uses a **distinct
  name** with `Path` **scoped to our prefix**. Their `_mybespokeroom_session` arrives on every
  request and is **never** an authentication signal — its signing key is in their git history.
- **`Path` scoping decides where a session exists at all.** `@fastify/session` builds a session only
  for URLs *under* the cookie's `Path`; elsewhere `request.session` is a bare object whose
  `regenerate()` is undefined. So every session-bearing route is served **under the prefix**. A
  mis-set prefix does not just break links — it makes the session silently not exist.

## 5. Images

This application is image-heavy by nature — room photographs, generated renders, before-and-after
comparisons, product imagery. `next/image` throughout.

Size every image explicitly. The Sign in photograph is the LCP element of the first screen a user
sees and ships at 1600×1067, so it needs dimensions, a `srcset` and fetch priority or it shifts the
form as it loads.

## 6. Fonts — resolve this, do not do both

There are two self-hosting mechanisms in play and they are alternatives, not complements:
`docs/design/design-system.md` specifies `next/font`, while `design`'s `tokens/fonts.css`
already declares `@font-face` rules. Pick one deliberately when the client is scaffolded.

The binaries are tracked under `your resources pathfonts/` (open question 2b answered).
`pnpm fonts:sync` regenerates `your component library source pathstyles/fonts.css` from the extractor.

## 7. A route fetches; a screen renders

This is the seam the whole UI strategy rests on. `your UI package` screens are props-only — `your guard checks`
forbids `next/*`, `your API package`, `your contract package`, `your data layer` and `fetch` there — so the route is where data
meets design:

```tsx
// your app route path(portal)/rooms/page.tsx
import { RoomPickerScreen } from 'your UI package';

export default async function Page() {
  const rooms = await getRooms();          // fetch here, with the cookie forwarded (§4)
  return <RoomPickerScreen rooms={rooms} />; // the screen only renders
}
```

**Why bother, when the route could just render the markup?** Because a screen that cannot fetch has
exactly one appearance for a given set of props, which is what makes a pixel baseline against its
design frame meaningful (your project's architectural decisions). Move a `fetch` into
the screen and you have not just bent a rule — you have removed the only mechanism that verifies the
design was implemented correctly.

Three things belong to the route and never to the screen:

- **Data.** Fetch it, map it, pass it. A DTO from `your API package` is mapped to the screen's own prop
  types at this boundary — deliberately, so the design layer never depends on the API's shape.
- **Routing.** `next/link` applies `basePath`; pass it into a slot (§2).
- **Session and redirects.** Auth decisions are routing concerns.

Keep the route thin. If it grows a transaction and three branches, see §8.

## 7a. Where state goes — ask in this order

your project's architectural decisions decides it, with one scoped
exception (your project's architectural decisions);
state-and-data-flow.md is the working copy. The
first "yes" wins:

**The ladder is the short form.** The full model — the layered data flow from externals down to a
component, one-way arrows, and why each rule exists — is the **`state-management` skill**
(`references/data-flow.md`). Load that one when the question is architectural; this ladder when you
are deciding the next value.

**The server→client seam adds a pre-step.** SSR owns starting params — a store is never seeded by a
client fetch. Initial state is server-resolved and passed once into a provider or hook at the
boundary. Because `apiFetch` is `no-store`, a server read used more than once in one render needs
`cache()` — `getSession` is the template (`lib/session.ts`). Independent reads run in parallel.
Full contract: state-and-data-flow.md `#the-server-client-seam`.

**Ladder citations use the §7a anchor, not bare numbers.** A bare "rule 4" in another file is
ambiguous — §7a numbers it 0 (identity) while the owning document numbers it 4. Cite it as
"§7a rule 0" or "the session-slice rule" rather than a bare number.



0. **Is it the signed-in identity — status, user, sign-in/sign-out transitions?** → **the session
   store in `your store package`**: `SessionProvider` hydrated from the server-read session (`getSession()`
   in `your app source pathlib/session.ts`), `useSession()` / `useUser()`, transitions injected by the
   app. The store never fetches; a reload re-derives it from the next server render.
1. **In the database?** → Server Component reads it, Server Action mutates it, then **revalidate**.
   A mutation that does not invalidate leaves a stale render that reads exactly like a caching bug.
2. **Should a reload, a shared link or the back button reproduce it?** → **the URL**. The wizard
   step, filters, whether a dialog is open. This is what makes a bug report a link.
3. **Changes without the user acting?** → TanStack Query. Today that is the Generating screen only;
   it is not the general data layer.
4. **Read widely, written almost never?** → Context. The locale, toasts.
5. **Interaction state confined to one screen, more than two values?** → a Zustand slice **in the
   feature folder**. Never a global store. When the screen is a `your UI package` composite rather than a
   route, the feature folder is the screen's own folder in `your UI package path` and the slice is a
   **screen store** holding view state only — `state-management`'s composite rung and
   your project's architectural decisions own it. That is not a licence
   to import a store INTO a composite: `your store package` and this app's feature slices stay banned there.
6. **Otherwise** → props.

**A store holding a copy of a database row means the answer was 1** — the session store is the one
exception, and it holds only what the server resolved. Treat every URL value as attacker-controlled:
map codes through an allow-list, and put any redirect target through an open-redirect guard.

## 7b. Routing and navigation — what makes a click become a URL

Routes are the directory tree under `app/`; segments map to layouts and pages; `?step=` and other
search params are how filters and wizard state travel in a URL. The conventions that matter here:

- **The URL is the committed truth; the UI may preview.** The wizard's bug history is this lesson:
  an optimistic URL written *before* the persist landed let the force-dynamic server page re-render
  mid-flight, hit the forward lock, and hard-reload. Move the URL only when the server has committed
  (the flow hook is the pattern).
- **`<Link>` for destinations, `router.push/replace` for programmatic changes.** A Click → `router`
  → RSC is a *soft* navigation; a full `window.location` or a server `redirect()` is *hard* — which
  is a reload with a flash. Know which one your code produces: a `redirect()` from a Server
  Component with a stale dependency is the hidden reload.
- **`router.back()` and the back button** restore the URL, so the screen must follow `?step=` — the
  URL-seed/follow effect in `use-wizard-flow.ts` is the template. Read and write it through
  `useRouterState` (`src/lib/use-router-state.ts`), which pairs the allow-list with the writer and
  hands back the **parsed value**; keying a follow effect on the `useSearchParams()` object instead
  makes it re-run when the URL has not moved, and the screen jumps backwards.
- **Parallel (`@slot`) and intercepting (`(.)`) routes** compose; the share view and portal layout
  each get their own route group (§3) and never one conditional layout.

The deep dive — prefetching, the router cache, soft vs hard navigation, when a navigation becomes a
full load — is `references/nextjs-production-readiness.md` §2.

## 7c. Rendering — RSC by default, client at the leaves, streaming for the slow

Three rendering modes exist; the file chooses one, not the route:

- **Server Components (the default).** Fetch, read the DB, hold secrets, run business code that
  must not ship to the browser. `"use client"` is a boundary you push down, not up (§1).
- **Client Components.** Interaction, `useState`/stores, effects, `next/navigation`. They still
  render server-side once (SSR) — a client leaf is not an excuse to fetch in an effect.
- **Async server components + `Suspense` segment-by-segment.** A slow query becomes a skeleton,
  not a spinner for the whole page. `loading.tsx` is the fallback for the route segment it sits in.

The page is `force-dynamic` when the data is per-request (sessions, the wizard) — the wizard page
declares it; a static marketing page should not. In-memory or persistent caching is a deliberate
choice, not a default chase: `revalidatePath` after a mutation, `cache: 'no-store'` on the API hop,
and let the router cache do its job (§7d).

**Why `cache: 'no-store'`: every response is session data.** Because `apiFetch` is `no-store`,
Next's request memoization is disabled — two callers of the same URL in one render make two HTTP
calls. The fix is React `cache()` wrapping the reader (`getSession` is the template). This is
why "fetch where it's needed" is unsafe without it: use `cache()` or memoize in a `*-data.ts`
module, never rely on Next's fetch deduplication. The deep detail — RSC/SSR/CSR trade-offs, streaming,
`Suspense` keys, and the `metadata`/caching interactions — is `references/nextjs-production-readiness.md` §3.

## 7d. Preloading, prefetching and the router cache

Next preloads aggressively and that is a feature, not a leak:

- **`<Link>` prefetches** the segment's RSC payload on hover/entry (production) and heavily in dev;
  the router cache then serves the visited route with no round-trip. Do not disable prefetch to
  "save bandwidth" without measuring — and do not *require* it: a deep link, a reload or a hard
  navigation always re-fetches.
- **`/<path>?step=x` is one cache entry per query string**, not one per page. Navigating `upload →`
  `style` re-runs the route when the server is re-rendering it; the wizard is `force-dynamic`, so a
  `?step=` change after a persist re-fetches the page — that is the intended, small cost.
- **`router.prefetch` / `preload`** are for the deliberate, user-expected next step

- **`preload` + React `cache()` + `server-only`** is the safe "fetch where needed" pattern.
  A `preload` function calls a `cache()`-wrapped reader; the reader is `server-only`. Multiple
  callers share one call per request. The `preload` pattern also enables parallel-into-sequential
  reads: the layout calls `preload()` then `await useQuery({ queryFn: theCachedReader })`,
  so the second starts fetching before the component that needs it has even rendered. (a hover-intended
  link). Prefetching the whole tree is an anti-pattern on a mobile-first product.

The interplay of prefetch, the router cache and `revalidatePath` is `references/nextjs-production-readiness.md` §4.

## 7e. Integrations — every outward boundary behind the Next server

The Next server is the single client of your web framework (§4), so integrations stack up here:

- **The API hop** — a Server Component or Action calls `apiEndpoint(...)` with the session cookie
  forwarded (`forwardedCredentials`, `your API package`); the browser never reaches your web framework (your project's architectural decisions).
- **The S3 upload hop** — the room photograph goes **direct to LocalStack/AWS S3** with a presigned
  PUT the proxy mints (§4); the object is never streamed through the app.
- **The generation queue** — accepting a generation enqueues an SQS message; the worker (`apps/ai`)
  loads the row, runs the pipeline and publishes the result. The client's Job is polling the
  status read surface, never holding a queue handle.
- **Images and fonts** — `next/image` with explicit sizes (§5); fonts either `next/font` or the
  tracked `@font-face` set, one mechanism (§6).
- **Third-party** — the upstream system mock in-process, the style catalogue ingest, and any future partner
  calls all sit behind the proxy boundary (integration tests prove the seams — `test`).

Every integration that crosses the wire has a contract in `your contract package` (or `your integration package` for
server-to-server), tested at its seam. A new external hop without a schema behind the proxy is a
feature that has not been specified yet.

## 7f. Access control — asking the question from a route

The model, and the four ways it is got wrong, are
security.md §Authorization
— read that first, and load `security` for a review. `your contract package pathsrc/policy.ts` is the
vocabulary; `your app source pathlib/access.ts` is how a component asks. `ADMIN.md` rule 15 is the
console's copy of the headline: **ask `can(subject, Permission.X)`, never `role === '…'`.**

What is specific to a route, and therefore lives here:

- **Server Components compute; Client Components ask a hook.** `can(subjectOf(session), Permission.X)`
  in a `page.tsx` (`admin/companies/page.tsx` is the pattern), `useUserPolicy()` / `useCan(permission,
  resource)` below a `"use client"` boundary. Both reach the policy through `your API package` — naming
  `your contract package` from the client fails `client-api-boundary`.
- **Pass the decision, not the subject.** A screen receives `canCreate={…}`, not a session it can
  interrogate: `your UI package` composites are presentational, and a policy call inside one is the same
  violation as a `fetch` (§7). It also removes the appearance a pixel baseline could not predict.
- **A refused read is a routing outcome, and each status has a different one.** 401 → the sign-in
  redirect carrying where they were going; 403/404 → `notFound()`; anything else → the segment's
  boundary. A page that renders an empty list for a refusal is telling the reader there is nothing
  there, which is the same defect class as rendering an outage as "you have no projects" (§4).
- **Say at each gate that it is a courtesy.** The comment is load-bearing: without it the next
  reader assumes the client check or the server check is redundant, and removes one.

## 8. Business logic does not live here

Not in a route handler, not in a server action. It belongs in `api/src/services/` as plain functions
with no framework types in their signatures — see `api-contract`. A route handler that grew a
transaction and three branches is a service that ended up in the wrong workspace.

## 9. Responsive is a functional requirement

Capitalised in the brief: **"DESIGNS MUST BE RESPONSIVE - WE WILL BE USING ON TABLETS (IN BOTH
PORTRAIT AND LANDSCAPE), MOBILES AND DESKTOP"**. Treat it as testable behaviour, not polish.

Use `dvh`, never `vh` — mobile browser chrome makes `vh` wrong. The four-viewport matrix and the
1024px tablet-landscape trap live in `design` §2.4 and
`RESPONSIVE.md`; that is the one copy.
