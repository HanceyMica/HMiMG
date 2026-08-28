/**
 * 网站设置状态管理模块
 *
 * 本模块负责管理网站的公共设置信息，包括：
 * - 网站标题 (websiteTitle)
 * - 默认语言 (defaultLanguage)
 * - 是否允许新用户注册 (allowRegistration)
 *
 * 这些设置属于公开信息，所有访客都可以获取，通常用于：
 * - 展示网站名称
 * - 设置多语言支持的默认语言
 * - 控制注册功能的可见性
 *
 * 设置数据从后端 API (/settings/public) 获取，并缓存在前端状态中。
 */

import { defineStore } from 'pinia'
// 导入 API 封装模块，用于发送 HTTP 请求到后端
import api from '@/lib/api'

/**
 * 使用 defineStore 创建 settings 状态库
 *
 * defineStore 接受两个参数：
 * 1. 状态库的唯一标识名称（字符串）
 * 2. 一个对象或函数，定义状态、getters 和 actions
 */
export const useSettingsStore = defineStore('settings', {
  /**
   * state 函数返回一个对象，定义状态的初始值
   * 这里定义了三个公开设置项
   */
  state: () => ({
    // 网站标题，用于显示在浏览器标签页和网站各处
    websiteTitle: 'HMiMG',
    // 默认语言代码，'zh' 表示中文，'en' 表示英文等
    defaultLanguage: 'zh',
    // 是否允许新用户注册，true 表示开放注册，false 表示关闭注册
    allowRegistration: true
  }),

  /**
   * actions 对象包含所有修改状态的方法（类似 Vuex 中的 mutations，但更灵活）
   * actions 可以是普通函数或异步函数
   */
  actions: {
    /**
     * 从后端获取公开设置信息
     *
     * 此方法会调用后端 API 获取最新的网站设置，
     * 并更新本地状态。如果请求失败，会在控制台输出错误信息。
     *
     * 使用场景：
     * - 应用启动时加载网站设置
     * - 用户访问设置页面时刷新设置
     * - 网站标题等变更后同步更新
     */
    async fetchPublicSettings() {
      try {
        // 发送 GET 请求到后端的 /settings/public 端点
        const res = await api.get('/settings/public')

        // 如果后端返回了网站标题，则更新本地状态
        if (res.data.website_title) {
          this.websiteTitle = res.data.website_title
        }

        // 如果后端返回了默认语言设置，则更新本地状态
        if (res.data.default_language) {
          this.defaultLanguage = res.data.default_language
        }

        // 如果后端返回了注册允许设置（注意：需要检查 undefined，因为可能是 false）
        if (res.data.allow_registration !== undefined) {
          this.allowRegistration = res.data.allow_registration
        }
      } catch (e) {
        // 请求失败时打印错误信息，便于调试
        console.error('Failed to fetch settings:', e)
      }
    },

    /**
     * 设置网站标题
     *
     * 用于在运行时动态更新网站标题，
     * 例如管理员在后台修改网站名称后，前端可以调用此方法实时更新。
     *
     * @param {string} title - 新的网站标题
     */
    setWebsiteTitle(title) {
      this.websiteTitle = title
    }
  }
})
