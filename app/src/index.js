import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import { Provider } from 'react-redux';
import registerServiceWorker from './registerServiceWorker';
import store from './store';
import { Route } from 'react-router';
import { BrowserRouter } from 'react-router-dom'
import Home from './components/Home';
import Login from './components/Login';
import Logout from './components/Logout';
import Register from './components/Register';

ReactDOM.render((
    <Provider store={store}>
        <BrowserRouter>
            <App>
                <Route exact path="/" component={Home} />
                <Route path="/login" component={Login} />
                <Route path="/register" component={Register} />
                <Route path="/logout" component={Logout} />
            </App>
        </BrowserRouter>
    </Provider>
), document.getElementById('root'));

registerServiceWorker();
