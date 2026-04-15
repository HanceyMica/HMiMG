/**
 * Pinia 状态管理库初始化模块
 *
 * 本文件负责创建和导出 Pinia 实例，这是 Vue 3 应用程序状态管理的核心。
 *
 * Pinia 是 Vue 官方推荐的新一代状态管理库，相比 Vuex 更加轻量且 API 更加简洁。
 * 它提供了更好的 TypeScript 支持、更直观的 API，以及自动的代码分割能力。
 *
 * 使用方法：
 * 在 main.js 中通过 app.use(pinia) 将 Pinia 实例挂载到 Vue 应用上。
 * 然后在各组件中通过 import { useXxxStore } from '@/store' 来获取状态库。
 */

// 从 pinia 包中导入 createPinia 函数，用于创建 Pinia 实例
import { createPinia } from 'pinia'

// 创建 Pinia 实例
const pinia = createPinia()

// 导出 Pinia 实例，供 main.js 挂载到 Vue 应用
export default pinia
