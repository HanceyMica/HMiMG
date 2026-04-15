/**
 * 路由配置文件
 *
 * 本文件负责配置 Vue Router 路由管理，包括：
 * - 定义应用的所有路由路径
 * - 设置路由守卫进行身份验证
 * - 配置懒加载路由组件以优化性能
 *
 * 路由结构：
 * - 根路径 '/' 使用 MainLayout 布局组件，包含需要认证的子路由
 * - 登录页面 '/login' 为独立路由
 * - 所有未匹配的路径显示 404 页面
 */

import { createRouter, createWebHistory } from 'vue-router'
import Cookies from 'js-cookie'

/**
 * 路由配置数组
 *
 * 每个路由对象包含：
 * - path: 路由路径
 * - component: 路由组件（使用懒加载优化初始加载性能）
 * - meta: 元信息，用于存储额外数据（如 requiresAuth 标记需要认证的路由）
 * - children: 子路由数组，仅适用于布局路由
 */
const routes = [
  {
    // 根路径 - 使用主布局组件，需要身份验证
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    // meta.requiresAuth: true 表示访问此路由需要登录
    meta: { requiresAuth: true },
    // 子路由 - 这些路由会在 MainLayout 的 <router-view /> 中渲染
    children: [
      { path: '', component: () => import('@/pages/Home.vue') },                          // 首页
      { path: 'library', component: () => import('@/pages/Library.vue') },                // 图库页面
      { path: 'album/:id', component: () => import('@/pages/Album.vue') },                // 相册详情页，:id 为动态参数
      { path: 'collection/:id', component: () => import('@/pages/Collection.vue') },     // 合集详情页，:id 为动态参数
      { path: 'image/:id', component: () => import('@/pages/ImageDetails.vue') },         // 图片详情页，:id 为动态参数
      { path: 'admin', component: () => import('@/pages/Admin.vue') },                    // 后台管理页面
      { path: 'about', component: () => import('@/pages/About.vue') },                    // 关于页面
    ],
  },
  {
    // 登录页面 - 独立路由，不使用主布局
    path: '/login',
    component: () => import('@/pages/Login.vue'),
  },
  {
    // 404 页面 - 捕获所有未匹配的路由路径
    // vue-router 特有的动态路径语法，匹配任意路径
    path: '/:pathMatch(.*)*',
    component: () => import('@/pages/NotFound.vue'),
  },
]

/**
 * 创建 Vue Router 实例
 *
 * createWebHistory(): 使用 HTML5 History API，生成干净的 URL 路径
 * 例如: /album/1 而不是 /#/album/1
 */
const router = createRouter({
  history: createWebHistory(),
  routes,
})

/**
 * 全局前置路由守卫
 *
 * 每次路由跳转前都会执行此守卫
 * 用于检查用户是否已登录（通过 Cookie 中的 token 判断）
 * 如果用户未登录且访问需要认证的路由，则重定向到登录页面
 *
 * @param {Route} to - 目标路由对象，包含路由路径、参数等信息
 * @param {Route} from - 来源路由对象
 * @param {Function} next - 允许访问路由的回调函数
 */
router.beforeEach((to, from, next) => {
  // 从 Cookie 中获取登录令牌
  const token = Cookies.get('token')

  // 如果目标路由需要认证（meta.requiresAuth 为 true）但用户未登录
  if (to.meta.requiresAuth && !token) {
    // 重定向到登录页面，并携带原始目标路径和未授权原因
    // 这样登录成功后可以跳回原页面
    next({ path: '/login', query: { redirect: to.fullPath, reason: 'unauthorized' } })
  } else {
    // 允许访问目标路由
    next()
  }
})

// 导出路由实例，供 main.js 中注册到 Vue 应用
export default router
