const API_BASE_URL = 'http://localhost:8080'

let token = null;

function status(response) {
  console.log('status ' + response.status);
    if (response.status >= 200 && response.status < 300 && !response.redirected) {
      return Promise.resolve(response)
    } else {
      return Promise.reject(new Error(response.body))
    }
}
  
function json(response) {
    return response.json()
}

const requests = {
    del: url => fetch(url,{ method: "DELETE", headers: {"Authorization": token}}).then(status).then(json),
    get: url => fetch(url,{ method: "GET", headers: {"Authorization": token}}).then(status).then(json),
    put: (url, body) => fetch(url,{ method: "PUT", body: body, headers: {"Authorization": token}}).then(status).then(json),
    post: (url, body) => fetch(url,{ method: "POST", body: body}).then(status).then(json)
  };

const Auth = {
    current: () =>  requests.get(`${API_BASE_URL}/user`),
    all: () =>      requests.get(`${API_BASE_URL}/users`),
    login: (username, password) =>  requests.post(`${API_BASE_URL}/login`, JSON.stringify({ username, password })),
    register: (username, email, password) => requests.post(`${API_BASE_URL}/register`, JSON.stringify({ username, email, password }))
  };

export default {
   Auth,
  setToken: _token => { token = _token; }
};