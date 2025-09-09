// Utility functions for authentication

/**
 * Get cookie value by name
 * @param {string} name - Cookie name
 * @returns {string|null} - Cookie value or null if not found
 */
function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    return parts.pop().split(';').shift();
  }
  return null;
}

/**
 * Check if user is authenticated by looking for token in localStorage or cookie
 * If token exists in cookie but not in localStorage, copy it to localStorage
 * @returns {string|null} - Token if authenticated, null otherwise
 */
function getAuthToken() {
  // First check localStorage
  let token = localStorage.getItem('access_token');
  
  if (!token) {
    // Check cookie if localStorage doesn't have token
    token = getCookie('access_token');
    if (token) {
      // Copy token from cookie to localStorage for consistency
      localStorage.setItem('access_token', token);
    }
  }
  
  return token;
}

/**
 * Redirect to login page if not authenticated
 * @param {string} redirectPath - Path to redirect to if not authenticated (default: '../index.html')
 */
function requireAuth(redirectPath = '../index.html') {
  const token = getAuthToken();
  if (!token) {
    location.replace(redirectPath);
  }
  return token;
}

/**
 * Clear authentication data and redirect to login
 * @param {string} redirectPath - Path to redirect to (default: '../index.html')
 */
function logout(redirectPath = '../index.html') {
  localStorage.removeItem('access_token');
  // Clear cookie by setting it to expire
  document.cookie = 'access_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  location.replace(redirectPath);
}

// Export functions for use in other scripts
window.AuthUtils = {
  getCookie,
  getAuthToken,
  requireAuth,
  logout
};