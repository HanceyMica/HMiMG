<template>
  <div class="text-center">
    <div class="hero-section py-8 py-sm-16">
      <p class="text-h4 text-sm-h3 font-weight-bold mb-4 mt-4 mt-sm-8">
        {{ $t('home.welcome', { title: settingsStore.websiteTitle }) }}
      </p>
      <v-row justify="center">
        <v-col cols="12" md="8" lg="6">
          <v-hover v-slot="{ isHovering, props }">
            <v-card
              v-bind="props"
              :elevation="isHovering ? 8 : 2"
              class="pa-8 mb-12 border-dashed rounded-xl cursor-pointer upload-card"
              @click="triggerFileInput"
            >
              <v-icon size="64" color="primary" class="mb-4">mdi-cloud-upload-outline</v-icon>
              <div class="text-h6 mb-2 upload-text">{{ $t('home.dragDropText') }}</div>
              <div class="text-body-2 text-grey-darken-1">{{ $t('home.dragDropHint') }}</div>
              
              <v-file-input
                ref="fileInput"
                v-model="files"
                multiple
                class="d-none"
                @update:model-value="onFilesSelected"
              ></v-file-input>
            </v-card>
          </v-hover>
        </v-col>
      </v-row>

      <div class="d-flex justify-center gap-4">
        <v-btn
          to="/library"
          color="primary"
          size="large"
          prepend-icon="mdi-image-multiple"
          class="text-none px-8 rounded-pill"
        >
          {{ $t('common.library') }}
        </v-btn>
        <v-btn
          v-if="isAdmin"
          to="/admin"
          variant="outlined"
          size="large"
          prepend-icon="mdi-cog"
          class="text-none px-8 rounded-pill"
        >
          {{ $t('common.admin') }}
        </v-btn>
      </div>
    </div>

    <!-- Upload Modal -->
    <v-dialog v-model="showUploadModal" max-width="500">
      <v-card class="rounded-lg">
        <v-card-title class="pa-4 bg-primary text-white">
          {{ $t('home.selectAlbumToUpload') }}
        </v-card-title>
        <v-card-text class="pa-6">
          <div class="mb-4">{{ files.length }} files selected</div>
          <v-select
            v-model="uploadForm.albumId"
            :items="albums"
            item-title="name"
            item-value="id"
            :label="$t('home.selectAlbum')"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          ></v-select>
          <v-btn
            color="primary"
            block
            size="large"
            :loading="uploading"
            :disabled="!uploadForm.albumId"
            @click="handleUpload"
          >
            {{ $t('home.upload') }}
          </v-btn>
        </v-card-text>
      </v-card>
    </v-dialog>

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
import { useRouter } from 'vue-router'
import api, { buildUploadedFileUrl } from '@/lib/api'
import { useUserStore } from '@/store/user'
import { useSettingsStore } from '@/store/settings'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const userStore = useUserStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()
const isAdmin = computed(() => userStore.user?.role === 'admin')

const albums = ref([])
const files = ref([])
const uploading = ref(false)
const showUploadModal = ref(false)
const fileInput = ref(null)

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

const uploadForm = reactive({
  albumId: null
})

const fetchData = async () => {
  try {
    const res = await api.get('/albums')
    albums.value = res.data
    // Ensure settings are fetched
    if (settingsStore.websiteTitle === 'HMiMG') {
      await settingsStore.fetchPublicSettings()
    }
  } catch (e) {
    console.error(e)
  }
}

const getImageUrl = (path) => {
  return buildUploadedFileUrl(path)
}

const triggerFileInput = () => {
  const input = document.querySelector('.v-file-input input')
  if (input) input.click()
}

const onFilesSelected = (selectedFiles) => {
  if (selectedFiles && selectedFiles.length > 0) {
    showUploadModal.value = true
  }
}

const handleUpload = async () => {
  if (files.value.length === 0) return
  uploading.value = true
  const formData = new FormData()
  formData.append('albumId', uploadForm.albumId)
  files.value.forEach(file => {
    formData.append('images', file)
  })

  try {
    const res = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const albumId = uploadForm.albumId
    files.value = []
    uploadForm.albumId = null
    showUploadModal.value = false
    window.dispatchEvent(new CustomEvent('hmimg:images-uploaded', {
      detail: {
        albumId,
        ids: res.data?.ids || []
      }
    }))
    showNotify(t('home.uploadSuccess'))
    fetchData()
    if (albumId) {
      router.push(`/album/${albumId}`)
    }
  } catch (e) {
    showNotify(t('home.uploadFailed'), 'error')
  } finally {
    uploading.value = false
  }
}

onMounted(fetchData)
</script>

<style scoped>
.border-dashed {
  border: 2px dashed rgba(var(--v-border-color), 0.3) !important;
}

.upload-card {
  transition: all 0.3s ease;
}

.v-theme--dark .upload-text {
  color: rgba(255, 255, 255, 0.9) !important;
}

.v-theme--light .upload-text {
  color: rgba(0, 0, 0, 0.87) !important;
}

.gap-4 {
  gap: 16px;
}
.transition-swing {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}
</style>
