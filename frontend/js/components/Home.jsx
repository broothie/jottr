import React, {useEffect, useState} from 'react'
import {Link} from 'react-router-dom'
import * as api from '../api'
import * as Cookie from '../cookie'

export default function Home() {
  const [jots, setJots] = useState(null)

  useEffect(() => {
    api.bulkGetJots(...Cookie.getJotIds()).then(setJots)
  }, [])

  return <div className="home-page">
    <div className="home">
      <div className="welcome">
        <p>👋 welcome to jottr!</p>
        <Link className="button" to="/">new</Link>
      </div>

      {jots && <div className="recent-jots">
        <strong>recent jots</strong>
        {jots.map((jot) => (
          <Link className="link" to={`/jot/${jot.id}`} key={jot.id}>
            {jot.title}
          </Link>
        ))}
      </div>}
    </div>
  </div>
}
