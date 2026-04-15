<!--
  Collection.vue - 收藏夹详情页面

  此页面用于展示单个收藏夹的详细信息，包括收藏夹名称、描述以及其包含的子项目。
  子项目可以是子收藏夹或相册，通过item.type区分。
  支持管理员用户编辑收藏夹名称和描述。
-->
<template>
  <!-- 主容器：仅当收藏夹数据加载完成后显示 -->
  <div v-if="collection">
    <!-- 页面头部：包含返回按钮、收藏夹名称及管理员设置入口 -->
    <div class="d-flex align-center mb-4">
      <!-- 返回按钮：点击返回上一页 -->
      <v-btn icon="mdi-arrow-left" variant="text" @click="$router.back()"></v-btn>
      <!-- 收藏夹名称 -->
      <h1 class="ml-2">{{ collection.name }}</h1>
      <!-- 管理员专属：编辑收藏夹按钮 -->
      <v-btn v-if="isAdmin" icon="mdi-cog" variant="text" class="ml-2" @click="showEditDialog = true"></v-btn>
    </div>

    <!-- 收藏夹描述文本 -->
    <p class="mb-6 text-body-1">{{ collection.description }}</p>

    <!-- 编辑对话框：用于修改收藏夹名称和描述（仅管理员可见） -->
    <v-dialog v-model="showEditDialog" max-width="500">
      <v-card class="rounded-lg">
        <!-- 对话框标题栏 -->
        <v-card-title class="pa-4 bg-primary text-white">
          {{ $t('common.edit') }} {{ $t('admin.collection') }}
        </v-card-title>
        <!-- 对话框内容区域：包含编辑表单 -->
        <v-card-text class="pa-6">
          <!-- 收藏夹名称输入框 -->
          <v-text-field
            v-model="editForm.name"
            :label="$t('admin.name')"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          ></v-text-field>
          <!-- 收藏夹描述输入框（多行文本） -->
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

    <!-- 子项目网格布局：响应式展示收藏夹中的所有子项目（子收藏夹或相册） -->
    <v-row>
      <v-col
        v-for="item in collection.children"
        :key="item.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <!-- 子项目卡片：根据类型跳转到对应详情页 -->
        <v-card :to="item.type === 'album' ? `/album/${item.id}` : `/collection/${item.id}`">
          <!-- 有封面图片时显示图片 -->
          <v-img
            v-if="item.cover_image"
            :src="getImageUrl(item.cover_image)"
            height="150"
            cover
          ></v-img>
          <!-- 无封面图片时显示占位图标 -->
          <v-sheet
            v-else
            height="150"
            :color="item.type === 'album' ? 'grey-lighten-2' : 'surface-variant'"
            class="d-flex align-center justify-center"
          >
            <!-- 相册类型显示图片图标，收藏夹类型显示文件夹图标 -->
            <v-icon size="48" color="grey">{{ item.type === 'album' ? 'mdi-image' : 'mdi-folder' }}</v-icon>
          </v-sheet>
          <!-- 子项目名称标题 -->
          <v-card-title>
             <!-- 收藏夹类型在名称前显示文件夹图标 -->
             <v-icon v-if="item.type === 'collection'" start size="small">mdi-folder-outline</v-icon>
             {{ item.name }}
          </v-card-title>
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
 * Collection.vue 逻辑说明
 *
 * 此页面组件负责：
 * 1. 从路由参数获取收藏夹ID，通过API加载收藏夹详情
 * 2. 管理编辑对话框的显示/隐藏状态
 * 3. 处理收藏夹信息的更新请求
 * 4. 根据用户角色控制管理员功能的可见性
 * 5. 展示收藏夹中的子项目（子收藏夹或相册），并支持导航
 */

// ---------- 依赖导入 ----------
import { ref, reactive, onMounted, watch, computed } from 'vue'
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
 * collection - 收藏夹详情数据
 * 包含收藏夹的名称、描述、以及子项目列表(children)等信息，初始化为null
 */
const collection = ref(null)

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
 * 包含name(收藏夹名称)和description(收藏夹描述)两个字段
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
 * fetchData - 加载收藏夹详情数据
 *
 * 功能说明：
 * - 从路由参数获取收藏夹ID
 * - 调用API获取收藏夹详情
 * - 成功响应后更新collection数据
 * - 同时初始化编辑表单的默认值
 *
 * API端点：GET /collections/{id}
 */
const fetchData = async () => {
  try {
    const res = await api.get(`/collections/${route.params.id}`)
    collection.value = res.data
    editForm.name = res.data.name
    editForm.description = res.data.description
  } catch (e) {
    // 错误处理：静默处理，保留当前数据状态
  }
}

/**
 * handleUpdate - 处理收藏夹信息更新提交
 *
 * 功能说明：
 * - 首先设置saving标志防止重复提交
 * - 调用PUT API更新收藏夹信息
 * - 成功后关闭编辑对话框，显示成功提示，并刷新数据
 * - 失败时显示错误提示信息
 *
 * API端点：PUT /collections/{id}
 * 请求体：{ name, description }
 */
const handleUpdate = async () => {
  saving.value = true
  try {
    await api.put(`/collections/${route.params.id}`, editForm)
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

// ---------- 生命周期钩子 ----------

/**
 * onMounted - 组件挂载完成时执行
 * 功能：初始化加载数据
 */
onMounted(fetchData)

/**
 * watch - 监听路由参数变化
 * 功能：当路由中的收藏夹ID发生变化时（如从收藏夹A切换到收藏夹B），自动重新加载数据
 */
watch(() => route.params.id, fetchData)
</script>
