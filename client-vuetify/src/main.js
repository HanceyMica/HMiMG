/**
 * Vue 应用入口文件
 *
 * 本文件是 Vue 3 应用的起点，负责：
 * - 创建 Vue 应用实例
 * - 注册各种插件（如 Vuetify、Vue Router、Vue I18n 等）
 * - 将应用挂载到 DOM 中的 #app 元素
 *
 * 执行流程：
 * 1. 导入 createApp 函数创建应用实例
 * 2. 导入根组件 App.vue
 * 3. 导入插件注册函数 registerPlugins
 * 4. 调用 registerPlugins(app) 注册所有插件
 * 5. 将应用挂载到 #app 元素
 */

import { createApp } from 'vue'
// 导入根组件 - App.vue 是应用的顶层组件，所有其他组件都在其内部渲染
import App from './App.vue'
// 导入插件注册函数 - 集中管理所有插件的注册
import { registerPlugins } from '@/plugins'

// 创建 Vue 3 应用实例
const app = createApp(App)

// 调用插件注册函数，传入 app 实例
// registerPlugins 内部会依次注册 Vuetify、Vue Router、Vue I18n 等插件
registerPlugins(app)

// 将应用挂载到 DOM 中的 #app 元素
// 这会触发应用的挂载流程，开始渲染组件树
app.mount('#app')
