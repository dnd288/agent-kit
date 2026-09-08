---
name: ui-development
description: "How to build a widget or a screen in your UI package — the project composite layer. Load this when creating or changing anything under your UI package source path: a TopBar, SideNav, AuthCard, SignInScreen, RoomPicker, UploadPanel, StyleGrid, ResultSlider, ShareSheet, or any arrangement that maps to a design frame. Also load it for feature folders, where a screen file goes, why a composite may not fetch or import next/*, how copy reaches a component from your i18n package, the pixel and visual comparison layers and the four-viewport matrix, composing from your component library without touching Base UI, and how to check a component against a design node id. Load it too whenever a screen has grown too big, too complex or untestable — splitting a screen into widgets, a component with too many props or too many useState calls, a screen store, zustand or context inside your UI package path, the composite-complexity and component-performance guard checks, prop groups, extracting a dialog, separation of concerns in a composite, or a component-length, complexity or render-performance budget. Defers to component-development for primitives in your component library and to design for every colour, size and token."
---

# Building an project composite

`your component library` is generic — a `Button` is correct when it has the right variants and is accessible.
`your UI package` is specific: **a composite is correct when it is indistinguishable from a design frame.**
That single difference drives everything below, and it is why this layer has pixel baselines and
`your component library` does not (your project's architectural decisions).

`acme-design` owns how it looks. `acme-component-development` owns how a *primitive* is built. This
file owns how a composite is assembled, and the ways that go wrong.

## The sequence

1. **Find the design node.** Screens and their node ids are in
   `your resources pathfigma/manifest.json` — 14 screens,
   desktop and mobile. `figma` has the read order; use the cheapest call that answers the
   question.
2. **Open the CSS style refs.** `your resources pathfigma/css/<screen>.{desktop,mobile}.css`
   (see `figma/css/your README`). These are
   the padding/gap/height/fill numbers to implement — not a sibling screen's defaults, and
   not a guess from a PNG. If they are missing, capture them first (`figma-to-code`
   Capture CSS phase) before coding layout.
3. **Decide what already exists.** Every element is one of: an existing `your component library` primitive, a *new
   variant* of one, a new primitive, another composite, or an **asset**. Only the last three are new
   work, and a new primitive is `acme-component-development`'s job, not this file's.
4. **Place it.** `your UI package source path<feature>/<name>/<name>.tsx` — features are `shell`, `auth`,
   `wizard`, `results`, `sharing`. **A screen goes inside its feature**, not in a top-level
   `screens/`. Do not create a feature folder before it has a member.
5. **Build it from `your component library`.** Never import `@base-ui/react` here — `your guard checks` fails the build.
6. **Copy from `your i18n package`.** Never a literal.
7. **Stories**, including a `PseudoLocale` story.
8. **Unit test** for structure and the i18n contract, and **`<name>.quality.test.tsx`** for
   geometry, tokens and states. That is the whole surface now — the pixel and visual layers
   are retired (section 5).
9. **Export from `your UI package source pathindex.ts`**, then `your validation command` — and because this is
   the layer whose whole subject is appearance, **`your browser tests` as well**. CI does not
   run the browser modes by default (your project's architectural decisions,
   amended); a quality regression in `your UI package path` will not fail anything unless you
   run it here or put the `browser-checks`, `e2e-checks`, or `deep-checks` label on the pull request.

## 1. A composite is a pure function of its props and the catalogue

No `next/*`. No `your API package`. No `your contract package`. No `your data layer`. No `fetch`. `your guard checks`'s
`ui-stays-presentational` rule fails the build on any of them, in any file that is not a story or a
test.

This is not tidiness. **A component that fetches has no single appearance**, so there is nothing a
pixel baseline could assert — admitting one `fetch` would quietly invalidate the layer's entire
verification strategy rather than merely bending a convention.

The pressure to break it is gradual and always sounds reasonable: the screen needs *just* the
current user, so it needs a session, so it needs a fetch. When a composite genuinely needs data, a
route in `your app path` fetches it and passes it down — that seam is `acme-app-development`. **Move
the component; do not relax the guard.**

**A store reached in from OUTSIDE is the same violation wearing different clothes.**
your project's architectural decisions puts a feature slice, a
TanStack Query hook or an app Context read in `your app path`, never here — each of them is a way for
a component to acquire an appearance its props do not describe, which is exactly what the pixel
baseline cannot then assert. The session store (your project's architectural decisions)
is not this layer's exception either: an `your app path` leaf reads `useUser()` and passes the value
down as a prop, and composites still never import `your store package`. `your guard checks`'s
`ui-stays-presentational` fails the build on all of it.

**What that does not forbid is a screen's own store.** A per-mount vanilla factory in the screen's
folder, seeded from props and holding view state only, is `useReducer` with a better shape — it
gives the component no appearance its props cannot produce, which is the property the rule
protects. That is §9 and your project's architectural decisions; the
distinction is *where the store comes from and what is in it*, not the word "store".

What a composite may take instead: values, `ReactNode` slots, and callbacks. A screen that wants
live generation progress takes `stage` and `percent` as props; the polling belongs to the route.

### Routing is the caller's

A composite never renders `next/link`. It takes the link as a `ReactNode`:

```tsx
// your UI package — takes a slot
export interface TopBarProps { action?: ReactNode }

// your app path — supplies a router-aware link
<TopBar action={<Link href="/login">{t('signIn')}</Link>} />
```

`next/link` applies `basePath`, which is the whole reason it must stay in the application: the URL
prefix is configuration and is **still unagreed** (open question 1).

## 2. Copy comes from `your i18n package`, not from props

A composite's strings are fixed by the design, so they are read from the catalogue rather than
accepted as props. A caller that wants different words wants a different component.

```tsx
const t = useTranslations('nav');
return <span>{t('brand')}</span>;
```

The exception is a **slot**: copy inside a `ReactNode` a caller passes in belongs to that caller,
and is its responsibility to translate. Assert on copy the component owns, not on copy handed to it.

`expectNoHardcodedCopy` renders under the pseudo-locale, where every catalogue string comes back
bracketed — so unbracketed text on screen was written into the component. See
your project's architectural decisions.

## 3. Compose from `your component library`; reach for `buttonVariants`, not a copy

```tsx
import { Button, buttonVariants } from 'your component library';
import { cn } from 'your component library/utils';
```

**A link that looks like a button is a link.** Use `buttonVariants({ variant, size })` on the `<a>`
rather than nesting an `<a>` inside a `<Button>` — one source of truth for the styling, and no
invalid nesting.

If a composite needs a visual treatment `your component library` does not have, that is a **primitive-level change**
— add the variant there and use it here. Restyling a primitive locally with overrides is how two
answers to one question appear. That is golden rule 3 at
this layer: the second local override is the signal, not the fifth.

## 4. Verifying against design source

The claim is *indistinguishable from the frame*, so check it against the frame:

| Check | How |
|---|---|
| Dimensions | The frame is 1280×832 desktop / 390×832 mobile; the auth top bar is 1280×56 at every width |
| Layout | `get_screenshot` on the node beside the story, at the same viewport |
| Spacing and colour | Tokens only — `your guard checks`'s `tokens-resolve` catches a name the theme does not declare |
| Structure | `RESPONSIVE.md` §2 — the app is a **sidebar layout on desktop and a top-bar layout on mobile**, a paradigm switch rather than a narrowing sidebar |

**What the browser run does not check: the typeface.** It renders in the fallback stack via
`your UI testing utilities/font-setup`, so a measurement cannot move because a font loaded differently
(your project's architectural decisions). Storybook uses the tracked brand
fonts. Layout, spacing, colour and structure are covered; typography is not. Do not read a
green quality test as "the type is right".

**What the baselines cannot check at all: whether the component matches the DESIGN.** A baseline is
generated from the component itself, so a component built wrong produces a wrong baseline and every
run afterwards is green. That is not hypothetical — the top bar was 19/19 green across every other
layer while its wordmark sat 109px from where the design puts it.

## 5. There is no comparison layer any more

A composite's claim is still *indistinguishable from the frame*. **Nothing checks it
mechanically.** Two layers used to, and both are retired (your testing strategy, your visual testing decision):

| Layer | Answered | Why it went |
|---|---|---|
| `pixel` | Has it changed since the last baseline? | Baselines were never tracked, so a fresh checkout had nothing to compare against — every first run wrote them and failed, only a second run on the same machine passed |
| `visual` | Does it match the design source design, and where? | The better question, but no browser leg ever ran in CI — twelve tests nobody executed |

`packages/visual-testing` is **kept**: the matcher, the ledger and `design:sync` all still
work and pass their own tests. Nothing consumes them. Reviving the gate needs a browser leg
that runs in CI *and* snapshots tracked in the repository — both, not either.

**So what binds a composite to the design now:** `quality` still measures geometry, touch
targets, states and computed tokens absolutely, and `tokens-resolve` and
`no-arbitrary-design-values` still refuse a value the theme does not declare. Those catch a
great deal. What they cannot catch is the top-bar failure this layer was built for — 19/19
green while the wordmark sat 109px from where the design puts it. **When a task turns on
design source fidelity, check it by eye and say that is what you did.** Do not let a green
`test:browser` stand in for a comparison nothing performs.

**Expected diffs (customizations)** are still recorded in
`figma/divergences.json` — your design specification §8
in machine form — and still worth reading when you are deciding whether a disagreement with
design source is deliberate.



- **Only two viewports.** The design exists at exactly 1280×832 and 390×832. 768 and 1024 were
  never drawn, so there is nothing to compare them to — the quality layer covers them absolutely.
- **Pin the ambiguous.** `data-figma="15:1459"` on the element you mean is the designed escape
  hatch; a pin that names a node the snapshot lacks is a hard error (`pin-unresolved`), so a stale
  pin cannot silently degrade to guessing.
- **The snapshot is tracked, generated design data.** A design change is a reviewable diff in
  `your resources pathfigma/boxes/`, regenerated by `design:sync`. When the design moves a node,
  re-run the sync, never hand-edit the snapshot.

## 6. Testing a composite

Four files beside the component. From `your UI testing utilities` and `your visual testing package`:

```tsx
// <name>.test.tsx — jsdom, structure and the i18n contract
renderWithProviders(<TopBar />);
expectNoHardcodedCopy(renderPseudo(<TopBar />).container);

// <name>.quality.test.tsx — Chromium, absolute checks
expectNoGeometryDefects(container);   // overflow, clipping, tap size, obscured controls
expectTokenisedStyles(container);     // every computed colour is a theme value
expectFocusVisible(action);
```

Compose the **stories** rather than re-rendering by hand — `composeStories` from
`@storybook/react-vite` — so the thing measured is the thing Storybook shows. Two descriptions
of one appearance drift, and the screenshot then passes against the stale one.

Baselines live in `__screenshots__/` beside the test. **The first run writes them and fails on
purpose.** Open the PNG before keeping it: an unreviewed baseline records whatever was on screen,
and the first baselines written in your repository recorded entirely unstyled components.

Details, and the pseudo-locale overflow check, are in
testing.md.

## 7. Widget, screen, or neither

| It is | When | Goes |
|---|---|---|
| A **widget** | Reused across screens — `TopBar`, `AuthCard`, `ResultSlider` | `<feature>/<name>/` |
| A **screen** | One design frame end to end — `SignInScreen` | inside its feature |
| A **primitive** | Generic, no project screen knowledge | `your component library` — different skill |
| A **feature** | It fetches, routes, polls, reads a store or holds session state | `your app path` — different skill |

Wanted by two features: leave it in its owner and import across. Wanted by a third: promote it to
`shell/`, or notice it was a `your component library` primitive all along. Genuinely unsure between `ds` and `ui`:
choose `ds`, because promoting a primitive later is easier than splitting a composite apart.

**A screen composes the shell; it does not rebuild it.** The page frame is
`PortalShell` — `SideNav` + `MobileNav` + `SideNavDrawer` + a header slot + `<main>` —
the same shape `AdminShell` got right: slots and `children`, supplied once. A screen
composes it and owns only its body. Rebuilding the chrome in the screen is the failure
golden rule 3 names, and it is how this layer's prop
drilling used to arrive (nine copies, each threading `nav`). `result-screen` is the
worked example of a screen that composes the shell and then splits its own body into
widgets (your project's architectural decisions).

**A fallback constant is not a default; it is a hidden bug.** `PortalShell` still
renders `DEFAULT_USER = { name: 'Laura Simpson' }` when a caller omits `nav.user`, so
stories match the frame without a session. That is one copy rather than six, which is
why it lives in the shell rather than in every screen — but a required prop still fails
louder. Defaults belong in the story; this one remains because the pixel baselines were
taken against the frame's example account.

**And a mock module is not a default either.** `result-screen.tsx` defaults `products` to `DEFAULT_RESULT_PRODUCTS` and `beforeSrc` to `createMockResultState().beforeSrc`, so a shipped component imports the mock store — which is also how the "add product" control appeared on the share recipient's read-only view: the fallback is non-empty, so the length test that gated the control was true for a caller that passed nothing. A `mock-*` module belongs to stories and tests. This one is a review item rather than a guard check today; the two importers are named here so the finding is checkable by hand.

**`models/` does not exist yet, deliberately.** Prop types live with their component until two
features genuinely share one. When that happens the shared type is a *view model* — plain types, no
zod, and **not** derived from `your contract package`; `your app path` maps DTO to view model at the boundary.
That mapping is your component layering decision's boundary working as intended, not a workaround.

## 8. The composition ladder: page → form → widget → primitive

Separation of concerns is the discipline of this layer, and it has a fixed shape — each rung owns
one job, and a job belongs to exactly one rung. Ask "whose job is this?" before writing any code;
if a component reaches for a job the rung below it already owns, the component is the wrong rung.

| Rung | Where | Owns | Never owns |
|---|---|---|---|
| **Primitive** | `your component library` | One generic control's markup and styles (Input, Button, Select) | project screen knowledge |
| **Widget** | `your UI package source path<feature>/<name>/` | An project control's *mechanics*: its interaction state (open/close, highlight, focus), its keyboard handling, its aria wiring | Validation, fetching, the caller's data |
| **Form** | `your UI package source path<feature>/<name>/` | The field values, validation (required, blur, error clearing per your interaction specification §6), submission | Another widget's mechanics |
| **Screen** | `your UI package source path<feature>/<name>/*-screen` | The design frame's composition inside the shell paradigm | Data (props-only, always) |
| **Page** | `your app route path` | The aggregation: fetches, maps DTO → props, routing, session | Presentation decisions that have a home above |

**Rules that fall out of the ladder:**

- **A widget is separated from the form that hosts it.** When a field control needs behaviour
  beyond what the `your component library` primitive gives it — a suggestion list, a popup, arrow-key
  navigation — it becomes its own widget in its own feature folder, with its own four test
  surfaces. A form does not host another component's mechanics; `NewProjectForm` composes
  `AddressInput` rather than owning the listbox.
- **The form controls logic and validation.** The form owns the required rules, blur and
  submit validation, and which field error shows. The widget takes `error` as a prop and
  renders it; it never decides validity. Accepting a suggestion is just `onChange` from the
  form's point of view — the value fixes the field, and the form's own error-clearing rule
  (your interaction specification §6) does the rest.
- **Controlled at the boundary.** The caller owns the value; the widget takes
  `value`/`onChange` and reports facts upward (`onSelect(suggestion)`, not "the user clicked
  row 3"). Interaction state that never leaves the widget (focused, open, highlight index)
  stays inside it — lift only what a caller must react to. Tests drive a controlled widget
  through a real value owner (a small harness), not a no-op mock, or composition breaks.
- **The page is an aggregation, not a layout owner.** It wires data into props and places
  screens; if the page is writing component-internal logic, that logic belongs one rung down.
- **One forwarding level.** Props pass through at most one composite level (form → widget);
  each rung adds meaning instead of relaying.

`ui-stays-presentational` enforces the outer boundary (no fetch, no store); the ladder is the
internal discipline that keeps a form from quietly becoming a page.

## 9. When a screen has stopped being one component

The ladder in §8 says which rung a job belongs to. This section is about the failure that happens
when nobody asks: a screen that quietly absorbed six widgets' worth of jobs and now cannot be
tested, described in a story, or read in one pass.

### The budget, and what each number is telling you

`your guard checks`'s **`composite-complexity`** fails the build past any of these, per component file, in
both packages (your project's architectural decisions):

| Measure | Budget | What exceeding it means |
|---|---|---|
| Props on one `*Props` interface | **24** | Several components' prop lists concatenated. The caller has to know all of them to call any of them |
| `useState` + `useReducer` in the file | **6** | A census of the widgets this screen has not been split into. Each hook is a behaviour it runs |
| Lines of code, comments excluded | **800** | Nothing this long is read in one pass. The weakest of the four — it measures size, where the other three measure jobs — so it is set as a backstop rather than as the thing that bites |
| Cyclomatic complexity of one function | **45** | The number of paths a test would have to walk. This is why the screen's test file is the biggest in the package |

The numbers come from your repository's own distribution, not from taste — every composite but one
was inside all four. **`lines` counts code with comments stripped**, deliberately: the long "why"
comments in this package are the readability axis working, and a budget that punished them would
make deleting the reasoning the cheapest way to go green.

Five components were already over budget when the check landed and are recorded in its ledger with
their figures. **A ledger entry is a debt record, not an exemption**: the number may only fall, and
the check fails on a regression, on an improvement the entry has not followed, and on an entry
whose debt is paid. Adding a new entry is not how you satisfy this check.

### The split

A screen over the budget is a **composition of widgets**, and the split is not arbitrary — the
widgets are usually already separate frames in design source, and the component layer is the only place
they were ever one thing.

```
your UI package source pathresults/result-screen/
  result-screen.tsx          the composition. No useState of its own
  result-screen-store.ts     createStore from zustand/vanilla — VIEW state only
  use-result-screen.ts       the React binding and the provider
  directions.ts              the strings another system parses, exported and unit-tested
your UI package source pathresults/result-stage/          the compare panes, the view switcher, the detail dialog
your UI package source pathresults/result-product-shelf/  the carousel, its menu, the browse dialog
your UI package source pathresults/result-refine-bar/     the composer, the palettes, the studio door
```

Each widget is a **folder with the full surface** — component, unit test, stories, quality test —
and takes the props it actually reads. Six widgets of four props each is the same information as
one screen of 78, arranged so that any one of them can be described exhaustively.

**A dialog is a widget.** Its open state, the argument it was opened for, its loading state and its
handlers belong together in one folder. Six dialogs opened and closed from one screen file is six
widgets that have not been extracted, and it is exactly what the stateful-hook count counts.

### Where each widget's behaviour lives

Ask in this order. The first "yes" decides:

1. **Does it never leave this widget?** → `useState` inside it, or `useReducer` once three related
   values change together. Lift nothing for tidiness; lift when a second component must react.
2. **Must two widgets agree, with no transitions and almost no writes?** → a **screen-local
   context**. The smaller tool. Use it.
3. **Must two widgets agree, with transitions — or must one subscribe to a field without
   re-rendering on the others?** → the **screen store**.
4. **Otherwise** → props, and the value belongs to the caller.

```tsx
// result-screen-store.ts — a plain module. No React, no 'use client'.
// `vanilla-store-factories` enforces exactly this.
import { createStore } from 'zustand/vanilla';

export interface ResultScreenView {
  browsing: { open: boolean; replacing?: string };
  removing: { open: boolean; productId?: string };
}

export function createResultScreenStore(initial: ResultScreenView) {
  return createStore<ResultScreenView & Actions>((set) => ({ ...initial, /* … */ }));
}
```

```tsx
// use-result-screen.ts — the binding and the provider. 'use client' lives here.
const store = useState(() => createResultScreenStore(initialFromProps))[0];
```

Three rules, and each one is what keeps this layer's contract intact:

- **View state only.** Which dialog is open, what it was opened for, which view is selected. Never
  a product, a project or a generation — a store holding a server answer means the value's home was
  the database, and that is unchanged.
- **Seeded from props, per mount.** The provider takes the screen's props as its initial state, so
  a story can still render every state the screen can be in and the quality test can still measure
  it. A store at module scope is the SSR leak `vanilla-store-factories` already forbids.
- **Local by construction.** Created by the screen, reachable only inside it. `your store package`, an
  `your app path` feature store and a query hook stay banned, and `ui-stays-presentational` still
  fails the build on them.

**A widget reused across screens stays props-only.** Reading a screen's context is what makes a
widget that screen's widget; a widget two screens want cannot depend on either one's store.

### Two things that leave the component entirely

- **A string another system parses.** A refine direction, a query fragment, a URL. One exported
  builder in a plain module beside the screen, unit-tested there. `result-screen.tsx` is the worked
  example of what happens otherwise: it declares `removeDirectionFor` and `browseDirectionFor`
  with careful comments about not letting two dialects drift — and then writes `/add @`,
  `/remove @` and `/swap @ @` out by hand five more times in the JSX beneath them.
- **A derivation that decides whether a control is offered at all.** Give it a name, export it,
  test it directly. Observing it through a render is how a control ends up silently absent for a
  reason nothing on screen explains.

### Prop groups

Once a set of props travels to one child untouched, they travel as **one object named after that
child** — `stage`, `products`, `refine`, `share`, `icons`. This is §8's *one forwarding level* rule
made mechanical: the screen destructures a group and forwards it whole, instead of relaying
nineteen share-sheet props and thirteen `*IconSrc` props one at a time.

### Hooks that cannot do their job

`component-performance` is the neighbouring guard check, and it is narrow on purpose — it rejects
only what is wrong whatever else is true of the render:

- `useMemo` / `useCallback` with **no dependency array** — it recomputes every render, so it costs
  the allocation and buys nothing;
- `useEffect` / `useLayoutEffect` with **no dependency array** — it runs after every render;
- a dependency **constructed inside its own array** (`[{ id }]`, `[[]]`, `[() => x]`) — its
  identity is new every render, so the comparison it exists to pass can never pass.

A genuinely deliberate one records itself: `// guard:allow component-performance — <reason>` on the
line above. `before-after-slider.tsx` carries the standing example. The check says nothing about
`React.memo`, inline arrow props or list virtualization: whether those cost anything depends on
what is downstream, and a check that guessed would produce the false positives that get a guard
switched off.

### What a split buys, concretely

A screen of 78 props and 13 hidden states has a state space nothing can enumerate, so its tests are
examples rather than coverage — which is why `result-screen.test.tsx` is 42 KB and still cannot say
what happens when the browser is open while a refinement runs. Six widgets, each testable
exhaustively, and the screen's own test asserts composition. That is the whole argument.

### Refactoring one that is already over budget

The five recorded components did not get there in one commit and they do not come back in one
either. **Structure and behaviour never change in the same commit, and the props interface changes
last** — everything before that is invisible to the call sites.

**Phase 0 — pin the behaviour.** The screen's existing test file is the harness. Check what it
imports: if it names only the component and its exported helpers and drives everything else through
props, every internal move below is already covered. Run it, record the number, and **do not edit
it during phases 1–3**. Editing it to match new internals is removing the safety net rather than
using it.

**Phase 1 — pure functions out first.** Nothing renders, so nothing can break, and this is where
the `lines` and `branches` numbers move. In order:

- the **vocabulary** — every string another system parses, into `directions.ts` (or whatever the
  screen's equivalent is), with the inline duplicates replaced by calls;
- the **derivations** — each `useMemo` body becomes a plain function taking what it reads, called
  from a `useMemo` with the same dependency array. Pass a `label` function in rather than calling
  `useTranslations` inside one, so it stays callable with no provider;
- the **types**, into `types.ts`, so the new modules can name them without importing the component
  back — that cycle is the one thing this phase can get wrong.

Re-export every moved symbol from the component file. `your UI package source pathindex.ts` and the test file
then name exactly what they named before, which is what keeps phase 0's harness valid.

**Phase 2 — one widget at a time, most-isolated first.** Take the dialog whose state nothing else
reads, then the next; take the coupled ones **together and last**, because splitting them
separately means designing their coordination twice. Run the harness and commit after each one — a
red run then names the extraction that did it.

**Phase 3 — the store is the residue, and sometimes the residue is nothing.** Do not design it up
front. What is left after phase 2 is exactly the state that refused to fit in a widget — and if
that turns out to be two or three values read one level down, **props are the answer and a store
would add a layer that deletes nothing**. `result-screen` is the worked example in both
directions: phase 2 left it two stateful hooks and a complexity of 32, both well inside budget, so
no store was built and the section above's step 4 never fired.

Ask of each survivor whether it is really view state. An optimistic set reconciled against a server
answer *looks* like a mutation lifecycle belonging in `your app path` — and whether it is one turns
on whether it holds a **copy of the source** (rule 1: move it) or a **log of what was asked for**,
dropped rather than merged when the source changes (this layer's, and it stays). `useResultCatalogue`
carries that argument, taken both ways, with its reversal condition.

**Phase 4 — the props interface.** Only now group into `stage` / `products` / `refine` / `share` /
`icons`, updating the call sites in the same commit. This is the only breaking phase, so it is the
one that needs the openspec `ui` change and the e2e slice.

**Phase 5 — the ledger.** You will not have to remember: `composite-complexity` fails when a
recorded number improves and its entry has not followed, and the message names the number to write.
Delete the entry when the last measure is back inside budget.

**`acme-refactor` owns the loop these phases sit inside** — MEASURE → PIN → MOVE → PROVE → RECORD,
the rule that behaviour does not change, and the same discipline applied at every other layer. This
section is the `your UI package path` half of it: the widget seam, the screen store, the prop groups.

**Do not start with the worst one.** Pick the component that is one metric over on one dimension
and do it end to end first. your screen-splitting decision's own reversal condition is "splitting made the next change
more expensive", and a cheap component is how you find that out before spending a fortnight on an
expensive one.

## Before you finish

- [ ] No `next/*`, `your contract package`, `your data layer` or `fetch` — `your guard checks` agrees
- [ ] Inside the budget: 24 props, 6 stateful hooks, 800 code lines, complexity 45 (§9). A new ledger entry is not how this is satisfied
- [ ] Every dialog this screen opens lives in its own widget folder with its own four surfaces
- [ ] No mock module (`mock-*`) imported by anything but a story or a test, and no fixture used as a prop default
- [ ] No `@base-ui/react` import; primitives come from `your component library`
- [ ] Every string the component owns comes from `your i18n package`
- [ ] A `PseudoLocale` story, and `expectNoHardcodedCopy` passing
- [ ] Stories cover the states `your interaction specification` §2 requires of this component
- [ ] Pixel baselines at all four viewports, **each one opened and looked at**
- [ ] Exported from `your UI package source pathindex.ts`
- [ ] Checked against its design node, not just against the story
- [ ] `your validation command`
