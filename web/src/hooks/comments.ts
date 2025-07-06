import { useMutation, useQueryClient } from '@tanstack/react-query'
import comments from '../api/actions/comments'
import { toast } from 'react-toastify'

export const useCreateComment = () => {
  return useMutation({
    mutationFn: ({ postID, comment }: { postID: number | undefined; comment: string }) =>
      comments.createComment(postID, comment),
  })
}

export const useDeleteComment = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ postID, commentID }: { postID: number; commentID: number }) =>
      comments.deleteComment(postID, commentID),
    onSuccess: (_, vars) => {
      toast.success('Comment deleted successfully!')
      queryClient.invalidateQueries({ queryKey: ['feed'] })
      queryClient.invalidateQueries({ queryKey: ['post', vars.postID] })
    },
  })
}
