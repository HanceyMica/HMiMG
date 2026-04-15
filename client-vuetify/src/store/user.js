/**
 * 用户状态管理模块
 *
 * 本模块负责管理用户的登录状态和认证信息，包括：
 * - 用户对象 (user) - 存储登录用户的详细信息
 * - 认证令牌 (token) - 用于 API 请求的身份验证 token
 *
 * 用户数据通过 Cookies 进行持久化存储，这样在页面刷新后
 * 用户的登录状态可以自动恢复，无需重新登录。
 *
 * 注意： Cookies 的有效期设置为 7 天，过期后用户需要重新登录。
 */

import { defineStore } from 'pinia'
// 导入 js-cookie 库，用于操作浏览器 Cookies
import Cookies from 'js-cookie'

/**
 * 使用 defineStore 创建 user 状态库
 *
 * 该状态库管理当前用户的登录状态和认证凭证
 */
export const useUserStore = defineStore('user', {
  /**
   * state 函数返回初始状态
   *
   * 初始化时从 Cookies 中读取用户信息和 token：
   * - 如果 Cookies 中存在 'user'，则解析为用户对象
   * - 如果 Cookies 中存在 'token'，则作为认证令牌
   *
   * 使用三目运算符确保在数据不存在时返回合适的默认值
   */
  state: () => ({
    /**
     * 当前登录用户对象
     * 如果 Cookies 中有保存的用户信息则解析使用，否则为 null
     * 用户对象通常包含 id, username, email, avatar 等字段
     */
    user: Cookies.get('user') ? JSON.parse(Cookies.get('user')) : null,

    /**
     * 认证令牌 token
     * 用于在 API 请求中验证用户身份
     * 如果 Cookies 中存在则读取，否则为 null
     */
    token: Cookies.get('token') || null,
  }),

  /**
   * actions 对象包含所有修改状态的方法
   */
  actions: {
    /**
     * 设置用户信息和认证令牌
     *
     * 当用户登录成功时调用此方法，保存用户数据和 token。
     * 数据会同时存储到 Pinia 状态和 Cookies 中，
     * 以实现页面刷新后登录状态的持久化。
     *
     * @param {Object} user - 用户对象，包含用户的详细信息
     * @param {string} token - 认证令牌，用于 API 请求身份验证
     *
     * Cookies 存储选项：
     * - expires: 7 表示 cookie 在 7 天后过期
     * - 这确保用户在 7 天内无需重新登录
     */
    setUser(user, token) {
      // 更新 Pinia 状态
      this.user = user
      this.token = token

      // 如果提供了 token，则存储到 Cookie
      if (token) {
        Cookies.set('token', token, { expires: 7 })
      }

      // 如果提供了用户对象，则将其 JSON 序列化后存储到 Cookie
      // JSON.stringify 是必要的，因为 Cookie 只支持字符串类型
      if (user) {
        Cookies.set('user', JSON.stringify(user), { expires: 7 })
      }
    },

    /**
     * 用户登出
     *
     * 清除用户的所有认证信息，包括：
     * - 清空 Pinia 状态中的 user 和 token
     * - 删除 Cookie 中保存的 token 和 user
     *
     * 调用此方法后，用户需要重新登录才能访问受保护的资源。
     */
    logout() {
      // 清空 Pinia 状态
      this.user = null
      this.token = null

      // 删除 Cookie 中的认证信息
      Cookies.remove('token')
      Cookies.remove('user')
    },
  },
})
