import { authHeader } from '../_helpers';

const API_BASE_URL = 'http://localhost:8080'

export const userService = {
    login,
    logout,
    register,
    users,
    user
  };

function register(username, email, password) {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, password })
};

return fetch(`${API_BASE_URL}/register`, requestOptions)
    .then(response => {
        if (!response.ok) { 
            return Promise.reject(response.statusText);
        }

        return response.json();
    })
    .then(user => {
        if (user && user.token) {
          localStorage.setItem('jwt', user.token);
          localStorage.setItem('user', JSON.stringify(user));
        }
        
        return user;
    });
}

function login(username, password) {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
};

return fetch(`${API_BASE_URL}/login`, requestOptions)
    .then(response => {
        if (!response.ok) { 
            return Promise.reject(response.statusText);
        }

        return response.json();
    })
    .then(user => {
        if (user && user.token) {
          localStorage.setItem('jwt', user.token);
          localStorage.setItem('user', JSON.stringify(user));
        }
        
        return user;
    });
}

function logout() {
  localStorage.removeItem('jwt');
  localStorage.removeItem('user');
}

function user() {
  const requestOptions = { 
    method: 'GET',
    headers: authHeader() };

return fetch(`${API_BASE_URL}/user`, requestOptions)
    .then(response => {
        if (!response.ok) { 
            return Promise.reject(response.statusText);
        }

        return response.json();
    })
    .then(user => {
        if (user && user.token) {
          localStorage.setItem('user', JSON.stringify(user));
        }
        
        return user;
    });
}

function users() {
  const requestOptions = { 
    method: 'GET',
    headers: authHeader() };

return fetch(`${API_BASE_URL}/users`, requestOptions)
    .then(response => {
        if (!response.ok) { 
            return Promise.reject(response.statusText);
        }

        return response.json();
    });
}