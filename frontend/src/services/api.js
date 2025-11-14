import axios from 'axios';

const API_BASE_URL = 'http://localhost:8081';

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

// Auth endpoints
export const authAPI = {
  login: (usernameOrEmail, password) => 
    api.post('/users/login', { username_or_email: usernameOrEmail, password: password }),
  
  register: (userData) => 
    api.post('/users', userData),
};

// Properties endpoints
export const propertiesAPI = {
  search: (params) => 
    axios.get('http://localhost:8083/search', { params }),
  
  getById: (id) => 
    axios.get(`http://localhost:8082/api/properties/${id}`),
  
  create: (propertyData) => 
    axios.post('http://localhost:8082/api/properties', propertyData),
};

export default api;
