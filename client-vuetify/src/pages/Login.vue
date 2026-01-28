<template>
  <v-container class="fill-height login-bg" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4" lg="3">
        <v-card class="elevation-4 rounded-xl pa-4 glass-card position-relative">
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

          <v-card-item class="text-center py-8">
            <v-icon size="64" color="primary" class="mb-4">mdi-image-multiple</v-icon>
            <v-card-title class="text-h4 font-weight-bold">
              {{ settingsStore.websiteTitle }}
            </v-card-title>
            <v-card-subtitle class="text-subtitle-1 mt-2">
              {{ isRegister ? $t('login.registerBtn') : $t('login.title') }}
            </v-card-subtitle>
          </v-card-item>

          <v-card-text>
            <v-form @submit.prevent="handleSubmit" class="mt-4">
              <v-text-field
                v-model="form.username"
                :label="$t('login.username')"
                prepend-inner-icon="mdi-account-outline"
                variant="outlined"
                density="comfortable"
                required
                class="mb-2"
              ></v-text-field>
              
              <v-text-field
                v-model="form.password"
                :label="$t('login.password')"
                prepend-inner-icon="mdi-lock-outline"
                type="password"
                variant="outlined"
                density="comfortable"
                required
                class="mb-2"
              ></v-text-field>

              <v-expand-transition>
                <div v-if="isRegister">
                  <v-text-field
                    v-model="form.email"
                    :label="$t('login.email')"
                    prepend-inner-icon="mdi-email-outline"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-text-field>
                  <v-text-field
                    v-model="form.phone"
                    :label="$t('login.phone')"
                    prepend-inner-icon="mdi-phone-outline"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-text-field>
                </div>
              </v-expand-transition>

              <v-alert
                v-if="error"
                type="error"
                variant="tonal"
                density="compact"
                class="mb-4 rounded-lg"
                closable
                @click:close="error = ''"
              >
                {{ error }}
              </v-alert>

              <v-btn
                type="submit"
                color="primary"
                block
                size="large"
                class="rounded-lg text-none font-weight-bold mt-4"
                :loading="loading"
                elevation="2"
              >
                {{ isRegister ? $t('login.registerBtn') : $t('login.loginBtn') }}
              </v-btn>
            </v-form>
          </v-card-text>

          <v-card-actions class="justify-center pb-6" v-if="settingsStore.allowRegistration">
            <v-btn
              variant="text"
              color="primary"
              class="text-none"
              @click="isRegister = !isRegister"
            >
              {{ isRegister ? $t('login.toLogin') : $t('login.toRegister') }}
            </v-btn>
          </v-card-actions>

          <v-divider class="mx-4"></v-divider>
          
          <div class="text-center py-4 text-grey text-caption">
            {{ settingsStore.websiteTitle }} © 2026 HanceyMica
          </div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Global Snackbar for Notifications -->
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
  </v-container>
</template>

<style scoped>
.login-bg {
  background-position: center center !important;
  background-repeat: no-repeat !important;
  background-attachment: fixed !important;
  background-size: cover !important;
  transition: background-image 0.3s ease-in-out, background-color 0.3s ease-in-out;
}

.v-theme--light .login-bg {
  background-image: url('/images/slogan_light.png') !important;
  background-color: #f0f2f5 !important;
}

.v-theme--dark .login-bg {
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
</style>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useSettingsStore } from '@/store/settings'
import { useI18n } from 'vue-i18n'
import { useTheme } from 'vuetify'
import api from '@/lib/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()
const theme = useTheme()

const isRegister = ref(false)
const loading = ref(false)
const error = ref('')

const isDark = computed(() => theme.global.current.value.dark)

const toggleTheme = () => {
  const newTheme = isDark.value ? 'light' : 'dark'
  theme.global.name.value = newTheme
  localStorage.setItem('theme_manual', 'true')
  localStorage.setItem('theme', newTheme)
}

onMounted(async () => {
  await settingsStore.fetchPublicSettings()
  
  if (route.query.reason === 'unauthorized') {
    showNotify(t('common.notLoggedIn'), 'error')
    // Clear query to prevent re-triggering on refresh
    router.replace({ query: {} })
  }
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
  icon: 'mdi-check-circle'
})

const showNotify = (text, color = 'success') => {
  snackbar.text = text
  snackbar.color = color
  snackbar.icon = color === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'
  snackbar.show = true
}

const form = reactive({
  username: '',
  password: '',
  email: '',
  phone: ''
})

const handleSubmit = async () => {
  if (loading.value) return
  loading.value = true
  error.value = ''
  try {
    if (isRegister.value) {
      await api.post('/register', form)
      isRegister.value = false
      showNotify(t('login.registerSuccess'))
    } else {
      const res = await api.post('/login', {
        username: form.username,
        password: form.password
      })
      userStore.setUser(res.data.user, res.data.token)
      router.push('/')
    }
  } catch (err) {
    error.value = err.response?.data?.error || t('login.failed')
  } finally {
    loading.value = false
  }
}
</script>
