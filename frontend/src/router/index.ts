import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/episodes'
    },
    {
      path: '/episodes',
      name: 'episodes',
      component: () => import('@/views/EpisodesView.vue')
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: () => import('@/views/TasksView.vue')
    }
  ]
})

export default router
