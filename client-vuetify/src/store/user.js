import { defineStore } from 'pinia'
import Cookies from 'js-cookie'

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
        Cookies.set('token', token, { expires: 7 })
      }
      if (user) {
        Cookies.set('user', JSON.stringify(user), { expires: 7 })
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
