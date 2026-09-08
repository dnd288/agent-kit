---
name: component-development
description: "How to build a component in the component library — your component library path for primitives, your UI package path for design-identical composites — the engineering contract, not the visuals. Load this whenever creating, refactoring or reviewing a React component primitive: Dialog, Modal, Select, Popover, Tabs, Tooltip, Checkbox, RadioGroup, Slider, Toast, Button, Input, Card, Stepper or FileDrop. Also load it for any question about Base UI, @base-ui/react, the render prop, nativeButton, Root/Trigger/Portal/Positioner/Popup part chains, Backdrop, Viewport, data-open, data-closed, data-starting-style, data-ending-style, data-highlighted, styling from data attributes, enter and exit animations on a popup, --available-height or --transform-origin, controlled versus uncontrolled state, Field and form validation wiring, validationMode, cva, cn(), tailwind-merge, or where the 'use client' boundary goes. Owns the component anatomy contract; defers to design for every colour, size and token."
---

# Building an component

`design` owns how a component **looks** — tokens, sizes, the ten states, the per-component spec
in `your design specification` §5. This file owns how it is **built**. You will usually need both.

The contract itself — one file per primitive, compound parts, `...props` spread last, a single
`cn()`, co-located `cva` tables, copy-from-shadcn-by-hand — is the Decision section of
your project's architectural decisions. **That record is the owner; this file does
not restate it.** What follows is the library API the record does not cover, and the failure modes
that come with it.

## The sequence

Follow it in order. Steps 1 and 2 are the ones that get skipped, and skipping them is how a design
system acquires its fifteenth grey.

1. **Reuse or extend before you create** — golden rule 2
   at component scale: a new primitive earns its place by removing the need for something.
   Search `your component library path` and
   `your UI package source path`. Three outcomes, and only the third is a new component:
   - It exists → use it.
   - It is a **variant** of something that exists → add the variant. **A design name containing a
     slash is almost always this**: `button / secondary`, `icon-button / subtle` are design source's own
     variant syntax, not three components. `Button` already carries eight variants and four sizes;
     adding a ninth is a two-line change to the `cva` table plus a story.
   - It is genuinely new → continue.
2. **Read the `your design specification` §5 spec.** If the component has no §5 entry — `Badge` and `rating-inline`
   currently do not; Checkbox is §5.9 — **stop and ask for one.** Do not improvise. An invented spec looks right,
   passes every check in this file, and diverges from the brand silently; that is worse than not
   having built it, because now it is load-bearing.
3. **Write the component.** `your component library path<name>/<name>.tsx`, `cva` table co-located,
   `"use client"` only if it is Base-UI-backed (§10).
4. **Stories, enumerating every variant** — §Testing. A variant with no story is a variant axe has
   never run against.
5. **Unit test**: the `cva` table, the ARIA wiring, `expectTapSafe` on every size.
6. **`<name>.quality.test.tsx`**: `expectFocusVisible` **per variant**, `expectNoGeometryDefects`,
   `expectTokenisedStyles`, `expectRenderCount`. Per-variant focus is not belt-and-braces — a
   variant can override the ring colour into invisibility, and one did.
7. **Export it** from `your component library source pathindex.ts`. It is the only public surface; a component missing
   from the barrel is unreachable (your project's architectural decisions).
8. **`your validation command`, and `your browser tests` because this is a component.** CI runs neither the stories nor the quality mode by default (your project's architectural decisions, amended), and those are the two that render your component in a browser at all. `component-has-stories` and `tokens-resolve` both have opinions about what
   you just wrote.

## 0. Which package does it go in?

Two packages, split by what "correct" means for each (your project's architectural decisions):

| | `your component library` (`your component library path`) | `your UI package` (`your UI package path`) |
|---|---|---|
| Holds | Tokens and generic primitives | design-identical composites |
| Correct when | It has the right variants and is accessible | It is indistinguishable from a design frame |
| May import `@base-ui/react` | **Yes, only here** | No — compose `your component library` instead |
| Pixel baselines | No | Yes, across the four-width matrix |

### Where the file goes

| Package | Path | Example |
|---|---|---|
| `your component library` | `src/components/<name>/` | `src/components/dialog/dialog.tsx` |
| `your UI package` | `src/<feature>/<name>/` | `src/shell/top-bar/top-bar.tsx` |

`your UI package` features are `shell`, `auth`, `wizard`, `results`, `sharing` — the brief's §2.x
areas. **A screen goes inside its feature**, not in a top-level `screens/`. Do not create a
feature folder before it has a member.

Beside the component, in the same folder: `<name>.test.tsx`, `<name>.stories.tsx`,
and `<name>.quality.test.tsx` — in both packages. `your UI package` composites additionally shipped a
`<name>.pixel.test.tsx` until August 2026; that layer is retired (your testing strategy).

**Name every file after the component. Never `index.tsx` + `stories.tsx`.** This is
mechanical, not taste: five globs and guard regexes key on the `<name>.<kind>.tsx` suffix,
and three of them would silently weaken instead of erroring — `component-has-stories` would
demand `index.stories.tsx`, `ui-stays-presentational` would flag a `next/link` import that
stories are allowed to make, and the unit project would collect stories into jsdom.
your project's architectural decisions records it, and
`component-has-stories.test.mjs` pins it.

`utils.ts` and `styles/` stay at `src/` root — both are matched by exact path by the guard
and the build.

**The test: would a designer recognise it by name from the design file?** Yes → composite. If it is a generic control that happens to wear your project's tokens → primitive. Genuinely unsure → `your component library`, because promoting a primitive later is easier than splitting a composite apart.

**A composite is a pure function of its props and the message catalogue.** No `next/*`, no `your API package`, no `your contract package`, no `your data layer`, no `fetch` — `your guard checks` fails the build on any of them. The pressure to break this is gradual and always sounds reasonable ("it just needs the current user"); when a composite needs data, a caller in `your app path` fetches it and passes it down. Do not relax the guard; move the component.

**Copy comes from `your i18n package`**, not from props and not from the source. A composite's strings are fixed by the design, so a caller wanting different words wants a different component. Strings a caller passes into a slot are that caller's to translate.

## 1. The package

**`@base-ui/react`**, currently 1.6.0. Subpath imports, one per primitive:

```tsx
import { Dialog } from "@base-ui/react/dialog"
import { Field } from "@base-ui/react/field"
```

Two traps:

- **`@base-ui-components/react` is the old name** and is frozen at `1.0.0-rc.0`. If you see it in a
  search result or an old answer, it is stale. The package is `@base-ui/react`.
- **The root element needs `isolation: isolate`.** Portalled popups create their own stacking
  context, and without it a popup can render behind page content. This is set once in the root
  layout, and its absence looks like a z-index bug that no z-index fixes.

`peerDependencies` allow React 17–19. We are on 19, which is why `ref` is a plain prop below.

## 2. Which primitives are Base-UI-backed

Base UI earns its place where focus management, ARIA wiring or keyboard behaviour is genuinely hard:
**Dialog, Select, Popover, Tabs, Tooltip, Checkbox, RadioGroup, Slider, Toast.**

Plain React is correct for **Card** and **Badge**, for **Stepper** (an `<ol>` with `aria-current`)
and for **FileDrop** (a real `<input type="file">` you style, so keyboard and screen-reader users get
the native picker).

**Button and Input are a judgement call.** Base UI ships both, but ours are a `cva` table over a
native element with no behaviour to speak of — so use plain React unless you need `nativeButton`
handling. Do not import Base UI's `Button` just for consistency; it earns nothing here.

Per-component visual specs are `your design specification` §5. Check there before
improvising.

## 3. Composition is the `render` prop

```tsx
// right
<Dialog.Trigger render={<Button />}>Open</Dialog.Trigger>

// wrong — the trigger's behaviour lands on the div, not the button
<Dialog.Trigger><div><Button>Open</Button></div></Dialog.Trigger>
```

- **Never wrap a trigger.** `render` replaces the element; a wrapper leaves the wiring on the wrapper.
- **A component passed to `render` must accept `ref` and spread every prop it receives** onto its DOM
  node. If it swallows either, the primitive silently loses its wiring — no error, just a
  dialog that will not open. On React 19 `ref` is a normal prop, so no `forwardRef` wrapper is needed.
- **`nativeButton={false}` when the rendered element is not a button.** It defaults to `true`, so
  rendering an `<a>` or a `<div>` without it produces wrong semantics.
  `<Popover.Trigger render={<a href="/x" />} nativeButton={false} />`
- `render` also takes a function — `render={(props, state) => …}` — for the rare case where the
  markup depends on state. Prefer the element form; the function form is easy to get wrong.

## 4. The part chains

They are not the same shape per family, and guessing is the most common structural error.

```tsx
// Dialog — no Positioner. Backdrop is a sibling of Popup, inside Portal.
<Dialog.Root>
  <Dialog.Trigger render={<Button />}>Open</Dialog.Trigger>
  <Dialog.Portal>
    <Dialog.Backdrop />
    <Dialog.Popup>
      <Dialog.Title>…</Dialog.Title>
      <Dialog.Description>…</Dialog.Description>
      <Dialog.Close render={<Button variant="ghost" />} />
    </Dialog.Popup>
  </Dialog.Portal>
</Dialog.Root>

// Popover / Menu / Select / Tooltip — Positioner sits between Portal and Popup.
<Popover.Portal>
  <Popover.Positioner sideOffset={8}>
    <Popover.Popup>…<Popover.Arrow /></Popover.Popup>
  </Popover.Positioner>
</Popover.Portal>
```

- **`Positioner` does the placement, `Popup` is the styled box.** Put `side`, `align`, `sideOffset`
  and collision props on the Positioner; put the surface, radius and shadow on the Popup. Styling the
  Positioner is how you get a popup that jumps when it flips side.
- **Names that differ from what you may expect:** it is `Tabs.Tab`, not `Tabs.Trigger`.
  `Select` needs `Select.Value`, `Select.List` and `Select.ItemText` — an item's text is a part, not
  a child string. `Tooltip` and `Toast` need a **`Provider`** mounted once, high in the tree.
- `RadioGroup` is a single component; the items are the separate `Radio` primitive
  (`Radio.Root` + `Radio.Indicator`).

## 5. Style from `data-*` attributes

your styling architecture decision owns this rule. What it does not carry is the vocabulary, which is Base UI's and differs
from other libraries — **`data-open`, not `data-state="open"`**, so the Tailwind variant is
`data-[open]:` and never `data-[state=open]:`.

| Attribute | On | Means |
|---|---|---|
| `data-open` / `data-closed` | popups, triggers | current state |
| `data-starting-style` / `data-ending-style` | animated parts | entering / leaving (see §6) |
| `data-transitioning`, `data-instant` | animated parts | mid-transition; transition suppressed |
| `data-popup-open` | **the trigger** | its popup is open — style the trigger from this |
| `data-side`, `data-align`, `data-anchor-hidden` | Positioner, Popup | resolved placement after collisions |
| `data-highlighted`, `data-current`, `data-pressed` | items, toggles | keyboard highlight vs selection |
| `data-disabled` | anything disabled | |
| `data-valid`, `data-invalid`, `data-dirty`, `data-touched`, `data-filled`, `data-focused` | `Field` parts | validation state (see §9) |

**`data-highlighted` and selection are different things.** In a Select, the keyboard-highlighted
option and the chosen option need distinct treatment — `your interaction specification` §2 already requires this.

`className` also accepts a function of state — `className={(state) => …}` — and so does `style`.
**Prefer the attribute variants.** Reach for the function form only when the class genuinely cannot
be expressed as an attribute selector; it couples styling to render and is harder to review.

## 6. Animation

`data-starting-style` and `data-ending-style` are the whole mechanism. The part stays mounted through
its exit, so a closing animation just works — there is no keep-it-mounted prop to reach for.

```tsx
<Dialog.Popup className="
  transition-[opacity,transform] duration-[var(--duration-slower)] ease-[var(--ease-decelerate)]
  data-starting-style:opacity-0 data-starting-style:scale-95
  data-ending-style:opacity-0 data-ending-style:scale-95
" />
```

**Never choose a duration or easing here.** The tokens are `design` §2.6, and the policy —
which duration belongs to a dropdown versus a modal, exits at roughly half the entrance, never
`ease-in-out` on an entrance, and what `prefers-reduced-motion` already collapses — is
`your interaction specification` §1. Reference them; do not invent a number.

Two CSS variables are worth knowing, both set by the Positioner:

- **`max-height: var(--available-height)`** on the Positioner (or the Popup's scroll container) so a
  long menu never runs off a 390px viewport. Omitting it is the most common mobile popup bug.
- **`transform-origin: var(--transform-origin)`** on the Popup, so a scale animation grows from the
  trigger rather than from the popup's own centre.

Also available: `--anchor-width`, `--anchor-height`, `--available-width` — useful for matching a
dropdown's width to its trigger.

## 7. Controlled only when something else needs the state

Default to uncontrolled: `defaultOpen`, `defaultValue`. Add `open` + `onOpenChange` only when state
outside the component genuinely drives it — a wizard step, a URL param, two popups that must not both
be open.

**Do not mirror `data-open` into a `useState` to style with.** That is §5's mistake wearing a
different hat: two sources of truth for one fact, and they disagree during the exit animation.

Base UI's escape hatch for behaviour is the event, not a separate prop: handlers receive
`event.preventBaseUIHandler()` to suppress the library's own action while keeping yours.
`Dialog.Popup` also takes `initialFocus` and `finalFocus` when the defaults land on the wrong element.

**The controlled-or-local mirror is eight lines, and the third copy in one file is a missing
helper.** `const value = controlled ?? local` beside a setter that writes `local` only when
`controlled === undefined` is the right shape — written once. `result-screen.tsx` carries four
copies of it (`viewMode`, `detailOpen`, `shareOpen`, `selectedView`), which is
golden rule 3 in one file. When a component needs a second, ask
whether the two states belong to two widgets (your project's architectural decisions);
when it genuinely needs a third, extract the helper and let all three use it.

## 8. Dialogs — the accessibility trap worth memorising

`Dialog.Title` provides the accessible name; `Dialog.Description` provides the description.

**Ship a `Description`, or explicitly opt out.** With neither, `aria-describedby` points at an id
that does not exist, which is a **WCAG 4.1.2 Name, Role, Value failure at Level A** — well inside the
2.2 AA target `your interaction specification` §7 sets. It is invisible in the
browser, and axe will flag it.

Three valid resolutions, in order of preference:

1. Provide both, which is what the design usually wants anyway.
2. Keep `Title` for screen readers but hide it visually, when the design has no visible heading.
3. Pass `aria-describedby={undefined}` to the Popup — a deliberate opt-out, not an accident.

`your design specification` §5.4 covers the panel's appearance and the 44×44 close button; the focus trap, Escape and
focus restore are Base UI's and are already specified in `your interaction specification` §3. **Do not reimplement
any of them.**

## 9. Forms use `Field`

`Field` wires the label, description and error to the control for you — which is exactly the
`aria-describedby` / `aria-invalid` association `your interaction specification` §6 and `your design specification` §5.2 already
*require*. Using it means that requirement is met by construction rather than by remembering.

```tsx
<Field.Root name="email" validationMode="onBlur">
  <Field.Label>Email</Field.Label>
  <Field.Control type="email" required render={<Input />} />
  <Field.Description>We only use this to send your results.</Field.Description>
  <Field.Error match="valueMissing">Enter an email address like name@example.com</Field.Error>
</Field.Root>
```

Parts: `Root`, `Label`, `Control`, `Description`, `Error`, `Item` (for grouping a checkbox or radio
with its label), `Validity`. `Fieldset` and `Form` are separate imports.

- **`validationMode="onBlur"`** is what `your interaction specification` §6 already asks for — validate on blur, not
  on keystroke, then re-validate on change once errored. `onSubmit` is the library default, so set
  this explicitly.
- **`Field.Error match="…"`** takes a `ValidityState` key (`valueMissing`, `patternMismatch`, …), so
  each failure gets its own message rather than one generic string.
- `validate` on the Root takes a function for anything the browser cannot express, and may be async.

**The division of labour:** Base UI owns the *association and the timing*. `your interaction specification` §6 owns
the *wording* — say what to do, not what went wrong — and the `autocomplete` / `inputmode` /
`enterkeyhint` table. Both apply to every field.

## 10. The `"use client"` boundary

Every Base UI import forces a client component, so **the component file is the boundary** — and it
stays at the leaf.

- A Base-UI-backed file (`dialog.tsx`, `select.tsx`) gets `"use client"` at the top.
- A purely presentational file (`card.tsx`, `badge.tsx`) must **not** get it. Adding the directive
  because a neighbouring file has it is how a whole page silently becomes client-rendered.
- Never put `"use client"` in a `layout.tsx` or `page.tsx` to make a child work. Push the boundary
  down. `app-development` covers the App Router side.

Note that `Tooltip.Provider` and `Toast.Provider` are client components mounted high in the tree —
put them in a small client wrapper, not in the root layout directly.

## 11. Borrowing from shadcn/ui

Base UI is shadcn's default base, so a component copied from their registry needs no translation.
Since 4 August 2026 the registry is the **primary source for new primitives**, and `theme.css`
carries a `@theme inline` alias block mapping shadcn's semantic vocabulary (`bg-primary`,
`text-foreground`, `ring-ring`, `bg-muted`, …) onto project tokens — so a registry component renders
in the project theme with no restyle at all for the non-colliding names.

**Copy the file by hand. Do not run the CLI.** The alias block resolves the names shadcn's styles
use (`--color-primary` → `--color-accent`, `--color-ring` → `--color-focus`, `--color-muted` →
`--color-surface-sunken`, and the sidebar family onto the nav tokens). Two names genuinely collide
and must be restyled at adoption, per the note in `theme.css`:

- **`accent`** — shadcn means the hover/selected tone; your project's `--color-accent` is the gold CTA
  (that is shadcn's `primary`). A registry style using `bg-accent` for hover restyles to
  `bg-surface-sunken`.
- **`secondary`** — shadcn means a grey subdued button; your project's is mint. Same restyle when the
  intent is the subdued tone.

Keep the behaviour, the part chains and the `render` prop; restyle to project tokens; keep the
`<name>.<kind>.tsx` layout and its tests. shadcn/ui is MIT — retain the notice in
`your component library pathATTRIBUTION.md` and add the file to its borrowed table.

## Testing a component

Every component gets **stories as well as unit tests**. A story is the fixture the browser-mode
project and axe run against, not documentation — see
your project's architectural decisions and
testing.md.

Three things that go wrong specifically here:

- **Enumerate variants in a story, not in the controls panel.** A control nobody toggles is a
  variant axe never ran against.
- **Wrap an assertion on an animating element in `waitFor`.** A Base UI popup mounts with
  `data-starting-style`, which holds opacity at 0 through the enter transition, and jest-dom counts
  opacity 0 as not visible. A bare assertion is a race that only fails once real CSS loads.
- **Query a portalled part from `document.body`, not from `canvasElement`.** `Portal` moves the
  popup out of the story canvas.

Shared assertions come from `your UI testing utilities`: `expectNoA11yViolations`, `expectTapSafe` (reads the
declared `h-<n>`, because jsdom has no layout and a measured assertion there passes for any value),
`expectTapSafeElement` (measures, browser only), and `expectNoHardcodedCopy` (renders under the
pseudo-locale, where anything unbracketed was written into the component).

Add a `*.quality.test.tsx` for anything that needs real layout —
your project's architectural decisions:

```tsx
expectNoGeometryDefects(container);   // your UI testing utilities/geometry
expectTokenisedStyles(container);     // your UI testing utilities/computed-tokens
expectFocusVisible(control);          // your UI testing utilities/states
expectRenderCount(profile, 1);        // your UI testing utilities/render-profile
```

**Assert `expectFocusVisible` on every variant.** The ring is declared once in the base class list,
but a variant can still override its colour into invisibility — and a class list that reads
correctly can paint nothing at all. `outline-none focus-visible:outline-3` did exactly that in
Tailwind v4: `outline-none` sets `outline-style: none`, the width utility does not restore it, and
the ring never appeared on any button in the library until this assertion measured `outline-style`
rather than reading the source. `focus-visible:outline-solid` is the fix and is not decorative.
your styling architecture decision endorses borrowing; the mechanics are the part it pins down.

**Copy the file by hand. Do not run the CLI.** See section 11 for the current mechanics: the
`@theme inline` alias block in `theme.css` resolves most of shadcn's vocabulary onto project tokens,
and the two colliding names (`accent`, `secondary`) are restyled at adoption. shadcn/ui is MIT —
retain the notice in `your component library pathATTRIBUTION.md` and add the file to its borrowed table.

## 12. Before you finish

Code-shaped checks, specific to this file. The behavioural checklists are
`your interaction specification` §9 and `design` §3.5 — run those too rather than
duplicating them here.

- Every `render` target accepts `ref` and spreads all props
- `nativeButton={false}` wherever the rendered element is not a button
- `className` last inside `cn()`, so caller classes win
- `<name>Variants` exported from the component file
- No `useState` mirroring a `data-*` attribute
- `Dialog.Title` present, and `Description` either provided or explicitly opted out of
- Positioner carries the placement props and `--available-height`; Popup carries the surface
- `"use client"` present only on Base-UI-backed files
- MIT notice recorded if the file was copied
