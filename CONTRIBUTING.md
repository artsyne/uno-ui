# Contributing to Uno UI

Thank you for your interest in contributing to Uno UI! We welcome contributions from everyone. This document provides guidelines for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How Can I Contribute?

### Reporting Bugs

Before submitting a bug report:

- **Check existing issues** to avoid duplicates
- **Use a clear and descriptive title**
- **Describe the exact steps** to reproduce the problem
- **Provide specific examples** to demonstrate the steps
- **Describe the behavior** you observed and what you expected
- **Include screenshots** if applicable
- **Specify your environment** (browser, OS, Alpine.js version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the proposed feature
- **Explain why this enhancement would be useful**
- **List any similar features** in other libraries
- **Include mockups or examples** if possible

### Pull Requests

#### Development Workflow

1. **Fork the repository** and create your branch from `main`:

   ```bash
   git checkout -b feature/my-feature
   # or
   git checkout -b fix/issue-123
   ```

2. **Install dependencies**:

   ```bash
   npm install
   ```

3. **Make your changes**:
   - Follow the [coding style](#coding-style)
   - Add tests for your changes
   - Update documentation as needed

4. **Run tests**:

   ```bash
   npm run test:run     # Run all unit tests once
   npm run lint         # Linting
   ```

5. **Commit your changes**:

   ```bash
   git commit -m "feat: add new feature"
   # or
   git commit -m "fix: resolve issue #123"
   ```

   Use [Conventional Commits](https://www.conventionalcommits.org/) format:
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation changes
   - `test:` for adding or updating tests
   - `refactor:` for code refactoring
   - `perf:` for performance improvements
   - `chore:` for maintenance tasks

6. **Push to your fork** and submit a pull request

#### Pull Request Guidelines

- **One feature per PR**: Keep pull requests focused on a single feature or fix
- **Update documentation**: Include relevant documentation updates
- **Add tests**: Ensure your changes are covered by tests
- **Pass all tests**: All tests must pass before merging
- **Follow code style**: Use ESLint and Prettier
- **Update CHANGELOG**: Add an entry to `CHANGELOG.md` for user-facing changes
- **Link issues**: Reference related issues in the PR description

## Coding Style

### JavaScript/TypeScript

- **ES6+ syntax**: Use modern JavaScript features
- **ESLint**: Follow the project's ESLint configuration
- **Prettier**: Format code with Prettier
- **JSDoc comments**: Document public APIs with JSDoc

Example:

```javascript
/**
 * Show a notification message
 *
 * @param {string} message - The notification message
 * @param {Object} [options] - Configuration options
 * @param {string} [options.type='info'] - Notification type (info, success, error, warning)
 * @param {number} [options.duration=3000] - Display duration in milliseconds
 * @returns {Object} Notification item instance
 *
 * @example
 * Alpine.store('uno').notify.show('Hello World', { type: 'success' })
 */
show(message, options = {}) {
  // Implementation
}
```

### Naming Conventions

- **camelCase** for variables and functions
- **PascalCase** for classes and components
- **UPPER_CASE** for constants
- **kebab-case** for file names

### Architecture Guidelines

Uno UI follows a **Store-First, Directive-Driven** architecture:

1. **Store-First**: State management via Alpine.js stores

   ```javascript
   Alpine.store('uno').componentName.method()
   ```

2. **Directive-Driven**: Declarative HTML interface

   ```html
   <button x-uno-componentName="{ option: value }">Click</button>
   ```

3. **4-Level Configuration**: Support multiple configuration levels
   - Plugin Defaults
   - Global User Config (via `plugins`)
   - Direct Global Config
   - Instance Config (highest priority)

4. **Plugin Registry**: Use unified plugin registration
   ```javascript
   import { registerPlugin } from '../../core/plugin-registry'
   registerPlugin(Alpine, 'componentName', componentStore)
   ```

## Testing

### Unit Tests

- Use **Vitest** for unit tests
- Test files: `tests/**/*.test.js`
- Coverage goal: 100% for core functionality
- All components and utilities must have comprehensive unit tests

Example:

```javascript
import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('Component Store', () => {
  let store
  let Alpine

  beforeEach(() => {
    // Setup mocks
    Alpine = {
      store: vi.fn(),
      nextTick: vi.fn((fn) => setTimeout(fn, 0))
    }

    // Create store instance
    store = createComponentStore(Alpine)
  })

  it('should initialize with default state', () => {
    expect(store.items).toEqual([])
    expect(store.visible).toBe(false)
  })

  it('should update state correctly', () => {
    store.show({ content: 'Test' })
    expect(store.visible).toBe(true)
  })
})
```

**Test Coverage Requirements:**

- ✅ All public methods must be tested
- ✅ Edge cases and error handling must be covered
- ✅ Configuration merge logic must be verified
- ✅ Callbacks and event handlers must be tested
- ✅ Integration with Uno.getConfig() must be validated

## Documentation

### API Documentation

- Document all public APIs with JSDoc
- Include parameter types and descriptions
- Provide usage examples
- Update relevant `.md` files in `doc/`

### README Updates

Update `README.md` if your changes:

- Add new features
- Change installation steps
- Modify usage examples
- Update supported environments

### TypeScript Definitions

- Update `types/index.d.ts` for API changes
- Ensure type definitions match implementation
- Test TypeScript usage in `doc/typescript-usage.md`

## Component Development

### Creating a New Component

1. **Create component directory**:

   ```
   src/components/mycomponent/
   ├── index.js           # Plugin factory
   ├── mycomponent.store.js   # Store logic
   ├── mycomponent.directive.js  # Directive registration
   └── mycomponent.template.js   # HTML templates
   ```

2. **Implement Store-First pattern**:

   ```javascript
   // mycomponent.store.js
   export const createMyComponentStore = (Alpine) => ({
     // State
     items: [],

     // Methods
     doSomething(options) {
       const config = Uno.getConfig('mycomponent', options)
       // Implementation
     }
   })
   ```

3. **Register plugin**:

   ```javascript
   // index.js
   import { Uno } from '../../core/uno'
   import { registerPlugin } from '../../core/plugin-registry'
   import { createMyComponentStore } from './mycomponent.store'
   import { registerMyComponentDirectives } from './mycomponent.directive'

   export default (Alpine, options = {}) => {
     // Register plugin defaults
     Uno.registerPluginDefaults('mycomponent', {
       /* defaults */
     })

     // Register store
     const store = createMyComponentStore(Alpine)
     registerPlugin(Alpine, 'mycomponent', store)

     // Register directives
     registerMyComponentDirectives(Alpine)
   }
   ```

4. **Add tests**:
   - Unit tests: `tests/components/mycomponent.test.js`
   - Test all public methods, edge cases, and configuration merge
   - Ensure 100% coverage of core functionality

5. **Create documentation**:
   - API docs: `doc/mycomponent-api.md`
   - Update `README.md`
   - Add TypeScript types to `types/index.d.ts`

6. **Add examples**:
   - Create `examples/components/mycomponent.html`

## Performance Guidelines

- **Minimize DOM operations**: Batch updates, use RAF when needed
- **Avoid memory leaks**: Clean up event listeners, use WeakMap/WeakSet
- **Lazy initialization**: Initialize components only when needed
- **Bundle size**: Keep the library lightweight (<25KB gzip)

## Security Guidelines

- **XSS Prevention**: Sanitize all user-provided HTML
- **CSP Support**: Support Content Security Policy with nonce
- **Dependency Security**: Minimize dependencies, use trusted sources
- **Input Validation**: Validate all user inputs

## Release Process

1. **Update version** in `package.json`
2. **Update CHANGELOG.md** with all changes
3. **Run all tests**: `npm run test:run`
4. **Build**: `npm run build`
5. **Create release commit**: `git commit -m "chore: release v0.x.0"`
6. **Tag release**: `git tag v0.x.0`
7. **Push**: `git push && git push --tags`

## Getting Help

- **GitHub Discussions**: Ask questions and discuss ideas
- **GitHub Issues**: Report bugs and request features
- **Documentation**: Read the docs in `doc/` and `README.md`

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Uno UI! 🎉
