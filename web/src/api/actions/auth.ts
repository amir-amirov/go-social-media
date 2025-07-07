import fetchData from '..'

export interface LoginPayload {
  email: string
  password: string
}
export interface SignInPayload {
  username: string
  email: string
  password: string
}
export interface User {
  id: number
  name: string /*…*/
}

const login = async (creds: LoginPayload) => {
  const { response, status } = await fetchData({
    url: '/auth/login',
    body: {
      email: creds.email,
      password: creds.password,
    },
    method: 'POST',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message || 'Login failed')
  }

  return response.data
}

const signin = async (creds: SignInPayload) => {
  const { response, status } = await fetchData({
    url: '/auth/user',
    body: {
      username: creds.username,
      email: creds.email,
      password: creds.password,
    },
    method: 'POST',
  })

  if (status !== 201) {
    throw new Error(response?.message)
  }

  return response.data
}

const verifyEmail = async (code: string) => {
  const { response, status } = await fetchData({
    url: `/users/activate/${code}`,
    method: 'PUT',
  })

  if (status !== 200) {
    throw new Error('Email verification failed')
  }

  return response.data
}

const resendCode = async (email: string | undefined) => {
  if (!email) return
  const { response, status } = await fetchData({
    url: `/auth/token`,
    method: 'POST',
    body: {
      email,
    },
  })

  if (status !== 200) {
    throw new Error('Email verification failed')
  }

  return response.data
}

const auth = {
  login,
  signin,
  verifyEmail,
  resendCode,
}

export default auth
