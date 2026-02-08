import { ref } from 'vue'
import { apiClient } from '@/services/api'
import type { DownloadTask, CreateTaskRequest } from '@/types/task'

export function useTasks() {
  const tasks = ref<DownloadTask[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchTasks() {
    loading.value = true
    error.value = null
    try {
      tasks.value = await apiClient.getTasks()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch tasks'
    } finally {
      loading.value = false
    }
  }

  async function createTask(url: string): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const request: CreateTaskRequest = { url }
      const newTask = await apiClient.createTask(request)
      tasks.value.unshift(newTask)
      return true
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create task'
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    tasks,
    loading,
    error,
    fetchTasks,
    createTask
  }
}
