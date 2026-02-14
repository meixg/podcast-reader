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
          <span>{{ formatDuration(episode) }}</span>
          <span>{{ formatPublishTime(episode) }}</span>
          <span>{{ formatFileSize(episode.fileSize) }}</span>
          <span>{{ formatDate(episode.downloadDate) }}</span>
        </div>
      </div>
      <div class="flex-shrink-0 flex items-center gap-2">
        <button
          v-if="props.episode.sourceUrl"
          @click.stop="openSourcePage"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-colors"
          title="打开小宇宙页面"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
          </svg>
          原页面
        </button>
        <button
          @click.stop="downloadAudio"
          class="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
          title="下载音频"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          下载
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Episode } from '@/types/episode'

const props = defineProps<{
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

function formatDuration(episode: Episode): string {
  // Use metadata duration if available, otherwise fall back to episode duration
  if (episode.metadata?.duration) {
    return episode.metadata.duration
  }
  return episode.duration || '--'
}

function formatPublishTime(episode: Episode): string {
  // Use metadata publish_time if available
  if (episode.metadata?.publish_time) {
    return episode.metadata.publish_time
  }
  return '--'
}

function getAudioUrl(filePath: string): string {
  return `http://localhost:8080/static/${filePath}`
}

function downloadAudio(): void {
  const audioUrl = getAudioUrl(props.episode.filePath)
  const link = document.createElement('a')
  link.href = audioUrl
  // Use episode title as filename, sanitize it for filesystem compatibility
  const sanitizedTitle = props.episode.title.replace(/[<>:"/\\|?*]/g, '_')
  link.download = `${sanitizedTitle}.m4a`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function openSourcePage(): void {
  if (props.episode.sourceUrl) {
    window.open(props.episode.sourceUrl, '_blank')
  }
}
</script>
