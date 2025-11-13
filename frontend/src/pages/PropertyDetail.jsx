import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { propertiesAPI } from '../services/api';
import { MapPin, Users, Bed, Bath, ArrowLeft, Check } from 'lucide-react';

function PropertyDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [property, setProperty] = useState(null);
  const [loading, setLoading] = useState(true);
  const [booking, setBooking] = useState(false);

  useEffect(() => {
    loadProperty();
  }, [id]);

  const loadProperty = async () => {
    try {
      const response = await propertiesAPI.getById(id);
      setProperty(response.data.data);
    } catch (err) {
      console.error('Error al cargar propiedad:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleBooking = async () => {
    setBooking(true);
    // Simular reserva (aquí irá la llamada real a la API)
    setTimeout(() => {
      setBooking(false);
      navigate('/congrats');
    }, 1500);
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
        <div className="max-w-7xl mx-auto px-4 py-4">
          <button
            onClick={() => navigate('/search')}
            className="flex items-center gap-2 text-gray-600 hover:text-primary transition"
          >
            <ArrowLeft size={20} />
            Volver a búsqueda
          </button>
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
                  <span className="text-lg">{property.city}, {property.country}</span>
                </div>

                <div className="flex gap-6 mb-8 pb-8 border-b">
                  <div className="flex items-center gap-2">
                    <Users size={20} className="text-gray-600" />
                    <span>{property.maxGuests} huéspedes</span>
                  </div>
                  {property.bedrooms > 0 && (
                    <div className="flex items-center gap-2">
                      <Bed size={20} className="text-gray-600" />
                      <span>{property.bedrooms} habitaciones</span>
                    </div>
                  )}
                  {property.bathrooms > 0 && (
                    <div className="flex items-center gap-2">
                      <Bath size={20} className="text-gray-600" />
                      <span>{property.bathrooms} baños</span>
                    </div>
                  )}
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
                      ${property.pricePerNight}
                    </p>
                    <p className="text-gray-600">por noche</p>
                  </div>

                  <button
                    onClick={handleBooking}
                    disabled={booking}
                    className="w-full bg-primary text-white py-4 rounded-lg font-bold text-lg hover:bg-gray-800 transition disabled:opacity-50"
                  >
                    {booking ? 'Reservando...' : 'Reservar'}
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
    </div>
  );
}

export default PropertyDetail;
