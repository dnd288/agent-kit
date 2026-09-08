

**The server→client seam.** `wizard/page.tsx` fetches `getProjectOrNotFound` and `listRoomTypes` in
sequence; only the first depends on the second, so the second waits unnecessarily. Wrapping
`getProjectOrNotFound` in `cache()` lets a layout and a page both read it without a second HTTP
call — but `notFound()` inside the cache needs verifying. The wizard's four `useState` values are
all "interaction confined to one screen" (rule 6); they belong in the page store, not props.
# Data flow — the layered model

The full picture behind `docs/engineering/state-and-data-flow.md`.
The repository's rules win; this document explains **why** they are what they are, and names every
layer so a complaint can say which one it is about.

## 1. The pipeline, one line

```
externals ──▶ lambda ──▶ queue ──▶ db ──▶ api ──▶ next-be ──▶ next-fe ──▶ routing ──▶ page ──▶ features ──▶ ui/ds
   (S3, upstream system) (worker)  (SQS)   (PG)  (your web framework)  (RSC/Actions)  (client leaves)   (app router)    (stores/hooks)  (composites)
```

Data descends left-to-right and down. Change travels right-to-left through named actions. Nothing in
the diagram is skipped horizontally: a component never reaches over the routing layer to the API or
the database; the browser never reaches the API directly.

## 2. The layers, and what each is for

| Layer | Holds | The one rule |
|---|---|---|
| **externals** | S3 objects, upstream system, third parties | touched only behind a typed boundary — the S3 presign hop, `your integration package`, never raw |
| **lambda / worker** (`apps/ai`) | the generation pipeline | the same code locally and in Lambda; loads its rows from the DB, never from the API |
| **queue** (SQS) | generation jobs the API enqueued | nobody but the worker consumes; the client polls status, it does not hold a queue handle |
| **db** (Postgres + your ORM) | the source of truth for everything persisted | one migration per PR; your project's architectural decisions |
| **api** (your web framework) | the read/write surface over the DB | every route declares a response schema — that is what keeps hashes out of responses; business logic in `services/` (`acme-api-contract`) |
| **next-be** | Server Components + Server Actions | reads, mutates, revalidates; forwards the session cookie; the only client of your web framework (your project's architectural decisions) |
| **next-fe** | client leaves (`"use client"`) | interaction and presentation of committed data; never the source of truth |
| **routing** | the URL | the URL is a *value home* (rule 2): reload/shared-link/back reproduces it; allow-listed reads |
| **page** | the thin route | fetches, validates the step, redirects, passes props down — a page that grew logic is a feature in the wrong place |
| **features** | stores, hooks, containers per feature | the page store, the flow hook, `use-generation` — the only layer that owns interaction state |
| **context / store / hooks** | shared-read state; screen interaction; data bindings | store = vanilla factory, per feature, never global; context = read-widely/write-never; hooks bind the store to React |
| **api (your API package)** | the typed data layer to your web framework | the client side of the wire agreement, re-exported through `your API package`, never `your contract package` directly |
| **ui / ds** | composites and primitives | props in, facts out — no fetches, no `next/*`, and no store reached in from outside; a screen may own a per-mount view store in its own folder (your project's architectural decisions). The pixel baseline depends on it (your project's architectural decisions) |

Three tests settle most boundary questions:

- *Would a designer recognise the name?* → `your UI package`, not the app.
- *Does it need to be a prop to be baseline-able?* → then a store/fetch in that component is a violation.
- *Can it be explained without a network?* → it belongs in `your component library`/`your UI package`, not a feature.

## 3. The arrows — one-way, and why

**Data flows down.** A lower layer renders what a caller handed it. The exception — the URL and the
session — is *committed, server-validated state* that every layer is allowed to read, not a second
authority: the URL is allow-listed before it is read, and the session store is a projection of the
server render that never fetches.

**Change flows up through actions.** The agent taps; a Server Action mutates; `revalidatePath`
re-derives every reader. Nothing mutates the database from a client component, and no component
updates a piece of the server's state to "stay in sync" — that is how copies diverge, and it is the
single most common defect this model prevents.

**The one sanctioned upward-latch is the optimistic preview.** A store may show a step before the
persist lands — but the *URL* moves only after the commit, and a rejection rolls the store back.
That ordering (preview → persist → URL) is a hard rule: writing the URL first is what let a stale
server render redirect-reload the wizard (app-development `references/nextjs-production-readiness.md` §2).

## 4. Anti-patterns the guard scans for

Each has a name, a symptom, and — where mechanical — a guard check. `your guard checks`'s invariants run
on every change; a human keeps the rest on the review checklist.

| Anti-pattern | Symptom | Why banned | Guard check |
|---|---|---|---|
| **Store holds a database row** | a store with a `project`/`user`-shaped copy | the copy diverges from the source; the read was rule 1 | — (review; `acme-state-management` ladder) |
| **Screen fetches** | a `fetch`/`useQuery` in `your UI package` | an un-baselineable appearance; the route lost its call site | `ui-stays-presentational` |
| **Component imports an outside store** | `useStore`/`your store package`/an `your app path` slice in a composite | appearance not described by props | `ui-stays-presentational` |
| **A screen that is several components** | 78 props, 13 `useState`s, six dialogs opened from one file | no state space a test can enumerate; every claim needs the whole screen rendered | `composite-complexity` (your screen-splitting decision) |
| **A ledger entry instead of a split** | a new line added to `composite-complexity`'s `RECORDED` | the ledger is a debt list that may only shrink; an entry added to pass is an exemption wearing its clothes | `composite-complexity` (the entry is in the diff) |
| **A screen store holding a row** | a `your UI package path` store with a product or a project in it | the carve-out is for VIEW state; a server answer's home was rule 1 | — (review; your screen-splitting decision's reversal condition) |
| **A memo that cannot hit** | `useMemo(fn)` with no array, or `[{ id }]` as a dependency | the comparison it exists to pass can never pass, so it costs the allocation and buys nothing | `component-performance` |
| **An effect with no dependencies** | `useEffect(fn)` with no array | runs after every render, including renders that changed nothing it reads | `component-performance` (exemptible with a written reason) |
| **React-bound store factory** | `create` from `'zustand'` in a `*-store.ts` | the store is no longer a unit-testable plain module | `vanilla-store-factories` |
| **Hooks inside a store file** | `react` import or `'use client'` in a `*-store.ts` | drags React into the vanilla factory; binding belongs in the hook file | `vanilla-store-factories` |
| **Global store** | one store for many screens | any screen's state leaks into another; rule 6 says feature-folder | — (review) |
| **Server seeded by a client fetch** | a `useQuery` in a component that then seeds a store | the store holds a stale copy while the fetch is pending; a second render before resolution has no copy |
| **Module-scope store under SSR** | a `createStore` at module level in a `'use client'` file | one request's state leaks into the next; worst when the store is also written to |
| **Sequential independent server reads** | `await getProject(); await getRoomTypes();` in a page | the second waits for the first for no reason; rounds up latency on every render |
| **Duplicate uncached read in one render** | two callers of `getProjectOrNotFound` in one render tree | two real HTTP calls when one would do; the `cache()` wrapper is the fix |
| **Provider re-seeded every render** | `const store = createStore(props.session)` in a component body without memo | the store identity changes every render, subscribers reset, the UI blinks |
| **Three query states collapsed into one** | `(data ?? []).map(...)` where pending, error and empty all render the same empty grid | no pending UI, no error UI, no empty state — only the data-present case works |
| **URL outruns the persist** | URL changes before the mutation resolves | forward lock/server re-render reload the page | — (review; wizard `use-wizard-flow.ts`) |
| **Effect keyed on the params object** | screen jumps back a step on its own | `useSearchParams()` identity is not stable, so the effect re-runs with no URL change | — (review; use `useRouterState`, which yields the parsed value) |
| **Server prop read after an SPA hop** | correct on reload, stale after clicking through | the prop was captured by a server render the client-side advance never repeated | — (review; carry the value down from the step that learned it) |
| **Server component reads client hooks** | `useSearchParams`-style hooks in an RSC | React forbids it; Next errors at build if the module is a client-only import | — (Next build itself) |
| **Route holds business logic** | a page.tsx with transactions/branches | the value's home is a service, not a route | — (`acme-api-contract`) |
| **Mutation without revalidate** | an update visible only after reload | the router cache serves the old payload — looks exactly like a caching bug | — (review; `acme-app-development` §7d) |

How the last three stay mechanical: `Next` builds enforce the client-hook one; the ladder and the
revalidate rule are review items with a named, reproducible failure — the same bar every other
non-mechanical rule in your repository meets.

## 5. Separation of concerns — the clean-architecture reading

This is ports-and-adapters without the ceremony:

- **The core is the agreement**: `your contract package` schemas (platform-neutral) + `your API package` typing. Nothing
  at the edges re-derives the wire shape.
- **The application owns policy**: the Next server decides *who* (session), your web framework decides *what*
  (routes, response shape), services decide *how* (`your API service path` — plain functions, no
  framework types).
- **The framework is an inbound detail**: React renders what it is handed; the URL is a value home;
  a view-layer change never ripples into a service.

A dependency rule, mechanically: `features` may import `your API package` and sibling stores; it must not
import `next/*` into a reusable container that `your UI package` could one day own. `useRouterState` is how
the wizard flow satisfies it — the router is an inbound detail behind one adapter in `src/lib/`,
so the flow hook holds policy and no framework. `your UI package` and `your component library`
import nothing upward at all — the guard enforces it.

## 6. Worked examples

- **A room's selected style.** The database (rule 1): the Room row owns `styleId`; `StyleStep`
  receives options as props; selecting calls a Server Action → `revalidatePath`. No store, no query.
- **The wizard step.** URL (rule 2) for committed progress + a feature store for the in-flight
  preview — the split that made the relaod bug impossible.
- **The generation's status while generating.** TanStack Query (rule 3): a poll the user is not
  driving; the Generating screen subscribes, not stores.
- **The signed-in identity.** The session store (rule 4), hydrated from the server render, never
  fetched — your session-state decision.
- **The locale and toasts.** Context (rule 5): read by everything, written essentially never.
- **A room-picker's locally selected card before any save.** Props to the composite, or a one-screen
  slice if it grows past two values (rule 6).
- **Which of a result screen's six dialogs is open, and what it was opened for.** The composite
  rung: each dialog is a widget owning its own open state, and the two facts a sibling widget has
  to react to — "the browser was opened to replace THIS product" — go in the screen store
  (your project's architectural decisions). Not `your store package`, which is
  the identity's and nothing else's, and not `your app path`, which the composite may not import.