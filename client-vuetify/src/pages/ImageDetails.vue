<!--
  图片详情页组件 (ImageDetails.vue)

  功能说明：
  - 展示单张图片的详细信息，包括原始文件名、大小、类型、上传时间等
  - 支持图片下载功能
  - 支持图片重命名功能（仅管理员或上传者可用）
  - 支持图片删除功能（仅管理员可用）
  - 支持在相册内切换上一张/下一张图片
  - 使用全局 Snackbar 组件显示操作结果通知

  页面布局：
  - 顶部：返回按钮和上一张/下一张导航按钮
  - 主体：左侧大图展示区域，右侧信息卡片（含详情和操作按钮）
-->
<template>
  <!-- 当图片数据加载完成后显示页面内容 -->
  <div v-if="image">
    <!-- 顶部导航栏：返回按钮 + 图片导航按钮 -->
    <div class="d-flex align-center justify-space-between mb-4 flex-wrap ga-2">
      <!-- 返回上一页按钮 -->
      <v-btn variant="text" prepend-icon="mdi-arrow-left" @click="router.back()">
        {{ $t('image.back') }}
      </v-btn>

      <!-- 上一张/下一张图片导航按钮组 -->
      <div class="d-flex ga-2">
        <!-- 上一张按钮：当前为第一张时禁用 -->
        <v-btn
          variant="outlined"
          prepend-icon="mdi-chevron-left"
          :disabled="!previousImage"
          @click="goToImage(previousImage)"
        >
          {{ $t('image.previous') }}
        </v-btn>
        <!-- 下一张按钮：当前为最后一张时禁用 -->
        <v-btn
          variant="outlined"
          append-icon="mdi-chevron-right"
          :disabled="!nextImage"
          @click="goToImage(nextImage)"
        >
          {{ $t('image.next') }}
        </v-btn>
      </div>
    </div>

    <!-- 主内容区域：图片展示 + 详情信息 -->
    <v-row>
      <!-- 左侧：图片展示区域（占8/12宽度，中等屏幕以上） -->
      <v-col cols="12" md="8">
        <v-img
          :src="getImageUrl(image.path)"
          class="rounded-lg elevation-4"
          max-height="80vh"
        ></v-img>
      </v-col>

      <!-- 右侧：图片详情卡片（占4/12宽度，中等屏幕以上） -->
      <v-col cols="12" md="4">
        <v-card class="fill-height">
          <!-- 卡片标题：详情 -->
          <v-card-title>{{ $t('image.details') }}</v-card-title>

          <!-- 详情列表：展示图片各项元数据 -->
          <v-list density="compact">
            <v-list-item :title="$t('image.originalName')" :subtitle="image.original_name"></v-list-item>
            <v-list-item :title="$t('image.size')" :subtitle="formatSize(image.size)"></v-list-item>
            <v-list-item :title="$t('image.type')" :subtitle="image.mimetype"></v-list-item>
            <v-list-item :title="$t('image.album')" :subtitle="image.album_name || $t('home.noAlbum')"></v-list-item>
            <v-list-item :title="$t('image.uploadedAt')" :subtitle="formatDate(image.created_at)"></v-list-item>
          </v-list>

          <!-- 重命名表单：仅管理员或上传者可见 -->
          <v-card-text v-if="canEditName" class="pt-0">
            <v-text-field
              v-model="renameForm.originalName"
              :label="$t('image.originalName')"
              variant="outlined"
              density="comfortable"
              hide-details
              class="mb-3"
            ></v-text-field>
            <v-btn
              color="primary"
              variant="outlined"
              block
              :loading="savingName"
              @click="handleRename"
            >
              {{ $t('image.update') }}
            </v-btn>
          </v-card-text>

          <!-- 卡片底部操作按钮区域 -->
          <v-card-actions class="pa-4 flex-column align-stretch">
            <!-- 下载按钮 -->
            <v-btn
              color="primary"
              variant="elevated"
              class="mb-2"
              prepend-icon="mdi-download"
              @click="handleDownload"
            >
              {{ $t('image.download') }}
            </v-btn>

            <!-- 删除按钮：仅管理员可见 -->
            <v-btn
              v-if="isAdmin"
              color="grey-darken-1"
              variant="outlined"
              prepend-icon="mdi-delete"
              @click="handleDelete"
            >
              {{ $t('image.delete') }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- 全局消息提示条：用于显示操作成功/失败等反馈信息 -->
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
/*
 * ImageDetails.vue - 图片详情页面逻辑
 *
 * 主要功能：
 * 1. 从 URL 参数获取图片 ID，通过 API 获取图片详情
 * 2. 根据图片所属相册 ID 获取同相册图片列表，用于上下张导航
 * 3. 提供图片下载、重命名、删除等操作
 * 4. 管理页面级别的通知消息（snackbar）
 */

import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/store/user'
import api, { buildUploadedFileUrl } from '@/lib/api'

// =============================================
// 路由与 Store 相关
// =============================================

const route = useRoute()           // 当前路由对象，用于获取 URL 参数（图片 ID）
const router = useRouter()         // 路由实例，用于页面导航
const { t } = useI18n()            // 国际化函数，用于翻译文本
const userStore = useUserStore()   // 用户状态管理，包含当前用户信息

// =============================================
// 响应式状态（Ref）
// =============================================

/**
 * image: 当前查看的图片详情对象
 * 结构示例：{ id, original_name, path, size, mimetype, album_id, album_name, uploaded_by, created_at }
 */
const image = ref(null)

/**
 * imageList: 当前相册中的所有图片列表，用于实现上下张切换
 * 每个元素结构与 image 相同
 */
const imageList = ref([])

/**
 * savingName: 重命名操作进行中的加载状态标志
 * true 时重命名按钮显示加载动画
 */
const savingName = ref(false)

// =============================================
// 计算属性（Computed）
// =============================================

/**
 * isAdmin: 判断当前用户是否为管理员
 * 用于控制管理员专属功能（如删除图片）的显示
 */
const isAdmin = computed(() => userStore.user?.role === 'admin')

/**
 * currentImageIndex: 当前图片在 imageList 中的索引位置
 * 用于确定上一张和下一张图片
 */
const currentImageIndex = computed(() =>
  imageList.value.findIndex(item => String(item.id) === String(route.params.id))
)

/**
 * previousImage: 上一张图片对象
 * 如果当前为第一张或不存在则返回 null
 */
const previousImage = computed(() => {
  if (currentImageIndex.value <= 0) return null
  return imageList.value[currentImageIndex.value - 1] || null
})

/**
 * nextImage: 下一张图片对象
 * 如果当前为最后一张或不存在则返回 null
 */
const nextImage = computed(() => {
  if (currentImageIndex.value < 0 || currentImageIndex.value >= imageList.value.length - 1) return null
  return imageList.value[currentImageIndex.value + 1] || null
})

/**
 * canEditName: 判断当前用户是否有权限重命名此图片
 * 条件：用户已登录且（是管理员或是图片上传者）
 */
const canEditName = computed(() => {
  if (!image.value || !userStore.user) return false
  return userStore.user.role === 'admin' || Number(image.value.uploaded_by) === Number(userStore.user.id)
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
 * 使用示例：showNotify('操作成功！') 或 showNotify('出错了！', 'error')
 */
const showNotify = (text, color = 'success') => {
  snackbar.text = text
  snackbar.color = color
  snackbar.icon = color === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'
  snackbar.show = true
}

// =============================================
// 重命名表单状态
// =============================================

/**
 * renameForm: 重命名表单的数据对象
 * - originalName: 用户输入的新文件名
 *
 * 注意：此表单为响应式，绑定到 v-model
 */
const renameForm = reactive({
  originalName: ''
})

// =============================================
// API 数据获取
// =============================================

/**
 * fetchData: 异步获取图片详情和相册图片列表
 *
 * 执行流程：
 * 1. 调用 GET /images/:id 获取当前图片详情
 * 2. 根据返回的 album_id 获取同相册所有图片
 * 3. 更新 image 和 imageList 响应式变量
 *
 * 错误处理：静默处理，不显示错误（页面会保持空白）
 */
const fetchData = async () => {
  try {
    // 获取当前图片详情
    const res = await api.get(`/images/${route.params.id}`)
    image.value = res.data
    // 设置重命名表单的初始值为当前图片名
    renameForm.originalName = res.data.original_name || ''

    // 获取同相册的所有图片，用于上下张导航
    const listRes = await api.get(`/images?albumId=${res.data.album_id}`)
    imageList.value = listRes.data
  } catch (e) {
    // 静默处理错误
  }
}

// =============================================
// 辅助函数
// =============================================

/**
 * getImageUrl: 根据图片路径构建完整的访问 URL
 * @param {string} path - 图片相对路径
 * @returns {string} 完整的图片访问 URL
 *
 * 使用 api 模块的 buildUploadedFileUrl 函数处理路径
 */
const getImageUrl = (path) => {
  return buildUploadedFileUrl(path)
}

/**
 * formatSize: 将字节数格式化为人类可读的大小字符串
 * @param {number} bytes - 字节数
 * @returns {string} 格式化后的大小字符串，如 "1.5 MB"
 *
 * 支持的单位：B, KB, MB, GB，自动选择最合适的单位
 */
const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * formatDate: 将日期字符串格式化为本地化的日期时间字符串
 * @param {string} date - ISO 格式的日期字符串
 * @returns {string} 本地化的日期时间，如 "2024/1/15 下午3:30:00"
 */
const formatDate = (date) => new Date(date).toLocaleString()

// =============================================
// 事件处理函数
// =============================================

/**
 * handleDownload: 处理图片下载操作
 *
 * 实现方式：创建一个隐藏的 <a> 标签，设置 href 和 download 属性
 * 然后触发点击事件进行下载
 *
 * 注意：download 属性指定下载后的文件名
 */
const handleDownload = () => {
  const url = getImageUrl(image.value.path)
  const link = document.createElement('a')
  link.href = url
  link.download = image.value.original_name
  link.click()
}

/**
 * goToImage: 导航到指定图片的详情页
 * @param {Object} target - 目标图片对象，必须包含 id 属性
 *
 * 如果 target 存在且有 id，则使用 router.push 跳转到该图片的详情页
 */
const goToImage = (target) => {
  if (target?.id) {
    router.push(`/image/${target.id}`)
  }
}

/**
 * handleRename: 处理图片重命名请求
 *
 * 验证：检查新文件名是否为空
 * 请求：发送 PUT /images/:id 请求更新 original_name
 * 成功：更新本地状态和列表中的图片名，显示成功消息
 * 失败：显示错误消息
 */
const handleRename = async () => {
  // 验证：文件名不能为空
  if (!renameForm.originalName.trim()) {
    showNotify(t('image.nameRequired'), 'error')
    return
  }

  savingName.value = true  // 显示加载状态
  try {
    // 发送更新请求
    const res = await api.put(`/images/${route.params.id}`, {
      original_name: renameForm.originalName.trim()
    })

    // 更新当前图片对象，保留 album_name 等字段
    image.value = {
      ...image.value,
      ...res.data,
      album_name: image.value?.album_name
    }

    // 更新图片列表中对应项的名称（用于上下张导航时显示正确的名称）
    imageList.value = imageList.value.map(item => (
      String(item.id) === String(route.params.id)
        ? { ...item, original_name: res.data.original_name }
        : item
    ))

    showNotify(t('image.updateSuccess'))  // 显示成功消息
  } catch (e) {
    // 显示错误消息，优先使用后端返回的错误信息
    showNotify(e.response?.data?.error || t('image.updateFailed'), 'error')
  } finally {
    savingName.value = false  // 隐藏加载状态
  }
}

/**
 * handleDelete: 处理图片删除请求
 *
 * 流程：
 * 1. 弹出确认对话框，用户确认后继续
 * 2. 发送 DELETE /images/:id 请求删除图片
 * 3. 删除成功后返回上一页
 * 4. 失败则显示错误消息
 */
const handleDelete = async () => {
  // 弹出确认对话框
  if (!confirm('Are you sure you want to delete this image?')) return

  try {
    await api.delete(`/images/${route.params.id}`)
    router.back()  // 删除成功后返回上一页
  } catch (e) {
    showNotify(e.response?.data?.error || 'Failed to delete image', 'error')
  }
}

// =============================================
// 生命周期钩子
// =============================================

/**
 * onMounted: 组件挂载后立即获取数据
 * 调用 fetchData() 获取图片详情和列表
 */
onMounted(fetchData)

/**
 * watch: 监听路由参数变化，当图片 ID 改变时重新获取数据
 *
 * 场景：用户通过上下张按钮切换图片时，URL 参数会改变
 * 此时需要重新获取新图片的详情数据
 */
watch(() => route.params.id, fetchData)
</script>
