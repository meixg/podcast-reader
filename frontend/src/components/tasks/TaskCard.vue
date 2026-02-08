<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
    <div class="flex items-start justify-between">
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2">
          <TaskStatusBadge :status="task.status" />
          <span class="text-xs text-gray-500">{{ formatDate(task.createdAt) }}</span>
        </div>
        <p class="mt-2 text-sm text-gray-900 truncate">{{ task.url }}</p>
        <div v-if="task.status === 'downloading' && task.progress !== undefined" class="mt-2">
          <div class="flex items-center justify-between text-xs text-gray-600 mb-1">
            <span>Downloading...</span>
            <span>{{ task.progress }}%</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-2">
            <div
              class="bg-blue-600 h-2 rounded-full transition-all duration-300"
              :style="{ width: `${task.progress}%` }"
            ></div>
          </div>
        </div>
        <p v-if="task.errorMessage" class="mt-1 text-sm text-red-600">{{ task.errorMessage }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import TaskStatusBadge from './TaskStatusBadge.vue'
import type { DownloadTask } from '@/types/task'

defineProps<{
  task: DownloadTask
}>()

function formatDate(date: string): string {
  return new Date(date).toLocaleString()
}
</script>
