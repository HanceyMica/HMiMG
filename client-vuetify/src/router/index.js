import { createRouter, createWebHistory } from 'vue-router'
import Cookies from 'js-cookie'
import { useInstallStore } from '@/store/install'

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
      { path: 'admin', component: () => import('@/pages/Admin.vue'), meta: { requiresAdmin: true } },
      { path: 'about', component: () => import('@/pages/About.vue') },
    ],
  },
  {
    path: '/login',
    component: () => import('@/pages/Login.vue'),
  },
  {
    path: '/install',
    component: () => import('@/pages/Install.vue'),
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

function getCurrentUser() {
  try {
    const raw = Cookies.get('user')
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

router.beforeEach(async (to, from, next) => {
  // 安装守卫：未安装强制进向导；已安装屏蔽 /install
  const installStore = useInstallStore()
  if (!installStore.checked) {
    await installStore.fetchStatus()
  }
  if (!installStore.installed && to.path !== '/install') {
    next('/install')
    return
  }
  if (installStore.installed && to.path === '/install') {
    next('/')
    return
  }

  const token = Cookies.get('token')
  if (to.meta.requiresAuth && !token) {
    next({ path: '/login', query: { redirect: to.fullPath, reason: 'unauthorized' } })
    return
  }
  if (to.meta.requiresAdmin) {
    const user = getCurrentUser()
    if (!user || user.role !== 'admin') {
      next({ path: '/' })
      return
    }
  }
  next()
})

export default router
