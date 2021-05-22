import "bootstrap/dist/css/bootstrap.css";
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import ListsPage from "./pages/Lists/Lists";

const App = () => {
  return (
    <Router>
      <div className="App"></div>
      <Switch>
        <Route path="/">
          <ListsPage />
        </Route>
      </Switch>
    </Router>
  );
};

export default App;
