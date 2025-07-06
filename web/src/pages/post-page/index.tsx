import React, { useState } from 'react'
import {
  Box,
  Typography,
  Chip,
  Avatar,
  Stack,
  Divider,
  Paper,
  CircularProgress,
  IconButton,
  Menu,
  MenuItem,
} from '@mui/material'
import dayjs from 'dayjs'
import { useParams, useNavigate } from 'react-router-dom'
import MoreVertIcon from '@mui/icons-material/MoreVert'
import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'

import CommentInput from '../../components/CommentInput'
import DeleteConfirmDialog from '../../components/DeleteConfirmDialog/DeleteConfirmDialog'
import EditPostModal from '../../components/EditPostModal/EditPostModal'

import { useGetPost, useDeletePost, useUpdatePost } from '../../hooks/posts'
import { useCreateComment, useDeleteComment } from '../../hooks/comments'
import { useUser } from '../../store/user'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'react-toastify'

export default function PostPage() {
  const { id } = useParams()
  const postID = Number(id)

  const { data: post, isPending, isError, refetch } = useGetPost(postID)
  const { user: currentUser } = useUser()
  const queryClient = useQueryClient()
  const navigate = useNavigate()

  const { mutate: deletePost, isPending: isDeleting } = useDeletePost()
  const { mutate: updatePost, isPending: isUpdating } = useUpdatePost()
  const { mutate: createComment, isPending: isCreatingComment } = useCreateComment()
  const { mutate: deleteComment, isPending: isDeletingComment } = useDeleteComment()

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [editOpen, setEditOpen] = useState(false)
  const [commentDeleteOpen, setCommentDeleteOpen] = useState(false)
  const [commentToDelete, setCommentToDelete] = useState<number | null>(null)

  const isAuthor = currentUser?.id === post?.user_id

  const handleMenuClick = (e: React.MouseEvent<HTMLElement>) => setAnchorEl(e.currentTarget)
  const closeMenu = () => setAnchorEl(null)

  const openDelete = () => {
    closeMenu()
    setDeleteOpen(true)
  }

  const openEdit = () => {
    closeMenu()
    setEditOpen(true)
  }

  const confirmDelete = () => {
    deletePost(postID, {
      onSuccess: () => navigate('/feed'),
    })
    setDeleteOpen(false)
  }

  const saveEdit = (data: { title: string; content: string; tags: string[] }) => {
    updatePost(
      { postID, updatedPost: data },
      {
        onSuccess: () => refetch(),
      }
    )
    setEditOpen(false)
  }

  const addComment = (comment: string) => {
    createComment(
      { postID, comment },
      {
        onSuccess: () => {
          toast.success('Comment created')
          queryClient.invalidateQueries({ queryKey: ['post', postID] })
        },
        onError: () => toast.error('Comment not created'),
      }
    )
  }

  const askDeleteComment = (id: number) => {
    setCommentToDelete(id)
    setCommentDeleteOpen(true)
  }

  const confirmDeleteComment = () => {
    if (commentToDelete == null) return
    deleteComment({ postID, commentID: commentToDelete })
    setCommentDeleteOpen(false)
  }

  if (isPending) {
    return (
      <Box display="flex" justifyContent="center" mt={12}>
        <CircularProgress />
      </Box>
    )
  }

  if (isError || !post) {
    return (
      <Box display="flex" justifyContent="center" mt={12}>
        <Typography variant="h5">Something went wrong…</Typography>
      </Box>
    )
  }

  return (
    <>
      <Box maxWidth="md" mx="auto" my={4} px={2}>
        <Paper elevation={3} sx={{ p: 4, borderRadius: 3 }}>
          <Stack direction="row" justifyContent="space-between" alignItems="flex-start">
            <Typography variant="h4" fontWeight="bold">
              {post.title}
            </Typography>

            {isAuthor && (
              <>
                <IconButton
                  aria-label="post actions"
                  aria-controls={anchorEl ? 'post-menu' : undefined}
                  aria-haspopup="true"
                  onClick={handleMenuClick}
                >
                  <MoreVertIcon />
                </IconButton>

                <Menu
                  id="post-menu"
                  anchorEl={anchorEl}
                  open={Boolean(anchorEl)}
                  onClose={closeMenu}
                  anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                  transformOrigin={{ vertical: 'top', horizontal: 'right' }}
                >
                  <MenuItem onClick={openEdit}>
                    <EditIcon fontSize="small" sx={{ mr: 1 }} />
                    Edit
                  </MenuItem>
                  <MenuItem onClick={openDelete}>
                    <DeleteIcon fontSize="small" color="error" sx={{ mr: 1 }} />
                    Delete
                  </MenuItem>
                </Menu>
              </>
            )}
          </Stack>

          <Stack direction="row" spacing={2} alignItems="center" mb={2} mt={1}>
            <Avatar src={post.user.avatar || undefined}>
              {post.user.username[0].toUpperCase()}
            </Avatar>
            <Box>
              <Typography variant="subtitle1">{post.user.username}</Typography>
              <Typography variant="caption" color="text.secondary">
                Posted on {dayjs(post.created_at).format('MMMM D, YYYY')}
              </Typography>
            </Box>
          </Stack>

          {/*  Content  */}
          <Typography variant="body1" sx={{ whiteSpace: 'pre-line' }}>
            {post.content}
          </Typography>

          {/*  Tags  */}
          <Stack direction="row" spacing={1} mt={3} mb={2}>
            {post.tags?.map((tag: any) => (
              <Chip key={tag} label={`#${tag}`} variant="outlined" />
            ))}
          </Stack>

          <Divider sx={{ my: 3 }} />

          {/*  Comments  */}
          <Typography variant="h6" fontWeight="bold" gutterBottom>
            {post.comments.length} Comments
          </Typography>

          <CommentInput
            onSubmit={addComment}
            userAvatar={currentUser?.avatar}
            loading={isCreatingComment}
          />

          <Stack spacing={2}>
            {post.comments.map((comment: any) => (
              <Paper key={comment.id} variant="outlined" sx={{ p: 2 }}>
                <Stack direction="row" spacing={2}>
                  <Avatar src={comment.user.avatar || undefined}>
                    {comment.user.username[0].toUpperCase()}
                  </Avatar>
                  <Box flexGrow={1}>
                    <Typography variant="subtitle2">{comment.user.username}</Typography>
                    <Typography variant="caption" color="text.secondary" display="block">
                      {dayjs(comment.created_at).format('MMM D, YYYY HH:mm')}
                    </Typography>
                    <Typography mt={1}>{comment.content}</Typography>
                  </Box>
                  {isAuthor && (
                    <IconButton
                      size="small"
                      color="error"
                      aria-label="delete comment"
                      onClick={() => askDeleteComment(comment.id)}
                      sx={{ alignSelf: 'start' }}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  )}
                </Stack>
              </Paper>
            ))}
          </Stack>
        </Paper>
      </Box>

      <DeleteConfirmDialog
        open={deleteOpen}
        loading={isDeleting}
        onClose={() => setDeleteOpen(false)}
        onConfirm={confirmDelete}
      />

      <DeleteConfirmDialog
        title="Delete comment"
        subtitle="This action can’t be undone. Are you absolutely sure you want to delete this comment?"
        open={commentDeleteOpen}
        loading={isDeletingComment}
        onClose={() => setCommentDeleteOpen(false)}
        onConfirm={confirmDeleteComment}
      />

      <EditPostModal
        open={editOpen}
        loading={isUpdating}
        onClose={() => setEditOpen(false)}
        onSave={saveEdit}
        post={post}
      />
    </>
  )
}
