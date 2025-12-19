/**
 * Uno - DOM Utils
 *
 * @module uno.utils.dom
 */

export function escapeHTML(input) {
  return String(input ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

const onceMarks = new WeakMap()

export function once(target, key) {
  if (!target) return false
  const k = String(key)
  let set = onceMarks.get(target)
  if (!set) {
    set = new Set()
    onceMarks.set(target, set)
  }
  if (set.has(k)) return false
  set.add(k)
  return true
}

export function sanitizeHTML(html, options = {}) {
  const {
    allowedTags = [
      'b',
      'i',
      'strong',
      'em',
      'u',
      'br',
      'p',
      'ul',
      'ol',
      'li',
      'code',
      'pre',
      'a',
      'span',
      'div'
    ],
    allowedAttrs = { a: ['href', 'target', 'rel'], '*': ['class'] }
  } = options

  const template = document.createElement('template')
  template.innerHTML = String(html ?? '')

  const walker = document.createTreeWalker(template.content, NodeFilter.SHOW_ELEMENT)
  const toRemove = []
  const toStrip = [] // Tags to remove completely including content

  while (walker.nextNode()) {
    const el = walker.currentNode
    const tag = el.tagName.toLowerCase()

    // Tags to strip completely
    if (['script', 'style', 'iframe', 'object', 'embed'].includes(tag)) {
      toStrip.push(el)
      continue
    }

    // 中文说明：默认白名单策略。非白名单标签将被降级为纯文本，避免 XSS。
    if (!allowedTags.includes(tag)) {
      toRemove.push(el)
      continue
    }

    for (const attr of [...el.attributes]) {
      const name = attr.name.toLowerCase()
      const allow = (allowedAttrs[tag] || []).concat(allowedAttrs['*'] || [])
      if (!allow.includes(name)) el.removeAttribute(attr.name)
    }

    // 中文说明：href 只允许安全协议/站内链接，防止 javascript: 等注入。
    if (tag === 'a') {
      const href = (el.getAttribute('href') || '').trim()
      const isSafe =
        href.startsWith('/') ||
        href.startsWith('https://') ||
        href.startsWith('http://') ||
        href.startsWith('#')
      if (!isSafe) el.removeAttribute('href')
      if (!el.getAttribute('rel')) el.setAttribute('rel', 'noopener noreferrer')
    }
  }

  for (const node of toRemove) {
    node.replaceWith(document.createTextNode(node.textContent || ''))
  }

  return template.innerHTML
}

export function setHTML(el, html, { trusted = false, alpine } = {}) {
  if (!el) return
  el.innerHTML = trusted ? String(html ?? '') : sanitizeHTML(html)
  alpine?.initTree?.(el)
}
