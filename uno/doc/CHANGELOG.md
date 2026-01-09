# Documentation Changelog

## [v1.1.0] - 2025-01-09

### Documentation Overhaul
Comprehensive update to all component API documentation to align with source code v2.0.

- **Standardization**: All component docs now follow a consistent structure:
    - Overview
    - Basic Usage
    - API Parameters (Table with Type, Default, Required columns)
    - Events
    - Methods
    - Examples
    - Version History
- **Language**: All documentation is now in English.

### Component-Specific Updates

#### Box (`box-api.md`)
- Added `type` parameter details (dialog, confirm).
- Documented `onOk`, `onCancel`, `onBeforeClose`, `onClose` callbacks.
- Added HTMX integration parameters (`hxUrl`, `hxMethod`).
- Clarified `width` vs `widthClass`.

#### Notify (`notify-api.md`)
- Added Global Progress Bar methods (`startProgress`, `stopProgress`).
- Documented HTMX zero-config integration events.
- Added `closable` and `duration` options.

#### Tab (`tab-api.md`)
- Documented `window.unoTabs` helper.
- Added `hash` synchronization details.
- Clarified `trigger` binding parameters (`ajaxUrl`, `hxUrl`).

#### Tooltip (`tooltip-api.md`)
- Added `trigger` modes (`hover`, `click`, `focus`).
- Documented Smart Positioning (`placement`, `offset`).
- Added HTMX dynamic content loading examples.

#### Valid (`valid-api.md`)
- Complete list of built-in validators (`email`, `pwd`, etc.).
- Documented `htmx:confirm` interception logic.
- Added `registerValidator` for custom rules.
