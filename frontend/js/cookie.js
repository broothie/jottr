import Cookies from 'js-cookie'

const jotIdCookieName = 'jot_ids'

export function addJotId(jotId) {
  const jotIds = new Set(getJotIds())
  jotIds.add(jotId)
  set(jotIdCookieName, Array.from(jotIds))
}

export function getJotIds() {
  return get(jotIdCookieName, [])
}

export function set(name, value) {
  return Cookies.set(name, btoa(JSON.stringify(value)))
}

export function get(name, fallback = {}) {
  return JSON.parse(atob(Cookies.get(name) || btoa(JSON.stringify(fallback))))
}
