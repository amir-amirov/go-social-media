import {
  Avatar,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Slide,
  Stack,
  TextField,
  Typography,
} from '@mui/material'
import { TransitionProps } from '@mui/material/transitions'
import React, { forwardRef, useState } from 'react'
import { useCreatePost } from '../../hooks/posts'
import { useUser } from '../../store/user'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'react-toastify'

const DEFAULT_AVATAR =
  'https://firebasestorage.googleapis.com/v0/b/auth-2c46a.appspot.com/o/user-profile-icon-profile-avatar-user-icon-male-icon-face-icon-profile-icon-free-png.webp?alt=media&token=a96e4694-e39e-4d96-80be-034d19eebb52'

const Transition = forwardRef(function Transition(
  props: TransitionProps & { children: React.ReactElement },
  ref: React.Ref<unknown>
) {
  return <Slide direction="up" ref={ref} {...props} />
})

interface Props {
  open: boolean
  onClose: () => void
}

export default function CreatePostModal({ open, onClose }: Props) {
  const { user } = useUser()
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const { mutate, isPending } = useCreatePost()
  const queryClient = useQueryClient()

  const handleSubmit = async () => {
    if (!title.trim() || !content.trim()) return

    mutate(
      { title, content },
      {
        onSuccess: () => {
          toast.success('Post created successfully!')
          queryClient.invalidateQueries({ queryKey: ['feed'] })
          setTitle('')
          setContent('')
          onClose()
        },
        onError: () => {
          toast.error('Post is NOT created!')
        },
      }
    )
  }

  return (
    <Dialog
      open={open}
      onClose={onClose}
      TransitionComponent={Transition}
      fullWidth
      maxWidth="sm"
      PaperProps={{ sx: { borderRadius: 3, p: 2 } }}
    >
      <DialogTitle>
        <Stack direction="row" spacing={2} alignItems="center">
          <Avatar src={user?.avatar ?? DEFAULT_AVATAR} />
          <Typography variant="h6">Create Post</Typography>
        </Stack>
      </DialogTitle>

      <DialogContent sx={{ pt: 1 }}>
        <Stack spacing={2}>
          <TextField
            placeholder="Title"
            variant="outlined"
            fullWidth
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={isPending}
          />
          <TextField
            placeholder="Content"
            variant="outlined"
            fullWidth
            multiline
            minRows={4}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            disabled={isPending}
          />
        </Stack>
      </DialogContent>

      <DialogActions sx={{ justifyContent: 'space-between', px: 3 }}>
        <Button onClick={onClose} color="inherit" disabled={isPending}>
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={handleSubmit}
          disabled={isPending || !title || !content}
        >
          {isPending ? 'Posting...' : 'Post'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
