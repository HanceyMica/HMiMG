<!--
  登录/注册页面组件 (Login.vue)

  功能说明：
  - 提供用户登录功能（用户名 + 密码）
  - 提供用户注册功能（用户名 + 密码 + 邮箱 + 电话）
  - 支持明暗主题切换
  - 显示网站标题和版权信息
  - 支持国际化和多语言

  页面布局：
  - 全屏背景（带动态主题效果）
  - 居中显示玻璃态登录卡片
  - 卡片顶部：主题切换按钮、网站图标、标题、副标题
  - 卡片中部：登录/注册表单
  - 卡片底部：注册/登录切换链接、版权信息

  安全说明：
  - 登录/注册失败时显示错误提示，不会暴露敏感信息
  - 未授权用户访问受保护页面时会重定向到登录页并显示提示
-->
<template>
  <!-- 全屏容器：背景图片根据主题变化，fill-height 使其占满整个视口 -->
  <v-container class="fill-height login-bg" fluid>
    <!-- 垂直居中的行 -->
    <v-row align="center" justify="center">
      <!-- 登录卡片：响应式宽度（小屏占12列，大屏占3列） -->
      <v-col cols="12" sm="8" md="4" lg="3">
        <v-card class="elevation-4 rounded-xl pa-4 glass-card position-relative">
          <!-- 主题切换按钮：固定在卡片右上角 -->
          <div class="position-absolute top-0 right-0 pa-4">
            <v-btn
              icon
              variant="text"
              :color="isDark ? 'yellow' : 'primary'"
              @click="toggleTheme"
            >
              <v-icon>{{ isDark ? 'mdi-weather-night' : 'mdi-weather-sunny' }}</v-icon>
            </v-btn>
          </div>

          <!-- 卡片头部区域：网站图标、标题、副标题 -->
          <v-card-item class="text-center py-8">
            <!-- 网站图标 -->
            <v-icon size="64" color="primary" class="mb-4">mdi-image-multiple</v-icon>
            <!-- 网站标题：来自设置 Store -->
            <v-card-title class="text-h4 font-weight-bold">
              {{ settingsStore.websiteTitle }}
            </v-card-title>
            <!-- 副标题：根据是登录还是注册状态显示不同文本 -->
            <v-card-subtitle class="text-subtitle-1 mt-2">
              {{ isRegister ? $t('login.registerBtn') : $t('login.title') }}
            </v-card-subtitle>
          </v-card-item>

          <!-- 卡片内容区域：登录/注册表单 -->
          <v-card-text>
            <!-- 表单：阻止默认提交行为，调用 handleSubmit 处理 -->
            <v-form @submit.prevent="handleSubmit" class="mt-4">
              <!-- 用户名字段：登录和注册都需要 -->
              <v-text-field
                v-model="form.username"
                :label="$t('login.username')"
                prepend-inner-icon="mdi-account-outline"
                variant="outlined"
                density="comfortable"
                required
                class="mb-2"
              ></v-text-field>

              <!-- 密码字段：登录和注册都需要 -->
              <v-text-field
                v-model="form.password"
                :label="$t('login.password')"
                prepend-inner-icon="mdi-lock-outline"
                type="password"
                variant="outlined"
                density="comfortable"
                required
                class="mb-2"
              ></v-text-field>

              <!-- 注册专属字段：邮箱和电话，仅在注册时显示 -->
              <!-- v-expand-transition 实现展开/收起动画效果 -->
              <v-expand-transition>
                <div v-if="isRegister">
                  <!-- 邮箱字段 -->
                  <v-text-field
                    v-model="form.email"
                    :label="$t('login.email')"
                    prepend-inner-icon="mdi-email-outline"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-text-field>
                  <!-- 电话字段 -->
                  <v-text-field
                    v-model="form.phone"
                    :label="$t('login.phone')"
                    prepend-inner-icon="mdi-phone-outline"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-text-field>
                </div>
              </v-expand-transition>

              <!-- 错误提示框：当存在错误信息时显示 -->
              <!-- closable 允许用户关闭此提示框 -->
              <v-alert
                v-if="error"
                type="error"
                variant="tonal"
                density="compact"
                class="mb-4 rounded-lg"
                closable
                @click:close="error = ''"
              >
                {{ error }}
              </v-alert>

              <!-- 提交按钮：登录或注册，根据当前状态显示 -->
              <v-btn
                type="submit"
                color="primary"
                block
                size="large"
                class="rounded-lg text-none font-weight-bold mt-4"
                :loading="loading"
                elevation="2"
              >
                {{ isRegister ? $t('login.registerBtn') : $t('login.loginBtn') }}
              </v-btn>
            </v-form>
          </v-card-text>

          <!-- 卡片操作区域：登录/注册切换链接，仅在允许注册时显示 -->
          <v-card-actions class="justify-center pb-6" v-if="settingsStore.allowRegistration">
            <v-btn
              variant="text"
              color="primary"
              class="text-none"
              @click="isRegister = !isRegister"
            >
              <!-- 根据当前状态显示对应的切换文本 -->
              {{ isRegister ? $t('login.toLogin') : $t('login.toRegister') }}
            </v-btn>
          </v-card-actions>

          <!-- 分隔线 -->
          <v-divider class="mx-4"></v-divider>

          <!-- 版权信息：显示网站标题和开发者信息 -->
          <div class="text-center py-4 text-grey text-caption">
            {{ settingsStore.websiteTitle }} © 2026 HanceyMica
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- 全局消息提示条：用于显示操作结果通知（如登录失败、注册成功等） -->
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :timeout="3000"
      elevation="24"
      rounded="pill"
    >
      <div class="d-flex align-center">
        <v-icon start class="mr-2">{{ snackbar.icon }}</v-icon>
        {{ snackbar.text }}
      </div>
      <template v-slot:actions>
        <v-btn variant="text" @click="snackbar.show = false">Close</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<style scoped>
/*
 * 登录页面样式
 *
 * 主要效果：
 * 1. 背景根据主题变化（浅色/深色主题使用不同的背景图和底色）
 * 2. 登录卡片使用玻璃态效果（毛玻璃背景 + 半透明边框）
 */

/* 全屏背景容器样式 */
.login-bg {
  background-position: center center !important;  /* 背景图居中 */
  background-repeat: no-repeat !important;         /* 不重复平铺 */
  background-attachment: fixed !important;        /* 背景图固定，不随滚动滚动 */
  background-size: cover !important;               /* 背景图覆盖整个容器 */
  /* 背景切换时有过渡动画效果 */
  transition: background-image 0.3s ease-in-out, background-color 0.3s ease-in-out;
}

/* 浅色主题下的背景样式 */
.v-theme--light .login-bg {
  /* 使用浅色主题的背景图 */
  background-image: url('/images/slogan_light.png') !important;
  background-color: #f0f2f5 !important;  /* 备用背景色：浅灰蓝 */
}

/* 深色主题下的背景样式 */
.v-theme--dark .login-bg {
  /* 使用深色主题的背景图 */
  background-image: url('/images/slogan_dark.png') !important;
  background-color: #000000 !important;  /* 备用背景色：黑色 */
}

/*
 * 玻璃态卡片效果
 * 使用 backdrop-filter 实现模糊背景效果
 * 兼容 WebKit 内核浏览器（如 Safari）
 */
.glass-card {
  backdrop-filter: blur(20px) saturate(160%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(160%) !important;
}

/* 浅色主题下玻璃卡片样式：白色半透明 + 浅色边框 */
.v-theme--light .glass-card {
  background-color: rgba(255, 255, 255, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.4) !important;
}

/* 深色主题下玻璃卡片样式：深色半透明 + 浅色边框 */
.v-theme--dark .glass-card {
  background-color: rgba(30, 30, 30, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
}
</style>

<script setup>
/*
 * Login.vue - 登录/注册页面逻辑
 *
 * 主要功能：
 * 1. 用户登录：验证用户名密码，成功后跳转到首页
 * 2. 用户注册：提交注册信息（用户名、密码、邮箱、电话）
 * 3. 主题切换：支持浅色/深色主题切换，并持久化到本地存储
 * 4. 权限验证：处理未授权访问的重定向
 *
 * 状态管理：
 * - userStore: 存储用户登录信息
 * - settingsStore: 存储网站公开设置（如网站标题、是否允许注册）
 */

import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useSettingsStore } from '@/store/settings'
import { useI18n } from 'vue-i18n'
import { useTheme } from 'vuetify'
import api from '@/lib/api'

// =============================================
// 路由与 Store 相关
// =============================================

const router = useRouter()           // 路由实例，用于页面导航
const route = useRoute()             // 当前路由对象，用于获取查询参数
const userStore = useUserStore()      // 用户状态管理，包含登录用户信息和 Token
const settingsStore = useSettingsStore()  // 网站设置状态管理
const { t } = useI18n()              // 国际化函数，用于翻译文本
const theme = useTheme()             // Vuetify 主题实例，用于切换主题

// =============================================
// 响应式状态（Ref）
// =============================================

/**
 * isRegister: 控制页面显示登录还是注册表单
 * false = 显示登录表单，true = 显示注册表单
 */
const isRegister = ref(false)

/**
 * loading: 表单提交进行中的加载状态标志
 * true 时提交按钮会显示加载动画
 */
const loading = ref(false)

/**
 * error: 错误信息文本，当登录/注册失败时显示
 * 为空字符串时不显示错误提示框
 */
const error = ref('')

// =============================================
// 计算属性（Computed）
// =============================================

/**
 * isDark: 判断当前是否为深色主题
 * 用于控制主题切换按钮的图标和颜色
 */
const isDark = computed(() => theme.global.current.value.dark)

// =============================================
// 主题切换功能
// =============================================

/**
 * toggleTheme: 切换浅色/深色主题
 *
 * 实现流程：
 * 1. 判断当前主题状态，取反得到新主题
 * 2. 使用 theme.global.name.value = newTheme 切换 Vuetify 主题
 * 3. 将主题设置标记为"手动设置"并存入 localStorage
 * 4. 持久化当前主题到 localStorage，用于下次访问时恢复
 */
const toggleTheme = () => {
  const newTheme = isDark.value ? 'light' : 'dark'
  theme.global.name.value = newTheme
  localStorage.setItem('theme_manual', 'true')
  localStorage.setItem('theme', newTheme)
}

// =============================================
// 生命周期钩子
// =============================================

/**
 * onMounted: 组件挂载后执行的初始化操作
 *
 * 执行流程：
 * 1. 调用 settingsStore.fetchPublicSettings() 获取网站公开设置
 *    （包括网站标题、是否允许注册等）
 * 2. 检查 URL 查询参数是否有 reason=unauthorized
 *    - 如果有，说明用户尝试访问未授权页面，需要登录
 *    - 显示提示消息后，清除查询参数（避免刷新时重复提示）
 */
onMounted(async () => {
  // 获取网站公开设置
  await settingsStore.fetchPublicSettings()

  // 检查是否为未授权重定向
  if (route.query.reason === 'unauthorized') {
    showNotify(t('common.notLoggedIn'), 'error')
    // 清除查询参数，防止刷新时重复触发
    router.replace({ query: {} })
  }
})

// =============================================
// 通知消息状态（Snackbar）
// =============================================

/**
 * snackbar: 全局通知消息的状态对象
 * - show: 控制消息条的显示/隐藏
 * - text: 要显示的消息文本
 * - color: 消息条的颜色（success/error/info）
 * - icon: 消息条左侧的图标名称
 */
const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
  icon: 'mdi-check-circle'
})

/**
 * showNotify: 显示通知消息的辅助函数
 * @param {string} text - 要显示的消息文本
 * @param {string} color - 消息颜色类型，默认为 'success'
 *
 * 使用示例：showNotify('登录成功！') 或 showNotify('用户名或密码错误', 'error')
 */
const showNotify = (text, color = 'success') => {
  snackbar.text = text
  snackbar.color = color
  snackbar.icon = color === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'
  snackbar.show = true
}

// =============================================
// 表单数据
// =============================================

/**
 * form: 登录/注册表单的数据对象
 * - username: 用户名（登录和注册都需要）
 * - password: 密码（登录和注册都需要）
 * - email: 邮箱（仅注册时需要）
 * - phone: 电话号码（仅注册时需要）
 *
 * 使用 reactive 使其成为深度响应式对象
 */
const form = reactive({
  username: '',
  password: '',
  email: '',
  phone: ''
})

// =============================================
// 表单提交处理
// =============================================

/**
 * handleSubmit: 处理登录/注册表单提交
 *
 * 执行流程：
 * 1. 防止重复提交（检查 loading 状态）
 * 2. 设置加载状态，清空之前的错误信息
 * 3. 根据 isRegister 状态决定执行登录还是注册
 *
 * 登录流程：
 *   - 发送 POST /login 请求，包含用户名和密码
 *   - 成功后调用 userStore.setUser() 保存用户信息
 *   - 跳转到首页 /
 *
 * 注册流程：
 *   - 发送 POST /register 请求，包含所有表单字段
 *   - 成功后切换回登录状态，显示注册成功消息
 *
 * 失败处理：
 *   - 显示错误信息，优先使用后端返回的错误描述
 */
const handleSubmit = async () => {
  // 防止重复提交
  if (loading.value) return

  loading.value = true   // 开始加载状态
  error.value = ''       // 清空错误信息

  try {
    if (isRegister.value) {
      // ====================
      // 注册流程
      // ====================
      await api.post('/register', form)
      // 注册成功后切换回登录状态
      isRegister.value = false
      showNotify(t('login.registerSuccess'))
    } else {
      // ====================
      // 登录流程
      // ====================
      const res = await api.post('/login', {
        username: form.username,
        password: form.password
      })
      // 保存用户信息和 Token 到 Store
      userStore.setUser(res.data.user, res.data.token)
      // 跳转到首页
      router.push('/')
    }
  } catch (err) {
    // 显示错误信息
    // 优先使用后端返回的错误信息，否则使用国际化默认文本
    error.value = err.response?.data?.error || t('login.failed')
  } finally {
    // 无论成功或失败，都要关闭加载状态
    loading.value = false
  }
}
</script>
