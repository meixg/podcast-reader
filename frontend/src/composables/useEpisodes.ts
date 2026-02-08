import { ref, computed } from 'vue'
import { apiClient } from '@/services/api'
import type { Episode, PaginatedEpisodes } from '@/types/episode'

export function useEpisodes() {
  const episodes = ref<Episode[]>([])
  const total = ref(0)
  const page = ref(1)
  const pageSize = ref(20)
  const totalPages = ref(0)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const hasNextPage = computed(() => page.value < totalPages.value)
  const hasPrevPage = computed(() => page.value > 1)

  async function fetchEpisodes() {
    loading.value = true
    error.value = null
    try {
      const result: PaginatedEpisodes = await apiClient.getEpisodes(page.value, pageSize.value)
      episodes.value = result.episodes
      total.value = result.total
      totalPages.value = result.totalPages
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch episodes'
    } finally {
      loading.value = false
    }
  }

  function nextPage() {
    if (hasNextPage.value) {
      page.value++
      fetchEpisodes()
    }
  }

  function prevPage() {
    if (hasPrevPage.value) {
      page.value--
      fetchEpisodes()
    }
  }

  function goToPage(newPage: number) {
    if (newPage >= 1 && newPage <= totalPages.value) {
      page.value = newPage
      fetchEpisodes()
    }
  }

  function setPageSize(newSize: number) {
    pageSize.value = newSize
    page.value = 1
    fetchEpisodes()
  }

  return {
    episodes,
    total,
    page,
    pageSize,
    totalPages,
    loading,
    error,
    hasNextPage,
    hasPrevPage,
    fetchEpisodes,
    nextPage,
    prevPage,
    goToPage,
    setPageSize
  }
}
