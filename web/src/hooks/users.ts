import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import users from '../api/actions/users'
import { useUser } from '../store/user'

export const useTopUsers = () =>
  useQuery({
    queryKey: ['top'],
    queryFn: () => users.getTopUsers(),
    staleTime: 5 * 60_000,
  })

export const useFollowUser = () => {
  const queryClient = useQueryClient()
  const { user } = useUser()
  return useMutation({
    mutationFn: (userID: number) => users.followUser(userID),
    onSuccess: () => {
      //   queryClient.invalidateQueries({ queryKey: ['top'] })
      queryClient.invalidateQueries({ queryKey: ['feed'] })
      queryClient.invalidateQueries({ queryKey: ['follow-stats', user?.id] })
    },
  })
}

export const useUnFollowUser = () => {
  // const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (userID: number) => users.unfollowUser(userID),
    onSuccess: () => {
      //   queryClient.invalidateQueries({ queryKey: ['top'] })
      //   queryClient.invalidateQueries({ queryKey: ['feed'] })
    },
  })
}

export const useGetFollowStats = (userID: number | undefined) =>
  useQuery({
    queryKey: ['follow-stats', userID],
    queryFn: () => users.getFollowStats(userID),
    enabled: !!userID,
    staleTime: 5 * 60_000,
  })
