<template>
  <nav class="bg-white shadow-sm">
    <div class="container mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center space-x-8">
          <h1 class="text-xl font-bold text-gray-900">Podcast Reader</h1>
          <div class="flex space-x-4">
            <router-link
              to="/episodes"
              class="px-3 py-2 rounded-md text-sm font-medium"
              :class="isActive('/episodes') ? 'bg-gray-900 text-white' : 'text-gray-700 hover:bg-gray-100'"
            >
              Podcasts
            </router-link>
            <router-link
              to="/tasks"
              class="px-3 py-2 rounded-md text-sm font-medium"
              :class="isActive('/tasks') ? 'bg-gray-900 text-white' : 'text-gray-700 hover:bg-gray-100'"
            >
              Tasks
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

function isActive(path: string): boolean {
  return route.path === path
}
</script>
