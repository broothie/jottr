import React, {useEffect, useState} from 'react'
import {Link} from 'react-router-dom'
import * as Api from '../api'
import * as Cookie from '../cookie'
import setTitle from "../title"

export default function Home() {
  const [jots, setJots] = useState(null)

  function jotsPresent() {
    return jots && jots.length !== 0
  }

  function clearJots(jots) {
    const cookieJotIds = Cookie.getJotIds()
    const existingJotIds = new Set(jots.map((jot) => jot.id))

    const jotIdsToRemove = cookieJotIds.filter((cookieJotId) => !existingJotIds.has(cookieJotId))
    Cookie.removeJotIds(...jotIdsToRemove)
  }

  useEffect(() => {
    setTitle('home')

    Api.bulkGetJots(...Cookie.getJotIds())
      .then((jots) => {
        setJots(jots)
        clearJots(jots)
      })
  }, [])

  return <div className="home-page">
    <div className="home">
      <div className="welcome">
        <strong>👋 <em>welcome to jottr!</em></strong>
        <Link className="button" to="/">new jot</Link>
      </div>

      {
        jotsPresent() && <div className="recent-jots">
          <strong>recent jots</strong>
          {jots.map((jot) => <Link className="link" to={`/jot/${jot.id}`} key={jot.id}>{jot.title}</Link>)}
        </div>
      }
    </div>
  </div>
}
