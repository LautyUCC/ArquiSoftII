import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Search from './pages/Search';
import PropertyDetail from './pages/PropertyDetail';
import Congrats from './pages/Congrats';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace />} />
        <Route path="/login" element={<Login />} />
        <Route 
          path="/search" 
          element={
            <ProtectedRoute>
              <Search />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/property/:id" 
          element={
            <ProtectedRoute>
              <PropertyDetail />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/congrats" 
          element={
            <ProtectedRoute>
              <Congrats />
            </ProtectedRoute>
          } 
        />
      </Routes>
    </Router>
  );
}

export default App;
