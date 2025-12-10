# Phase 1 Retrospective: Infrastructure & Static Resource Build

**Date**: 2025-12-09  
**Phase**: 1 - Infrastructure

## 1. What Went Well

- **Build Speed**: `esbuild` proved to be extremely fast (~76ms), providing an instant feedback loop.
- **Modular Structure**: The separation of `core`, `stores`, and `modules` is clean. The shift to CDN for libraries made the `bundle.js` incredibly lightweight (~3.5KB), containing only our business logic.
- **Strangler Fig Implementation**: We successfully introduced the new stack (CDN-based) alongside the old one without breaking existing functionality.

## 2. Challenges Encountered

### 2.1 Bundling vs CDN

- **Issue**: Initial implementation bundled Alpine and HTMX, resulting in a 238KB bundle. The user preferred a CDN approach for better caching and separation of concerns.
- **Resolution**: Quickly pivoted to CDN links in `header.html` and updated `init.js` to rely on global `window.Alpine` and `window.htmx` objects. Removed library imports from the build process.
- **Lesson**: Clarify deployment preferences (Bundled vs CDN) early in the architectural planning.

### 2.2 Directory Path Confusion

- **Issue**: Initially created the modular JS structure in an incorrect path (`/uno/website/...`) due to a misunderstanding of the project root context.
- **Resolution**: Corrected the working directory to `/hebe` and recreated the structure in `/hebe/website/static/src/js`.
- **Lesson**: Always verify absolute paths before creating file structures.

### 2.3 Template Linter Errors

- **Issue**: `header.html` contains legacy code causing linter errors.
- **Resolution**: Ignored for now as they are out of scope.

## 3. Improvements for Future Phases

- **Strict Path Verification**: Use `ls` or `Read` to verify parent directories before writing new files.
- **CDN Versioning**: Monitor CDN versions. Currently using `htmx@2.0.4` and `alpinejs@3.14.8`.
- **Global Types**: If we introduce TypeScript later, we'll need to declare `window.Alpine` and `window.htmx` types.

## 4. Conclusion

Phase 1 is complete. The infrastructure is lean (CDN + small bundle) and ready for component migration.
