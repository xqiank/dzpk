import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', name: 'login', component: () => import('@/views/Login.vue') },
  {
    path: '/',
    component: () => import('@/components/Layout.vue'),
    children: [
      { path: 'dashboard', name: 'dashboard', component: () => import('@/views/Dashboard.vue'), meta: { title: '总览' } },
      { path: 'tables', name: 'tables', component: () => import('@/views/Tables.vue'), meta: { title: '牌桌监控' } },
      { path: 'rooms', name: 'rooms', component: () => import('@/views/Rooms.vue'), meta: { title: '房间管理' } },
      { path: 'users', name: 'users', component: () => import('@/views/Users.vue'), meta: { title: '用户管理' } },
      { path: 'robots', name: 'robots', component: () => import('@/views/Robots.vue'), meta: { title: '机器人' } },
      { path: 'hands', name: 'hands', component: () => import('@/views/Hands.vue'), meta: { title: '牌局记录' } },
      { path: 'ledger', name: 'ledger', component: () => import('@/views/Ledger.vue'), meta: { title: '财务流水' } },
      { path: 'finance', name: 'finance', component: () => import('@/views/Finance.vue'), meta: { title: '充提审核' } },
      { path: 'settings', name: 'settings', component: () => import('@/views/Settings.vue'), meta: { title: '系统设置' } },
      { path: 'security', name: 'security', component: () => import('@/views/Security.vue'), meta: { title: '安全审计' } },
      { path: 'selfcheck', name: 'selfcheck', component: () => import('@/views/SelfCheck.vue'), meta: { title: '系统自检' } }
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/dashboard' }
]

const router = createRouter({ history: createWebHistory('/admin/'), routes })
router.beforeEach((to) => {
  const token = localStorage.getItem('u5_admin_token')
  if (!token && to.name !== 'login') return { name: 'login' }
  return true
})
export default router
