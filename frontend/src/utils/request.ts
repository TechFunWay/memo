import axios from 'axios'

const api = axios.create({
  baseURL: `${import.meta.env.BASE_URL}api`,
  timeout: 10000,
})

// 请求拦截器
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = token
  }
  return config
})

// 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const hadToken = Boolean(localStorage.getItem('token'))
    const requestUrl = String(error.config?.url || '')
    const isLoginRequest = ['/auth/login', '/auth/fnos/login', '/auth/fnos/bind'].some((path) => requestUrl.endsWith(path))
    if (error.response?.status === 401 && hadToken && !isLoginRequest) {
      localStorage.removeItem('token')
      window.location.href = `${import.meta.env.BASE_URL}login`
    }
    return Promise.reject(error)
  }
)

export default api
