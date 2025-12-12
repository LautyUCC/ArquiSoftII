import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { bookingsAPI, propertiesAPI } from '../services/api';
import { Calendar, MapPin, DollarSign, LogOut, Home, Edit, X, Trash2 } from 'lucide-react';

function MyBookings() {
  const navigate = useNavigate();
  const [bookings, setBookings] = useState([]);
  const [properties, setProperties] = useState([]);
  const [propertyDetails, setPropertyDetails] = useState({}); // Cache de detalles de propiedades
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedPropertyId, setSelectedPropertyId] = useState('');
  const [viewMode, setViewMode] = useState('my'); // 'my' o 'property'
  
  // Estados para edición y cancelación
  const [editingBooking, setEditingBooking] = useState(null);
  const [editDates, setEditDates] = useState({ startDate: '', endDate: '' });
  const [editLoading, setEditLoading] = useState(false);
  const [editError, setEditError] = useState('');
  const [cancelLoading, setCancelLoading] = useState({});
  const [actionError, setActionError] = useState('');
  const [actionSuccess, setActionSuccess] = useState('');
  
  // Obtener usuario del localStorage
  const user = JSON.parse(localStorage.getItem('user') || '{}');
  // Comparación robusta: convertir a mayúsculas para evitar problemas de case
  const isAdmin = user.role?.toUpperCase() === 'ADMIN';

  useEffect(() => {
    if (isAdmin && viewMode === 'property' && selectedPropertyId) {
      loadPropertyBookings(selectedPropertyId);
    } else {
      loadMyBookings();
    }
  }, [viewMode, selectedPropertyId, isAdmin]);

  // Cargar propiedades para el selector (solo ADMIN)
  useEffect(() => {
    if (isAdmin) {
      loadProperties();
    }
  }, [isAdmin]);

  const loadProperties = async () => {
    try {
      const response = await propertiesAPI.search({ q: '' });
      setProperties(response.data.results || []);
    } catch (err) {
      console.error('Error al cargar propiedades:', err);
    }
  };

  const loadMyBookings = async () => {
    setLoading(true);
    setError('');

    try {
      const response = await bookingsAPI.getMyBookings();
      // La respuesta puede venir directamente como array o dentro de un objeto data
      const bookingsData = Array.isArray(response.data) ? response.data : (response.data?.data || []);
      setBookings(bookingsData);
      
      // Cargar detalles de propiedades para las reservas
      await loadPropertyDetails(bookingsData);
      
      if (bookingsData.length === 0) {
        setError('No tienes reservas aún');
      }
    } catch (err) {
      console.error('Error al cargar reservas:', err);
      
      if (err.response?.status === 401) {
        setError('Sesión expirada. Por favor inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.code === 'ERR_NETWORK') {
        setError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setError('Error al cargar las reservas. Por favor intenta de nuevo.');
      }
    } finally {
      setLoading(false);
    }
  };

  const loadPropertyDetails = async (bookingsData) => {
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

  const loadPropertyBookings = async (propertyId) => {
    if (!propertyId) {
      setBookings([]);
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await bookingsAPI.getPropertyBookings(propertyId);
      const bookingsData = response.data || [];
      setBookings(bookingsData);
      
      // Cargar detalles de la propiedad
      if (bookingsData.length > 0 && !propertyDetails[propertyId]) {
        try {
          const propResponse = await propertiesAPI.getById(propertyId);
          setPropertyDetails(prev => ({ ...prev, [propertyId]: propResponse.data }));
        } catch (err) {
          console.error(`Error al cargar propiedad ${propertyId}:`, err);
        }
      }
      
      if (bookingsData.length === 0) {
        setError('No hay reservas para esta propiedad');
      }
    } catch (err) {
      console.error('Error al cargar reservas de propiedad:', err);
      
      if (err.response?.status === 401) {
        setError('Sesión expirada. Por favor inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 403) {
        setError('No tienes permisos para ver estas reservas. Se requiere rol ADMIN.');
      } else if (err.response?.status === 400) {
        setError('Propiedad no encontrada');
      } else if (err.code === 'ERR_NETWORK') {
        setError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setError('Error al cargar las reservas. Por favor intenta de nuevo.');
      }
    } finally {
      setLoading(false);
    }
  };

  const handlePropertyChange = (e) => {
    const propertyId = e.target.value;
    setSelectedPropertyId(propertyId);
    if (propertyId) {
      setViewMode('property');
    } else {
      setViewMode('my');
      loadMyBookings();
    }
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate('/');
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

  const handleOpenEdit = (booking) => {
    setEditingBooking(booking);
    // Convertir fechas a formato YYYY-MM-DD para los inputs
    const startDate = booking.checkIn ? new Date(booking.checkIn).toISOString().split('T')[0] : '';
    const endDate = booking.checkOut ? new Date(booking.checkOut).toISOString().split('T')[0] : '';
    setEditDates({ startDate, endDate });
    setEditError('');
    setActionError('');
    setActionSuccess('');
  };

  const handleCloseEdit = () => {
    setEditingBooking(null);
    setEditDates({ startDate: '', endDate: '' });
    setEditError('');
    setActionError('');
    setActionSuccess('');
  };

  const handleUpdateBooking = async () => {
    if (!editingBooking) return;

    setEditLoading(true);
    setEditError('');
    setActionError('');
    setActionSuccess('');

    // Validaciones
    if (!editDates.startDate || !editDates.endDate) {
      setEditError('Por favor selecciona ambas fechas');
      setEditLoading(false);
      return;
    }

    const startDate = new Date(editDates.startDate);
    const endDate = new Date(editDates.endDate);
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    if (startDate < today) {
      setEditError('La fecha de inicio no puede ser en el pasado');
      setEditLoading(false);
      return;
    }

    if (endDate <= startDate) {
      setEditError('La fecha de fin debe ser posterior a la fecha de inicio');
      setEditLoading(false);
      return;
    }

    try {
      // Enviar PATCH con startDate y endDate
      await bookingsAPI.updateBooking(editingBooking.id, {
        startDate: editDates.startDate,
        endDate: editDates.endDate,
        // También enviar como checkIn/checkOut para compatibilidad
        checkIn: editDates.startDate,
        checkOut: editDates.endDate
      });

      setActionSuccess('Reserva actualizada exitosamente');
      handleCloseEdit();
      
      // Recargar reservas después de un breve delay
      setTimeout(() => {
        if (viewMode === 'property' && selectedPropertyId) {
          loadPropertyBookings(selectedPropertyId);
        } else {
          loadMyBookings();
        }
        setActionSuccess('');
      }, 1500);
    } catch (err) {
      console.error('Error al actualizar reserva:', err);
      
      if (err.response?.status === 401 || err.response?.status === 403) {
        setEditError('Tu sesión ha expirado o no tienes permisos. Por favor, inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 409) {
        setEditError('La propiedad ya está reservada en esas fechas. Por favor, selecciona otras fechas.');
      } else if (err.response?.status === 400) {
        const errorMessage = err.response?.data?.message || err.response?.data?.error || 'Error de validación. Por favor, verifica los datos ingresados.';
        setEditError(errorMessage);
      } else {
        setEditError('Error al actualizar la reserva. Por favor, intenta de nuevo.');
      }
    } finally {
      setEditLoading(false);
    }
  };

  const handleCancelBooking = async (bookingId) => {
    if (!window.confirm('¿Estás seguro de que deseas cancelar esta reserva?')) {
      return;
    }

    setCancelLoading(prev => ({ ...prev, [bookingId]: true }));
    setActionError('');
    setActionSuccess('');

    try {
      await bookingsAPI.deleteBooking(bookingId);
      
      setActionSuccess('Reserva cancelada exitosamente');
      
      // Actualizar la lista localmente marcando como CANCELLED
      setBookings(prev => prev.map(b => 
        b.id === bookingId ? { ...b, status: 'CANCELLED' } : b
      ));
      
      // También actualizar en propertyDetails si es necesario
      // (no es necesario, pero mantiene consistencia)
      
      setTimeout(() => {
        setActionSuccess('');
        // Recargar reservas
        if (viewMode === 'property' && selectedPropertyId) {
          loadPropertyBookings(selectedPropertyId);
        } else {
          loadMyBookings();
        }
      }, 1500);
    } catch (err) {
      console.error('Error al cancelar reserva:', err);
      
      if (err.response?.status === 401 || err.response?.status === 403) {
        setActionError('Tu sesión ha expirado o no tienes permisos. Por favor, inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 404) {
        setActionError('Reserva no encontrada');
      } else {
        setActionError('Error al cancelar la reserva. Por favor, intenta de nuevo.');
      }
    } finally {
      setCancelLoading(prev => ({ ...prev, [bookingId]: false }));
    }
  };

  return (
    <div className="min-h-screen bg-secondary">
      {/* Header */}
      <header className="bg-white shadow-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <h1 className="text-2xl font-bold text-primary">Spotly</h1>
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
        <div className="bg-white rounded-2xl shadow-lg p-6 mb-6">
          <h2 className="text-3xl font-bold text-primary mb-6">Mis Reservas</h2>

          {/* Selector de modo (solo ADMIN) */}
          {isAdmin && (
            <div className="mb-6 p-4 bg-gray-50 rounded-lg">
              <div className="flex items-center gap-4 flex-wrap">
                <div className="flex items-center gap-2">
                  <input
                    type="radio"
                    id="view-my"
                    name="view-mode"
                    checked={viewMode === 'my'}
                    onChange={() => {
                      setViewMode('my');
                      setSelectedPropertyId('');
                      loadMyBookings();
                    }}
                    className="w-4 h-4 text-primary"
                  />
                  <label htmlFor="view-my" className="text-gray-700 font-medium cursor-pointer">
                    Mis Reservas
                  </label>
                </div>
                <div className="flex items-center gap-2">
                  <input
                    type="radio"
                    id="view-property"
                    name="view-mode"
                    checked={viewMode === 'property'}
                    onChange={() => {
                      setViewMode('property');
                    }}
                    className="w-4 h-4 text-primary"
                  />
                  <label htmlFor="view-property" className="text-gray-700 font-medium cursor-pointer">
                    Reservas por Propiedad
                  </label>
                </div>
              </div>

              {/* Selector de propiedad (solo cuando está en modo property) */}
              {viewMode === 'property' && (
                <div className="mt-4">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Seleccionar Propiedad
                  </label>
                  <select
                    value={selectedPropertyId}
                    onChange={handlePropertyChange}
                    className="w-full md:w-1/2 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                  >
                    <option value="">-- Selecciona una propiedad --</option>
                    {properties.map((property) => (
                      <option key={property.id} value={property.id}>
                        {property.title} - {property.city}, {property.country}
                      </option>
                    ))}
                  </select>
                </div>
              )}
            </div>
          )}

          {/* Loading State */}
          {loading && (
            <div className="text-center py-12">
              <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
              <p className="mt-4 text-gray-600">Cargando reservas...</p>
            </div>
          )}

          {/* Error State */}
          {!loading && error && (
            <div className="text-center py-12">
              <p className="text-xl text-gray-600 mb-4">{error}</p>
              {error !== 'No tienes reservas aún' && error !== 'No hay reservas para esta propiedad' && (
                <button
                  onClick={() => {
                    if (viewMode === 'property' && selectedPropertyId) {
                      loadPropertyBookings(selectedPropertyId);
                    } else {
                      loadMyBookings();
                    }
                  }}
                  className="text-primary hover:underline font-medium"
                >
                  Intentar de nuevo
                </button>
              )}
            </div>
          )}

          {/* Bookings Table */}
          {!loading && !error && bookings.length > 0 && (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-200">
                    <th className="text-left py-3 px-4 font-semibold text-gray-700">Propiedad</th>
                    <th className="text-left py-3 px-4 font-semibold text-gray-700">Check-in</th>
                    <th className="text-left py-3 px-4 font-semibold text-gray-700">Check-out</th>
                    <th className="text-left py-3 px-4 font-semibold text-gray-700">Precio Total</th>
                    <th className="text-left py-3 px-4 font-semibold text-gray-700">Estado</th>
                    {isAdmin && viewMode === 'property' && (
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Usuario</th>
                    )}
                    {viewMode === 'my' && (
                      <th className="text-left py-3 px-4 font-semibold text-gray-700">Acciones</th>
                    )}
                  </tr>
                </thead>
                <tbody>
                  {bookings.map((booking) => {
                    const isConfirmed = booking.status?.toUpperCase() === 'CONFIRMED';
                    const isCancelled = booking.status?.toUpperCase() === 'CANCELLED';
                    
                    return (
                      <tr key={booking.id} className="border-b border-gray-100 hover:bg-gray-50 transition">
                        <td className="py-4 px-4">
                          <div className="flex items-center gap-3">
                            {propertyDetails[booking.propertyId]?.images?.[0] ? (
                              <img
                                src={propertyDetails[booking.propertyId].images[0]}
                                alt={propertyDetails[booking.propertyId]?.title || 'Propiedad'}
                                className="w-16 h-16 object-cover rounded-lg"
                              />
                            ) : (
                              <MapPin size={24} className="text-gray-400" />
                            )}
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
                          <div className="flex items-center gap-2">
                            <Calendar size={16} className="text-gray-400" />
                            <span className="text-gray-700">{formatDate(booking.checkIn)}</span>
                          </div>
                        </td>
                        <td className="py-4 px-4">
                          <div className="flex items-center gap-2">
                            <Calendar size={16} className="text-gray-400" />
                            <span className="text-gray-700">{formatDate(booking.checkOut)}</span>
                          </div>
                        </td>
                        <td className="py-4 px-4">
                          <div className="flex items-center gap-2">
                            <DollarSign size={16} className="text-gray-400" />
                            <span className="font-semibold text-primary">
                              ${booking.totalPrice?.toFixed(2) || '0.00'}
                            </span>
                          </div>
                        </td>
                        <td className="py-4 px-4">
                          {getStatusBadge(booking.status)}
                        </td>
                        {isAdmin && viewMode === 'property' && (
                          <td className="py-4 px-4">
                            <span className="text-gray-700">Usuario {booking.userId}</span>
                          </td>
                        )}
                        {viewMode === 'my' && (
                          <td className="py-4 px-4">
                            <div className="flex items-center gap-2">
                              {isConfirmed && (
                                <button
                                  onClick={() => handleOpenEdit(booking)}
                                  className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition"
                                  title="Editar fechas"
                                >
                                  <Edit size={18} />
                                </button>
                              )}
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
                            </div>
                          </td>
                        )}
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}

          {/* Empty State */}
          {!loading && !error && bookings.length === 0 && (
            <div className="text-center py-12">
              <Calendar size={64} className="mx-auto text-gray-300 mb-4" />
              <p className="text-xl text-gray-600 mb-2">
                {viewMode === 'property' 
                  ? 'No hay reservas para esta propiedad' 
                  : 'No tienes reservas aún'}
              </p>
              {viewMode === 'my' && (
                <button
                  onClick={() => navigate('/search')}
                  className="text-primary hover:underline font-medium"
                >
                  Explorar propiedades
                </button>
              )}
            </div>
          )}

          {/* Mensajes de acción */}
          {actionError && (
            <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-red-600 text-sm">{actionError}</p>
            </div>
          )}
          {actionSuccess && (
            <div className="mt-4 p-4 bg-green-50 border border-green-200 rounded-lg">
              <p className="text-green-600 text-sm">{actionSuccess}</p>
            </div>
          )}
        </div>
      </div>

      {/* Modal de Edición de Fechas */}
      {editingBooking && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl max-w-md w-full">
            <div className="sticky top-0 bg-white border-b px-6 py-4 flex items-center justify-between">
              <h2 className="text-2xl font-bold text-primary">Editar Fechas de Reserva</h2>
              <button
                onClick={handleCloseEdit}
                className="text-gray-500 hover:text-gray-700 transition"
                disabled={editLoading}
              >
                <X size={24} />
              </button>
            </div>

            <div className="p-6">
              {editError && (
                <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                  <p className="text-red-600 text-sm">{editError}</p>
                </div>
              )}

              <div className="space-y-4 mb-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Fecha de Inicio *
                  </label>
                  <input
                    type="date"
                    value={editDates.startDate}
                    onChange={(e) => setEditDates({ ...editDates, startDate: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    min={new Date().toISOString().split('T')[0]}
                    disabled={editLoading}
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Fecha de Fin *
                  </label>
                  <input
                    type="date"
                    value={editDates.endDate}
                    onChange={(e) => setEditDates({ ...editDates, endDate: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    min={editDates.startDate || new Date().toISOString().split('T')[0]}
                    disabled={editLoading}
                  />
                </div>
              </div>

              <div className="flex gap-4">
                <button
                  onClick={handleCloseEdit}
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition"
                  disabled={editLoading}
                >
                  Cancelar
                </button>
                <button
                  onClick={handleUpdateBooking}
                  className="flex-1 px-4 py-2 bg-primary text-white rounded-lg hover:bg-gray-800 transition disabled:opacity-50"
                  disabled={editLoading}
                >
                  {editLoading ? (
                    <span className="flex items-center justify-center gap-2">
                      <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent"></div>
                      Guardando...
                    </span>
                  ) : (
                    'Guardar Cambios'
                  )}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default MyBookings;

