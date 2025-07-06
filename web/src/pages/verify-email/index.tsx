import { useEffect, useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import CenteredCard from '../../components/CenteredCard'
import { Alert, Box, Button, CircularProgress, Typography } from '@mui/material'
import { toast } from 'react-toastify'
import { useVerifyEmail } from '../../hooks/auth'
import ReactCodeInput from 'react-code-input'

const VerifyEmail: React.FC = () => {
  const [code, setCode] = useState('')
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  const { state } = useLocation()

  const email = state?.email as string | undefined

  const verifyEmailMutation = useVerifyEmail()

  useEffect(() => {
    if (!email) {
      navigate('/signup')
    }
  }, [email, navigate])

  const handleSubmit: React.FormEventHandler = (e) => {
    e.preventDefault()
    if (!code || code.length != 6) {
      setError('Verification code is required')
      toast.error('Verification code is required')
      return
    }
    verifyEmailMutation.mutate(code, {
      onSuccess: () => {
        navigate('/login')
        toast.success('Account created successfully!')
      },
      onError: (error: any) => {
        setError(error?.message || 'Something went wrong')
      },
    })
  }

  useEffect(() => {
    if (code.length === 6 && !verifyEmailMutation.isPending) {
      handleSubmit({ preventDefault: () => {} } as React.FormEvent)
    }
  }, [code])

  return (
    <CenteredCard
      title="Enter Code"
      subtitle={`We have sent a 6‑digit code to your email ${email ?? ''}`}
    >
      {error && <Alert severity="error">{error}</Alert>}
      <Box component="form" mt={2} onSubmit={handleSubmit}>
        <Box display="flex" justifyContent="center" mt={2} mb={1}>
          <ReactCodeInput
            type="tel"
            fields={6}
            onChange={(value) => setCode(value)}
            disabled={verifyEmailMutation.isPending}
            inputStyle={{
              fontFamily: 'sans-serif',
              margin: '0 6px',
              MozAppearance: 'textfield',
              width: '40px',
              borderRadius: '8px',
              fontSize: '24px',
              height: '48px',
              paddingLeft: '7px',
              backgroundColor: 'white',
              color: '#000',
              border: '1px solid #ccc',
              textAlign: 'center',
            }}
            name={''}
            inputMode={'numeric'}
          />
        </Box>

        <Button
          variant="contained"
          color="primary"
          size="large"
          fullWidth
          sx={{ mt: 2, py: 1.2, borderRadius: 2 }}
          type="submit"
          disabled={verifyEmailMutation.isPending}
        >
          {verifyEmailMutation.isPending ? (
            <CircularProgress size={24} sx={{ color: 'white' }} />
          ) : (
            'Verify'
          )}
        </Button>
        <Typography align="center" mt={2}>
          Didn't receive the code? <Link to="/signup">Resend</Link>
        </Typography>
      </Box>
    </CenteredCard>
  )
}
export default VerifyEmail
