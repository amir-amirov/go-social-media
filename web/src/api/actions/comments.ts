import fetchData from '..'

type CreateCommentPayload = {
  content: string
}

const createComment = async (postID: number | undefined, comment: string) => {
  if (!postID) {
    throw new Error('invalid post id')
  }
  let payload: CreateCommentPayload = {
    content: comment,
  }
  const { response, status } = await fetchData({
    url: `/posts/${postID}/comments`,
    method: 'POST',
    body: payload,
  })

  if (status !== 201) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const deleteComment = async (postID: number | undefined, commentID: number | undefined) => {
  if (!postID || !commentID) {
    throw new Error('invalid id of post or comment')
  }
  const { response, status } = await fetchData({
    url: `/posts/${postID}/comments/${commentID}`,
    method: 'DELETE',
  })

  if (status !== 200) {
    throw new Error(response?.data?.message)
  }

  return response.data
}

const comments = {
  createComment,
  deleteComment,
}

export default comments
