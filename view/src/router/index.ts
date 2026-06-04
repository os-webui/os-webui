import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import NotFoundPage from '@/pages/NotFoundPage.vue'
import PluginPage from '@/pages/PluginPage.vue'
import ConfigPage from '@/pages/ConfigPage.vue'
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
      meta: {
        title: ['main.plugin'],
      },
    },
    {
      path: '/config/:id',
      component: ConfigPage,
      meta: {
        title: ['main.config'],
      },
    },
    {
      path: '/about',
      component: AboutView,
      meta: {
        title: ['main.about'],
      },
    },
    {
      path: '/store',
      component: StorePage,
      meta: {
        title: ['main.store'],
      },
    },
    {
      path: '/:pathMatch(.*)*',
      component: NotFoundPage,
      meta: {
        title: ['main.notFound'],
      },
    }
  ],
})
export default router
