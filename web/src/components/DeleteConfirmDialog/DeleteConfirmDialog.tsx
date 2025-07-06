import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Button,
} from '@mui/material'

interface DeleteConfirmDialogProps {
  title?: string
  subtitle?: string
  loading: boolean
  open: boolean
  onClose: () => void
  onConfirm: () => void
}

export default function DeleteConfirmDialog({
  title,
  subtitle,
  loading,
  open,
  onClose,
  onConfirm,
}: DeleteConfirmDialogProps) {
  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{title ?? 'Delete post'}</DialogTitle>

      <DialogContent>
        <DialogContentText>
          {subtitle ??
            'This action can’t be undone. Are you absolutely sure you want to delete this post?'}
        </DialogContentText>
      </DialogContent>

      <DialogActions sx={{ pr: 3, pb: 2 }}>
        <Button onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        <Button variant="contained" color="error" onClick={onConfirm} disabled={loading}>
          Delete
        </Button>
      </DialogActions>
    </Dialog>
  )
}
