import { useInfiniteQuery, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import posts, { Post } from '../api/actions/posts'
import { toast } from 'react-toastify'

const LIMIT = 5

export const useFeedInfinite = (search: string) =>
  useInfiniteQuery({
    queryKey: ['feed', search],
    queryFn: ({ pageParam = 1 }) => posts.getFeed(pageParam, search),
    initialPageParam: 1,
    getNextPageParam: (lastPage, pages) => {
      return lastPage.length < LIMIT ? undefined : pages.length + 1
    },
    staleTime: 5 * 60 * 1000,
  })

export const useUserPostsInfinite = (userID: number | undefined) =>
  useInfiniteQuery({
    queryKey: ['posts', userID],
    queryFn: ({ pageParam = 1 }) => posts.getUserPosts(pageParam, userID),
    enabled: !!userID,
    initialPageParam: 1,
    getNextPageParam: (lastPage, pages) => {
      return lastPage.length < LIMIT ? undefined : pages.length + 1
    },
  })

export const useGetPost = (postID: number) =>
  useQuery({
    queryKey: ['post', postID],
    queryFn: () => posts.getPost(postID),
    enabled: !!postID,
  })

export const useCreatePost = () => {
  return useMutation({
    mutationFn: (newPost: { title: string; content: string }) => posts.createPost(newPost),
  })
}

export const useUpdatePost = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ postID, updatedPost }: { postID: number; updatedPost: Post }) =>
      posts.updatePost(postID, updatedPost),
    onSuccess: (_, postID) => {
      toast.success('Post updated successfully!')
      queryClient.invalidateQueries({ queryKey: ['feed'] })
      queryClient.invalidateQueries({ queryKey: ['posts'] })
      queryClient.invalidateQueries({ queryKey: ['post', postID] })
    },
  })
}

export const useDeletePost = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (postID: number) => posts.deletePost(postID),
    onSuccess: () => {
      toast.success('Post deleted successfully!')
      queryClient.invalidateQueries({ queryKey: ['feed'] })
      queryClient.invalidateQueries({ queryKey: ['posts'] })
    },
  })
}
