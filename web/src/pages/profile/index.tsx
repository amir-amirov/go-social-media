import { useEffect } from 'react'
import Typography from '@mui/material/Typography'
import CircularProgress from '@mui/material/CircularProgress'
import Box from '@mui/material/Box'
import Stack from '@mui/material/Stack'
import PostCard from '../../components/PostCard'
import { useUser } from '../../store/user'
import { useInView } from 'react-intersection-observer'
import { Avatar, Container, Paper } from '@mui/material'
import { useUserPostsInfinite } from '../../hooks/posts'
import { useGetFollowStats } from '../../hooks/users'

export default function ProfilePage() {
  const { user } = useUser()
  const { ref, inView } = useInView({
    rootMargin: '200px',
  })

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = useUserPostsInfinite(
    user?.id
  )

  const { data: followStats, isPending: isFetchingFollowStats } = useGetFollowStats(user?.id)

  const allPosts = data?.pages?.flatMap((page) => page)

  useEffect(() => {
    if (inView && hasNextPage) {
      fetchNextPage()
    }
  }, [inView, hasNextPage, fetchNextPage])

  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" mt={6}>
        <CircularProgress />
      </Box>
    )
  }

  if (!data)
    return (
      <Container maxWidth="xs" sx={{ display: 'flex', justifyContent: 'center', mt: 6 }}>
        <Typography color="textPrimary">User not found.</Typography>
      </Container>
    )

  return (
    <Container maxWidth="xs">
      {/* USER INFO CARD */}
      <Paper
        elevation={3}
        sx={{
          p: 3,
          mb: 4,
          borderRadius: 4,
          backgroundColor: '#f5f8fa',
          display: 'flex',
          alignItems: 'center',
          gap: 2,
        }}
      >
        <Avatar
          src={user?.avatar || undefined}
          alt={user?.username}
          sx={{ width: 64, height: 64 }}
        />
        <Box>
          <Typography variant="h6" color="textPrimary">
            {user?.username}
          </Typography>
          {!isFetchingFollowStats && (
            <Typography variant="body2" color="textSecondary">
              {followStats?.followers ?? 0} followers • {followStats?.following ?? 0} following
            </Typography>
          )}
        </Box>
      </Paper>

      <Typography variant="h6" mt={4} mb={2} color="black">
        Posts
      </Typography>

      <Stack spacing={3}>
        {allPosts && allPosts.length ? (
          allPosts?.map(({ post, comments_count }) => (
            <PostCard key={post.id} post={{ ...post, comments_count }} />
          ))
        ) : (
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              height: '100%',
              width: '100%',
            }}
          >
            <Typography color="primary">No posts yet</Typography>
          </Box>
        )}

        <Box ref={ref} display="flex" justifyContent="center" mt={2}>
          {isFetchingNextPage && <CircularProgress />}
        </Box>
      </Stack>
    </Container>
  )
}
