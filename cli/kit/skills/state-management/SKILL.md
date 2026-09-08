---
name: state-management
description: "project state and data-flow conventions — where every value lives, across the whole pipeline. Load this when deciding where a piece of state belongs (server, database, Server Action, URL, TanStack Query, Context, a Zustand store, props), when adding or reviewing a store, when the server→client seam needs a rule (memoized readers, per-mount seeding, parallel fetches), or when a store is being seeded from a client fetch (a page store, a feature store, the session store in your store package), when wiring a hook or a feature container to data, when a value arrives in the wrong place or a screen is not re-rendering after a mutation, or when questioning the layering of the app: the one-way data flow from externals (object storage, the queue, external systems) through the database and the API into the Next server, then down through routing and pages into features — context, stores, hooks — and finally props into your UI package and your component library. Also load it for separation-of-concerns and clean-architecture questions, for splitting an oversized screen into widgets and deciding whether their shared behaviour belongs in a screen store or a screen-local context, for zustand or context inside your UI package path, and for the anti-patterns your guard checks scan for. It complements app-development (the Next-side mechanics) and api-contract (the wire)."
---

# project state and data flow

The one doctrine: **every value has exactly one home, decided by asking one ladder of questions, and
data flows one way — from the edges inward, responses out through actions.** A value whose home is
unclear is a bug in waiting; this skill is how to place it.

The repository's decision record is
`docs/adr/0025-state-management-and-data-flow.md`,
the working copy is
`docs/engineering/state-and-data-flow.md`, and it
is the *authority this skill routes to*. When this skill and the document disagree, the document
wins and this skill is a bug.

## The ladder

Ask in order; the first "yes" decides (state-and-data-flow.md `#the-one-rule`):

1. **The database?** → Server Component reads it, Server Action mutates it, then **revalidate**.
2. **A reload, a shared link or the back button reproducing it?** → **the URL** (read through an
   allow-list; `?step=` is the template).
3. **Changes without the user acting?** → **TanStack Query** (today: the Generating screen's job
   polling).
4. **The signed-in identity?** → the **session store in `your store package`** (your project's architectural decisions) —
   the one scoped reversal of your state-management decision, a projection of the server render.
5. **Read widely, written almost never?** → **React Context** (locale, toasts).
6. **Interaction confined to one screen, more than two related values?** → a **Zustand store in the
   feature folder** — a page store. Never a global store. When the screen is a `your UI package` composite,
   the feature folder is the screen's own folder in `your UI package path` and the store is a **screen
   store** (below).
7. **Otherwise** → **props**, and the value probably belongs to the component above.

**If a store would hold a copy of a database row, the answer was 1.** The one exception is the
identity itself (rule 4).
## The server→client seam

The page fetches on the server; the client renders. Where the seam lives is decided by
state-and-data-flow.md `#the-server-client-seam`
— routing, not restating. The four rules most needed at the skill level:

- **SSR owns starting params.** A store is never seeded by a client fetch.
- **A store is created per mount, never at module scope.** A singleton leaks one request's state into the next under SSR.
- **A server read used more than once in one render is a `cache()`-wrapped reader.** Because `apiFetch` is `no-store`, two callers in one render make two HTTP calls without it.
- **Independent server reads run in parallel.** Sequential awaits of independent reads are a waterfall.
- **Pending, error and empty are three distinct renders.** Collapsing them into one renders an identical empty grid for all three.

## Hydration contract

The session store is the template: server-resolved, passed as a prop, derived synchronously so the
first client render agrees with the server render (no hydration mismatch). See
your project's architectural decisions for the full contract.

**A store that fetches its own initial state is your session-state decision's reversal condition in the wild.** The store
is a projection; it never has a second authority.



## One-way data flow

Data descends through the layers of `references/data-flow.md §2` — never up, never sideways. A lower
layer never reaches over a higher one for a value, and a component never fetches what a caller can
pass. Change travels the other way: the agent acts, a Server Action mutates, `revalidatePath`
re-derives the read side. The screen that receives props does not own the data's life cycle; that is
what makes a pixel baseline able to assert the design (your project's architectural decisions).

## The store contract

A zustand store is a **vanilla factory** — the `your store package` shape: `createStore` from
`zustand/vanilla`, a `create*Store(initial)` factory, React bound separately in a hook file. That is
what makes it unit-testable with no React. The guard check **`vanilla-store-factories`** enforces
the mechanical half — a `*-store.ts` file stays a plain module, whichever kind of store it is:

- it imports no React-bound `zustand` (the binding is `createStore` from `zustand/vanilla`), no
  `react`, and carries no `'use client'` — the hook file owns all three;
- the check is bans-only by design: a session store or a storage client that is nothing to do with
  zustand is simply untouched.

The store holds state and transitions; **routing is not the store's concern.** The wizard's
`use-wizard-flow.ts` is the template: the store previews, the persist lands, and only then does the
flow move the URL (the optimistic-URL race that hard-reloaded the wizard is the anti-pattern this
ordering exists to prevent — `references/data-flow.md §4`).

The URL itself is reached through **`useRouterState`**, the one client module importing
`next/navigation` — and it yields the *parsed value*, not the params object, because an effect keyed
on that object's identity re-runs when the URL has not moved. See
state-and-data-flow.md `#the-url`, which owns
the rule and the failure it prevents.

## The composite rung — state inside `your UI package path`

Rung 6 assumed a feature folder in `your app path`. A `your UI package` screen has none, which for a long
time left it exactly one answer — `useState`, in the screen — and thirteen of them in
`result-screen.tsx` is what that produced. your project's architectural decisions
adds the rung; state-and-data-flow.md `#inside-a-composite`
is the owning document. Ask in order:

1. **Never leaves one widget?** → `useState` there, `useReducer` once three related values move
   together. Nothing is lifted for tidiness.
2. **Two widgets must agree, no transitions, almost no writes?** → a **screen-local context**.
3. **Two widgets must agree, with transitions — or one must subscribe to a field without
   re-rendering on the others?** → a **screen store**: `<screen>-store.ts` beside the component,
   `createStore` from `zustand/vanilla`, bound in `use-<screen>.ts`, provided by the screen.
   The same vanilla-factory contract as every other store, so `vanilla-store-factories` covers it
   unchanged.
4. **Otherwise** → props.

Three constraints are what make this a carve-out rather than a hole in
your project's architectural decisions:

- **View state only** — which dialog is open, what it was opened for. A store holding a domain entity
  or a server-side resource means the answer was rule 1, exactly as before.
- **Seeded from props, per mount** — so a story still describes every state the screen can be in.
- **Local by construction** — created by the screen, unreachable from outside it. `your store package`, an
  `your app path` feature store and a query hook stay banned; `ui-stays-presentational` is unchanged
  and still fails on them.

**How you find out you needed it:** `your guard checks`'s `composite-complexity` budget — 24 props, 6
stateful hooks, 800 code lines, complexity 45. A stateful-hook count over budget is a census of the
widgets the screen has not been split into. `ui-development` §9 owns the split itself.

## Where to look when a value is wrong

Follow the value's home: a database row → the route/action that reads/writes it; a URL param → the
page that parses it and the flow that wrote it; a store → the slice and the server render it
projects; props → the caller. `docs/engineering/state-and-data-flow.md` has the directory map; the
full model is `references/data-flow.md`.

## Where to look when the page is wrong — failures, by layer

The error-handling architecture is a data-flow concern too: `state-and-data-flow.md`'s
§Failures, by layer
owns it. The three lines that matter here:

- **Classification happens where the typed error is visible** (the server fetch site), never in a
  client `error.tsx` — Next serializes server-render errors, so the class never crosses.
- **A down API is not signed out** (`getSession` rethrows `ApiUnavailableError`; the portal gate
  renders the unavailable state instead of redirecting) and not an empty wall (the wall page
  classifies instead of soft-failing).
- **Across a Server Action the classification is the `{ unavailable: true }` marker**
  (`your app source pathlib/api-unavailable.ts`) — never a message to match.

When a screen "renders wrong when the API is down", the fix is at the fetch site, not in the
boundary component.

## Read next

- `references/data-flow.md` — the layered pipeline, every layer's responsibility, the anti-pattern catalogue.
- `docs/engineering/state-and-data-flow.md` — the
  ladder and the directory map (owning document).
- `app-development` §7a — the same ladder from the Next side, with the mechanics.
- `api-contract` — the wire the values cross.
- `your state-management decision`, your project's architectural decisions, your project's architectural decisions — the reasoning behind the rules.