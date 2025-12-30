export function isPlainObject(value) {
  return Object.prototype.toString.call(value) === '[object Object]'
}

export function mergeDeep(base, patch) {
  if (!isPlainObject(base)) return isPlainObject(patch) ? { ...patch } : patch
  if (!isPlainObject(patch)) return patch === undefined ? { ...base } : patch

  const out = { ...base }
  for (const [key, value] of Object.entries(patch)) {
    if (isPlainObject(value) && isPlainObject(out[key])) {
      out[key] = mergeDeep(out[key], value)
      continue
    }
    if (isPlainObject(value)) {
      out[key] = mergeDeep({}, value)
      continue
    }
    // Handle string append (e.g. "+class")
    if (
      typeof value === 'string' &&
      value.trim().indexOf('+') === 0 &&
      typeof out[key] === 'string'
    ) {
      out[key] = `${out[key]} ${value.trim().slice(1)}`.trim()
      continue
    }
    out[key] = value
  }
  return out
}
