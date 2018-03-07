import UserList from './UserList';
import React from 'react';
import { connect } from 'react-redux';

const mapStateToProps = state => ({
  users: state.users
});

const MainView = props => {
  return (
    <div className="col-md-9">
      <UserList
        users={props.users} 
      />
    </div>
  );
};

export default connect(mapStateToProps, () => ({}))(MainView);