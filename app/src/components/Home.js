import MainView from './MainView';
import React from 'react';
import { connect } from 'react-redux';
import agent from '../agent';
import { USERS_LOADED } from '../constants/actionTypes';

const mapStateToProps = state => ({
  appName: state.appName
});

const mapDispatchToProps = dispatch => ({
  onLoad: (payload) =>
  dispatch({ type: USERS_LOADED, payload }),
});
  

class Home extends React.Component {
  
  componentWillMount() {
    agent.Auth.all()
    .then(function(response){
      this.props.onLoad(response);  
    })
    .catch(function(error) {
      console.log('There has been a problem with your fetch operation: ', error.message);
    })
  }

  render() {
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

export default connect(mapStateToProps, mapDispatchToProps)(Home);