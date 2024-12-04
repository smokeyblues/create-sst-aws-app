import { Link, useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import { Auth } from "aws-amplify";
import { AppContext, AppContextType } from "./lib/contextLib";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import { onError } from "./lib/errorLib";
import Routes from "./Routes.tsx";
import "./App.css";


function App() {
  const nav = useNavigate();
  const [isAuthenticating, setIsAuthenticating] = useState(true);
  const [isAuthenticated, userHasAuthenticated] = useState(false);

  useEffect(() => {
    onLoad();
  }, []);
  
  async function onLoad() {
    try {
      await Auth.currentSession();
      userHasAuthenticated(true);
    } catch (error) {
      if (error !== "No current user") {
        onError(error);
      }
    }
  
    setIsAuthenticating(false);
  }

  async function handleLogout() {
    await Auth.signOut();
  
    userHasAuthenticated(false);

    nav("/login");
  }

  return (
    !isAuthenticating && (
      <div className="App container py-3">
        <Navbar collapseOnSelect bg="light" expand="md" className="mb-3 px-3">
          <Nav.Link as={Link} to="/">
            <Navbar.Brand className="fw-bold text-muted">Scratch</Navbar.Brand>
          </Nav.Link>
          <Navbar.Toggle />
          <Navbar.Collapse className="justify-content-end">
            <Nav activeKey={window.location.pathname}>
              {isAuthenticated ? (
              <>
                <Nav.Link as={Link} to="/settings">Settings</Nav.Link>
                <Nav.Link onClick={handleLogout}>Logout</Nav.Link>
              </>
                ) : (
                  <>
                    <Nav.Link as={Link} to="signup">Signup</Nav.Link>
                    <Nav.Link as={Link} to="/login">Login</Nav.Link>
                  </>
                )}
                
                
            </Nav>
          </Navbar.Collapse>
        </Navbar>
        <AppContext.Provider
          value={{ isAuthenticated, userHasAuthenticated } as AppContextType}
        >
          <Routes />
        </AppContext.Provider>
      </div>
    )
  );
}

export default App;