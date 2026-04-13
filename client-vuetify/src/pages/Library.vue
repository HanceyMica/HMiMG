<template>
  <div>
    <div class="d-flex align-center mb-6">
      <v-icon size="32" color="primary" class="mr-3">mdi-image-multiple-outline</v-icon>
      <h1 class="text-h4 font-weight-bold">{{ $t('common.library') }}</h1>
    </div>
    
    <v-tabs v-model="tab" color="primary" class="mb-6 border-b">
      <v-tab value="collections" prepend-icon="mdi-folder-multiple">{{ $t('common.collections') }}</v-tab>
      <v-tab value="albums" prepend-icon="mdi-image-album">{{ $t('common.albums') }}</v-tab>
      <v-tab value="photos" prepend-icon="mdi-image-multiple">{{ $t('common.photos') }}</v-tab>
    </v-tabs>

    <v-window v-model="tab">
      <v-window-item value="collections">
        <v-row>
          <v-col v-for="c in collections" :key="c.id" cols="12" sm="6" md="4" lg="3">
            <v-hover v-slot="{ isHovering, props }">
              <v-card
                v-bind="props"
                :elevation="isHovering ? 8 : 1"
                :to="`/collection/${c.id}`"
                class="rounded-lg transition-swing border"
                flat
              >
                <v-card-item class="bg-surface-variant text-white">
                  <template v-slot:prepend>
                    <v-icon size="32">mdi-folder-outline</v-icon>
                  </template>
                  <v-card-title class="text-subtitle-1">{{ c.name }}</v-card-title>
                </v-card-item>
                <v-card-text class="pa-4 text-body-2 text-grey-darken-1">
                  {{ c.description || 'No description' }}
                </v-card-text>
              </v-card>
            </v-hover>
          </v-col>
        </v-row>
      </v-window-item>

      <v-window-item value="albums">
        <v-row>
          <v-col v-for="a in albums" :key="a.id" cols="12" sm="6" md="4" lg="3">
            <v-hover v-slot="{ isHovering, props }">
              <v-card
                v-bind="props"
                :elevation="isHovering ? 8 : 1"
                :to="`/album/${a.id}`"
                class="rounded-lg overflow-hidden transition-swing border"
                flat
              >
                <v-img
                  v-if="a.cover_image"
                  :src="getImageUrl(a.cover_image)"
                  height="180"
                  cover
                  class="bg-grey-lighten-2"
                >
                  <template v-slot:placeholder>
                    <v-row class="fill-height ma-0" align="center" justify="center">
                      <v-progress-circular indeterminate color="grey-lighten-4"></v-progress-circular>
                    </v-row>
                  </template>
                </v-img>
                <v-sheet v-else height="180" color="grey-lighten-4" class="d-flex align-center justify-center">
                  <v-icon size="48" color="grey-lighten-1">mdi-image-outline</v-icon>
                </v-sheet>
                <v-card-title class="text-subtitle-1 font-weight-bold">{{ a.name }}</v-card-title>
                <v-card-text class="text-caption text-grey text-truncate pt-0">
                  {{ a.description || 'No description' }}
                </v-card-text>
              </v-card>
            </v-hover>
          </v-col>
        </v-row>
      </v-window-item>

      <v-window-item value="photos">
        <v-row>
          <v-col
            v-for="img in images"
            :key="img.id"
            cols="12"
            sm="6"
            md="4"
            lg="3"
          >
            <v-hover v-slot="{ isHovering, props }">
              <v-card
                v-bind="props"
                :elevation="isHovering ? 8 : 1"
                :to="`/image/${img.id}`"
                class="rounded-lg overflow-hidden transition-swing border"
                flat
              >
                <v-img
                  :src="getImageUrl(img.path)"
                  cover
                  height="180"
                  class="bg-grey-lighten-2"
                >
                  <template v-slot:placeholder>
                    <v-row class="fill-height ma-0" align="center" justify="center">
                      <v-progress-circular indeterminate color="grey-lighten-4"></v-progress-circular>
                    </v-row>
                  </template>
                </v-img>
                <v-card-title class="text-subtitle-2 text-truncate pa-3">
                  {{ img.original_name }}
                </v-card-title>
              </v-card>
            </v-hover>
          </v-col>
        </v-row>
      </v-window-item>
    </v-window>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api, { buildUploadedFileUrl } from '@/lib/api'

const tab = ref('collections')
const albums = ref([])
const collections = ref([])
const images = ref([])

const fetchData = async () => {
  try {
    const [albRes, colRes, imgRes] = await Promise.all([
      api.get('/albums'),
      api.get('/collections'),
      api.get('/images')
    ])
    albums.value = albRes.data
    collections.value = colRes.data
    images.value = imgRes.data
  } catch (e) {}
}

const getImageUrl = (path) => {
  return buildUploadedFileUrl(path)
}

onMounted(fetchData)
</script>

<style scoped>
.transition-swing {
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}
</style>
