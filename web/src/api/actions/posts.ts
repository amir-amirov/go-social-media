import fetchData from '..'

export type Post = {
  title: string
  content: string
  tags?: string[]
}

const getUserPosts = async (page: number, userID: number | undefined) => {
  if (!userID) {
    return
  }
  const { response, status } = await fetchData({
    url: `/posts/users/${userID}?limit=5&offset=${5 * (page - 1)}`,
    method: 'GET',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const getFeed = async (page: number, search: string) => {
  const params = new URLSearchParams({
    limit: '5',
    offset: String(5 * (page - 1)),
  })

  if (search) params.append('search', search)

  const { response, status } = await fetchData({
    url: `/users/feed?${params.toString()}`,
    method: 'GET',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const getPost = async (postID: number) => {
  const { response, status } = await fetchData({
    url: `/posts/${postID}`,
    method: 'GET',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const createPost = async (post: Post) => {
  const { response, status } = await fetchData({
    url: '/posts',
    method: 'POST',
    body: post,
  })

  if (status !== 201) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const updatePost = async (postID: number, post: Post) => {
  const { response, status } = await fetchData({
    url: `/posts/${postID}`,
    method: 'PATCH',
    body: post,
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const deletePost = async (postID: number) => {
  const { response, status } = await fetchData({
    url: `/posts/${postID}`,
    method: 'DELETE',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const posts = {
  getUserPosts,
  getFeed,
  getPost,
  createPost,
  updatePost,
  deletePost,
}

export default posts
