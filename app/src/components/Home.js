import MainView from './MainView';
import React from 'react';
import { connect } from 'react-redux';
import agent from '../agent';

const Promise = global.Promise;

const mapStateToProps = state => ({
  appName: state.appName
});

const mapDispatchToProps = dispatch => ({
  onLoad: (payload) =>
  dispatch({ type: 'USERS_LOADED', payload }),
});
  

class Home extends React.Component {
  
  componentWillMount() {
    this.props.onLoad(agent.Users.all());
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