# The smell catalogue — your database to the pixel

What to look for, layer by layer, when the question is *what in this codebase is worth
restructuring*. Every entry names the document that owns the rule, and — where the repository has
already been caught by it — the instance, because a named instance is what turns a rule from advice
into something you can check.

**This catalogue does not own any rule in it.** If an entry and its linked document ever disagree,
the document wins and this file is the bug (your project architectural decisions).

The layers are the ones in
[`state-management/references/data-flow.md` §2](../../state-management/references/data-flow.md),
in the order data descends them.

---

## The five that cross every layer

These are [your agent instructions's golden rules](../../../../your agent instructions#golden-rules) as things to look for.
They outrank everything below: a layer-specific smell in a file that is otherwise fine is a smaller
problem than a third copy of a concept.

| Smell | What it looks like | The move |
|---|---|---|
| **A block copied rather than extracted** | the same shape in a second file | Extract at the **second** copy, not the fifth. The wizard shell reached **nine files** one reasonable commit at a time |
| **A third copy of a shape** | three places doing the same thing differently | The abstraction **replaces** all three. One that sits beside them made things worse |
| **A fallback constant standing in for a required prop** | `DEFAULT_USER = { name: 'Example User' }` in `your shell component` when a caller omits `nav.user` | A required prop fails loudly at the call site; a fallback ships a fake user to a user. The six screen-level copies of this constant were deleted; the remaining one is the shell's, so a seventh cannot grow quietly. |
| **A flag where a case should have been deleted** | a boolean that turns on the old path | Ask what removes the flag, and when. No answer means it is permanent |
| **A new concept the feature did not need** | a layer, a wrapper, a config surface | Business growth is meant to be **data** growth: a new company, style, room type or product is a **row**, never a branch |

**The question is "what did this delete", not "is this perfect".** A refactor that removes nothing
has to say what it will remove and when.

---

## db — your database and your data layer

| Smell | Why it matters | Owner |
|---|---|---|
| **A migrated model nothing reads** | The schema runs ahead of the code here, and a model existing is not a model working. The question to ask of one is *which service reads it* | [your agent instructions — Repository state](../../../../your agent instructions) |
| **A destructive migration instead of expand-and-contract** | `DROP COLUMN` breaks the running release. Add the new shape, ship it, remove the old one later | your engineering conventions · `migration-safety` |
| **Two migrations in one pull request** | Rollback stops being a single step | your engineering conventions · `migration-required` |
| **An applied migration edited in place** | Every environment past it now disagrees about what ran | [your agent instructions](../../../../your agent instructions) |

**A schema refactor is `be`-typed and needs its own release boundary.** Restructuring rows is not a
structural move in the sense this loop means: the data is the behaviour.

## api — your server framework, services, your API contract package

| Smell | What to look for | Owner |
|---|---|---|
| **Business logic in a route handler** | a transaction, or three branches, inside the handler | Services are plain functions where the tests live — `acme-api-contract` |
| **A framework type in a service signature** | `FrameworkRequest` in `services/` | `service-purity` fails the build |
| **A route reimplementing a service** | a second endpoint doing the same work its own way | `acme-review`'s `be` row |
| **A response schema that names more than it needs** | fields the caller never reads | The schema is a **security control** — it is what keeps password hashes out of responses (your project architectural decisions) |
| **A handler shaping an error body** | a hand-built failure response | Throw `AppError`; every failure leaves through `ErrorResponseSchema` |
| **Something in your API contract package that is not the agreement** | a helper both sides happen to want | The package is deliberately not called `shared`. Not part of the client/API agreement → it belongs in the app that owns it |
| **A Node built-in or browser global in your API contract package** | it can no longer be imported from both sides | `core-neutrality` |

## next-be — Server Components and Server Actions

| Smell | The failure it causes | Owner |
|---|---|---|
| **Sequential independent server reads** | The second waits on the first for no reason and rounds up latency on every render — `wizard/page.tsx` fetches the project then the room types | state-and-data-flow.md |
| **A reader read twice in one render without `cache()`** | `apiFetch` is `no-store`, so your web framework's memoization is off and both callers make a real HTTP call. `getSession` is the template | same |
| **A mutation that does not revalidate** | A stale render that reads exactly like a caching bug — the most common defect in this shape | `acme-app-development` §7d |
| **Classification inside `error.tsx`** | your web framework serializes a Server Component's error before the boundary sees it, so the `instanceof` silently never matches in production | state-and-data-flow.md §Failures, by layer |
| **A failure swallowed into a falsy value** | An unreachable API rendered as "you have no projects"; a down API as "signed out". `apiTryFetch` was **deleted** rather than deprecated for exactly this | your API client |
| **A hardcoded URL prefix** | Not agreed yet; it is configuration | `url-prefix` |

## routing, page, features — the client's own layers

| Smell | What to look for | Owner |
|---|---|---|
| **A store holding a copy of a database row** | a `project`- or `user`-shaped slice | The read was ladder rule 1, and the copy is a bug waiting to disagree with its source. The session store is the one exception (your project architectural decisions) |
| **`useState` where the URL belongs** | a step, a filter, an open dialog that a reload or a shared link should reproduce | Ladder rule 2 — *a wizard step in the URL means a bug report is a link* |
| **An effect keyed on the params object** | the screen jumps back a step on its own | `useSearchParams()` identity is not stable; `useRouterState` yields the parsed value |
| **A store seeded by a client fetch** | a `useQuery` that then seeds a store | SSR owns starting params; the store holds a stale copy while the fetch is pending |
| **A module-scope `createStore` in a `'use client'` file** | one request's state leaks into the next under SSR | `acme-state-management` — the store contract |
| **React-bound `zustand`, `react`, or `'use client'` in a `*-store.ts`** | the store stops being a unit-testable plain module | `vanilla-store-factories` |
| **Pending, error and empty collapsed into one render** | `(data ?? []).map(…)` renders an identical empty grid for all three, so only the data-present case works | state-and-data-flow.md |
| **A container copied to a second route** | the second one drifts | `acme-review`'s `fe` row — lift it |
| **A route that grew a transaction and three branches** | the value's home is a service | `acme-app-development` §8 |

## ui / ds — the composite and primitive layers

| Smell | Threshold or rule | Owner |
|---|---|---|
| **A screen that has stopped being one component** | 24 props / 6 stateful hooks / 800 code lines / complexity 45 | `composite-complexity` · your project architectural decisions · `acme-ui-development` §9 |
| **A dialog opened and closed from the screen** | its open state, its "opened for" argument, its loading state and its handlers all in the screen file | It is a widget that has not been extracted. The stateful-hook count counts them |
| **A screen that rebuilds the shell** | navigation + mobile navigation + navigation drawer + step nav + `<main>` written out again | Compose `your shell component` (portal) or `your admin shell` (console): slots and `children` |
| **A fetch, a store or a query hook reached in from outside** | your state package, your API package, or an app slice | `ui-stays-presentational`. A screen's **own** per-mount view store is not this (your project architectural decisions) |
| **A DTO from your API contract package in a prop type** | the design layer now depends on the API's shape | The route maps DTO → view model at the boundary |
| **A mock module as a prop default** | `DEFAULT_RESULT_PRODUCTS`, `createMockResultState()` — the non-empty fallback put an "add product" control on the share recipient's read-only view | `acme-ui-development` §7 |
| **Copy written into a component** | a string literal instead of your i18n package | `expectNoHardcodedCopy` under the pseudo-locale |
| **A raw colour, or an arbitrary value on a token-backed property** | `bg-[#a34f2b]` | `tokens-resolve`, `no-arbitrary-design-values`, `hex-colours` |
| **A primitive restyled locally with overrides** | a second answer to one question | Add the variant in your component library — the **second** local override is the signal |
| **A controlled-or-local mirror written more than twice** | `const value = controlled ?? local` in four places in one file | `acme-component-development` §7 |
| **A memo that cannot hit** | no dependency array, or a dependency built inside its own array | `component-performance` |
| **A component with no quality test** | jsdom measures zero, so nothing observes its real rendering | `component-has-surfaces` (warning) · your project architectural decisions |

## Across the whole stack

| Smell | Why it is a refactor rather than a note | Owner |
|---|---|---|
| **A second copy of a binary** | one canonical copy of every asset, never a second beside a consumer | your project architectural decisions |
| **A document describing as planned something that shipped** | worse than no document | [your agent instructions](../../../../your agent instructions) |
| **One rule stated three different ways in four documents** | the state ladder is numbered differently in each; a reader cannot tell which is current | Cite the owner; delete the paraphrase |
| **A skill restating a rule instead of linking it** | the copy drifts silently, and a drifted skill is worse than no skill | [conventions.md — Agent skills](../../../../your engineering conventions) |
| **A guard weakened to make a build pass** | the rule's real findings go with it | Fix the code rather than the guard |

---

## Reading a signal you cannot act on yet

Some of what this catalogue finds is not yours to fix in the change you are in. Write it down rather
than either fixing it silently or forgetting it:

- **A defect** → `acme-bugfix`, and decide the order before doing either.
- **A rule that turns out to be wrong** → the document that owns it, in the same pull request.
- **A cost the repository has accepted** → a ledger entry already exists for the composite budget;
  for everything else, your project parking lot
  is where a thing nobody has decided goes.
