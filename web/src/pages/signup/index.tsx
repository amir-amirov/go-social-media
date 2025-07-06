import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import CenteredCard from '../../components/CenteredCard'
import { Alert, Box, Button, CircularProgress, TextField, Typography } from '@mui/material'
import { useSignIn } from '../../hooks/auth'

const SignUp: React.FC = () => {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  const signInMutation = useSignIn()

  const handleSubmit: React.FormEventHandler = (e) => {
    e.preventDefault()
    if (!username || !email || !password) {
      setError('All fields are required')
      return
    }

    signInMutation.mutate(
      {
        username: username.toLowerCase(),
        email: email.toLowerCase(),
        password: password.toLowerCase(),
      },
      {
        onSuccess: () => {
          navigate('/verify', { state: { email } })
        },
        onError: (error: any) => {
          if (error?.message == 'a user with that email already exists') {
            navigate('/verify', { state: { email } })
          }
          setError(error?.message || 'Something went wrong')
        },
      }
    )
  }

  return (
    <CenteredCard title="Create Account" subtitle="Sign up to get started">
      {error && <Alert severity="error">{error}</Alert>}
      <Box component="form" mt={2} onSubmit={handleSubmit}>
        <TextField
          label="Username"
          variant="outlined"
          fullWidth
          margin="normal"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          disabled={signInMutation.isPending}
        />
        <TextField
          label="Email"
          type="email"
          variant="outlined"
          fullWidth
          margin="normal"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          disabled={signInMutation.isPending}
        />
        <TextField
          label="Password"
          type="password"
          variant="outlined"
          fullWidth
          margin="normal"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          disabled={signInMutation.isPending}
        />
        <Button
          variant="contained"
          color="primary"
          size="large"
          fullWidth
          sx={{ mt: 2, py: 1.2, borderRadius: 2 }}
          type="submit"
          disabled={signInMutation.isPending}
        >
          {signInMutation.isPending ? (
            <CircularProgress size={24} sx={{ color: 'white' }} />
          ) : (
            'Continue'
          )}
        </Button>
        <Typography align="center" mt={2}>
          Already have an account? <Link to="/login">Log in</Link>
        </Typography>
      </Box>
    </CenteredCard>
  )
}

export default SignUp
