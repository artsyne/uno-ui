;(function () {
  // Uno Configuration
  // Theme logic has been moved to theme.js to prevent FOUC (Flash of Wrong Theme)

  // Configure Uno using setConfig()
  // Note: We use a defensive pattern here to ensure Uno is loaded
  // because examples.js is loaded synchronously while Uno is deferred.
  const initExamples = () => {
    if (window.Uno) {
      console.log('[Examples] Configuring Uno...')
      window.Uno.setConfig({
        debug: true
      })
    } else {
      console.warn('[Examples] Uno not found during initialization')
    }
  }

  // If Uno is already available (unlikely if deferred), run immediately
  if (window.Uno) {
    initExamples()
  } else {
    // Wait for DOMContentLoaded which fires after defer scripts execute
    document.addEventListener('DOMContentLoaded', initExamples)
  }

  // Uno will auto-register with Alpine when it loads
})()
