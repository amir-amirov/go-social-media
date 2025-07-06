import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Avatar,
  Menu,
  MenuItem,
  IconButton,
  Box,
} from '@mui/material'
import { Link as RouterLink, useNavigate } from 'react-router-dom'
import { useState } from 'react'
import { useUser } from '../../store/user'
import LogoutIcon from '@mui/icons-material/Logout'
import PersonIcon from '@mui/icons-material/Person'
import { useQueryClient } from '@tanstack/react-query'

export default function Navbar() {
  const { user, resetUser } = useUser()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const avatar = user?.avatar || '/default-avatar.png' // fallback

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
  const open = Boolean(anchorEl)

  const handleAvatarClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleClose = () => {
    setAnchorEl(null)
  }

  const handleProfile = () => {
    navigate(`/profile/${user?.username || 'me'}`)
    handleClose()
  }

  const handleLogout = () => {
    resetUser()
    queryClient.clear()
    navigate('/login')
    handleClose()
  }

  return (
    <AppBar position="fixed">
      <Toolbar>
        <Typography variant="h6" sx={{ flexGrow: 1 }} onClick={() => navigate('/feed')}>
          HABAR
        </Typography>

        <Button
          color="inherit"
          component={RouterLink}
          to="/feed"
          sx={{
            color: 'white',
            '&:hover': {
              color: 'lightblue',
            },
          }}
        >
          Feed
        </Button>

        {user && (
          <Box>
            <IconButton onClick={handleAvatarClick} sx={{ p: 0 }}>
              <Avatar src={avatar} alt="User avatar" />
            </IconButton>
            <Menu
              anchorEl={anchorEl}
              open={open}
              onClose={handleClose}
              onClick={handleClose}
              PaperProps={{
                elevation: 4,
                sx: {
                  mt: 1.5,
                  overflow: 'visible',
                  filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.15))',
                  minWidth: 140,
                },
              }}
              transformOrigin={{ horizontal: 'right', vertical: 'top' }}
              anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
            >
              <MenuItem onClick={handleProfile}>
                <PersonIcon fontSize="small" sx={{ mr: 1 }} />
                Profile
              </MenuItem>
              <MenuItem onClick={handleLogout} sx={{ color: 'red' }}>
                <LogoutIcon fontSize="small" sx={{ mr: 1 }} />
                Log out
              </MenuItem>
            </Menu>
          </Box>
        )}
      </Toolbar>
    </AppBar>
  )
}
