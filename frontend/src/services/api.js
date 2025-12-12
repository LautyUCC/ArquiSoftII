import axios from 'axios';

// Usar rutas relativas para que pasen por nginx
// Nginx hará proxy a las APIs correspondientes
const API_BASE_URL = '/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para agregar token JWT
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth endpoints - /api/users/* → users-api
export const authAPI = {
  login: (usernameOrEmail, password) =>
      api.post('/users/login', { usernameOrEmail: usernameOrEmail, password: password }),

  register: (userData) =>
      api.post('/users', userData),
};

// Admin endpoints - /api/admin/* → users-api (solo ADMIN)
export const adminAPI = {
  // Obtener todos los usuarios (solo ADMIN)
  getAllUsers: () =>
      api.get('/admin/users'),

  // Actualizar usuario, incluyendo role (solo ADMIN)
  updateUser: (userId, userData) =>
      api.put(`/admin/users/${userId}`, userData),
};

// Properties endpoints - /api/properties/* → properties-api
// Search endpoint - /api/search → search-api
export const propertiesAPI = {
  search: ({ q, ...params }) =>
      api.get('/search', {
        params: { query: q, ...params }
      }),

  getById: (id) =>
      api.get(`/properties/${id}`),

  create: (propertyData) =>
      api.post('/properties', propertyData),

  update: (id, propertyData) =>
      api.put(`/properties/${id}`, propertyData),

  delete: (id) =>
      api.delete(`/properties/${id}`),

  // Obtener todas las propiedades (solo ADMIN)
  getAll: () =>
      api.get('/admin/properties'),
};

// Upload endpoints - /api/upload → properties-api
export const uploadAPI = {
  // Subir una imagen (cualquier usuario autenticado)
  // FormData con campo "image"
  uploadImage: (file) => {
    const formData = new FormData();
    formData.append('image', file);
    
    const token = localStorage.getItem('token');
    return axios.post(`${API_BASE_URL}/upload`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': token ? `Bearer ${token}` : '',
      },
    });
  },
};

// Bookings endpoints - /api/bookings/* → properties-api
export const bookingsAPI = {
  // Crear una nueva reserva (cualquier usuario autenticado)
  // Body: { propertyId, startDate, endDate } o { propertyId, checkIn, checkOut }
  createBooking: (bookingData) =>
      api.post('/bookings', bookingData),

  // Obtener mis reservas (cualquier usuario autenticado)
  getMyBookings: () =>
      api.get('/bookings/my'),

  // Obtener reservas de una propiedad (cualquier usuario autenticado, solo CONFIRMED)
  getPropertyBookings: (propertyId) =>
      api.get(`/bookings/property/${propertyId}`),

  // Actualizar fechas de una reserva (PATCH)
  // Body: { startDate, endDate } o { checkIn, checkOut }
  updateBooking: (bookingId, bookingData) =>
      api.patch(`/bookings/${bookingId}`, bookingData),

  // Cancelar una reserva (DELETE - marca como CANCELLED)
  deleteBooking: (bookingId) =>
      api.delete(`/bookings/${bookingId}`),

  // Obtener todas las reservas (solo ADMIN)
  getAllBookings: () =>
      api.get('/admin/bookings'),
};

export default api;