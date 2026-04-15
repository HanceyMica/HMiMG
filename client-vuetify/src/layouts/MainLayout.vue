<!--
  MainLayout.vue - 应用主布局组件

  本组件定义了应用的主要页面布局结构，包括：
  - 顶部导航栏 (AppBar)：包含Logo、导航菜单、主题切换、用户登录/登出
  - 移动端侧边导航抽屉 (Navigation Drawer)：在小屏幕设备上替代顶部导航
  - 主内容区域 (v-main)：显示路由页面内容
  - 页脚 (Footer)：显示版权信息

  布局特点：
  - 响应式设计：根据屏幕尺寸自动切换桌面/移动端布局
  - 毛玻璃效果：导航栏和抽屉使用磨砂玻璃视觉效果
  - 多语言支持：所有文案通过 i18n 国际化系统获取
-->
<template>
  <!-- v-layout: Vuetify 提供的布局容器组件 -->
  <v-layout>
    <!-- ============================================
         移动端导航抽屉
         v-model="drawer": 控制抽屉的显示/隐藏状态
         temporary: 临时模式，点击外部或导航后自动关闭
         class="glass-drawer": 应用毛玻璃样式
         ============================================ -->
    <v-navigation-drawer v-model="drawer" temporary class="glass-drawer">
      <!-- 导航列表 - 第一部分：主要导航菜单 -->
      <v-list nav>
        <!-- 首页导航项 -->
        <v-list-item
          prepend-icon="mdi-home"
          to="/"
          :title="$t('common.home')"
          @click="drawer = false"
        ></v-list-item>

        <!-- 图库导航项 -->
        <v-list-item
          prepend-icon="mdi-image-multiple"
          to="/library"
          :title="$t('common.library')"
          @click="drawer = false"
        ></v-list-item>

        <!-- 后台管理导航项（仅管理员可见） -->
        <!-- v-if="isAdmin": 只有管理员角色才显示此选项 -->
        <v-list-item
          v-if="isAdmin"
          prepend-icon="mdi-shield-check"
          to="/admin"
          :title="$t('common.admin')"
          @click="drawer = false"
        ></v-list-item>

        <!-- 关于页面导航项 -->
        <v-list-item
          prepend-icon="mdi-information"
          to="/about"
          :title="$t('common.about')"
          @click="drawer = false"
        ></v-list-item>
      </v-list>

      <!-- 分割线 - 分隔主要导航和辅助功能 -->
      <v-divider></v-divider>

      <!-- 导航列表 - 第二部分：辅助功能 -->
      <v-list nav>
        <!-- 主题切换按钮 -->
        <!-- 根据当前主题显示太阳或月亮图标 -->
        <v-list-item
          :prepend-icon="isDark ? 'mdi-weather-night' : 'mdi-weather-sunny'"
          :title="isDark ? $t('common.dark') : $t('common.light')"
          @click="toggleTheme"
        ></v-list-item>

        <!-- 已登录用户：显示登出按钮 -->
        <v-list-item
          v-if="user"
          prepend-icon="mdi-logout"
          :title="$t('common.logout')"
          @click="handleLogout"
        ></v-list-item>

        <!-- 未登录用户：显示登录按钮 -->
        <v-list-item
          v-else
          prepend-icon="mdi-login"
          to="/login"
          :title="$t('login.loginBtn')"
          @click="drawer = false"
        ></v-list-item>
      </v-list>
    </v-navigation-drawer>

    <!-- ============================================
         顶部导航栏
         flat: 移除默认的 elevation 阴影
         border: 显示底部边框
         height="64": 固定高度为 64px
         ============================================ -->
    <v-app-bar flat border height="64">
      <v-container class="d-flex align-center pa-0 fill-height">
        <!-- 移动端菜单按钮 - 仅在小屏幕显示 -->
        <!-- d-md-none: 在中等及以上尺寸屏幕隐藏 -->
        <v-app-bar-nav-icon class="d-md-none mr-2" @click="drawer = !drawer"></v-app-bar-nav-icon>

        <!-- ============================================
             Logo 和网站标题区域
             ============================================ -->
        <div class="mr-4">
          <!-- router-link: Vue Router 提供的链接组件，点击跳转到首页 -->
          <router-link to="/" class="d-flex align-center text-decoration-none">
            <!-- 网站图标 -->
            <v-avatar size="32" class="mr-2 rounded-lg overflow-hidden border">
              <img
                src="/images/favicon.png"
                alt="HMiMG"
                style="width: 100%; height: 100%; object-fit: cover;"
              >
            </v-avatar>

            <!-- 网站标题 -->
            <!-- d-none d-sm-flex: 在超小屏幕隐藏，仅在 sm 及以上尺寸显示 -->
            <span class="text-h6 font-weight-bold d-none d-sm-flex site-title">
              <!-- 从 settingsStore 获取网站标题 -->
              {{ settingsStore.websiteTitle }}
            </span>
          </router-link>
        </div>

        <!-- ============================================
             桌面端导航菜单
             d-none d-md-flex: 仅在中等及以上尺寸屏幕显示
             ============================================ -->
        <div class="d-none d-md-flex align-center">
          <!-- 首页按钮 -->
          <v-btn
            to="/"
            variant="text"
            class="text-none mx-1"
            :active="route.path === '/'"
          >
            {{ $t('common.home') }}
          </v-btn>

          <!-- 图库按钮 -->
          <!-- active 条件包含多个路径：/library 及其子路径 -->
          <v-btn
            to="/library"
            variant="text"
            class="text-none mx-1"
            :active="route.path.startsWith('/library') || route.path.startsWith('/album') || route.path.startsWith('/collection') || route.path.startsWith('/image')"
          >
            {{ $t('common.library') }}
          </v-btn>

          <!-- 后台管理按钮（仅管理员） -->
          <v-btn
            v-if="isAdmin"
            to="/admin"
            variant="text"
            class="text-none mx-1"
            :active="route.path.startsWith('/admin')"
          >
            {{ $t('common.admin') }}
          </v-btn>

          <!-- 关于按钮 -->
          <v-btn
            to="/about"
            variant="text"
            class="text-none mx-1"
            :active="route.path.startsWith('/about')"
          >
            {{ $t('common.about') }}
          </v-btn>
        </div>

        <!-- 填充剩余空间，推动右侧内容到边缘 -->
        <v-spacer></v-spacer>

        <!-- ============================================
             桌面端右侧操作按钮区域
             d-none d-md-flex: 仅在中等及以上尺寸屏幕显示
             ============================================ -->
        <div class="d-none d-md-flex align-center">
          <!-- 主题切换按钮 -->
          <v-btn
            icon
            variant="text"
            class="mr-2"
            :color="isDark ? 'yellow' : 'primary'"
            @click="toggleTheme"
          >
            <!-- 根据主题显示不同的图标 -->
            <v-icon>{{ isDark ? 'mdi-weather-night' : 'mdi-weather-sunny' }}</v-icon>
          </v-btn>

          <!-- 已登录用户：显示登出按钮 -->
          <v-btn
            v-if="user"
            variant="flat"
            size="small"
            class="rounded-lg px-4 logout-btn font-weight-bold"
            prepend-icon="mdi-logout-variant"
            @click="handleLogout"
          >
            {{ $t('common.logout') }}
          </v-btn>

          <!-- 未登录用户：显示登录按钮 -->
          <v-btn
            v-else
            to="/login"
            color="primary"
            variant="flat"
            size="small"
            class="rounded-pill px-4"
          >
            {{ $t('login.loginBtn') }}
          </v-btn>
        </div>

        <!-- ============================================
             移动端右侧操作区域
             d-md-none: 仅在小屏幕显示
             ============================================ -->
        <div class="d-md-none">
          <!-- 未登录用户显示登录图标按钮 -->
          <v-btn v-if="!user" to="/login" icon="mdi-login" variant="text" size="small"></v-btn>
        </div>
      </v-container>
    </v-app-bar>

    <!-- ============================================
         主内容区域
         v-main: Vuetify 提供的主内容容器组件
         ============================================ -->
    <v-main>
      <v-container class="py-4 py-sm-8">
        <!-- 内容纸张组件 - 包含页面内容 -->
        <v-sheet class="pa-4 pa-sm-6 rounded-lg elevation-1 min-h-main">
          <!-- router-view: 显示当前路由对应的页面组件 -->
          <router-view />
        </v-sheet>
      </v-container>
    </v-main>

    <!-- ============================================
         页脚
         app: 将 footer 固定在应用底部
         border: 显示顶部边框
         height="60": 固定高度为 60px
         ============================================ -->
    <v-footer app border height="60">
      <div class="text-center w-100 text-grey text-caption text-sm-body-2">
        <!-- 版权信息，通过 i18n 格式化 -->
        <!-- {title} 会被替换为实际的网站标题 -->
        {{ $t('common.copyright', { title: settingsStore.websiteTitle }) }}
      </div>
    </v-footer>
  </v-layout>
</template>

<script setup>
/**
 * Script Setup 部分
 * 使用 Vue 3 Composition API 定义组件的响应式状态和逻辑
 */

// ============================================
// 响应式API和生命周期钩子导入
// ============================================
import { ref, computed, onMounted, watch } from 'vue'

// ============================================
// Pinia Store 导入
// ============================================
// userStore: 管理用户登录状态和用户信息
import { useUserStore } from '@/store/user'
// settingsStore: 管理网站设置（如网站标题、默认语言等）
import { useSettingsStore } from '@/store/settings'

// ============================================
// Vue Router 导入
// ============================================
// useRouter: 获取路由器实例，用于编程式导航
import { useRouter, useRoute } from 'vue-router'

// ============================================
// Vue I18n 导入
// ============================================
// useI18n: 获取国际化函数
import { useI18n } from 'vue-i18n'

// ============================================
// Vuetify 导入
// ============================================
// useTheme: 获取主题管理函数
import { useTheme } from 'vuetify'

// ============================================
// API 库导入
// ============================================
// api: 封装了 HTTP 请求的工具库，用于与后端通信
import api from '@/lib/api'

// ============================================
// Store 实例化
// ============================================
const userStore = useUserStore()
const settingsStore = useSettingsStore()

// ============================================
// Router 实例化
// ============================================
const router = useRouter()
const route = useRoute()

// ============================================
// I18n 和 Theme 实例化
// ============================================
const { locale, t } = useI18n()
const theme = useTheme()

// ============================================
// 响应式状态定义
// ============================================
// drawer: 控制移动端导航抽屉的显示/隐藏
const drawer = ref(false)

// ============================================
// 计算属性 (Computed Properties)
// ============================================

/**
 * user 计算属性
 * 从 userStore 获取当前登录用户信息
 * 返回 user 对象或 undefined（未登录时）
 */
const user = computed(() => userStore.user)

/**
 * isAdmin 计算属性
 * 判断当前用户是否为管理员
 * 通过检查用户的 role 字段是否为 'admin'
 */
const isAdmin = computed(() => userStore.user?.role === 'admin')

/**
 * isDark 计算属性
 * 判断当前是否为暗色主题
 * 通过 Vuetify 主题系统获取
 */
const isDark = computed(() => theme.global.current.value.dark)

/**
 * pageTitle 计算属性
 * 动态生成当前页面的浏览器标签页标题
 * 格式为：页面名称 - 网站标题
 */
const pageTitle = computed(() => {
  // 获取网站标题
  const websiteTitle = settingsStore.websiteTitle
  // 获取当前路由路径
  const path = route.path

  // 根路径显示网站标题
  if (path === '/') return websiteTitle

  // 根据路径判断页面类型
  let pageName = ''
  if (path === '/library') pageName = t('common.library')
  else if (path.startsWith('/album/')) pageName = t('common.album')
  else if (path.startsWith('/collection/')) pageName = t('common.collection')
  else if (path === '/admin') pageName = t('common.admin')
  else if (path === '/about') pageName = t('common.about')
  else if (path === '/login') pageName = t('login.title')
  else if (path.startsWith('/image/')) pageName = t('image.details')
  else pageName = t('notFound.title')

  // 返回格式化的标题
  return `${pageName} - ${websiteTitle}`
})

// ============================================
// 方法定义 (Methods)
// ============================================

/**
 * fetchData 异步函数
 * 组件挂载时获取必要的初始化数据
 */
const fetchData = async () => {
  try {
    // 获取公开的网站设置（网站标题、默认语言等）
    await settingsStore.fetchPublicSettings()

    // 如果后端有设置默认语言，且本地存储中没有用户偏好
    // 则使用后端设置的默认语言
    if (settingsStore.defaultLanguage && !localStorage.getItem('lang')) {
      locale.value = settingsStore.defaultLanguage
    }
  } catch (e) {
    // 静默处理错误，不影响用户体验
  }
}

/**
 * handleLogout 函数
 * 处理用户点击登出按钮的操作
 */
const handleLogout = () => {
  // 调用 store 的 logout 方法清除用户状态
  userStore.logout()
  // 关闭移动端导航抽屉
  drawer.value = false
  // 跳转到登录页面
  router.push('/login')
}

/**
 * toggleTheme 函数
 * 切换应用的主题（明/暗）
 */
const toggleTheme = () => {
  // 计算新主题名称
  const newTheme = isDark.value ? 'light' : 'dark'
  // 更新 Vuetify 全局主题
  theme.global.name.value = newTheme
  // 标记用户手动设置了主题（而非跟随系统）
  localStorage.setItem('theme_manual', 'true')
  // 保存主题设置到本地存储
  localStorage.setItem('theme', newTheme)
}

// ============================================
// 监听器 (Watchers)
// ============================================

/**
 * 监听页面标题变化
 * 当路由变化导致页面标题改变时，
 * 自动更新浏览器标签页的标题
 */
watch(pageTitle, (newTitle) => {
  document.title = newTitle
}, { immediate: true })

// ============================================
// 生命周期钩子 (Lifecycle Hooks)
// ============================================

/**
 * onMounted
 * 组件挂载完成后执行的钩子
 * 这里用于获取初始化数据
 */
onMounted(fetchData)
</script>

<style scoped>
/**
 * Scoped 样式
 * 这些样式仅影响当前组件，不会泄漏到其他组件
 */

/* 主内容区域的最小高度 */
.min-h-main {
  /* 确保页面在内容较少时也能正常显示 */
  min-height: 400px;
}

/**
 * 毛玻璃效果导航抽屉
 *
 * 使用 backdrop-filter 实现磨砂玻璃视觉效果：
 * - blur(20px): 20像素的高斯模糊
 * - saturate(160%): 160% 的色彩饱和度
 */
.glass-drawer {
  backdrop-filter: blur(20px) saturate(160%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(160%) !important;
  /* 使用 Vuetify 主题的 surface 颜色并设置透明度 */
  background-color: rgba(var(--v-theme-surface), 0.7) !important;
}

/**
 * 网站标题过渡动画
 * 平滑的颜色过渡效果
 */
.site-title {
  transition: color 0.3s ease;
}

/* 浅色主题下网站标题的颜色 */
.v-theme--light .site-title {
  color: #000000 !important;
}

/* 深色主题下网站标题的颜色 */
.v-theme--dark .site-title {
  color: #ffffff !important;
}

/**
 * 登出按钮样式
 *
 * 特点：
 * - 位置相对定位，overflow hidden 用于noise效果
 * - 使用 cubic-bezier 实现平滑的过渡动画
 * - 移除了默认的 text-transform 和 letter-spacing
 */
.logout-btn {
  position: relative;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  text-transform: none;
  letter-spacing: 0.5px;
}

/* noise 效果叠加层 - 使用 SVG 噪点纹理 */
.logout-btn::before {
  content: "";
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  /* 内联 SVG 作为背景图片 */
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
  opacity: 0.15;
  pointer-events: none;
  /* 使用 overlay 混合模式 */
  mix-blend-mode: overlay;
}

/* ============================================
   浅色模式下的登出按钮样式
   黑色背景 + 白色文字 + 阴影效果
   ============================================ */
.v-theme--light .logout-btn {
  background-color: #1a1a1a !important;
  color: #ffffff !important;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3), inset 0 1px 1px rgba(255, 255, 255, 0.1) !important;
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
}

/* 浅色模式下悬停效果 */
.v-theme--light .logout-btn:hover {
  background-color: #000000 !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.4) !important;
}

/* ============================================
   深色模式下的登出按钮样式
   白色背景 + 黑色文字 + 发光效果
   ============================================ */
.v-theme--dark .logout-btn {
  background-color: #ffffff !important;
  color: #000000 !important;
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.3), 0 0 2px rgba(255, 255, 255, 0.5) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
}

/* 深色模式下悬停效果 */
.v-theme--dark .logout-btn:hover {
  background-color: #f0f0f0 !important;
  transform: translateY(-2px);
  box-shadow: 0 0 20px rgba(255, 255, 255, 0.6), 0 0 10px rgba(255, 255, 255, 0.4) !important;
}
</style>
