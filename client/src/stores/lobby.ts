import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { socket } from '@/lib/ws'

export interface LobbyRoom {
  roomNo: number
  name: string
  gameType: string
  gameName: string
  seatCount: number
  online: number
  sb: number
  bb: number
  minBuyin: number
  maxBuyin: number
  vpip: number
  expireAt: number
  isOfficial: boolean
  clubId: number
  cardCost: number
  insurance: boolean
  status: string
  betLimit: string
  lastHandAt?: number
  handNo?: number
  pot?: number
  avgPot?: number
}

export const useLobby = defineStore('lobby', () => {
  const rooms = ref<LobbyRoom[]>([])
  // 实时变化标记：数字跳动动画用（1.4s 后自动清空）
  const bumps = ref<Record<number, { online?: number; vpip?: number; full?: number }>>({})
  let bumpTimer = 0
  const saved = localStorage.getItem('u5_lobby_filter') as any
const filter = ref<'all' | 'nlh' | 'short6' | 'omaha' | 'niuniu'>(saved || 'all')
  const loaded = ref(false)

  const visible = computed(() => {
    const list = filter.value === 'all' ? rooms.value : rooms.value.filter(r => r.gameType === filter.value)
    const mode = localStorage.getItem('u5_lobby_sort') || 'order'
    if (mode === 'order') return list
    if (mode === 'pot') {
      return [...list].sort((a, b) => (b.pot || 0) - (a.pot || 0) || (b.online / b.seatCount) - (a.online / a.seatCount))
    }
    if (mode === 'hands') {
      return [...list].sort((a, b) => (b.handNo || 0) - (a.handNo || 0) || (b.online / b.seatCount) - (a.online / a.seatCount))
    }
    return [...list].sort((a, b) => {
      const af = a.online >= a.seatCount ? 1 : 0
      const bf = b.online >= b.seatCount ? 1 : 0
      if (af !== bf) return af - bf
      return (b.online / b.seatCount) - (a.online / a.seatCount) || a.roomNo - b.roomNo
    })
  })

  const counts = computed(() => {
    const c: Record<string, number> = { all: rooms.value.length, nlh: 0, short6: 0, omaha: 0, niuniu: 0 }
    rooms.value.forEach(r => { c[r.gameType] = (c[r.gameType] || 0) + 1 })
    return c
  })

  async function fetch() {
    try {
      rooms.value = await socket.call<LobbyRoom[]>('lobby.rooms')
      loaded.value = true
    } catch {}
  }

  function apply(list: LobbyRoom[]) {
    const prev = new Map(rooms.value.map(r => [r.roomNo, r]))
    const now = Date.now()
    const next: Record<number, { online?: number; vpip?: number; full?: number }> = {}
    for (const r of list) {
      const p = prev.get(r.roomNo)
      if (!p) continue
      const b: { online?: number; vpip?: number; full?: number } = {}
      if (Number(r.online) !== Number(p.online)) b.online = now
      if (Number(r.vpip) !== Number(p.vpip)) b.vpip = now
      if (Number(r.online) >= Number(r.seatCount) && Number(p.online) < Number(p.seatCount)) b.full = now
      if (Object.keys(b).length) next[r.roomNo] = b
    }
    rooms.value = list
    loaded.value = true
    if (Object.keys(next).length) {
      bumps.value = next
      window.clearTimeout(bumpTimer)
      bumpTimer = window.setTimeout(() => { bumps.value = {} }, 1400)
    }
  }

  function bumpOf(roomNo: number) {
    return bumps.value[roomNo] || {}
  }

  function setFilter(f: typeof filter.value) {
    filter.value = f
    localStorage.setItem('u5_lobby_filter', f)
  }

  return { rooms, filter, visible, counts, loaded, fetch, apply, setFilter, bumps, bumpOf }
})
