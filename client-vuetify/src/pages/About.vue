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
import { computed } from 'vue'
import { useTheme } from 'vuetify'

const theme = useTheme()
const isDark = computed(() => theme.global.current.value.dark)
const logoSrc = computed(() => isDark.value ? '/images/slogan_dark.png' : '/images/slogan_light.png')
</script>
