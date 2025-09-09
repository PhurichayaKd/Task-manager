const API_BASE = 'http://localhost:8080';

// Show error message
function showError(message) {
  const errorDiv = document.createElement('div');
  errorDiv.className = 'error-message';
  errorDiv.innerHTML = `<p class="error-text">${message}</p>`;
  
  const form = document.getElementById('form-google-signup');
  const existingError = form.querySelector('.error-message');
  if (existingError) {
    existingError.remove();
  }
  
  form.insertBefore(errorDiv, form.querySelector('button'));
  
  setTimeout(() => {
    errorDiv.remove();
  }, 5000);
}

// Get user info from URL parameters or token
function getUserInfoFromToken() {
  const token = localStorage.getItem('access_token');
  if (!token) {
    window.location.href = './login.html';
    return null;
  }
  
  // Decode JWT token to get user info
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return {
      email: payload.email || '',
      name: payload.name || '',
      avatar: payload.avatar || ''
    };
  } catch (e) {
    console.error('Error decoding token:', e);
    return null;
  }
}

// Display user info
function displayUserInfo(userInfo) {
  const emailElement = document.getElementById('user-email');
  const avatarElement = document.getElementById('user-avatar');
  const nameInput = document.getElementById('display-name');
  
  if (userInfo.email) {
    emailElement.textContent = userInfo.email;
  }
  
  if (userInfo.avatar) {
    avatarElement.src = userInfo.avatar;
    avatarElement.style.display = 'block';
  }
  
  if (userInfo.name) {
    nameInput.value = userInfo.name;
  }
}

// Handle form submission
function handleGoogleSignupComplete(event) {
  event.preventDefault();
  
  const displayName = document.getElementById('display-name').value.trim();
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value.trim();
  const token = localStorage.getItem('access_token');
  
  if (!displayName || !username || !password) {
    showError('กรุณากรอกข้อมูลให้ครบถ้วน');
    return;
  }
  
  if (password.length < 6) {
    showError('รหัสผ่านต้องมีอย่างน้อย 6 ตัวอักษร');
    return;
  }
  
  // Update user profile
  fetch(`${API_BASE}/api/users/profile`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      name: displayName,
      username: username,
      password: password
    })
  })
  .then(response => response.json())
  .then(data => {
    if (data.success) {
      // Redirect to dashboard
      window.location.href = '/home.html';
    } else {
      showError(data.error || 'เกิดข้อผิดพลาดในการอัปเดตข้อมูล');
    }
  })
  .catch(error => {
    console.error('Update profile error:', error);
    showError('เกิดข้อผิดพลาดในการเชื่อมต่อ');
  });
}

// Initialize page
document.addEventListener('DOMContentLoaded', function() {
  // Check if user is authenticated
  const token = localStorage.getItem('access_token');
  if (!token) {
    window.location.href = './login.html';
    return;
  }
  
  // Get and display user info
  const userInfo = getUserInfoFromToken();
  if (userInfo) {
    displayUserInfo(userInfo);
  }
  
  // Handle form submission
  const form = document.getElementById('form-google-signup');
  if (form) {
    form.addEventListener('submit', handleGoogleSignupComplete);
  }
  
  const completeBtn = document.getElementById('btn-complete');
  if (completeBtn) {
    completeBtn.addEventListener('click', handleGoogleSignupComplete);
  }
});