# Changelog

All notable changes to this project will be documented in this file.

## [0.7.0] - 2026-01-20

### 🚀 Major Improvements

- **Lazy Initialization**: The core plugins (`box`, `notify`, `tab`, `tooltip`) have been refactored to support lazy initialization. This means:
  - Global DOM elements (containers, overlays) are no longer injected immediately upon application startup.
  - Resources are allocated only when a component is actually used (e.g., calling `box.open()`, `notify.show()`, or using directives like `x-uno-tab`).
  - Improved initial page load performance and reduced DOM footprint.

### ⚠️ Breaking Changes

- **Internal API**: The `init()` methods on store instances (e.g., `Alpine.store('uno').box.init()`) no longer trigger immediate DOM injection. They are now effectively no-ops or handle internal state setup only. If you were manually calling these methods to force DOM injection, this will no longer work (nor is it necessary).
- **Tab Component**: The Tab component now only starts listening for `hashchange` events after the first tab is registered (via directive or JS). This prevents Uno from interfering with hash routing on pages that don't use tabs.

### 🐛 Fixes

- Fixed an issue where `ensureContainer` was being called eagerly, causing unnecessary DOM mutations.

## [0.6.0] - 2025-12-18

### ✨ Features

- **TypeScript Support**: Full TypeScript definitions included.
- **Valid Component**: New validation engine.
- **I18n**: Enhanced internationalization support.

### 📦 Chore

- Updated dependencies.
