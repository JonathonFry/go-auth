import React from 'react';

const UserView = props => {
  const user = props.user;

  return (
    <div className="article-preview">
      <div className="article-meta">
        <div className="info">
          <span className="author">
            {user.username}
          </span>
          <br/>
          <span className="email">
            {user.email_address}
          </span>
          <br/>
          <span className="password">
            {user.password}
          </span>
        </div>
      </div>
    </div>
  );
}

export default UserView;