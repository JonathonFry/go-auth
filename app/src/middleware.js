import agent from './agent';
import {
  LOGIN,
  LOGOUT,
  REGISTER,
  ASYNC_START
} from './constants/actionTypes';

const promiseMiddleware = store => next => action => {
    if (isPromise(action.payload)) {
      store.dispatch({ type: ASYNC_START, subtype: action.type });
      action.payload.then(
        res => {
          action.payload = res;
          store.dispatch(action);
        },
        error => {
          action.error = true;
          action.payload = error.message;
          store.dispatch(action);
        }
      );
  
      return;
    }
  
    next(action);
};

const localStorageMiddleware = store => next => action => {
  if (action.type === REGISTER || action.type === LOGIN) {
    if (!action.error) {
      window.localStorage.setItem('jwt', action.payload.Token);
      agent.setToken(action.payload.Token);
    }
  } else if (action.type === LOGOUT) {
    window.localStorage.setItem('jwt', '');
    agent.setToken(null);
  }

  next(action);
};
  
function isPromise(v) {
    return v && typeof v.then === 'function';
}

  
export {
    promiseMiddleware,
    localStorageMiddleware
};