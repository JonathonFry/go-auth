import React from 'react';
import { Link } from 'react-router-dom';

class Header extends React.Component {
  render() {
    if (this.props.currentUser) {
      return (
        <nav className="navbar navbar-light">
          <div className="container">
            <Link to="/" className="navbar-brand">
              {this.props.appName.toLowerCase()}
            </Link>
  
            <ul className="nav navbar-nav pull-xs-right">
              <li className="nav-item">
                <Link to="/user" className="nav-link">
                  {this.props.currentUser.username}
                </Link>
              </li>
              <li className="nav-item">
                <Link to="/logout" className="nav-link">
                  logout
                </Link>
              </li>
            </ul>
          </div>
        </nav>
      );
    } 
    return (
      <nav className="navbar navbar-light">
        <div className="container">
          <Link to="/" className="navbar-brand">
            {this.props.appName.toLowerCase()}
          </Link>

          <ul className="nav navbar-nav pull-xs-right">
            <li className="nav-item">
              <Link to="/" className="nav-link">
                Home
              </Link>
            </li>

            <li className="nav-item">
              <Link to="/login" className="nav-link">
                Sign in
              </Link>
            </li>
            <li className="nav-item">
              <Link to="/register" className="nav-link">
                Register
              </Link>
            </li>
          </ul>
        </div>
      </nav>
    );
  }
}

export default Header;