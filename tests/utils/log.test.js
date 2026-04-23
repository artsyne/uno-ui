/**
 * Unit tests for log utility
 */
import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { log } from '../../src/utils/log.js'

describe('log', () => {
  let consoleLogSpy
  let consoleWarnSpy
  let consoleErrorSpy

  beforeEach(() => {
    // Reset debug mode
    log.setDebug(false)
    log.setReporter(null)

    // Mock console methods
    consoleLogSpy = vi.spyOn(console, 'log').mockImplementation(() => {})
    consoleWarnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})
    consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  afterEach(() => {
    consoleLogSpy.mockRestore()
    consoleWarnSpy.mockRestore()
    consoleErrorSpy.mockRestore()
  })

  describe('setDebug()', () => {
    it('should enable debug mode', () => {
      log.setDebug(true)
      log.info('test message')
      expect(consoleLogSpy).toHaveBeenCalledWith('[Uno Info]', 'test message')
    })

    it('should disable debug mode', () => {
      log.setDebug(false)
      log.info('test message')
      expect(consoleLogSpy).not.toHaveBeenCalled()
    })
  })

  describe('setReporter()', () => {
    it('should set custom error reporter', () => {
      const reporter = vi.fn()
      log.setReporter(reporter)
      log.error('error message')
      expect(reporter).toHaveBeenCalledWith('error message')
    })

    it('should replace existing reporter', () => {
      const reporter1 = vi.fn()
      const reporter2 = vi.fn()

      log.setReporter(reporter1)
      log.error('message 1')
      expect(reporter1).toHaveBeenCalledWith('message 1')

      log.setReporter(reporter2)
      log.error('message 2')
      expect(reporter1).not.toHaveBeenCalledWith('message 2')
      expect(reporter2).toHaveBeenCalledWith('message 2')
    })
  })

  describe('info()', () => {
    it('should not log when debug is off', () => {
      log.setDebug(false)
      log.info('test message')
      expect(consoleLogSpy).not.toHaveBeenCalled()
    })

    it('should log when debug is on', () => {
      log.setDebug(true)
      log.info('test message', { data: 123 })
      expect(consoleLogSpy).toHaveBeenCalledWith('[Uno Info]', 'test message', { data: 123 })
    })

    it('should handle multiple arguments', () => {
      log.setDebug(true)
      log.info('arg1', 'arg2', 'arg3')
      expect(consoleLogSpy).toHaveBeenCalledWith('[Uno Info]', 'arg1', 'arg2', 'arg3')
    })
  })

  describe('warn()', () => {
    it('should not log when debug is off', () => {
      log.setDebug(false)
      log.warn('warning message')
      expect(consoleWarnSpy).not.toHaveBeenCalled()
    })

    it('should log when debug is on', () => {
      log.setDebug(true)
      log.warn('warning message', { warning: true })
      expect(consoleWarnSpy).toHaveBeenCalledWith('[Uno Warn]', 'warning message', {
        warning: true
      })
    })

    it('should handle multiple arguments', () => {
      log.setDebug(true)
      log.warn('warn1', 'warn2')
      expect(consoleWarnSpy).toHaveBeenCalledWith('[Uno Warn]', 'warn1', 'warn2')
    })
  })

  describe('error()', () => {
    it('should not log to console when debug is off and no reporter', () => {
      log.setDebug(false)
      log.error('error message')
      expect(consoleErrorSpy).not.toHaveBeenCalled()
    })

    it('should log to console when debug is on', () => {
      log.setDebug(true)
      log.error('error message', new Error('test'))
      expect(consoleErrorSpy).toHaveBeenCalledWith(
        '[Uno Error]',
        'error message',
        expect.any(Error)
      )
    })

    it('should call custom reporter when set', () => {
      const reporter = vi.fn()
      log.setReporter(reporter)
      log.error('error message', { code: 500 })
      expect(reporter).toHaveBeenCalledWith('error message', { code: 500 })
    })

    it('should call reporter and console when both are available', () => {
      log.setDebug(true)
      const reporter = vi.fn()
      log.setReporter(reporter)

      log.error('error')
      expect(reporter).toHaveBeenCalledWith('error')
      expect(consoleErrorSpy).toHaveBeenCalledWith('[Uno Error]', 'error')
    })

    it('should handle reporter error gracefully', () => {
      log.setDebug(true)
      const failingReporter = vi.fn().mockImplementation(() => {
        throw new Error('Reporter failed')
      })
      log.setReporter(failingReporter)

      log.error('test error')
      expect(failingReporter).toHaveBeenCalledWith('test error')
      expect(consoleErrorSpy).toHaveBeenCalledWith(
        '[Uno Error] Reporter failed:',
        expect.any(Error)
      )
      expect(consoleErrorSpy).toHaveBeenCalledWith('[Uno Error]', 'test error')
    })

    it('should handle reporter error silently when debug is off', () => {
      log.setDebug(false)
      const failingReporter = vi.fn().mockImplementation(() => {
        throw new Error('Reporter failed')
      })
      log.setReporter(failingReporter)

      log.error('test error')
      expect(failingReporter).toHaveBeenCalledWith('test error')
      expect(consoleErrorSpy).not.toHaveBeenCalled()
    })

    it('should handle multiple arguments in error', () => {
      log.setDebug(true)
      log.error('error1', 'error2', { code: 500 })
      expect(consoleErrorSpy).toHaveBeenCalledWith('[Uno Error]', 'error1', 'error2', { code: 500 })
    })
  })
})
