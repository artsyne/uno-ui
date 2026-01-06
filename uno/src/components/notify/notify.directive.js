export const registerNotifyDirectives = (Alpine) => {
  Alpine.directive('uno-notify', (el, { expression }, { evaluateLater }) => {
    // Usage: <button x-uno-notify="{ message: 'Message', type: 'success' }">
    // Usage: <button x-uno-notify="'Message'">

    if (!expression) {
      console.warn(
        '[Uno] x-uno-notify directive requires a configuration object or message string.'
      )
      return
    }

    const getExpression = evaluateLater(expression)

    const handler = () => {
      const store = Alpine.store('uno')?.notify
      if (!store) return

      getExpression((result) => {
        let message = ''
        let config = {}

        if (typeof result === 'string') {
          message = result
        } else if (typeof result === 'object' && result !== null) {
          message = result.message || ''
          config = result
        } else {
          console.warn('[Uno] x-uno-notify directive value must be a string or object', result)
          return
        }

        store.show(message, config)
      })
    }

    el.addEventListener('click', handler)
  })
}
