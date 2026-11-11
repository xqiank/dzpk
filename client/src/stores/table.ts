import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { socket } from '@/lib/ws'
import { sfx } from '@/lib/sfx'

export interface SeatState {
  seat: number
  uid: number
  nick: string
  avatar: string
  vip: number
  robot: boolean
  stack: number
  bet: number
  committed: number
  folded: boolean
  allin: boolean
  acted: boolean
  lastAction: string
  lastAmount: number
  deadline: number
  cards: string[]
  cardCount: number
  vpip?: number
  pfr?: number
  muckShown?: boolean
  muckFee?: number
  show: boolean
  handName?: string
  bid?: number
  stake?: number
  win: number
  mult: number
  insured: boolean
  sittingOut: boolean
  hands: number
  net: number
  buyin: number
  connected: boolean
  isMe: boolean
}

export interface TableState {
  roomNo: number
  name: string
  gameType: string
  seatCount: number
  sb: number
  bb: number
  ante: number
  betLimit: string
  pot: number
  board: string[]
  street: number
  handNo: number
  handId: number
  dealer: number
  banker: number
  toAct: number
  state: string
  deadline: number
  expireAt: number
  mySeat: number
  myCards: string[]
  actions: null | {
    canCheck: boolean
    canFold: boolean
    canCall: boolean
    callAmount: number
    canRaise: boolean
    minRaiseTo: number
    maxRaiseTo: number
    canAllIn: boolean
    isRaiseOfBet: boolean
  }
  seats: SeatState[]
  winners: { seat: number; nick: string; amount: number; handName: string; cards: string[]; best?: string[] }[]
  sidePots: { amount: number; seats: number[] }[]
  insurance: null | {
    seat: number; equityP: number; premium: number; payout: number; deadline: number; potSize: number;
    cover?: number; boardCount?: number; outs?: number; outCards?: string[]; odds?: number;
    plans?: { key: string; label: string; premium: number; payout: number }[]; planKey?: string; claimFeeBps?: number;
  }
  serverTime: number
  deckLeft: number
  insuranceEnabled: boolean
  minBuyin: number
  maxBuyin: number
  isOfficial: boolean
  clubId: number
  hostUid?: number
  observers: number
  seedHash: string
  seed: string
  newCards: { seat: number; card: string; index: number; street: number }[]
  actionEv: null | { seat: number; action: string; amount: number; street: number; ts: number }
  winEv: null | { seat: number; amount: number; pot: number }
  version: number
  bankerBid?: number
  nextHandAt?: number
  resultHoldMs?: number
  lastHandAt?: number
  lastHandId?: number
}

export const useTable = defineStore('table', () => {
  const state = ref<TableState | null>(null)
  const prev = ref<TableState | null>(null)
  const messages = ref<{ uid: number; nick: string; avatar: string; content: string; time: number; vip?: number; type?: string; asset?: string; target?: number; data?: string; dur?: number; vid?: string; listeners?: number[] }[]>([])
  const report = ref<any>(null)
  const loading = ref(false)
  const joined = ref(false)

  const mySeat = computed(() => state.value?.seats.find(s => s.isMe) || null)
  const isMyTurn = computed(() => !!(state.value && state.value.mySeat >= 0 && state.value.toAct === state.value.mySeat && state.value.state === 'playing'))
  const contenders = computed(() => (state.value?.seats || []).filter(s => !s.folded))

  async function join(roomNo: number, buyin: number, seat = -1, observe = false, code = '') {
    loading.value = true
    try {
      await socket.call(observe ? 'table.observe' : 'table.join', { roomNo, buyin, seat, code })
      joined.value = true
    } finally {
      loading.value = false
    }
  }

  async function leave() {
    try { await socket.call('table.leave') } catch {}
    joined.value = false
    state.value = null
  }

  async function act(action: string, amount = 0) {
    try {
      await socket.call('table.action', { action, amount })
    } catch (e: any) {
      throw e
    }
  }

  async function rebuy(amount: number) {
    await socket.call('table.rebuy', { amount })
  }

  async function decideInsurance(accept: boolean, plan = '') {
    await socket.call('table.insurance', { accept, plan })
  }

  /** 付费查看对手盖牌（官方「查看对手盖牌」） */
  async function revealMuck(seat: number) {
    return await socket.call('table.revealMuck', { seat })
  }

  async function send(text: string, type = 'text', asset = '', target = 0) {
    await socket.call('table.chat', { text, type, asset, target })
  }

  async function sitOut(v: boolean) {
    await socket.call('table.sitout', { sitOut: v })
  }

  async function loadReport() {
    report.value = await socket.call('table.report')
  }

  // 座位距离衰减：离自己越远声音越轻
  function seatGain(st: TableState, seat: number) {
    if (!st || seat < 0 || st.mySeat < 0) return 0.85
    const n = st.seatCount || 6
    const d = Math.min(Math.abs(seat - st.mySeat), n - Math.abs(seat - st.mySeat))
    return Math.max(0.45, 1 - d * 0.07)
  }

  function apply(next: TableState) {
    const old = state.value
    if (old && old.version === next.version && old.roomNo === next.roomNo) return
    prev.value = old
    state.value = next
    if (next.newCards?.length) sfx.deal()
    if (next.actionEv && old?.actionEv?.ts !== next.actionEv.ts) {
      const a = next.actionEv.action
      const gain = seatGain(next, next.actionEv.seat)
      if (a === 'raise') { sfx.play('raise', { gain }); sfx.voice('raise') }
      else if (a === 'bet') sfx.play('bet', { gain })
      else if (a === 'call') { sfx.play('chips', { gain }); sfx.voice('call') }
      else if (a === 'allin') sfx.allin()
      else if (a === 'check') { sfx.play('check', { gain }); sfx.voice('check') }
      else if (a === 'fold') { sfx.play('fold', { gain }); sfx.voice('fold') }
    }
    if (next.winEv && old?.winEv?.pot !== next.winEv.pot) {
      sfx.pot()
      const mine = next.winEv.seat === next.mySeat
      setTimeout(() => (mine ? sfx.win() : null), 260)
    }
    if (next.pot > (old?.pot || 0) && next.state === 'playing' && next.street === (old?.street ?? -1)) sfx.chip()
    if (next.state === 'playing' && old?.state !== 'playing' && next.handNo !== old?.handNo) sfx.deal()
    if (next.toAct === next.mySeat && old?.toAct !== next.mySeat && next.state === 'playing') sfx.turn()
  }

  function pushMessage(m: any) {
    messages.value.push(m)
    if (messages.value.length > 80) messages.value.shift()
    sfx.tap()
  }

  function reset() {
    state.value = null
    prev.value = null
    messages.value = []
    report.value = null
    joined.value = false
  }

  return {
    state, prev, messages, report, loading, joined, mySeat, isMyTurn, contenders,
    join, leave, act, rebuy, decideInsurance, revealMuck, send, sitOut, loadReport, apply, pushMessage, reset
  }
})
