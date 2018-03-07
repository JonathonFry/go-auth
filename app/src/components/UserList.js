import React from 'react';

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
      {
        props.users.map(user => {
          return (
            <h2>{user.username}</h2>
          );
        })
      }
    </div>
  );
};

export default UserList;