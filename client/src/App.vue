<template>
  <DesktopLanding v-if="showLanding" @enter="enterDesktop" />
  <div v-else class="shell" :class="{ 'desktop-frame': isDesktop }">
    <router-view v-slot="{ Component }">
      <transition :name="trans" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
    <BottomNav v-if="showNav" />
    <div class="toasts">
      <div v-for="t in ui.toasts" :key="t.id" class="toast" :class="t.kind">{{ t.text }}</div>
    </div>
    <transition name="fade">
      <div v-if="netBanner" class="net-banner" :class="netBanner.kind">
        <span class="nb-dot"></span>
        <b>{{ netBanner.title }}</b>
        <span class="nb-sub">{{ netBanner.sub }}</span>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import BottomNav from '@/components/BottomNav.vue'
import DesktopLanding from '@/components/DesktopLanding.vue'
import { useAuth } from '@/stores/auth'
import { useUI } from '@/stores/ui'
import { useLobby } from '@/stores/lobby'
import { useTable } from '@/stores/table'
import { socket, netState } from '@/lib/ws'
import { sfx } from '@/lib/sfx'

const route = useRoute()
// 断线重连提示：掉线 → 显示重试次数；恢复 → 弹一次「已恢复」
const netBanner = ref<{ kind: string; title: string; sub: string } | null>(null)
let recoverTimer = 0
watch(() => [netState.status, netState.attempt] as const, ([st, at]) => {
  if (!auth.token) { netBanner.value = null; return }
  if (st === 'closed' || st === 'connecting') {
    if (netState.everOpen) {
      netBanner.value = {
        kind: 'down',
        title: '网络连接中断，正在重连…',
        sub: at > 0 ? `第 ${at} 次尝试 · 牌桌状态会自动同步` : '正在建立连接'
      }
    }
  } else if (st === 'open' && netBanner.value) {
    const ms = netState.lastRecoverMs
    netBanner.value = { kind: 'up', title: '已重新连接', sub: ms > 0 ? `${(ms / 1000).toFixed(1)} 秒恢复 · 状态已同步` : '状态已同步' }
    clearTimeout(recoverTimer)
    recoverTimer = window.setTimeout(() => { netBanner.value = null }, 2600)
  }
})
const auth = useAuth()
const ui = useUI()
const lobby = useLobby()
;(window as any).__u5lobby = lobby
;(window as any).__u5auth = auth
const table = useTable()

// 电脑端先给介绍页（对齐官方：手机进游戏、电脑看介绍），可一键切到网页版
const isDesktop = ref(false)
const desktopOK = ref(localStorage.getItem('u5_desktop_mode') === '1')
const showLanding = computed(() => isDesktop.value && !desktopOK.value)
function enterDesktop() {
  desktopOK.value = true
  detectDesktop()
}
function detectDesktop() {
  const touch = 'ontouchstart' in window || navigator.maxTouchPoints > 0
  isDesktop.value = window.innerWidth >= 900 && !touch
  document.body.classList.toggle('desk-bg', isDesktop.value && !showLanding.value)
}

const showNav = computed(() => !['login', 'table', 'community', 'hand', 'private-create'].includes(String(route.name)))
const trans = computed(() => (String(route.name) === 'table' ? 'zoom' : 'fade'))

let offs: (() => void)[] = []

// 全局点击音（官方大多数可点元素都有 click 音；同名音效 60ms 去重，不会叠加）
const TAP_SELECTOR = 'button,.btn,.nav-item,.tab,.seg,.quick,.room,.course,.menu-item,.list-item,.pack,.amt,.cloth-item,.av-pick,.t-round,.tb-btn,.side-tool,.act-btn,.cert-tag'
function onGlobalPointerDown(e: PointerEvent) {
  const el = (e.target as HTMLElement | null)?.closest?.(TAP_SELECTOR) as HTMLElement | null
  if (!el) return
  const disabled = el.hasAttribute('disabled') || el.getAttribute('aria-disabled') === 'true' ||
    el.classList.contains('switch') || el.dataset.silent === '1'
  if (disabled) return
  sfx.tap()
}

onMounted(async () => {
  detectDesktop()
  window.addEventListener('resize', detectDesktop)
  sfx.init()
  document.addEventListener('pointerdown', onGlobalPointerDown, { passive: true })
  await auth.bootstrap()
  if (auth.token) {
    lobby.fetch()
    setTimeout(() => socket.call('user.me').then(p => { auth.profile = p as any }).catch(() => {}), 600)
  }
  offs = [
    socket.on('__open', () => {
      lobby.fetch()
      socket.call('user.me').then(p => { auth.profile = p as any }).catch(() => {})
      if (location.pathname.startsWith('/table/')) socket.call('table.resume').catch(() => {})
    }),
    socket.on('lobby.rooms', (list) => lobby.apply(list)),
    socket.on('user.balance', (d) => auth.applyBalance(d)),
    socket.on('table.state', (st) => table.apply(st)),
    socket.on('chat.message', (m) => table.pushMessage(m)),
    socket.on('notice', (n) => ui.toast(n.text, 'info', 3600)),
    socket.on('kick', (d) => { ui.err(d.reason || '账号在别处登录'); auth.logout() }),
    socket.on('global.chat', (m) => { window.dispatchEvent(new CustomEvent('u5-world', { detail: m })) })
  ]
})

onBeforeUnmount(() => {
  offs.forEach(f => f())
  document.removeEventListener('pointerdown', onGlobalPointerDown)
})
</script>

<style>
.fade-enter-active, .fade-leave-active { transition: opacity .18s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.zoom-enter-active { transition: transform .34s cubic-bezier(.2, .9, .3, 1), opacity .26s ease; }
.zoom-enter-from { transform: scale(.93); opacity: 0; }
.zoom-leave-active { transition: opacity .16s ease; }
.zoom-leave-to { opacity: 0; }
.net-banner {
  position: fixed; top: 10px; left: 50%; transform: translateX(-50%);
  display: flex; align-items: center; gap: 9px; padding: 7px 14px; border-radius: 12px;
  font-size: 12px; font-weight: 700; white-space: nowrap; z-index: 1200; backdrop-filter: blur(6px);
  max-width: 94vw;
}
.net-banner b { font-weight: 800; white-space: nowrap; }
.net-banner.down { background: rgba(60, 34, 8, .92); border: 1px solid rgba(255, 190, 90, .42); color: #ffe0a8; }
.net-banner.up { background: rgba(8, 46, 28, .92); border: 1px solid rgba(120, 240, 170, .42); color: #b6f5cf; }
.net-banner .nb-dot { width: 7px; height: 7px; border-radius: 50%; background: currentColor; animation: nbPulse 1s infinite; }
.net-banner.up .nb-dot { animation: none; }
.net-banner .nb-sub::before { content: "·"; margin-right: 6px; opacity: .55; }
.net-banner .nb-sub { font-size: 11px; font-weight: 500; opacity: .72; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
@keyframes nbPulse { 0%,100% { opacity: 1 } 50% { opacity: .25 } }
.net-flag {
  position: fixed; top: 10px; left: 50%; transform: translateX(-50%);
  font-size: 11px; color: #ffd48a; background: rgba(60, 40, 10, .9);
  padding: 4px 12px; border-radius: 10px; border: 1px solid rgba(255, 200, 120, .3); z-index: 999;
}
</style>
