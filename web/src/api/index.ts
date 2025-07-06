import { store } from '../store'
import { resetUser } from '../store/user/reducers'
import axios from './axios'
import { toast } from 'react-toastify'

interface FetchDataOptions {
  method?: string
  url: string
  body?: any
  headers?: any
  viewNotify?: boolean
  dontCanceled?: boolean | undefined
}

export const fetchData = async ({
  method = 'GET',
  headers,
  url,
  body,
  viewNotify = true,
  dontCanceled,
}: FetchDataOptions): Promise<any> => {
  console.log('REQUEST_DATA:', `URL: ${url}`, body)

  try {
    const state = store.getState()
    const { token } = state.user

    // if (token) {
    //   const expirationTime = tokenCreatedAt + expiresIn * 1000;
    //   const currentTimestamp = Date.now();
    //   if (expirationTime <= currentTimestamp) {
    //     // Токен просрочен
    //     console.log('Токен просрочен');
    //     await authApi.actionRefreshToken();
    //   }
    // }

    const response = await axios({
      headers: token
        ? { ...headers, Authorization: 'Bearer ' + token, dontCanceled }
        : { ...headers, dontCanceled },
      method,
      url,
      data: body,
    })

    console.log('RESPONSE_DATA:', `URL: ${url}`, response?.data)

    return { status: response.status, response: response?.data }
  } catch (error: any) {
    if (error.response?.status === 401) {
      store.dispatch(resetUser)
      // dispatchAction({type: "RESET_ALL"});
      return {
        status: error.response?.status ?? 403, // Если статус неизвестен, возвращаем 500
        response: error.response?.data,
      }
    }
    if (error !== 'dontShow') {
      // Проверка, что объект ошибки содержит ответ от сервера
      const errorMessage = error.response?.data?.message ?? error.message

      if (viewNotify) {
        // Попытка показать кастомный тост
        // toast.show({
        //   type: 'errorCustom',// 'errorCustom', 'error',
        //   text1: 'Внимание',
        //   text2: errorMessage,
        // });
        toast.error(errorMessage)
      }

      console.log('Error:', `URL: ${url}`, error)

      // Возвращение статуса ошибки и данных из ответа, если они есть
      return {
        status: error.response?.status ?? 500, // Если статус неизвестен, возвращаем 500
        response: error.response?.data,
      }
    }
    return {
      status: 500,
      response: '',
    }
  }
}

export default fetchData
