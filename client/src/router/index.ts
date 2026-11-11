import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/cash' },
  { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
  { path: '/cash', name: 'cash', component: () => import('@/views/CashGames.vue') },
  { path: '/lobby', redirect: '/cash' },
  { path: '/private', name: 'private', component: () => import('@/views/Private.vue') },
  { path: '/private/create', name: 'private-create', component: () => import('@/views/PrivateCreate.vue') },
  { path: '/club', name: 'club', component: () => import('@/views/Club.vue') },
  { path: '/records', name: 'records', component: () => import('@/views/Records.vue') },
  { path: '/records/hand/:id', name: 'hand', component: () => import('@/views/HandReview.vue') },
  { path: '/profile', name: 'profile', component: () => import('@/views/Profile.vue') },
  { path: '/table/:roomNo', name: 'table', component: () => import('@/views/Table.vue') },
  { path: '/community', name: 'community', component: () => import('@/views/Community.vue') },
  { path: '/bills', name: 'bills', component: () => import('@/views/Bills.vue') },
  { path: '/:pathMatch(.*)*', redirect: '/cash' }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
})

router.beforeEach((to) => {
  const token = localStorage.getItem('u5_token')
  if (!token && to.name !== 'login') return { name: 'login', query: { redirect: to.fullPath } }
  return true
})

export default router
