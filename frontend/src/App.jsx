import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Search from './pages/Search';
import PropertyDetail from './pages/PropertyDetail';
import Congrats from './pages/Congrats';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/login" replace />} />
        <Route path="/login" element={<Login />} />
        <Route path="/search" element={<Search />} />
        <Route path="/property/:id" element={<PropertyDetail />} />
        <Route path="/congrats" element={<Congrats />} />
      </Routes>
    </Router>
  );
}

export default App;
