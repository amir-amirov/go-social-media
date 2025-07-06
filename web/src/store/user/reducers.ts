import * as types from './types'
import { User } from './types'

export const setUser: types.BaseContract<User | null> = (state, action) => {
  return {
    ...state,
    user: action.payload,
  }
}

export const removeUserDetails = (state: types.UserState) => {
  return {
    ...state,
    user: null,
  }
}

export const setToken: types.BaseContract<string> = (state, action) => {
  return {
    ...state,
    token: action.payload,
  }
}

export const removeToken = (state: types.UserState) => {
  return {
    ...state,
    token: '',
  }
}

export const resetUser = (): types.UserState => ({
  user: null,
  token: '',
})
