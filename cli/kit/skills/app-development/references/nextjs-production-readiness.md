# Next.js production readiness

The deep material behind `app-development`. Each section states the concept, the rule **this
repository** applies, the trap that has cost time here, and what to read next. The repository's own
decisions win; the Next.js docs are cited for the mechanics, because they move.

## 1. Rendering models — who runs what, where

A Next.js app can render a page four ways. The app uses the first two; the modes are the *file's*
choice, never the route's:

| Mode | Runs in | Use for | This repo |
|---|---|---|---|
| **RSC (Server Component)** | server | fetch, DB reads, secrets, business code | the default — §1 of `SKILL.md` |
| **CSR (client component, SSR'd once)** | server for the first render, then browser | interaction: state, effects, `next/navigation` | `"use client"` leaves, pushed down |
| **SSG / static** | build time | content that never changes per request | only future marketing/landing routes |
| **On-demand / ISR** | server, cached | content that changes rarely but can be invalidated | not used yet; revisit with a real need |

The switch for the page (not the file) is `export const dynamic`:

- `'force-dynamic'` — re-render on every request. **The session- and wizard-bearing routes use this**:
  `getProjectOrNotFound` reads the session cookie and the persisted step per request.
- the default (static) — for a page with no request-dependent data, which this product has almost none
  of yet. Do not sprinkle `force-dynamic`; wait for a demonstrable per-request dependency.

**Trap:** a page that declares static output but fetches in a Server Component with the default
`fetch` cache gets cached HTML that shows another agent's session if the session cookie is read at
revalidate time. The session is per-request; so is the route.

## 2. Routing and navigation — what makes a click become a URL

### 2.1 Soft vs hard navigation

A click on `<Link>` or a `router.push/replace` is a **soft navigation**: Next fetches the RSC
payload and patches the tree — no reload, state preserved. A **hard navigation** is a full
`window.location` load. Four things produce a hard navigation; know each:

1. `window.location` / `<a href>` with a full page (outside `Link`).
2. A `redirect()` returned by a **Server Component during a soft navigation** — Next follows it as a
   full navigation. This is the wizard's reload: the URL named a step ahead of the persisted one,
   the server page re-rendered, the forward lock redirected, the page fully reloaded onto the old
   step. The fix is ordering, not redirects: the URL moves only after the persist lands
   (`features/wizard/use-wizard-flow.ts`).
3. A full navigation from a Server Action (certain `redirect()` paths) — prefer actions that mutate
   and revalidate, letting the client flow own the URL (`advanceWizardStep` deliberately does not
   redirect).
4. A stale module boundary (missing `Link`, a `key` remount) — read as "the page reloaded".

### 2.2 `Link` vs `router`

- `<Link>` is for destinations the markup can name — the sidebar, a card → its project.
- `router.push/replace` is for programmatic changes — advancing a wizard, redirecting after an
  action. `replace` keeps history clean for steps; `push` for a new address.
- A `Link` swallows the click and prefetches by default; a bare `<a>` in client markup opts out of
  both. When a nav item must not prefetch (a heavy destination), use `<Link prefetch={false}>` —
  but measure before defaulting to it.

### 2.3 The URL as state

Per your project's architectural decisions rule 2, any value a
reload, a back button or a shared link should reproduce lives in the URL (`?step=…` in the wizard).
The conventions:

- Read the URL through **`useRouterState`** (`your app source pathlib/use-router-state.ts`), the one
  client module importing `next/navigation`. It takes the allow-list as its `parse` — URL input is
  attacker-controlled, and `parseWizardStepParam` is the template — and returns the parsed value
  rather than the params object, whose identity is not stable across re-renders.
- Writing a URL from a client component goes through `wizardStepRoute`-style pure helpers, never a
  string literal scattered at call sites (and never the prefix — §2 of `SKILL.md`).
- **The screen follows the URL**: a back/forward effect re-anchors `?step=` changes without a
  remount. The wizard's `use-wizard-flow.ts` is the pattern: the store previews instantly, commits
  to the URL only when the server agrees, and still follows the URL when the *browser* moved.

### 2.4 Route groups, parallel and intercepting routes

- **Route groups `(portal)` / `(share)`** split audiences by layout and auth — the portal requires a
  session, the share view must not. Never one layout with a conditional (§3 of `SKILL.md`).
- **Parallel routes `@slot`** render independent segments side by side in a shared layout (a settings
  pane + its detail). **Intercepting routes `(.)`** show a nested view in context (a modal over a
  list) while the URL stays deep. Both are composition tools for later; the design has not asked for
  them yet, and adding one is an `fe` change with a spec.

## 3. Streaming, Suspense and the loading shell

- `loading.tsx` next to a page renders while that segment's server components await — the fallback,
  not the data. It belongs beside the segment it describes.
- An **async server component inside `<Suspense>`** streams: the browser paints the skeleton, the
  slow query lands in a second, independent chunk. This is how a lint-boundary query does not stall
  the hero. Do not `await` two slow queries at the page top if the first paint can explain itself with
  a skeleton.
- `error.tsx` boundaries catch a segment's render errors; the closest boundary wins. A fetch that
  can fail (the API hop, §4 of `SKILL.md`) gets a boundary near where it is consumed, with copy the
  agent can act on — never the default stack.

**Trap:** `Suspense` around a *client* data hook does not stream; a client fetch shows its own
loading. Stream at the server component that owns the `await`.

## 4. Preloading, prefetching, and the three caches

Three caches interlock; know which is which or a "missing update" hides in one:

1. **The router cache** (client) — recently visited segments' RSC payloads. `<Link>` prefetches into
   it; the back button reads it. It is per-session and in-memory; a reload clears it.
2. **The data cache** (`fetch` cache, server) — keyed by URL + options. `cache: 'no-store'` bypasses
   it (the API hop uses this because every call is per-request). `revalidatePath`/`revalidateTag`
   purge it.
3. **The full route cache** (build-time, static routes) — `force-dynamic` opts a route out.

Rules that follow:

- **After a mutation, `revalidatePath` the route(s) that read the value** (`actions.ts` in the app
  does this) — or the router cache serves the old payload and the screen "did not update", which is
  indistinguishable from a bug in your mutation. This is the most common defect in this shape.
- **`prefetch` on hover** is the default for `Link`; it makes the next page feel instant because the
  browser went home with the RSC payload before the click. Keep it.
- **A `?step=` change on a `force-dynamic` page re-fetches the route** — that is the committed cost of
  the URL being the truth. It is small; do not "optimise" it by caching the wizard.

## 5. Integrations — the boundaries

The Next server is the only client of your web framework (your project's architectural decisions),
so every outward hop has one gate:

| Hopped boundary | Mechanism | The rule |
|---|---|---|
| your web framework API | `your API package` typed service → `apiEndpoint()`, cookie forwarded | Server Components & Actions only; browser never (§4 of `SKILL.md`) |
| S3 upload | presigned `PUT` minted by the proxy | direct from the browser; never streamed through the app |
| Generation queue | SQS message enqueued on accept | the worker (`apps/ai`) owns the pipeline; the client polls the status |
| upstream system | `your integration package` signed internal client | server-to-server only, base URL from config |
| Images/fonts | `next/image` sizing; one font mechanism | §5–6 of `SKILL.md` |

Every wire crossing has a zod contract in `your contract package` (or `your integration package` for server-to-server payloads)
and a seam test (`test`) — an integration nobody can prove is an integration nobody should ship.

## 6. Where the code goes — the app's own layering

```
app/(portal)/<route>/page.tsx      server route: fetch, validate, redirect, thin
app/(portal)/<route>/actions.ts    server actions: mutate + revalidate
features/<feature>/                the feature folder: stores, hooks, containers
lib/                               shared client lib: base-path, queries, transition
your UI package source path…                  props-only composites (your component layering decision)
your component library source path…                  primitives + tokens
```

Push work down and left; a route that grows logic is a feature that escaped its folder (`state-management`
owns the full data-flow picture). The route renders, the feature holds interaction, the packages draw.

## 7. Production-readiness checklist

Run this before marking a screen "done":

- **Session-correct.** Server components reading the session forward the cookie; session-bearing
  routes are under the configured prefix; nothing trusts the upstream system cookie.
- **URL-truthful.** A reload reproduces the screen; the URL never runs ahead of the database; URL
  values go through an allow-list; back/forward follow.
- **Loading and error shells.** Every slow segment has `loading.tsx`/`Suspense`; every fallible fetch
  has an `error.tsx` with actionable copy.
- **Caching-aware.** Mutations revalidate the routes that read them; per-request data is `no-store`;
  nothing was cached "for speed" that is actually per-session.
- **Image- and font-honest.** `next/image` everywhere, sizes set, one font mechanism.
- **Keyword-tested.** The four-viewport matrix holds (`design` §2.4); `dvh`, never `vh`.
- **Guard-clean.** `your guard checks` passes, and the anti-pattern checks (`vanilla-store-factories`,
  `ui-stays-presentational`, `client-api-boundary`, …) have nothing to say about your change.