import { useState } from 'react'
import { Box, Typography, TextField, Button, Alert, CircularProgress } from '@mui/material'
import CenteredCard from '../../components/CenteredCard'
import { Link, useNavigate } from 'react-router-dom'
import { toast } from 'react-toastify'
import { useLogin } from '../../hooks/auth'

const Login: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)

  const loginMutation = useLogin()
  const navigate = useNavigate()

  const handleSubmit: React.FormEventHandler = (e) => {
    e.preventDefault()

    if (!email || !password) {
      setError('Please fill in all fields')
      toast.error('Please fill in all fields!')
      return
    }
    loginMutation.mutate(
      { email: email.toLowerCase(), password: password.toLowerCase() },
      {
        onSuccess: () => {
          navigate('/feed')
          toast.success('Logged in successfully')
        },
        onError: (error: any) => {
          setError(error?.message || 'Something went wrong')
          toast.error(error?.message || 'Something went wrong')
        },
      }
    )
  }

  return (
    <CenteredCard title="Welcome Back! 👋" subtitle="Log in with your credentials">
      {error && <Alert severity="error">{error}</Alert>}
      <Box component="form" mt={2} onSubmit={handleSubmit}>
        <TextField
          label="Email"
          type="email"
          variant="outlined"
          fullWidth
          margin="normal"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          disabled={loginMutation.isPending}
        />
        <TextField
          label="Password"
          type="password"
          variant="outlined"
          fullWidth
          margin="normal"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          disabled={loginMutation.isPending}
        />
        <Button
          variant="contained"
          color="primary"
          size="large"
          fullWidth
          sx={{ mt: 2, py: 1.2, borderRadius: 2 }}
          type="submit"
          disabled={loginMutation.isPending}
        >
          {loginMutation.isPending ? (
            <CircularProgress size={24} sx={{ color: 'white' }} />
          ) : (
            'Login'
          )}
        </Button>
        <Typography align="center" mt={2}>
          New here? <Link to="/signup">Create an account</Link>
        </Typography>
      </Box>
    </CenteredCard>
  )
}

export default Login
