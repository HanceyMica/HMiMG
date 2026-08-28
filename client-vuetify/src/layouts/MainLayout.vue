<template>
  <v-layout>
    <!-- Navigation Drawer for Mobile -->
    <v-navigation-drawer v-model="drawer" temporary class="glass-drawer">
      <v-list nav>
        <v-list-item prepend-icon="mdi-home" to="/" :title="$t('common.home')" @click="drawer = false"></v-list-item>
        <v-list-item prepend-icon="mdi-image-multiple" to="/library" :title="$t('common.library')" @click="drawer = false"></v-list-item>
        <v-list-item v-if="isAdmin" prepend-icon="mdi-shield-check" to="/admin" :title="$t('common.admin')" @click="drawer = false"></v-list-item>
        <v-list-item prepend-icon="mdi-information" to="/about" :title="$t('common.about')" @click="drawer = false"></v-list-item>
      </v-list>
      <v-divider></v-divider>
      <v-list nav>
        <v-list-item
          :prepend-icon="isDark ? 'mdi-weather-night' : 'mdi-weather-sunny'"
          :title="isDark ? $t('common.dark') : $t('common.light')"
          @click="toggleTheme"
        ></v-list-item>
        <v-list-item v-if="user" prepend-icon="mdi-logout" :title="$t('common.logout')" @click="handleLogout"></v-list-item>
        <v-list-item v-else prepend-icon="mdi-login" to="/login" :title="$t('login.loginBtn')" @click="drawer = false"></v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar flat border height="64">
      <v-container class="d-flex align-center pa-0 fill-height">
        <!-- Mobile Menu Toggle -->
        <v-app-bar-nav-icon class="d-md-none mr-2" @click="drawer = !drawer"></v-app-bar-nav-icon>

        <!-- Favicon / Logo -->
        <div class="mr-4">
          <router-link to="/" class="d-flex align-center text-decoration-none">
            <v-avatar size="32" class="mr-2 rounded-lg overflow-hidden border">
              <img src="/images/favicon.png" alt="HMiMG" style="width: 100%; height: 100%; object-fit: cover;">
            </v-avatar>
            <span 
              class="text-h6 font-weight-bold d-none d-sm-flex site-title"
            >
              {{ settingsStore.websiteTitle }}
            </span>
          </router-link>
        </div>

        <!-- Desktop Navigation -->
        <div class="d-none d-md-flex align-center">
          <v-btn to="/" variant="text" class="text-none mx-1" :active="route.path === '/'">{{ $t('common.home') }}</v-btn>
          <v-btn to="/library" variant="text" class="text-none mx-1" :active="route.path.startsWith('/library') || route.path.startsWith('/album') || route.path.startsWith('/collection') || route.path.startsWith('/image')">{{ $t('common.library') }}</v-btn>
          <v-btn v-if="isAdmin" to="/admin" variant="text" class="text-none mx-1" :active="route.path.startsWith('/admin')">{{ $t('common.admin') }}</v-btn>
          <v-btn to="/about" variant="text" class="text-none mx-1" :active="route.path.startsWith('/about')">{{ $t('common.about') }}</v-btn>
        </div>

        <v-spacer></v-spacer>

        <!-- Right Side Actions (Desktop) -->
        <div class="d-none d-md-flex align-center">
          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn
                v-bind="props"
                icon
                variant="text"
                class="mr-2"
                :title="$t('common.language')"
              >
                <v-icon>mdi-translate</v-icon>
              </v-btn>
            </template>
            <v-list density="compact">
              <v-list-item
                v-for="lang in languages"
                :key="lang.value"
                :active="locale === lang.value"
                :title="lang.label"
                @click="switchLanguage(lang.value)"
              ></v-list-item>
            </v-list>
          </v-menu>

          <v-btn
            icon
            variant="text"
            class="mr-2"
            :color="isDark ? 'yellow' : 'primary'"
            @click="toggleTheme"
          >
            <v-icon>{{ isDark ? 'mdi-weather-night' : 'mdi-weather-sunny' }}</v-icon>
          </v-btn>

          <v-btn
            v-if="user"
            variant="flat"
            size="small"
            class="rounded-lg px-4 logout-btn font-weight-bold"
            prepend-icon="mdi-logout-variant"
            @click="handleLogout"
          >
            {{ $t('common.logout') }}
          </v-btn>
          <v-btn v-else to="/login" color="primary" variant="flat" size="small" class="rounded-pill px-4">
            {{ $t('login.loginBtn') }}
          </v-btn>
        </div>

        <!-- Right Side Actions (Mobile) -->
        <div class="d-md-none">
          <v-menu>
            <template v-slot:activator="{ props }">
              <v-btn v-bind="props" icon="mdi-translate" variant="text" size="small"></v-btn>
            </template>
            <v-list density="compact">
              <v-list-item
                v-for="lang in languages"
                :key="lang.value"
                :active="locale === lang.value"
                :title="lang.label"
                @click="switchLanguage(lang.value)"
              ></v-list-item>
            </v-list>
          </v-menu>
          <v-btn v-if="!user" to="/login" icon="mdi-login" variant="text" size="small"></v-btn>
        </div>
      </v-container>
    </v-app-bar>

    <v-main>
      <v-container class="py-4 py-sm-8">
        <v-sheet class="pa-4 pa-sm-6 rounded-lg elevation-1 min-h-main">
          <router-view />
        </v-sheet>
      </v-container>
    </v-main>

    <v-footer app border height="60">
      <div class="text-center w-100 text-grey text-caption text-sm-body-2">
        {{ $t('common.copyright', { title: settingsStore.websiteTitle }) }}
      </div>
    </v-footer>
  </v-layout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useUserStore } from '@/store/user'
import { useSettingsStore } from '@/store/settings'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useTheme } from 'vuetify'
import api from '@/lib/api'

const userStore = useUserStore()
const settingsStore = useSettingsStore()
const router = useRouter()
const route = useRoute()
const { locale, t } = useI18n()
const theme = useTheme()

const drawer = ref(false)
const user = computed(() => userStore.user)
const isAdmin = computed(() => userStore.user?.role === 'admin')
const isDark = computed(() => theme.global.current.value.dark)

const languages = [
  { label: 'English', value: 'en' },
  { label: '简体中文', value: 'zh' },
  { label: '日本語', value: 'ja' }
]

const switchLanguage = (lang) => {
  locale.value = lang
  localStorage.setItem('lang', lang)
}

// Dynamic Page Title
const pageTitle = computed(() => {
  const websiteTitle = settingsStore.websiteTitle
  const path = route.path
  
  if (path === '/') return websiteTitle
  
  let pageName = ''
  if (path === '/library') pageName = t('common.library')
  else if (path.startsWith('/album/')) pageName = t('common.album')
  else if (path.startsWith('/collection/')) pageName = t('common.collection')
  else if (path === '/admin') pageName = t('common.admin')
  else if (path === '/about') pageName = t('common.about')
  else if (path === '/login') pageName = t('login.title')
  else if (path.startsWith('/image/')) pageName = t('image.details')
  else pageName = t('notFound.title')

  return `${pageName} - ${websiteTitle}`
})

const fetchData = async () => {
  try {
    await settingsStore.fetchPublicSettings()
    
    // Set initial locale from backend if available and not set in localstorage
    if (settingsStore.defaultLanguage && !localStorage.getItem('lang')) {
      locale.value = settingsStore.defaultLanguage
    }
  } catch (e) {}
}

watch(pageTitle, (newTitle) => {
  document.title = newTitle
}, { immediate: true })

const handleLogout = () => {
  userStore.logout()
  drawer.value = false
  router.push('/login')
}

const toggleTheme = () => {
  const newTheme = isDark.value ? 'light' : 'dark'
  theme.global.name.value = newTheme
  localStorage.setItem('theme_manual', 'true')
  localStorage.setItem('theme', newTheme)
}

onMounted(fetchData)
</script>

<style scoped>
.min-h-main {
  min-height: 400px;
}
.glass-drawer {
  backdrop-filter: blur(20px) saturate(160%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(160%) !important;
  background-color: rgba(var(--v-theme-surface), 0.7) !important;
}

.site-title {
  transition: color 0.3s ease;
}

.v-theme--light .site-title {
  color: #000000 !important;
}

.v-theme--dark .site-title {
  color: #ffffff !important;
}

.logout-btn {
  position: relative;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  text-transform: none;
  letter-spacing: 0.5px;
}

/* Noise effect overlay */
.logout-btn::before {
  content: "";
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
  opacity: 0.15;
  pointer-events: none;
  mix-blend-mode: overlay;
}

/* Light Mode: Black with Shadow */
.v-theme--light .logout-btn {
  background-color: #1a1a1a !important;
  color: #ffffff !important;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3), inset 0 1px 1px rgba(255, 255, 255, 0.1) !important;
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
}

.v-theme--light .logout-btn:hover {
  background-color: #000000 !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.4) !important;
}

/* Dark Mode: White with Glow */
.v-theme--dark .logout-btn {
  background-color: #ffffff !important;
  color: #000000 !important;
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.3), 0 0 2px rgba(255, 255, 255, 0.5) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
}

.v-theme--dark .logout-btn:hover {
  background-color: #f0f0f0 !important;
  transform: translateY(-2px);
  box-shadow: 0 0 20px rgba(255, 255, 255, 0.6), 0 0 10px rgba(255, 255, 255, 0.4) !important;
}
</style>
