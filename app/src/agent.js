const API_BASE_URL = 'http://localhost:8080'

function status(response) {
    console.log(response)
    if (response.status >= 200 && response.status < 300 && !response.redirected) {
        console.log('response success')
      return Promise.resolve(response)
    } else {
        console.log('response rejected')
      return Promise.reject(new Error(response.body))
    }
}
  
function json(response) {
    return response.json()
}

const requests = {
    del: url => fetch(url,{ method: "DELETE" }).then(status).then(json),
    get: url => fetch(url,{ method: "GET" }).then(status).then(json),
    put: (url, body) => fetch(url,{ method: "PUT", body: body}).then(status).then(json),
    post: (url, body) => fetch(url,{ method: "POST", body: body}).then(status).then(json)
  };

const Auth = {
    current: () =>  requests.get(`${API_BASE_URL}/user`),
    all: () =>      requests.get(`${API_BASE_URL}/users`),
    login: (username, password) =>  requests.post(`${API_BASE_URL}/login`, { username, password }),
    register: (username, email, password) =>  requests.post(`${API_BASE_URL}/register`,{ username, email, password })
  };

export default {
   Auth
};