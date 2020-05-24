import React from 'react';
import './App.css';


import {
  BrowserRouter as Router,
  Route,
  Switch
} from "react-router-dom";

import Authentication from './pages/authentication/authentication';
import Home from './pages/home/home';
import FourOFour from './pages/four-o-four/four-o-four';


function App() {
  return (
     <Router>
        <Switch>
          <Route path="/" exact component={Authentication} />
          <Route path="/home/" component={Home} />
          <Route component={FourOFour} />
        </Switch>
      </Router>
  );
}

export default App;
