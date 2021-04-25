import React, {useEffect} from 'react'
import {useHistory} from 'react-router-dom'
import * as api from '../api'

export default function New() {
  const history = useHistory()

  useEffect(() => {
    api.createJot().then(({ id }) => history.replace(`/jot/${id}`))
  }, [])

  return null
}
