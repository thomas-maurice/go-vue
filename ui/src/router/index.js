import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'

import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
      meta: {
        requiresAuth: true,
      }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: {
        requiresAuth: false,
      }
    },
    {
      path: '/auth/callback/:name',
      name: 'oidc-callback',
      component: () => import('../views/OIDCCallback.vue'),
      meta: {
        requiresAuth: false,
      }
    },
    {
      path: '/profile',
      name: 'profile-self',
      component: () => import('../views/ProfileSelfView.vue'),
      meta: {
        requiresAuth: true,
      }
    },
  ],
})

router.beforeEach(async (to, from) => {
  const userStore = useUserStore()
  if (!userStore.logged_in && to.meta.requiresAuth) {
    return { name: 'login' }
  }
})

export default router
