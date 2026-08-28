<template>
  <div v-if="collection">
    <div class="d-flex align-center mb-4">
      <v-btn icon="mdi-arrow-left" variant="text" @click="$router.back()"></v-btn>
      <h1 class="ml-2">{{ collection.name }}</h1>
      <v-btn v-if="isAdmin" icon="mdi-cog" variant="text" class="ml-2" @click="showEditDialog = true"></v-btn>
    </div>
    <p class="mb-6 text-body-1">{{ collection.description }}</p>

    <!-- Edit Dialog -->
    <v-dialog v-model="showEditDialog" max-width="500">
      <v-card class="rounded-lg">
        <v-card-title class="pa-4 bg-primary text-white">
          {{ $t('common.edit') }} {{ $t('admin.collection') }}
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
        v-for="item in collection.children"
        :key="item.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card :to="item.type === 'album' ? `/album/${item.id}` : `/collection/${item.id}`">
          <v-img
            v-if="item.cover_image"
            :src="getImageUrl(item.cover_image)"
            height="150"
            cover
          ></v-img>
          <v-sheet v-else height="150" :color="item.type === 'album' ? 'grey-lighten-2' : 'surface-variant'" class="d-flex align-center justify-center">
            <v-icon size="48" color="grey">{{ item.type === 'album' ? 'mdi-image' : 'mdi-folder' }}</v-icon>
          </v-sheet>
          <v-card-title>
             <v-icon v-if="item.type === 'collection'" start size="small">mdi-folder-outline</v-icon>
             {{ item.name }}
          </v-card-title>
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
        <v-btn variant="text" @click="snackbar.show = false">{{ $t('common.close') }}</v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useI18n } from 'vue-i18n'
import api, { buildUploadedFileUrl } from '@/lib/api'

const route = useRoute()
const userStore = useUserStore()
const { t } = useI18n()
const isAdmin = computed(() => userStore.user?.role === 'admin')

const collection = ref(null)
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
    const res = await api.get(`/collections/${route.params.id}`)
    collection.value = res.data
    editForm.name = res.data.name
    editForm.description = res.data.description
  } catch (e) {
    showNotify(e.response?.data?.error || t('common.loadFailed'), 'error')
  }
}

const handleUpdate = async () => {
  saving.value = true
  try {
    await api.put(`/collections/${route.params.id}`, editForm)
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
  return buildUploadedFileUrl(path)
}

onMounted(fetchData)
watch(() => route.params.id, fetchData)
</script>
