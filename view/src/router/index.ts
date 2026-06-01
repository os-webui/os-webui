import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import PluginPage from '@/pages/PluginPage.vue'
import StorePage from '@/pages/StorePage.vue'
import AboutView from '@/views/AboutView.vue'
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: HomePage,
    },
    {
      path: '/plugin/:id',
      component: PluginPage,
    },
    {
      path: '/about',
      component: AboutView,
      // component: () => import('../views/AboutView.vue'), // lazy-loaded
    },
    {
      path: '/store',
      component: StorePage,
    },
  ],
})

export default router
