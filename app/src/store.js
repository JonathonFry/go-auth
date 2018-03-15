import { applyMiddleware, createStore } from 'redux';
import { promiseMiddleware, localStorageMiddleware } from './middleware';
import {
    APP_LOAD,
    REGISTER,
    UPDATE_FIELD_AUTH,
    USERS_LOADED,
    LOGIN,
    ASYNC_START,
    REDIRECT
  } from './constants/actionTypes';

const defaultState = {
    appName: 'go-auth',
    users: null,
    username: '',
    email: '',
    password: '',
    token: null,
    currentUser: null,
    redirectTo: null,
};
const reducer = function(state = defaultState, action) {
    switch (action.type) {
        case APP_LOAD:
            return {
                ...state,
                token: action.token || null,
                appLoaded: true,
                currentUser: action.payload ? action.payload : null
            };
        case REDIRECT:
            return { ...state, redirectTo: null };
        case ASYNC_START:
            if (action.subtype === LOGIN || action.subtype === REGISTER) {
                return { ...state, inProgress: true };
            } else {
                return { ...state };
            }
        case USERS_LOADED:
           return { ...state, users: action.payload}
        case UPDATE_FIELD_AUTH:
            return { ...state, [action.key]: action.value };
        case REGISTER:
            return {
              ...state,
              inProgress: false,
              errors: action.error ? action.payload.errors : null
            };
        case LOGIN:
            return {
              ...state,
              inProgress: false,
              redirectTo: '/',
              errors: action.error ? action.payload.errors : null
            };
        default: 
            return state;
   }
};

const middleware = applyMiddleware(promiseMiddleware, localStorageMiddleware);

const store = createStore(reducer, middleware);

export default store;