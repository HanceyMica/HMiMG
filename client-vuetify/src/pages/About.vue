<template>
  <v-row justify="center">
    <v-col cols="12" md="8">
      <v-card class="pa-6">
        <v-card-title class="text-h4 mb-4 text-center">{{ $t('about.title') }}</v-card-title>
        
        <div class="d-flex justify-center mb-6">
          <v-img
            :src="logoSrc"
            max-width="400"
            alt="HMiMG Slogan"
            contain
          ></v-img>
        </div>

        <v-card-text class="text-body-1">
          <p class="mb-6">{{ $t('about.description') }}</p>
          
          <h3 class="text-h6 mb-2">{{ $t('about.features') }}</h3>
          <v-list density="compact">
            <v-list-item prepend-icon="mdi-check-circle" :title="$t('about.feature1')"></v-list-item>
            <v-list-item prepend-icon="mdi-check-circle" :title="$t('about.feature2')"></v-list-item>
            <v-list-item prepend-icon="mdi-check-circle" :title="$t('about.feature3')"></v-list-item>
            <v-list-item prepend-icon="mdi-check-circle" :title="$t('about.feature4')"></v-list-item>
          </v-list>
          
          <v-divider class="my-6"></v-divider>
          
          <div class="d-flex align-center mb-6">
            <v-avatar size="40" class="mr-3">
              <v-img src="/images/avatar.jpg" alt="HanceyMica"></v-img>
            </v-avatar>
            <span class="text-h6 font-weight-medium">HanceyMica</span>
          </div>

          <div class="d-flex align-center justify-center gap-6">
            <v-btn
              icon="mdi-github"
              variant="text"
              size="large"
              href="https://github.com/HanceyMica/HMiMG"
              target="_blank"
              color="grey-darken-3"
              title="GitHub"
            ></v-btn>
            <v-btn
              icon="mdi-twitter"
              variant="text"
              size="large"
              href="https://x.com/HanceyMica"
              target="_blank"
              color="blue-lighten-2"
              title="X (Twitter)"
            ></v-btn>
            <v-btn
              icon="mdi-television-classic"
              variant="text"
              size="large"
              href="https://space.bilibili.com/7035664"
              target="_blank"
              color="pink-lighten-2"
              title="Bilibili"
            ></v-btn>
          </div>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<style scoped>
.gap-6 {
  gap: 24px;
}
</style>

<script setup>
/**
 * About.vue - 关于页面
 *
 * 功能说明：
 * 1. 展示网站 Logo 和 slogan 图片（支持深色/浅色主题自适应）
 * 2. 显示系统功能特性列表
 * 3. 展示作者信息和社交媒体链接
 *
 * 使用说明：
 * - 纯展示型页面，无需特殊权限
 * - Logo 会根据当前主题自动切换深色/浅色版本
 */

import { computed } from 'vue'
import { useTheme } from 'vuetify'

// 获取当前 Vuetify 主题实例
const theme = useTheme()

/**
 * 根据当前主题判断是否为深色模式
 * 用于动态切换 Logo 样式
 */
const isDark = computed(() => theme.global.current.value.dark)

/**
 * 根据主题返回对应的 Logo 图片路径
 * - 深色模式：使用 slogan_dark.png
 * - 浅色模式：使用 slogan_light.png
 */
const logoSrc = computed(() => isDark.value ? '/images/slogan_dark.png' : '/images/slogan_light.png')
</script>
