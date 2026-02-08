import { ref, onMounted, onUnmounted } from 'vue'

export function useTaskPolling(fetchCallback: () => Promise<void>, intervalMs = 2500) {
  const isPolling = ref(false)
  let intervalId: number | null = null

  function startPolling() {
    if (isPolling.value) return

    isPolling.value = true
    intervalId = window.setInterval(() => {
      fetchCallback()
    }, intervalMs)
  }

  function stopPolling() {
    if (intervalId !== null) {
      clearInterval(intervalId)
      intervalId = null
    }
    isPolling.value = false
  }

  onMounted(() => {
    startPolling()
  })

  onUnmounted(() => {
    stopPolling()
  })

  return {
    isPolling,
    startPolling,
    stopPolling
  }
}
