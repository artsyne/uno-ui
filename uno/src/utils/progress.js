/**
 * Uno - Progress Bar Utils
 *
 * @module uno.core.progress
 * @version 0.1.0
 * @license MIT
 */

let loopHandle = null
let running = false

/**
 * Start a progress bar loop simulation
 * @param {HTMLElement} fillEl - Progress bar element
 * @param {number} [durationMs=1000] - Animation duration
 */
export function startProgressLoop(fillEl, durationMs = 1000) {
  if (!fillEl) return
  running = true
  const run = () => {
    if (!running) return
    fillEl.style.transition = 'none'
    fillEl.style.width = '0%'
    void fillEl.offsetHeight
    setTimeout(() => {
      if (!running) return
      fillEl.style.transition = `width ${durationMs}ms ease-in-out`
      fillEl.style.width = '100%'
      loopHandle = setTimeout(run, durationMs + 200)
    }, 50)
  }
  run()
}

/**
 * Stop the progress loop
 * @param {HTMLElement} fillEl - Progress bar element
 */
export function stopProgressLoop(fillEl) {
  running = false
  if (loopHandle) {
    clearTimeout(loopHandle)
    loopHandle = null
  }
  if (fillEl) fillEl.style.width = '0%'
}
