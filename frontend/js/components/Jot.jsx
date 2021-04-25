import React, {useEffect, useState} from 'react'
import {Link, useHistory, useRouteMatch} from 'react-router-dom'
import * as Api from '../api'
import Quill from 'quill'
import * as Cookie from '../cookie'
import setTitle from '../title'

const inputDelayMilliseconds = 500
const whitespaceRegexp = /^\s*$/

const SAVED = 'saved'
const NOT_SAVED = 'not saved'
const SAVING = 'saving...'

const quillConfig = { theme: 'snow', modules: { toolbar: '#toolbar' } }

export default function Jot() {
  let typingTimeout, quill
  const history = useHistory()
  const routeMatch = useRouteMatch('/jot/:jotId')
  const jotId = routeMatch.params.jotId

  const [savedStatus, setSavedStatus] = useState(SAVED)
  const [saveOnClose, setSaveOnClose] = useState(true)

  // Save jot to db
  function save(updateTitle = true) {
    if (!quill) return

    setSavedStatus(SAVING)
    const delta = quill.getContents()
    const title = quill.getText()
      .split('\n')
      .find((line) => !whitespaceRegexp.test(line))

    Api.updateJot(jotId, { title, delta })
      .then(() => {
        setSavedStatus(SAVED)
        if (updateTitle) setTitle(title || jotId)
      })
  }

  // Start Quill
  function initializeQuill() {
    quill = new Quill('#quill', quillConfig)

    quill.on('text-change', (_delta, _oldContents, source) => {
      if (source !== 'user') return

      clearTimeout(typingTimeout)
      setSavedStatus(NOT_SAVED)
      typingTimeout = setTimeout(() => save(), inputDelayMilliseconds)
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
    return () => saveOnClose && save(false)
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

    <div id="toolbar">
      <span className="ql-formats">
        <button className="ql-bold"/>
        <button className="ql-italic"/>
        <button className="ql-underline"/>
        <button className="ql-strike"/>
        <button className="ql-code-block"/>
      </span>

      <span className="ql-formats">
        <button className="ql-header" value="1"/>
        <button className="ql-header" value="2"/>
      </span>

      <span className="ql-formats">
        <button className="ql-list" value="ordered"/>
        <button className="ql-list" value="bullet"/>
      </span>

      <span className="ql-formats">
        <button className="ql-clean"/>
      </span>
    </div>
  </div>
}
