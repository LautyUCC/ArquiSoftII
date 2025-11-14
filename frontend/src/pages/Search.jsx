import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { propertiesAPI } from '../services/api';
import { MapPin, Users, LogOut } from 'lucide-react';

function Search() {
  const navigate = useNavigate();
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    loadProperties();
  }, []);

  const loadProperties = async (query = '') => {
    setLoading(true);
    setError('');

    try {
      const response = await propertiesAPI.search({ q: query });
      setProperties(response.data.results || []);

      if (response.data.results.length === 0) {
        setError('No se encontraron propiedades');
      }
    } catch (err) {
      console.error('Error al cargar propiedades:', err);

      if (err.code === 'ERR_NETWORK') {
        setError('Error de conexión. Verifica que el servidor esté activo.');
      } else {
        setError('Error al cargar las propiedades. Por favor intenta de nuevo.');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    loadProperties(searchQuery);
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate('/');
  };

  return (
      <div className="min-h-screen bg-secondary">
        {/* Header */}
        <header className="bg-white shadow-sm sticky top-0 z-50">
          <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
            <h1 className="text-2xl font-bold text-primary">Spotly</h1>
            <button
                onClick={handleLogout}
                className="flex items-center gap-2 text-gray-600 hover:text-primary transition"
            >
              <LogOut size={20} />
              Salir
            </button>
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
                        onClick={() => loadProperties('')}
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
                        onClick={() => navigate(`/property/${property.id}`)}
                        className="bg-white rounded-2xl overflow-hidden shadow-lg hover:shadow-2xl transition cursor-pointer group"
                    >
                      {/* Image */}
                      <div className="aspect-[4/3] bg-gray-200 relative overflow-hidden">
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
                      </div>

                      {/* Info */}
                      <div className="p-6">
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
        </div>
      </div>
  );
}

export default Search;