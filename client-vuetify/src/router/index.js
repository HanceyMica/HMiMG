import { createRouter, createWebHistory } from 'vue-router'
import Cookies from 'js-cookie'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', component: () => import('@/pages/Home.vue') },
      { path: 'library', component: () => import('@/pages/Library.vue') },
      { path: 'album/:id', component: () => import('@/pages/Album.vue') },
      { path: 'collection/:id', component: () => import('@/pages/Collection.vue') },
      { path: 'image/:id', component: () => import('@/pages/ImageDetails.vue') },
      { path: 'admin', component: () => import('@/pages/Admin.vue') },
      { path: 'about', component: () => import('@/pages/About.vue') },
    ],
  },
  {
    path: '/login',
    component: () => import('@/pages/Login.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    component: () => import('@/pages/NotFound.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const token = Cookies.get('token')
  if (to.meta.requiresAuth && !token) {
    next({ path: '/login', query: { redirect: to.fullPath, reason: 'unauthorized' } })
  } else {
    next()
  }
})

export default router
