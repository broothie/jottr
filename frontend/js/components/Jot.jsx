import React, {useEffect, useState} from 'react'
import {Link, useHistory} from 'react-router-dom'
import * as Api from '../api'
import Quill from 'quill'
import * as Cookie from '../cookie'
import setSubtitle from '../title'

const inputDelayMilliseconds = 500
const whitespaceRegexp = /^\s*$/

const SAVED = 'saved'
const NOT_SAVED = 'not saved'
const SAVING = 'saving...'

const quillConfig = {
  theme: 'snow',
  modules: { toolbar: '#toolbar' },
  placeholder: 'jot something...',
}

export default function Jot(props) {
  let typingTimeout, quill
  let saveOnClose = true
  const history = useHistory()
  const jotId = props.match.params.jotId

  const [savedStatus, setSavedStatus] = useState(SAVED)

  function getTitle() {
    if (!quill) return ''

    return quill
      .getText()
      .split('\n')
      .find((line) => !whitespaceRegexp.test(line))
      || ''
  }

  function updateTitle() {
    const title = getTitle()
    setSubtitle(title || jotId)

    if (title) {
      const normalizedTitle = title
        .replace(/ /g, '-')
        .replace(/[^\w-]/g, '')

      history.push(`/${jotId}/${normalizedTitle}`)
    } else {
      history.push(`/${jotId}`)
    }
  }

  function initializeQuill(jot) {
    quill = new Quill('#quill', quillConfig)
    quill.setContents(jot.delta)
    quill.setSelection(quill.getText().length)
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
      .catch(() => history.push('/home'))
      .then((jot) => {
        initializeQuill(jot)
        updateTitle()
        Cookie.addJotIds(jotId)
      })
      .catch(console.log)
  }

  // Save jot to db
  function save(shouldUpdateTitle = true) {
    if (!quill) return

    setSavedStatus(SAVING)
    const delta = quill.getContents()
    const title = getTitle()

    Api.updateJot(jotId, { title, delta })
      .then(() => {
        setSavedStatus(SAVED)
        if (shouldUpdateTitle) updateTitle()
      })
  }

  // Delete jot
  function deleteJot() {
    Api.deleteJot(jotId)
      .then(() => {
        saveOnClose = false
        Cookie.removeJotIds(jotId)
        history.push('/home')
      })
  }

  useEffect(() => {
    getJot()

    return function umount() {
      if (saveOnClose) save(false)
    }
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

    <div className="toolbar-container">
      <div id="toolbar">
        <span className="toolbar-section ql-formats">
          <button className="ql-bold"/>
          <button className="ql-italic"/>
          <button className="ql-underline"/>
          <button className="ql-strike"/>
          <button className="ql-code-block"/>
        </span>

          <span className="toolbar-section ql-formats">
          <button className="ql-header" value="1"/>
          <button className="ql-header" value="2"/>
        </span>

          <span className="toolbar-section ql-formats">
          <button className="ql-list" value="ordered"/>
          <button className="ql-list" value="bullet"/>
        </span>

          <span className="toolbar-section ql-formats">
          <button className="ql-blockquote"/>
          <button className="ql-link"/>
        </span>

          <span className="toolbar-section ql-formats">
          <button className="ql-clean"/>
        </span>
      </div>
    </div>
  </div>
}
