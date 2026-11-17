# Accessibility

This page describes the accessibility behaviour that is implemented in the current source and
covered by the automated tests. It is not a WCAG conformance statement. Component scans and
keyboard tests are evidence for specific behaviours, not proof that an application is
accessible.

## What the tests cover

The repository runs component-level axe scans on hydrated Alpine DOM for all five components.
The shared helper disables rules that need a real page or layout engine:

| Disabled rule | Reason |
| :--- | :--- |
| `color-contrast` | Requires real layout and computed colors. |
| `region`, `landmark-one-main`, `page-has-heading-one`, `bypass`, `html-has-lang`, `document-title` | Host-page responsibilities, not component responsibilities. |

Two token-level contrast suites assert the intended color tokens (the floating Box close
button and the Valid character counter), but they do not compute contrast ratios. Keyboard
contracts are covered by component tests for Box, Tab, Tooltip and Valid. See
[Test reports](test-reports.md) for commands and scope.

## Box

| Behaviour | Implementation |
| :--- | :--- |
| Dialog semantics | `role="dialog"`, `:aria-modal`, `:aria-describedby` pointing at the content container. |
| Accessible name | Resolves `labelledBy` > the title `h3` (`aria-labelledby`) > a localized `aria-label` fallback (`box.dialog`); the unused attribute is bound to `null`. `prompt` uses a visible `<label for>` and a localized title (`box.prompt`). The close button gets `aria-label` from `x-t.aria-label="box.close"`. |
| Initial focus | First visible input/textarea/select when `autoFocus` is `true`; otherwise the title heading, or the container. Focus never lands on the floating close button. |
| Focus trap | `x-uno-trap` cycles Tab/Shift+Tab, including DOM-order-aware wrapping from the neutral target. |
| Focus restore | Closing restores focus to the opener when focus is still inside the dialog or has fallen to `body`. |
| Stacking | Underlay dialogs receive `inert` and `aria-hidden="true"`. A topmost modal also isolates the page background (the `<body>` children outside `#uno-box-manager`) with native `inert`; close restores only the nodes the component inented, and `closeAll()` restores the background too. |
| Dismissal | `closeOnEsc` and `closeOnOverlayClick` are opt-in (`false` by default); Escape is consumed only by the topmost open box, and an open tooltip consumes it first. The close button is always present unless `showCloseButton: false`. |
| Buttons | OK, Cancel and extra buttons are real `<button type="button">` elements with visible text. |

`aria-describedby` points at the content container even when the dialog has no body text; keep
that in mind for very short dialogs.

## Notify

| Behaviour | Implementation |
| :--- | :--- |
| Live region | `role="alert"` + `aria-live="assertive"` for `error`/`warning`; `role="status"` + `aria-live="polite"` for `success`/`info`. |
| Announcement | The live region sits on the message text and the HTMX content wrapper, not the whole item, so the message is announced when it appears but action buttons are not re-announced. |
| Close button | Icon-only button with a visually hidden label from `box.close`. |
| Focus | Notifications never steal focus; closing a focused notification hands focus to another visible notification. |
| Persistence | `error()` is persistent by default (`duration: 0`); auto-dismiss pauses on hover/focus and resumes with the remaining time. |
| Progress bar | Visual only; it has no `progressbar` role in the current template. |

## Tooltip

| Behaviour | Implementation |
| :--- | :--- |
| Popup role | The container has `role="tooltip"`. |
| Hover/focus | Hints: no focus movement, no `aria-expanded`, no trap. Hover also shows on keyboard focus. All triggers associate the trigger's `aria-describedby` with the tooltip id as a token, removed on hide. |
| Click popup | Trigger receives `aria-haspopup="true"`, `aria-expanded="true"`, and `aria-describedby="uno-tooltip-content"`; focus moves to the first focusable element in the popup when one exists. |
| Keyboard | Tab cycles inside the popup. Arrow keys/Home/End move between focusable items **only when the content declares `role="menu"`**. Escape and outside click close it; focus returns to the trigger when it was inside the popup. |
| Menu markup | The author must declare `role="menu"` and `role="menuitem"`; the directive does not add them. |

Known limitation: the popup keeps `role="tooltip"` in click mode, even when it behaves like a
menu or disclosure. This is the current implementation, not a full APG menu pattern. If you
need a full menu, use a component designed for that pattern.

## Tab

| Behaviour | Implementation |
| :--- | :--- |
| Tab semantics | The directive sets `role="tab"`, `id="tab-{name}"` when absent, `aria-controls`, `aria-selected`, and roving `tabindex` (active `0`, others `-1`). |
| Panel relationship | Each panel declares `id="panel-{name}" role="tabpanel" aria-labelledby="tab-{name}"`. Shared-panel mode uses `role="tabpanel" data-tab-group="{group}"` and updates `aria-labelledby` on activation. |
| Keyboard | ArrowLeft/ArrowRight rotate within the group; Home/End jump to the first/last tab. Focus follows and activation reuses the click path. RTL containers reverse the horizontal arrows. Enter/Space activate the focused tab (Space synthesized on anchors, scroll prevented, repeat-guarded). |
| Tab list label | The author supplies an accessible label on the `role="tablist"` element. |

Standalone panels are your markup: the store does not render or manage `x-show` for you.

## Valid

| Behaviour | Implementation |
| :--- | :--- |
| Error element | Generated errors get `role="alert"` and the stable id `${field}-error`; author-provided id/role are preserved. |
| Error association | The input's `aria-describedby` receives the error id, appended to existing tokens, and `aria-invalid` is set or removed at every error mutation. |
| Counter | `showCounter` links the counter through `aria-describedby` and gives it a full-sentence `aria-label` (i18n `valid.counter`), while the visible text stays compact. |
| Error summary | Two or more failing fields produce a focusable `role="alert"` summary at the top of the form. Each link carries the field name (label / `aria-label`, falling back to the field name) plus the message and anchors the real control, minting a stable id when the control has none. Inline errors remain. |
| Scroll/focus | `validateForm()` focuses then scrolls the summary or the first failing field unless `scrollOnError: false`; scrolling honors `prefers-reduced-motion`. |
| Server errors | A `422` JSON body passed to `showErrors()` updates the same error state and `aria-invalid`. |

## Manual acceptance

Automated checks cannot replace a real browser and assistive technology. Before a release,
check at least:

1. Keyboard-only: open/close every component, reach every action, confirm visible focus.
2. Screen reader: one desktop (NVDA/JAWS) and one mobile (VoiceOver/TalkBack) pass over the
   five component demos.
3. Zoom: 200% browser zoom and a 320 CSS px viewport for dialogs, toasts and popovers.
4. Reduced motion: decide whether your application needs a motion preference override; Uno's CSS
   transitions do not currently read `prefers-reduced-motion`, though the Valid error-summary
   scroll does honor it.
5. Error flows: client-side and `422` server errors, including the summary and focus order.

## Reporting accessibility issues

Open an issue with the browser, assistive technology, component, reproduction and expected
behaviour. Security-sensitive reports should follow [SECURITY.md](https://github.com/artsyne/uno-ui/blob/main/SECURITY.md).

## Related

- [Test reports](test-reports.md)
- [Box API](box-api.md#accessibility)
- [Notify API](notify-api.md#accessibility)
- [Tooltip API](tooltip-api.md#accessibility)
- [Tab API](tab-api.md#accessibility)
- [Valid API](valid-api.md#accessibility)
