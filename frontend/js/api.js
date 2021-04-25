
export function createJot() {
  return fetch('/api/jots', { method: 'POST' })
    .then(handleResponse)
    .then((response) => response.json())
}

export function getJot(jotId) {
  return fetch(`/api/jots/${jotId}`)
    .then(handleResponse)
    .then((response) => response.json())
}

export function updateJot(jotId, body = {}) {
  return fetch(`/api/jots/${jotId}`, { method: 'PATCH', body: JSON.stringify(body) })
    .then(handleResponse)
}

export function deleteJot(jotId) {
  return fetch(`/api/jots/${jotId}`, { method: 'DELETE' })
    .then(handleResponse)
}

export function bulkGetJots(...jotIds) {
  return fetch(`/api/bulk/jots?jot_ids=${jotIds.join(',')}`)
    .then(handleResponse)
    .then((response) => response.json())
}

function handleResponse(response) {
  if (!response.ok) throw Error(response.statusText)
  return response
}
