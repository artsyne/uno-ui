import { defineConfig } from 'vitepress'

export default defineConfig({
  lang: 'en',
  title: 'Uno UI',
  description:
    'Lightweight, Store-First UI components for Alpine.js and Tailwind CSS 4, with optional HTMX 2 support.',
  base: '/uno-ui/',
  cleanUrls: true,
  lastUpdated: true,
  head: [
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:title', content: 'Uno UI' }],
    [
      'meta',
      {
        property: 'og:description',
        content:
          'Lightweight, Store-First UI components for Alpine.js and Tailwind CSS 4, with optional HTMX 2 support.'
      }
    ],
    ['meta', { property: 'og:url', content: 'https://artsyne.github.io/uno-ui/' }],
    ['meta', { name: 'twitter:card', content: 'summary' }]
  ],
  themeConfig: {
    nav: [
      { text: 'Guide', link: '/getting-started' },
      { text: 'Components', link: '/getting-started' },
      { text: 'Docs', link: '/README' },
      { text: 'GitHub', link: 'https://github.com/artsyne/uno-ui' }
    ],
    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Getting Started', link: '/getting-started' },
          { text: 'Core Concepts', link: '/concepts' },
          { text: 'Configuration', link: '/configuration' },
          { text: 'Theming and Tailwind', link: '/theming' },
          { text: 'i18n', link: '/i18n' }
        ]
      },
      {
        text: 'Components',
        items: [
          { text: 'Box', link: '/box-api' },
          { text: 'Notify', link: '/notify-api' },
          { text: 'Tooltip', link: '/tooltip-api' },
          { text: 'Tab', link: '/tab-api' },
          { text: 'Valid', link: '/valid-api' }
        ]
      },
      {
        text: 'Integrations',
        items: [
          { text: 'HTMX', link: '/htmx' },
          { text: 'CSP', link: '/csp' },
          { text: 'Accessibility', link: '/accessibility' },
          { text: 'TypeScript', link: '/typescript-usage' }
        ]
      },
      {
        text: 'Operations',
        items: [
          { text: 'Troubleshooting', link: '/troubleshooting' },
          { text: 'Error Handling', link: '/error-handling' },
          { text: 'Migration', link: '/migration' },
          { text: 'Test Reports', link: '/test-reports' }
        ]
      },
      {
        text: 'Project',
        items: [
          { text: 'Documentation Hub', link: '/README' },
          {
            text: 'Release Notes',
            link: 'https://github.com/artsyne/uno-ui/blob/main/RELEASE.md'
          }
        ]
      }
    ],
    search: {
      provider: 'local'
    },
    socialLinks: [{ icon: 'github', link: 'https://github.com/artsyne/uno-ui' }],
    editLink: {
      pattern: 'https://github.com/artsyne/uno-ui/edit/main/docs/:path',
      text: 'Edit this page on GitHub'
    },
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 artsyne'
    }
  }
})
