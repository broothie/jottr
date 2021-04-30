import React from 'react'
import ReactDOM from 'react-dom'
import App from './components/App'

window.onfocus = () => fetch('/api/ping')

document.addEventListener('DOMContentLoaded', () => {
  const root = document.getElementById('root')
  ReactDOM.render(<App/>, root)
})
