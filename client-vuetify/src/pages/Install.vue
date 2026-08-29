<template>
  <v-container class="fill-height install-bg" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="10" md="8" lg="6">
        <v-card class="elevation-4 rounded-xl pa-4 glass-card">
          <v-card-item class="text-center py-6">
            <v-icon size="56" color="primary" class="mb-3">mdi-image-multiple</v-icon>
            <v-card-title class="text-h4 font-weight-bold">
              {{ $t('install.title') }}
            </v-card-title>
            <v-card-subtitle class="text-subtitle-1 mt-2">
              {{ $t('install.subtitle') }}
            </v-card-subtitle>
          </v-card-item>

          <v-stepper v-model="step" :items="stepItems" class="rounded-lg bg-transparent">
            <!-- Step 1: Welcome / Environment -->
            <template v-slot:item.1>
              <v-card flat>
                <v-list density="compact">
                  <v-list-item :title="$t('install.env.version')" :subtitle="installStore.version || '-'">
                    <template v-slot:append>
                      <v-icon color="success">mdi-check-circle</v-icon>
                    </template>
                  </v-list-item>
                  <v-list-item :title="$t('install.env.upload')">
                    <template v-slot:append>
                      <v-icon :color="installStore.uploadWritable ? 'success' : 'error'">
                        {{ installStore.uploadWritable ? 'mdi-check-circle' : 'mdi-close-circle' }}
                      </v-icon>
                    </template>
                  </v-list-item>
                  <v-list-item :title="$t('install.env.db')" :subtitle="installStore.dbError || undefined">
                    <template v-slot:append>
                      <v-icon :color="installStore.hasDb ? 'success' : 'warning'">
                        {{ installStore.hasDb ? 'mdi-check-circle' : 'mdi-help-circle' }}
                      </v-icon>
                    </template>
                  </v-list-item>
                </v-list>
                <v-alert v-if="installStore.dbError" type="warning" variant="tonal" density="compact" class="mt-2 mb-4 rounded-lg">
                  {{ installStore.dbError }}
                </v-alert>
                <v-alert type="info" variant="tonal" density="compact" class="mb-4 rounded-lg">
                  {{ $t('install.env.httpsHint') }}
                </v-alert>
                <v-btn color="primary" block size="large" @click="goNext">
                  {{ $t('install.common.next') }}
                </v-btn>
              </v-card>
            </template>

            <!-- Step 2: Database -->
            <template v-slot:item.2>
              <v-card flat>
                <v-alert v-if="installStore.hasDb" type="success" variant="tonal" density="compact" class="mb-4 rounded-lg">
                  {{ $t('install.db.alreadyConfigured') }}
                </v-alert>
                <v-form @submit.prevent="submitDatabase">
                  <v-select
                    v-model="dbForm.driver"
                    :items="dbDrivers"
                    item-title="label"
                    item-value="value"
                    :label="$t('install.db.driver')"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-select>
                  <v-row dense>
                    <v-col cols="8">
                      <v-text-field
                        v-model="dbForm.host"
                        :label="$t('install.db.host')"
                        variant="outlined"
                        density="comfortable"
                        class="mb-2"
                        required
                      ></v-text-field>
                    </v-col>
                    <v-col cols="4">
                      <v-text-field
                        v-model.number="dbForm.port"
                        :label="$t('install.db.port')"
                        type="number"
                        variant="outlined"
                        density="comfortable"
                        class="mb-2"
                        required
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  <v-text-field
                    v-model="dbForm.user"
                    :label="$t('install.db.user')"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                    required
                  ></v-text-field>
                  <v-text-field
                    v-model="dbForm.password"
                    :label="$t('install.db.password')"
                    type="password"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-text-field>
                  <v-text-field
                    v-model="dbForm.name"
                    :label="$t('install.db.name')"
                    variant="outlined"
                    density="comfortable"
                    class="mb-4"
                    required
                  ></v-text-field>

                  <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mb-4 rounded-lg" closable @click:close="error = ''">
                    {{ error }}
                  </v-alert>
                  <v-alert v-if="dbTestOk" type="success" variant="tonal" density="compact" class="mb-4 rounded-lg">
                    {{ $t('install.db.testOk') }}
                  </v-alert>

                  <div class="d-flex ga-2">
                    <v-btn variant="outlined" size="large" :loading="testing" :disabled="dbReady" @click="testDatabase">
                      {{ $t('install.db.test') }}
                    </v-btn>
                    <v-btn type="submit" color="primary" size="large" class="flex-grow-1" :loading="saving" :disabled="dbReady">
                      {{ $t('install.db.submit') }}
                    </v-btn>
                  </div>
                </v-form>
              </v-card>
            </template>

            <!-- Step 3: Admin -->
            <template v-slot:item.3>
              <v-card flat>
                <v-form @submit.prevent="submitAdmin">
                  <v-text-field
                    v-model="adminForm.username"
                    :label="$t('install.admin.username')"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                    required
                  ></v-text-field>
                  <v-text-field
                    v-model="adminForm.password"
                    :label="$t('install.admin.password')"
                    type="password"
                    variant="outlined"
                    density="comfortable"
                    :hint="$t('install.admin.passwordHint')"
                    persistent-hint
                    class="mb-2"
                    required
                  ></v-text-field>
                  <v-text-field
                    v-model="adminForm.confirmPassword"
                    :label="$t('install.admin.confirmPassword')"
                    type="password"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                    required
                  ></v-text-field>
                  <v-text-field
                    v-model="adminForm.email"
                    :label="$t('install.admin.email')"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-text-field>
                  <v-text-field
                    v-model="adminForm.phone"
                    :label="$t('install.admin.phone')"
                    variant="outlined"
                    density="comfortable"
                    class="mb-4"
                  ></v-text-field>

                  <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mb-4 rounded-lg" closable @click:close="error = ''">
                    {{ error }}
                  </v-alert>

                  <v-btn type="submit" color="primary" block size="large" :loading="saving">
                    {{ $t('install.common.next') }}
                  </v-btn>
                </v-form>
              </v-card>
            </template>

            <!-- Step 4: Site -->
            <template v-slot:item.4>
              <v-card flat>
                <v-form @submit.prevent="submitSite">
                  <v-text-field
                    v-model="siteForm.website_title"
                    :label="$t('install.site.title')"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-text-field>
                  <v-select
                    v-model="siteForm.default_language"
                    :items="languages"
                    item-title="label"
                    item-value="value"
                    :label="$t('install.site.language')"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-select>
                  <v-text-field
                    v-model.number="siteForm.max_users"
                    :label="$t('install.site.maxUsers')"
                    type="number"
                    variant="outlined"
                    density="comfortable"
                    class="mb-2"
                  ></v-text-field>
                  <v-switch
                    v-model="siteForm.allow_registration"
                    :label="$t('install.site.allowRegistration')"
                    color="primary"
                    inset
                    hide-details
                    class="mb-4"
                  ></v-switch>

                  <v-alert v-if="error" type="error" variant="tonal" density="compact" class="mb-4 rounded-lg" closable @click:close="error = ''">
                    {{ error }}
                  </v-alert>

                  <v-btn type="submit" color="primary" block size="large" :loading="saving">
                    {{ $t('install.site.finish') }}
                  </v-btn>
                </v-form>
              </v-card>
            </template>

            <!-- Step 5: Done -->
            <template v-slot:item.5>
              <v-card flat class="text-center">
                <v-icon size="72" color="success" class="mb-4">mdi-check-circle-outline</v-icon>
                <div class="text-h6 mb-2">{{ $t('install.done.title') }}</div>
                <div class="text-body-2 text-grey mb-6">{{ $t('install.done.hint') }}</div>
                <v-btn color="primary" block size="large" to="/login" @click="finishInstall">
                  {{ $t('install.done.goLogin') }}
                </v-btn>
              </v-card>
            </template>
          </v-stepper>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/lib/api'
import { useInstallStore } from '@/store/install'

const { t } = useI18n()
const installStore = useInstallStore()

const step = ref(1)
const error = ref('')
const saving = ref(false)
const testing = ref(false)
const dbTestOk = ref(false)

const dbDrivers = [
  { label: 'MySQL', value: 'mysql' },
  { label: 'PostgreSQL', value: 'postgres' }
]

const languages = [
  { label: 'English', value: 'en' },
  { label: '简体中文', value: 'zh' },
  { label: '日本語', value: 'ja' }
]

const stepItems = computed(() => [
  t('install.step.welcome'),
  t('install.step.database'),
  t('install.step.admin'),
  t('install.step.site'),
  t('install.step.done')
])

const dbForm = reactive({ driver: 'mysql', host: '127.0.0.1', port: 3306, user: '', password: '', name: 'hmimg_db' })
const adminForm = reactive({ username: '', password: '', confirmPassword: '', email: '', phone: '' })
const siteForm = reactive({ website_title: 'HMiMG', default_language: 'zh', max_users: 100, allow_registration: false })

const dbReady = computed(() => installStore.hasDb)

const goNext = () => {
  // 数据库已配置（.env 已有）则跳过数据库步骤
  step.value = installStore.hasDb ? 3 : 2
}

const testDatabase = async () => {
  testing.value = true
  error.value = ''
  dbTestOk.value = false
  try {
    await api.post('/install/database', { ...dbForm, test_only: true })
    dbTestOk.value = true
  } catch (e) {
    error.value = e.response?.data?.error || t('install.db.testFailed')
  } finally {
    testing.value = false
  }
}

const submitDatabase = async () => {
  saving.value = true
  error.value = ''
  try {
    await api.post('/install/database', dbForm)
    await installStore.fetchStatus()
    step.value = 3
  } catch (e) {
    error.value = e.response?.data?.error || t('install.db.failed')
  } finally {
    saving.value = false
  }
}

const submitAdmin = async () => {
  if (adminForm.password !== adminForm.confirmPassword) {
    error.value = t('install.admin.mismatch')
    return
  }
  if (adminForm.password.length < 8) {
    error.value = t('install.admin.passwordHint')
    return
  }
  saving.value = true
  error.value = ''
  try {
    await api.post('/install/admin', adminForm)
    step.value = 4
  } catch (e) {
    error.value = e.response?.data?.error || t('install.admin.failed')
  } finally {
    saving.value = false
  }
}

const submitSite = async () => {
  saving.value = true
  error.value = ''
  try {
    await api.post('/install/site', siteForm)
    await installStore.fetchStatus()
    step.value = 5
  } catch (e) {
    error.value = e.response?.data?.error || t('install.site.failed')
  } finally {
    saving.value = false
  }
}

const finishInstall = () => {
  installStore.checked = false
}

onMounted(async () => {
  if (!installStore.checked) {
    await installStore.fetchStatus()
  }
  if (installStore.installed) return
  // 断点续装：根据后端进度定位起始步骤
  if (installStore.step === 'admin') step.value = 4
  else if (installStore.step === 'database' || installStore.hasDb) step.value = 3
  else step.value = 1
})
</script>

<style scoped>
.install-bg {
  background-color: #f0f2f5;
}

.glass-card {
  backdrop-filter: blur(20px) saturate(160%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(160%) !important;
}

.v-theme--dark .install-bg {
  background-color: #121212;
}

.v-theme--light .glass-card {
  background-color: rgba(255, 255, 255, 0.85) !important;
}

.v-theme--dark .glass-card {
  background-color: rgba(30, 30, 30, 0.85) !important;
}
</style>
