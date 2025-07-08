import axios from 'axios'

const headers = {
  //   'Accept-Language': 'RU',
  'Content-Type': 'application/json',
}

// Используем Map для хранения AbortController
const abortControllers = new Map<string, AbortController>()

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: headers,
})

axiosInstance.interceptors.request.use((config) => {
  if (!config.headers?.dontCanceled) {
    const urlRequest = config.url?.split('?')[0]
    const key = `${config.method}-${urlRequest}`

    // Отменяем предыдущий запрос с тем же ключом
    if (abortControllers.has(key)) {
      abortControllers.get(key)?.abort()
    }

    // Создаем новый AbortController
    const controller = new AbortController()
    config.signal = controller.signal

    // Сохраняем AbortController в Map
    abortControllers.set(key, controller)
  }

  return config
})

axiosInstance.interceptors.response.use(
  (response) => {
    if (!response.config.headers?.dontCanceled) {
      const urlRequest = response.config.url?.split('?')[0]
      const key = `${response.config.method}-${urlRequest}`
      abortControllers.delete(key)
    }

    return response
  },
  (error) => {
    // Проверяем наличие config в ошибке
    if (error.config && !error.config.headers?.dontCanceled) {
      const urlRequest = error.config.url?.split('?')[0]
      const key = `${error.config.method}-${urlRequest}`
      abortControllers.delete(key)
    }

    if (axios.isCancel(error)) {
      console.log('Request canceled:', error.message)
      return Promise.reject('dontShow')
    }

    // Возвращаем ошибку, если это не отмена
    return Promise.reject(error)
  }
)

export default axiosInstance
