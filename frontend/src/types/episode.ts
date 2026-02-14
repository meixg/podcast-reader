export interface PodcastMetadata {
  duration?: string
  publish_time?: string
  episode_title?: string
  podcast_name?: string
  source_url?: string
  extracted_at: string
}

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
  sourceUrl?: string
  metadata?: PodcastMetadata
}

export interface PaginatedEpisodes {
  episodes: Episode[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}
