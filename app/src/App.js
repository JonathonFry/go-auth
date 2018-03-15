import React, { Component } from 'react';
import { withRouter } from 'react-router'
import { push } from 'react-router-redux'
import { connect } from 'react-redux';
import './App.css';
import store from './store'
import Header from './components/Header';
import agent from './agent';
import { APP_LOAD, REDIRECT } from './constants/actionTypes';

const mapStateToProps = state => ({
  appLoaded: state.appLoaded,
  appName: state.appName,
  users: state.users,
  currentUser: state.currentUser,
  redirectTo: state.redirectTo
});

const mapDispatchToProps = dispatch => ({
  onLoad: (payload, token) =>
    dispatch({ type: APP_LOAD, payload, token }),
  onRedirect: () =>
    dispatch({ type: REDIRECT })
});

class App extends Component {

  componentWillReceiveProps(nextProps) {
    if (nextProps.redirectTo) {
      store.dispatch(push(nextProps.redirectTo));
      this.props.onRedirect();
    }
  }

  componentWillMount() {
    const token = window.localStorage.getItem('jwt');
    if (token) {
      agent.setToken(token);
    }

    this.props.onLoad(token ? agent.Auth.current() : null, token);
  }

  render() {
    return (
      <div>
        <Header appName={this.props.appName} currentUser={this.props.currentUser}/>
        {this.props.children}
      </div>
    );
  }
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(App));
