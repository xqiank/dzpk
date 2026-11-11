<template>
  <div class="layout" :class="{ collapsed }">
    <aside class="sidebar">
      <div class="logo">
        <span class="mark">♠</span>
        <div class="col logo-text">
          <b>U5 德扑</b>
          <span class="muted tiny">运营后台 · PRO</span>
        </div>
        <button class="collapse-btn" :title="collapsed ? '展开菜单' : '收起菜单'" @click="toggleCollapse">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round">
            <path :d="collapsed ? 'M9 6l6 6-6 6' : 'M15 6l-6 6 6 6'" />
          </svg>
        </button>
      </div>
      <nav class="nav">
        <template v-for="group in navGroups" :key="group.title">
          <div class="nav-group-title">{{ group.title }}</div>
          <router-link v-for="m in group.items" :key="m.path" :to="m.path" class="nav-item" active-class="active" :title="m.label">
            <span class="ico" v-html="m.icon"></span>
            <span class="nav-label">{{ m.label }}</span>
            <span v-if="m.badge" class="nav-badge">{{ m.badge }}</span>
          </router-link>
        </template>
      </nav>
      <div class="side-foot">
        <div class="online"><i></i><span class="nav-label">在线 {{ online }}</span></div>
        <button class="btn btn-sm nav-label" @click="logout">退出</button>
      </div>
    </aside>
    <main class="main">
      <header class="topbar">
        <div class="crumbs">
          <span class="crumb-muted">控制台</span>
          <span class="crumb-sep">/</span>
          <b class="crumb-cur">{{ title }}</b>
          <span class="muted tiny hide-sm">{{ subtitle }}</span>
        </div>
        <div class="row">
          <span class="live-chip"><i></i>{{ online }} 在线</span>
          <span class="tag gold mono">{{ now }}</span>
          <button class="btn btn-sm" @click="refresh">刷新</button>
          <button class="btn btn-sm" :title="fullscreen ? '退出全屏' : '全屏'" @click="toggleFullscreen">{{ fullscreen ? '退出全屏' : '全屏' }}</button>
        </div>
      </header>
      <section class="content">
        <router-view :key="route.fullPath" />
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { api, logout as doLogout } from '@/lib/api'

const route = useRoute()
const online = ref(0)
const now = ref('')
let timer = 0
let clock = 0

function wrapSvg(p: string) {
  return `<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">${p}</svg>`
}
const ICONS_RAW = {
  dashboard: '<rect x="3" y="3" width="7" height="9" rx="1.5"/><rect x="14" y="3" width="7" height="5" rx="1.5"/><rect x="14" y="12" width="7" height="9" rx="1.5"/><rect x="3" y="16" width="7" height="5" rx="1.5"/>',
  tables: '<circle cx="12" cy="12" r="8.5"/><path d="M12 3.5v17M3.5 12h17"/>',
  rooms: '<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 9h18M9 9v11"/>',
  users: '<circle cx="9" cy="8" r="3.4"/><path d="M2.5 19c1-3 3.4-4.6 6.5-4.6s5.5 1.6 6.5 4.6"/><path d="M17 8h4M19 6v4"/>',
  robots: '<rect x="4" y="7" width="16" height="12" rx="3"/><path d="M12 3v4M9 12h.01M15 12h.01M9 16h6"/>',
  hands: '<path d="M6 3h8l4 4v14H6z"/><path d="M14 3v4h4M9 12h6M9 16h6"/>',
  ledger: '<ellipse cx="12" cy="7" rx="7.5" ry="3"/><path d="M4.5 7v10c0 1.7 3.4 3 7.5 3s7.5-1.3 7.5-3V7"/>',
  selfcheck: '<path d="M4 5h16v4H4zM4 13h9v6H4zM17 13h3v6h-3z"/>',
  security: '<path d="M12 3l7 3v5.5c0 4.2-2.9 7.4-7 9.5-4.1-2.1-7-5.3-7-9.5V6l7-3z"/><path d="M9 12l2 2 4-4"/>',
  settings: '<circle cx="12" cy="12" r="3"/><path d="M19 12a7 7 0 00-.1-1.2l2-1.5-2-3.4-2.3 1a7 7 0 00-2-1.2L14.2 3h-4l-.4 2.7a7 7 0 00-2 1.2l-2.3-1-2 3.4 2 1.5a7 7 0 000 2.4l-2 1.5 2 3.4 2.3-1a7 7 0 002 1.2l.4 2.7h4l.4-2.7a7 7 0 002-1.2l2.3 1 2-3.4-2-1.5c.1-.4.1-.8.1-1.2z"/>'
}
type NavItem = { path: string; label: string; icon: string; badge?: string }
type NavGroup = { title: string; items: NavItem[] }
const ICONS = Object.fromEntries(Object.entries(ICONS_RAW).map(([k, v]) => [k, wrapSvg(v as string)])) as Record<string, string>
const navGroups: NavGroup[] = [
  { title: '概览', items: [{ path: '/dashboard', label: '总览', icon: ICONS.dashboard }] },
  {
    title: '运营',
    items: [
      { path: '/tables', label: '牌桌监控', icon: ICONS.tables },
      { path: '/rooms', label: '房间管理', icon: ICONS.rooms },
      { path: '/users', label: '用户管理', icon: ICONS.users },
      { path: '/robots', label: '机器人', icon: ICONS.robots },
      { path: '/hands', label: '牌局记录', icon: ICONS.hands }
    ]
  },
  { title: '财务', items: [
    { path: '/ledger', label: '财务流水', icon: ICONS.ledger },
    { path: '/finance', label: '充提审核', icon: ICONS.ledger }
  ] },
  {
    title: '系统',
    items: [
      { path: '/selfcheck', label: '系统自检', icon: ICONS.selfcheck },
      { path: '/security', label: '安全审计', icon: ICONS.security },
      { path: '/settings', label: '系统设置', icon: ICONS.settings }
    ]
  }
]
const collapsed = ref(localStorage.getItem('u5_admin_collapsed') === '1')
const fullscreen = ref(false)
function toggleCollapse() {
  collapsed.value = !collapsed.value
  localStorage.setItem('u5_admin_collapsed', collapsed.value ? '1' : '0')
}
function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen?.().then(() => { fullscreen.value = true }).catch(() => {})
  } else {
    document.exitFullscreen?.().then(() => { fullscreen.value = false }).catch(() => {})
  }
}

const title = computed(() => (route.meta.title as string) || '总览')
const subtitle = computed(() => {
  const map: Record<string, string> = {
    总览: '平台核心指标、实时在线与最近牌局',
    牌桌监控: '实时查看每张桌的牌局进程与座位状态',
    房间管理: '官方房/私人房配置、盲注与买入区间',
    用户管理: '余额、房卡、封禁与账号安全',
    机器人: '机器人风格、房间分配与产能',
    牌局记录: '每手牌的行动线与公平哈希存证',
    财务流水: '充值、带入、带出与后台调整明细',
    安全审计: '限流、注入探测、异常登录记录',
    系统自检: '服务健康、运行时指标与最近一次全量验收',
    系统设置: '公告、全服广播与运营参数'
  }
  return map[title.value] || ''
})

async function refresh() {
  try {
    const d = await api.get<{ online: number }>('dashboard')
    online.value = d.online
  } catch {}
}

onMounted(() => {
  refresh()
  timer = window.setInterval(refresh, 5000)
  const tick = () => { now.value = new Date().toLocaleTimeString('zh-CN', { hour12: false }) }
  tick()
  clock = window.setInterval(tick, 1000)
  window.addEventListener('admin-refresh', refresh)
})
onBeforeUnmount(() => { clearInterval(timer); clearInterval(clock) })

function logout() { doLogout() }
</script>

<style scoped>
.layout { display: flex; height: 100vh; }
.sidebar { width: 224px; flex: none; background: linear-gradient(180deg, rgba(10, 26, 18, .96), rgba(5, 12, 8, .98)); border-right: 1px solid rgba(50, 100, 74, .35); display: flex; flex-direction: column; transition: width .2s ease; }
.layout.collapsed .sidebar { width: 68px; }
.logo { display: flex; align-items: center; gap: 10px; padding: 15px 14px 13px; border-bottom: 1px solid #14261b; }
.mark { width: 34px; height: 34px; flex: none; border-radius: 9px; background: radial-gradient(circle at 34% 28%, rgba(46, 224, 122, .3), rgba(9, 30, 20, .9)); border: 1px solid #255c3d; color: var(--gold); display: flex; align-items: center; justify-content: center; font-size: 18px; box-shadow: 0 4px 14px rgba(0, 0, 0, .45); }
.logo-text { min-width: 0; }
.logo-text b { font-size: 14px; letter-spacing: .3px; }
.logo-text .tiny { font-size: 10px; color: #6f9a83; }
.collapse-btn { margin-left: auto; width: 24px; height: 24px; border-radius: 7px; background: rgba(255, 255, 255, .05); border: 1px solid rgba(120, 200, 160, .2); color: #a9c3b4; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.collapse-btn:hover { color: #d7ffe9; border-color: rgba(120, 200, 160, .45); }
.layout.collapsed .logo-text, .layout.collapsed .collapse-btn { display: none; }
.layout.collapsed .logo { justify-content: center; padding: 15px 0 13px; }
.nav { flex: 1; padding: 8px 8px 12px; display: flex; flex-direction: column; gap: 1px; overflow-y: auto; }
.nav-group-title { font-size: 10px; letter-spacing: 1.2px; color: #4f7a63; font-weight: 700; padding: 10px 10px 4px; text-transform: uppercase; }
.layout.collapsed .nav-group-title { text-align: center; font-size: 8.5px; letter-spacing: 0; padding: 8px 0 3px; }
.nav-item { position: relative; display: flex; align-items: center; gap: 10px; height: 38px; padding: 0 11px; border-radius: 10px; color: #a9c3b4; font-size: 13px; transition: background .15s, color .15s, border-color .15s; border: 1px solid transparent; }
.nav-item:hover { background: #0f2016; color: var(--text); }
.nav-item.active { background: linear-gradient(90deg, rgba(46, 224, 122, .18), rgba(46, 224, 122, .03)); color: #b7ffd4; border-color: #1f4a33; }
.nav-item.active::before { content: ""; position: absolute; left: -1px; top: 8px; bottom: 8px; width: 3px; border-radius: 0 3px 3px 0; background: linear-gradient(180deg, #7dffb0, #12a35f); box-shadow: 0 0 10px rgba(46, 224, 122, .8); }
.nav-badge { margin-left: auto; font-size: 9.5px; padding: 1px 6px; border-radius: 8px; background: rgba(46, 224, 122, .18); color: #7dffb0; }
.ico { width: 16px; height: 16px; display: inline-flex; flex: none; }
.ico :deep(svg) { width: 16px; height: 16px; }
.layout.collapsed .nav-label, .layout.collapsed .nav-badge { display: none; }
.layout.collapsed .nav-item { justify-content: center; padding: 0; }
.side-foot { padding: 11px 12px; border-top: 1px solid #14261b; display: flex; align-items: center; justify-content: space-between; gap: 6px; }
.online { display: flex; align-items: center; gap: 6px; font-size: 11px; color: #8fd8b4; }
.online i { width: 6px; height: 6px; border-radius: 50%; background: #2ee07a; box-shadow: 0 0 8px #2ee07a; animation: pulse 1.8s infinite ease-in-out; }
@keyframes pulse { 0%, 100% { opacity: .55; transform: scale(.9) } 50% { opacity: 1; transform: scale(1.1) } }
.layout.collapsed .side-foot { flex-direction: column; gap: 8px; padding: 10px 6px; }
.main { flex: 1; min-width: 0; display: flex; flex-direction: column; }
.topbar { height: 60px; flex: none; display: flex; align-items: center; justify-content: space-between; padding: 0 18px; border-bottom: 1px solid #14261b;
  background: linear-gradient(180deg, rgba(9, 22, 15, .96), rgba(7, 17, 12, .96)); backdrop-filter: blur(8px); }
.crumbs { display: flex; align-items: center; gap: 8px; min-width: 0; }
.crumb-muted { font-size: 12px; color: #6f9a83; }
.crumb-sep { color: #33513f; }
.crumb-cur { font-size: 14.5px; }
.hide-sm { max-width: 340px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.live-chip { display: inline-flex; align-items: center; gap: 5px; height: 24px; padding: 0 9px; border-radius: 8px; font-size: 11px; font-weight: 700;
  background: rgba(46, 224, 122, .12); border: 1px solid rgba(46, 224, 122, .32); color: #7dffb0; }
.live-chip i { width: 5px; height: 5px; border-radius: 50%; background: #2ee07a; box-shadow: 0 0 7px #2ee07a; }
.content { flex: 1; overflow-y: auto; padding: 16px 18px 30px; }
.row { display: flex; align-items: center; gap: 8px; }
.col { display: flex; flex-direction: column; }
.tiny { font-size: 10.5px; }
@media (max-width: 860px) {
  .sidebar { width: 68px; }
  .sidebar .logo-text, .sidebar .nav-label, .sidebar .nav-badge, .sidebar .collapse-btn { display: none; }
  .sidebar .nav-item { justify-content: center; padding: 0; }
  .hide-sm { display: none; }
}
</style>
