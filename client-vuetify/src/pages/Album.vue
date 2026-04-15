<!--
  Album.vue - 相册详情页面

  此页面用于展示单个相册的详细信息，包括相册名称、描述以及其中的所有图片。
  支持管理员用户编辑相册名称和描述，并监听全局的上传事件以在图片上传后自动刷新列表。
-->
<template>
  <!-- 主容器：仅当相册数据加载完成后显示 -->
  <div v-if="album">
    <!-- 页面头部：包含返回按钮、相册名称及管理员设置入口 -->
    <div class="d-flex align-center mb-4">
      <!-- 返回按钮：点击返回上一页 -->
      <v-btn icon="mdi-arrow-left" variant="text" @click="$router.back()"></v-btn>
      <!-- 相册名称 -->
      <h1 class="ml-2">{{ album.name }}</h1>
      <!-- 管理员专属：编辑相册按钮 -->
      <v-btn v-if="isAdmin" icon="mdi-cog" variant="text" class="ml-2" @click="showEditDialog = true"></v-btn>
    </div>

    <!-- 相册描述文本 -->
    <p class="mb-6 text-body-1">{{ album.description }}</p>

    <!-- 编辑对话框：用于修改相册名称和描述（仅管理员可见） -->
    <v-dialog v-model="showEditDialog" max-width="500">
      <v-card class="rounded-lg">
        <!-- 对话框标题栏 -->
        <v-card-title class="pa-4 bg-primary text-white">
          {{ $t('common.edit') }} {{ $t('admin.album') }}
        </v-card-title>
        <!-- 对话框内容区域：包含编辑表单 -->
        <v-card-text class="pa-6">
          <!-- 相册名称输入框 -->
          <v-text-field
            v-model="editForm.name"
            :label="$t('admin.name')"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          ></v-text-field>
          <!-- 相册描述输入框（多行文本） -->
          <v-textarea
            v-model="editForm.description"
            :label="$t('admin.description')"
            variant="outlined"
            rows="3"
            class="mb-4"
          ></v-textarea>
          <!-- 提交更新按钮 -->
          <v-btn color="primary" block size="large" :loading="saving" @click="handleUpdate">
            {{ $t('image.update') }}
          </v-btn>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- 图片网格布局：响应式展示相册中的所有图片 -->
    <v-row>
      <v-col
        v-for="img in images"
        :key="img.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <!-- 图片卡片：点击可跳转到图片详情页 -->
        <v-card :to="`/image/${img.id}`">
          <!-- 图片预览：固定高度200px，cover模式裁剪适应 -->
          <v-img
            :src="getImageUrl(img.path)"
            cover
            height="200"
          ></v-img>
          <!-- 图片副标题：显示原始文件名 -->
          <v-card-subtitle class="pt-2">{{ img.original_name }}</v-card-subtitle>
        </v-card>
      </v-col>
    </v-row>

    <!-- 全局通知提示条：用于显示操作成功或错误信息 -->
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :timeout="3000"
      elevation="24"
      rounded="pill"
    >
      <!-- 提示内容：图标 + 文字 -->
      <div class="d-flex align-center">
        <v-icon start class="mr-2">{{ snackbar.icon }}</v-icon>
        {{ snackbar.text }}
      </div>
      <!-- 关闭按钮 -->
      <template v-slot:actions>
        <v-btn variant="text" @click="snackbar.show = false">Close</v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script setup>
/*
 * Album.vue 逻辑说明
 *
 * 此页面组件负责：
 * 1. 从路由参数获取相册ID，通过API加载相册详情和图片列表
 * 2. 管理编辑对话框的显示/隐藏状态
 * 3. 处理相册信息的更新请求
 * 4. 监听全局图片上传事件，在当前相册有图片上传时自动刷新列表
 * 5. 根据用户角色控制管理员功能的可见性
 */

// ---------- 依赖导入 ----------
import { ref, reactive, onMounted, onBeforeUnmount, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useI18n } from 'vue-i18n'
import api, { buildUploadedFileUrl } from '@/lib/api'

// ---------- 路由与状态管理 ----------
const route = useRoute()           // 当前路由实例，用于获取URL参数
const userStore = useUserStore()   // 用户状态存储，提供用户信息和角色
const { t } = useI18n()            // 国际化函数，用于翻译文本

// ---------- 计算属性 ----------

/**
 * isAdmin - 判断当前用户是否为管理员
 * 依赖用户Store中的role字段，角色为'admin'时返回true
 */
const isAdmin = computed(() => userStore.user?.role === 'admin')

// ---------- 响应式数据 ----------

/**
 * album - 相册详情数据
 * 包含相册的名称、描述等信息，初始化为null
 */
const album = ref(null)

/**
 * images - 相册中的图片列表
 * 数组类型，每个元素包含图片的id、path、original_name等属性
 */
const images = ref([])

/**
 * showEditDialog - 控制编辑对话框的显示与隐藏
 * true时显示编辑对话框，false时隐藏
 */
const showEditDialog = ref(false)

/**
 * saving - 防止重复提交标志
 * 当为true时表示正在保存数据，禁用提交按钮
 */
const saving = ref(false)

/**
 * snackbar - 通知提示条的状态管理
 * 包含show(是否显示)、text(提示文本)、color(颜色)、icon(图标)
 */
const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
  icon: 'mdi-check-circle'
})

/**
 * editForm - 编辑表单数据
 * 包含name(相册名称)和description(相册描述)两个字段
 */
const editForm = reactive({
  name: '',
  description: ''
})

// ---------- 方法 ----------

/**
 * showNotify - 显示通知提示条
 * @param {string} text - 提示文本内容
 * @param {string} color - 提示颜色，默认为'success'成功绿色，可传'error'错误红色
 *
 * 功能说明：
 * - 根据color参数设置提示图标的类型
 * - 成功时显示绿色勾选图标，错误时显示橙色警告图标
 * - 自动设置show为true以显示提示条
 */
const showNotify = (text, color = 'success') => {
  snackbar.text = text
  snackbar.color = color
  snackbar.icon = color === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'
  snackbar.show = true
}

/**
 * fetchData - 加载相册数据和图片列表
 *
 * 功能说明：
 * - 从路由参数获取相册ID
 * - 并行发起两个API请求获取相册详情和图片列表
 * - 成功响应后更新album和images数据
 * - 同时初始化编辑表单的默认值
 *
 * API端点：
 * - GET /albums/{id} - 获取相册详情
 * - GET /images?albumId={id} - 获取指定相册的所有图片
 */
const fetchData = async () => {
  try {
    const id = route.params.id
    const [albRes, imgRes] = await Promise.all([
      api.get(`/albums/${id}`),
      api.get(`/images?albumId=${id}`)
    ])
    album.value = albRes.data
    editForm.name = albRes.data.name
    editForm.description = albRes.data.description
    images.value = imgRes.data
  } catch (e) {
    // 错误处理：静默处理，保留当前数据状态
  }
}

/**
 * handleUpdate - 处理相册信息更新提交
 *
 * 功能说明：
 * - 首先设置saving标志防止重复提交
 * - 调用PUT API更新相册信息
 * - 成功后关闭编辑对话框，显示成功提示，并刷新数据
 * - 失败时显示错误提示信息
 *
 * API端点：PUT /albums/{id}
 * 请求体：{ name, description }
 */
const handleUpdate = async () => {
  saving.value = true
  try {
    await api.put(`/albums/${route.params.id}`, editForm)
    showEditDialog.value = false
    showNotify(t('image.updateSuccess'))
    fetchData()
  } catch (e) {
    showNotify(e.response?.data?.error || t('image.updateFailed'), 'error')
  } finally {
    saving.value = false
  }
}

/**
 * getImageUrl - 构建图片的完整访问URL
 * @param {string} path - 图片的相对路径
 * @returns {string} 图片的完整访问URL
 *
 * 功能说明：
 * - 使用buildUploadedFileUrl将相对路径转换为完整的可访问URL
 * - 解决图片路径在不同环境下可能无法正确访问的问题
 */
const getImageUrl = (path) => {
  return buildUploadedFileUrl(path)
}

/**
 * handleImagesUploaded - 处理全局图片上传成功事件
 * @param {CustomEvent} event - 上传事件，包含detail.albumId信息
 *
 * 功能说明：
 * - 监听全局的'hmimg:images-uploaded'事件
 * - 检查上传的图片是否属于当前相册
 * - 如果是当前相册的图片上传成功，则自动刷新列表
 *
 * 使用场景：
 * - 用户在其他页面（如上传页面）上传图片后，
 * - 当前页面能自动感知并刷新，无需手动刷新
 */
const handleImagesUploaded = (event) => {
  const uploadedAlbumId = String(event.detail?.albumId || '')
  if (uploadedAlbumId && uploadedAlbumId === String(route.params.id)) {
    fetchData()
  }
}

// ---------- 生命周期钩子 ----------

/**
 * onMounted - 组件挂载完成时执行
 * 功能：初始化加载数据，并注册全局事件监听器
 */
onMounted(() => {
  fetchData()
  // 注册全局图片上传事件监听器
  window.addEventListener('hmimg:images-uploaded', handleImagesUploaded)
})

/**
 * onBeforeUnmount - 组件卸载前执行
 * 功能：移除全局事件监听器，防止内存泄漏
 */
onBeforeUnmount(() => {
  window.removeEventListener('hmimg:images-uploaded', handleImagesUploaded)
})

/**
 * watch - 监听路由参数变化
 * 功能：当路由中的相册ID发生变化时（如从相册A切换到相册B），自动重新加载数据
 */
watch(() => route.params.id, fetchData)
</script>
