<template>
  <div v-if="album">
    <div class="d-flex align-center mb-4">
      <v-btn icon="mdi-arrow-left" variant="text" @click="$router.back()"></v-btn>
      <h1 class="ml-2">{{ album.name }}</h1>
      <v-btn v-if="isAdmin" icon="mdi-cog" variant="text" class="ml-2" @click="showEditDialog = true"></v-btn>
    </div>
    <p class="mb-6 text-body-1">{{ album.description }}</p>

    <!-- Edit Dialog -->
    <v-dialog v-model="showEditDialog" max-width="500">
      <v-card class="rounded-lg">
        <v-card-title class="pa-4 bg-primary text-white">
          {{ $t('common.edit') }} {{ $t('admin.album') }}
        </v-card-title>
        <v-card-text class="pa-6">
          <v-text-field v-model="editForm.name" :label="$t('admin.name')" variant="outlined" density="comfortable" class="mb-4"></v-text-field>
          <v-textarea v-model="editForm.description" :label="$t('admin.description')" variant="outlined" rows="3" class="mb-4"></v-textarea>
          <v-btn color="primary" block size="large" :loading="saving" @click="handleUpdate">
            {{ $t('image.update') }}
          </v-btn>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-row>
      <v-col
        v-for="img in images"
        :key="img.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card :to="`/image/${img.id}`">
          <v-img
            :src="getImageUrl(img.path)"
            cover
            height="200"
          ></v-img>
          <v-card-subtitle class="pt-2">{{ img.original_name }}</v-card-subtitle>
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
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useI18n } from 'vue-i18n'
import api from '@/lib/api'

const route = useRoute()
const userStore = useUserStore()
const { t } = useI18n()
const isAdmin = computed(() => userStore.user?.role === 'admin')

const album = ref(null)
const images = ref([])
const showEditDialog = ref(false)
const saving = ref(false)

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

const editForm = reactive({
  name: '',
  description: ''
})

const fetchData = async () => {
  try {
    const id = route.params.id
    const [albRes, imgRes] = await Promise.all([
      api.get(`/albums/${id}`),
      api.get(`/images?albumId=${id}`)
    ])
    album.value = albRes.data
    editForm.name = albRes.data.name
    editForm.description = albRes.data.description
    images.value = imgRes.data
  } catch (e) {}
}

const handleUpdate = async () => {
  saving.value = true
  try {
    await api.put(`/albums/${route.params.id}`, editForm)
    showEditDialog.value = false
    showNotify(t('image.updateSuccess'))
    fetchData()
  } catch (e) {
    showNotify(e.response?.data?.error || t('image.updateFailed'), 'error')
  } finally {
    saving.value = false
  }
}

const getImageUrl = (path) => {
  const base = (import.meta.env.VITE_API_URL || 'http://localhost:9108/api').replace('/api', '')
  return `${base}/${path}`
}

onMounted(fetchData)
</script>
