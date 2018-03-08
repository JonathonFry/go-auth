import { applyMiddleware, createStore } from 'redux';
import { promiseMiddleware } from './middleware';

const defaultState = {
    appName: 'go-auth',
    users: null
};
const reducer = function(state = defaultState, action) {
    switch (action.type) {
       case 'USERS_LOADED':
           return { ...state, users: action.payload}
       default: 
           return state;
   }
};

const middleware = applyMiddleware(promiseMiddleware);

const store = createStore(reducer, middleware);

export default store;