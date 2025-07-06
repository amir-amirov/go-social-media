import { Avatar, Paper, InputBase, IconButton } from '@mui/material'
import { useState } from 'react'
import SendIcon from '@mui/icons-material/Send'

const CommentInput = ({
  onSubmit,
  userAvatar,
  loading,
}: {
  onSubmit: (comment: string) => void
  userAvatar?: string
  loading: boolean
}) => {
  const [comment, setComment] = useState('')

  const handleSubmit = () => {
    if (comment.trim()) {
      onSubmit(comment)
      setComment('')
    }
  }

  return (
    <Paper
      elevation={0}
      sx={{
        display: 'flex',
        alignItems: 'center',
        p: 1,
        borderRadius: 3,
        border: '1px solid #ccc',
        backgroundColor: 'background.paper',
        marginBottom: 3,
      }}
    >
      <Avatar src={userAvatar} sx={{ width: 32, height: 32, mr: 1 }} />
      <InputBase
        disabled={loading}
        placeholder="Write a comment..."
        fullWidth
        value={comment}
        onChange={(e) => setComment(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault()
            handleSubmit()
          }
        }}
        sx={{
          fontSize: 14,
        }}
        multiline
        maxRows={3}
      />
      <IconButton
        onClick={handleSubmit}
        disabled={!comment.trim() || loading}
        size="small"
        sx={{ ml: 1 }}
      >
        <SendIcon fontSize="small" color={comment == '' ? 'disabled' : 'primary'} />
      </IconButton>
    </Paper>
  )
}

export default CommentInput
