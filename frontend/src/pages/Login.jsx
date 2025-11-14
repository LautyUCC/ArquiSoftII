import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authAPI } from '../services/api';

function Login() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    usernameOrEmail: '',
    password: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState({
    usernameOrEmail: '',
    password: ''
  });

  // Validar formato de email
  const isValidEmail = (email) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  // Validar campo username/email
  const validateUsernameOrEmail = (value) => {
    if (!value.trim()) {
      return 'Por favor ingresa tu usuario o email';
    }
    // Si contiene @, validar formato de email
    if (value.includes('@') && !isValidEmail(value)) {
      return 'Por favor ingresa un email válido';
    }
    return '';
  };

  // Validar campo password
  const validatePassword = (value) => {
    if (!value.trim()) {
      return 'Por favor ingresa tu contraseña';
    }
    if (value.length < 6) {
      return 'La contraseña debe tener al menos 6 caracteres';
    }
    return '';
  };

  // Manejar cambio en username/email
  const handleUsernameOrEmailChange = (e) => {
    const value = e.target.value;
    setFormData({...formData, usernameOrEmail: value});
    // Validar en tiempo real
    if (fieldErrors.usernameOrEmail || value.trim()) {
      setFieldErrors({
        ...fieldErrors,
        usernameOrEmail: validateUsernameOrEmail(value)
      });
    }
  };

  // Manejar cambio en password
  const handlePasswordChange = (e) => {
    const value = e.target.value;
    setFormData({...formData, password: value});
    // Validar en tiempo real
    if (fieldErrors.password || value.trim()) {
      setFieldErrors({
        ...fieldErrors,
        password: validatePassword(value)
      });
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    // Validar todos los campos
    const usernameError = validateUsernameOrEmail(formData.usernameOrEmail);
    const passwordError = validatePassword(formData.password);

    setFieldErrors({
      usernameOrEmail: usernameError,
      password: passwordError
    });

    // Si hay errores de validación, no continuar
    if (usernameError || passwordError) {
      return;
    }

    setLoading(true);

    try {
      const response = await authAPI.login(formData.usernameOrEmail, formData.password);
      
      if (response.data.token) {
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('user', JSON.stringify(response.data.user));
        navigate('/search');
      } else {
        setError('Error al iniciar sesión. Por favor intenta de nuevo.');
      }
    } catch (err) {
      console.error('Error en login:', err);
      
      if (err.response?.status === 401) {
        setError('Usuario o contraseña incorrectos');
      } else if (err.response?.status === 404) {
        setError('Usuario no encontrado');
      } else if (err.code === 'ERR_NETWORK') {
        setError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setError(err.response?.data?.message || 'Error al iniciar sesión. Por favor intenta de nuevo.');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-secondary flex items-center justify-center px-4">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-primary mb-2">Spotly</h1>
          <p className="text-gray-600">Propiedades de lujo</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-600 text-sm">{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Usuario o Email
            </label>
            <input
              type="text"
              value={formData.usernameOrEmail}
              onChange={handleUsernameOrEmailChange}
              onBlur={(e) => {
                setFieldErrors({
                  ...fieldErrors,
                  usernameOrEmail: validateUsernameOrEmail(e.target.value)
                });
              }}
              className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition ${
                fieldErrors.usernameOrEmail 
                  ? 'border-red-300 focus:ring-red-500' 
                  : 'border-gray-300'
              }`}
              placeholder="admin o usuario@email.com"
              disabled={loading}
            />
            {fieldErrors.usernameOrEmail && (
              <p className="mt-1 text-sm text-red-600">{fieldErrors.usernameOrEmail}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Contraseña
            </label>
            <input
              type="password"
              value={formData.password}
              onChange={handlePasswordChange}
              onBlur={(e) => {
                setFieldErrors({
                  ...fieldErrors,
                  password: validatePassword(e.target.value)
                });
              }}
              className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition ${
                fieldErrors.password 
                  ? 'border-red-300 focus:ring-red-500' 
                  : 'border-gray-300'
              }`}
              placeholder="••••••••"
              disabled={loading}
            />
            {fieldErrors.password && (
              <p className="mt-1 text-sm text-red-600">{fieldErrors.password}</p>
            )}
            {formData.password.length > 0 && formData.password.length < 6 && (
              <p className="mt-1 text-sm text-gray-500">
                Mínimo 6 caracteres ({formData.password.length}/6)
              </p>
            )}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-primary text-white py-3 rounded-lg font-medium hover:bg-gray-800 transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? 'Iniciando sesión...' : 'Iniciar Sesión'}
          </button>
        </form>

        <p className="text-center text-sm text-gray-600 mt-6">
          ¿No tienes cuenta?{' '}
          <button
            onClick={() => navigate('/register')}
            className="text-primary font-medium hover:underline"
          >
            Regístrate
          </button>
        </p>
      </div>
    </div>
  );
}

export default Login;
