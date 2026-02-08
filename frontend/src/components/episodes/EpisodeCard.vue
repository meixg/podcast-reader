<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 hover:shadow-md transition-shadow cursor-pointer" @click="$emit('click')">
    <div class="flex gap-4">
      <div v-if="episode.coverImagePath" class="flex-shrink-0">
        <img :src="getCoverUrl(episode.coverImagePath)" :alt="episode.title" class="w-24 h-24 rounded object-cover" />
      </div>
      <div class="flex-1 min-w-0">
        <h3 class="text-lg font-semibold text-gray-900 truncate">{{ episode.title }}</h3>
        <p class="text-sm text-gray-600 mt-1">{{ episode.podcastName }}</p>
        <div class="flex items-center gap-4 mt-2 text-sm text-gray-500">
          <span>{{ episode.duration }}</span>
          <span>{{ formatFileSize(episode.fileSize) }}</span>
          <span>{{ formatDate(episode.downloadDate) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Episode } from '@/types/episode'

defineProps<{
  episode: Episode
}>()

defineEmits<{
  click: []
}>()

function getCoverUrl(path: string): string {
  return `http://localhost:8080/static/${path}`
}

function formatFileSize(bytes: number): string {
  const mb = bytes / (1024 * 1024)
  return `${mb.toFixed(1)} MB`
}

function formatDate(date: string): string {
  return new Date(date).toLocaleDateString()
}
</script>
