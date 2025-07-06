import {
  Avatar,
  Box,
  Button,
  Divider,
  IconButton,
  Paper,
  Skeleton,
  Stack,
  Typography,
} from '@mui/material'
import CloseIcon from '@mui/icons-material/Close'
import { useState } from 'react'
import { useFollowUser, useTopUsers, useUnFollowUser } from '../../hooks/users'

export default function FollowSuggestions() {
  const { data, isLoading } = useTopUsers()
  const [dismissed, setDismissed] = useState<number[]>([])
  const { mutate: follow } = useFollowUser()
  const { mutate: unfollow } = useUnFollowUser()

  // show only first 3 not‑yet‑dismissed users
  const suggestions =
    data?.filter((u: any) => !dismissed.includes(u.id) && !u.is_followed).slice(0, 3) ?? []

  const handleDismiss = (id: number) => setDismissed((prev) => [...prev, id])

  return (
    <Paper elevation={3} sx={{ p: 2, borderRadius: 3, position: 'sticky', top: 80 }}>
      <Typography variant="h6" fontWeight="bold" gutterBottom>
        Add to your feed
      </Typography>

      {isLoading &&
        [1, 2, 3].map((i) => (
          <Box key={i} mb={2}>
            <Skeleton variant="circular" width={40} height={40} />
            <Skeleton variant="text" sx={{ fontSize: 14 }} width={120} />
            <Skeleton variant="rectangular" height={32} sx={{ borderRadius: 2, mt: 1 }} />
            <Divider sx={{ my: 2 }} />
          </Box>
        ))}

      {suggestions && suggestions.length > 0 ? (
        suggestions.map((u: any) => (
          <Box key={u.id} mb={2}>
            <Stack direction="row" spacing={1} alignItems="center">
              <Avatar src={u.avatar || undefined}>{u.username[0].toUpperCase()}</Avatar>
              <Box flexGrow={1}>
                <Typography fontWeight="medium" noWrap>
                  {u.username}
                </Typography>
                <Typography variant="caption" color="text.secondary" noWrap>
                  {u.headline}
                </Typography>
              </Box>

              <IconButton size="small" onClick={() => handleDismiss(u.id)}>
                <CloseIcon fontSize="inherit" />
              </IconButton>
            </Stack>

            <Button
              fullWidth
              variant={u.is_followed ? 'outlined' : 'contained'}
              startIcon={u.is_followed ? undefined : '+'}
              sx={{ mt: 1, textTransform: 'none' }}
              onClick={() =>
                u.is_followed
                  ? unfollow(u.id)
                  : follow(u.id, { onSuccess: () => (u.is_followed = true) })
              }
              disabled={u.is_followed}
            >
              {u.is_followed ? 'Following' : 'Follow'}
            </Button>

            <Divider sx={{ mt: 2 }} />
          </Box>
        ))
      ) : isLoading ? (
        <></>
      ) : (
        <Box>
          <Typography>No users</Typography>
        </Box>
      )}
    </Paper>
  )
}
