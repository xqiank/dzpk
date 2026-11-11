import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { socket, netState } from '@/lib/ws'

export interface Profile {
  uid: number
  nick: string
  avatar: string
  vip: number
  balance: number
  frozen: number
  cards: number
  signature: string
  hasPayPwd: boolean
  guest?: boolean
  username?: string
  isAdmin?: boolean
  stats?: Record<string, any>
  clubs?: any[]
}

export const useAuth = defineStore('auth', () => {
  const token = ref(localStorage.getItem('u5_token') || '')
  const profile = ref<Profile | null>(null)
  const ready = ref(false)
  const online = computed(() => netState.status === 'open')

  async function bootstrap() {
    if (!token.value) { ready.value = true; return }
    socket.connect(token.value)
    try {
      profile.value = await socket.call<Profile>('user.me')
    } catch { /* 保持 token，重连后再试 */ }
    ready.value = true
  }

  function setSession(t: string, p: Partial<Profile>) {
    token.value = t
    localStorage.setItem('u5_token', t)
    socket.connect(t)
    if (p && (p as any).uid) profile.value = p as Profile
    refresh()
  }

  async function refresh() {
    try { profile.value = await socket.call<Profile>('user.me') } catch {}
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('u5_token')
    profile.value = null
    socket.close()
  }

  function applyBalance(patch: { balance?: number; cards?: number; frozen?: number }) {
    if (!profile.value) return
    if (patch.balance !== undefined) profile.value.balance = patch.balance
    if (patch.cards !== undefined) profile.value.cards = patch.cards
    if (patch.frozen !== undefined) profile.value.frozen = patch.frozen
  }

  return { token, profile, ready, online, bootstrap, setSession, refresh, logout, applyBalance }
})
