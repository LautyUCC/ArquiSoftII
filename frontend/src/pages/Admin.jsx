import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { adminAPI, propertiesAPI, bookingsAPI, uploadAPI } from '../services/api';
import { Users, Home, LogOut, Edit, Trash2, Plus, X, Check, AlertCircle, Calendar } from 'lucide-react';

function Admin() {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState('users'); // 'users', 'properties' o 'bookings'
  
  // Estados para Usuarios
  const [users, setUsers] = useState([]);
  const [usersLoading, setUsersLoading] = useState(false);
  const [usersError, setUsersError] = useState('');
  const [usersSuccess, setUsersSuccess] = useState('');
  const [editingUser, setEditingUser] = useState(null);
  const [userRoleChange, setUserRoleChange] = useState({});

  // Estados para Propiedades
  const [properties, setProperties] = useState([]);
  const [propertiesLoading, setPropertiesLoading] = useState(false);
  const [propertiesError, setPropertiesError] = useState('');
  const [propertiesSuccess, setPropertiesSuccess] = useState('');
  const [showCreateProperty, setShowCreateProperty] = useState(false);
  const [editingProperty, setEditingProperty] = useState(null);
  const [deleteConfirm, setDeleteConfirm] = useState(null);

  // Formulario de nueva propiedad
  const [newProperty, setNewProperty] = useState({
    title: '',
    description: '',
    price: '',
    location: '',
    ownerId: '',
    capacity: '',
    available: true,
    amenities: [],
    images: []
  });

  // Formulario de edición de propiedad
  const [editPropertyData, setEditPropertyData] = useState({});

  // Estados para Reservas
  const [bookings, setBookings] = useState([]);
  const [bookingsLoading, setBookingsLoading] = useState(false);
  const [bookingsError, setBookingsError] = useState('');
  const [bookingsSuccess, setBookingsSuccess] = useState('');
  const [propertyDetails, setPropertyDetails] = useState({}); // Cache de detalles de propiedades
  const [cancelLoading, setCancelLoading] = useState({});

  useEffect(() => {
    if (activeTab === 'users') {
      loadUsers();
    } else if (activeTab === 'properties') {
      loadProperties();
    } else if (activeTab === 'bookings') {
      loadBookings();
    }
  }, [activeTab]);

  // ============================================
  // USUARIOS
  // ============================================

  const loadUsers = async () => {
    setUsersLoading(true);
    setUsersError('');
    setUsersSuccess('');

    try {
      const response = await adminAPI.getAllUsers();
      setUsers(response.data || []);
    } catch (err) {
      console.error('Error al cargar usuarios:', err);
      
      if (err.response?.status === 401) {
        setUsersError('Sesión expirada. Por favor inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 403) {
        setUsersError('No tienes permisos para ver usuarios. Se requiere rol ADMIN.');
      } else if (err.code === 'ERR_NETWORK') {
        setUsersError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setUsersError('Error al cargar usuarios. Por favor intenta de nuevo.');
      }
    } finally {
      setUsersLoading(false);
    }
  };

  const handleRoleChange = async (userId, newRole) => {
    setUsersError('');
    setUsersSuccess('');

    try {
      await adminAPI.updateUser(userId, { role: newRole });
      setUsersSuccess(`Rol del usuario actualizado a ${newRole} exitosamente`);
      setUserRoleChange({ ...userRoleChange, [userId]: newRole });
      
      // Ocultar mensaje de éxito después de 3 segundos
      setTimeout(() => {
        setUsersSuccess('');
      }, 3000);
      
      // Recargar usuarios después de un breve delay
      setTimeout(() => {
        loadUsers();
      }, 1000);
    } catch (err) {
      console.error('Error al actualizar rol:', err);
      
      if (err.response?.status === 403) {
        setUsersError('No tienes permisos para cambiar roles. Se requiere rol ADMIN.');
      } else if (err.response?.status === 400) {
        setUsersError(err.response?.data?.error || 'Error al actualizar el rol del usuario');
      } else {
        setUsersError('Error al actualizar el rol. Por favor intenta de nuevo.');
      }
    }
  };

  // ============================================
  // PROPIEDADES
  // ============================================

  const loadProperties = async () => {
    setPropertiesLoading(true);
    setPropertiesError('');
    setPropertiesSuccess('');

    try {
      const response = await propertiesAPI.getAll();
      setProperties(response.data || []);
    } catch (err) {
      console.error('Error al cargar propiedades:', err);
      
      if (err.response?.status === 401) {
        setPropertiesError('Sesión expirada. Por favor inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 403) {
        setPropertiesError('No tienes permisos para ver propiedades. Se requiere rol ADMIN.');
      } else if (err.code === 'ERR_NETWORK') {
        setPropertiesError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setPropertiesError('Error al cargar propiedades. Por favor intenta de nuevo.');
      }
    } finally {
      setPropertiesLoading(false);
    }
  };

  const handleCreateProperty = async (e) => {
    e.preventDefault();
    setPropertiesError('');
    setPropertiesSuccess('');

    // Validaciones
    if (!newProperty.title || !newProperty.description || !newProperty.price || 
        !newProperty.location || !newProperty.ownerId || !newProperty.capacity) {
      setPropertiesError('Por favor completa todos los campos requeridos');
      return;
    }

    const propertyData = {
      title: newProperty.title,
      description: newProperty.description,
      price: parseFloat(newProperty.price),
      location: newProperty.location,
      ownerId: newProperty.ownerId,
      capacity: parseInt(newProperty.capacity),
      available: newProperty.available,
      amenities: newProperty.amenities,
      images: newProperty.images
    };

    try {
      await propertiesAPI.create(propertyData);
      setPropertiesSuccess('Propiedad creada exitosamente');
      
      // Ocultar mensaje de éxito después de 3 segundos
      setTimeout(() => {
        setPropertiesSuccess('');
      }, 3000);
      
      setShowCreateProperty(false);
      setNewProperty({
        title: '',
        description: '',
        price: '',
        location: '',
        ownerId: '',
        capacity: '',
        available: true,
        amenities: [],
        images: []
      });
      loadProperties();
    } catch (err) {
      console.error('Error al crear propiedad:', err);
      
      if (err.response?.status === 403) {
        setPropertiesError('No tienes permisos para crear propiedades. Se requiere rol ADMIN.');
      } else if (err.response?.status === 400) {
        setPropertiesError(err.response?.data?.error || 'Error al crear la propiedad');
      } else {
        setPropertiesError('Error al crear la propiedad. Por favor intenta de nuevo.');
      }
    }
  };

  const handleUpdateProperty = async (propertyId) => {
    setPropertiesError('');
    setPropertiesSuccess('');

    try {
      await propertiesAPI.update(propertyId, editPropertyData);
      setPropertiesSuccess('Propiedad actualizada exitosamente');
      
      // Ocultar mensaje de éxito después de 3 segundos
      setTimeout(() => {
        setPropertiesSuccess('');
      }, 3000);
      
      setEditingProperty(null);
      setEditPropertyData({});
      loadProperties();
    } catch (err) {
      console.error('Error al actualizar propiedad:', err);
      
      if (err.response?.status === 403) {
        setPropertiesError('No tienes permisos para editar propiedades. Se requiere rol ADMIN.');
      } else if (err.response?.status === 400) {
        setPropertiesError(err.response?.data?.error || 'Error al actualizar la propiedad');
      } else {
        setPropertiesError('Error al actualizar la propiedad. Por favor intenta de nuevo.');
      }
    }
  };

  const handleDeleteProperty = async (propertyId) => {
    setPropertiesError('');
    setPropertiesSuccess('');

    try {
      await propertiesAPI.delete(propertyId);
      setPropertiesSuccess('Propiedad eliminada exitosamente');
      
      // Ocultar mensaje de éxito después de 3 segundos
      setTimeout(() => {
        setPropertiesSuccess('');
      }, 3000);
      
      setDeleteConfirm(null);
      loadProperties();
    } catch (err) {
      console.error('Error al eliminar propiedad:', err);
      
      if (err.response?.status === 403) {
        setPropertiesError('No tienes permisos para eliminar propiedades. Se requiere rol ADMIN.');
      } else if (err.response?.status === 404) {
        setPropertiesError('Propiedad no encontrada');
      } else {
        setPropertiesError('Error al eliminar la propiedad. Por favor intenta de nuevo.');
      }
    }
  };

  // ============================================
  // RESERVAS
  // ============================================

  const loadBookings = async () => {
    setBookingsLoading(true);
    setBookingsError('');
    setBookingsSuccess('');

    try {
      const response = await bookingsAPI.getAllBookings();
      const bookingsData = Array.isArray(response.data) ? response.data : (response.data?.data || []);
      setBookings(bookingsData);
      
      // Cargar detalles de propiedades para las reservas
      await loadPropertyDetailsForBookings(bookingsData);
      
      if (bookingsData.length === 0) {
        setBookingsError('No hay reservas registradas');
      }
    } catch (err) {
      console.error('Error al cargar reservas:', err);
      
      if (err.response?.status === 401) {
        setBookingsError('Sesión expirada. Por favor inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 403) {
        setBookingsError('No tienes permisos para ver reservas. Se requiere rol ADMIN.');
      } else if (err.code === 'ERR_NETWORK') {
        setBookingsError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setBookingsError('Error al cargar las reservas. Por favor intenta de nuevo.');
      }
    } finally {
      setBookingsLoading(false);
    }
  };

  const loadPropertyDetailsForBookings = async (bookingsData) => {
    const propertyIds = [...new Set(bookingsData.map(b => b.propertyId))];
    const details = { ...propertyDetails };
    
    // Cargar detalles solo de propiedades que no tenemos en cache
    for (const propertyId of propertyIds) {
      if (!details[propertyId]) {
        try {
          const response = await propertiesAPI.getById(propertyId);
          details[propertyId] = response.data;
        } catch (err) {
          console.error(`Error al cargar propiedad ${propertyId}:`, err);
          details[propertyId] = { title: `Propiedad ${propertyId}`, location: 'N/A' };
        }
      }
    }
    
    setPropertyDetails(details);
  };

  const handleCancelBooking = async (bookingId) => {
    if (!window.confirm('¿Estás seguro de que deseas cancelar esta reserva?')) {
      return;
    }

    setCancelLoading(prev => ({ ...prev, [bookingId]: true }));
    setBookingsError('');
    setBookingsSuccess('');

    try {
      await bookingsAPI.deleteBooking(bookingId);
      
      setBookingsSuccess('Reserva cancelada exitosamente');
      
      // Actualizar la lista localmente marcando como CANCELLED
      setBookings(prev => prev.map(b => 
        b.id === bookingId ? { ...b, status: 'CANCELLED' } : b
      ));
      
      setTimeout(() => {
        setBookingsSuccess('');
        loadBookings();
      }, 1500);
    } catch (err) {
      console.error('Error al cancelar reserva:', err);
      
      if (err.response?.status === 401 || err.response?.status === 403) {
        setBookingsError('Tu sesión ha expirado o no tienes permisos. Por favor, inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 404) {
        setBookingsError('Reserva no encontrada');
      } else {
        setBookingsError('Error al cancelar la reserva. Por favor, intenta de nuevo.');
      }
    } finally {
      setCancelLoading(prev => ({ ...prev, [bookingId]: false }));
    }
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const getStatusBadge = (status) => {
    const statusUpper = status?.toUpperCase();
    const statusColors = {
      CONFIRMED: 'bg-green-100 text-green-800',
      PENDING: 'bg-yellow-100 text-yellow-800',
      CANCELLED: 'bg-red-100 text-red-800'
    };

    const statusLabels = {
      CONFIRMED: 'Confirmada',
      PENDING: 'Pendiente',
      CANCELLED: 'Cancelada'
    };

    return (
      <span className={`px-3 py-1 rounded-full text-xs font-medium ${statusColors[statusUpper] || 'bg-gray-100 text-gray-800'}`}>
        {statusLabels[statusUpper] || status}
      </span>
    );
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate('/');
  };

  const getRoleBadge = (role) => {
    const roleColors = {
      ADMIN: 'bg-purple-100 text-purple-800',
      USER: 'bg-blue-100 text-blue-800',
      OWNER: 'bg-green-100 text-green-800'
    };

    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${roleColors[role] || 'bg-gray-100 text-gray-800'}`}>
        {role}
      </span>
    );
  };

  return (
    <div className="min-h-screen bg-secondary">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <h1 className="text-2xl font-bold text-primary">Spotly - Administración</h1>
            <button
              onClick={() => navigate('/search')}
              className="flex items-center gap-2 text-gray-600 hover:text-primary transition"
            >
              <Home size={20} />
              Inicio
            </button>
          </div>
          <button
            onClick={handleLogout}
            className="flex items-center gap-2 text-gray-600 hover:text-primary transition"
          >
            <LogOut size={20} />
            Salir
          </button>
        </div>
      </header>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Tabs */}
        <div className="bg-white rounded-2xl shadow-lg mb-6">
          <div className="flex border-b border-gray-200">
            <button
              onClick={() => {
                setActiveTab('users');
                setUsersError('');
                setUsersSuccess('');
              }}
              className={`flex-1 px-6 py-4 font-medium transition ${
                activeTab === 'users'
                  ? 'text-primary border-b-2 border-primary'
                  : 'text-gray-600 hover:text-primary'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <Users size={20} />
                Usuarios
              </div>
            </button>
            <button
              onClick={() => {
                setActiveTab('properties');
                setPropertiesError('');
                setPropertiesSuccess('');
              }}
              className={`flex-1 px-6 py-4 font-medium transition ${
                activeTab === 'properties'
                  ? 'text-primary border-b-2 border-primary'
                  : 'text-gray-600 hover:text-primary'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <Home size={20} />
                Propiedades
              </div>
            </button>
            <button
              onClick={() => {
                setActiveTab('bookings');
                setBookingsError('');
                setBookingsSuccess('');
              }}
              className={`flex-1 px-6 py-4 font-medium transition ${
                activeTab === 'bookings'
                  ? 'text-primary border-b-2 border-primary'
                  : 'text-gray-600 hover:text-primary'
              }`}
            >
              <div className="flex items-center justify-center gap-2">
                <Calendar size={20} />
                Reservas
              </div>
            </button>
          </div>
        </div>

        {/* Sección Usuarios */}
        {activeTab === 'users' && (
          <div className="bg-white rounded-2xl shadow-lg p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold text-primary">Gestión de Usuarios</h2>
            </div>

            {/* Mensajes */}
            {usersError && (
              <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg flex items-center gap-2">
                <AlertCircle size={20} className="text-red-600" />
                <p className="text-red-600">{usersError}</p>
              </div>
            )}

            {usersSuccess && (
              <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg flex items-center gap-2">
                <Check size={20} className="text-green-600" />
                <p className="text-green-600">{usersSuccess}</p>
              </div>
            )}

            {/* Loading */}
            {usersLoading && (
              <div className="text-center py-12">
                <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                <p className="mt-4 text-gray-600">Cargando usuarios...</p>
              </div>
            )}

            {/* Tabla de Usuarios */}
            {!usersLoading && users.length > 0 && (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">ID</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Email</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Rol Actual</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Cambiar Rol</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map((user) => (
                      <tr key={user.id} className="border-b border-gray-100 hover:bg-gray-50 transition">
                        <td className="py-4 px-4">{user.id}</td>
                        <td className="py-4 px-4 font-medium">{user.email}</td>
                        <td className="py-4 px-4">
                          {getRoleBadge(userRoleChange[user.id] || user.role)}
                        </td>
                        <td className="py-4 px-4">
                          <select
                            value={userRoleChange[user.id] || user.role}
                            onChange={(e) => handleRoleChange(user.id, e.target.value)}
                            className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                            disabled={user.role === 'ADMIN' && userRoleChange[user.id] !== 'USER'}
                          >
                            <option value="USER">USER</option>
                            <option value="ADMIN">ADMIN</option>
                            {user.role === 'OWNER' && <option value="OWNER">OWNER</option>}
                          </select>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            {/* Empty State */}
            {!usersLoading && users.length === 0 && !usersError && (
              <div className="text-center py-12">
                <Users size={64} className="mx-auto text-gray-300 mb-4" />
                <p className="text-xl text-gray-600">No hay usuarios registrados</p>
              </div>
            )}
          </div>
        )}

        {/* Sección Propiedades */}
        {activeTab === 'properties' && (
          <div className="bg-white rounded-2xl shadow-lg p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold text-primary">Gestión de Propiedades</h2>
              <button
                onClick={() => setShowCreateProperty(true)}
                className="flex items-center gap-2 bg-primary text-white px-4 py-2 rounded-lg hover:bg-gray-800 transition"
              >
                <Plus size={20} />
                Crear Propiedad
              </button>
            </div>

            {/* Mensajes */}
            {propertiesError && (
              <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg flex items-center gap-2">
                <AlertCircle size={20} className="text-red-600" />
                <p className="text-red-600">{propertiesError}</p>
              </div>
            )}

            {propertiesSuccess && (
              <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg flex items-center gap-2">
                <Check size={20} className="text-green-600" />
                <p className="text-green-600">{propertiesSuccess}</p>
              </div>
            )}

            {/* Modal Crear Propiedad */}
            {showCreateProperty && (
              <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
                <div className="bg-white rounded-2xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                  <div className="p-6 border-b border-gray-200 flex items-center justify-between">
                    <h3 className="text-xl font-bold text-primary">Crear Nueva Propiedad</h3>
                    <button
                      onClick={() => {
                        setShowCreateProperty(false);
                        setNewProperty({
                          title: '',
                          description: '',
                          price: '',
                          location: '',
                          ownerId: '',
                          capacity: '',
                          available: true,
                          amenities: [],
                          images: []
                        });
                      }}
                      className="text-gray-400 hover:text-gray-600"
                    >
                      <X size={24} />
                    </button>
                  </div>
                  <form onSubmit={handleCreateProperty} className="p-6 space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Título *</label>
                      <input
                        type="text"
                        value={newProperty.title}
                        onChange={(e) => setNewProperty({...newProperty, title: e.target.value})}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Descripción *</label>
                      <textarea
                        value={newProperty.description}
                        onChange={(e) => setNewProperty({...newProperty, description: e.target.value})}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        rows="3"
                        required
                      />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Precio *</label>
                        <input
                          type="number"
                          step="0.01"
                          value={newProperty.price}
                          onChange={(e) => setNewProperty({...newProperty, price: e.target.value})}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                          required
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Capacidad *</label>
                        <input
                          type="number"
                          value={newProperty.capacity}
                          onChange={(e) => setNewProperty({...newProperty, capacity: e.target.value})}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                          required
                        />
                      </div>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Ubicación *</label>
                      <input
                        type="text"
                        value={newProperty.location}
                        onChange={(e) => setNewProperty({...newProperty, location: e.target.value})}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Owner ID *</label>
                      <input
                        type="text"
                        value={newProperty.ownerId}
                        onChange={(e) => setNewProperty({...newProperty, ownerId: e.target.value})}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        required
                      />
                    </div>
                    <div>
                      <label className="flex items-center gap-2">
                        <input
                          type="checkbox"
                          checked={newProperty.available}
                          onChange={(e) => setNewProperty({...newProperty, available: e.target.checked})}
                          className="w-4 h-4 text-primary"
                        />
                        <span className="text-sm text-gray-700">Disponible</span>
                      </label>
                    </div>
                    <div className="flex gap-4">
                      <button
                        type="submit"
                        className="flex-1 bg-primary text-white py-2 rounded-lg hover:bg-gray-800 transition"
                      >
                        Crear Propiedad
                      </button>
                      <button
                        type="button"
                        onClick={() => {
                          setShowCreateProperty(false);
                          setNewProperty({
                            title: '',
                            description: '',
                            price: '',
                            location: '',
                            ownerId: '',
                            capacity: '',
                            available: true,
                            amenities: [],
                            images: []
                          });
                        }}
                        className="flex-1 bg-gray-200 text-gray-700 py-2 rounded-lg hover:bg-gray-300 transition"
                      >
                        Cancelar
                      </button>
                    </div>
                  </form>
                </div>
              </div>
            )}

            {/* Modal Editar Propiedad */}
            {editingProperty && (
              <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
                <div className="bg-white rounded-2xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
                  <div className="p-6 border-b border-gray-200 flex items-center justify-between">
                    <h3 className="text-xl font-bold text-primary">Editar Propiedad</h3>
                    <button
                      onClick={() => {
                        setEditingProperty(null);
                        setEditPropertyData({});
                      }}
                      className="text-gray-400 hover:text-gray-600"
                    >
                      <X size={24} />
                    </button>
                  </div>
                  <form
                    onSubmit={(e) => {
                      e.preventDefault();
                      handleUpdateProperty(editingProperty.id);
                    }}
                    className="p-6 space-y-4"
                  >
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Título</label>
                      <input
                        type="text"
                        defaultValue={editingProperty.title}
                        onChange={(e) => setEditPropertyData({...editPropertyData, title: e.target.value})}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Descripción</label>
                      <textarea
                        defaultValue={editingProperty.description}
                        onChange={(e) => setEditPropertyData({...editPropertyData, description: e.target.value})}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        rows="3"
                      />
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Precio</label>
                        <input
                          type="number"
                          step="0.01"
                          defaultValue={editingProperty.price}
                          onChange={(e) => setEditPropertyData({...editPropertyData, price: parseFloat(e.target.value)})}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        />
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Capacidad</label>
                        <input
                          type="number"
                          defaultValue={editingProperty.capacity}
                          onChange={(e) => setEditPropertyData({...editPropertyData, capacity: parseInt(e.target.value)})}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        />
                      </div>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Ubicación</label>
                      <input
                        type="text"
                        defaultValue={editingProperty.location}
                        onChange={(e) => setEditPropertyData({...editPropertyData, location: e.target.value})}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="flex items-center gap-2">
                        <input
                          type="checkbox"
                          defaultChecked={editingProperty.available}
                          onChange={(e) => setEditPropertyData({...editPropertyData, available: e.target.checked})}
                          className="w-4 h-4 text-primary"
                        />
                        <span className="text-sm text-gray-700">Disponible</span>
                      </label>
                    </div>

                    {/* Imágenes */}
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Imágenes <span className="text-gray-500 text-xs">(JPG, PNG, GIF, WEBP, máx. 5MB)</span>
                      </label>
                      <div className="mb-2">
                        <input
                          type="file"
                          accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
                          onChange={handleEditImageUpload}
                          disabled={propertiesLoading}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed"
                        />
                        <p className="mt-1 text-xs text-gray-500">
                          Selecciona una imagen para agregarla a la propiedad.
                        </p>
                      </div>
                      {(editPropertyData.images || editingProperty?.images || []).length > 0 && (
                        <div className="space-y-2 mt-4">
                          <p className="text-sm font-medium text-gray-700">
                            Imágenes ({((editPropertyData.images || editingProperty?.images || []).length)}):
                          </p>
                          <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
                            {(editPropertyData.images || editingProperty?.images || []).map((image, index) => (
                              <div key={index} className="relative group">
                                <img
                                  src={image}
                                  alt={`Imagen ${index + 1}`}
                                  className="w-full h-24 object-cover rounded-lg border border-gray-200"
                                  onError={(e) => {
                                    e.target.src = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="100" height="100"%3E%3Crect width="100" height="100" fill="%23f3f4f6"/%3E%3Ctext x="50" y="50" text-anchor="middle" fill="%239ca3af" font-size="12"%3EError%3C/text%3E%3C/svg%3E';
                                  }}
                                />
                                <button
                                  type="button"
                                  onClick={() => handleRemoveEditImage(image)}
                                  className="absolute top-1 right-1 bg-red-600 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition"
                                  title="Eliminar imagen"
                                >
                                  <X size={14} />
                                </button>
                              </div>
                            ))}
                          </div>
                        </div>
                      )}
                    </div>

                    <div className="flex gap-4">
                      <button
                        type="submit"
                        disabled={propertiesLoading}
                        className="flex-1 bg-primary text-white py-2 rounded-lg hover:bg-gray-800 transition disabled:opacity-50"
                      >
                        {propertiesLoading ? 'Guardando...' : 'Guardar Cambios'}
                      </button>
                      <button
                        type="button"
                        onClick={() => {
                          setEditingProperty(null);
                          setEditPropertyData({});
                        }}
                        disabled={propertiesLoading}
                        className="flex-1 bg-gray-200 text-gray-700 py-2 rounded-lg hover:bg-gray-300 transition disabled:opacity-50"
                      >
                        Cancelar
                      </button>
                    </div>
                  </form>
                </div>
              </div>
            )}

            {/* Modal Confirmar Eliminación */}
            {deleteConfirm && (
              <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
                <div className="bg-white rounded-2xl shadow-xl max-w-md w-full p-6">
                  <h3 className="text-xl font-bold text-primary mb-4">Confirmar Eliminación</h3>
                  <p className="text-gray-600 mb-6">
                    ¿Estás seguro de que deseas eliminar la propiedad "{deleteConfirm.title}"? Esta acción no se puede deshacer.
                  </p>
                  <div className="flex gap-4">
                    <button
                      onClick={() => handleDeleteProperty(deleteConfirm.id)}
                      className="flex-1 bg-red-600 text-white py-2 rounded-lg hover:bg-red-700 transition"
                    >
                      Eliminar
                    </button>
                    <button
                      onClick={() => setDeleteConfirm(null)}
                      className="flex-1 bg-gray-200 text-gray-700 py-2 rounded-lg hover:bg-gray-300 transition"
                    >
                      Cancelar
                    </button>
                  </div>
                </div>
              </div>
            )}

            {/* Loading */}
            {propertiesLoading && (
              <div className="text-center py-12">
                <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                <p className="mt-4 text-gray-600">Cargando propiedades...</p>
              </div>
            )}

            {/* Tabla de Propiedades */}
            {!propertiesLoading && properties.length > 0 && (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Título</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Ubicación</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Precio</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Capacidad</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Disponible</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Acciones</th>
                    </tr>
                  </thead>
                  <tbody>
                    {properties.map((property) => (
                      <tr key={property.id} className="border-b border-gray-100 hover:bg-gray-50 transition">
                        <td className="py-4 px-4 font-medium">{property.title}</td>
                        <td className="py-4 px-4">{property.location}</td>
                        <td className="py-4 px-4">${property.price?.toFixed(2) || '0.00'}</td>
                        <td className="py-4 px-4">{property.capacity}</td>
                        <td className="py-4 px-4">
                          <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                            property.available 
                              ? 'bg-green-100 text-green-800' 
                              : 'bg-red-100 text-red-800'
                          }`}>
                            {property.available ? 'Sí' : 'No'}
                          </span>
                        </td>
                        <td className="py-4 px-4">
                          <div className="flex items-center gap-2">
                            <button
                              onClick={() => {
                                setEditingProperty(property);
                                setEditPropertyData({
                                  images: property.images || []
                                });
                              }}
                              className="p-2 text-blue-600 hover:bg-blue-50 rounded transition"
                              title="Editar"
                            >
                              <Edit size={18} />
                            </button>
                            <button
                              onClick={() => setDeleteConfirm(property)}
                              className="p-2 text-red-600 hover:bg-red-50 rounded transition"
                              title="Eliminar"
                            >
                              <Trash2 size={18} />
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            {/* Empty State */}
            {!propertiesLoading && properties.length === 0 && !propertiesError && (
              <div className="text-center py-12">
                <Home size={64} className="mx-auto text-gray-300 mb-4" />
                <p className="text-xl text-gray-600 mb-4">No hay propiedades registradas</p>
                <button
                  onClick={() => setShowCreateProperty(true)}
                  className="text-primary hover:underline font-medium"
                >
                  Crear primera propiedad
                </button>
              </div>
            )}
          </div>
        )}

        {/* Sección Reservas */}
        {activeTab === 'bookings' && (
          <div className="bg-white rounded-2xl shadow-lg p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold text-primary">Gestión de Reservas</h2>
            </div>

            {/* Mensajes */}
            {bookingsError && (
              <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg flex items-center gap-2">
                <AlertCircle size={20} className="text-red-600" />
                <p className="text-red-600">{bookingsError}</p>
              </div>
            )}

            {bookingsSuccess && (
              <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg flex items-center gap-2">
                <Check size={20} className="text-green-600" />
                <p className="text-green-600">{bookingsSuccess}</p>
              </div>
            )}

            {/* Loading */}
            {bookingsLoading && (
              <div className="text-center py-12">
                <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
                <p className="mt-4 text-gray-600">Cargando reservas...</p>
              </div>
            )}

            {/* Tabla de Reservas */}
            {!bookingsLoading && bookings.length > 0 && (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-200">
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Usuario</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Propiedad</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Check-in</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Check-out</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Precio Total</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Estado</th>
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Acciones</th>
                    </tr>
                  </thead>
                  <tbody>
                    {bookings.map((booking) => {
                      const isCancelled = booking.status?.toUpperCase() === 'CANCELLED';
                      
                      return (
                        <tr key={booking.id} className="border-b border-gray-100 hover:bg-gray-50 transition">
                          <td className="py-4 px-4">
                            <span className="text-gray-700">Usuario {booking.userId}</span>
                          </td>
                          <td className="py-4 px-4">
                            <div className="flex items-center gap-3">
                              {propertyDetails[booking.propertyId]?.images?.[0] ? (
                                <img
                                  src={propertyDetails[booking.propertyId].images[0]}
                                  alt={propertyDetails[booking.propertyId]?.title || 'Propiedad'}
                                  className="w-12 h-12 object-cover rounded-lg"
                                />
                              ) : null}
                              <div>
                                <p className="font-medium text-gray-900">
                                  {propertyDetails[booking.propertyId]?.title || `Propiedad ${booking.propertyId}`}
                                </p>
                                <p className="text-sm text-gray-500">
                                  {propertyDetails[booking.propertyId]?.location || booking.propertyId}
                                </p>
                              </div>
                            </div>
                          </td>
                          <td className="py-4 px-4">
                            <span className="text-gray-700">{formatDate(booking.checkIn)}</span>
                          </td>
                          <td className="py-4 px-4">
                            <span className="text-gray-700">{formatDate(booking.checkOut)}</span>
                          </td>
                          <td className="py-4 px-4">
                            <span className="font-semibold text-primary">
                              ${booking.totalPrice?.toFixed(2) || '0.00'}
                            </span>
                          </td>
                          <td className="py-4 px-4">
                            {getStatusBadge(booking.status)}
                          </td>
                          <td className="py-4 px-4">
                            {!isCancelled && (
                              <button
                                onClick={() => handleCancelBooking(booking.id)}
                                disabled={cancelLoading[booking.id]}
                                className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition disabled:opacity-50"
                                title="Cancelar reserva"
                              >
                                {cancelLoading[booking.id] ? (
                                  <div className="animate-spin rounded-full h-4 w-4 border-2 border-red-600 border-t-transparent"></div>
                                ) : (
                                  <Trash2 size={18} />
                                )}
                              </button>
                            )}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}

            {/* Empty State */}
            {!bookingsLoading && bookings.length === 0 && !bookingsError && (
              <div className="text-center py-12">
                <Calendar size={64} className="mx-auto text-gray-300 mb-4" />
                <p className="text-xl text-gray-600">No hay reservas registradas</p>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default Admin;

