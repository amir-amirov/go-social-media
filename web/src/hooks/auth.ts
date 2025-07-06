import { useMutation } from '@tanstack/react-query'
import auth from '../api/actions/auth'
import { UserState } from '../store/user/types'
import { useUser } from '../store/user'

export const useLogin = () => {
  const { setUser, setToken } = useUser()
  return useMutation({
    mutationFn: auth.login,
    onSuccess: (data: UserState) => {
      setUser(data.user)
      setToken(data.token)
    },
  })
}

export const useSignIn = () =>
  useMutation({
    mutationFn: auth.signin,
  })

export const useVerifyEmail = () =>
  useMutation({
    mutationFn: auth.verifyEmail,
  })
