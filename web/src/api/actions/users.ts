import fetchData from '..'

const getTopUsers = async () => {
  const { response, status } = await fetchData({
    url: `/users/top`,
    method: 'GET',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const followUser = async (userID: number) => {
  const { response, status } = await fetchData({
    url: `/users/${userID}/follow`,
    method: 'PUT',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const unfollowUser = async (userID: number) => {
  const { response, status } = await fetchData({
    url: `/users/${userID}/unfollow`,
    method: 'PUT',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const getFollowStats = async (userID: number | undefined) => {
  if (!userID) return
  const { response, status } = await fetchData({
    url: `/users/${userID}/follow-stats`,
    method: 'GET',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const users = {
  getTopUsers,
  followUser,
  unfollowUser,
  getFollowStats,
}

export default users
