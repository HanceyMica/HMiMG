import { defineStore } from 'pinia'
import api from '@/lib/api'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    websiteTitle: 'HMiMG',
    defaultLanguage: 'zh',
    allowRegistration: true
  }),
  actions: {
    async fetchPublicSettings() {
      try {
        const res = await api.get('/settings/public')
        if (res.data.website_title) {
          this.websiteTitle = res.data.website_title
          document.title = res.data.website_title
        }
        if (res.data.default_language) {
          this.defaultLanguage = res.data.default_language
        }
        if (res.data.allow_registration !== undefined) {
          this.allowRegistration = res.data.allow_registration
        }
      } catch (e) {
        console.error('Failed to fetch settings:', e)
      }
    },
    setWebsiteTitle(title) {
      this.websiteTitle = title
      document.title = title
    }
  }
})
