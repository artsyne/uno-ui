/**
 * HTMX event-detail response helpers.
 *
 * htmx 2 exposes response data on `event.detail.xhr` (XMLHttpRequest:
 * status / responseText / getResponseHeader). htmx 4 replaces XHR with
 * fetch() and moves it to `event.detail.ctx` (ctx.response.status /
 * ctx.text / ctx.response.headers).
 *
 * All reads of response data from htmx events must go through these
 * helpers so the htmx 4 migration only needs to touch this file.
 * Keep in sync with website/static/src/js/utils/htmx-compat.js.
 */

const getXhrHeader = (xhr, name) => {
  try {
    return xhr.getResponseHeader(name)
  } catch {
    return null
  }
}

/**
 * Normalize response info from an htmx event.
 *
 * @param {Object} event - htmx event (htmx 2 / htmx 4 detail layouts supported)
 * @returns {?{status: number, text: string, header: function(string): ?string}}
 *   Null when no response data is available; callers must then keep their
 *   previous "no xhr" fallback behavior.
 */
export const getResponseInfo = (event) => {
  const detail = event?.detail

  // htmx 2: detail.xhr (XMLHttpRequest)
  const xhr = detail?.xhr
  if (xhr) {
    return {
      status: xhr.status,
      text: xhr.responseText ?? '',
      header: (name) => getXhrHeader(xhr, name)
    }
  }

  // htmx 4: detail.ctx (fetch-based)
  const response = detail?.ctx?.response
  if (response) {
    return {
      status: response.status,
      text: detail.ctx.text ?? '',
      header: (name) => {
        try {
          return response.headers?.get(name) ?? null
        } catch {
          return null
        }
      }
    }
  }

  return null
}

/**
 * Extract the htmx 2 XMLHttpRequest handle from an event.
 * Returns undefined under htmx 4 (fetch), so WeakSet-based request
 * deduplication is simply skipped there.
 *
 * @param {Object} event - htmx event
 * @returns {XMLHttpRequest|undefined}
 */
export const extractXhr = (event) => event?.detail?.xhr
