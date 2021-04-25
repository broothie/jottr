import React, {useEffect, useState} from 'react'
import {Link, useHistory, useRouteMatch} from 'react-router-dom'
import * as api from '../api'
import Quill from "quill";
import * as Cookie from '../cookie'

const inputDelayMilliseconds = 250
const whitespaceRegexp = /^\s*$/

const SAVED = 'saved'
const NOT_SAVED = 'not saved'
const SAVING = 'saving...'

const quillConfig = {
  theme: 'bubble',
  modules: {
    toolbar: [
      ['bold', 'italic', 'underline', 'strike', 'code-block'],
      [{ list: 'ordered' }, { list: 'bullet' }],
      [{ header: 1 }, { header: 2 }],
      ['clean']
    ]
  }
}

export default function Jot() {
  let quill
  const history = useHistory()
  const routeMatch = useRouteMatch('/jot/:jotId')
  const jotId = routeMatch.params.jotId
  const [savedStatus, setSavedStatus] = useState(SAVED)

  // Save jot to db
  function save() {
    if (!quill) return

    const delta = quill.getContents()
    const title = quill.getText()
      .split('\n')
      .find((line) => !whitespaceRegexp.test(line))

    setSavedStatus(SAVING)
    api.updateJot(jotId, { title, delta }).then(() => setSavedStatus(SAVED))
  }

  // Start Quill
  function initializeQuill() {
    quill = new Quill('#quill', quillConfig)

    let lastTextChangeAt = Date.now()
    quill.on('text-change', (_delta, _oldContents, source) => {
      if (source !== 'user') return

      const currentTextChangeAt = Date.now()
      if (currentTextChangeAt > lastTextChangeAt) lastTextChangeAt = currentTextChangeAt

      setSavedStatus(NOT_SAVED)
      setTimeout(() => {
        if (lastTextChangeAt < Date.now() - inputDelayMilliseconds) save()
      }, inputDelayMilliseconds)
    })
  }

  // Get jot from db
  function getJot() {
    api.getJot(jotId)
      .then(({ delta }) => quill.setContents(delta))
      .then(() => quill.setSelection(quill.getText().length))
      .then(() => Cookie.addJotId(jotId))
      .catch(() => history.push('/home'))
  }

  // Focus quill editor
  function focusEditor() {
    quill?.focus()
  }

  // Delete jot
  function deleteJot() {
    api.deleteJot(jotId).then(() => history.push('/home'))
  }

  // Lifecycle
  useEffect(() => {
    // After mount
    initializeQuill()
    getJot()

    // Before unmount
    return save
  }, [])

  // Markup
  return <div className="jot-page">
    <div className="nav-bar">
      <p className="saved-status">{savedStatus}</p>
      <Link className="button" to="/">new</Link>
      <button className="button" onClick={deleteJot}>delete</button>
      <Link className="button" to="/home">home</Link>
    </div>

    <div className="quill-container">
      <div id="quill" onClick={focusEditor}/>
    </div>
  </div>
}
