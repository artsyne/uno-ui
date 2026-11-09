# Uno Component Documentation

Welcome to the documentation for the **Uno** component library. Uno is a lightweight, modular UI toolkit designed for the **Alpine.js + TailwindCSS + HTMX** stack.

## Components

| Component | Description | Documentation |
| :--- | :--- | :--- |
| **Box** | Modals, Alerts, Confirms, Prompts. | [Box API](./box-api.md) |
| **Notify** | Toast notifications, Global Progress Bar. | [Notify API](./notify-api.md) |
| **Tab** | Tabs with Hash Sync and Deep Linking. | [Tab API](./tab-api.md) |
| **Tooltip** | Smart positioning tooltips and popovers. | [Tooltip API](./tooltip-api.md) |
| **Valid** | Form validation with HTMX integration. | [Valid API](./valid-api.md) |

## General Guides

- [Configuration Guide](./config-api.md) - How to configure global themes and defaults.
- [TypeScript Usage](./typescript-usage.md) - Using Uno with TypeScript.
- [Error Handling](./error-handling.md) - Common errors and solutions.
- [FAQ](./faq.md) - Frequently Asked Questions.

## Architecture Philosophy

Uno follows a **"Store-First, Directive-Driven"** architecture:

1.  **Store-First**: State and logic reside in global Alpine stores (`Alpine.store('uno').*`). This allows programmatic control from anywhere (e.g., `Alpine.store('uno').notify.success('Saved')`).
2.  **Directive-Driven**: Directives (`x-uno-*`) provide a declarative API for HTML, handling event listeners and binding to the store automatically.
3.  **HTMX-Native**: Components are designed to work seamlessly with HTMX lifecycles (e.g., Valid intercepts `htmx:confirm`, Notify listens to `htmx:beforeRequest`).

## Directory Structure

```text
uno/
├── src/
│   ├── components/       # Component source code
│   │   ├── box/
│   │   ├── notify/
│   │   ├── tab/
│   │   ├── tooltip/
│   │   └── valid/
│   ├── core/             # Core logic (uno.js, layer.js)
│   └── utils/            # Utilities (dom, log, position)
└── doc/                  # Documentation (You are here)
```

## HTMX version

Uno targets HTMX 2.x and only relies on XHR event details. The version used by the test suite is pinned in `package.json` and `package-lock.json`; applications that vendor HTMX should upgrade it together with Uno.
