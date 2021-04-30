import React, {useEffect} from 'react'
import {useHistory} from 'react-router-dom'
import * as Api from '../api'

export default function New() {
  const history = useHistory()

  useEffect(() => {
    Api.createJot()
      .then(({ id }) => history.replace(`/${id}`))
  }, [])

  return null
}
