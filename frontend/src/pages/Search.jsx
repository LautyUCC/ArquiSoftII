import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { propertiesAPI, uploadAPI } from '../services/api';
import { MapPin, Users, LogOut, Plus, X, Check, AlertCircle, Calendar, Trash2, ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-react';

function Search() {
  const navigate = useNavigate();
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalResults, setTotalResults] = useState(0);
  const pageSize = 6; // 6 propiedades por página
  
  // Estados para crear propiedad (solo ADMIN)
  const [isAdmin, setIsAdmin] = useState(false);
  const [showCreateProperty, setShowCreateProperty] = useState(false);
  const [createLoading, setCreateLoading] = useState(false);
  const [createError, setCreateError] = useState('');
  const [createSuccess, setCreateSuccess] = useState('');
  const [deleteConfirm, setDeleteConfirm] = useState(null);
  const [deleteLoading, setDeleteLoading] = useState(false);
  
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
    amenityInput: '', // Input temporal para agregar amenidades
    images: []
  });

  useEffect(() => {
    // Verificar si el usuario es ADMIN
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        console.log('Usuario desde localStorage:', user);
        console.log('Role del usuario:', user.role);
        console.log('Role en mayúsculas:', user.role?.toUpperCase());
        // Comparación robusta: convertir a mayúsculas para evitar problemas de case
        const isUserAdmin = user.role?.toUpperCase() === 'ADMIN';
        console.log('Es ADMIN?', isUserAdmin);
        setIsAdmin(isUserAdmin);
        // Si es admin, usar su ID como ownerId por defecto
        if (isUserAdmin && user.id) {
          setNewProperty(prev => ({ ...prev, ownerId: user.id.toString() }));
        }
      } catch (err) {
        console.error('Error parseando usuario:', err);
      }
    } else {
      console.log('No se encontró usuario en localStorage');
    }
    loadProperties('', 1);
  }, []);

  const loadProperties = async (query = '', page = 1) => {
    setLoading(true);
    setError('');

    try {
      console.log('🔍 Buscando propiedades con query:', query, 'página:', page);
      const response = await propertiesAPI.search({ 
        q: query,
        page: page,
        pageSize: pageSize
      });
      console.log('📦 Respuesta completa de búsqueda:', response);
      console.log('📦 response.data:', response.data);
      
      // La respuesta de search-api viene directamente como SearchResponse
      // { results: [...], totalResults: N, page: X, pageSize: Y, totalPages: Z }
      const results = response.data?.results || [];
      const total = response.data?.totalResults || 0;
      const totalPagesCount = response.data?.totalPages || 1;
      const currentPageNum = response.data?.page || 1;
      
      console.log(`✅ Encontradas ${total} propiedades (${results.length} en esta página):`, results);
      setProperties(results);
      setTotalResults(total);
      setTotalPages(totalPagesCount);
      setCurrentPage(currentPageNum);

      if (results.length === 0) {
        setError('No se encontraron propiedades');
      } else {
        setError(''); // Limpiar error si hay resultados
      }
    } catch (err) {
      console.error('❌ Error al cargar propiedades:', err);
      console.error('Detalles del error:', err.response?.data);

      if (err.code === 'ERR_NETWORK') {
        setError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setError(err.response?.data?.message || 'Error al cargar las propiedades. Por favor intenta de nuevo.');
      }
      setProperties([]);
      setTotalResults(0);
      setTotalPages(1);
      setCurrentPage(1);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    setCurrentPage(1); // Resetear a la primera página al buscar
    loadProperties(searchQuery, 1);
  };

  const handlePageChange = (newPage) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setCurrentPage(newPage);
      loadProperties(searchQuery, newPage);
      // Scroll al inicio de la página
      window.scrollTo({ top: 0, behavior: 'smooth' });
    }
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate('/');
  };

  // ============================================
  // CREAR PROPIEDAD (solo ADMIN)
  // ============================================

  const handleCreateProperty = async (e) => {
    e.preventDefault();
    setCreateError('');
    setCreateSuccess('');
    setCreateLoading(true);

    // Validaciones
    if (!newProperty.title || !newProperty.description || !newProperty.price || 
        !newProperty.location || !newProperty.ownerId || !newProperty.capacity) {
      setCreateError('Por favor completa todos los campos requeridos');
      setCreateLoading(false);
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
      setCreateSuccess('Propiedad creada exitosamente');
      
      // Ocultar mensaje de éxito después de 3 segundos
      setTimeout(() => {
        setCreateSuccess('');
      }, 3000);
      
      // Cerrar modal y limpiar formulario
      setShowCreateProperty(false);
      setNewProperty({
        title: '',
        description: '',
        price: '',
        location: '',
        ownerId: newProperty.ownerId, // Mantener el ownerId
        capacity: '',
        available: true,
        amenities: [],
        amenityInput: '',
        images: []
      });
      
      // Recargar propiedades después de un breve delay (dar tiempo a que se indexe en Solr)
      setTimeout(() => {
        loadProperties(searchQuery, currentPage);
      }, 2000);
    } catch (err) {
      console.error('Error al crear propiedad:', err);
      
      if (err.response?.status === 401) {
        setCreateError('Sesión expirada. Por favor inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 403) {
        setCreateError('No tienes permisos para crear propiedades. Se requiere rol ADMIN.');
      } else if (err.response?.status === 400) {
        setCreateError(err.response?.data?.message || err.response?.data?.error || 'Error al crear la propiedad');
      } else {
        setCreateError('Error al crear la propiedad. Por favor intenta de nuevo.');
      }
    } finally {
      setCreateLoading(false);
    }
  };

  const handleAddAmenity = () => {
    if (newProperty.amenityInput.trim() && !newProperty.amenities.includes(newProperty.amenityInput.trim())) {
      setNewProperty({
        ...newProperty,
        amenities: [...newProperty.amenities, newProperty.amenityInput.trim()],
        amenityInput: ''
      });
    }
  };

  const handleRemoveAmenity = (amenity) => {
    setNewProperty({
      ...newProperty,
      amenities: newProperty.amenities.filter(a => a !== amenity)
    });
  };

  const handleImageUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    // Validar tipo de archivo
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
      setCreateError('Formato de imagen no permitido. Solo se aceptan: JPG, JPEG, PNG, GIF, WEBP');
      return;
    }

    // Validar tamaño (máximo 5MB)
    if (file.size > 5 * 1024 * 1024) {
      setCreateError('El archivo es demasiado grande. Tamaño máximo: 5MB');
      return;
    }

    setCreateError('');
    setCreateLoading(true);

    try {
      const response = await uploadAPI.uploadImage(file);
      const imageUrl = response.data.url;
      
      if (!newProperty.images.includes(imageUrl)) {
        setNewProperty({
          ...newProperty,
          images: [...newProperty.images, imageUrl]
        });
      }
    } catch (err) {
      console.error('Error al subir imagen:', err);
      if (err.response?.status === 401) {
        setCreateError('Sesión expirada. Por favor inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 400) {
        setCreateError(err.response?.data?.message || 'Error al subir la imagen');
      } else {
        setCreateError('Error al subir la imagen. Por favor intenta de nuevo.');
      }
    } finally {
      setCreateLoading(false);
      // Limpiar el input para permitir subir el mismo archivo nuevamente
      e.target.value = '';
    }
  };

  const handleRemoveImage = (image) => {
    setNewProperty({
      ...newProperty,
      images: newProperty.images.filter(img => img !== image)
    });
  };

  // ============================================
  // ELIMINAR PROPIEDAD (solo ADMIN)
  // ============================================

  const handleDeleteProperty = async (propertyId) => {
    setDeleteLoading(true);
    setCreateError('');
    setCreateSuccess('');

    try {
      await propertiesAPI.delete(propertyId);
      setCreateSuccess('Propiedad eliminada exitosamente');
      
      setTimeout(() => {
        setCreateSuccess('');
      }, 3000);
      
      setDeleteConfirm(null);
      loadProperties(searchQuery, currentPage);
    } catch (err) {
      console.error('Error al eliminar propiedad:', err);
      
      if (err.response?.status === 401) {
        setCreateError('Sesión expirada. Por favor inicia sesión nuevamente.');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 2000);
      } else if (err.response?.status === 403) {
        setCreateError('No tienes permisos para eliminar propiedades. Se requiere rol ADMIN.');
      } else if (err.response?.status === 404) {
        setCreateError('Propiedad no encontrada');
      } else {
        setCreateError('Error al eliminar la propiedad. Por favor intenta de nuevo.');
      }
    } finally {
      setDeleteLoading(false);
    }
  };

  return (
      <div className="min-h-screen bg-secondary">
        {/* Header */}
        <header className="bg-white shadow-sm sticky top-0 z-50">
          <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
            <h1 className="text-2xl font-bold text-primary">Spotly</h1>
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
                  onClick={() => setShowCreateProperty(true)}
                  className="flex items-center gap-2 bg-primary text-white px-4 py-2 rounded-lg hover:bg-gray-800 transition"
                >
                  <Plus size={20} />
                  Crear Propiedad
                </button>
              )}
              <button
                  onClick={handleLogout}
                  className="flex items-center gap-2 text-gray-600 hover:text-primary transition"
              >
                <LogOut size={20} />
                Salir
              </button>
            </div>
          </div>
        </header>

        {/* Search Bar */}
        <div className="bg-white shadow-md">
          <div className="max-w-7xl mx-auto px-4 py-6">
            <form onSubmit={handleSearch} className="flex gap-4">
              <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Buscar propiedades por ubicación, título..."
                  className="flex-1 px-6 py-4 border border-gray-300 rounded-full focus:ring-2 focus:ring-primary focus:border-transparent text-lg"
              />
              <button
                  type="submit"
                  disabled={loading}
                  className="bg-primary text-white px-8 py-4 rounded-full font-medium hover:bg-gray-800 transition disabled:opacity-50"
              >
                {loading ? 'Buscando...' : 'Buscar'}
              </button>
            </form>
          </div>
        </div>

        {/* Content */}
        <div className="max-w-7xl mx-auto px-4 py-12">
          {/* Loading State */}
          {loading && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {[1, 2, 3].map((n) => (
                    <div key={n} className="bg-white rounded-2xl overflow-hidden shadow-lg animate-pulse">
                      <div className="aspect-[4/3] bg-gray-200"></div>
                      <div className="p-6">
                        <div className="h-6 bg-gray-200 rounded mb-4"></div>
                        <div className="h-4 bg-gray-200 rounded mb-2 w-2/3"></div>
                        <div className="h-4 bg-gray-200 rounded w-1/2"></div>
                      </div>
                    </div>
                ))}
              </div>
          )}

          {/* Error State */}
          {!loading && error && (
              <div className="text-center py-20">
                <p className="text-2xl text-gray-600 mb-4">{error}</p>
                {properties.length === 0 && error === 'No se encontraron propiedades' && (
                    <button
                        onClick={() => {
                          setCurrentPage(1);
                          loadProperties('', 1);
                        }}
                        className="text-primary hover:underline font-medium"
                    >
                      Ver todas las propiedades
                    </button>
                )}
              </div>
          )}

          {/* Properties Grid */}
          {!loading && !error && properties.length > 0 && (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {properties.map((property) => (
                    <div
                        key={property.id}
                        className="bg-white rounded-2xl overflow-hidden shadow-lg hover:shadow-2xl transition group relative"
                    >
                      {/* Image */}
                      <div 
                        className="aspect-[4/3] bg-gray-200 relative overflow-hidden cursor-pointer"
                        onClick={() => navigate(`/property/${property.id}`)}
                      >
                        {property.images && property.images.length > 0 && property.images[0] ? (
                            <img
                                src={property.images[0]}
                                alt={property.title}
                                className="w-full h-full object-cover group-hover:scale-105 transition duration-300"
                            />
                        ) : (
                            <div className="w-full h-full flex items-center justify-center text-gray-400">
                              <MapPin size={48} />
                            </div>
                        )}
                        {/* Botón eliminar (solo ADMIN) */}
                        {isAdmin && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              setDeleteConfirm(property);
                            }}
                            className="absolute top-2 right-2 bg-red-600 text-white p-2 rounded-full opacity-0 group-hover:opacity-100 transition hover:bg-red-700 z-10"
                            title="Eliminar propiedad"
                          >
                            <Trash2 size={18} />
                          </button>
                        )}
                      </div>

                      {/* Info */}
                      <div className="p-6 cursor-pointer" onClick={() => navigate(`/property/${property.id}`)}>
                        <h3 className="text-xl font-bold text-primary mb-2 group-hover:text-gray-800 transition">
                          {property.title}
                        </h3>

                        <div className="flex items-center text-gray-600 mb-3">
                          <MapPin size={16} className="mr-1"/>
                          <span className="text-sm">{property.city}, {property.country}</span>
                        </div>

                        <p className="text-gray-600 text-sm mb-4 line-clamp-2">
                          {property.description}
                        </p>

                        <div className="flex items-center justify-between">
                          <div className="flex items-center text-gray-600 text-sm">
                            <Users size={16} className="mr-1"/>
                            <span>{property.maxGuests} huéspedes</span>
                          </div>
                          <div>
                      <span className="text-2xl font-bold text-primary">
                        ${Math.round(property.pricePerNight)}
                      </span>
                            <span className="text-sm text-gray-600"> / noche</span>
                          </div>
                        </div>
                      </div>
                    </div>
                ))}
              </div>
          )}

          {/* Paginación */}
          {!loading && !error && totalPages > 1 && (
            <div className="mt-8 flex flex-col items-center gap-4">
              {/* Información de paginación */}
              <div className="text-sm text-gray-600">
                Mostrando {properties.length} de {totalResults} propiedades (Página {currentPage} de {totalPages})
              </div>

              {/* Controles de paginación */}
              <div className="flex items-center gap-2">
                {/* Botón Primera Página */}
                <button
                  onClick={() => handlePageChange(1)}
                  disabled={currentPage === 1}
                  className="p-2 rounded-lg border border-gray-300 hover:bg-gray-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                  title="Primera página"
                >
                  <ChevronsLeft size={20} />
                </button>

                {/* Botón Página Anterior */}
                <button
                  onClick={() => handlePageChange(currentPage - 1)}
                  disabled={currentPage === 1}
                  className="p-2 rounded-lg border border-gray-300 hover:bg-gray-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                  title="Página anterior"
                >
                  <ChevronLeft size={20} />
                </button>

                {/* Números de página */}
                <div className="flex items-center gap-1">
                  {(() => {
                    const pages = [];
                    const maxVisible = 5; // Máximo de números de página visibles
                    let startPage = Math.max(1, currentPage - Math.floor(maxVisible / 2));
                    let endPage = Math.min(totalPages, startPage + maxVisible - 1);

                    // Ajustar si estamos cerca del final
                    if (endPage - startPage < maxVisible - 1) {
                      startPage = Math.max(1, endPage - maxVisible + 1);
                    }

                    // Mostrar primera página si no está visible
                    if (startPage > 1) {
                      pages.push(
                        <button
                          key={1}
                          onClick={() => handlePageChange(1)}
                          className="px-3 py-1 rounded-lg border border-gray-300 hover:bg-gray-50 transition"
                        >
                          1
                        </button>
                      );
                      if (startPage > 2) {
                        pages.push(
                          <span key="ellipsis-start" className="px-2 text-gray-400">
                            ...
                          </span>
                        );
                      }
                    }

                    // Mostrar páginas visibles
                    for (let i = startPage; i <= endPage; i++) {
                      pages.push(
                        <button
                          key={i}
                          onClick={() => handlePageChange(i)}
                          className={`px-3 py-1 rounded-lg border transition ${
                            i === currentPage
                              ? 'bg-primary text-white border-primary'
                              : 'border-gray-300 hover:bg-gray-50'
                          }`}
                        >
                          {i}
                        </button>
                      );
                    }

                    // Mostrar última página si no está visible
                    if (endPage < totalPages) {
                      if (endPage < totalPages - 1) {
                        pages.push(
                          <span key="ellipsis-end" className="px-2 text-gray-400">
                            ...
                          </span>
                        );
                      }
                      pages.push(
                        <button
                          key={totalPages}
                          onClick={() => handlePageChange(totalPages)}
                          className="px-3 py-1 rounded-lg border border-gray-300 hover:bg-gray-50 transition"
                        >
                          {totalPages}
                        </button>
                      );
                    }

                    return pages;
                  })()}
                </div>

                {/* Botón Página Siguiente */}
                <button
                  onClick={() => handlePageChange(currentPage + 1)}
                  disabled={currentPage === totalPages}
                  className="p-2 rounded-lg border border-gray-300 hover:bg-gray-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                  title="Página siguiente"
                >
                  <ChevronRight size={20} />
                </button>

                {/* Botón Última Página */}
                <button
                  onClick={() => handlePageChange(totalPages)}
                  disabled={currentPage === totalPages}
                  className="p-2 rounded-lg border border-gray-300 hover:bg-gray-50 transition disabled:opacity-50 disabled:cursor-not-allowed"
                  title="Última página"
                >
                  <ChevronsRight size={20} />
                </button>
              </div>
            </div>
          )}
        </div>

        {/* Modal Crear Propiedad (solo ADMIN) */}
        {showCreateProperty && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-2xl shadow-xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
              <div className="p-6 border-b border-gray-200 flex items-center justify-between sticky top-0 bg-white">
                <h3 className="text-xl font-bold text-primary">Crear Nueva Propiedad</h3>
                <button
                  onClick={() => {
                    setShowCreateProperty(false);
                    setCreateError('');
                    setCreateSuccess('');
                    setNewProperty({
                      title: '',
                      description: '',
                      price: '',
                      location: '',
                      ownerId: newProperty.ownerId,
                      capacity: '',
                      available: true,
                      amenities: [],
                      amenityInput: '',
                      images: [],
                      imageInput: ''
                    });
                  }}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <X size={24} />
                </button>
              </div>

              <form onSubmit={handleCreateProperty} className="p-6 space-y-4">
                {/* Mensajes de éxito/error */}
                {createError && (
                  <div className="p-4 bg-red-50 border border-red-200 rounded-lg flex items-center gap-2">
                    <AlertCircle size={20} className="text-red-600" />
                    <p className="text-red-600">{createError}</p>
                  </div>
                )}

                {createSuccess && (
                  <div className="p-4 bg-green-50 border border-green-200 rounded-lg flex items-center gap-2">
                    <Check size={20} className="text-green-600" />
                    <p className="text-green-600">{createSuccess}</p>
                  </div>
                )}

                {/* Título */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Título * <span className="text-gray-500 text-xs">(Nombre de la propiedad)</span>
                  </label>
                  <input
                    type="text"
                    value={newProperty.title}
                    onChange={(e) => setNewProperty({...newProperty, title: e.target.value})}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    placeholder="Ej: Apartamento moderno en el centro"
                    required
                  />
                </div>

                {/* Descripción */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Descripción * <span className="text-gray-500 text-xs">(Detalles de la propiedad)</span>
                  </label>
                  <textarea
                    value={newProperty.description}
                    onChange={(e) => setNewProperty({...newProperty, description: e.target.value})}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    rows="4"
                    placeholder="Describe la propiedad, sus características, ubicación, etc."
                    required
                  />
                </div>

                {/* Precio y Capacidad */}
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Precio por Noche * <span className="text-gray-500 text-xs">(USD)</span>
                    </label>
                    <input
                      type="number"
                      step="0.01"
                      min="0"
                      value={newProperty.price}
                      onChange={(e) => setNewProperty({...newProperty, price: e.target.value})}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      placeholder="120000.00"
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Capacidad * <span className="text-gray-500 text-xs">(Huéspedes)</span>
                    </label>
                    <input
                      type="number"
                      min="1"
                      value={newProperty.capacity}
                      onChange={(e) => setNewProperty({...newProperty, capacity: e.target.value})}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      placeholder="4"
                      required
                    />
                  </div>
                </div>

                {/* Ubicación */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Ubicación * <span className="text-gray-500 text-xs">(Ciudad, País)</span>
                  </label>
                  <input
                    type="text"
                    value={newProperty.location}
                    onChange={(e) => setNewProperty({...newProperty, location: e.target.value})}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    placeholder="Ej: Bogotá, Colombia"
                    required
                  />
                </div>

                {/* Owner ID */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Owner ID * <span className="text-gray-500 text-xs">(ID del propietario)</span>
                  </label>
                  <input
                    type="text"
                    value={newProperty.ownerId}
                    onChange={(e) => setNewProperty({...newProperty, ownerId: e.target.value})}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                    placeholder="1"
                    required
                  />
                  <p className="text-xs text-gray-500 mt-1">Por defecto se usa tu ID de usuario</p>
                </div>

                {/* Amenidades */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Amenidades <span className="text-gray-500 text-xs">(Opcional)</span>
                  </label>
                  <div className="flex gap-2 mb-2">
                    <input
                      type="text"
                      value={newProperty.amenityInput}
                      onChange={(e) => setNewProperty({...newProperty, amenityInput: e.target.value})}
                      onKeyPress={(e) => {
                        if (e.key === 'Enter') {
                          e.preventDefault();
                          handleAddAmenity();
                        }
                      }}
                      className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
                      placeholder="Ej: wifi, pool, parking, kitchen"
                    />
                    <button
                      type="button"
                      onClick={handleAddAmenity}
                      className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition"
                    >
                      Agregar
                    </button>
                  </div>
                  {newProperty.amenities.length > 0 && (
                    <div className="flex flex-wrap gap-2">
                      {newProperty.amenities.map((amenity, index) => (
                        <span
                          key={index}
                          className="inline-flex items-center gap-1 px-3 py-1 bg-primary/10 text-primary rounded-full text-sm"
                        >
                          {amenity}
                          <button
                            type="button"
                            onClick={() => handleRemoveAmenity(amenity)}
                            className="text-primary hover:text-red-600"
                          >
                            <X size={14} />
                          </button>
                        </span>
                      ))}
                    </div>
                  )}
                </div>

                {/* Imágenes */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Imágenes <span className="text-gray-500 text-xs">(Opcional - JPG, PNG, GIF, WEBP, máx. 5MB)</span>
                  </label>
                  <div className="mb-2">
                    <input
                      type="file"
                      accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
                      onChange={handleImageUpload}
                      disabled={createLoading}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent disabled:opacity-50 disabled:cursor-not-allowed"
                    />
                    <p className="mt-1 text-xs text-gray-500">
                      Selecciona una imagen desde tu computadora. Puedes agregar múltiples imágenes.
                    </p>
                  </div>
                  {newProperty.images.length > 0 && (
                    <div className="space-y-2 mt-4">
                      <p className="text-sm font-medium text-gray-700">Imágenes agregadas ({newProperty.images.length}):</p>
                      <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
                        {newProperty.images.map((image, index) => (
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
                              onClick={() => handleRemoveImage(image)}
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

                {/* Disponible */}
                <div>
                  <label className="flex items-center gap-2">
                    <input
                      type="checkbox"
                      checked={newProperty.available}
                      onChange={(e) => setNewProperty({...newProperty, available: e.target.checked})}
                      className="w-4 h-4 text-primary rounded focus:ring-primary"
                    />
                    <span className="text-sm text-gray-700">Propiedad disponible para reservas</span>
                  </label>
                </div>

                {/* Botones */}
                <div className="flex gap-4 pt-4 border-t border-gray-200">
                  <button
                    type="submit"
                    disabled={createLoading}
                    className="flex-1 bg-primary text-white py-3 rounded-lg hover:bg-gray-800 transition disabled:opacity-50 disabled:cursor-not-allowed font-medium"
                  >
                    {createLoading ? 'Creando...' : 'Crear Propiedad'}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setShowCreateProperty(false);
                      setCreateError('');
                      setCreateSuccess('');
                      setNewProperty({
                        title: '',
                        description: '',
                        price: '',
                        location: '',
                        ownerId: newProperty.ownerId,
                        capacity: '',
                        available: true,
                        amenities: [],
                        amenityInput: '',
                        images: [],
                        imageInput: ''
                      });
                    }}
                    className="flex-1 bg-gray-200 text-gray-700 py-3 rounded-lg hover:bg-gray-300 transition font-medium"
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
                  disabled={deleteLoading}
                  className="flex-1 bg-red-600 text-white py-2 rounded-lg hover:bg-red-700 transition disabled:opacity-50"
                >
                  {deleteLoading ? 'Eliminando...' : 'Eliminar'}
                </button>
                <button
                  onClick={() => setDeleteConfirm(null)}
                  disabled={deleteLoading}
                  className="flex-1 bg-gray-200 text-gray-700 py-2 rounded-lg hover:bg-gray-300 transition disabled:opacity-50"
                >
                  Cancelar
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    );
  }

  export default Search;