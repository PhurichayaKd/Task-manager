// API Base URL - ใช้ port 8080 สำหรับ Docker API container
const API_BASE = 'https://task-manager-production-6c61.up.railway.app';

// Helper function to show error messages
function showError(message) {
  const errorDiv = document.getElementById('error-message');
  const errorText = errorDiv?.querySelector('.error-text');
  if (errorDiv && errorText) {
    errorText.textContent = message;
    errorDiv.style.display = 'block';
    // Hide error after 5 seconds
    setTimeout(() => {
      errorDiv.style.display = 'none';
    }, 5000);
  } else {
    alert(message);
  }
}

// Global function for Google login
window.handleGoogleLogin = function() {
  console.log('handleGoogleLogin called');
  console.log('API_BASE:', API_BASE);
  console.log('Current location:', window.location.href);
  
  try {
    const nextUrl = ''; // ไม่กำหนด next URL เพื่อให้ใช้ onboardIfNew logic
    const onboardIfNew = '1'; // ส่งผู้ใช้ใหม่ไปหน้า profile
    const timestamp = Date.now(); // เพิ่ม timestamp เพื่อป้องกัน cache
    const fullUrl = `${API_BASE}/api/auth/google/login?onboardIfNew=${onboardIfNew}&t=${timestamp}&cache=${Math.random()}`;
    
    console.log('Redirecting to:', fullUrl);
    
    const healthUrl = `${API_BASE}/healthz?t=${Date.now()}`;
    console.log('Testing health endpoint:', healthUrl);
    
    // ตรวจสอบว่า API พร้อมใช้งานก่อน redirect
    fetch(healthUrl, {
      method: 'GET',
      cache: 'no-cache',
      headers: {
        'Cache-Control': 'no-cache'
      }
    })
      .then(response => {
        console.log('Health check response status:', response.status);
        console.log('Health check response ok:', response.ok);
        if (response.ok) {
          console.log('Health check passed, redirecting to Google login...');
          window.location.href = fullUrl;
        } else {
          throw new Error(`API not available - status: ${response.status}`);
        }
      })
      .catch(error => {
        console.error('API connection error details:', error);
        console.error('Error type:', typeof error);
        console.error('Error message:', error.message);
        alert('ไม่สามารถเชื่อมต่อกับเซิร์ฟเวอร์ได้ กรุณาลองใหม่อีกครั้ง');
      });
      
  } catch (error) {
    console.error('Google login error:', error);
    alert('เกิดข้อผิดพลาดในการเข้าสู่ระบบ กรุณาลองใหม่อีกครั้ง');
  }
}

// DOMContentLoaded event listener
document.addEventListener('DOMContentLoaded', function() {
  // Check if user is already authenticated (for create account page)
if (window.location.pathname.includes('create_account.html')) {
    // Wait for AuthUtils to load
    setTimeout(() => {
      if (window.AuthUtils) {
        const token = window.AuthUtils.getAuthToken();
        if (token) {
          // User is already authenticated, redirect to dashboard
          window.location.href = '/home.html';
          return;
        }
      }
    }, 100);
  }
  
  // Handle Continue button in signup page
  const btnContinue = document.getElementById('btn-continue');
  if (btnContinue) {
    btnContinue.addEventListener('click', function() {
      const email = document.getElementById('signup-email')?.value.trim();
      if (!email) {
        alert('กรุณากรอกอีเมล');
        return;
      }
      
      // Validate email format
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(email)) {
        alert('รูปแบบอีเมลไม่ถูกต้อง');
        return;
      }
      
      // Store email and redirect to create account page
      sessionStorage.setItem('tm_signup_email', email);
      window.location.href = '../create_account.html?email=' + encodeURIComponent(email);
    });
  }
  
  // Handle Login button in login page
  const btnLogin = document.getElementById('btn-login');
  if (btnLogin) {
    btnLogin.addEventListener('click', function() {
       const usernameOrEmail = document.getElementById('login-username')?.value.trim();
       const password = document.getElementById('login-password')?.value.trim();
       
       if (!usernameOrEmail) {
         showError('กรุณากรอกชื่อผู้ใช้หรืออีเมล');
         return;
       }
       
       if (!password) {
         showError('กรุณากรอกรหัสผ่าน');
         return;
       }
       
       // Call login API directly
       fetch(`${API_BASE}/api/auth/login`, {
         method: 'POST',
         headers: {
           'Content-Type': 'application/json'
         },
         body: JSON.stringify({
           email: usernameOrEmail,
           password: password
         })
       })
       .then(response => response.json())
       .then(data => {
         if (data.success && data.token) {
           // Store token in localStorage
           localStorage.setItem('access_token', data.token);
           // Redirect to dashboard
           window.location.href = '/home.html';
         } else {
           showError('ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง');
         }
       })
       .catch(error => {
         console.error('Login error:', error);
         showError('เกิดข้อผิดพลาดในการเข้าสู่ระบบ');
       });
    });
  }
  
  // Handle login password page
  if (window.location.pathname.includes('login-password.html')) {
    const email = sessionStorage.getItem('tm_login_email');
    if (!email) {
      // No email stored, redirect back to login
      window.location.href = './login.html';
      return;
    }
    
    // Display the email
    const displayEmail = document.getElementById('display-email');
    if (displayEmail) {
      displayEmail.textContent = email;
    }
    
    // Handle login button
    const btnLogin = document.getElementById('btn-login');
    if (btnLogin) {
      btnLogin.addEventListener('click', function() {
        const username = document.getElementById('login-username')?.value.trim();
        const password = document.getElementById('login-password')?.value.trim();
        
        if (!username) {
          showError('กรุณากรอกชื่อผู้ใช้');
          return;
        }
        
        if (!password) {
          showError('กรุณากรอกรหัสผ่าน');
          return;
        }
        
        // Call login API
        fetch(`${API_BASE}/api/auth/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            email: username,
            password: password
          })
        })
        .then(response => response.json())
        .then(data => {
          if (data.success && data.token) {
            // Store token in localStorage
            localStorage.setItem('access_token', data.token);
            // Clear session storage
            sessionStorage.removeItem('tm_login_email');
            // Redirect to dashboard
            window.location.href = '/home.html';
          } else {
            showError(data.error || 'ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง');
          }
        })
        .catch(error => {
          console.error('Login error:', error);
          showError('เกิดข้อผิดพลาดในการเข้าสู่ระบบ กรุณาลองใหม่อีกครั้ง');
        });
      });
    }
  }
  
  // Handle signup details form
  const formDetails = document.getElementById('form-details');
  if (formDetails) {
    formDetails.addEventListener('submit', function(e) {
      e.preventDefault();
      
      const email = sessionStorage.getItem('tm_signup_email');
      const fullName = document.getElementById('full-name')?.value.trim();
      const username = document.getElementById('username')?.value.trim();
      const password = document.getElementById('password')?.value.trim();
      
      if (!email || !fullName || !username || !password) {
        alert('กรุณากรอกข้อมูลให้ครบถ้วน');
        return;
      }
      
      // Remove password length requirement
      
      // Call register API
      fetch(`${API_BASE}/api/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: email,
          username: username,
          name: fullName,
          password: password
        })
      })
      .then(response => response.json())
      .then(data => {
        if (data.success && data.token) {
          // Store token in localStorage
          localStorage.setItem('access_token', data.token);
          // Clear stored email
          sessionStorage.removeItem('tm_signup_email');
          // Redirect to dashboard
          window.location.href = '/home.html';
        } else {
          alert(data.error || 'สมัครสมาชิกไม่สำเร็จ');
        }
      })
      .catch(error => {
        console.error('Register error:', error);
        alert('เกิดข้อผิดพลาดในการสมัครสมาชิก');
      });
    });
  }
  
  // Display stored email in signup details page
  const emailPreview = document.getElementById('email-preview');
  if (emailPreview) {
    const storedEmail = sessionStorage.getItem('tm_signup_email');
    if (storedEmail) {
      emailPreview.textContent = `อีเมล: ${storedEmail}`;
    }
  }
});

// Global function to handle email check
window.handleEmailCheck = function(email) {
  return fetch(`${API_BASE}/api/auth/login/check`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ email: email })
  })
  .then(response => response.json())
  .then(data => {
    return data.exists;
  })
  .catch(error => {
    console.error('Email check error:', error);
    return false;
  });
}