import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./auth/AuthContext";
import ProtectedRoute from "./auth/ProtectedRoute";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Profile from "./pages/Profile";
import Settings from "./pages/Settings";
import Vocab from "./pages/Vocab";
import Analysis from "./pages/Analysis";

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="App">
          <header className="App-header"></header>
          <Routes>
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <Home />
                </ProtectedRoute>
              }
            />
            <Route
              path="/vocab"
              element={
                <ProtectedRoute>
                  <Vocab />
                </ProtectedRoute>
              }
            />
            <Route
              path="/vocabulary"
              element={
                <ProtectedRoute>
                  <Vocab />
                </ProtectedRoute>
              }
            />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/profile" element={<Profile />} />
            <Route path="/settings" element={<Settings />}/>
            <Route path="/analysis" element={<Analysis />}/>
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;