/** HTMX 2 XHR response helpers; version follows Uno lockfile. */

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
 * @param {Object} event - HTMX 2 event
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

  return null
}

/**
 * Extract the htmx 2 XMLHttpRequest handle from an event.
 *
 * @param {Object} event - htmx event
 * @returns {XMLHttpRequest|undefined}
 */
export const extractXhr = (event) => event?.detail?.xhr
