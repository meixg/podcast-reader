<template>
  <Modal :show="show" title="New Download" @close="handleClose">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label for="url" class="block text-sm font-medium text-gray-700">Episode URL</label>
        <input
          id="url"
          v-model="url"
          type="text"
          placeholder="https://www.xiaoyuzhoufm.com/episode/..."
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          :class="{ 'border-red-500': validationError }"
        />
        <p v-if="validationError" class="mt-1 text-sm text-red-600">{{ validationError }}</p>
      </div>

      <div class="flex justify-end gap-2">
        <button
          type="button"
          @click="handleClose"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="loading"
          class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700 disabled:opacity-50"
        >
          {{ loading ? 'Creating...' : 'Create' }}
        </button>
      </div>
    </form>
  </Modal>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Modal from '@/components/common/Modal.vue'

const props = defineProps<{
  show: boolean
  loading: boolean
}>()

const emit = defineEmits<{
  close: []
  submit: [url: string]
}>()

const url = ref('')
const validationError = ref('')

function validateUrl(): boolean {
  validationError.value = ''

  if (!url.value.trim()) {
    validationError.value = 'URL is required'
    return false
  }

  if (!url.value.includes('xiaoyuzhoufm.com')) {
    validationError.value = 'Must be a Xiaoyuzhou FM URL'
    return false
  }

  return true
}

function handleSubmit() {
  if (validateUrl()) {
    emit('submit', url.value)
  }
}

function handleClose() {
  url.value = ''
  validationError.value = ''
  emit('close')
}
</script>
