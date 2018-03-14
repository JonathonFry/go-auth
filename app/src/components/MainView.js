import React from 'react';
import { connect } from 'react-redux';
import UserView from './UserView';

const mapStateToProps = state => ({
  currentUser: state.currentUser,
  users: state.users
});

const MainView = props => {
  if (props.currentUser){
    return (
      <div className="col-md-9">
      <UserView user={props.currentUser}
        />
      </div>
    );
  } else {
    return (
      <div className="col-md-9">
      <p>login</p>
      </div>
    );
  }
};

export default connect(mapStateToProps, () => ({}))(MainView);