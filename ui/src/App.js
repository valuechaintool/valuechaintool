import React from "react";
import Login from './Login';
import Navbar from 'react-bootstrap/Navbar';
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useRouteMatch,
  useParams
} from "react-router-dom";

import 'bootstrap/dist/css/bootstrap.min.css';

export default function App() {
    return (
        <div>
            <Navbar className="d-flex flex-column flex-md-row align-items-center p-3 px-md-4 mb-3 border-bottom box-shadow">
                <Navbar.Brand href="/">ValueChainTool</Navbar.Brand>
            </Navbar>

            <Router>
                <Switch>
                    <Route path={`/topics/:topicId`}>
                        <Topic />
                    </Route>
                    <Route path="/topics">
                        <Topics />
                    </Route>
                    <Route path="/login">
                        <Login />
                    </Route>
                    <Route path="/">
                        <Home />
                    </Route>
                </Switch>
            </Router>
        </div>
    );
}

function Home() {
  return <h2>Home</h2>;
}

function Topics() {
  let match = useRouteMatch();

  return (
    <div>
      <h2>Topics</h2>

      <ul>
        <li>
          <Link to={`${match.url}/components`}>Components</Link>
        </li>
        <li>
          <Link to={`${match.url}/props-v-state`}>
            Props v. State
          </Link>
        </li>
      </ul>

      {/* The Topics page has its own <Switch> with more routes
          that build on the /topics URL path. You can think of the
          2nd <Route> here as an "index" page for all topics, or
          the page that is shown when no topic is selected */}
      <Switch>
        <Route path={match.path}>
          <h3>Please select a topic.</h3>
        </Route>
      </Switch>
    </div>
  );
}

function Topic() {
  let { topicId } = useParams();
  return <h3>Requested topic ID: {topicId}</h3>;
}
