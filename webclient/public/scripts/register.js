document.getElementById('registerForm').addEventListener('submit', function(e) {
  e.preventDefault();

  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;
  const confirmPassword = document.getElementById('confirmPassword').value;

  if (password !== confirmPassword) {
    const messageEl = document.getElementById('registerMessage');
    messageEl.textContent = 'Passwords do not match';
    messageEl.className = 'response error';
    messageEl.style.display = 'block';
    return;
  }

  fetch('http://localhost:8080/user/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password })
  })
  .then(response => {
    if (!response.ok) {
      return response.json().then(err => { throw new Error(err.error || 'Registration failed') });
    }
    return response.json();
  })
  .then(data => {
    const messageEl = document.getElementById('registerMessage');
    messageEl.textContent = 'Registration successful! Logging you in...';
    messageEl.className = 'response success';
    messageEl.style.display = 'block';

    // Store token and user ID for auto-login
    localStorage.setItem('authToken', data.token);
    localStorage.setItem('userId', data.user_id);

    // Clear form
    document.getElementById('registerForm').reset();

    // Auto-redirect to app after 1 second
    setTimeout(() => {
      window.location.href = '/app';
    }, 1000);
  })
  .catch(error => {
    const messageEl = document.getElementById('registerMessage');
    messageEl.textContent = error.message;
    messageEl.className = 'response error';
    messageEl.style.display = 'block';
  });
});

// Check if user is already logged in
if (localStorage.getItem('authToken')) {
  // Verify token is still valid
  fetch('http://localhost:8080/user/me', {
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('authToken')
    }
  })
  .then(response => {
    if (response.ok) {
      window.location.href = '/app';
    } else {
      // Token invalid, clear it
      localStorage.removeItem('authToken');
      localStorage.removeItem('userId');
    }
  })
  .catch(() => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('userId');
  });
}
