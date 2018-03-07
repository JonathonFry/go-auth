import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { Provider } from 'react-redux';
import registerServiceWorker from './registerServiceWorker';
import { applyMiddleware, createStore } from 'redux';
import { promiseMiddleware } from './middleware';

const defaultState = {
     appName: 'go-auth',
     users: null
};
const reducer = function(state = defaultState, action) {
     switch(action.type) {
        case 'USERS_LOADED':
            return { ...state, users: action.payload}
        default: 
            return state;
    }
};

const store = createStore(reducer, applyMiddleware(promiseMiddleware));

ReactDOM.render((
    <Provider store={store}>
        <App />
    </Provider>
), document.getElementById('root'));

registerServiceWorker();
