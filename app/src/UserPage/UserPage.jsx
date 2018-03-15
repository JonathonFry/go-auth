import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import { userActions } from '../_actions';

class UserPage extends React.Component {
    componentDidMount() {
        this.props.dispatch(userActions.user());
    }

    render() {
        const { loading, user, error } = this.props;
        return (
            <div className="col-md-6 col-md-offset-3">
            {user &&
                <div>
                    <h1>Hi {user.username}!</h1>
                    <p>Email address: {user.email_address}</p>
                </div>
            } 
                
                {loading && <em>Loading user...</em>}
                {error && <span className="text-danger">ERROR: {error}</span>}
                <p>
                    <Link to="/login">Logout</Link>
                </p>
            </div>
        );
    }
}

function mapStateToProps(state) {
    const { loading, user, error}  = state.users;
    return {
        loading,
        user,
        error
    };
}

const connectedUserPage = connect(mapStateToProps)(UserPage);
export { connectedUserPage as UserPage };