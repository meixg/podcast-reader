<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Download Tasks</h1>
      <button
        @click="showCreateModal = true"
        class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Download
      </button>
    </div>

    <div v-if="loading && tasks.length === 0" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      <p class="mt-4 text-gray-600">Loading tasks...</p>
    </div>

    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-600">{{ error }}</p>
      <button @click="fetchTasks" class="mt-4 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700">
        Retry
      </button>
    </div>

    <TaskList v-else :tasks="tasks" />

    <CreateTaskModal
      :show="showCreateModal"
      :loading="loading"
      @close="showCreateModal = false"
      @submit="handleCreateTask"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTasks } from '@/composables/useTasks'
import { useTaskPolling } from '@/composables/useTaskPolling'
import TaskList from '@/components/tasks/TaskList.vue'
import CreateTaskModal from '@/components/tasks/CreateTaskModal.vue'

const { tasks, loading, error, fetchTasks, createTask } = useTasks()
const showCreateModal = ref(false)

useTaskPolling(fetchTasks)

async function handleCreateTask(url: string) {
  const success = await createTask(url)
  if (success) {
    showCreateModal.value = false
  }
}

onMounted(() => {
  fetchTasks()
})
</script>
