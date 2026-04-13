import axios from 'axios'
import Cookies from 'js-cookie'

const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:9108/api'

const api = axios.create({
  baseURL: apiBaseUrl,
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  const token = Cookies.get('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const buildUploadedFileUrl = (path) => {
  const normalizedPath = String(path || '')
    .split('/')
    .filter(Boolean)
    .map(segment => encodeURIComponent(segment))
    .join('/')

  return normalizedPath ? `${apiBaseUrl.replace(/\/$/, '')}/files/${normalizedPath}` : ''
}

export default api
