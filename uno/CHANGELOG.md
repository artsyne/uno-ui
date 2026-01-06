# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2026-01-06

### Added

- **TypeScript Support**: Full TypeScript type definitions included (`types/index.d.ts`, `types/alpinejs.d.ts`)
  - Complete type definitions for all components (Box, Notify, Tooltip, Tab, Valid)
  - Type-safe configuration interfaces
  - IntelliSense support in IDEs
  - Compile-time type checking
  - Alpine.js store type extensions
  - Alpine.js minimal type definitions included (no external dependencies)
  - No additional setup required
- TypeScript Usage Guide documentation (`doc/typescript-usage.md`)
- tsconfig.example.json for reference configuration
- Updated README with TypeScript examples and features

### Changed

- Version bumped to 0.6.0
- package.json updated with type declarations
- Enhanced README with TypeScript support section

### Technical Details

**Type Coverage:**

- Core: `UnoInstance`, `UnoOptions`, `DomManager`, `HtmxUtilities`
- Box: `BoxOptions`, `BoxItem`, `BoxStore`, `BoxPosition`, `BoxType`
- Notify: `NotifyOptions`, `NotifyItem`, `NotifyStore`, `NotifyPosition`, `NotifyType`
- Tooltip: `TooltipOptions`, `TooltipStore`, `TooltipPlacement`, `TooltipTrigger`
- Tab: `TabOptions`, `TabStore`
- Valid: `ValidOptions`, `ValidStore`, `ValidationRule`, `ValidationResult`
- I18n: `I18nStore`, `I18nMessages`
- UI Configuration: Complete type definitions for all styling options

**Developer Experience:**

- Zero-cost abstraction: Type definitions don't affect runtime bundle size
- Full autocomplete and IntelliSense support
- Type-safe configuration at all levels
- Compatible with both JavaScript and TypeScript projects

## [0.5.0] - 2026-01-05

### Added

- Store-First, Directive-Driven architecture
- 9-point positioning system for Box component
- Performance benchmarks (~0.19ms per operation)
- Security: XSS protection with automatic HTML sanitization
- HTMX native integration
- Comprehensive E2E test suite (195 tests, 100% passing)
- Modular component structure
- Unified plugin registration system

### Removed

- Backward compatibility: `trustedHTML` option removed for security
- Legacy API: Direct store mutations replaced with plugin registry

### Changed

- All components refactored to Store-First pattern
- Components registered under `Alpine.store('uno')` namespace
- Configuration system upgraded to 4-level merge strategy
- Documentation updated and expanded

### Breaking Changes

- All components now accessed via `Alpine.store('uno').{component}`
- `trustedHTML` configuration removed (use HTMX for dynamic content)
- Legacy direct store access deprecated

### Security

- Removed `trustedHTML` to prevent XSS vulnerabilities
- All user content automatically sanitized
- Image tags removed from sanitization to prevent `onerror` XSS

### Performance

- Tooltip operations: ~0.19ms per show/hide
- Bundle size: 19.88 kB (gzipped)
- Zero layout thrashing with RAF optimization

## [0.4.0] - 2025-12-XX

### Added

- Initial public release
- Basic component library (Box, Notify, Tooltip, Tab, Valid)
- Alpine.js integration
- Tailwind CSS styling

---

## Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backward compatible manner
- **PATCH** version for backward compatible bug fixes

## Links

- [Documentation](./README.md)
- [TypeScript Guide](./doc/typescript-usage.md)
- [API Reference](./doc/)
- [GitHub Repository](#)
