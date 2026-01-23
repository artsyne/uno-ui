---
name: frontend-design
description: Create distinctive, production-grade frontend interfaces with high design quality. Use this skill when the user asks to build web components, pages, artifacts, posters, or applications.
license: Apache-2.0
---

# Frontend Design

This skill guides creation of distinctive, production-grade frontend interfaces that avoid generic "AI slop" aesthetics. Implement real working code with exceptional attention to aesthetic details and creative choices.

## Usage

The user provides frontend requirements: a component, page, application, or interface to build. They may include context about the purpose, audience, or technical constraints.

## Design Thinking

Before coding, understand the context and align with the project's **Brand Guidelines**:

1.  **Purpose**: What problem does this interface solve? Who uses it?
2.  **Tone**: Strictly follow the **Minimalist** ethos defined in `brand-guidelines`. Structure is king; decoration is noise.
3.  **Constraints**: Technical requirements (framework, performance, accessibility).
4.  **Differentiation**: In a minimalist system, differentiation comes from **perfect spacing**, **subtle interactions**, and **flawless typography** (using the approved font stack).

**CRITICAL**: Execute the minimalist direction with precision. The key is intentionality—every pixel must justify its existence.

Then implement working code (HTML/CSS/JS) that is:

- Production-grade and functional
- Visually consistent with the Brand Guidelines
- Meticulously refined in every detail

## Frontend Aesthetics Guidelines

Focus on executing the Brand Guidelines with excellence:

- **Typography**: strict adherence to the type scale and font family defined in Brand Guidelines. Use whitespace and hierarchy rather than size to distinguish elements.
- **Color & Theme**: Use the semantic and structural colors defined in Brand Guidelines. Do not introduce new colors unless strictly necessary for a specific semantic meaning not covered by the system.
- **Motion**: Use subtle, fast transitions (duration-200 or less) as defined in Brand Patterns. Avoid bouncy or slow animations that distract from content.
- **Spatial Composition**: Use the grid and spacing scale. Generous negative space is the primary design element.
- **Backgrounds & Visual Details**: Flat, clean backgrounds. Use borders instead of shadows to define depth (Bordered Layout).

## Negative Constraints (NEVER)

- NEVER deviate from the defined Color Palette and Typography.
- NEVER use heavy box-shadows or skeuomorphic elements.
- NEVER use "cluttered" layouts; when in doubt, add more whitespace.
- NEVER use generic "AI slop" code structures (see Technical Implementation Guidelines below).

## Implementation Note

Match implementation complexity to the aesthetic vision. Maximalist designs need elaborate code with extensive animations and effects. Minimalist or refined designs need restraint, precision, and careful attention to spacing, typography, and subtle details.

## Technical Implementation Guidelines (Best Practices)

### 1. Tailwind CSS 4+
- **Utility-First**: Use utility classes for 99% of styling. Avoid custom CSS classes unless creating complex animations or utilizing features not exposed by Tailwind.
- **Dynamic Values**: Leverage Tailwind v4's arbitrary value support (e.g., `w-[123px]`) sparingly; prefer theme-consistent values.
- **Dark Mode**: Always implement `dark:` variants for every color utility.
- **Variables**: Use CSS variables for theme colors to allow runtime switching (e.g., `bg-[--bg-primary]`).

### 2. Alpine.js 3.15+
- **Logic Separation**: Avoid massive `x-data="{...}"` objects in HTML. Extract logic into `Alpine.data()` or dedicated `<script>` blocks.

  ```html
  <!-- BAD -->
  <div x-data="{ count: 0, init() { ... }, update() { ... } }">...</div>

  <!-- GOOD -->
  <div x-data="counter">...</div>
  <script>
    document.addEventListener('alpine:init', () => {
      Alpine.data('counter', () => ({
        count: 0,
        init() { ... }
      }))
    })
  </script>
  ```

- **State Management**: Use `Alpine.store()` for global state shared between components.

### 3. HTMX 2+

- **Server-Driven UI**: Let the server render HTML fragments. Do not return JSON unless specifically needed for client-side processing.
- **Clean Interactions**: Use `hx-trigger` modifiers like `delay:500ms` for search inputs to reduce server load.
- **Feedback**: Always provide visual feedback during requests using `.htmx-indicator` or `hx-indicator`.
- **Progressive Enhancement**: Ensure forms and links work reasonably well (or degrade gracefully) if JS fails, though HTMX is primary.
```

### Negative Constraints (NEVER)

- **NEVER** use heavy box-shadows or skeuomorphism.
- **NEVER** use "cluttered" layouts; whitespace is a feature.
- **NEVER** return JSON to an HTMX request unless specifically required by a client-side extension.
