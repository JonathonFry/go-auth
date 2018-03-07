import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { Provider } from 'react-redux';
import registerServiceWorker from './registerServiceWorker';
import store from './store';
import { Router, Route, IndexRoute } from 'react-router';
import { HashRouter } from 'react-router-dom';
import Home from './components/Home';

ReactDOM.render((
    <Provider store={store}>
        <HashRouter>
            <Route path="/" component={App}>
                <IndexRoute component={Home} />
            </Route>
        </HashRouter>
    </Provider>
), document.getElementById('root'));

registerServiceWorker();
