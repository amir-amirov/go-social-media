import { Outlet } from 'react-router-dom'
import Navbar from '../components/Navbar'
import Container from '@mui/material/Container'
import theme from '../theme'
import { Box } from '@mui/material'

export default function MainLayout() {
  return (
    <Box sx={{ background: theme.palette.background, minHeight: '100vh' }}>
      <Navbar />
      <Container maxWidth="md" sx={{ pt: 10, height: '100%' }}>
        <Outlet />
      </Container>
    </Box>
  )
}
