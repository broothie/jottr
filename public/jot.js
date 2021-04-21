
const jotId = document.currentScript.dataset.jotId
const inputDelayMilliseconds = 250

document.addEventListener('DOMContentLoaded', () => {
  const jotBodyContainer = document.getElementById('jot-body-container')
  const jotBody = document.getElementById('jot-body')
  const savedStatus = document.getElementById('saved-status')
  let lastInputAt = Date.now()

  scrollToBottom(jotBodyContainer)
  setCursorToEnd(jotBody)

  jotBody.addEventListener('input', () => {
    savedStatus.innerText = 'not saved'
    const currentInputAt = Date.now()
    if (currentInputAt > lastInputAt) lastInputAt = currentInputAt

    setTimeout(() => {
      if (lastInputAt > Date.now() - inputDelayMilliseconds) return

      savedStatus.innerText = 'saving...'
      const title = jotBody.innerText.split('\n')[0]
      const body = JSON.stringify({ title , body: jotBody.innerHTML })

      fetch(`/api/jots/${jotId}`, { method: 'put', body })
        .then(() => {
          savedStatus.innerText = 'saved'
          setTitle(title)
        })
    }, inputDelayMilliseconds)
  })

  document.getElementById('delete-jot').addEventListener('click', () => {
    if (!confirm(`Are you sure you want to delete this jot (id = ${jotId})`)) return

    fetch(`/api/jots/${jotId}`, { method: 'delete' }).then(() => location.href = '/home')
  })

  window.onbeforeunload = () => {
    if (jotBody.innerText.length === 0) fetch(`/api/jots/${jotId}`, { method: 'delete' })
  }
})

function setTitle(title) {
  if (!title || title.length === 0) title = jotId
  document.title = `jottr - ${title}`
}

function scrollToBottom(element) {
  element.scrollTop = element.scrollHeight
}

function setCursorToEnd(contentEditable) {
  const range = document.createRange()
  range.selectNodeContents(contentEditable)
  range.collapse(false)
  const selection = window.getSelection()
  selection.removeAllRanges()
  selection.addRange(range)
}
