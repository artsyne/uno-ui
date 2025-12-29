import { describe, it, expect, vi } from 'vitest'
import { debounce, throttle } from '../../src/utils/performance'

describe('Performance Utils', () => {
  it('should debounce function calls', async () => {
    vi.useFakeTimers()
    const fn = vi.fn()
    const debounced = debounce(fn, 100)

    debounced()
    debounced()
    debounced()

    expect(fn).not.toHaveBeenCalled()
    
    vi.runAllTimers()
    
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('should throttle function calls', () => {
    vi.useFakeTimers()
    const fn = vi.fn()
    const throttled = throttle(fn, 100)

    throttled()
    throttled()
    throttled()

    expect(fn).toHaveBeenCalledTimes(1)
    
    vi.advanceTimersByTime(100)
    throttled()
    expect(fn).toHaveBeenCalledTimes(2)
  })
})
