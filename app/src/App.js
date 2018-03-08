import React, { Component } from 'react';
import { withRouter } from 'react-router'
import { connect } from 'react-redux';
// import logo from './logo.svg';
import './App.css';
import Header from './components/Header';
import PropTypes from 'prop-types';

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

App.contextTypes = {
  router: PropTypes.object.isRequired
};

export default withRouter(connect(mapStateToProps, () => ({}) )(App));
