<template>
  <v-container class="fill-height not-found-bg" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-4 rounded-xl pa-8 glass-card text-center position-relative">
          <!-- Theme Toggle in Card Top Right -->
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

          <v-img
            src="/images/404.png"
            max-width="320"
            class="mx-auto mb-8 not-found-img"
            contain
          ></v-img>
          
          <p class="text-h5 mb-12 custom-text-color">{{ $t('notFound.description') }}</p>
          
          <v-btn
            color="primary"
            size="large"
            to="/"
            prepend-icon="mdi-home"
            class="text-none px-12 rounded-pill"
            elevation="2"
          >
            {{ $t('notFound.backHome') }}
          </v-btn>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.not-found-bg {
  background-position: center center !important;
  background-repeat: no-repeat !important;
  background-attachment: fixed !important;
  background-size: cover !important;
  transition: background-image 0.3s ease-in-out, background-color 0.3s ease-in-out;
}

.v-theme--light .not-found-bg {
  background-image: url('/images/slogan_light.png') !important;
  background-color: #f0f2f5 !important;
}

.v-theme--dark .not-found-bg {
  background-image: url('/images/slogan_dark.png') !important;
  background-color: #000000 !important;
}

.glass-card {
  backdrop-filter: blur(20px) saturate(160%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(160%) !important;
}

.v-theme--light .glass-card {
  background-color: rgba(255, 255, 255, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.4) !important;
}

.v-theme--dark .glass-card {
  background-color: rgba(30, 30, 30, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
}

.v-theme--light .custom-text-color {
  color: #000000 !important;
}

.v-theme--dark .custom-text-color {
  color: #ffffff !important;
}

.v-theme--dark .not-found-img {
  filter: invert(1) brightness(2);
}
</style>

<script setup>
import { computed } from 'vue'
import { useTheme } from 'vuetify'

const theme = useTheme()
const isDark = computed(() => theme.global.current.value.dark)

const toggleTheme = () => {
  const newTheme = isDark.value ? 'light' : 'dark'
  theme.global.name.value = newTheme
  localStorage.setItem('theme_manual', 'true')
  localStorage.setItem('theme', newTheme)
}
</script>
