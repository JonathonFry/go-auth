import React, { Component } from 'react';
import { withRouter } from 'react-router'
import { connect } from 'react-redux';
import './App.css';
import Header from './components/Header';
import agent from './agent';
import { APP_LOAD } from './constants/actionTypes';

const mapStateToProps = state => ({
  appLoaded: state.appLoaded,
  appName: state.appName,
  users: state.users,
  currentUser: state.currentUser
});

const mapDispatchToProps = dispatch => ({
  onLoad: (payload, token) =>
    dispatch({ type: APP_LOAD, payload, token }),
});

class App extends Component {

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
