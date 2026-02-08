export type TaskStatus = 'pending' | 'downloading' | 'completed' | 'failed'

export interface DownloadTask {
  id: string
  url: string
  status: TaskStatus
  createdAt: string
  completedAt?: string
  progress?: number
  errorMessage?: string
  episodeId?: string
}

export interface CreateTaskRequest {
  url: string
}

export interface APIError {
  error: string
  code: string
  details?: Record<string, any>
}
