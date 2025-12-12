import { Navigate } from 'react-router-dom';

/**
 * AdminRoute - Componente que protege rutas solo para ADMIN
 * 
 * Verifica si hay un token en localStorage Y que el usuario tenga rol ADMIN.
 * Si no hay token, redirige a login.
 * Si hay token pero no es ADMIN, redirige a /search.
 * 
 * @param {React.Component} children - El componente a renderizar si el usuario es ADMIN
 * @returns {React.Component} - El componente hijo o un Navigate
 */
function AdminRoute({ children }) {
  // Verificar si hay un token en localStorage
  const token = localStorage.getItem('token');
  
  // Si no hay token, redirigir a login
  if (!token) {
    return <Navigate to="/login" replace />;
  }

  // Verificar el rol del usuario
  const userStr = localStorage.getItem('user');
  if (!userStr) {
    return <Navigate to="/login" replace />;
  }

  try {
    const user = JSON.parse(userStr);
    
    // Comparación robusta: convertir a mayúsculas para evitar problemas de case
    // Si el usuario no es ADMIN, redirigir a /search
    if (user.role?.toUpperCase() !== 'ADMIN') {
      return <Navigate to="/search" replace />;
    }
  } catch (err) {
    console.error('Error parseando usuario:', err);
    return <Navigate to="/login" replace />;
  }

  // Si hay token y es ADMIN, renderizar el componente hijo
  return children;
}

export default AdminRoute;

