import { useAppDispatch, useAppSelector as useSelector } from '../index'
import { userActions } from './slice'
import { User } from './types'

export const useUser = () => {
  const dispatch = useAppDispatch()

  const setUser = (user: User | null) => {
    dispatch(userActions.setUser(user))
  }

  const setToken = (token: string) => {
    dispatch(userActions.setToken(token))
  }

  return {
    user: useSelector(({ user }) => user.user),
    token: useSelector(({ user }) => user.token),
    setUser,
    setToken,
    removeUserDetails: () => dispatch(userActions.removeUserDetails()),
    removeToken: () => dispatch(userActions.removeToken()),
    resetUser: () => dispatch(userActions.resetUser()),
  }
}
