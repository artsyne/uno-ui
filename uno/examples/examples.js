;(function () {
  // Inject Tailwind v4 config for class-based dark mode
  const style = document.createElement('style')
  style.type = 'text/tailwindcss'
  style.textContent = `
    @variant dark (&:where(.dark, .dark *));
  `
  document.head.appendChild(style)

  window.tailwind = window.tailwind || {}
  window.tailwind.config = { darkMode: 'class' }

  const applyTheme = (isDark) => {
    document.documentElement.classList.toggle('dark', isDark)
  }

  const getPreferredThemeIsDark = () => {
    try {
      if (localStorage.theme === 'dark') return true
      if (localStorage.theme === 'light') return false
    } catch (err) {
      void err
    }

    return !!(window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
  }

  applyTheme(getPreferredThemeIsDark())

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
})()
