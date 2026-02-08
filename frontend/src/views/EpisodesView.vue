<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Podcasts</h1>
      <select
        v-model="pageSize"
        @change="handlePageSizeChange"
        class="rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
      >
        <option :value="20">20 per page</option>
        <option :value="50">50 per page</option>
        <option :value="100">100 per page</option>
      </select>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      <p class="mt-4 text-gray-600">Loading episodes...</p>
    </div>

    <div v-else-if="error" class="text-center py-12">
      <p class="text-red-600">{{ error }}</p>
      <button @click="fetchEpisodes" class="mt-4 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700">
        Retry
      </button>
    </div>

    <template v-else>
      <EpisodeList :episodes="episodes" @select="handleSelectEpisode" />

      <Pagination
        v-if="totalPages > 1"
        :current-page="page"
        :total-pages="totalPages"
        :total="total"
        :page-size="pageSize"
        @prev="prevPage"
        @next="nextPage"
        @goto="goToPage"
      />
    </template>

    <ShowNotesModal
      :show="showModal"
      :episode="selectedEpisode"
      @close="showModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useEpisodes } from '@/composables/useEpisodes'
import EpisodeList from '@/components/episodes/EpisodeList.vue'
import ShowNotesModal from '@/components/episodes/ShowNotesModal.vue'
import Pagination from '@/components/common/Pagination.vue'
import type { Episode } from '@/types/episode'

const {
  episodes,
  total,
  page,
  pageSize,
  totalPages,
  loading,
  error,
  fetchEpisodes,
  nextPage,
  prevPage,
  goToPage,
  setPageSize
} = useEpisodes()

const showModal = ref(false)
const selectedEpisode = ref<Episode | null>(null)

function handleSelectEpisode(episode: Episode) {
  selectedEpisode.value = episode
  showModal.value = true
}

function handlePageSizeChange() {
  setPageSize(pageSize.value)
}

onMounted(() => {
  fetchEpisodes()
})
</script>
