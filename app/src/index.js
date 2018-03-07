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
     users: []
};
const reducer = function(state = defaultState, action) {
    return state;
};

const store = createStore(reducer, applyMiddleware(promiseMiddleware));

ReactDOM.render((
    <Provider store={store}>
        <App />
    </Provider>
), document.getElementById('root'));

registerServiceWorker();
