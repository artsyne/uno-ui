import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { once, sanitizeHTML, setHTML, DomManager } from '../../src/utils/dom'

describe('DOM Utils - once()', () => {
  let testEl

  beforeEach(() => {
    testEl = document.createElement('div')
  })

  it('should return true for first call with a key', () => {
    expect(once(testEl, 'init')).toBe(true)
  })

  it('should return false for repeated calls with same key', () => {
    once(testEl, 'init')
    expect(once(testEl, 'init')).toBe(false)
    expect(once(testEl, 'init')).toBe(false)
  })

  it('should allow different keys on same element', () => {
    expect(once(testEl, 'init')).toBe(true)
    expect(once(testEl, 'setup')).toBe(true)
    expect(once(testEl, 'init')).toBe(false)
    expect(once(testEl, 'setup')).toBe(false)
  })

  it('should return false for null/undefined target', () => {
    expect(once(null, 'key')).toBe(false)
    expect(once(undefined, 'key')).toBe(false)
  })

  it('should handle numeric keys correctly', () => {
    expect(once(testEl, 123)).toBe(true)
    expect(once(testEl, 123)).toBe(false)
  })
})

describe('DOM Utils - sanitizeHTML()', () => {
  it('should remove script tags', () => {
    const html = '<div>Safe</div><script>alert("xss")</script>'
    const result = sanitizeHTML(html)
    expect(result).not.toContain('<script>')
    expect(result).toContain('<div>Safe</div>')
  })

  it('should remove script tags with attributes', () => {
    const html = '<script type="text/javascript">alert(1)</script>'
    const result = sanitizeHTML(html)
    expect(result).not.toContain('<script>')
  })

  it('should remove img tags (to prevent onerror XSS)', () => {
    const html = '<img src="x" onerror="alert(1)">'
    const result = sanitizeHTML(html)
    expect(result).not.toContain('<img')
  })

  it('should remove on* event handlers', () => {
    const html = '<div onclick="alert(1)" onmouseover="alert(2)">Test</div>'
    const result = sanitizeHTML(html)
    expect(result).not.toContain('onclick')
    expect(result).not.toContain('onmouseover')
    expect(result).toContain('<div')
    expect(result).toContain('Test')
  })

  it('should remove javascript: protocol in href', () => {
    const html = '<a href="javascript:alert(1)">Click</a>'
    const result = sanitizeHTML(html)
    expect(result).not.toContain('javascript:')
  })

  it('should preserve safe HTML', () => {
    const html = '<div class="test"><p>Hello <strong>World</strong></p></div>'
    const result = sanitizeHTML(html)
    expect(result).toBe(html)
  })

  it('should handle empty/null input', () => {
    expect(sanitizeHTML('')).toBe('')
    expect(sanitizeHTML(null)).toBe('')
    expect(sanitizeHTML(undefined)).toBe('')
  })

  it('should remove mixed case event handlers', () => {
    const html = '<div OnClick="alert(1)" ONLOAD="alert(2)">Test</div>'
    const result = sanitizeHTML(html)
    expect(result).not.toContain('OnClick')
    expect(result).not.toContain('ONLOAD')
  })
})

describe('DOM Utils - setHTML()', () => {
  let container
  let mockAlpine

  beforeEach(() => {
    container = document.createElement('div')
    mockAlpine = {
      initTree: vi.fn()
    }
  })

  it('should set sanitized HTML by default', () => {
    const html = '<div>Safe</div><script>alert("xss")</script>'
    setHTML(container, html)
    expect(container.innerHTML).toContain('<div>Safe</div>')
    expect(container.innerHTML).not.toContain('<script>')
  })

  it('should set trusted HTML when trusted=true', () => {
    const html = '<div onclick="test()">Click</div>'
    setHTML(container, html, { trusted: true })
    expect(container.innerHTML).toContain('onclick')
  })

  it('should call Alpine.initTree if provided', () => {
    setHTML(container, '<div>Test</div>', { alpine: mockAlpine })
    expect(mockAlpine.initTree).toHaveBeenCalledWith(container)
  })

  it('should handle null element gracefully', () => {
    expect(() => setHTML(null, '<div>Test</div>')).not.toThrow()
  })

  it('should handle null/undefined HTML', () => {
    setHTML(container, null)
    expect(container.innerHTML).toBe('')

    setHTML(container, undefined)
    expect(container.innerHTML).toBe('')
  })
})

describe('DomManager - ensureContainer()', () => {
  afterEach(() => {
    document.querySelectorAll('[id^="test-"]').forEach((el) => el.remove())
  })

  it('should create container if not exists', () => {
    const container = DomManager.ensureContainer('test-container')
    expect(container).toBeDefined()
    expect(container.id).toBe('test-container')
    expect(document.getElementById('test-container')).toBe(container)
  })

  it('should return existing container if already exists', () => {
    const first = DomManager.ensureContainer('test-existing')
    const second = DomManager.ensureContainer('test-existing')
    expect(first).toBe(second)
  })

  it('should set HTML content', () => {
    const container = DomManager.ensureContainer('test-html', '<div class="inner">Content</div>')
    expect(container.innerHTML).toContain('Content')
  })

  it('should set attributes', () => {
    const container = DomManager.ensureContainer('test-attrs', '', {
      'data-test': 'value',
      class: 'custom-class'
    })
    expect(container.getAttribute('data-test')).toBe('value')
    expect(container.className).toBe('custom-class')
  })

  it('should auto-inject CSP nonce if available', () => {
    window.Uno = {
      options: {
        csp: {
          nonce: 'test-nonce-123'
        }
      }
    }

    const container = DomManager.ensureContainer('test-nonce')
    expect(container.getAttribute('nonce')).toBe('test-nonce-123')

    delete window.Uno
  })

  it('should not override explicit nonce in attributes', () => {
    window.Uno = {
      options: {
        csp: {
          nonce: 'auto-nonce'
        }
      }
    }

    const container = DomManager.ensureContainer('test-explicit-nonce', '', {
      nonce: 'explicit-nonce'
    })
    expect(container.getAttribute('nonce')).toBe('explicit-nonce')

    delete window.Uno
  })

  it('should append to custom parent if provided', () => {
    const parent = document.createElement('div')
    const container = DomManager.ensureContainer('test-parent', '', {}, parent)
    expect(container.parentElement).toBe(parent)
  })

  it('should append to body by default', () => {
    const container = DomManager.ensureContainer('test-default-parent')
    expect(container.parentElement).toBe(document.body)
  })
})

describe('DomManager - injectStyles()', () => {
  afterEach(() => {
    document.querySelectorAll('style[id^="test-"]').forEach((el) => el.remove())
  })

  it('should inject CSS styles', () => {
    DomManager.injectStyles('test-style', '.test { color: red; }')
    const style = document.getElementById('test-style')
    expect(style).toBeDefined()
    expect(style.tagName).toBe('STYLE')
    expect(style.textContent).toBe('.test { color: red; }')
  })

  it('should not inject duplicate styles', () => {
    DomManager.injectStyles('test-duplicate', '.test { color: blue; }')
    DomManager.injectStyles('test-duplicate', '.test { color: green; }')

    const styles = document.querySelectorAll('#test-duplicate')
    expect(styles.length).toBe(1)
    expect(styles[0].textContent).toBe('.test { color: blue; }') // First one kept
  })

  it('should inject CSP nonce if available', () => {
    window.Uno = {
      options: {
        csp: {
          nonce: 'style-nonce-456'
        }
      }
    }

    DomManager.injectStyles('test-style-nonce', '.test {}')
    const style = document.getElementById('test-style-nonce')
    expect(style.getAttribute('nonce')).toBe('style-nonce-456')

    delete window.Uno
  })

  it('should append to document.head', () => {
    DomManager.injectStyles('test-head', '.test {}')
    const style = document.getElementById('test-head')
    expect(style.parentElement).toBe(document.head)
  })
})

describe('DomManager - exports', () => {
  it('should export setHTML', () => {
    expect(DomManager.setHTML).toBe(setHTML)
  })

  it('should export sanitize', () => {
    expect(DomManager.sanitize).toBe(sanitizeHTML)
  })
})
