import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Login from './pages/login'
import SignUp from './pages/signup'
import VerifyEmail from './pages/verify-email'
import MainLayout from './layouts/MainLayout'
import ProfilePage from './pages/profile'
import PostPage from './pages/post-page'
import FeedPage from './pages/feed'
import NotFoundPage from './pages/404'
import ProtectedRoutes from './components/ProtectedRoutes/ProtectedRoutes'

function App() {
  return (
    <Router>
      <Routes>
        {/* Authentication routes without navbar */}
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/verify" element={<VerifyEmail />} />

        {/* Main app routes with navbar */}
        <Route element={<ProtectedRoutes />}>
          <Route path="/" element={<MainLayout />}>
            <Route index element={<FeedPage />} />
            <Route path="feed" element={<FeedPage />} />
            <Route path="profile/:username" element={<ProfilePage />} />
            <Route path="post/:id" element={<PostPage />} />
          </Route>
        </Route>

        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </Router>
  )
}

export default App
