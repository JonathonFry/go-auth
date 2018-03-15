import React from 'react';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom'
import { LOGOUT } from '../constants/actionTypes';
import { Redirect } from "react-router-dom";




class Logout extends React.Component {
    componentWillMount () {
        this.props.dispatch({ type: LOGOUT });
    }

    render() {
        return (
          <Redirect to="/" />
        );
    }
};

export default withRouter(connect()(Logout))