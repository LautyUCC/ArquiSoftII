import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authAPI } from '../services/api';

function Register() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    firstName: '',
    lastName: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState({
    username: '',
    email: '',
    password: '',
    firstName: '',
    lastName: ''
  });

  // Validar formato de email
  const isValidEmail = (email) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  // Validar campo username
  const validateUsername = (value) => {
    if (!value.trim()) {
      return 'Por favor ingresa un nombre de usuario';
    }
    if (value.length < 3) {
      return 'El nombre de usuario debe tener al menos 3 caracteres';
    }
    if (value.length > 50) {
      return 'El nombre de usuario no puede tener más de 50 caracteres';
    }
    return '';
  };

  // Validar campo email
  const validateEmail = (value) => {
    if (!value.trim()) {
      return 'Por favor ingresa tu email';
    }
    if (!isValidEmail(value)) {
      return 'Por favor ingresa un email válido';
    }
    return '';
  };

  // Validar campo password
  const validatePassword = (value) => {
    if (!value.trim()) {
      return 'Por favor ingresa una contraseña';
    }
    if (value.length < 6) {
      return 'La contraseña debe tener al menos 6 caracteres';
    }
    return '';
  };

  // Validar campo firstName
  const validateFirstName = (value) => {
    if (!value.trim()) {
      return 'Por favor ingresa tu nombre';
    }
    return '';
  };

  // Validar campo lastName
  const validateLastName = (value) => {
    if (!value.trim()) {
      return 'Por favor ingresa tu apellido';
    }
    return '';
  };

  // Manejar cambio en campos
  const handleChange = (field, value) => {
    setFormData({...formData, [field]: value});
    
    // Validar en tiempo real
    let error = '';
    switch (field) {
      case 'username':
        error = validateUsername(value);
        break;
      case 'email':
        error = validateEmail(value);
        break;
      case 'password':
        error = validatePassword(value);
        break;
      case 'firstName':
        error = validateFirstName(value);
        break;
      case 'lastName':
        error = validateLastName(value);
        break;
    }
    
    if (fieldErrors[field] || value.trim()) {
      setFieldErrors({
        ...fieldErrors,
        [field]: error
      });
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    // Validar todos los campos
    const usernameError = validateUsername(formData.username);
    const emailError = validateEmail(formData.email);
    const passwordError = validatePassword(formData.password);
    const firstNameError = validateFirstName(formData.firstName);
    const lastNameError = validateLastName(formData.lastName);

    setFieldErrors({
      username: usernameError,
      email: emailError,
      password: passwordError,
      firstName: firstNameError,
      lastName: lastNameError
    });

    // Si hay errores de validación, no continuar
    if (usernameError || emailError || passwordError || firstNameError || lastNameError) {
      return;
    }

    setLoading(true);

    try {
      // Llamar al endpoint de registro a través de Nginx
      const response = await authAPI.register({
        username: formData.username,
        email: formData.email,
        password: formData.password,
        firstName: formData.firstName,
        lastName: formData.lastName
      });
      
      // Si el registro es exitoso, intentar loguear automáticamente
      if (response.status === 201) {
        try {
          // Intentar login automático con las credenciales
          const loginResponse = await authAPI.login(formData.email, formData.password);
          
          if (loginResponse.data.token) {
            localStorage.setItem('token', loginResponse.data.token);
            localStorage.setItem('user', JSON.stringify(loginResponse.data.user));
            navigate('/search');
          } else {
            // Si el login automático falla, redirigir a login
            navigate('/login');
          }
        } catch (loginError) {
          // Si el login automático falla, redirigir a login
          console.error('Error en login automático:', loginError);
          navigate('/login');
        }
      } else {
        setError('Error al registrar. Por favor intenta de nuevo.');
      }
    } catch (err) {
      console.error('Error en registro:', err);
      
      // Manejar errores de validación del backend
      if (err.response?.status === 400) {
        const errorMessage = err.response?.data?.error || 'Error de validación';
        setError(errorMessage);
        
        // Intentar extraer errores de campos específicos del mensaje
        const errorText = errorMessage.toLowerCase();
        if (errorText.includes('username')) {
          setFieldErrors(prev => ({...prev, username: errorMessage}));
        } else if (errorText.includes('email')) {
          setFieldErrors(prev => ({...prev, email: errorMessage}));
        } else if (errorText.includes('password')) {
          setFieldErrors(prev => ({...prev, password: errorMessage}));
        } else if (errorText.includes('firstname') || errorText.includes('nombre')) {
          setFieldErrors(prev => ({...prev, firstName: errorMessage}));
        } else if (errorText.includes('lastname') || errorText.includes('apellido')) {
          setFieldErrors(prev => ({...prev, lastName: errorMessage}));
        }
      } else if (err.response?.status === 409 || err.response?.status === 422) {
        // Conflict o Unprocessable Entity - usuario o email ya existe
        const errorMessage = err.response?.data?.error || 'El usuario o email ya existe';
        setError(errorMessage);
        
        if (errorMessage.toLowerCase().includes('username')) {
          setFieldErrors(prev => ({...prev, username: errorMessage}));
        } else if (errorMessage.toLowerCase().includes('email')) {
          setFieldErrors(prev => ({...prev, email: errorMessage}));
        }
      } else if (err.code === 'ERR_NETWORK') {
        setError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setError(err.response?.data?.error || 'Error al registrar. Por favor intenta de nuevo.');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-secondary flex items-center justify-center px-4 py-8">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-primary mb-2">Spotly</h1>
          <p className="text-gray-600">Crea tu cuenta</p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-600 text-sm">{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Nombre de Usuario
            </label>
            <input
              type="text"
              value={formData.username}
              onChange={(e) => handleChange('username', e.target.value)}
              onBlur={(e) => {
                setFieldErrors({
                  ...fieldErrors,
                  username: validateUsername(e.target.value)
                });
              }}
              className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition ${
                fieldErrors.username 
                  ? 'border-red-300 focus:ring-red-500' 
                  : 'border-gray-300'
              }`}
              placeholder="usuario123"
              disabled={loading}
            />
            {fieldErrors.username && (
              <p className="mt-1 text-sm text-red-600">{fieldErrors.username}</p>
            )}
            {formData.username.length > 0 && formData.username.length < 3 && (
              <p className="mt-1 text-sm text-gray-500">
                Mínimo 3 caracteres ({formData.username.length}/3)
              </p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              type="email"
              value={formData.email}
              onChange={(e) => handleChange('email', e.target.value)}
              onBlur={(e) => {
                setFieldErrors({
                  ...fieldErrors,
                  email: validateEmail(e.target.value)
                });
              }}
              className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition ${
                fieldErrors.email 
                  ? 'border-red-300 focus:ring-red-500' 
                  : 'border-gray-300'
              }`}
              placeholder="usuario@email.com"
              disabled={loading}
            />
            {fieldErrors.email && (
              <p className="mt-1 text-sm text-red-600">{fieldErrors.email}</p>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Nombre
              </label>
              <input
                type="text"
                value={formData.firstName}
                onChange={(e) => handleChange('firstName', e.target.value)}
                onBlur={(e) => {
                  setFieldErrors({
                    ...fieldErrors,
                    firstName: validateFirstName(e.target.value)
                  });
                }}
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition ${
                  fieldErrors.firstName 
                    ? 'border-red-300 focus:ring-red-500' 
                    : 'border-gray-300'
                }`}
                placeholder="Juan"
                disabled={loading}
              />
              {fieldErrors.firstName && (
                <p className="mt-1 text-sm text-red-600">{fieldErrors.firstName}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Apellido
              </label>
              <input
                type="text"
                value={formData.lastName}
                onChange={(e) => handleChange('lastName', e.target.value)}
                onBlur={(e) => {
                  setFieldErrors({
                    ...fieldErrors,
                    lastName: validateLastName(e.target.value)
                  });
                }}
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent transition ${
                  fieldErrors.lastName 
                    ? 'border-red-300 focus:ring-red-500' 
                    : 'border-gray-300'
                }`}
                placeholder="Pérez"
                disabled={loading}
              />
              {fieldErrors.lastName && (
                <p className="mt-1 text-sm text-red-600">{fieldErrors.lastName}</p>
              )}
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Contraseña
            </label>
            <input
              type="password"
              value={formData.password}
              onChange={(e) => handleChange('password', e.target.value)}
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
            {loading ? 'Registrando...' : 'Registrarse'}
          </button>
        </form>

        <p className="text-center text-sm text-gray-600 mt-6">
          ¿Ya tienes cuenta?{' '}
          <button
            onClick={() => navigate('/login')}
            className="text-primary font-medium hover:underline"
          >
            Inicia sesión
          </button>
        </p>
      </div>
    </div>
  );
}

export default Register;

