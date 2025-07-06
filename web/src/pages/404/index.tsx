import { Box, Typography, Button } from '@mui/material'
import { Link as RouterLink } from 'react-router-dom'
import theme from '../../theme'

export default function NotFoundPage() {
  return (
    <Box
      height="100vh"
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="center"
      textAlign="center"
      px={2}
      bgcolor={theme.palette.background}
    >
      <Typography variant="h1" fontWeight="bold" color="primary" gutterBottom>
        404
      </Typography>
      <Typography variant="h5" gutterBottom color="primary">
        Oops! Page not found.
      </Typography>
      <Typography variant="body1" color="textSecondary" mb={4}>
        The page you’re looking for doesn’t exist or has been moved.
      </Typography>
      <Button variant="contained" component={RouterLink} to="/" size="large">
        Go to Home
      </Button>
    </Box>
  )
}
