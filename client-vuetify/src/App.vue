<!--
  App.vue - Vue 应用根组件

  本组件是整个 Vue 应用的顶层容器，负责：
  - 提供 Vuetify 应用的基本结构 (<v-app>)
  - 作为路由视图的容器，显示当前路由对应的组件
  - 管理全局主题（明/暗模式）的初始化和切换
  - 设置全局背景样式和毛玻璃效果
  - 监听系统主题变化，自动切换应用主题
-->
<template>
  <!-- v-app: Vuetify 提供的应用容器组件 -->
  <!-- 它会自动处理应用的主题、布局和全局样式 -->
  <v-app>
    <!-- router-view: Vue Router 提供的路由视图组件 -->
    <!-- 根据当前 URL 显示对应路由的组件 -->
    <router-view />
  </v-app>
</template>

<script setup>
/**
 * 响应式API相关导入
 * - watch: 用于监听响应式数据的变化
 * - onMounted: 组件挂载完成后的生命周期钩子
 * - onUnmounted: 组件卸载前的生命周期钩子
 */
import { watch, onMounted, onUnmounted } from 'vue'

// Vuetify 主题管理hook
import { useTheme } from 'vuetify'

// Vue Router 路由hook，用于获取当前路由信息
import { useRoute } from 'vue-router'

// 获取 Vuetify 主题实例，用于操作主题切换
const theme = useTheme()

// 获取当前路由对象，可监听路由变化
const route = useRoute()

/**
 * 更新 body 元素的样式类
 *
 * 根据当前主题（明/暗）设置 body 的 CSS 类，
 * 以便应用全局背景图片和颜色
 *
 * 注意：在 SSR 环境下（服务端渲染）document 对象不存在，
 * 需要进行判断避免报错
 */
const updateBodyStyles = () => {
  // 服务端渲染时直接返回，避免报错
  if (typeof document === 'undefined') return

  // 获取当前主题是否为暗色模式
  const isDark = theme.global.current.value.dark

  // 移除可能存在的旧样式类
  document.body.classList.remove('light-mode', 'dark-mode')
  document.documentElement.classList.remove('light-mode', 'dark-mode')

  // 添加与当前主题对应的样式类（html 上的镜像供窗口滚动条样式使用）
  document.body.classList.add(isDark ? 'dark-mode' : 'light-mode')
  document.documentElement.classList.add(isDark ? 'dark-mode' : 'light-mode')
}

/**
 * 监听主题变化
 *
 * 当 Vuetify 主题发生变化时（如用户切换明/暗模式），
 * 自动更新 body 的样式类以应用对应的全局背景
 *
 * immediate: true - 立即执行一次，确保初始化时也应用了正确的样式
 *
 * 注意：首次回调是初始化同步，此时 onMounted 尚未恢复用户保存的
 * 主题偏好，不能把默认主题写入 localStorage（会覆盖手动设置），
 * 因此仅在非首次变化时才持久化
 */
let themeWatchInitialized = false
watch(
  () => theme.global.current.value.dark,
  () => {
    updateBodyStyles()
    if (themeWatchInitialized) {
      localStorage.setItem('theme', theme.global.current.value.dark ? 'dark' : 'light')
    }
    themeWatchInitialized = true
  },
  { immediate: true }
)

/**
 * 处理系统主题变化的回调函数
 *
 * 当用户的操作系统主题发生变化时（如在系统设置中切换主题），
 * 此函数会被调用，自动将应用主题同步为系统主题
 *
 * @param {MediaQueryListEvent} e - 媒体查询事件对象
 */
const handleSystemThemeChange = (e) => {
  // 根据系统偏好设置新的主题
  const newTheme = e.matches ? 'dark' : 'light'

  // 更新 Vuetify 主题
  theme.global.name.value = newTheme

  // 清除手动偏好设置
  // 这样下次访问时会继续跟随系统主题，而不是使用手动设置的
  localStorage.removeItem('theme_manual')

  // 同步保存到本地存储
  localStorage.setItem('theme', newTheme)
}

/**
 * 滚动条自动隐藏：捕获阶段监听全局滚动（scroll 不冒泡，但捕获可截获
 * 任意容器的滚动），滚动时给 <html> 加 .scrolling 类，静置后移除。
 * 配合全局 CSS：默认 thumb 透明，.scrolling 时渐显。
 */
let scrollHideTimer = null
const handleAnyScroll = () => {
  document.documentElement.classList.add('scrolling')
  clearTimeout(scrollHideTimer)
  scrollHideTimer = setTimeout(() => {
    document.documentElement.classList.remove('scrolling')
  }, 800)
}

/**
 * 组件挂载时的初始化逻辑
 */
onMounted(() => {
  // 获取系统主题偏好查询对象
  // matches 为 true 表示系统使用暗色主题
  const systemDark = window.matchMedia('(prefers-color-scheme: dark)')

  // 1. 注册系统主题变化监听器
  // 当用户操作系统主题时，自动同步应用到网站
  systemDark.addEventListener('change', handleSystemThemeChange)

  // 2. 注册全局滚动监听（捕获阶段，覆盖窗口与所有内部容器）
  window.addEventListener('scroll', handleAnyScroll, { capture: true, passive: true })

  // 3. 初始化主题设置
  // 检查用户是否有手动设置过主题偏好
  const savedTheme = localStorage.getItem('theme')
  const isManual = localStorage.getItem('theme_manual')

  if (isManual && savedTheme) {
    // 如果用户手动设置过主题，使用用户设置
    theme.global.name.value = savedTheme
  } else {
    // 否则跟随系统主题设置
    theme.global.name.value = systemDark.matches ? 'dark' : 'light'
  }

  // 3. 确保 body 元素有初始的样式类
  // 这会应用正确的背景图片和颜色
  const isDark = theme.global.current.value.dark
  document.body.classList.add(isDark ? 'dark-mode' : 'light-mode')
})

/**
 * 组件卸载前清理
 *
 * 移除注册的事件监听器，防止内存泄漏
 * 良好的实践：onMounted 中注册的事件应在 onUnmounted 中移除
 */
onUnmounted(() => {
  const systemDark = window.matchMedia('(prefers-color-scheme: dark)')
  systemDark.removeEventListener('change', handleSystemThemeChange)
  window.removeEventListener('scroll', handleAnyScroll, { capture: true })
  clearTimeout(scrollHideTimer)
})
</script>

<style>
/* ============================================
   全局样式部分

   注意：这些是全局CSS样式，会影响整个应用
   scoped 样式只影响当前组件，但这里的样式会影响全局
   ============================================ */

/* 重置 HTML 和 body 的默认样式 */
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}

/* 重复定义（保留原文件内容） */
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}

/**
 * body 元素背景样式
 *
 * 关键配置：
 * - background-attachment: fixed - 背景图片固定，不随滚动滚动
 * - background-size: cover - 背景图片覆盖整个页面
 * - transition: 平滑过渡主题切换时的变化
 */
body {
  background-position: center center !important;
  background-repeat: no-repeat !important;
  background-attachment: fixed !important;
  background-size: cover !important;
  /* 主题切换时的平滑过渡效果 */
  transition: background-image 0.3s ease-in-out, background-color 0.3s ease-in-out;
}

/* ============================================
   滚动条样式：两端圆角，随明暗主题联动
   - WebKit (Chrome/Edge/Safari): ::-webkit-scrollbar 系列
   - Firefox: scrollbar-width / scrollbar-color
   - 窗口滚动条在 html 上，Vuetify 主题变量不在其作用域内，
     由 html.light-mode / html.dark-mode 镜像类提供颜色
   - 自动隐藏：默认 thumb 透明，滚动时 App.vue 给 <html> 加
     .scrolling 类，thumb 渐显，静置 800ms 后渐隐
   ============================================ */

/* 窗口滚动条主题色（内部容器直接用 Vuetify 主题变量） */
html.light-mode {
  --sb-thumb: rgba(0, 0, 0, 0.35);
  --sb-thumb-hover: rgba(0, 0, 0, 0.5);
}

html.dark-mode {
  --sb-thumb: rgba(255, 255, 255, 0.35);
  --sb-thumb-hover: rgba(255, 255, 255, 0.55);
}

::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  /* 两端圆角；透明边框 + content-box 让圆角更自然 */
  border-radius: 999px;
  border: 2px solid transparent;
  background-clip: content-box;
  /* 默认隐藏，滚动时由 html.scrolling 规则点亮 */
  background-color: transparent;
  transition: background-color 0.25s ease;
}

/* 滚动中：内部容器 thumb 显示（Vuetify 主题变量在此作用域内可用） */
html.scrolling ::-webkit-scrollbar-thumb {
  background-color: rgb(var(--v-theme-surface-variant, 158, 158, 158));
}

html.scrolling ::-webkit-scrollbar-thumb:hover {
  background-color: rgb(var(--v-theme-primary, 103, 58, 183));
}

/* 滚动中：窗口滚动条（html 自身，主题变量不可达，用镜像变量） */
html.scrolling::-webkit-scrollbar-thumb {
  background-color: var(--sb-thumb);
}

html.scrolling::-webkit-scrollbar-thumb:hover {
  background-color: var(--sb-thumb-hover);
}

/* Firefox：默认透明，滚动中显示 */
* {
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;
}

html.scrolling,
html.scrolling * {
  scrollbar-color: rgb(var(--v-theme-surface-variant, 158, 158, 158)) transparent;
}

html.scrolling {
  scrollbar-color: var(--sb-thumb) transparent;
}

/* 浅色模式下的背景设置 */
body.light-mode {
  /* 使用浅色模式的slogan图片作为背景 */
  background-image: url('/images/slogan_light.png') !important;
  background-color: #f0f2f5 !important;
}

/* 深色模式下的背景设置 */
body.dark-mode {
  /* 使用深色模式的slogan图片作为背景 */
  background-image: url('/images/slogan_dark.png') !important;
  background-color: #000000 !important;
}

/**
 * Vuetify 组件透明化
 *
 * 将 Vuetify 的主要容器组件设置为透明，
 * 以便显示 body 的背景图片
 */
.v-application,
.v-application__wrap,
.v-layout,
.v-main {
  background: transparent !important;
}

/**
 * 毛玻璃效果（Glassmorphism）
 *
 * 为卡片、纸张、导航抽屉和对话框添加磨砂玻璃效果
 * 使用 backdrop-filter 实现模糊和饱和度调整
 */
.v-card,
.v-sheet.rounded-lg.elevation-1,
.v-navigation-drawer,
.v-dialog > .v-card {
  backdrop-filter: blur(20px) saturate(160%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(160%) !important;
}

/* 浅色主题下的毛玻璃效果 */
body.light-mode .v-card,
body.light-mode .v-sheet.rounded-lg.elevation-1,
body.light-mode .v-navigation-drawer {
  /* 半透明白色背景 + 细微白色边框 */
  background-color: rgba(255, 255, 255, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.4) !important;
  /* 微妙的阴影效果 */
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.07) !important;
}

/* 深色主题下的毛玻璃效果 */
body.dark-mode .v-card,
body.dark-mode .v-sheet.rounded-lg.elevation-1,
body.dark-mode .v-navigation-drawer {
  /* 半透明深色背景 + 细微白色边框 */
  background-color: rgba(30, 30, 30, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  /* 深色阴影效果 */
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.3) !important;
}

/**
 * 顶部导航栏 (AppBar) 毛玻璃效果
 *
 * 相比其他组件使用较轻的模糊效果
 */
.v-app-bar {
  backdrop-filter: blur(12px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(12px) saturate(180%) !important;
}

/* 浅色模式下 AppBar 的样式 */
body.light-mode .v-app-bar {
  background-color: rgba(255, 255, 255, 0.4) !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.3) !important;
}

/* 深色模式下 AppBar 的样式 */
body.dark-mode .v-app-bar {
  background-color: rgba(15, 15, 15, 0.4) !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
}

/**
 * 页脚 (Footer) 毛玻璃效果
 *
 * 使用最轻的模糊效果和透明度
 */
.v-footer {
  background-color: rgba(255, 255, 255, 0.1) !important;
  backdrop-filter: blur(10px) saturate(150%) !important;
  -webkit-backdrop-filter: blur(10px) saturate(150%) !important;
  border-top: 1px solid rgba(255, 255, 255, 0.1) !important;
}

/* 浅色模式下 Footer 的样式 */
body.light-mode .v-footer {
  background-color: rgba(255, 255, 255, 0.3) !important;
  color: #333 !important;
}

/* 深色模式下 Footer 的样式 */
body.dark-mode .v-footer {
  background-color: rgba(0, 0, 0, 0.3) !important;
  color: #ccc !important;
}

/**
 * 毛玻璃容器内的文字可读性增强
 *
 * 通过添加细微的文字阴影来提高对比度
 */
.v-card-title, .v-card-text, .v-list-item-title {
  text-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

/* 深色模式下的文字阴影（更明显以确保可读性） */
body.dark-mode .v-card-title,
body.dark-mode .v-card-text {
  text-shadow: 0 1px 2px rgba(0,0,0,0.5);
}
</style>
