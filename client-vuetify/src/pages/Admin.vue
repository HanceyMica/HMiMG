<template>
  <div class="admin-page">
    <div class="d-flex align-center mb-8">
      <v-icon size="40" color="primary" class="mr-4">mdi-shield-check-outline</v-icon>
      <h1 class="text-h4 font-weight-bold">{{ $t('admin.dashboard') }}</h1>
    </div>
    
    <v-tabs v-model="activeTab" color="primary" class="mb-6 border-b">
      <v-tab value="system" prepend-icon="mdi-tune">{{ $t('admin.systemSettings') }}</v-tab>
      <v-tab value="organize" prepend-icon="mdi-folder-cog">{{ $t('admin.organize') }}</v-tab>
      <v-tab value="account" prepend-icon="mdi-account-cog">{{ $t('admin.accountSettings') }}</v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <!-- System Settings -->
      <v-window-item value="system">
        <v-card border flat class="pa-6 rounded-lg">
          <v-form @submit.prevent="handleUpdateSettings">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="settings.website_title"
                  :label="$t('admin.websiteTitle')"
                  variant="outlined"
                  density="comfortable"
                  prepend-inner-icon="mdi-format-title"
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="settings.max_users"
                  :label="$t('admin.maxUsers')"
                  type="number"
                  variant="outlined"
                  density="comfortable"
                  prepend-inner-icon="mdi-account-group"
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="settings.default_language"
                  :items="languages"
                  item-title="label"
                  item-value="value"
                  :label="$t('admin.defaultLanguage')"
                  variant="outlined"
                  density="comfortable"
                  prepend-inner-icon="mdi-translate"
                ></v-select>
              </v-col>
              <v-col cols="12">
                <v-switch
                  v-model="settings.allow_registration"
                  :label="$t('admin.allowRegistration')"
                  color="primary"
                  inset
                  hide-details
                ></v-switch>
              </v-col>
            </v-row>
            <v-btn type="submit" color="primary" class="mt-6 px-8" size="large" :loading="savingSettings">
              {{ $t('admin.saveSettings') }}
            </v-btn>
          </v-form>
        </v-card>
      </v-window-item>

      <!-- Organize (Albums/Collections) -->
      <v-window-item value="organize">
        <v-row>
          <v-col cols="12" md="6">
            <v-card border flat class="pa-6 rounded-lg h-100">
              <div class="text-h6 font-weight-bold mb-4 d-flex align-center">
                <v-icon start color="primary">mdi-image-plus</v-icon>
                {{ $t('admin.createAlbum') }}
              </div>
              <v-form @submit.prevent="handleCreateAlbum">
                <v-text-field v-model="newAlbum.name" :label="$t('admin.name')" variant="outlined" density="comfortable"></v-text-field>
                <v-textarea v-model="newAlbum.description" :label="$t('admin.description')" variant="outlined" rows="3"></v-textarea>
                <v-btn type="submit" color="primary" block size="large">{{ $t('admin.create') }}</v-btn>
              </v-form>
            </v-card>
          </v-col>
          <v-col cols="12" md="6">
            <v-card border flat class="pa-6 rounded-lg h-100">
              <div class="text-h6 font-weight-bold mb-4 d-flex align-center">
                <v-icon start color="secondary">mdi-folder-plus</v-icon>
                {{ $t('admin.createCollection') }}
              </div>
              <v-form @submit.prevent="handleCreateCollection">
                <v-text-field v-model="newCollection.name" :label="$t('admin.name')" variant="outlined" density="comfortable"></v-text-field>
                <v-textarea v-model="newCollection.description" :label="$t('admin.description')" variant="outlined" rows="3"></v-textarea>
                <v-btn type="submit" color="secondary" block size="large">{{ $t('admin.create') }}</v-btn>
              </v-form>
            </v-card>
          </v-col>
          <v-col cols="12">
            <v-card border flat class="pa-6 rounded-lg">
              <div class="text-h6 font-weight-bold mb-4 d-flex align-center">
                <v-icon start color="success">mdi-file-tree</v-icon>
                {{ $t('admin.organize') }}
              </div>
              <v-form @submit.prevent="handleAddToCollection">
                <v-row>
                  <v-col cols="12" md="4">
                    <v-select
                      v-model="organizeForm.collectionId"
                      :items="collections"
                      item-title="name"
                      item-value="id"
                      :label="$t('admin.targetCollection')"
                      variant="outlined"
                      density="comfortable"
                    ></v-select>
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-radio-group v-model="organizeForm.itemType" inline label="Type">
                      <v-radio :label="$t('admin.album')" value="album"></v-radio>
                      <v-radio :label="$t('admin.collection')" value="collection"></v-radio>
                    </v-radio-group>
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-select
                      v-model="organizeForm.itemName"
                      :items="organizeForm.itemType === 'album' ? albums : collections"
                      item-title="name"
                      item-value="name"
                      :label="$t('admin.itemName')"
                      variant="outlined"
                      density="comfortable"
                    ></v-select>
                  </v-col>
                </v-row>
                <v-btn type="submit" color="success" class="px-8" size="large">{{ $t('admin.add') }}</v-btn>
              </v-form>
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>

      <!-- Account Settings -->
      <v-window-item value="account">
        <v-card border flat class="pa-6 rounded-lg max-w-800 mx-auto">
          <v-form @submit.prevent="handleUpdateProfile">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field v-model="profile.username" :label="$t('login.username')" variant="outlined" density="comfortable" prepend-inner-icon="mdi-account"></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="profile.email" :label="$t('login.email')" variant="outlined" density="comfortable" prepend-inner-icon="mdi-email"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field v-model="profile.phone" :label="$t('login.phone')" variant="outlined" density="comfortable" prepend-inner-icon="mdi-phone"></v-text-field>
              </v-col>
            </v-row>
            
            <v-divider class="my-6"></v-divider>
            <div class="text-subtitle-1 font-weight-bold mb-4">{{ $t('admin.passwordHelp') }}</div>
            
            <v-row>
              <v-col cols="12">
                <v-text-field v-model="profile.oldPassword" :label="$t('admin.oldPassword')" type="password" variant="outlined" density="comfortable" prepend-inner-icon="mdi-lock-reset"></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="profile.password" :label="$t('admin.newPassword')" type="password" variant="outlined" density="comfortable" prepend-inner-icon="mdi-lock-plus"></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="profile.confirmPassword" :label="$t('admin.confirmPassword')" type="password" variant="outlined" density="comfortable" prepend-inner-icon="mdi-lock-check"></v-text-field>
              </v-col>
            </v-row>
            <v-btn type="submit" color="primary" class="mt-6 px-8" size="large">{{ $t('admin.updateProfile') }}</v-btn>
          </v-form>
        </v-card>
      </v-window-item>
    </v-window>

    <!-- Global Snackbar for Notifications -->
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
  </div>
</template>

<script setup>
/**
 * Admin.vue - 管理员控制面板页面
 *
 * 功能说明：
 * 1. 系统设置管理 - 网站标题、最大用户数、默认语言、是否允许注册
 * 2. 组织管理 - 创建相册、创建集合、将相册/集合添加到集合
 * 3. 账户设置 - 修改个人资料、修改密码
 *
 * 使用说明：
 * - 需要管理员权限才能访问此页面
 * - 所有表单提交均通过 API 与后端交互
 * - 操作成功/失败通过 Snackbar 通知组件反馈给用户
 */

// 引入 Vue 响应式 API
import { ref, reactive, onMounted } from 'vue'

// 引入 API 封装模块，用于与后端进行 HTTP 通信
import api from '@/lib/api'

// 引入用户状态管理，用于获取当前用户信息
import { useUserStore } from '@/store/user'

// 引入设置状态管理，用于同步公共设置
import { useSettingsStore } from '@/store/settings'

// 引入 Vue Router，用于页面跳转和路由控制
import { useRouter } from 'vue-router'

// 引入国际化功能，用于多语言支持
import { useI18n } from 'vue-i18n'

// ==================== 状态管理 ====================

// 用户 store 实例，用于访问当前登录用户信息
const userStore = useUserStore()

// 设置 store 实例，用于刷新和同步系统设置
const settingsStore = useSettingsStore()

// 路由实例，用于页面跳转
const router = useRouter()

// 国际化实例，用于获取翻译文本
const { t } = useI18n()

// ==================== 页面状态 ====================

// 当前激活的标签页，可选值：'system'(系统设置)、'organize'(组织管理)、'account'(账户设置)
const activeTab = ref('system')

// 是否正在保存设置的加载状态，用于防止重复提交
const savingSettings = ref(false)

// 相册列表，用于下拉选择和数据展示
const albums = ref([])

// 集合列表，用于下拉选择和数据展示
const collections = ref([])

// 支持的语言选项列表
// 用于系统设置中选择默认界面语言
const languages = [
  { label: 'English', value: 'en' },
  { label: '简体中文', value: 'zh' },
  { label: '日本語', value: 'ja' }
]

// ==================== 通知组件状态 ====================

/**
 * Snackbar 通知组件的状态管理
 * 用于显示操作成功或失败的提示信息
 * @property {boolean} show - 控制通知的显示/隐藏
 * @property {string} text - 通知的文本内容
 * @property {string} color - 通知的颜色主题 ('success' | 'error' | 'warning' 等)
 * @property {string} icon - 通知左侧的图标名称
 */
const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
  icon: 'mdi-check-circle'
})

/**
 * 显示通知提示的通用方法
 * @param {string} text - 要显示的通知文本
 * @param {string} color - 通知颜色，默认为 'success'
 */
const showNotify = (text, color = 'success') => {
  snackbar.text = text
  snackbar.color = color
  snackbar.icon = color === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'
  snackbar.show = true
}

// ==================== 表单数据 ====================

/**
 * 系统设置表单数据
 * @property {string} website_title - 网站标题
 * @property {number} max_users - 最大用户数量限制
 * @property {boolean} allow_registration - 是否允许新用户注册
 * @property {string} default_language - 默认界面语言
 */
const settings = reactive({
  website_title: '',
  max_users: 100,
  allow_registration: true,
  default_language: 'zh'
})

// 新建相册表单数据
// 用于在"组织管理"标签页中创建新的相册
const newAlbum = reactive({ name: '', description: '' })

// 新建集合表单数据
// 用于在"组织管理"标签页中创建新的集合
const newCollection = reactive({ name: '', description: '' })

/**
 * 组织管理表单数据
 * @property {number|null} collectionId - 目标集合的 ID
 * @property {string} itemType - 要添加的项目类型 ('album' | 'collection')
 * @property {string} itemName - 要添加的相册或集合的名称
 */
const organizeForm = reactive({ collectionId: null, itemType: 'album', itemName: '' })

/**
 * 个人资料表单数据
 * @property {string} username - 用户名
 * @property {string} email - 邮箱地址
 * @property {string} phone - 手机号码
 * @property {string} oldPassword - <PASSWORD>
 * @property {string} password - <PASSWORD>
 * @property {string} confirmPassword - 确认新密码
 */
const profile = reactive({
  username: userStore.user?.username || '',
  email: userStore.user?.email || '',
  phone: userStore.user?.phone || '',
  oldPassword: '',
  password: '',
  confirmPassword: ''
})

// ==================== 数据获取 ====================

/**
 * 页面加载时获取所有必要数据
 * 同时请求系统设置、相册列表和集合列表
 * 使用 Promise.all 并行请求以提高性能
 */
const fetchData = async () => {
  try {
    // 并行获取三项数据
    const [settRes, albRes, colRes] = await Promise.all([
      api.get('/settings'),   // 获取系统设置
      api.get('/albums'),     // 获取相册列表
      api.get('/collections') // 获取集合列表
    ])

    // 合并设置数据，注意 allow_registration 需要从字符串转换为布尔值
    Object.assign(settings, {
      ...settRes.data,
      allow_registration: settRes.data.allow_registration === 'true'
    })

    // 更新相册列表
    albums.value = albRes.data

    // 更新集合列表
    collections.value = colRes.data
  } catch (e) {
    // 错误处理已省略，静默失败
  }
}

// ==================== 表单提交处理 ====================

/**
 * 更新系统设置
 * 将表单数据提交到后端，并刷新公共设置缓存
 */
const handleUpdateSettings = async () => {
  savingSettings.value = true
  try {
    // 注意：后端要求 allow_registration 为字符串类型，所以需要转换
    await api.put('/settings', {
      ...settings,
      allow_registration: String(settings.allow_registration)
    })
    // 刷新公共设置缓存，确保其他组件能获取到最新设置
    await settingsStore.fetchPublicSettings()
    showNotify(t('admin.settingsUpdated'))
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.settingsFailed'), 'error')
  } finally {
    savingSettings.value = false
  }
}

/**
 * 创建新相册
 * 提交后清空表单并刷新相册列表
 */
const handleCreateAlbum = async () => {
  try {
    await api.post('/albums', newAlbum)
    // 清空表单
    newAlbum.name = ''
    newAlbum.description = ''
    showNotify(t('admin.albumCreated'))
    // 刷新相册列表
    fetchData()
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.failedCreateAlbum'), 'error')
  }
}

/**
 * 创建新集合
 * 提交后清空表单并刷新集合列表
 */
const handleCreateCollection = async () => {
  try {
    await api.post('/collections', newCollection)
    // 清空表单
    newCollection.name = ''
    newCollection.description = ''
    showNotify(t('admin.collectionCreated'))
    // 刷新集合列表
    fetchData()
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.failedCreateCollection'), 'error')
  }
}

/**
 * 将相册或集合添加到指定集合中
 * 用于组织管理功能
 */
const handleAddToCollection = async () => {
  try {
    await api.post('/collections/add', organizeForm)
    showNotify(t('admin.addedSuccess'))
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.failedAdd'), 'error')
  }
}

/**
 * 更新个人资料和修改密码
 * 包含密码验证逻辑
 * 如果修改了密码，需要重新登录
 */
const handleUpdateProfile = async () => {
  // 验证两次输入的新密码是否一致
  if (profile.password && profile.password !== profile.confirmPassword) {
    showNotify(t('admin.passwordMismatch'), 'error')
    return
  }
  try {
    const res = await api.put('/admin/update', profile)
    if (res.data.passwordChanged) {
      // 如果修改了密码，提示用户并跳转登录页
      showNotify(t('admin.changePasswordSuccess'))
      setTimeout(() => {
        userStore.logout()
        router.push('/login')
      }, 2000)
    } else {
      showNotify(t('admin.profileUpdated'))
    }
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.profileFailed'), 'error')
  }
}

// ==================== 生命周期 ====================

// 组件挂载时获取初始数据
onMounted(fetchData)
</script>

<style scoped>
.max-w-800 {
  max-width: 800px;
}
</style>
