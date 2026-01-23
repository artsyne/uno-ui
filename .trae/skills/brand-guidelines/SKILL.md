---
name: brand-guidelines
description: Applies the project's official minimalistic design system, focusing on Tailwind CSS 4 utility classes, grayscale aesthetics, and dynamic theming.
license: Apache-2.0
---

# Brand Guidelines & Design System

This skill provides the official design language for the **iDiary.ai** project. It enforces a **minimalist, content-focused** aesthetic implemented strictly via **Tailwind CSS 4 Utility Classes**.

## Core Philosophy

- **Minimalism**: Content is king. UI should recede.
- **Grayscale Foundation**: Structure is built on `gray-50` to `gray-900`. Color is used sparingly for semantic meaning or user-selected themes.
- **Utility-First**: No custom CSS classes unless absolutely necessary. Use Tailwind v4 utilities.
- **Dark Mode First**: All designs must fully support `dark:` variants.

## Typography

**Font Family**:

- Primary: `'Inter', 'Noto Sans SC', system-ui, sans-serif`
- Use `font-sans` utility.

**Scale**:

- Headings: `text-3xl font-bold` (H1), `text-2xl font-semibold` (H2).
- Body: `text-base` or `text-sm`.
- Relaxed Reading: Use `leading-relaxed` or `tracking-wide` for long-form content.

## Color Palette

**Structural Colors (Neutral)**:

- **Background**: `bg-white` (Light) / `dark:bg-gray-900` or `dark:bg-black` (Dark).
- **Surface (Cards/Forms)**: `bg-gray-50` (Light) / `dark:bg-gray-800` (Dark).
- **Borders**: `border-gray-200` (Light) / `dark:border-gray-700` (Dark).
- **Text**:
  - Primary: `text-gray-900` (Light) / `dark:text-gray-100` (Dark).
  - Secondary: `text-gray-600` (Light) / `dark:text-gray-400` (Dark).
  - Muted: `text-gray-500` (Light) / `dark:text-gray-500` (Dark).

**Semantic Colors**:

- **Action/Link**: `text-blue-600 dark:text-blue-400`.
- **Danger/Error**: `text-red-600 dark:text-red-400` or `bg-red-600`.
- **Success**: `text-green-600 dark:text-green-400`.

**Dynamic Theming**:

- The project uses classes `.dt1` through `.dt16` for user-customizable accent colors.
- When designing generic components, stick to **Neutrals**. Do not hardcode "Purple" or "Orange" unless it's a specific semantic state.

## Component Patterns

### 1. Cards (Container)

Minimalist bordered container.

```html
<div
  class="relative border border-gray-200 dark:border-gray-700 p-4 dark:bg-gray-800 break-words transition-colors duration-200"
>
  class="relative border border-gray-200 dark:border-gray-700 p-4
  dark:bg-gray-800 transition-colors"
  <h3 class="text-gray-900 dark:text-gray-100">Card Title</h3>
  <div class="border-b border-gray-300 dark:border-gray-600 mb-2 pb-2">
    <div class="text-gray-700 dark:text-gray-300">Content goes here...</div>
  </div>
  <p class="text-gray-700 dark:text-gray-300">Content goes here...</p>
</div>
### 2. Buttons Flat, bordered, subtle interaction. **Default (Secondary)**: ###
Primary Action Button
<button
  class="border border-gray-300 dark:border-gray-600 px-4 py-2 text-sm tracking-widest text-black dark:text-white bg-gray-50 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 transition cursor-pointer"
>
  CANCEL
</button>
```

**Primary (Action)**:

- Typically relies on the specific context, but often `btn-red` is used for primary actions in this project.

````html
<button
  class="border border-transparent bg-red-600 dark:bg-red-500 px-4 py-2 text-sm tracking-widest text-white hover:bg-red-700 dark:hover:bg-red-600 transition cursor-pointer"
>
  ```html
  <button
    class="bg-red-600 dark:bg-red-500 text-white px-4 py-2 text-sm tracking-widest hover:bg-red-700 transition cursor-pointer"
  >
    ### 3. Form Inputs Clean, flat background. ```html
    <input
      type="text"
      class="block w-full border border-transparent bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 p-2.5 focus:border-gray-500 dark:focus:border-gray-400 focus:bg-white dark:focus:bg-gray-700 focus:outline-none focus:ring-0"
    />
  </button>
</button>
````

## Layout Guidelines

- **Container**: Use `container mx-auto px-3`.
- **Spacing**: Use consistent spacing scale (`p-4`, `m-4`, `gap-6`).
- **Grid**: Prefer `grid` over `flex` for page layouts. `grid-cols-1 md:grid-cols-2 lg:grid-cols-3`.

## Do's and Don'ts

- **✅ Do**: Use `dark:` variants for EVERY color utility.
- **✅ Do**: Use `border-transparent` on inputs to prevent layout shifts on focus.
- **✅ Do**: Use `tracking-widest` for buttons to match the "iDiary.ai" aesthetic.
- **❌ Don't**: Use heavy box-shadows (`shadow-xl`). The design is flat/bordered.
- **❌ Don't**: Use rounded corners larger than `rounded` or `rounded-lg` (except for pills/avatars).
- **❌ Don't**: Mix serif and sans-serif fonts randomly.
