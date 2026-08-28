import { defineStore } from 'pinia'
import Cookies from 'js-cookie'

const cookieOptions = () => ({
  expires: 7,
  sameSite: 'lax',
  secure: window.location.protocol === 'https:',
})

export const useUserStore = defineStore('user', {
  state: () => ({
    user: Cookies.get('user') ? JSON.parse(Cookies.get('user')) : null,
    token: Cookies.get('token') || null,
  }),
  actions: {
    setUser(user, token) {
      this.user = user
      this.token = token
      if (token) {
        Cookies.set('token', token, cookieOptions())
      }
      if (user) {
        Cookies.set('user', JSON.stringify(user), cookieOptions())
      }
    },
    logout() {
      this.user = null
      this.token = null
      Cookies.remove('token')
      Cookies.remove('user')
    },
  },
})
