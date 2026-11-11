<template>
  <div class="page hand-review handreview-official">
    <header class="hr-head">
      <button class="hr-back" @click="router.back()"><Icon name="back" :size="15" /></button>
      <div class="hr-titles">
        <b class="hr-title">手牌回顾</b>
        <span class="hr-sub">{{ hand?.roomNo ? `房间 #${hand.roomNo} · #${hand.handNo} 局` : '加载中…' }}</span>
      </div>
      <span class="hr-tag">{{ gameLabel }}</span>
      <button class="hr-share" @click="shareCard"><Icon name="send" :size="14" /><span>分享</span></button>
    </header>

    <div v-if="!hand" class="card empty">
      <span class="sm muted">加载中…</span>
    </div>

    <template v-else>
      <div class="replayer-tabs">
        <button class="tab-btn" :class="{ active: tab === 'replay' }" @click="tab = 'replay'">回放</button>
        <button class="tab-btn" :class="{ active: tab === 'fairness' }" @click="tab = 'fairness'">公平存证</button>
      </div>

      <template v-if="tab === 'replay'">
        <div class="card card-pad mini-stage">
          <div class="mini-table-stage">
            <div class="stage-left">
              <span class="stage-label">{{ currentStreetName }}</span>
              <div class="street-stepper-bar">
                <button v-for="(st, si) in streetSteps" :key="'st' + si" class="step-chip"
                        :class="{ active: streetStep === si }" @click="jumpStreet(si)">{{ st }}</button>
              </div>
            </div>
            <div class="pot-banner">
              <span>底池</span>
              <strong>{{ money(currentPot) }}</strong>
            </div>
          </div>
          <div class="board-cards-row">
            <PlayingCard v-for="(c, i) in visibleBoard" :key="'b' + i" :card="c" size="md" />
            <div v-for="i in (5 - visibleBoard.length)" :key="'ph' + i" class="board-ph"></div>
            <span v-if="!board.length" class="empty-board-tip">翻牌前结束</span>
          </div>
          <div class="stepper-controls">
            <button class="rp-btn" :disabled="step <= 0" @click="step = 0">⏮</button>
            <button class="rp-btn" :disabled="step <= 0" @click="step--">◀</button>
            <button class="rp-btn play" @click="togglePlay">{{ playing ? '⏸' : '▶' }}</button>
            <button class="rp-btn" :disabled="step >= maxStep" @click="step++">▶</button>
            <span class="step-num-indicator">{{ step + 1 }}/{{ maxStep + 1 }}</span>
            <input class="rp-slider" type="range" min="0" :max="maxStep" step="1" v-model.number="step" />
          </div>
        </div>

        <div class="card card-pad">
          <b class="sm gold">行动时间线</b>
          <div class="actions-timeline-list">
            <div v-for="(a, ai) in acts" :key="'t' + ai" class="action-timeline-item"
                 :class="{ now: currentAct && currentAct.seq === a.seq, past: ai < step }" @click="step = ai + 1">
              <span class="act-seat">{{ a.seat + 1 }}座</span>
              <span class="act-name">{{ nickOf(a.uid) }}</span>
              <span class="act-type" :class="a.action">{{ actionText(a.action) }}</span>
              <span v-if="a.amount" class="act-amount">{{ money(a.amount) }}</span>
            </div>
            <div v-if="!acts.length" class="empty-list" style="padding:14px;text-align:center">本手暂无公开行动记录</div>
          </div>
        </div>

        <div class="card card-pad">
          <b class="sm gold">手牌结果</b>
          <div v-for="(p, i) in players" :key="i" class="res-row">
            <span class="seat-tag">{{ p.seat }}</span>
            <span class="nick">{{ p.nick }}</span>
            <div class="row gap-2">
              <PlayingCard v-for="(c, ci) in p.cards" :key="ci" :card="c" size="xs" />
            </div>
            <span class="tiny gold" style="margin-left:auto">
              {{ p.handName || '—' }}<em v-if="hand?.gameType === 'niuniu' && multOf(p.handName) > 1" class="mult">×{{ multOf(p.handName) }}</em>
            </span>
            <b :class="p.win > 0 ? 'green' : 'red'">{{ p.win > 0 ? '+' + money(p.win) : money(p.net) }}</b>
          </div>
        </div>
      </template>

      <template v-else>
        <div class="card card-pad proof-card">
          <b class="sm gold">开局承诺哈希（SHA-256）</b>
          <div class="proof-code-box">{{ hand?.seedHash || '—' }}</div>
          <b class="sm gold" style="margin-top:6px">揭示的原始种子</b>
          <div class="proof-code-box">{{ hand?.seed || '本手尚未揭示' }}</div>
          <div class="copy-proof-row">
            <button class="copy-proof-btn" @click="copyProof">复制凭证</button>
            <button class="copy-proof-btn" @click="verifyOpen = true">本机验证</button>
          </div>
          <p class="tiny muted" style="margin:8px 0 0">用种子按 HMAC-SHA256 派生随机流做 Fisher-Yates 洗牌，可复算出与牌局完全一致的牌序。</p>
        </div>
      </template>
    </template>

    <div v-if="verifyOpen" class="mask" @click.self="verifyOpen = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:10px"><b class="h2">公平哈希验证</b>
          <button class="close-x" @click="verifyOpen = false"><Icon name="close" :size="18" /></button></div>
        <div class="hash"><span class="tiny muted">种子哈希</span><code>{{ hand?.seedHash }}</code></div>
        <div class="hash"><span class="tiny muted">公开种子</span><code>{{ hand?.seed || '未公开' }}</code></div>
        <p class="tiny muted">用种子按 HMAC-SHA256 派生随机流做 Fisher-Yates 洗牌，可复算出与牌局完全一致的牌序。</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Icon from '@/components/Icon.vue'
import PlayingCard from '@/components/PlayingCard.vue'
import { socket } from '@/lib/ws'
import { money, actionText } from '@/lib/format'
import '@/styles/handreview-official.css'
import qrcode from 'qrcode-generator'
import { useUI } from '@/stores/ui'

const route = useRoute()
const router = useRouter()
const ui = useUI()
const hand = ref<any>(null)
const sharing = ref(false)

// ===== 牌谱分享图：Canvas 直出，零依赖 =====
function parseCards(str?: string) {
  if (!str) return [] as { rank: string; suit: string; red: boolean }[]
  const map: Record<string, string> = { h: '♥', d: '♦', c: '♣', s: '♠' }
  const out: { rank: string; suit: string; red: boolean }[] = []
  const s2 = String(str)
  for (let i = 0; i + 1 < s2.length; i += 2) {
    out.push({ rank: s2[i], suit: map[s2[i + 1]] || s2[i + 1], red: s2[i + 1] === 'h' || s2[i + 1] === 'd' })
  }
  return out
}

async function shareCard() {
  const h = hand.value
  if (!h || sharing.value) return
  sharing.value = true
  try {
    const W = 640, H = 960
    const cv = document.createElement('canvas')
    cv.width = W; cv.height = H
    const ctx = cv.getContext('2d')
    if (!ctx) throw new Error('画布不可用')

    // 背景
    const bg = ctx.createLinearGradient(0, 0, W, H)
    bg.addColorStop(0, '#0d3b2a')
    bg.addColorStop(.55, '#07231a')
    bg.addColorStop(1, '#04120c')
    ctx.fillStyle = bg
    ctx.fillRect(0, 0, W, H)
    // 顶部光晕
    const glow = ctx.createRadialGradient(W * .5, 120, 20, W * .5, 120, 420)
    glow.addColorStop(0, 'rgba(126,240,180,.22)')
    glow.addColorStop(1, 'rgba(126,240,180,0)')
    ctx.fillStyle = glow
    ctx.fillRect(0, 0, W, 480)

    ctx.textAlign = 'left'
    ctx.fillStyle = 'rgba(255,224,138,.9)'
    ctx.font = '700 20px -apple-system, "PingFang SC", sans-serif'
    ctx.fillText('U5 德扑 · 手牌回顾', 44, 74)

    ctx.fillStyle = 'rgba(255,255,255,.55)'
    ctx.font = '500 16px -apple-system, "PingFang SC", sans-serif'
    ctx.fillText(`房间 ${h.roomNo} · 第 ${h.handNo} 手 · ${gameLabel.value}`, 44, 104)

    // 公共牌
    const board = parseCards(h.board)
    const cw = 78, ch = 108, gap = 12
    const startX = (W - (cw * 5 + gap * 4)) / 2
    board.forEach((c, i) => {
      const x = startX + i * (cw + gap)
      const y = 156
      ctx.fillStyle = '#f7f3e6'
      roundRect(ctx, x, y, cw, ch, 10)
      ctx.fill()
      ctx.fillStyle = c.red ? '#d63a4a' : '#1b2b22'
      ctx.font = '800 34px "Times New Roman", serif'
      ctx.fillText(c.rank, x + 10, y + 42)
      ctx.font = '700 30px serif'
      ctx.fillText(c.suit, x + 10, y + 84)
      ctx.font = '700 18px serif'
      ctx.globalAlpha = .5
      ctx.fillText(c.suit, x + cw - 26, y + ch - 14)
      ctx.globalAlpha = 1
    })
    if (!board.length) {
      ctx.fillStyle = 'rgba(255,255,255,.4)'
      ctx.font = '500 18px -apple-system, "PingFang SC", sans-serif'
      ctx.fillText('翻牌前结束', 44, 200)
    }

    // 底池
    ctx.fillStyle = '#7ff0b0'
    ctx.font = '800 40px -apple-system, "PingFang SC", sans-serif'
    ctx.fillText(`底池 ${money(h.pot)}`, 44, 340)

    // 玩家列表
    let y = 400
    const allPlayers = [...(h.players || [])].sort((a: any, b: any) => (b.net || 0) - (a.net || 0))
    const players = allPlayers.slice(0, 6)
    const restCount = Math.max(0, allPlayers.length - players.length)
    for (const p of players) {
      const win = (p.net || 0) > 0
      ctx.fillStyle = win ? 'rgba(30,80,52,.9)' : 'rgba(255,255,255,.045)'
      roundRect(ctx, 36, y - 30, W - 72, 62, 14)
      ctx.fill()
      ctx.strokeStyle = win ? 'rgba(127,240,176,.55)' : 'rgba(255,255,255,.08)'
      ctx.lineWidth = 1.5
      roundRect(ctx, 36, y - 30, W - 72, 62, 14)
      ctx.stroke()

      if (win) {
        // 赢家：金色描边 + WIN 角标
        ctx.strokeStyle = 'rgba(255,224,138,.9)'
        ctx.lineWidth = 2
        roundRect(ctx, 36, y - 30, W - 72, 62, 14)
        ctx.stroke()
        ctx.fillStyle = 'rgba(255,224,138,.95)'
        roundRect(ctx, W - 128, y - 46, 52, 22, 11)
        ctx.fill()
        ctx.fillStyle = '#2a1c00'
        ctx.font = '800 12px -apple-system, "PingFang SC", sans-serif'
        ctx.textAlign = 'center'
        ctx.fillText('WIN', W - 102, y - 31)
        ctx.textAlign = 'left'
      }

      ctx.fillStyle = '#e8f6ee'
      ctx.font = '700 20px -apple-system, "PingFang SC", sans-serif'
      ctx.fillText(String(p.nick || '玩家').slice(0, 8), 56, y + 2)

      ctx.fillStyle = 'rgba(255,255,255,.55)'
      ctx.font = '500 15px -apple-system, "PingFang SC", sans-serif'
      ctx.fillText([p.handName, (p.cards || []).join(' ')].filter(Boolean).join(' · ').slice(0, 24), 56, y + 24)

      ctx.textAlign = 'right'
      const net = Number(p.net || 0)
      ctx.fillStyle = net === 0 ? 'rgba(255,255,255,.45)' : net > 0 ? '#7ff0b0' : '#ff9f9f'
      ctx.font = '800 22px -apple-system, "PingFang SC", sans-serif' 
      ctx.fillText(`${net >= 0 ? '+' : ''}${money(net)}`, W - 56, y + 10)
      ctx.textAlign = 'left'
      y += 76
    }

    if (restCount > 0) {
      ctx.fillStyle = 'rgba(255,255,255,.4)'
      ctx.font = '500 14px -apple-system, "PingFang SC", sans-serif'
      ctx.fillText(`另有 ${restCount} 位玩家未展示`, 44, y + 4)
    }

    // 二维码：扫码直达本手回放
    try {
      const url = `${location.origin}/records/hand/${h.id}`
      const qr = qrcode(0, 'M')
      qr.addData(url)
      qr.make()
      const n = qr.getModuleCount()
      const size = 96
      const cell = size / n
      const qx = W - size - 26
      const qy = 30
      ctx.fillStyle = '#ffffff'
      roundRect(ctx, qx - 8, qy - 8, size + 16, size + 16, 10)
      ctx.fill()
      ctx.fillStyle = '#0a2419'
      for (let r = 0; r < n; r++) {
        for (let c = 0; c < n; c++) {
          if (qr.isDark(r, c)) ctx.fillRect(qx + c * cell, qy + r * cell, Math.ceil(cell), Math.ceil(cell))
        }
      }
      ;(window as any).__u5shareQR = url
    } catch { /* 二维码失败不影响分享图 */ }

    // 房间号水印
    ctx.save()
    ctx.globalAlpha = .07
    ctx.fillStyle = '#ffffff'
    ctx.font = '900 120px -apple-system, "PingFang SC", sans-serif'
    ctx.textAlign = 'right'
    ctx.fillText(String(h.roomNo), W - 28, H - 120)
    ctx.restore()

    // 底部存证
    ctx.fillStyle = 'rgba(255,255,255,.42)'
    ctx.font = '500 15px ui-monospace, monospace'
    ctx.fillText(`种子哈希 ${String(h.seedHash || '').slice(0, 32)}…`, 44, H - 74)
    ctx.fillStyle = 'rgba(255,224,138,.75)'
    ctx.font = '700 16px -apple-system, "PingFang SC", sans-serif'
    ctx.fillText('U5 德州 · 每一手都可验证', 44, H - 42)

    const dataUrl = cv.toDataURL('image/png')
    ;(window as any).__u5lastShare = dataUrl.length
    const a = document.createElement('a')
    a.href = dataUrl
    a.download = `u5-hand-${h.roomNo}-${h.handNo}.png`
    a.click()
    ui.ok('牌谱图已生成')
  } catch (e: any) {
    ui.err(e.message || '生成失败')
  } finally {
    sharing.value = false
  }
}

function roundRect(ctx: CanvasRenderingContext2D, x: number, y: number, w: number, h: number, r: number) {
  ctx.beginPath()
  ctx.moveTo(x + r, y)
  ctx.arcTo(x + w, y, x + w, y + h, r)
  ctx.arcTo(x + w, y + h, x, y + h, r)
  ctx.arcTo(x, y + h, x, y, r)
  ctx.arcTo(x, y, x + w, y, r)
  ctx.closePath()
}
const verifyOpen = ref(false)

const players = computed<any[]>(() => hand.value?.players || [])
const NIU_MULT: Record<string, number> = { '牛七': 2, '牛八': 3, '牛九': 4, '牛牛': 5, '五花牛': 5, '炸弹': 6, '五小牛': 8 }
const multOf = (n: string) => NIU_MULT[n] || 1
const gameLabel = computed(() => ({ nlh: '长牌', short6: '短牌', omaha: '奥马哈', niuniu: '牛牛' } as any)[hand.value?.gameType] || '牌局')
const board = computed(() => parse(hand.value?.board || ''))

function parse(b: string) {
  const out: string[] = []
  for (let i = 0; i + 1 < b.length; i += 2) out.push(b.slice(i, i + 2))
  return out
}
function nickOf(uid: number) {
  return players.value.find((p: any) => p.uid === uid)?.nick || '玩家'
}

const streets = computed(() => {
  const list = [
    { name: '翻牌前', actions: [] as any[], cards: [] as string[], pot: 0, street: 0 },
    { name: '翻牌', actions: [] as any[], cards: board.value.slice(0, 3), pot: 0, street: 1 },
    { name: '转牌', actions: [] as any[], cards: board.value.slice(3, 4), pot: 0, street: 2 },
    { name: '河牌', actions: [] as any[], cards: board.value.slice(4, 5), pot: 0, street: 3 }
  ]
  ;(hand.value?.actions || []).forEach((a: any) => {
    const s = list[Math.min(3, a.street || 0)]
    s.actions.push(a)
    s.pot = Math.max(s.pot, a.pot || 0)
  })
  return list.filter(s => s.actions.length || s.cards.length)
})

// ===== 回放 =====
const step = ref(0)
const tab = ref<'replay' | 'fairness'>('replay')
const streetSteps = ['翻牌前', '翻牌圈', '转牌圈', '河牌圈', '摊牌']
const streetStep = computed(() => {
  const a = currentAct.value
  if (!a) return 0
  const st = Number(a.street || 0)
  return Math.min(4, st + 1)
})
function jumpStreet(idx: number) {
  if (idx === 0) { step.value = 0; return }
  const target = acts.value.findIndex((a: any) => Number(a.street || 0) === idx - 1)
  if (target < 0) { step.value = Math.min(maxStep.value, idx === 1 ? 0 : maxStep.value); return }
  step.value = target + 1
}
async function copyProof() {
  const h = hand.value
  if (!h) return
  const text = `U5牌局种子承诺\n牌局号：#${h.handNo}\n种子哈希：${h.seedHash || ''}\n洗牌揭示：${h.seed || '未揭示'}`
  try {
    await navigator.clipboard.writeText(text)
    ui.ok('凭证已复制')
  } catch {
    ui.err('复制失败')
  }
}
const playing = ref(false)
let playTimer: number | null = null
const acts = computed<any[]>(() => hand.value?.actions || [])
const maxStep = computed(() => acts.value.length)
const currentAct = computed<any | null>(() => (step.value > 0 ? acts.value[step.value - 1] : null))
const STREET_NAMES = ['翻牌前', '翻牌', '转牌', '河牌']
const currentStreet = computed(() => Math.min(3, currentAct.value?.street ?? 0))
const currentStreetName = computed(() => STREET_NAMES[currentStreet.value])
const streetTagClass = computed(() => ['tag-dim', 'tag-gold', 'tag-gold', 'tag-gold'][currentStreet.value])
const visibleBoard = computed(() => {
  const n = currentStreet.value === 0 ? 0 : currentStreet.value === 1 ? 3 : currentStreet.value === 2 ? 4 : 5
  return board.value.slice(0, n)
})
const currentPot = computed(() => currentAct.value?.pot ?? 0)
const currentText = computed(() => {
  const a = currentAct.value
  if (!a) return '发牌前置盲，等待行动'
  return `${nickOf(a.uid)} ${actionText(a.action)}${a.amount ? ' ' + money(a.amount) : ''}`
})
const curvePoints = computed(() => {
  const vals = [0, ...acts.value.map((a: any) => a.pot || 0)]
  const max = Math.max(1, ...vals)
  return vals.map((v, i) => {
    const x = (i / Math.max(1, vals.length - 1)) * 300
    const y = 42 - (v / max) * 38
    return `${x.toFixed(1)},${y.toFixed(1)}`
  }).join(' ')
})
const curveMarker = computed(() => {
  const vals = [0, ...acts.value.map((a: any) => a.pot || 0)]
  const max = Math.max(1, ...vals)
  const i = Math.min(step.value, vals.length - 1)
  return {
    x: ((i / Math.max(1, vals.length - 1)) * 300).toFixed(1),
    y: (42 - (vals[i] / max) * 38).toFixed(1)
  }
})
function togglePlay() {
  if (playing.value) {
    playing.value = false
    if (playTimer) { clearInterval(playTimer); playTimer = null }
    return
  }
  if (step.value >= maxStep.value) step.value = 0
  playing.value = true
  playTimer = window.setInterval(() => {
    if (step.value >= maxStep.value) {
      playing.value = false
      if (playTimer) { clearInterval(playTimer); playTimer = null }
      return
    }
    step.value++
  }, 900)
}
function jumpTo(a: any) {
  const idx = acts.value.findIndex(x => x.seq === a.seq)
  if (idx >= 0) step.value = idx + 1
}

onMounted(async () => {
  try {
    hand.value = await socket.call('records.hand', { id: Number(route.params.id) })
  } catch (e: any) {
    ui.err(e.message)
  }
})

onUnmounted(() => { if (playTimer) clearInterval(playTimer) })
</script>

<style scoped>
.hand-review { padding-top: 12px; display: flex; flex-direction: column; gap: 8px; }
.h-head { display: flex; align-items: center; gap: 10px; }
.head-btn { height: 34px; padding: 0 12px; border-radius: 10px; background: rgba(255,255,255,.06); border: 1px solid var(--line); font-size: 13px; display: flex; align-items: center; gap: 4px; }
.head-card { display: flex; flex-direction: column; gap: 10px; }
.board-row { display: flex; gap: 6px; align-items: center; min-height: 62px; }
.act-line, .res-row { display: flex; align-items: center; gap: 7px; padding: 5px 0; border-bottom: 1px solid rgba(122,196,156,.08); font-size: 12px; }
.seat-tag { width: 20px; height: 20px; border-radius: 6px; background: rgba(46,224,122,.14); color: #7ff0b0; font-size: 10px; display: flex; align-items: center; justify-content: center; font-weight: 700; flex: none; }
.nick { color: #d6ecdf; min-width: 62px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.act { padding: 1px 7px; border-radius: 6px; font-size: 11px; font-weight: 700; background: rgba(46,224,122,.14); color: #7ff0b0; }
.act.fold { background: rgba(255,120,120,.12); color: #ff9f9f; }
.act.raise, .act.bet, .act.allin { background: rgba(239,207,142,.16); color: #ffd98a; }
.act.check { background: rgba(140,190,255,.12); color: #9ccbff; }
.mult { font-style: normal; margin-left: 3px; color: #ffd98a; font-size: 10px; }
.pot { margin-left: auto; color: var(--muted-2); font-size: 11px; }
.close-x { width: 32px; height: 32px; border-radius: 10px; background: rgba(255,255,255,.06); display: flex; align-items: center; justify-content: center; color: var(--muted); }
.hash { display: flex; flex-direction: column; gap: 5px; padding: 10px; border-radius: 10px; background: rgba(4,16,10,.7); border: 1px solid var(--line); margin-bottom: 8px; }
.hash code { font-size: 11px; color: #9ef7c4; word-break: break-all; }
.card.empty { padding: 30px; text-align: center; }
.board-ph { width: 45px; height: 64px; border-radius: 5px; border: 1px dashed rgba(212,252,232,.16); background: rgba(255,255,255,.03); }
.replay-bar { display: flex; align-items: center; gap: 7px; }
.rp-btn {
  width: 30px; height: 28px; border-radius: 9px; font-size: 12px; color: #cfe9dc;
  background: rgba(255,255,255,.06); border: 1px solid rgba(122,196,156,.18);
  display: flex; align-items: center; justify-content: center;
}
.rp-btn:disabled { opacity: .35; }
.rp-btn.play { background: linear-gradient(180deg, #2ee08a, #12a862); color: #04220f; border-color: rgba(160,255,210,.4); font-weight: 800; }
.rp-count { font-size: 11px; color: var(--muted-2); font-variant-numeric: tabular-nums; min-width: 42px; }
.rp-slider { flex: 1; height: 22px; -webkit-appearance: none; background: transparent; }
.rp-slider::-webkit-slider-runnable-track { height: 4px; border-radius: 4px; background: linear-gradient(90deg, #1b7d55, #ffe08a); }
.rp-slider::-webkit-slider-thumb { -webkit-appearance: none; width: 16px; height: 16px; margin-top: -6px; border-radius: 50%; background: radial-gradient(circle at 35% 30%, #fff, #ffe08a 55%, #d4a738 100%); border: 1px solid rgba(4,40,24,.7); }
.rp-now { display: flex; align-items: center; gap: 8px; font-size: 12px; }
.rp-text { color: #dff3e8; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.rp-pot { margin-left: auto; color: #ffe08a; font-variant-numeric: tabular-nums; }
.pot-curve { width: 100%; height: 46px; display: block; }
.act-line.now { background: rgba(46,224,122,.1); border-radius: 8px; }
.act-line { cursor: pointer; }

.share-btn { margin-left: auto; color: #ffe0a8; border-color: rgba(255, 216, 117, .4); background: rgba(255, 216, 117, .1); }
</style>
