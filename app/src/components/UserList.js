import React from 'react';
import UserView from './UserView';

const UserList = props => {
  if (!props.users) {
    return (
      <div className="user-preview">Loading...</div>
    );
  }

  if (props.users.length === 0) {
    return (
      <div className="user-preview">
        No users registered here... yet.
      </div>
    );
  }

  return (
    <div>
      <h3>Users:</h3>
      {
        props.users.map(user => {
          return (
            <UserView user={user} key={user.username} />
          );
        })
      }
    </div>
  );
};

export default UserList;