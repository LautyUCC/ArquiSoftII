import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { propertiesAPI, bookingsAPI } from '../services/api';
import { MapPin, Users, ArrowLeft, Check, Edit, X, Calendar } from 'lucide-react';

function PropertyDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [property, setProperty] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isAdmin, setIsAdmin] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [editLoading, setEditLoading] = useState(false);
  const [editError, setEditError] = useState('');
  const [editSuccess, setEditSuccess] = useState('');
  const [editData, setEditData] = useState({
    title: '',
    description: '',
    price: '',
    location: '',
    capacity: '',
    available: true,
    amenities: [],
    images: []
  });
  const [booking, setBooking] = useState(false);
  const [bookingLoading, setBookingLoading] = useState(false);
  const [bookingSuccess, setBookingSuccess] = useState('');
  const [bookingData, setBookingData] = useState({
    checkIn: '',
    checkOut: '',
    guests: 1
  });
  const [bookingError, setBookingError] = useState('');
  const [propertyBookings, setPropertyBookings] = useState([]); // Reservas confirmadas de la propiedad
  const [bookingsLoading, setBookingsLoading] = useState(false);

  useEffect(() => {
    // Verificar si el usuario es ADMIN
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        setIsAdmin(user.role?.toUpperCase() === 'ADMIN');
      } catch (err) {
        console.error('Error parseando usuario:', err);
      }
    }
    loadProperty();
    loadPropertyBookings();
  }, [id]);

  const loadProperty = async () => {
    try {
      const response = await propertiesAPI.getById(id);
      const propertyData = response.data?.data || response.data;
      setProperty(propertyData);
      
      // Prellenar datos de edición
      if (propertyData) {
        setEditData({
          title: propertyData.title || '',
          description: propertyData.description || '',
          price: propertyData.price?.toString() || '',
          location: propertyData.location || '',
          capacity: propertyData.capacity?.toString() || '',
          available: propertyData.available !== undefined ? propertyData.available : true,
          amenities: propertyData.amenities || [],
          images: propertyData.images || []
        });
      }
    } catch (err) {
      console.error('Error al cargar propiedad:', err);
    } finally {
      setLoading(false);
    }
  };

  const loadPropertyBookings = async () => {
    setBookingsLoading(true);
    try {
      const response = await bookingsAPI.getPropertyBookings(id);
      const bookings = Array.isArray(response.data) ? response.data : (response.data?.data || []);
      setPropertyBookings(bookings);
    } catch (err) {
      console.error('Error al cargar reservas de la propiedad:', err);
      // No es crítico si falla, simplemente no bloquearemos fechas
      setPropertyBookings([]);
    } finally {
      setBookingsLoading(false);
    }
  };

  // Función para verificar si una fecha está en un rango de reservas
  const isDateInBookingRange = (date, booking) => {
    const checkDate = new Date(date);
    checkDate.setHours(0, 0, 0, 0);
    
    const checkIn = new Date(booking.checkIn);
    checkIn.setHours(0, 0, 0, 0);
    
    const checkOut = new Date(booking.checkOut);
    checkOut.setHours(0, 0, 0, 0);
    
    // checkOut es exclusivo, así que la fecha debe ser >= checkIn y < checkOut
    return checkDate >= checkIn && checkDate < checkOut;
  };

  // Función para verificar si un rango de fechas se solapa con alguna reserva
  const hasDateRangeOverlap = (startDate, endDate) => {
    if (!startDate || !endDate) return false;
    
    const start = new Date(startDate);
    start.setHours(0, 0, 0, 0);
    
    const end = new Date(endDate);
    end.setHours(0, 0, 0, 0);
    
    return propertyBookings.some(booking => {
      const bookingStart = new Date(booking.checkIn);
      bookingStart.setHours(0, 0, 0, 0);
      
      const bookingEnd = new Date(booking.checkOut);
      bookingEnd.setHours(0, 0, 0, 0);
      
      // Verificar solapamiento: el rango se solapa si:
      // - start < bookingEnd Y end > bookingStart
      return start < bookingEnd && end > bookingStart;
    });
  };

  // Función para verificar si una fecha específica está bloqueada
  const isDateBlocked = (dateString) => {
    if (!dateString) return false;
    return propertyBookings.some(booking => isDateInBookingRange(dateString, booking));
  };

  // Función para obtener los días bloqueados como string legible
  const getBlockedDatesInfo = () => {
    if (propertyBookings.length === 0) return null;
    
    return propertyBookings.map(booking => {
      const checkIn = new Date(booking.checkIn);
      const checkOut = new Date(booking.checkOut);
      checkOut.setDate(checkOut.getDate() - 1); // checkOut es exclusivo, mostrar hasta el día anterior
      
      return `${checkIn.toLocaleDateString('es-ES')} - ${checkOut.toLocaleDateString('es-ES')}`;
    }).join(', ');
  };

  const handleOpenEdit = () => {
    setShowEditModal(true);
    setEditError('');
    setEditSuccess('');
  };

  const handleCloseEdit = () => {
    setShowEditModal(false);
    setEditError('');
    setEditSuccess('');
    // Restaurar datos originales
    if (property) {
      setEditData({
        title: property.title || '',
        description: property.description || '',
        price: property.price?.toString() || '',
        location: property.location || '',
        capacity: property.capacity?.toString() || '',
        available: property.available !== undefined ? property.available : true,
        amenities: property.amenities || [],
        images: property.images || []
      });
    }
  };

  const handleEditSubmit = async (e) => {
    e.preventDefault();
    setEditLoading(true);
    setEditError('');
    setEditSuccess('');

    // Validaciones
    if (!editData.title.trim()) {
      setEditError('El título es requerido');
      setEditLoading(false);
      return;
    }

    if (!editData.description.trim()) {
      setEditError('La descripción es requerida');
      setEditLoading(false);
      return;
    }

    if (!editData.price || parseFloat(editData.price) <= 0) {
      setEditError('El precio debe ser mayor a 0');
      setEditLoading(false);
      return;
    }

    if (!editData.location.trim()) {
      setEditError('La ubicación es requerida');
      setEditLoading(false);
      return;
    }

    if (!editData.capacity || parseInt(editData.capacity) <= 0) {
      setEditError('La capacidad debe ser mayor a 0');
      setEditLoading(false);
      return;
    }

    try {
      const updateData = {
        title: editData.title.trim(),
        description: editData.description.trim(),
        price: parseFloat(editData.price),
        location: editData.location.trim(),
        capacity: parseInt(editData.capacity),
        available: editData.available,
        amenities: editData.amenities.filter(a => a.trim() !== ''),
        images: editData.images.filter(img => img.trim() !== '')
      };

      await propertiesAPI.update(id, updateData);
      setEditSuccess('Propiedad actualizada exitosamente');
      
      // Recargar propiedad después de un breve delay
      setTimeout(() => {
        loadProperty();
        setShowEditModal(false);
      }, 1500);
    } catch (err) {
      console.error('Error al actualizar propiedad:', err);
      if (err.response?.status === 401) {
        setEditError('Tu sesión ha expirado. Por favor, inicia sesión nuevamente.');
      } else if (err.response?.status === 403) {
        setEditError('No tienes permisos para editar propiedades. Solo los administradores pueden editar.');
      } else {
        setEditError(err.response?.data?.message || 'Error al actualizar la propiedad. Por favor intenta de nuevo.');
      }
    } finally {
      setEditLoading(false);
    }
  };

  const handleAddAmenity = () => {
    setEditData({
      ...editData,
      amenities: [...editData.amenities, '']
    });
  };

  const handleRemoveAmenity = (index) => {
    setEditData({
      ...editData,
      amenities: editData.amenities.filter((_, i) => i !== index)
    });
  };

  const handleAmenityChange = (index, value) => {
    const newAmenities = [...editData.amenities];
    newAmenities[index] = value;
    setEditData({
      ...editData,
      amenities: newAmenities
    });
  };

  const handleAddImage = () => {
    setEditData({
      ...editData,
      images: [...editData.images, '']
    });
  };

  const handleRemoveImage = (index) => {
    setEditData({
      ...editData,
      images: editData.images.filter((_, i) => i !== index)
    });
  };

  const handleImageChange = (index, value) => {
    const newImages = [...editData.images];
    newImages[index] = value;
    setEditData({
      ...editData,
      images: newImages
    });
  };

  const handleBooking = async () => {
    setBookingError('');
    setBookingSuccess('');
    setBookingLoading(true);

    // Validaciones
    if (!bookingData.checkIn) {
      setBookingError('Por favor selecciona la fecha de entrada');
      setBookingLoading(false);
      return;
    }

    if (!bookingData.checkOut) {
      setBookingError('Por favor selecciona la fecha de salida');
      setBookingLoading(false);
      return;
    }

    const checkIn = new Date(bookingData.checkIn);
    const checkOut = new Date(bookingData.checkOut);
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    if (checkIn < today) {
      setBookingError('La fecha de entrada no puede ser en el pasado');
      setBookingLoading(false);
      return;
    }

    if (checkOut <= checkIn) {
      setBookingError('La fecha de salida debe ser posterior a la entrada');
      setBookingLoading(false);
      return;
    }

    // Validar que el rango no se solape con reservas existentes
    if (hasDateRangeOverlap(bookingData.checkIn, bookingData.checkOut)) {
      setBookingError('El rango de fechas seleccionado incluye días ya reservados. Por favor, selecciona otras fechas.');
      setBookingLoading(false);
      return;
    }

    // Calcular número de noches
    const nights = Math.ceil((checkOut - checkIn) / (1000 * 60 * 60 * 24));

    // Validar que la reserva no sea mayor a 30 noches
    if (nights > 30) {
      setBookingError('La reserva no puede ser mayor a 30 noches');
      setBookingLoading(false);
      return;
    }

    if (bookingData.guests < 1) {
      setBookingError('Debe haber al menos 1 huésped');
      setBookingLoading(false);
      return;
    }

    if (bookingData.guests > property.capacity) {
      setBookingError(`Esta propiedad solo tiene capacidad para ${property.capacity} huéspedes`);
      setBookingLoading(false);
      return;
    }

    // Verificar que el usuario esté autenticado
    const token = localStorage.getItem('token');
    if (!token) {
      setBookingError('Debes iniciar sesión para realizar una reserva');
      setBookingLoading(false);
      setTimeout(() => {
        navigate('/login');
      }, 2000);
      return;
    }

    try {
      // Preparar datos para el API
      // El backend acepta startDate/endDate o checkIn/checkOut
      const bookingPayload = {
        propertyId: property.id,
        checkIn: bookingData.checkIn,
        checkOut: bookingData.checkOut,
        // También enviar como startDate/endDate para compatibilidad
        startDate: bookingData.checkIn,
        endDate: bookingData.checkOut
      };

      // Llamar al API
      const response = await bookingsAPI.createBooking(bookingPayload);

      // Si la respuesta es 201, mostrar éxito y redirigir
      if (response.status === 201 || response.data) {
        setBookingSuccess('¡Reserva creada exitosamente!');
        setBookingError('');
        
        // Redirigir a /my-bookings después de 1.5 segundos
        setTimeout(() => {
          navigate('/my-bookings');
        }, 1500);
      }
    } catch (err) {
      console.error('Error al crear reserva:', err);
      
      // Manejar diferentes códigos de error
      if (err.response?.status === 401 || err.response?.status === 403) {
        // Sesión expirada o sin permisos - redirigir a login
        setBookingError('Tu sesión ha expirado. Por favor, inicia sesión nuevamente.');
        setTimeout(() => {
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 409) {
        // Conflicto de solapamiento
        setBookingError('La propiedad ya está reservada en esas fechas. Por favor, selecciona otras fechas.');
      } else if (err.response?.status === 400) {
        // Error de validación
        const errorMessage = err.response?.data?.message || err.response?.data?.error || 'Error de validación. Por favor, verifica los datos ingresados.';
        setBookingError(errorMessage);
      } else {
        // Error genérico
        setBookingError('Error al crear la reserva. Por favor, intenta de nuevo.');
      }
    } finally {
      setBookingLoading(false);
    }
  };

  if (loading) {
    return (
        <div className="min-h-screen bg-secondary flex items-center justify-center">
          <div className="animate-spin rounded-full h-12 w-12 border-4 border-gray-200 border-t-primary"></div>
        </div>
    );
  }

  if (!property) {
    return (
        <div className="min-h-screen bg-secondary flex items-center justify-center">
          <p className="text-xl text-gray-600">Propiedad no encontrada</p>
        </div>
    );
  }

  return (
      <div className="min-h-screen bg-secondary">
        {/* Header */}
        <header className="bg-white shadow-sm sticky top-0 z-50">
          <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
            <button
                onClick={() => navigate('/search')}
                className="flex items-center gap-2 text-gray-600 hover:text-primary transition"
            >
              <ArrowLeft size={20} />
              Volver a búsqueda
            </button>
            
            <div className="flex items-center gap-4">
              <button
                onClick={() => navigate('/my-bookings')}
                className="flex items-center gap-2 text-gray-600 hover:text-primary transition"
              >
                <Calendar size={20} />
                Mis Reservas
              </button>
              {isAdmin && (
                <button
                    onClick={handleOpenEdit}
                    className="flex items-center gap-2 bg-primary text-white px-4 py-2 rounded-lg hover:bg-gray-800 transition"
                >
                  <Edit size={20} />
                  Editar Propiedad
                </button>
              )}
            </div>
          </div>
        </header>

        <div className="max-w-7xl mx-auto px-4 py-12">
          <div className="bg-white rounded-2xl shadow-xl overflow-hidden">
            {/* Image */}
            <div className="aspect-[21/9] bg-gray-200 relative">
              {property.images && property.images[0] ? (
                  <img
                      src={property.images[0]}
                      alt={property.title}
                      className="w-full h-full object-cover"
                  />
              ) : (
                  <div className="w-full h-full flex items-center justify-center text-gray-400">
                    <MapPin size={64} />
                  </div>
              )}
            </div>

            <div className="p-8 md:p-12">
              <div className="grid md:grid-cols-3 gap-12">
                {/* Left: Info */}
                <div className="md:col-span-2">
                  <h1 className="text-4xl font-bold text-primary mb-4">
                    {property.title}
                  </h1>

                  <div className="flex items-center text-gray-600 mb-6">
                    <MapPin size={20} className="mr-2" />
                    <span className="text-lg">{property.location}</span>
                  </div>

                  <div className="flex gap-6 mb-8 pb-8 border-b">
                    <div className="flex items-center gap-2">
                      <Users size={20} className="text-gray-600" />
                      <span>{property.capacity} huéspedes</span>
                    </div>
                  </div>

                  <div className="mb-8">
                    <h2 className="text-2xl font-bold text-primary mb-4">Descripción</h2>
                    <p className="text-gray-600 leading-relaxed">
                      {property.description}
                    </p>
                  </div>

                  {property.amenities && property.amenities.length > 0 && (
                      <div>
                        <h2 className="text-2xl font-bold text-primary mb-4">Amenidades</h2>
                        <div className="grid grid-cols-2 gap-4">
                          {property.amenities.map((amenity, index) => (
                              <div key={index} className="flex items-center gap-2">
                                <Check size={20} className="text-green-600" />
                                <span className="text-gray-700 capitalize">{amenity}</span>
                              </div>
                          ))}
                        </div>
                      </div>
                  )}
                </div>

                {/* Right: Booking Card */}
                <div className="md:col-span-1">
                  <div className="sticky top-24 bg-gray-50 rounded-2xl p-6 shadow-lg">
                    <div className="mb-6">
                      <p className="text-3xl font-bold text-primary">
                        ${Math.round(property.price)}
                      </p>
                      <p className="text-gray-600">por noche</p>
                    </div>

                    {bookingError && (
                        <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                          <p className="text-red-600 text-sm">{bookingError}</p>
                        </div>
                    )}

                    {bookingSuccess && (
                        <div className="mb-4 p-3 bg-green-50 border border-green-200 rounded-lg">
                          <p className="text-green-600 text-sm">{bookingSuccess}</p>
                        </div>
                    )}

                    <div className="space-y-4 mb-6">
                      {/* Mostrar información de fechas bloqueadas */}
                      {propertyBookings.length > 0 && (
                        <div className="p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                          <p className="text-xs font-medium text-yellow-800 mb-1">
                            📅 Fechas no disponibles:
                          </p>
                          <p className="text-xs text-yellow-700">
                            {getBlockedDatesInfo()}
                          </p>
                        </div>
                      )}

                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Fecha de entrada
                        </label>
                        <input
                            type="date"
                            value={bookingData.checkIn}
                            onChange={(e) => {
                              const selectedDate = e.target.value;
                              setBookingData({...bookingData, checkIn: selectedDate});
                              
                              // Validar si la fecha seleccionada está bloqueada
                              if (isDateBlocked(selectedDate)) {
                                setBookingError('Esta fecha no está disponible. Por favor, selecciona otra fecha.');
                              } else {
                                setBookingError('');
                                // Si checkOut está seleccionado y ahora hay solapamiento, limpiar checkOut
                                if (bookingData.checkOut && hasDateRangeOverlap(selectedDate, bookingData.checkOut)) {
                                  setBookingData(prev => ({...prev, checkOut: ''}));
                                  setBookingError('El rango seleccionado incluye fechas no disponibles. Por favor, selecciona nuevas fechas.');
                                }
                              }
                            }}
                            className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent ${
                              bookingData.checkIn && isDateBlocked(bookingData.checkIn)
                                ? 'border-red-500 bg-red-50'
                                : 'border-gray-300'
                            }`}
                            min={new Date().toISOString().split('T')[0]}
                        />
                        {bookingData.checkIn && isDateBlocked(bookingData.checkIn) && (
                          <p className="mt-1 text-xs text-red-600">
                            ⚠️ Esta fecha está reservada
                          </p>
                        )}
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Fecha de salida
                        </label>
                        <input
                            type="date"
                            value={bookingData.checkOut}
                            onChange={(e) => {
                              const selectedDate = e.target.value;
                              setBookingData({...bookingData, checkOut: selectedDate});
                              
                              // Validar si el rango completo se solapa
                              if (bookingData.checkIn && hasDateRangeOverlap(bookingData.checkIn, selectedDate)) {
                                setBookingError('El rango de fechas seleccionado incluye días ya reservados. Por favor, selecciona otras fechas.');
                              } else if (isDateBlocked(selectedDate)) {
                                setBookingError('Esta fecha no está disponible. Por favor, selecciona otra fecha.');
                              } else {
                                setBookingError('');
                              }
                            }}
                            className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent ${
                              bookingData.checkOut && (isDateBlocked(bookingData.checkOut) || 
                                (bookingData.checkIn && hasDateRangeOverlap(bookingData.checkIn, bookingData.checkOut)))
                                ? 'border-red-500 bg-red-50'
                                : 'border-gray-300'
                            }`}
                            min={bookingData.checkIn || new Date().toISOString().split('T')[0]}
                        />
                        {bookingData.checkOut && (
                          (isDateBlocked(bookingData.checkOut) || 
                           (bookingData.checkIn && hasDateRangeOverlap(bookingData.checkIn, bookingData.checkOut))) && (
                            <p className="mt-1 text-xs text-red-600">
                              ⚠️ Este rango incluye fechas reservadas
                            </p>
                          )
                        )}
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Número de huéspedes
                        </label>
                        <input
                            type="number"
                            value={bookingData.guests}
                            onChange={(e) => {
                              const value = parseInt(e.target.value) || 1;
                              setBookingData({...bookingData, guests: Math.min(value, property.capacity)});
                            }}
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                            min="1"
                            max={property.capacity}
                        />
                        <p className="text-xs text-gray-500 mt-1">
                          Máximo: {property.capacity} huéspedes
                        </p>
                      </div>
                    </div>

                    <button
                        onClick={handleBooking}
                        disabled={bookingLoading || bookingSuccess}
                        className="w-full bg-primary text-white py-4 rounded-lg font-bold text-lg hover:bg-gray-800 transition disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {bookingLoading ? (
                        <span className="flex items-center justify-center gap-2">
                          <div className="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
                          Reservando...
                        </span>
                      ) : bookingSuccess ? (
                        'Reserva Creada ✓'
                      ) : (
                        'Reservar'
                      )}
                    </button>

                    <p className="text-center text-sm text-gray-500 mt-4">
                      No se te cobrará todavía
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Modal de Edición (solo para ADMIN) */}
        {showEditModal && isAdmin && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b px-6 py-4 flex items-center justify-between">
              <h2 className="text-2xl font-bold text-primary">Editar Propiedad</h2>
              <button
                  onClick={handleCloseEdit}
                  className="text-gray-500 hover:text-gray-700 transition"
              >
                <X size={24} />
              </button>
            </div>

            <form onSubmit={handleEditSubmit} className="p-6 space-y-6">
              {editError && (
                <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
                  <p className="text-red-600">{editError}</p>
                </div>
              )}

              {editSuccess && (
                <div className="p-4 bg-green-50 border border-green-200 rounded-lg">
                  <p className="text-green-600">{editSuccess}</p>
                </div>
              )}

              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Título *
                  </label>
                  <input
                      type="text"
                      value={editData.title}
                      onChange={(e) => setEditData({...editData, title: e.target.value})}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Precio por noche *
                  </label>
                  <input
                      type="number"
                      step="0.01"
                      min="0"
                      value={editData.price}
                      onChange={(e) => setEditData({...editData, price: e.target.value})}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      required
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Descripción *
                </label>
                <textarea
                    value={editData.description}
                    onChange={(e) => setEditData({...editData, description: e.target.value})}
                    rows="4"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    required
                />
              </div>

              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Ubicación *
                  </label>
                  <input
                      type="text"
                      value={editData.location}
                      onChange={(e) => setEditData({...editData, location: e.target.value})}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Capacidad (huéspedes) *
                  </label>
                  <input
                      type="number"
                      min="1"
                      value={editData.capacity}
                      onChange={(e) => setEditData({...editData, capacity: e.target.value})}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      required
                  />
                </div>
              </div>

              <div>
                <label className="flex items-center gap-2">
                  <input
                      type="checkbox"
                      checked={editData.available}
                      onChange={(e) => setEditData({...editData, available: e.target.checked})}
                      className="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary"
                  />
                  <span className="text-sm font-medium text-gray-700">Disponible</span>
                </label>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Amenidades
                </label>
                <div className="space-y-2">
                  {editData.amenities.map((amenity, index) => (
                      <div key={index} className="flex gap-2">
                        <input
                            type="text"
                            value={amenity}
                            onChange={(e) => handleAmenityChange(index, e.target.value)}
                            placeholder="Ej: WiFi, Piscina, Cocina..."
                            className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        />
                        <button
                            type="button"
                            onClick={() => handleRemoveAmenity(index)}
                            className="px-4 py-2 bg-red-100 text-red-600 rounded-lg hover:bg-red-200 transition"
                        >
                          <X size={20} />
                        </button>
                      </div>
                  ))}
                  <button
                      type="button"
                      onClick={handleAddAmenity}
                      className="w-full px-4 py-2 border-2 border-dashed border-gray-300 rounded-lg text-gray-600 hover:border-primary hover:text-primary transition"
                  >
                    + Agregar Amenidad
                  </button>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Imágenes (URLs)
                </label>
                <div className="space-y-2">
                  {editData.images.map((image, index) => (
                      <div key={index} className="flex gap-2">
                        <input
                            type="url"
                            value={image}
                            onChange={(e) => handleImageChange(index, e.target.value)}
                            placeholder="https://ejemplo.com/imagen.jpg"
                            className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                        />
                        <button
                            type="button"
                            onClick={() => handleRemoveImage(index)}
                            className="px-4 py-2 bg-red-100 text-red-600 rounded-lg hover:bg-red-200 transition"
                        >
                          <X size={20} />
                        </button>
                      </div>
                  ))}
                  <button
                      type="button"
                      onClick={handleAddImage}
                      className="w-full px-4 py-2 border-2 border-dashed border-gray-300 rounded-lg text-gray-600 hover:border-primary hover:text-primary transition"
                  >
                    + Agregar Imagen
                  </button>
                </div>
              </div>

              <div className="flex gap-4 pt-4 border-t">
                <button
                    type="button"
                    onClick={handleCloseEdit}
                    className="flex-1 px-6 py-3 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition"
                    disabled={editLoading}
                >
                  Cancelar
                </button>
                <button
                    type="submit"
                    className="flex-1 px-6 py-3 bg-primary text-white rounded-lg hover:bg-gray-800 transition disabled:opacity-50"
                    disabled={editLoading}
                >
                  {editLoading ? 'Guardando...' : 'Guardar Cambios'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default PropertyDetail;