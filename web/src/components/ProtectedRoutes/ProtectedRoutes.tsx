import { Outlet, Navigate } from 'react-router-dom'
import { useUser } from '../../store/user'

const ProtectedRoutes = () => {
  const { token } = useUser()
  return token ? <Outlet /> : <Navigate to="/login" />
}

export default ProtectedRoutes
