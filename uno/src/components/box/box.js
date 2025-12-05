/**
 * Uno - Box Component
 *
 * @module uno.components.box
 * @version 0.1.0
 * @license MIT
 */

import { getPositionStyle } from '../../utils/overlayContainers'
import { lock, unlock } from '../../utils/bodyScroll'
import { get as httpGet } from '../../utils/ajax'
import { log } from '../../utils/log'

export default (Alpine) => {
  const cached = {}
  let idPrefix = 'uno-box-'

  function genId() {
    return Math.random().toString(36).slice(2, 10)
  }

  function sel(id, name) {
    return `#${idPrefix}${name}${id}`
  }

  function parseData(el) {
    const d = el.dataset
    const map = {
      boxTitle: 'title',
      boxContent: 'content',
      boxWidth: 'width',
      boxHeight: 'height',
      boxModal: 'modal',
      boxPosition: 'position',
      boxAjaxUrl: 'ajaxUrl',
      boxIframeUrl: 'iframeUrl',
      boxOkText: 'okButtonText',
      boxCancelText: 'cancelButtonText',
      boxCache: 'cache',
      boxId: 'boxId',
      boxWidthClass: 'widthClass',
      boxHeightClass: 'heightClass',
      boxContainerClass: 'containerClass',
      boxContentOverflow: 'contentOverflow',
      boxButtonAlign: 'buttonAlign'
    }
    const out = {}
    for (const k in map) {
      const v = d[k]
      if (v !== undefined) {
        const key = map[k]
        if (key === 'width' || key === 'height') out[key] = parseInt(v, 10)
        else if (key === 'modal' || key === 'cache') out[key] = v === 'true' || v === true
        else out[key] = v
      }
    }
    return out
  }

  function t(key, fallback) {
    try {
      const s = Alpine?.store('i18n')
      return (s && s.t(`box.${key}`)) || fallback
    } catch { return fallback }
  }

  function build(opts) {
    const id = opts.uniqueId
    const hasTitle = !!opts.title
    const containerClasses = (opts.containerClass || '').trim()
    const header = hasTitle ? `
      <div id="${idPrefix}header${id}" class="flex justify-between items-center p-4 border-b border-gray-300 dark:border-gray-600">
        <div id="${idPrefix}title${id}" class="text-lg font-semibold text-gray-900 dark:text-gray-100"></div>
        <button id="${idPrefix}close${id}" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200" title="${t('close','Close')}">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      </div>
    ` : ''
    const noTitleClose = !hasTitle ? `
      <button id="${idPrefix}close${id}" class="absolute top-2 right-2 z-[101] text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 bg-white dark:bg-gray-900 p-1" title="Close">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
      </button>
    ` : ''
    const wrapperClasses = hasTitle ? 'p-4 text-gray-900 dark:text-gray-100' : 'pt-12 px-4 pb-4 text-gray-900 dark:text-gray-100'
    const outer = document.createElement('div')
    outer.id = `${idPrefix}outer-container${id}`
    outer.className = 'hidden fixed inset-0 z-[103] p-4 flex items-center justify-center'
    const container = document.createElement('div')
    container.id = `${idPrefix}container${id}`
    container.className = `bg-white dark:bg-gray-900 border-8 border-gray-200 dark:border-gray-600 shadow-lg overflow-hidden transition-colors duration-300 ${containerClasses}`
    container.style.position = 'relative'
    container.innerHTML = `${header}${noTitleClose}<div id="${idPrefix}wrapper${id}" class="${wrapperClasses}"><div id="${idPrefix}content${id}" class="overflow-auto"></div></div><div id="${idPrefix}footer${id}" class="p-4 border-t border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-100 hidden"></div>`
    outer.appendChild(container)
    const overlay = document.createElement('div')
    overlay.id = `${idPrefix}overlay${id}`
    overlay.className = 'hidden fixed inset-0 bg-black z-[102]'
    overlay.style.opacity = String(opts.overlayOpacity ?? 0.5)
    document.body.appendChild(outer)
    document.body.appendChild(overlay)
    return { outer, container, overlay }
  }

  function applySizeAndOverflow(id, opts) {
    const container = document.querySelector(sel(id, 'container'))
    const content = document.querySelector(sel(id, 'content'))
    const wrapper = document.querySelector(sel(id, 'wrapper'))
    if (!container || !content) return
    if (opts.widthClass || opts.heightClass) {
      container.style.width = ''
      container.style.height = ''
      container.style.maxWidth = ''
      container.style.maxHeight = ''
      if (opts.widthClass) container.classList.add(...opts.widthClass.split(' '))
      if (opts.heightClass) {
        container.classList.add(...opts.heightClass.split(' '))
        if (wrapper) {
          wrapper.classList.add('flex', 'flex-col', 'h-full')
          content.classList.add('flex-1', 'min-h-0')
        }
      }
    } else {
      const md = 768
      const ww = window.innerWidth
      if (ww <= md) {
        container.style.width = 'calc(100% - 2rem)'
        container.style.height = 'calc(100% - 2rem)'
        if (wrapper) {
          wrapper.classList.add('flex', 'flex-col', 'h-full')
          content.classList.add('flex-1', 'min-h-0')
        } else {
          content.style.height = '100%'
        }
      } else {
        if (opts.width) container.style.width = `${opts.width}px`
        if (opts.height) content.style.height = `${opts.height}px`
      }
    }
    content.classList.remove('overflow-auto', 'overflow-hidden', 'overflow-visible')
    const ov = opts.contentOverflow || 'auto'
    if (ov === 'hidden') content.classList.add('overflow-hidden')
    else if (ov === 'visible') content.classList.add('overflow-visible')
    else content.classList.add('overflow-auto')
  }

  function setPosition(id, position) {
    const outer = document.querySelector(sel(id, 'outer-container'))
    if (!outer) return
    const style = getPositionStyle(position || 'center')
    Object.assign(outer.style, style)
    outer.classList.remove('hidden')
  }

  function bindClose(id, modal) {
    const btn = document.querySelector(sel(id, 'close'))
    if (btn) {
      btn.addEventListener('click', () => close({ uniqueId: id, modal }))
    }
  }

  async function loadContent(id, opts) {
    const content = document.querySelector(sel(id, 'content'))
    if (!content) return
    if (typeof opts.content === 'string') {
      content.innerHTML = opts.content
    } else if (typeof opts.targetId === 'string') {
      const src = document.getElementById(opts.targetId)
      if (src) content.innerHTML = src.innerHTML
    } else if (typeof opts.ajaxUrl === 'string') {
      content.innerHTML = `<div class="flex justify-center items-center h-full"><div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 dark:border-gray-100"></div><span class="ml-2 text-gray-900 dark:text-gray-100">${t('loading','Loading')}</span></div>`
      try {
        const data = await httpGet(opts.ajaxUrl)
        if (typeof data === 'object' && data !== null) {
          content.innerHTML = `<pre class="bg-gray-100 p-3 rounded overflow-x-auto whitespace-pre-wrap">${JSON.stringify(data, null, 2)}</pre>`
        } else {
          content.innerHTML = String(data)
        }
      } catch (e) {
        log.error('box ajax error', e)
      }
    } else if (typeof opts.iframeUrl === 'string') {
      content.innerHTML = `<iframe frameborder="0" class="w-full h-full" src="${opts.iframeUrl}"></iframe>`
    }
  }

  function renderButtons(id, type, opts) {
    const footer = document.querySelector(sel(id, 'footer'))
    if (!footer) return
    const align = ({ left: 'justify-start', center: 'justify-center', right: 'justify-end' }[opts.buttonAlign] || 'justify-end')
    footer.classList.remove('hidden')
    footer.className = `p-4 border-t border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-100 flex ${align} space-x-4`
    const okText = opts.okButtonText || t('ok','OK')
    const cancelText = opts.cancelButtonText || t('cancel','Cancel')
    const okNeeded = type === 'confirm' || type === 'message' || (type === 'dialog' && opts.showOkButton === true)
    const cancelNeeded = type === 'confirm' || (type === 'dialog' && opts.showCancelButton === true)
    footer.innerHTML = ''
    if (okNeeded) {
      const ok = document.createElement('button')
      ok.id = `${idPrefix}ok${id}`
      ok.className = 'px-4 py-2 bg-blue-500 text-white hover:bg-blue-600 transition-colors duration-200'
      ok.textContent = okText
      ok.addEventListener('click', () => {
        let shouldClose = true
        if (typeof opts.onOk === 'function') {
          const r = opts.onOk()
          shouldClose = r !== false
        }
        if (shouldClose) close({ uniqueId: id, modal: opts.modal !== false })
      })
      footer.appendChild(ok)
    }
    if (cancelNeeded) {
      const cancel = document.createElement('button')
      cancel.id = `${idPrefix}cancel${id}`
      cancel.className = 'px-4 py-2 bg-gray-300 dark:bg-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-400 dark:hover:bg-gray-500 transition-colors duration-200'
      cancel.textContent = cancelText
      cancel.addEventListener('click', () => {
        close({ uniqueId: id, modal: opts.modal !== false })
        if (typeof opts.onCancel === 'function') opts.onCancel()
      })
      footer.appendChild(cancel)
    }
    if (!okNeeded && !cancelNeeded) {
      footer.classList.add('hidden')
    }
  }

  /**
   * Open a box
   * @param {string} type - Box type (dialog, confirm, message)
   * @param {Object} opts - Options
   * @returns {Promise<string>} Box ID
   */
  async function open(type, opts) {
    const options = {
      modal: true,
      position: 'center',
      width: 400,
      height: undefined,
      overlayOpacity: 0.5,
      responsive: true,
      contentOverflow: 'auto',
      ...opts
    }
    if (typeof options.onBeforeOpen === 'function') options.onBeforeOpen()
    let id = options.boxId || genId()
    options.uniqueId = id
    if (options.cache && options.boxId && cached[options.boxId]) {
      id = cached[options.boxId].uniqueId
      document.querySelector(sel(id, 'outer-container'))?.classList.remove('hidden')
      document.querySelector(sel(id, 'overlay'))?.classList.toggle('hidden', options.modal === false)
      if (options.modal !== false) lock()
      renderButtons(id, type, options)
      bindClose(id, options.modal !== false)
      if (typeof options.onOpen === 'function') options.onOpen()
      return id
    }
    const nodes = build(options)
    if (options.cache && options.boxId) cached[options.boxId] = { uniqueId: id, modal: options.modal !== false }
    if (options.title) {
      const tEl = document.querySelector(sel(id, 'title'))
      if (tEl) tEl.textContent = options.title
    }
    if (options.modal !== false) {
      nodes.overlay.classList.remove('hidden')
      lock()
    }
    setPosition(id, options.position)
    applySizeAndOverflow(id, options)
    await loadContent(id, options)
    renderButtons(id, type, options)
    bindClose(id, options.modal !== false)
    if (typeof options.onOpen === 'function') options.onOpen()
    return id
  }

  /**
   * Close a box
   * @param {Object} params - { uniqueId, modal }
   */
  function close({ uniqueId, modal }) {
    document.querySelector(sel(uniqueId, 'outer-container'))?.classList.add('hidden')
    document.querySelector(sel(uniqueId, 'overlay'))?.classList.add('hidden')
    if (modal) unlock()
  }

  /**
   * Clear all cached boxes
   */
  function clearCache() {
    Object.keys(cached).forEach((key) => {
      const id = cached[key].uniqueId
      document.querySelector(sel(id, 'outer-container'))?.remove()
      document.querySelector(sel(id, 'overlay'))?.remove()
      delete cached[key]
    })
  }

  /**
   * Clear a specific cached box
   * @param {string} boxId 
   */
  function clearCacheById(boxId) {
    const rec = cached[boxId]
    if (!rec) return
    const id = rec.uniqueId
    document.querySelector(sel(id, 'outer-container'))?.remove()
    document.querySelector(sel(id, 'overlay'))?.remove()
    delete cached[boxId]
  }

  function bindDelegation() {
    document.addEventListener('click', (e) => {
      const t = e.target.closest('[data-box-type]')
      if (!t) return
      e.preventDefault()
      const type = t.getAttribute('data-box-type')
      if (!['dialog', 'confirm', 'message'].includes(type)) return
      const opts = parseData(t)
      open(type, opts)
    })
  }

  bindDelegation()

  Alpine.store('box', {
    open,
    close,
    clearCache,
    clearCacheById
  })
}
