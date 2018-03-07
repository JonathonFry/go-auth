import React from 'react';

class Header extends React.Component {
  render() {
    return (
      <nav className="navbar navbar-light">
        <div className="container">

          <h1 className="navbar-brand">
            {this.props.appName.toLowerCase()}
          </h1>
        </div>
      </nav>
    );
  }
}

export default Header;