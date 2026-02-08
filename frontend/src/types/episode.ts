export interface Episode {
  id: string
  title: string
  podcastName: string
  duration: string
  fileSize: number
  downloadDate: string
  showNotes: string
  filePath: string
  coverImagePath?: string
}

export interface PaginatedEpisodes {
  episodes: Episode[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}
