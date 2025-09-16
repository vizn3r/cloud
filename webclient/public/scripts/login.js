document.getElementById('loginForm').addEventListener('submit', function(e) {
  e.preventDefault();

  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  fetch('http://localhost:8080/user/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password })
  })
  .then(response => {
    if (!response.ok) {
      return response.json().then(err => { throw new Error(err.error || 'Login failed') });
    }
    return response.json();
  })
  .then(data => {
    // Store token in localStorage
    localStorage.setItem('authToken', data.token);
    localStorage.setItem('userId', data.user_id);

    // Redirect to /app
    window.location.href = '/app';
  })
  .catch(error => {
    const messageEl = document.getElementById('loginMessage');
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
