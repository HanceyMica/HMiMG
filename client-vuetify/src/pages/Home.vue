<template>
  <div class="text-center">
    <!-- 英雄区域：网站欢迎信息与上传入口 -->
    <div class="hero-section py-8 py-sm-16">
      <!-- 网站欢迎标题，根据 settingsStore 中的 websiteTitle 动态显示 -->
      <p class="text-h4 text-sm-h3 font-weight-bold mb-4 mt-4 mt-sm-8">
        {{ $t('home.welcome', { title: settingsStore.websiteTitle }) }}
      </p>

      <!-- 上传卡片区域：用户可以点击或拖拽文件进行上传 -->
      <v-row justify="center">
        <v-col cols="12" md="8" lg="6">
          <v-hover v-slot="{ isHovering, props }">
            <v-card
              v-bind="props"
              :elevation="isHovering ? 8 : 2"
              class="pa-8 mb-12 border-dashed rounded-xl cursor-pointer upload-card"
              @click="triggerFileInput"
            >
              <!-- 上传图标 -->
              <v-icon size="64" color="primary" class="mb-4">mdi-cloud-upload-outline</v-icon>
              <!-- 上传提示文字 -->
              <div class="text-h6 mb-2 upload-text">{{ $t('home.dragDropText') }}</div>
              <div class="text-body-2 text-grey-darken-1">{{ $t('home.dragDropHint') }}</div>

              <!-- 隐藏的文件输入框，用于选择要上传的图片文件 -->
              <v-file-input
                ref="fileInput"
                v-model="files"
                multiple
                class="d-none"
                @update:model-value="onFilesSelected"
              ></v-file-input>
            </v-card>
          </v-hover>
        </v-col>
      </v-row>

      <!-- 导航按钮区域：跳转到相册库或管理后台（仅管理员可见） -->
      <div class="d-flex justify-center gap-4">
        <!-- 相册库按钮 -->
        <v-btn
          to="/library"
          color="primary"
          size="large"
          prepend-icon="mdi-image-multiple"
          class="text-none px-8 rounded-pill"
        >
          {{ $t('common.library') }}
        </v-btn>
        <!-- 管理后台按钮（仅管理员可见） -->
        <v-btn
          v-if="isAdmin"
          to="/admin"
          variant="outlined"
          size="large"
          prepend-icon="mdi-cog"
          class="text-none px-8 rounded-pill"
        >
          {{ $t('common.admin') }}
        </v-btn>
      </div>
    </div>

    <!-- 上传弹窗：选择目标相册并确认上传 -->
    <v-dialog v-model="showUploadModal" max-width="500">
      <v-card class="rounded-lg">
        <!-- 弹窗标题 -->
        <v-card-title class="pa-4 bg-primary text-white">
          {{ $t('home.selectAlbumToUpload') }}
        </v-card-title>
        <!-- 弹窗内容 -->
        <v-card-text class="pa-6">
          <!-- 已选文件数量提示 -->
          <div class="mb-4">{{ files.length }} files selected</div>
          <!-- 相册选择下拉框 -->
          <v-select
            v-model="uploadForm.albumId"
            :items="albums"
            item-title="name"
            item-value="id"
            :label="$t('home.selectAlbum')"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          ></v-select>
          <!-- 上传确认按钮 -->
          <v-btn
            color="primary"
            block
            size="large"
            :loading="uploading"
            :disabled="!uploadForm.albumId"
            @click="handleUpload"
          >
            {{ $t('home.upload') }}
          </v-btn>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- 全局消息提示条：用于显示上传成功/失败等通知 -->
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
 * Home.vue - 首页组件
 *
 * 功能说明：
 * 1. 展示网站欢迎信息，支持自定义网站标题
 * 2. 提供图片上传入口，用户可选择文件并上传到指定相册
 * 3. 显示导航按钮，可跳转至相册库和管理后台
 * 4. 上传完成后通过自定义事件通知其他组件刷新数据
 */

import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import api, { buildUploadedFileUrl } from '@/lib/api'
import { useUserStore } from '@/store/user'
import { useSettingsStore } from '@/store/settings'
import { useI18n } from 'vue-i18n'

// Router 实例，用于上传成功后跳转到目标相册
const router = useRouter()

// 用户状态管理，包含当前用户信息及角色判断
const userStore = useUserStore()

// 网站设置状态管理，用于获取网站标题等配置信息
const settingsStore = useSettingsStore()

// i18n 翻译函数，用于获取多语言文本
const { t } = useI18n()

/**
 * 计算属性：判断当前用户是否为管理员
 * 只有管理员角色的用户才能看到并访问管理后台
 */
const isAdmin = computed(() => userStore.user?.role === 'admin')

// 相册列表，用于在上传弹窗中选择目标相册
const albums = ref([])

// 已选择的待上传文件列表
const files = ref([])

// 上传状态标志，表示当前是否正在进行文件上传
const uploading = ref(false)

// 控制上传弹窗的显示与隐藏
const showUploadModal = ref(false)

// 文件输入框的 DOM 引用，用于触发点击事件
const fileInput = ref(null)

/**
 * 消息提示条的状态对象
 * - show: 控制提示条是否显示
 * - text: 提示文本内容
 * - color: 提示条颜色（success 表示成功，error 表示错误）
 * - icon: 提示条左侧的图标
 */
const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
  icon: 'mdi-check-circle'
})

/**
 * 显示通知消息的函数
 * @param {string} text - 要显示的消息文本
 * @param {string} color - 消息类型，可选 'success'（成功，绿色）或 'error'（错误，红色）
 */
const showNotify = (text, color = 'success') => {
  snackbar.text = text
  snackbar.color = color
  snackbar.icon = color === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'
  snackbar.show = true
}

/**
 * 上传表单数据
 * - albumId: 用户选择的要上传到的目标相册 ID
 */
const uploadForm = reactive({
  albumId: null
})

/**
 * 获取相册列表数据
 * 同时确保网站设置信息已被加载（如果尚未加载则主动获取）
 */
const fetchData = async () => {
  try {
    // 从服务器获取相册列表
    const res = await api.get('/albums')
    albums.value = res.data

    // 确保设置已加载，如果网站标题仍为默认值则主动获取
    if (settingsStore.websiteTitle === 'HMiMG') {
      await settingsStore.fetchPublicSettings()
    }
  } catch (e) {
    // 错误处理：输出错误信息到控制台
    console.error(e)
  }
}

/**
 * 根据文件路径构建完整的图片访问 URL
 * @param {string} path - 文件的相对路径
 * @returns {string} 完整的文件访问 URL
 */
const getImageUrl = (path) => {
  return buildUploadedFileUrl(path)
}

/**
 * 触发隐藏的文件输入框的点击事件
 * 由于文件输入框被设置为隐藏（class="d-none"），需要通过 JS 模拟点击来打开文件选择器
 */
const triggerFileInput = () => {
  const input = document.querySelector('.v-file-input input')
  if (input) input.click()
}

/**
 * 当用户选择文件后的回调函数
 * @param {File[]} selectedFiles - 用户选择的文件数组
 * 如果有文件被选中，则打开上传弹窗让用户选择目标相册
 */
const onFilesSelected = (selectedFiles) => {
  if (selectedFiles && selectedFiles.length > 0) {
    showUploadModal.value = true
  }
}

/**
 * 处理文件上传的核心函数
 * 1. 将选中的文件和目标相册 ID 打包为 FormData
 * 2. 发送 POST 请求到服务器进行上传
 * 3. 上传成功后派发自定义事件通知其他组件刷新，关闭弹窗，并跳转到目标相册
 */
const handleUpload = async () => {
  // 如果没有选中文件则直接返回
  if (files.value.length === 0) return

  uploading.value = true

  // 构建 FormData 对象，用于发送 multipart/form-data 请求
  const formData = new FormData()
  formData.append('albumId', uploadForm.albumId)
  files.value.forEach(file => {
    formData.append('images', file)
  })

  try {
    // 发送上传请求
    const res = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    const albumId = uploadForm.albumId

    // 重置表单状态
    files.value = []
    uploadForm.albumId = null
    showUploadModal.value = false

    // 派发自定义事件，通知其他组件（如 Library 页面）刷新数据
    window.dispatchEvent(new CustomEvent('hmimg:images-uploaded', {
      detail: {
        albumId,
        ids: res.data?.ids || []
      }
    }))

    // 显示上传成功提示
    showNotify(t('home.uploadSuccess'))

    // 刷新相册列表
    fetchData()

    // 如果上传到了指定相册，则跳转到该相册页面
    if (albumId) {
      router.push(`/album/${albumId}`)
    }
  } catch (e) {
    // 上传失败时显示错误提示
    showNotify(t('home.uploadFailed'), 'error')
  } finally {
    // 无论成功或失败，都重置上传状态
    uploading.value = false
  }
}

// 组件挂载时获取初始数据（相册列表和网站设置）
onMounted(fetchData)
</script>

<style scoped>
/* 虚线边框样式，用于上传卡片区域 */
.border-dashed {
  border: 2px dashed rgba(var(--v-border-color), 0.3) !important;
}

/* 上传卡片的过渡动画效果 */
.upload-card {
  transition: all 0.3s ease;
}

/* 暗色主题下上传提示文字的颜色 */
.v-theme--dark .upload-text {
  color: rgba(255, 255, 255, 0.9) !important;
}

/* 亮色主题下上传提示文字的颜色 */
.v-theme--light .upload-text {
  color: rgba(0, 0, 0, 0.87) !important;
}

/* 按钮之间的间距 */
.gap-4 {
  gap: 16px;
}

/* 悬停动画过渡效果 */
.transition-swing {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}
</style>
