<template>
  <div class="admin-page">
    <div class="d-flex align-center mb-8">
      <v-icon size="40" color="primary" class="mr-4">mdi-shield-check-outline</v-icon>
      <h1 class="text-h4 font-weight-bold">{{ $t('admin.dashboard') }}</h1>
    </div>
    
    <v-tabs v-model="activeTab" color="primary" class="mb-6 border-b">
      <v-tab value="system" prepend-icon="mdi-tune">{{ $t('admin.systemSettings') }}</v-tab>
      <v-tab value="organize" prepend-icon="mdi-folder-cog">{{ $t('admin.organize') }}</v-tab>
      <v-tab value="users" prepend-icon="mdi-account-group">{{ $t('admin.users') }}</v-tab>
      <v-tab value="account" prepend-icon="mdi-account-cog">{{ $t('admin.accountSettings') }}</v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <!-- System Settings -->
      <v-window-item value="system">
        <v-card border flat class="pa-6 rounded-lg">
          <v-form @submit.prevent="handleUpdateSettings">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="settings.website_title"
                  :label="$t('admin.websiteTitle')"
                  variant="outlined"
                  density="comfortable"
                  prepend-inner-icon="mdi-format-title"
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="settings.max_users"
                  :label="$t('admin.maxUsers')"
                  type="number"
                  variant="outlined"
                  density="comfortable"
                  prepend-inner-icon="mdi-account-group"
                ></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="settings.default_language"
                  :items="languages"
                  item-title="label"
                  item-value="value"
                  :label="$t('admin.defaultLanguage')"
                  variant="outlined"
                  density="comfortable"
                  prepend-inner-icon="mdi-translate"
                ></v-select>
              </v-col>
              <v-col cols="12">
                <v-switch
                  v-model="settings.allow_registration"
                  :label="$t('admin.allowRegistration')"
                  color="primary"
                  inset
                  hide-details
                ></v-switch>
              </v-col>
            </v-row>
            <v-btn type="submit" color="primary" class="mt-6 px-8" size="large" :loading="savingSettings">
              {{ $t('admin.saveSettings') }}
            </v-btn>
          </v-form>
        </v-card>
      </v-window-item>

      <!-- Organize (Albums/Collections) -->
      <v-window-item value="organize">
        <v-row>
          <v-col cols="12" md="6">
            <v-card border flat class="pa-6 rounded-lg h-100">
              <div class="text-h6 font-weight-bold mb-4 d-flex align-center">
                <v-icon start color="primary">mdi-image-plus</v-icon>
                {{ $t('admin.createAlbum') }}
              </div>
              <v-form @submit.prevent="handleCreateAlbum">
                <v-text-field v-model="newAlbum.name" :label="$t('admin.name')" variant="outlined" density="comfortable"></v-text-field>
                <v-textarea v-model="newAlbum.description" :label="$t('admin.description')" variant="outlined" rows="3"></v-textarea>
                <v-btn type="submit" color="primary" block size="large">{{ $t('admin.create') }}</v-btn>
              </v-form>
            </v-card>
          </v-col>
          <v-col cols="12" md="6">
            <v-card border flat class="pa-6 rounded-lg h-100">
              <div class="text-h6 font-weight-bold mb-4 d-flex align-center">
                <v-icon start color="secondary">mdi-folder-plus</v-icon>
                {{ $t('admin.createCollection') }}
              </div>
              <v-form @submit.prevent="handleCreateCollection">
                <v-text-field v-model="newCollection.name" :label="$t('admin.name')" variant="outlined" density="comfortable"></v-text-field>
                <v-textarea v-model="newCollection.description" :label="$t('admin.description')" variant="outlined" rows="3"></v-textarea>
                <v-btn type="submit" color="secondary" block size="large">{{ $t('admin.create') }}</v-btn>
              </v-form>
            </v-card>
          </v-col>
          <v-col cols="12">
            <v-card border flat class="pa-6 rounded-lg">
              <div class="text-h6 font-weight-bold mb-4 d-flex align-center">
                <v-icon start color="success">mdi-file-tree</v-icon>
                {{ $t('admin.organize') }}
              </div>
              <v-form @submit.prevent="handleAddToCollection">
                <v-row>
                  <v-col cols="12" md="4">
                    <v-select
                      v-model="organizeForm.collectionId"
                      :items="collections"
                      item-title="name"
                      item-value="id"
                      :label="$t('admin.targetCollection')"
                      variant="outlined"
                      density="comfortable"
                    ></v-select>
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-radio-group v-model="organizeForm.itemType" inline label="Type">
                      <v-radio :label="$t('admin.album')" value="album"></v-radio>
                      <v-radio :label="$t('admin.collection')" value="collection"></v-radio>
                    </v-radio-group>
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-select
                      v-model="organizeForm.itemId"
                      :items="organizeForm.itemType === 'album' ? albums : collections"
                      item-title="name"
                      item-value="id"
                      :label="$t('admin.itemName')"
                      variant="outlined"
                      density="comfortable"
                    ></v-select>
                  </v-col>
                </v-row>
                <v-btn type="submit" color="success" class="px-8" size="large">{{ $t('admin.add') }}</v-btn>
              </v-form>
            </v-card>
          </v-col>
        </v-row>
      </v-window-item>

      <!-- Users Management -->
      <v-window-item value="users">
        <v-card border flat class="rounded-lg">
          <v-data-table
            :headers="userTableHeaders"
            :items="users"
            :loading="loadingUsers"
            class="rounded-lg"
          >
            <template v-slot:item.role="{ item }">
              <v-chip
                size="small"
                :color="item.role === 'admin' ? 'primary' : 'default'"
                variant="tonal"
              >
                {{ item.role === 'admin' ? $t('admin.roleAdmin') : $t('admin.roleUser') }}
              </v-chip>
            </template>
            <template v-slot:item.actions="{ item }">
              <div class="d-flex ga-2">
                <v-btn
                  size="small"
                  variant="outlined"
                  :disabled="item.id === userStore.user?.id"
                  @click="handleToggleRole(item)"
                >
                  {{ item.role === 'admin' ? $t('admin.roleUser') : $t('admin.roleAdmin') }}
                </v-btn>
                <v-btn
                  size="small"
                  color="error"
                  variant="outlined"
                  :disabled="item.id === userStore.user?.id"
                  @click="handleDeleteUser(item)"
                >
                  {{ $t('admin.deleteUser') }}
                </v-btn>
              </div>
            </template>
          </v-data-table>
        </v-card>
      </v-window-item>

      <!-- Account Settings -->
      <v-window-item value="account">
        <v-card border flat class="pa-6 rounded-lg max-w-800 mx-auto">
          <v-form @submit.prevent="handleUpdateProfile">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field v-model="profile.username" :label="$t('login.username')" variant="outlined" density="comfortable" prepend-inner-icon="mdi-account"></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="profile.email" :label="$t('login.email')" variant="outlined" density="comfortable" prepend-inner-icon="mdi-email"></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field v-model="profile.phone" :label="$t('login.phone')" variant="outlined" density="comfortable" prepend-inner-icon="mdi-phone"></v-text-field>
              </v-col>
            </v-row>
            
            <v-divider class="my-6"></v-divider>
            <div class="text-subtitle-1 font-weight-bold mb-4">{{ $t('admin.passwordHelp') }}</div>
            
            <v-row>
              <v-col cols="12">
                <v-text-field v-model="profile.oldPassword" :label="$t('admin.oldPassword')" type="password" variant="outlined" density="comfortable" prepend-inner-icon="mdi-lock-reset"></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="profile.password" :label="$t('admin.newPassword')" type="password" variant="outlined" density="comfortable" prepend-inner-icon="mdi-lock-plus"></v-text-field>
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field v-model="profile.confirmPassword" :label="$t('admin.confirmPassword')" type="password" variant="outlined" density="comfortable" prepend-inner-icon="mdi-lock-check"></v-text-field>
              </v-col>
            </v-row>
            <v-btn type="submit" color="primary" class="mt-6 px-8" size="large">{{ $t('admin.updateProfile') }}</v-btn>
          </v-form>
        </v-card>
      </v-window-item>
    </v-window>

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
import { ref, reactive, computed, onMounted } from 'vue'
import api from '@/lib/api'
import { useUserStore } from '@/store/user'
import { useSettingsStore } from '@/store/settings'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const userStore = useUserStore()
const settingsStore = useSettingsStore()
const router = useRouter()
const { t } = useI18n()

const activeTab = ref('system')
const savingSettings = ref(false)
const albums = ref([])
const collections = ref([])
const languages = [
  { label: 'English', value: 'en' },
  { label: '简体中文', value: 'zh' },
  { label: '日本語', value: 'ja' }
]

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

const settings = reactive({
  website_title: '',
  max_users: 100,
  allow_registration: true,
  default_language: 'zh'
})

const newAlbum = reactive({ name: '', description: '' })
const newCollection = reactive({ name: '', description: '' })
const organizeForm = reactive({ collectionId: null, itemType: 'album', itemId: null })

const users = ref([])
const loadingUsers = ref(false)
const userTableHeaders = computed(() => [
  { title: t('admin.username'), key: 'username' },
  { title: t('admin.email'), key: 'email' },
  { title: t('admin.phone'), key: 'phone' },
  { title: t('admin.role'), key: 'role' },
  { title: t('admin.createdAt'), key: 'created_at' },
  { title: t('admin.actions'), key: 'actions', sortable: false }
])

const fetchUsers = async () => {
  loadingUsers.value = true
  try {
    const res = await api.get('/admin/users')
    users.value = res.data
  } catch (e) {
    showNotify(e.response?.data?.error || t('common.loadFailed'), 'error')
  } finally {
    loadingUsers.value = false
  }
}

const handleToggleRole = async (user) => {
  const newRole = user.role === 'admin' ? 'user' : 'admin'
  try {
    await api.put(`/admin/users/${user.id}/role`, { role: newRole })
    showNotify(t('admin.roleUpdated'))
    fetchUsers()
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.roleUpdateFailed'), 'error')
  }
}

const handleDeleteUser = async (user) => {
  if (!confirm(t('admin.confirmDeleteUser'))) return
  try {
    await api.delete(`/admin/users/${user.id}`)
    showNotify(t('admin.userDeleted'))
    fetchUsers()
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.userDeleteFailed'), 'error')
  }
}

const profile = reactive({
  username: userStore.user?.username || '',
  email: userStore.user?.email || '',
  phone: userStore.user?.phone || '',
  oldPassword: '',
  password: '',
  confirmPassword: ''
})

const fetchData = async () => {
  try {
    const [settRes, albRes, colRes] = await Promise.all([
      api.get('/settings'),
      api.get('/albums'),
      api.get('/collections')
    ])
    Object.assign(settings, {
      ...settRes.data,
      allow_registration: settRes.data.allow_registration === 'true'
    })
    albums.value = albRes.data
    collections.value = colRes.data
  } catch (e) {
    showNotify(e.response?.data?.error || t('common.loadFailed'), 'error')
  }
}

onMounted(() => {
  fetchData()
  fetchUsers()
})

const handleUpdateSettings = async () => {
  savingSettings.value = true
  try {
    await api.put('/settings', {
      ...settings,
      allow_registration: String(settings.allow_registration)
    })
    await settingsStore.fetchPublicSettings()
    showNotify(t('admin.settingsUpdated'))
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.settingsFailed'), 'error')
  } finally {
    savingSettings.value = false
  }
}

const handleCreateAlbum = async () => {
  try {
    await api.post('/albums', newAlbum)
    newAlbum.name = ''
    newAlbum.description = ''
    showNotify(t('admin.albumCreated'))
    fetchData()
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.failedCreateAlbum'), 'error')
  }
}

const handleCreateCollection = async () => {
  try {
    await api.post('/collections', newCollection)
    newCollection.name = ''
    newCollection.description = ''
    showNotify(t('admin.collectionCreated'))
    fetchData()
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.failedCreateCollection'), 'error')
  }
}

const handleAddToCollection = async () => {
  try {
    await api.post('/collections/add', {
      collectionId: organizeForm.collectionId,
      itemType: organizeForm.itemType,
      itemId: organizeForm.itemId
    })
    showNotify(t('admin.addedSuccess'))
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.failedAdd'), 'error')
  }
}

const handleUpdateProfile = async () => {
  if (profile.password && profile.password !== profile.confirmPassword) {
    showNotify(t('admin.passwordMismatch'), 'error')
    return
  }
  try {
    const res = await api.put('/admin/update', profile)
    if (res.data.passwordChanged) {
      showNotify(t('admin.changePasswordSuccess'))
      setTimeout(() => {
        userStore.logout()
        router.push('/login')
      }, 2000)
    } else {
      showNotify(t('admin.profileUpdated'))
    }
  } catch (e) {
    showNotify(e.response?.data?.error || t('admin.profileFailed'), 'error')
  }
}
</script>

<style scoped>
.max-w-800 {
  max-width: 800px;
}
</style>
