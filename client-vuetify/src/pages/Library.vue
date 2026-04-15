<template>
  <div>
    <!-- 页面标题区域：显示相册库图标和国际化标题 -->
    <div class="d-flex align-center mb-6">
      <v-icon size="32" color="primary" class="mr-3">mdi-image-multiple-outline</v-icon>
      <h1 class="text-h4 font-weight-bold">{{ $t('common.library') }}</h1>
    </div>

    <!-- 标签页切换区域：支持在集合、相册和照片三种视图间切换 -->
    <v-tabs v-model="tab" color="primary" class="mb-6 border-b">
      <v-tab value="collections" prepend-icon="mdi-folder-multiple">{{ $t('common.collections') }}</v-tab>
      <v-tab value="albums" prepend-icon="mdi-image-album">{{ $t('common.albums') }}</v-tab>
      <v-tab value="photos" prepend-icon="mdi-image-multiple">{{ $t('common.photos') }}</v-tab>
    </v-tabs>

    <!-- 标签页内容区域：根据选中的标签页显示对应内容 -->
    <v-window v-model="tab">

      <!-- 集合（Collections）视图：展示所有图片集合 -->
      <v-window-item value="collections">
        <v-row>
          <v-col v-for="c in collections" :key="c.id" cols="12" sm="6" md="4" lg="3">
            <!-- 集合卡片：悬停时阴影增强，点击跳转到集合详情页 -->
            <v-hover v-slot="{ isHovering, props }">
              <v-card
                v-bind="props"
                :elevation="isHovering ? 8 : 1"
                :to="`/collection/${c.id}`"
                class="rounded-lg transition-swing border"
                flat
              >
                <!-- 卡片头部：显示文件夹图标和集合名称 -->
                <v-card-item class="bg-surface-variant text-white">
                  <template v-slot:prepend>
                    <v-icon size="32">mdi-folder-outline</v-icon>
                  </template>
                  <v-card-title class="text-subtitle-1">{{ c.name }}</v-card-title>
                </v-card-item>
                <!-- 卡片内容：显示集合描述 -->
                <v-card-text class="pa-4 text-body-2 text-grey-darken-1">
                  {{ c.description || 'No description' }}
                </v-card-text>
              </v-card>
            </v-hover>
          </v-col>
        </v-row>
      </v-window-item>

      <!-- 相册（Albums）视图：展示所有相册，包含封面图和基本信息 -->
      <v-window-item value="albums">
        <v-row>
          <v-col v-for="a in albums" :key="a.id" cols="12" sm="6" md="4" lg="3">
            <!-- 相册卡片：悬停时阴影增强，点击跳转到相册详情页 -->
            <v-hover v-slot="{ isHovering, props }">
              <v-card
                v-bind="props"
                :elevation="isHovering ? 8 : 1"
                :to="`/album/${a.id}`"
                class="rounded-lg overflow-hidden transition-swing border"
                flat
              >
                <!-- 相册封面图：如果存在封面图则显示，否则显示占位图标 -->
                <v-img
                  v-if="a.cover_image"
                  :src="getImageUrl(a.cover_image)"
                  height="180"
                  cover
                  class="bg-grey-lighten-2"
                >
                  <!-- 图片加载中的占位动画 -->
                  <template v-slot:placeholder>
                    <v-row class="fill-height ma-0" align="center" justify="center">
                      <v-progress-circular indeterminate color="grey-lighten-4"></v-progress-circular>
                    </v-row>
                  </template>
                </v-img>
                <!-- 无封面图时的占位显示 -->
                <v-sheet v-else height="180" color="grey-lighten-4" class="d-flex align-center justify-center">
                  <v-icon size="48" color="grey-lighten-1">mdi-image-outline</v-icon>
                </v-sheet>
                <!-- 相册名称 -->
                <v-card-title class="text-subtitle-1 font-weight-bold">{{ a.name }}</v-card-title>
                <!-- 相册描述（单行截断） -->
                <v-card-text class="text-caption text-grey text-truncate pt-0">
                  {{ a.description || 'No description' }}
                </v-card-text>
              </v-card>
            </v-hover>
          </v-col>
        </v-row>
      </v-window-item>

      <!-- 照片（Photos）视图：展示所有单独的图片 -->
      <v-window-item value="photos">
        <v-row>
          <v-col
            v-for="img in images"
            :key="img.id"
            cols="12"
            sm="6"
            md="4"
            lg="3"
          >
            <!-- 照片卡片：悬停时阴影增强，点击跳转到图片详情页 -->
            <v-hover v-slot="{ isHovering, props }">
              <v-card
                v-bind="props"
                :elevation="isHovering ? 8 : 1"
                :to="`/image/${img.id}`"
                class="rounded-lg overflow-hidden transition-swing border"
                flat
              >
                <!-- 图片缩略图 -->
                <v-img
                  :src="getImageUrl(img.path)"
                  cover
                  height="180"
                  class="bg-grey-lighten-2"
                >
                  <!-- 图片加载中的占位动画 -->
                  <template v-slot:placeholder>
                    <v-row class="fill-height ma-0" align="center" justify="center">
                      <v-progress-circular indeterminate color="grey-lighten-4"></v-progress-circular>
                    </v-row>
                  </template>
                </v-img>
                <!-- 图片原始文件名 -->
                <v-card-title class="text-subtitle-2 text-truncate pa-3">
                  {{ img.original_name }}
                </v-card-title>
              </v-card>
            </v-hover>
          </v-col>
        </v-row>
      </v-window-item>

    </v-window>
  </div>
</template>

<script setup>
/**
 * Library.vue - 相册库页面组件
 *
 * 功能说明：
 * 1. 以标签页形式展示三种视图：集合（Collections）、相册（Albums）、照片（Photos）
 * 2. 集合视图以文件夹形式展示所有集合
 * 3. 相册视图以卡片形式展示所有相册，包含封面图和描述信息
 * 4. 照片视图以卡片形式展示所有单独的图片
 * 5. 监听全局自定义事件 hmimg:images-uploaded，在图片上传成功后自动刷新数据
 */

import { ref, onMounted, onBeforeUnmount } from 'vue'
import api, { buildUploadedFileUrl } from '@/lib/api'

/**
 * 当前选中的标签页
 * 可选值：'collections'（集合）、'albums'（相册）、'photos'（照片）
 * 默认值为 'collections'
 */
const tab = ref('collections')

// 相册列表数据
const albums = ref([])

// 集合列表数据
const collections = ref([])

// 单独的照片列表数据
const images = ref([])

/**
 * 获取所有数据（相册、集合、照片）
 * 使用 Promise.all 并行请求三个接口，提高加载效率
 */
const fetchData = async () => {
  try {
    // 并行请求相册、集合和照片三个接口
    const [albRes, colRes, imgRes] = await Promise.all([
      api.get('/albums'),
      api.get('/collections'),
      api.get('/images')
    ])
    // 更新各个数据响应式变量
    albums.value = albRes.data
    collections.value = colRes.data
    images.value = imgRes.data
  } catch (e) {
    // 错误处理：静默处理，任何接口失败都不影响其他数据展示
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
 * 处理图片上传成功事件的回调函数
 * 当有图片上传成功时，刷新页面数据以显示新上传的图片
 */
const handleImagesUploaded = () => {
  fetchData()
}

/**
 * 组件挂载时的生命周期钩子
 * 1. 获取初始数据（相册、集合、照片列表）
 * 2. 注册全局事件监听器，监听图片上传成功事件
 */
onMounted(() => {
  fetchData()
  // 监听全局图片上传成功事件，Home.vue 上传图片后会派发此事件
  window.addEventListener('hmimg:images-uploaded', handleImagesUploaded)
})

/**
 * 组件卸载前的生命周期钩子
 * 移除全局事件监听器，防止内存泄漏
 */
onBeforeUnmount(() => {
  window.removeEventListener('hmimg:images-uploaded', handleImagesUploaded)
})
</script>

<style scoped>
/* 卡片悬停过渡动画效果 */
.transition-swing {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}
</style>
