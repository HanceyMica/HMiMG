import axios from 'axios'
import Cookies from 'js-cookie'

const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:9108/api'

const api = axios.create({
  baseURL: apiBaseUrl,
  timeout: 30000,
})

api.interceptors.request.use((config) => {
  const token = Cookies.get('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status
    const url = error.config?.url || ''
    if (status === 401 && !url.includes('/login') && !url.includes('/register')) {
      Cookies.remove('token')
      Cookies.remove('user')
      const current = window.location.pathname + window.location.search
      if (window.location.pathname !== '/login') {
        window.location.href = `/login?redirect=${encodeURIComponent(current)}&reason=unauthorized`
      }
    }
    return Promise.reject(error)
  }
)

export const buildUploadedFileUrl = (path) => {
  const normalizedPath = String(path || '')
    .split('/')
    .filter(Boolean)
    .map(segment => encodeURIComponent(segment))
    .join('/')

  return normalizedPath ? `${apiBaseUrl.replace(/\/$/, '')}/files/${normalizedPath}` : ''
}

export default api
