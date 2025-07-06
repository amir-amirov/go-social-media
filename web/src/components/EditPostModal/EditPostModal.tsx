import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Stack,
} from '@mui/material'
import { useState, useEffect } from 'react'

interface EditPostModalProps {
  loading: boolean
  open: boolean
  onClose: () => void
  onSave: (updated: { title: string; content: string; tags: string[] }) => void
  post: { title: string; content: string; tags: string[] }
}

export default function EditPostModal({
  loading,
  open,
  onClose,
  onSave,
  post,
}: EditPostModalProps) {
  const [title, setTitle] = useState(post.title)
  const [content, setContent] = useState(post.content)
  const [tagsInput, setTagsInput] = useState(post.tags?.join(', '))

  // reset when the dialog is opened again
  useEffect(() => {
    if (open) {
      setTitle(post.title)
      setContent(post.content)
      setTagsInput(post.tags?.join(', '))
    }
  }, [open, post])

  const handleSave = () => {
    onSave({
      title,
      content,
      tags: tagsInput
        ? tagsInput
            .split(',')
            .map((t) => t.trim())
            .filter(Boolean)
        : [],
    })
  }

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>Edit post</DialogTitle>

      <DialogContent>
        <Stack spacing={2} sx={{ mt: 1 }}>
          <TextField
            disabled={loading}
            label="Title"
            fullWidth
            value={title}
            onChange={(e) => setTitle(e.target.value)}
          />
          <TextField
            disabled={loading}
            label="Content"
            fullWidth
            multiline
            minRows={4}
            value={content}
            onChange={(e) => setContent(e.target.value)}
          />
          <TextField
            disabled={loading}
            variant="filled"
            label="Tags (comma‑separated)"
            fullWidth
            value={tagsInput}
            onChange={(e) => setTagsInput(e.target.value)}
          />
        </Stack>
      </DialogContent>

      <DialogActions sx={{ pr: 3, pb: 2 }}>
        <Button onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        <Button variant="contained" onClick={handleSave} disabled={loading}>
          Save
        </Button>
      </DialogActions>
    </Dialog>
  )
}
