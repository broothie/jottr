import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect
} from 'react-router-dom'
import Home from "./Home";
import New from "./New";
import Jot from "./Jot"

export default function App() {
  return <Router>
    <Switch>
      <Route path="/home">
        <Home/>
      </Route>

      <Route path="/jot/:jotId">
        <Jot/>
      </Route>

      <Route exact path="/">
        <New/>
      </Route>

      <Route>
        <Redirect to="/home"/>
      </Route>
    </Switch>
  </Router>
}
