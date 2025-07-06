import { CaseReducer, PayloadAction } from '@reduxjs/toolkit'

export type UserState = {
  token: string
  user: User | null
}

export type User = {
  id: number
  username: string
  email: string
  avatar: string
}

export type RegisterRequest = {
  username: string
  email: string
  password: string
}

export type LoginRequest = {
  email: string
  password: string
}

export type LoginResponse = {
  user: User
  token: string
}

// export type UpdateProfileRequest = {
//   id: number
//   username: string
//   email: string
//   avatar: string
// }

// export type UpdateProfileResponse = {
//   user: User
// }

// Contracts
export type BaseContract<T> = CaseReducer<UserState, PayloadAction<T>>
