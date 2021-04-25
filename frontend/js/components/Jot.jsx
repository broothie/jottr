import React, {useEffect, useState} from 'react'
import {Link, useHistory, useRouteMatch} from 'react-router-dom'
import * as Api from '../api'
import Quill from "quill";
import * as Cookie from '../cookie'
import setTitle from "../title";

const inputDelayMilliseconds = 500
const inputDelayOffsetMilliseconds = 10
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
  let typingTimeout, quill
  const history = useHistory()
  const routeMatch = useRouteMatch('/jot/:jotId')
  const jotId = routeMatch.params.jotId

  const [savedStatus, setSavedStatus] = useState(SAVED)
  const [saveOnClose, setSaveOnClose] = useState(true)

  // Save jot to db
  function save() {
    if (savedStatus === SAVED)
    if (!quill) return

    setSavedStatus(SAVING)
    const delta = quill.getContents()
    const title = quill.getText()
      .split('\n')
      .find((line) => !whitespaceRegexp.test(line))

    Api.updateJot(jotId, { title, delta })
      .then(() => {
        setSavedStatus(SAVED)
        setTitle(title || jotId)
      })
  }

  // Start Quill
  function initializeQuill() {
    quill = new Quill('#quill', quillConfig)

    quill.on('text-change', (_delta, _oldContents, source) => {
      if (source !== 'user') return

      clearTimeout(typingTimeout)
      setSavedStatus(NOT_SAVED)
      typingTimeout = setTimeout(save, inputDelayMilliseconds)
    })
  }

  // Get jot from db
  function getJot() {
    Api.getJot(jotId)
      .then(({ title, delta }) => {
        quill.setContents(delta)
        quill.setSelection(quill.getText().length)
        setTitle(title || jotId)
        Cookie.addJotIds(jotId)
      })
      .catch(() => history.push('/home'))
  }

  // Delete jot
  function deleteJot() {
    Api.deleteJot(jotId)
      .then(() => setSaveOnClose(false))
      .then(() => history.push('/home'))
  }

  // Lifecycle
  useEffect(() => {
    // After mount
    initializeQuill()
    getJot()

    // Before unmount
    return () => saveOnClose && save()
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
      <div id="quill"/>
    </div>
  </div>
}
