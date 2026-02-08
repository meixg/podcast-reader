<template>
  <Modal :show="show" :title="episode?.title || 'Show Notes'" @close="$emit('close')">
    <div v-if="loading" class="text-center py-8">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
      <p class="mt-2 text-gray-600">Loading show notes...</p>
    </div>
    <div v-else-if="error" class="text-center py-8">
      <p class="text-red-600">{{ error }}</p>
    </div>
    <div v-else class="prose max-w-none">
      <p class="whitespace-pre-wrap">{{ showNotes }}</p>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import Modal from '@/components/common/Modal.vue'
import { apiClient } from '@/services/api'
import type { Episode } from '@/types/episode'

const props = defineProps<{
  show: boolean
  episode: Episode | null
}>()

defineEmits<{
  close: []
}>()

const showNotes = ref('')
const loading = ref(false)
const error = ref<string | null>(null)

watch(() => props.show, async (newShow) => {
  if (newShow && props.episode) {
    await fetchShowNotes()
  }
})

async function fetchShowNotes() {
  if (!props.episode) return

  loading.value = true
  error.value = null
  try {
    const result = await apiClient.getShowNotes(props.episode.id)
    showNotes.value = result.showNotes
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load show notes'
  } finally {
    loading.value = false
  }
}
</script>
