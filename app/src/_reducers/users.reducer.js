import { userConstants } from '../_constants';

export function users(state = {}, action) {
  switch (action.type) {
    case userConstants.GETALL_REQUEST:
      return {
        loading: true
      };
    case userConstants.GETALL_SUCCESS:
      return {
        items: action.users
      };
    case userConstants.GETALL_FAILURE:
      return { 
        error: action.error
      };
    case userConstants.GET_REQUEST:
      return {
        loading: true
      };
    case userConstants.GET_SUCCESS:
      return {
        items: action.user
      };
    case userConstants.GET_FAILURE:
      return { 
        error: action.error
      };
    default:
      return state
  }
}