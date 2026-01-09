;(function () {
  // 1. Pre-configure Tailwind v4
  // Must run before the Tailwind CDN script to override defaults
  window.tailwind = window.tailwind || {}
  window.tailwind.config = {
    darkMode: 'class',
    theme: {
      extend: {}
    }
  }

  // 2. Inject Tailwind v4 config for class-based dark mode support
  // This ensures the @variant dark works correctly with the CDN version
  const style = document.createElement('style')
  style.type = 'text/tailwindcss'
  style.textContent = `
    @variant dark (&:where(.dark, .dark *));
  `
  document.head.appendChild(style)

  // 3. Immediate Theme Application (Blocks rendering to prevent flash)
  // This runs synchronously in <head>
  const getPreferredThemeIsDark = () => {
    try {
      if (localStorage.theme === 'dark') return true
      if (localStorage.theme === 'light') return false
    } catch (err) {
      void err
    }
    return !!(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
  }

  const applyTheme = (isDark) => {
    document.documentElement.classList.toggle('dark', isDark)
  }

  // Apply immediately
  applyTheme(getPreferredThemeIsDark())

  // 4. Setup Theme Toggle Logic
  // We attach a global listener for delegation
  if (typeof document !== 'undefined') {
    document.addEventListener('click', (event) => {
      const toggleButton = event.target.closest('[data-theme-toggle]')
      if (!toggleButton) return

      const nextIsDark = !document.documentElement.classList.contains('dark')
      applyTheme(nextIsDark)
      try {
        localStorage.theme = nextIsDark ? 'dark' : 'light'
      } catch (err) {
        void err
      }
    })
  }
})()
