import MainView from './MainView';
import React from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom'
import agent from '../agent';
import { USERS_LOADED } from '../constants/actionTypes';

const mapStateToProps = state => ({
  appName: state.appName,
  inProgress: state.inProgress
});

const mapDispatchToProps = dispatch => ({
  onLoad: (payload) =>
  dispatch({ type: USERS_LOADED, payload }),
});
  

class Home extends React.Component {
  
  componentWillMount() {
    var onLoad = this.props.onLoad;
    agent.Auth.all()
    .then(function(response) {
      onLoad(response);  
    })
    .catch(function(error) {
      console.log('There has been a problem with your fetch operation: ', error);
    })
  }

  render() {
    if (this.props.inProgress) {
      return (
        <div className="home-page">
        </div>
      );
    }
    return (
      <div className="home-page">
        <div className="container page">
          <div className="row">
            <MainView />
          </div>
        </div>

      </div>
    );
  }
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Home));