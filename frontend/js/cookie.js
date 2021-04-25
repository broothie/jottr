import Cookies from 'js-cookie'

const jotIdCookieName = 'jot_ids'

export function getJotIds() {
  return get(jotIdCookieName, [])
}

export function addJotIds(...jotIds) {
  const jotIdSet = new Set(getJotIds())
  jotIds.forEach((jotId) => jotIdSet.add(jotId))
  set(jotIdCookieName, Array.from(jotIdSet))
}

export function removeJotIds(...jotIds) {
  const jotIdSet = new Set(getJotIds())
  jotIds.forEach((jotId) => jotIdSet.delete(jotId))
  set(jotIdCookieName, Array.from(jotIdSet))
}

export function get(name, fallback = {}) {
  let value = Cookies.get(name)
  if (!value) return fallback

  return JSON.parse(atob(value))
}

export function set(name, value) {
  return Cookies.set(name, btoa(JSON.stringify(value)))
}
