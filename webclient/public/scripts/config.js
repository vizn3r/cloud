// API configuration - uses same hostname as the app with /v1 prefix
const API_BASE_URL = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
  ? 'http://localhost:8080/v1'
  : `${window.location.protocol}//${window.location.hostname}/v1`;

console.log('API base URL:', API_BASE_URL);