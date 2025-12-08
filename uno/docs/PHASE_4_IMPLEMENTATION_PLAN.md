# Phase 4 Implementation Plan: Structural Components

**Version**: 1.0
**Date**: 2025-12-05
**Status**: Active

## 1. Overview

Phase 4 focuses on modernizing structural components, specifically the **Tab** component. The goal is to transition from imperative DOM manipulation to a declarative, reactive architecture using `Alpine.data`.

## 2. Component Split & Architecture

The refactoring transforms the monolithic `tab.js` into a cohesive `Alpine.data` module.

### 2.1 Logical Components

1.  **Tab Container (Controller)**:
    *   Manages active state (`activeTab`).
    *   Handles Hash navigation logic.
    *   Manages AJAX content loading and caching.
    *   **Implementation**: `Alpine.data('tabs', ...)`

2.  **Tab Trigger (Interaction)**:
    *   Binds click events.
    *   Applies active/inactive classes.
    *   Updates ARIA attributes (`aria-selected`, `role="tab"`).
    *   **Implementation**: `trigger(tabId)` helper method.

3.  **Tab Panel (Content)**:
    *   Controls visibility (`x-show`).
    *   Manages ARIA attributes (`role="tabpanel"`).
    *   **Implementation**: `panel(tabId)` helper method.

## 3. Interface Design Specification

### 3.1 Initialization
```html
<div x-data="tabs('tab1', { hash: true })">
```
*   `defaultTab`: ID of the initially active tab.
*   `options`: Configuration object (e.g., enable/disable hash sync).

### 3.2 Trigger Binding
```html
<button x-bind="trigger('tab1', 'ajax-url', 'hash-fragment')">Tab 1</button>
```
*   `tabId`: Unique identifier for the tab.
*   `ajaxUrl` (Optional): URL to fetch content from.
*   `hash` (Optional): URL hash to sync with.

### 3.3 Panel Binding
```html
<div x-bind="panel('tab1')">Content</div>
```

## 4. Test Strategy

### 4.1 Automated E2E Testing (Playwright)
*   **File**: `e2e/tab.spec.ts`
*   **Scenarios**:
    *   Basic switching (click -> content change).
    *   Hash navigation (URL hash -> tab switch).
    *   Deep linking (Load page with hash -> correct tab open).
    *   AJAX loading (Fetch content -> render -> cache).
    *   Legacy compatibility (Verify `tabtarget` attributes still work via shim/init logic).

### 4.2 Manual Verification
*   **File**: `examples/components/tab.html`
*   **Checks**:
    *   Keyboard navigation (Tab key focus).
    *   Visual states (Active vs Inactive styles).
    *   Dark mode compatibility.

## 5. Quality Assurance Standards

1.  **Reactivity**: No `document.querySelector` for state changes; use `x-model` or internal Alpine state.
2.  **Accessibility**:
    *   `role="tablist"` on container (optional/manual).
    *   `role="tab"` on triggers.
    *   `role="tabpanel"` on content.
    *   `aria-selected` state management.
3.  **Performance**: AJAX content must be cached (no re-fetch on revisit).
4.  **Compatibility**: Must support existing HTML structure where possible, or provide a clear migration path.

## 6. Schedule & Progress

| Task | Description | Status |
|------|-------------|--------|
| **4.1** | Design `Alpine.data` Interface | ✅ Completed |
| **4.2** | Implement Core Logic (`tab.js`) | ✅ Completed |
| **4.3** | Update Example/Test Harness (`tab.html`) | ✅ Completed |
| **4.4** | Run E2E Tests & Fix Regressions | ✅ Completed |
| **4.5** | Final Code Review & Documentation | 🔄 In Progress |

## 7. Migration Guide

To migrate existing tabs to the new system:
1.  Add `x-data="tabs('default-id')"` to the wrapper.
2.  Replace `data-tab-target="..."` with `x-bind="trigger('...')"` (or keep legacy attributes, as the new implementation supports backward compatibility initialization).
3.  Ensure panels have matching IDs.
