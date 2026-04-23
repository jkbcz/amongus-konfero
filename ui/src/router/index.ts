import Submit from '@/pages/Submit.vue'
import Station from '@/pages/Station.vue'
import View from '@/pages/View.vue'
import { createRouter, createWebHistory } from 'vue-router'
import Admin from '@/pages/Admin.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: Submit,
    },
    {
      path: "/station",
      component: Station,
    },
    {
      path: "/view",
      component: View,
    },
    {
      path: "/admin",
      component: Admin,
    }
  ],
})

export default router
