/**
 * Uno - Validation Utils
 *
 * @module uno.utils.validation
 */

export const VALIDATORS = {
  validname: /^[a-zA-Z0-9_-]+$/,
  en: /^[a-zA-Z]+$/,
  domain: /^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$/,
  mobile: /^\+?[\d\s-]{7,}$/,
  int: /^-?\d+$/,
  email: /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/,
  date: /^(\d{4}-\d{2}-\d{2}|\d{2}\/\d{2}\/\d{4}|\d{2}\.\d{2}\.\d{4})$/,
  url: /^(https?:\/\/)?([\da-z.-]+)\.([a-z.-]{2,6})([/\w .-]*)*\/?$/i,
  pwd: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/,
  strongpwd: /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/
}

export function sanitizeRegexString(str) {
  const match = str.match(/^\/(.*?)\/([gimsuy]*)$/)
  if (match) {
    return { pattern: match[1], flags: match[2] }
  }
  return { pattern: str, flags: '' }
}

export function safeParseJSON(str) {
  try {
    if (!str) return []
    return JSON.parse(str)
  } catch {
    return []
  }
}
