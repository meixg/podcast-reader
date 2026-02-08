<template>
  <div v-if="episodes.length === 0" class="text-center py-12 bg-gray-50 rounded-lg">
    <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
    </svg>
    <h3 class="mt-2 text-sm font-medium text-gray-900">No episodes found</h3>
    <p class="mt-1 text-sm text-gray-500">Download some podcasts to get started.</p>
  </div>
  <div v-else class="space-y-4">
    <EpisodeCard
      v-for="episode in episodes"
      :key="episode.id"
      :episode="episode"
      @click="$emit('select', episode)"
    />
  </div>
</template>

<script setup lang="ts">
import EpisodeCard from './EpisodeCard.vue'
import type { Episode } from '@/types/episode'

defineProps<{
  episodes: Episode[]
}>()

defineEmits<{
  select: [episode: Episode]
}>()
</script>
