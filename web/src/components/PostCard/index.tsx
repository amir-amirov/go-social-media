import {
  Card,
  CardHeader,
  CardContent,
  Typography,
  Stack,
  Chip,
  Button,
  IconButton,
  Avatar,
  Menu,
  MenuItem,
  Box,
} from '@mui/material'
import CommentIcon from '@mui/icons-material/Comment'
import MoreVertIcon from '@mui/icons-material/MoreVert'
import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'
import { Link as RouterLink } from 'react-router-dom'
import { useState } from 'react'
import DeleteConfirmDialog from '../../components/DeleteConfirmDialog/DeleteConfirmDialog'
import EditPostModal from '../../components/EditPostModal/EditPostModal'
import { User } from '../../store/user/types'
import { useUser } from '../../store/user'
import { useDeletePost, useUpdatePost } from '../../hooks/posts'

interface Props {
  post: Post
}

interface Post {
  id: number
  title: string
  content: string
  tags: string[]
  created_at: string
  updated_at: string
  comments_count: number
  user: User
}

export default function PostCard({ post }: Props) {
  const { user } = post
  const { user: currentUser } = useUser()

  const isAuthor = currentUser?.id === user.id

  const { mutate: deletePost, isPending: isDeleting } = useDeletePost()
  const { mutate: updatePost, isPending: isUpdating } = useUpdatePost()

  // ── menu state ─────────────────────────────────────
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const openMenu = Boolean(anchorEl)
  const handleMenuClick = (e: React.MouseEvent<HTMLElement>) => setAnchorEl(e.currentTarget)
  const handleMenuClose = () => setAnchorEl(null)

  // ── dialog / modal state ───────────────────────────
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [editOpen, setEditOpen] = useState(false)

  // ── handlers ───────────────────────────────────────
  const handleDelete = () => {
    handleMenuClose()
    setDeleteOpen(true)
  }

  const handleDeleteConfirm = () => {
    deletePost(post.id)
    setDeleteOpen(false)
    // onDelete?.(post.id)
  }

  const handleEdit = () => {
    handleMenuClose()
    setEditOpen(true)
  }

  const handleSaveEdit = (data: { title: string; content: string; tags: string[] }) => {
    setEditOpen(false)
    // onUpdate?.(post.id, data)
    updatePost({ postID: post.id, updatedPost: data })
  }

  return (
    <>
      <Card elevation={3} sx={{ borderRadius: 3, p: 2, backgroundColor: 'background.paper' }}>
        <CardHeader
          avatar={
            // <Avatar src={user.avatar || ''} alt={user.username}>
            <Avatar src={''} alt={user.username}>
              {user.username[0].toUpperCase()}
            </Avatar>
          }
          title={
            <Typography variant="h6" fontWeight="bold">
              {post.title}
            </Typography>
          }
          subheader={
            <Typography variant="body2" color="text.secondary">
              {`by ${user.username} • ${new Date(post.created_at).toLocaleDateString()}`}
            </Typography>
          }
          action={
            isAuthor && (
              <>
                <IconButton
                  aria-label="post actions"
                  aria-controls={openMenu ? 'post-menu' : undefined}
                  aria-haspopup="true"
                  onClick={handleMenuClick}
                >
                  <MoreVertIcon />
                </IconButton>

                {/* Menu */}
                <Menu
                  id="post-menu"
                  anchorEl={anchorEl}
                  open={openMenu}
                  onClose={handleMenuClose}
                  anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                  transformOrigin={{ vertical: 'top', horizontal: 'right' }}
                >
                  <MenuItem onClick={handleEdit}>
                    <EditIcon fontSize="small" sx={{ mr: 1 }} />
                    Edit
                  </MenuItem>
                  <MenuItem onClick={handleDelete}>
                    <DeleteIcon fontSize="small" color="error" sx={{ mr: 1 }} />
                    Delete
                  </MenuItem>
                </Menu>
              </>
            )
          }
          sx={{ pb: 0 }}
        />

        <CardContent sx={{ pt: 1 }}>
          <Typography variant="body1" paragraph sx={{ whiteSpace: 'pre-line' }}>
            {post.content.length > 200 ? post.content.slice(0, 200) + '…' : post.content}
          </Typography>

          {post.tags?.length > 0 && (
            <Stack direction="row" spacing={1} sx={{ mb: 1, flexWrap: 'wrap' }}>
              {post.tags.map((tag, i) => (
                <Box key={i} sx={{ mb: 2 }}>
                  <Chip label={`#${tag}`} size="small" variant="outlined" />
                </Box>
              ))}
            </Stack>
          )}

          <Button
            component={RouterLink}
            to={`/post/${post.id}`}
            startIcon={<CommentIcon />}
            size="small"
            sx={{ textTransform: 'none' }}
          >
            {post.comments_count} comment{post.comments_count !== 1 ? 's' : ''}
          </Button>
        </CardContent>
      </Card>

      {/* Delete confirmation */}
      <DeleteConfirmDialog
        loading={isDeleting}
        open={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        onConfirm={handleDeleteConfirm}
      />

      {/* Edit post modal */}
      <EditPostModal
        loading={isUpdating}
        open={editOpen}
        onClose={() => setEditOpen(false)}
        onSave={handleSaveEdit}
        post={post}
      />
    </>
  )
}
