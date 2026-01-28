<template>
  <div v-if="image">
    <v-row>
      <v-col cols="12" md="8">
        <v-img
          :src="getImageUrl(image.path)"
          class="rounded-lg elevation-4"
          max-height="80vh"
        ></v-img>
      </v-col>
      <v-col cols="12" md="4">
        <v-card class="fill-height">
          <v-card-title>{{ $t('image.details') }}</v-card-title>
          <v-list density="compact">
            <v-list-item :title="$t('image.originalName')" :subtitle="image.original_name"></v-list-item>
            <v-list-item :title="$t('image.size')" :subtitle="formatSize(image.size)"></v-list-item>
            <v-list-item :title="$t('image.type')" :subtitle="image.mimetype"></v-list-item>
            <v-list-item :title="$t('image.album')" :subtitle="image.album_name || $t('home.noAlbum')"></v-list-item>
            <v-list-item :title="$t('image.uploadedAt')" :subtitle="formatDate(image.created_at)"></v-list-item>
          </v-list>
          
          <v-card-actions class="pa-4 flex-column align-stretch">
            <v-btn
              color="primary"
              variant="elevated"
              class="mb-2"
              prepend-icon="mdi-download"
              @click="handleDownload"
            >
              {{ $t('image.download') }}
            </v-btn>
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import api from '@/lib/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const image = ref(null)

const isAdmin = computed(() => userStore.user?.role === 'admin')

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

const fetchData = async () => {
  try {
    const res = await api.get(`/images/${route.params.id}`)
    image.value = res.data
  } catch (e) {}
}

const getImageUrl = (path) => {
  const base = (import.meta.env.VITE_API_URL || 'http://localhost:9108/api').replace('/api', '')
  return `${base}/${path}`
}

const formatSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (date) => new Date(date).toLocaleString()

const handleDownload = () => {
  const url = getImageUrl(image.value.path)
  const link = document.createElement('a')
  link.href = url
  link.download = image.value.original_name
  link.click()
}

const handleDelete = async () => {
  if (!confirm('Are you sure you want to delete this image?')) return
  try {
    await api.delete(`/images/${route.params.id}`)
    router.back()
  } catch (e) {
    showNotify(e.response?.data?.error || 'Failed to delete image', 'error')
  }
}

onMounted(fetchData)
</script>
