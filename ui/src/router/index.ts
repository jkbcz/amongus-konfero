import Player from "@/pages/Player.vue";
import Station from "@/pages/Station.vue";
import View from "@/pages/View.vue";
import { createRouter, createWebHistory } from "vue-router";
import Admin from "@/pages/Admin.vue";
import Login from "@/pages/Login.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: Player,
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
    },
    {
      path: "/login",
      component: Login,
    },
  ],
});

export default router;
