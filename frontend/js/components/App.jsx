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
      <Route exact path="/" component={New}/>
      <Route exact path="/home" component={Home}/>
      <Route path="/:jotId" component={Jot}/>

      <Route>
        <Redirect to="/home"/>
      </Route>
    </Switch>
  </Router>
}
