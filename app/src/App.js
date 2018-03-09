import React, { Component } from 'react';
import { withRouter } from 'react-router'
import { connect } from 'react-redux';
import './App.css';
import Header from './components/Header';

const mapStateToProps = state => ({
  appName: state.appName,
  users: state.users
});

class App extends Component {

  render() {
    return (
      <div>
        <Header appName={this.props.appName} />
        {this.props.children}
      </div>
    );
  }
}

export default withRouter(connect(mapStateToProps, () => ({}) )(App));
