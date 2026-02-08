<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <NavigationBar />
    <main class="container mx-auto px-4 py-8">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import NavigationBar from '@/components/common/NavigationBar.vue'
</script>
