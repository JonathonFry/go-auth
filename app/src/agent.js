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

const Users = {
    all : _ => fetch(`${API_BASE_URL}/users`)
    .then(status)
    .then(json)
};

export default {
   Users
};