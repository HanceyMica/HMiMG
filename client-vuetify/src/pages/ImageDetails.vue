<template>
  <div v-if="image">
    <div class="d-flex align-center justify-space-between mb-4 flex-wrap ga-2">
      <v-btn variant="text" prepend-icon="mdi-arrow-left" @click="router.back()">
        {{ $t('image.back') }}
      </v-btn>
      <div class="d-flex ga-2">
        <v-btn
          variant="outlined"
          prepend-icon="mdi-chevron-left"
          :disabled="!previousImage"
          @click="goToImage(previousImage)"
        >
          {{ $t('image.previous') }}
        </v-btn>
        <v-btn
          variant="outlined"
          append-icon="mdi-chevron-right"
          :disabled="!nextImage"
          @click="goToImage(nextImage)"
        >
          {{ $t('image.next') }}
        </v-btn>
      </div>
    </div>
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

          <v-card-text v-if="canEditName" class="pt-0">
            <v-text-field
              v-model="renameForm.originalName"
              :label="$t('image.originalName')"
              variant="outlined"
              density="comfortable"
              hide-details
              class="mb-3"
            ></v-text-field>
            <v-btn
              color="primary"
              variant="outlined"
              block
              :loading="savingName"
              @click="handleRename"
            >
              {{ $t('image.update') }}
            </v-btn>
          </v-card-text>

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
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/store/user'
import api, { buildUploadedFileUrl } from '@/lib/api'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const userStore = useUserStore()
const image = ref(null)
const imageList = ref([])
const savingName = ref(false)

const isAdmin = computed(() => userStore.user?.role === 'admin')
const currentImageIndex = computed(() => imageList.value.findIndex(item => String(item.id) === String(route.params.id)))
const previousImage = computed(() => {
  if (currentImageIndex.value <= 0) return null
  return imageList.value[currentImageIndex.value - 1] || null
})
const nextImage = computed(() => {
  if (currentImageIndex.value < 0 || currentImageIndex.value >= imageList.value.length - 1) return null
  return imageList.value[currentImageIndex.value + 1] || null
})
const canEditName = computed(() => {
  if (!image.value || !userStore.user) return false
  return userStore.user.role === 'admin' || Number(image.value.uploaded_by) === Number(userStore.user.id)
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

const renameForm = reactive({
  originalName: ''
})

const fetchData = async () => {
  try {
    const res = await api.get(`/images/${route.params.id}`)
    image.value = res.data
    renameForm.originalName = res.data.original_name || ''
    const listRes = await api.get(`/images?albumId=${res.data.album_id}`)
    imageList.value = listRes.data
  } catch (e) {}
}

const getImageUrl = (path) => {
  return buildUploadedFileUrl(path)
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

const goToImage = (target) => {
  if (target?.id) {
    router.push(`/image/${target.id}`)
  }
}

const handleRename = async () => {
  if (!renameForm.originalName.trim()) {
    showNotify(t('image.nameRequired'), 'error')
    return
  }

  savingName.value = true
  try {
    const res = await api.put(`/images/${route.params.id}`, {
      original_name: renameForm.originalName.trim()
    })
    image.value = {
      ...image.value,
      ...res.data,
      album_name: image.value?.album_name
    }
    imageList.value = imageList.value.map(item => (
      String(item.id) === String(route.params.id)
        ? { ...item, original_name: res.data.original_name }
        : item
    ))
    showNotify(t('image.updateSuccess'))
  } catch (e) {
    showNotify(e.response?.data?.error || t('image.updateFailed'), 'error')
  } finally {
    savingName.value = false
  }
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
watch(() => route.params.id, fetchData)
</script>
