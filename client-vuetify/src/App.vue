<template>
  <v-app>
    <router-view />
  </v-app>
</template>

<script setup>
import { watch, onMounted, onUnmounted } from 'vue'
import { useTheme } from 'vuetify'
import { useRoute } from 'vue-router'

const theme = useTheme()
const route = useRoute()

// Update body classes based on theme
const updateBodyStyles = () => {
  if (typeof document === 'undefined') return
  
  const isDark = theme.global.current.value.dark

  document.body.classList.remove('light-mode', 'dark-mode')
  document.body.classList.add(isDark ? 'dark-mode' : 'light-mode')
  
  localStorage.setItem('theme', isDark ? 'dark' : 'light')
}

watch(
  () => theme.global.current.value.dark,
  () => updateBodyStyles(),
  { immediate: true }
)

const handleSystemThemeChange = (e) => {
  // Always follow system theme change when it happens to ensure the feature is active
  const newTheme = e.matches ? 'dark' : 'light'
  theme.global.name.value = newTheme
  // Reset manual preference so it continues to follow system
  localStorage.removeItem('theme_manual')
  localStorage.setItem('theme', newTheme)
}

onMounted(() => {
  const systemDark = window.matchMedia('(prefers-color-scheme: dark)')
  
  // 1. Always listen for system theme changes
  systemDark.addEventListener('change', handleSystemThemeChange)

  // 2. Initial theme setup
  const savedTheme = localStorage.getItem('theme')
  const isManual = localStorage.getItem('theme_manual')

  if (isManual && savedTheme) {
    theme.global.name.value = savedTheme
  } else {
    theme.global.name.value = systemDark.matches ? 'dark' : 'light'
  }

  // Ensure initial body class
  const isDark = theme.global.current.value.dark
  document.body.classList.add(isDark ? 'dark-mode' : 'light-mode')
})

onUnmounted(() => {
  const systemDark = window.matchMedia('(prefers-color-scheme: dark)')
  systemDark.removeEventListener('change', handleSystemThemeChange)
})
</script>

<style>
/* Global Background Styles on Body */
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}

/* Global Background Styles on Body */
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}

body {
  background-position: center center !important;
  background-repeat: no-repeat !important;
  background-attachment: fixed !important;
  background-size: cover !important;
  transition: background-image 0.3s ease-in-out, background-color 0.3s ease-in-out;
}

body.light-mode {
  background-image: url('/images/slogan_light.png') !important;
  background-color: #f0f2f5 !important;
}

body.dark-mode {
  background-image: url('/images/slogan_dark.png') !important;
  background-color: #000000 !important;
}

/* Force Vuetify components to be transparent to show body background */
.v-application,
.v-application__wrap,
.v-layout,
.v-main {
  background: transparent !important;
}

/* Glassmorphism Effect for Cards, Sheets, and Dialogs */
.v-card, 
.v-sheet.rounded-lg.elevation-1,
.v-navigation-drawer,
.v-dialog > .v-card {
  backdrop-filter: blur(20px) saturate(160%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(160%) !important;
}

/* Light Theme Glass Effect */
body.light-mode .v-card,
body.light-mode .v-sheet.rounded-lg.elevation-1,
body.light-mode .v-navigation-drawer {
  background-color: rgba(255, 255, 255, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.4) !important;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.07) !important;
}

/* Dark Theme Glass Effect */
body.dark-mode .v-card,
body.dark-mode .v-sheet.rounded-lg.elevation-1,
body.dark-mode .v-navigation-drawer {
  background-color: rgba(30, 30, 30, 0.6) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.3) !important;
}

/* AppBar Glass Effect */
.v-app-bar {
  backdrop-filter: blur(12px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(12px) saturate(180%) !important;
}

body.light-mode .v-app-bar {
  background-color: rgba(255, 255, 255, 0.4) !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.3) !important;
}

body.dark-mode .v-app-bar {
  background-color: rgba(15, 15, 15, 0.4) !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1) !important;
}

/* Footer Glass Effect */
.v-footer {
  background-color: rgba(255, 255, 255, 0.1) !important;
  backdrop-filter: blur(10px) saturate(150%) !important;
  -webkit-backdrop-filter: blur(10px) saturate(150%) !important;
  border-top: 1px solid rgba(255, 255, 255, 0.1) !important;
}

body.light-mode .v-footer {
  background-color: rgba(255, 255, 255, 0.3) !important;
  color: #333 !important;
}

body.dark-mode .v-footer {
  background-color: rgba(0, 0, 0, 0.3) !important;
  color: #ccc !important;
}

/* Adjust text readability in glass containers */
.v-card-title, .v-card-text, .v-list-item-title {
  text-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

body.dark-mode .v-card-title, 
body.dark-mode .v-card-text {
  text-shadow: 0 1px 2px rgba(0,0,0,0.5);
}
</style>
