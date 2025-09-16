// API configuration - automatically switches based on hostname
const API_BASE_URL = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
  ? 'http://localhost:8080'
  : 'https://cloudapi.vizn3r.eu';

console.log('API base URL:', API_BASE_URL);