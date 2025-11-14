import { Navigate } from 'react-router-dom';

/**
 * ProtectedRoute - Componente que protege rutas privadas
 * 
 * Verifica si hay un token en localStorage. Si no hay token,
 * redirige al usuario a la página de login.
 * 
 * @param {React.Component} children - El componente a renderizar si el usuario está autenticado
 * @returns {React.Component} - El componente hijo o un Navigate a /login
 */
function ProtectedRoute({ children }) {
  // Verificar si hay un token en localStorage
  const token = localStorage.getItem('token');

  // Si no hay token, redirigir a login
  if (!token) {
    return <Navigate to="/login" replace />;
  }

  // Si hay token, renderizar el componente hijo
  return children;
}

export default ProtectedRoute;

