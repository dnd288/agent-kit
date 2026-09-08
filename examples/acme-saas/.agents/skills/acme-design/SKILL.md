---
name: design
description: "Design methodology for UI work: make design decisions from system intent, translate specs into accessible code, preserve token consistency, handle responsive/state/motion behavior, and review visual changes without hard-coding one project's brand details."
---

# Design methodology

Use this before building or restyling UI: screens, components, forms, layout, color, typography, motion, responsive behavior, or design-token work.

The goal is not to make a screen look plausible. The goal is to implement the design system in a way that stays consistent, accessible, and easy to change.

## Intent triage

Start by identifying the kind of design work.

| Work | Start with |
|---|---|
| New primitive | Existing component contract, token names, states, accessibility requirements |
| New composite/widget | Feature intent, layout rhythm, data/props boundary, responsive states |
| Screen styling | Page hierarchy, spacing system, breakpoint matrix, empty/loading/error states |
| Color choice | Semantic token purpose, contrast, fill-vs-text role, interaction states |
| Form work | Label/help/error model, focus states, validation timing, target size |
| Responsive fix | Content priority, breakpoint behavior, touch targets, overflow rules |
| Motion | What changed, whether motion clarifies it, reduced-motion behavior |
| Design-system change | Provenance, affected components, migration path, token audit |
| Visual bug | Source of truth, deliberate divergences, regression evidence |

Do not read or change everything. Follow the narrowest path that can answer the design question.

## Design authority

Every visual decision needs an authority:

1. The design system specification.
2. The component's documented contract.
3. The Figma/design source or approved mockup.
4. Accessibility standards and platform conventions.
5. Existing product patterns.

If authorities disagree, do not average them. State the conflict, choose the one the project says wins, and record any hard-to-reverse choice in the proper decision artifact.

When no design exists, compose from existing primitives and tokens. Avoid inventing new visual language for one surface.

## Tokens over values

Components use semantic tokens, not raw values.

- Use tokens by purpose: `surface`, `text`, `border`, `accent`, `focus`, `error`, `spacing`, `radius`, `shadow`, `motion`.
- Do not reference palette/ramp steps in component code unless the design system explicitly exposes them for that layer.
- Do not add a new token until an existing token fails by purpose, not by taste.
- A token change is a system change. Audit every component or surface that uses it.
- Generated token files stay generated. Change the source/extraction process, then regenerate and propagate.

Token names should describe roles, not colors. `text-danger` survives a palette change; `red-600` does not.

## Color, type, and spacing

Color role matters as much as color value.

- A color safe as a fill may fail as text on a light surface.
- A color safe as text may not work as a border or focus indicator.
- Decorative borders cannot be the only boundary if contrast is low.
- Error, warning, success, selected, disabled, and focus states need non-color cues where meaning matters.
- Check contrast in the actual pairing: foreground, background, font size, weight, and state.

Use the type scale as hierarchy, not a menu. Body copy and UI copy often use different sizes; keep that split consistent. Headings express document structure, not visual size alone. Decorative typefaces are accents only.

Layout should reveal hierarchy and reduce future branching. Prefer established containers, gutters, and spacing steps. Wide content must scroll within its own container rather than forcing body overflow. Empty, loading, error, and permission states are part of the layout.

## Responsive method

Test responsive behavior as a matrix, not as mobile plus desktop. Cover narrow mobile, tablet portrait, tablet landscape or small desktop, and desktop.

At each viewport, verify:

- Primary action remains reachable.
- Navigation has one clear mode.
- Tables/lists preserve the highest-priority content.
- Forms maintain label, field, help, and error relationships.
- Touch targets meet the project minimum.
- No horizontal page scroll appears.

If a layout needs several breakpoint-only exceptions, check whether the content priority or component split is wrong.

## Component states

A component is incomplete until its states are designed and coded.

Minimum state matrix: default, hover where available, focus-visible, active/pressed, disabled, loading/busy, error/invalid, and selected/current.

State changes should use tokens. Do not hide missing state design behind opacity hacks.

## Forms

Forms must be understandable before, during, and after validation.

- Every input has a programmatic label.
- Helper text and errors are associated with the field.
- Required/optional meaning is explicit and consistent.
- Validation timing should match user intent; avoid punishing typing.
- Error text says how to recover.
- Focus moves only when it helps the user recover or continue.
- Buttons communicate pending state and prevent duplicate submission when needed.

Do not rely on placeholder text as the only label.

## Motion

Motion earns its place by explaining change: overlay entry/exit, newly available controls, cause/effect after user action, and small affordance changes.

Avoid decorative loops, long transitions that block interaction, layout shifts that make content hard to track, and anything without a reduced-motion fallback. Use duration and easing tokens.

## Translating specs to code

Translate design into code in this order:

1. Identify the component or screen boundary.
2. Name the design authority and any unresolved gaps.
3. Map visual values to semantic tokens.
4. Define the state matrix.
5. Define responsive behavior and content priority.
6. Implement with existing primitives first.
7. Add tests/stories/examples for states that can regress.
8. Verify against the design source and accessibility requirements.

If the design needs data, the component API should describe that data. Do not let visual components fetch, infer, or own data unless the project architecture explicitly allows it.

## Maintaining the design system

Design systems decay through small exceptions.

Red flags:

- Raw hex, pixel, shadow, or timing values in component code.
- A new token used in only one place with no stated purpose.
- A variant that exists for one screen and deletes no duplication.
- Different components solving the same state differently.
- A screenshot match that breaks contrast, focus, keyboard access, or responsive behavior.
- A design-source mismatch called a bug without checking deliberate divergences.
- Generated token output edited by hand.

A new abstraction earns its place by replacing repeated patterns. If it sits beside them, it made the system larger without making change cheaper.

## Review checklist

- [ ] The design authority is named.
- [ ] Semantic tokens are used instead of raw values.
- [ ] Color pairings meet contrast requirements in each state.
- [ ] Typography follows the system scale and hierarchy.
- [ ] Layout uses established spacing, containers, radius, and elevation.
- [ ] Responsive behavior is checked across the viewport matrix.
- [ ] Interactive states include focus-visible and disabled behavior.
- [ ] Forms have labels, help/error associations, and recoverable errors.
- [ ] Motion uses system tokens and honors reduced motion.
- [ ] Empty, loading, error, and permission states are covered.
- [ ] Hard-to-reverse design choices are recorded.
- [ ] No generated design artifact was edited by hand.
