import { useNavigate } from 'react-router-dom';
import { CheckCircle, Home } from 'lucide-react';

function Congrats() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-secondary flex items-center justify-center px-4">
      <div className="max-w-md w-full bg-white rounded-2xl shadow-xl p-8 text-center">
        <div className="mb-6">
          <CheckCircle size={80} className="mx-auto text-green-500" />
        </div>

        <h1 className="text-3xl font-bold text-primary mb-4">
          ¡Reserva Exitosa!
        </h1>

        <p className="text-gray-600 mb-8 leading-relaxed">
          Tu reserva ha sido confirmada exitosamente. Recibirás un email con todos los detalles.
        </p>

        <button
          onClick={() => navigate('/search')}
          className="w-full bg-primary text-white py-3 rounded-lg font-medium hover:bg-gray-800 transition flex items-center justify-center gap-2"
        >
          <Home size={20} />
          Volver al inicio
        </button>
      </div>
    </div>
  );
}

export default Congrats;
