import { useEffect, useState } from 'react'
import {
  Box,
  CircularProgress,
  Container,
  InputAdornment,
  Paper,
  Stack,
  TextField,
  Typography,
  Fab,
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import SearchIcon from '@mui/icons-material/Search'
import { useInView } from 'react-intersection-observer'

import PostCard from '../../components/PostCard'
import CreatePostModal from '../../components/CreatePostForm'
import FollowSuggestions from '../../components/FollowSuggestions'

import { useFeedInfinite } from '../../hooks/posts'
import useDebounce from '../../hooks/debounce'

export default function FeedPage() {
  const [openPostModal, setOpenPostModal] = useState(false)
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search)

  /* infinite scroll */
  const { ref, inView } = useInView({ rootMargin: '200px' })
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } =
    useFeedInfinite(debouncedSearch)

  const allPosts = data?.pages.flatMap((page) => page) ?? []

  useEffect(() => {
    if (inView && hasNextPage) fetchNextPage()
  }, [inView, hasNextPage, fetchNextPage])

  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" mt={6}>
        <CircularProgress />
      </Box>
    )
  }

  return (
    <>
      {/* create‑post modal */}
      <CreatePostModal open={openPostModal} onClose={() => setOpenPostModal(false)} />

      {/* floating FAB */}
      <Fab
        color="primary"
        aria-label="add"
        onClick={() => setOpenPostModal(true)}
        sx={{ position: 'fixed', bottom: 32, right: 32, zIndex: 10 }}
      >
        <AddIcon />
      </Fab>

      {/* main layout */}
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Box
          display="flex"
          flexDirection={{ xs: 'column', md: 'row' }}
          gap={4}
          justifyContent="center"
        >
          {/* Feed Section */}
          <Box flex={1}>
            <Stack spacing={3}>
              {/* Header + Search */}
              <Paper elevation={3} sx={{ p: 2, borderRadius: 3 }}>
                <Typography variant="h5" fontWeight="bold">
                  Your Feed
                </Typography>
                <Typography variant="body2" color="text.secondary" mb={2}>
                  Posts from people you follow
                </Typography>
                <TextField
                  autoFocus
                  fullWidth
                  size="small"
                  placeholder="Search posts…"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  sx={{ background: 'white' }}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <SearchIcon fontSize="small" />
                      </InputAdornment>
                    ),
                  }}
                />
              </Paper>

              {/* Posts */}
              {allPosts.map(({ post, comments_count }) => (
                <PostCard key={post.id} post={{ ...post, comments_count }} />
              ))}

              {!allPosts.length && (
                <Box display="flex" justifyContent="center" p={4}>
                  <Typography color="primary">
                    {search ? 'No results.' : 'No posts yet!'}
                  </Typography>
                </Box>
              )}

              {/* Infinite scroll sentinel */}
              <Box ref={ref} display="flex" justifyContent="center" mt={2}>
                {isFetchingNextPage && <CircularProgress />}
              </Box>
            </Stack>
          </Box>

          {debouncedSearch === '' && (
            <Box flexShrink={0} sx={{ width: 300, display: { xs: 'none', md: 'block' } }}>
              <FollowSuggestions />
            </Box>
          )}
        </Box>
      </Container>
    </>
  )
}
