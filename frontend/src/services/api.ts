import type { PaginatedEpisodes } from '@/types/episode'
import type { DownloadTask, CreateTaskRequest, APIError } from '@/types/task'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

class APIClient {
  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    try {
      const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers
        }
      })

      if (!response.ok) {
        const error: APIError = await response.json()
        throw new Error(error.error || 'Request failed')
      }

      return await response.json()
    } catch (error) {
      if (error instanceof Error) {
        throw error
      }
      throw new Error('Network error')
    }
  }

  async getEpisodes(page = 1, pageSize = 20): Promise<PaginatedEpisodes> {
    return this.request<PaginatedEpisodes>(`/episodes?page=${page}&pageSize=${pageSize}`)
  }

  async getShowNotes(episodeId: string): Promise<{ showNotes: string }> {
    return this.request<{ showNotes: string }>(`/episodes/${episodeId}/shownotes`)
  }

  async getTasks(): Promise<DownloadTask[]> {
    return this.request<DownloadTask[]>('/tasks')
  }

  async createTask(request: CreateTaskRequest): Promise<DownloadTask> {
    return this.request<DownloadTask>('/tasks', {
      method: 'POST',
      body: JSON.stringify(request)
    })
  }
}

export const apiClient = new APIClient()
