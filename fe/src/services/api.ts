import axios from 'axios';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('accessToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => {
    // Handle the new response format: { result: { code, message }, data }
    if (response.data && response.data.result) {
      // If code is not 0, treat as error
      if (response.data.result.code !== 0) {
        return Promise.reject({
          response: {
            data: {
              message: response.data.result.message,
              code: response.data.result.code
            }
          }
        });
      }
      // Return only the data part for success responses
      return { ...response, data: response.data.data };
    }
    return response;
  },
  (error) => {
    // Handle error responses
    if (error.response?.status === 401) {
      // Handle unauthorized access
      localStorage.removeItem('accessToken');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    
    // Extract error message from the new format
    if (error.response?.data?.result?.message) {
      error.message = error.response.data.result.message;
    }
    
    return Promise.reject(error);
  }
);

export default api;