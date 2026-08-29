import { defineStore } from 'pinia'
import api from '@/lib/api'

export const useInstallStore = defineStore('install', {
  state: () => ({
    checked: false,
    installed: true,
    hasDb: false,
    dbError: '',
    step: '',
    version: '',
    uploadWritable: true
  }),
  actions: {
    async fetchStatus() {
      try {
        const res = await api.get('/install/status')
        this.installed = !!res.data.installed
        this.hasDb = !!res.data.has_db
        this.dbError = res.data.db_error || ''
        this.step = res.data.step || ''
        this.version = res.data.version || ''
        this.uploadWritable = !!res.data.upload_writable
        this.checked = true
      } catch (e) {
        console.error('Failed to fetch install status:', e)
        this.checked = true
      }
      return this.installed
    }
  }
})
