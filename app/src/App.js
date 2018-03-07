import React, { Component } from 'react';
import { connect } from 'react-redux';
// import logo from './logo.svg';
import './App.css';
import Header from './components/Header';
import Home from './components/Home';

const mapStateToProps = state => ({
  appName: state.appName,
  users: state.users
});

class App extends Component {

  render() {
    return (
      <div>
        <Header appName={this.props.appName} />
        <Home />
      </div>
    );
  }
}

export default connect(mapStateToProps, () => ({}) )(App);
