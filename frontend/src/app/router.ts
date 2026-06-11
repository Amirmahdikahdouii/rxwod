import { createRouter, createWebHistory } from 'vue-router'
import WODCreatePage from '@/pages/WODCreatePage.vue'
import WODListPage from '@/pages/WODListPage.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: WODCreatePage },
    { path: '/wods', component: WODListPage },
  ],
})
