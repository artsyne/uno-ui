# Uno UI Refactoring Master Plan

**Version**: 1.0  
**Status**: Draft  
**Goal**: Transform Uno UI into a high-quality, "Alpine-native" open-source library.

This plan is designed to be executed incrementally. Each task is self-contained and includes verification steps.

## 🛠 Phase 1: Foundation & Core Infrastructure

_Objective: Establish a solid base for refactoring, ensuring coding standards and core utilities are Alpine-friendly._

- [ ] **Step 1.1: Setup Linting & Formatting**
  - **Action**: Configure ESLint and Prettier for the project. Enforce consistent code style (StandardJS or Airbnb).
  - **Target**: Root directory, `package.json`.
  - **Verification**: Run `npm run lint` without errors.

- [ ] **Step 1.2: Refactor Core Utilities (i18n)**
  - **Action**: Convert `src/core/i18n.js` to a pure `Alpine.store('i18n')`.
  - **Changes**: Remove the standalone module logic if possible, or wrap it tightly so it's reactive.
  - **Verification**: Verify `t()` function works inside an Alpine component context.

- [ ] **Step 1.3: Standardize Directory Structure**
  - **Action**: Ensure all components follow `src/components/[name]/[name].js` (already mostly done).
  - **Refinement**: Create a shared `src/utils` for common helpers (e.g., DOM utils that Alpine doesn't provide).

## 🧩 Phase 2: Directive-based Components (Low Complexity)

_Objective: Convert simple DOM-manipulating components into elegant Alpine Directives._

- [x] **Step 2.1: Refactor Tooltip**
  - **Input**: `src/components/tooltip/tooltip.js`
  - **Action**: Convert from `Alpine.data` (or mixed logic) to `Alpine.directive('tooltip')`.
  - **New API**: `<button x-tooltip.top="'Hello'">` or `<button x-tooltip="{ content: 'Hello', placement: 'top' }">`.
  - **Implementation**: Use a lightweight positioning engine (like Floating UI) if `positioning.js` is insufficient, or refine `src/core/positioning.js`.
  - **Verification**: Hovering elements shows tooltip; supports dynamic content updates.

- [x] **Step 2.2: Refactor Validation (Valid)**
  - **Input**: `src/components/valid/valid.js`
  - **Action**: Remove manual `classList` manipulation. Create a directive `x-validate` or logic in `Alpine.data` that binds to `x-model`.
  - **Implementation**: Use `x-bind:class="{ 'border-red-500': error }"` pattern.
  - **Verification**: Form inputs visually update on invalid state without imperative JS calls.

## 📦 Phase 3: Store-driven Components (High Complexity)

_Objective: Move global UI elements (Modals, Toasts) to a state-driven architecture using `Alpine.store` and `x-teleport`._

- [x] **Step 3.1: Refactor Notify (Toast)**
  - **Input**: `src/components/notify/notify.js`
  - **Action**:
    1. Create `Alpine.store('notify')` with `show()`, `close()`.
    2. Create a single global rendering component that uses `x-for` to render toasts from the store.
    3. Remove `document.createElement` for individual toasts.
  - **Verification**: Calling `Alpine.store('notify').success('Hi')` adds a toast to the UI.

- [x] **Step 3.2: Refactor Box (Modal) - Part 1: State Management**
  - **Input**: `src/components/box/box.js`
  - **Action**: Create `Alpine.store('box')` to manage the stack of open modals.
  - **Goal**: Support multiple stacked modals (z-index management).

- [x] **Step 3.3: Refactor Box (Modal) - Part 2: Rendering**
  - **Action**: Implement the Modal Container using injected template and `x-for` iterating over the store.
  - **Change**: Instead of creating DOM nodes in JS, the JS only pushes data to the store. The Alpine view handles the rendering.
  - **HTMX**: Ensure the modal content can be loaded via HTMX.

## 📑 Phase 4: Structural Components

_Objective: Modernize layout components._

- [ ] **Step 4.1: Refactor Tabs**
  - **Input**: `src/components/tab/tab.js`
  - **Action**: Convert to `Alpine.data('tabs', { default: '...' })`.
  - **Logic**: Use `x-model` or internal state to track active tab. Use `x-show` for panels and `:class` for tab buttons.
  - **Verification**: Clicking tabs switches content; supports keyboard navigation (ARIA).

## 🚀 Phase 5: Ecosystem Integration & Polish

_Objective: Ensure seamless compatibility with HTMX and Tailwind._

- [ ] **Step 5.1: HTMX Events**
  - **Action**: Add listeners for `htmx:afterSwap` to re-initialize or update Alpine components if necessary (usually Alpine handles this, but custom logic might need checks).
  - **Feature**: Add `hx-` attributes support to Box/Modal for dynamic loading.

- [ ] **Step 5.2: Tailwind Config Extraction**
  - **Action**: Move hardcoded class strings (e.g., `bg-blue-500`) into a configuration object or CSS variables strategy to allow theming.
  - **Goal**: Users should be able to override primary colors without editing library code.

- [ ] **Step 5.3: Documentation & Examples**
  - **Action**: Update `examples/` folder to use the new API.
  - **Doc**: Generate a new `README.md` with "Getting Started".

## Execution Guidelines for AI

1.  **Read Context First**: Before refactoring a component, read its existing `.js` file and its corresponding test/example.
2.  **Test-First**: Check `e2e/` tests. Run them. If they fail after refactor, update the test or the code.
3.  **Minimal Changes**: Do not rewrite business logic (validation rules, etc.) unless necessary for the architecture shift.
4.  **Verify**: After each step, confirm the component still renders and functions in the browser (or via test).

---

_Plan generated by Trae AI_
