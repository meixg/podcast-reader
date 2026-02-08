<template>
  <span :class="badgeClasses" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
    {{ statusText }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TaskStatus } from '@/types/task'

const props = defineProps<{
  status: TaskStatus
}>()

const statusText = computed(() => {
  switch (props.status) {
    case 'pending':
      return 'Pending'
    case 'downloading':
      return 'Downloading'
    case 'completed':
      return 'Completed'
    case 'failed':
      return 'Failed'
    default:
      return props.status
  }
})

const badgeClasses = computed(() => {
  switch (props.status) {
    case 'pending':
      return 'bg-yellow-100 text-yellow-800'
    case 'downloading':
      return 'bg-blue-100 text-blue-800'
    case 'completed':
      return 'bg-green-100 text-green-800'
    case 'failed':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
})
</script>
