import { useNavigate, useLocation } from 'react-router-dom';
import { CheckCircle, Home, Calendar, Users, DollarSign } from 'lucide-react';
import { useEffect } from 'react';

function Congrats() {
  const navigate = useNavigate();
  const location = useLocation();
  const bookingData = location.state;

  // Proteger contra acceso directo
  useEffect(() => {
    if (!bookingData) {
      navigate('/search');
    }
  }, [bookingData, navigate]);

  if (!bookingData) {
    return null;
  }

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', { 
      day: 'numeric', 
      month: 'long', 
      year: 'numeric' 
    });
  };

  return (
    <div className="min-h-screen bg-secondary flex items-center justify-center px-4">
      <div className="max-w-2xl w-full bg-white rounded-2xl shadow-xl p-8">
        <div className="text-center mb-8">
          <CheckCircle size={80} className="mx-auto text-green-500 mb-4" />
          <h1 className="text-4xl font-bold text-primary mb-2">
            ¡Reserva Exitosa!
          </h1>
          <p className="text-gray-600">
            Tu reserva ha sido confirmada exitosamente
          </p>
        </div>

        {/* Resumen de Reserva */}
        <div className="bg-gray-50 rounded-xl p-6 mb-8 space-y-4">
          <h2 className="text-xl font-bold text-primary mb-4">
            Resumen de tu reserva
          </h2>

          <div className="flex items-start gap-3">
            <Home size={20} className="text-gray-600 mt-1" />
            <div>
              <p className="font-medium text-gray-900">{bookingData.property}</p>
            </div>
          </div>

          <div className="flex items-start gap-3">
            <Calendar size={20} className="text-gray-600 mt-1" />
            <div>
              <p className="text-gray-700">
                <span className="font-medium">Entrada:</span> {formatDate(bookingData.checkIn)}
              </p>
              <p className="text-gray-700">
                <span className="font-medium">Salida:</span> {formatDate(bookingData.checkOut)}
              </p>
              <p className="text-sm text-gray-500 mt-1">
                {bookingData.nights} {bookingData.nights === 1 ? 'noche' : 'noches'}
              </p>
            </div>
          </div>

          <div className="flex items-start gap-3">
            <Users size={20} className="text-gray-600 mt-1" />
            <div>
              <p className="text-gray-700">
                {bookingData.guests} {bookingData.guests === 1 ? 'huésped' : 'huéspedes'}
              </p>
            </div>
          </div>

          <div className="pt-4 border-t border-gray-200">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <DollarSign size={20} className="text-gray-600" />
                <span className="font-medium text-gray-900">Total:</span>
              </div>
              <span className="text-2xl font-bold text-primary">
                ${bookingData.totalPrice}
              </span>
            </div>
          </div>
        </div>

        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
          <p className="text-sm text-blue-800">
            📧 Recibirás un email de confirmación con todos los detalles de tu reserva.
          </p>
        </div>

        <button
          onClick={() => navigate('/search')}
          className="w-full bg-primary text-white py-4 rounded-lg font-medium hover:bg-gray-800 transition flex items-center justify-center gap-2"
        >
          <Home size={20} />
          Volver al inicio
        </button>
      </div>
    </div>
  );
}

export default Congrats;
