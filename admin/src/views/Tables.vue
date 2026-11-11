<template>
  <div class="tables-page">
    <div class="toolbar card">
      <div class="row">
        <select v-model="gameFilter" class="input" style="width:130px">
          <option value="">全部玩法</option><option value="nlh">长牌</option>
          <option value="short6">短牌</option><option value="omaha">奥马哈</option><option value="niuniu">牛牛</option>
        </select>
        <input v-model="kw" class="input" placeholder="搜索房号 / 房名" style="width:200px" />
        <label class="tiny muted row" style="gap:5px"><input type="checkbox" v-model="fullOnly" /> 只看满桌</label>
        <label class="tiny muted row" style="gap:5px"><input type="checkbox" v-model="busyOnly" /> 只看进行中</label>
      </div>
      <div class="row">
        <span class="sum-chip">共 <b>{{ shown.length }}</b> 桌</span>
        <span class="sum-chip">在座 <b>{{ totalSeated }}</b> / {{ totalSeats }}</span>
        <span class="sum-chip gold">底池 <b class="mono">{{ money(totalPot) }}</b></span>
      </div>
    </div>

    <div class="table-grid">
      <div v-for="t in shown" :key="t.roomNo" class="tcard" :class="t.state" @click="detail = t">
        <div class="tc-head">
          <div class="col">
            <b class="tc-title">{{ t.roomNo }} · {{ t.name }}</b>
            <span class="muted tiny">{{ gameName(t.gameType) }} · {{ t.seats }} 人桌 · {{ t.official ? '官方房' : '私人房' }}</span>
          </div>
          <span class="state-pill" :class="t.state">{{ stateName(t.state) }}</span>
        </div>

        <!-- 座位热力盘 -->
        <div class="heatmap" :class="'seats-' + (t.seats || 6)">
          <div v-for="i in (t.seats || 6)" :key="i" class="hm-seat"
               :class="seatClass(t, seatAt(t, i))" :title="seatTitle(t, i)">
            <span class="hm-no">{{ i }}</span>
            <i v-if="seatAt(t, i)" class="hm-dot"></i>
          </div>
          <div class="hm-center">
            <span class="muted tiny">底池</span>
            <b class="mono gold">{{ money(t.pot) }}</b>
            <span class="muted tiny mono">#{{ t.handNo }}</span>
          </div>
        </div>

        <div class="tc-stats">
          <div><span class="muted tiny">在座</span><b class="mono">{{ t.seated }}/{{ t.seats }}</b></div>
          <div><span class="muted tiny">旁观</span><b class="mono">{{ t.observers }}</b></div>
          <div><span class="muted tiny">机器人</span><b class="mono">{{ robotCount(t) }}</b></div>
        </div>
      </div>
    </div>
    <div v-if="!shown.length" class="card" style="text-align:center;padding:26px;color:var(--muted)">没有符合条件的牌桌</div>
  </div>

  <div v-if="detail" class="modal-mask" @click.self="detail = null">
    <div class="modal">
      <div class="between">
        <b>{{ detail.roomNo }} · {{ detail.name }}</b>
        <button class="btn btn-sm" @click="detail = null">关闭</button>
      </div>
      <p class="muted tiny">手数 #{{ detail.handNo }} · 底池 {{ (detail.pot / 100).toFixed(2) }} · 状态 {{ stateName(detail.state) }}</p>
      <table class="table">
        <thead><tr><th>座</th><th>玩家</th><th>筹码</th><th>本街下注</th><th>手牌</th><th>类型</th><th>状态</th></tr></thead>
        <tbody>
          <tr v-for="s in detail.seatsDetail || []" :key="s.seat">
            <td class="mono">{{ s.seat }}</td>
            <td>{{ s.nick }}</td>
            <td class="mono">{{ (s.stack / 100).toFixed(2) }}</td>
            <td class="mono">{{ (s.bet / 100).toFixed(2) }}</td>
            <td class="mono">{{ (s.cards || []).join(' ') || '—' }}</td>
            <td><span class="tag" :class="s.robot ? 'gray' : ''">{{ s.robot ? 'BOT' : '真人' }}</span></td>
            <td>
              <span v-if="s.folded" class="tag red">已弃牌</span>
              <span v-else-if="s.allin" class="tag gold">全下</span>
              <span v-else class="tag">在局</span>
            </td>
          </tr>
          <tr v-if="!(detail.seatsDetail || []).length"><td colspan="7" class="muted" style="text-align:center;padding:16px">该桌暂无在座玩家</td></tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { api, money } from '@/lib/api'
import { adminLive } from '@/lib/live'

const tables = ref<any[]>([])
const detail = ref<any>(null)
const gameFilter = ref('')
const fullOnly = ref(false)
const busyOnly = ref(false)
const kw = ref('')
const shown = computed(() => {
  let rows = tables.value
  if (gameFilter.value) rows = rows.filter(t => t.gameType === gameFilter.value)
  if (fullOnly.value) rows = rows.filter(t => t.seated >= t.seats)
  if (busyOnly.value) rows = rows.filter(t => t.state === 'playing' || t.state === 'showdown' || t.state === 'insurance')
  const k = kw.value.trim().toLowerCase()
  if (k) rows = rows.filter(t => String(t.roomNo).includes(k) || String(t.name || '').toLowerCase().includes(k))
  return rows
})
const totalSeated = computed(() => shown.value.reduce((a, b) => a + Number(b.seated || 0), 0))
const totalSeats = computed(() => shown.value.reduce((a, b) => a + Number(b.seats || 0), 0))
const totalPot = computed(() => shown.value.reduce((a, b) => a + Number(b.pot || 0), 0))
function seatAt(t: any, no: number) {
  return (t.seatsDetail || []).find((s: any) => Number(s.seat) === no - 1)
}
function seatClass(t: any, s: any) {
  if (!s) return 'empty'
  if (s.folded) return 'folded'
  if (s.allin) return 'allin'
  if (s.robot) return 'bot'
  return 'human'
}
function seatTitle(t: any, no: number) {
  const s = seatAt(t, no)
  if (!s) return `${no} 号位 · 空座`
  return `${no} 号位 · ${s.nick} · ${(Number(s.stack || 0) / 100).toFixed(2)}`
}
function robotCount(t: any) {
  return (t.seatsDetail || []).filter((s: any) => s.robot).length
}
let timer = 0
const gameName = (g: string) => ({ nlh: '长牌', short6: '短牌', omaha: '奥马哈', niuniu: '牛牛' } as any)[g] || g
const stateName = (s: string) => ({ playing: '进行中', showdown: '结算中', waiting: '等待中', insurance: '保险中' } as any)[s] || s

async function loadSeats() {
  try {
    const live = await api.get<any[]>('hands/live')
    const map = new Map(live.map((x: any) => [x.roomNo, x.seats]))
    tables.value = tables.value.map(t => ({ ...t, seatsDetail: map.get(t.roomNo) || [] }))
  } catch {}
}

async function load() {
  try {
    tables.value = await api.get<any[]>('tables')
    await loadSeats()
  } catch {}
}

// 座位明细也走 WS（admin.seats，2s 一帧）；HTTP 仅在首屏与断线兜底时拉一次
const seatMap = ref<Record<string, any[]>>({})
function applySeats(map: Record<string, any[]>) {
  if (!map || typeof map !== 'object') return
  tables.value = tables.value.map(t => ({ ...t, seatsDetail: map[String(t.roomNo)] || t.seatsDetail || [] }))
}

let offTables: (() => void) | null = null
let offSeats: (() => void) | null = null
onMounted(() => {
  load()
  offTables = adminLive.on('admin.tables', v => { if (Array.isArray(v)) { tables.value = v.map(t => ({ ...t, seatsDetail: seatMap.value[String(t.roomNo)] || [] })) } })
  offSeats = adminLive.on('admin.seats', v => { if (v && typeof v === 'object') { seatMap.value = v; applySeats(v) } })
  // WS 断开时才回落到 HTTP，避免长期轮询
  timer = window.setInterval(async () => {
    const t = tables.value.find(x => x.seatsDetail?.length) || tables.value[0]
    if (!t) return
    await loadSeats()
  }, 20000)
})
onBeforeUnmount(() => { clearInterval(timer); offTables?.(); offSeats?.() })
</script>

<style scoped>
.col { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.tables-page { display: flex; flex-direction: column; gap: 12px; }
.toolbar { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 12px 14px; flex-wrap: wrap; }
.sum-chip { font-size: 11.5px; color: #9fb8ab; padding: 3px 9px; border-radius: 8px; background: rgba(255, 255, 255, .04); border: 1px solid rgba(60, 120, 88, .3); }
.sum-chip b { color: #d7ffe9; }
.sum-chip.gold b { color: var(--gold); }
.table-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 12px; }
.tcard { position: relative; display: flex; flex-direction: column; gap: 10px; padding: 13px 14px; border-radius: 14px; cursor: pointer;
  background: linear-gradient(160deg, rgba(18, 42, 30, .92), rgba(9, 20, 14, .95));
  border: 1px solid rgba(60, 120, 88, .32); box-shadow: 0 10px 26px rgba(0, 0, 0, .3);
  transition: transform .16s ease, border-color .16s ease, box-shadow .16s ease; }
.tcard:hover { transform: translateY(-2px); border-color: rgba(46, 224, 122, .45); box-shadow: 0 14px 32px rgba(0, 0, 0, .42); }
.tcard.playing { border-color: rgba(46, 224, 122, .38); }
.tcard.showdown { border-color: rgba(255, 216, 117, .38); }
.tc-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 8px; }
.tc-title { font-size: 13.5px; }
.state-pill { flex: none; display: inline-flex; align-items: center; height: 21px; padding: 0 8px; border-radius: 7px; font-size: 10.5px; font-weight: 700;
  background: rgba(255, 255, 255, .06); color: #9fb8ab; border: 1px solid rgba(120, 180, 150, .2); }
.state-pill.playing { background: rgba(46, 224, 122, .16); color: #7dffb0; border-color: rgba(46, 224, 122, .35); }
.state-pill.showdown { background: rgba(255, 216, 117, .16); color: #ffd875; border-color: rgba(255, 216, 117, .35); }
.state-pill.insurance { background: rgba(125, 211, 252, .14); color: #7dd3fc; border-color: rgba(125, 211, 252, .32); }

.heatmap { position: relative; height: 132px; border-radius: 12px; background: radial-gradient(ellipse at 50% 55%, rgba(20, 96, 66, .32), rgba(4, 14, 10, .9));
  border: 1px solid rgba(60, 130, 96, .2); }
.hm-seat { position: absolute; width: 22px; height: 22px; border-radius: 50%; display: flex; align-items: center; justify-content: center;
  font-size: 9px; font-weight: 700; color: #93b3a3; background: rgba(255, 255, 255, .05); border: 1px solid rgba(120, 180, 150, .22); }
.hm-seat .hm-dot { position: absolute; inset: -3px; border-radius: 50%; border: 1px solid transparent; }
.hm-seat.human { background: rgba(46, 224, 122, .22); border-color: rgba(46, 224, 122, .6); color: #d7ffe9; box-shadow: 0 0 10px rgba(46, 224, 122, .35); }
.hm-seat.bot { background: rgba(125, 211, 252, .18); border-color: rgba(125, 211, 252, .5); color: #d6f0ff; }
.hm-seat.folded { opacity: .42; }
.hm-seat.allin { background: rgba(255, 216, 117, .22); border-color: rgba(255, 216, 117, .65); color: #ffe9b0; box-shadow: 0 0 10px rgba(255, 216, 117, .3); }
.hm-center { position: absolute; left: 50%; top: 50%; transform: translate(-50%, -50%); display: flex; flex-direction: column; align-items: center; gap: 1px; }
.hm-center b { font-size: 14px; }
/* 座位坐标（按人数一圈排布） */
.seats-6 .hm-seat:nth-child(1) { left: 50%; bottom: 6px; margin-left: -11px; }
.seats-6 .hm-seat:nth-child(2) { left: 12%; bottom: 26px; }
.seats-6 .hm-seat:nth-child(3) { left: 12%; top: 26px; }
.seats-6 .hm-seat:nth-child(4) { left: 50%; top: 6px; margin-left: -11px; }
.seats-6 .hm-seat:nth-child(5) { right: 12%; top: 26px; }
.seats-6 .hm-seat:nth-child(6) { right: 12%; bottom: 26px; }
.seats-2 .hm-seat:nth-child(1) { left: 50%; bottom: 8px; margin-left: -11px; }
.seats-2 .hm-seat:nth-child(2) { left: 50%; top: 8px; margin-left: -11px; }
.seats-3 .hm-seat:nth-child(1) { left: 50%; bottom: 8px; margin-left: -11px; }
.seats-3 .hm-seat:nth-child(2) { left: 16%; top: 24px; }
.seats-3 .hm-seat:nth-child(3) { right: 16%; top: 24px; }
.seats-4 .hm-seat:nth-child(1) { left: 50%; bottom: 8px; margin-left: -11px; }
.seats-4 .hm-seat:nth-child(2) { left: 14%; top: 50%; margin-top: -11px; }
.seats-4 .hm-seat:nth-child(3) { left: 50%; top: 8px; margin-left: -11px; }
.seats-4 .hm-seat:nth-child(4) { right: 14%; top: 50%; margin-top: -11px; }
.seats-5 .hm-seat:nth-child(1) { left: 50%; bottom: 6px; margin-left: -11px; }
.seats-5 .hm-seat:nth-child(2) { left: 18%; bottom: 30px; }
.seats-5 .hm-seat:nth-child(3) { left: 30%; top: 8px; }
.seats-5 .hm-seat:nth-child(4) { right: 30%; top: 8px; }
.seats-5 .hm-seat:nth-child(5) { right: 18%; bottom: 30px; }
.seats-7 .hm-seat:nth-child(1) { left: 50%; bottom: 6px; margin-left: -11px; }
.seats-7 .hm-seat:nth-child(2) { left: 22%; bottom: 16px; }
.seats-7 .hm-seat:nth-child(3) { left: 12%; top: 46px; }
.seats-7 .hm-seat:nth-child(4) { left: 28%; top: 10px; }
.seats-7 .hm-seat:nth-child(5) { right: 28%; top: 10px; }
.seats-7 .hm-seat:nth-child(6) { right: 12%; top: 46px; }
.seats-7 .hm-seat:nth-child(7) { right: 22%; bottom: 16px; }
.seats-8 .hm-seat:nth-child(1) { left: 50%; bottom: 6px; margin-left: -11px; }
.seats-8 .hm-seat:nth-child(2) { left: 24%; bottom: 14px; }
.seats-8 .hm-seat:nth-child(3) { left: 12%; bottom: 44px; }
.seats-8 .hm-seat:nth-child(4) { left: 12%; top: 44px; }
.seats-8 .hm-seat:nth-child(5) { left: 24%; top: 14px; }
.seats-8 .hm-seat:nth-child(6) { right: 24%; top: 14px; }
.seats-8 .hm-seat:nth-child(7) { right: 12%; top: 44px; }
.seats-8 .hm-seat:nth-child(8) { right: 12%; bottom: 44px; }
.seats-9 .hm-seat:nth-child(1) { left: 50%; bottom: 6px; margin-left: -11px; }
.seats-9 .hm-seat:nth-child(2) { left: 26%; bottom: 12px; }
.seats-9 .hm-seat:nth-child(3) { left: 12%; bottom: 40px; }
.seats-9 .hm-seat:nth-child(4) { left: 16%; top: 26px; }
.seats-9 .hm-seat:nth-child(5) { left: 34%; top: 6px; }
.seats-9 .hm-seat:nth-child(6) { right: 34%; top: 6px; }
.seats-9 .hm-seat:nth-child(7) { right: 16%; top: 26px; }
.seats-9 .hm-seat:nth-child(8) { right: 12%; bottom: 40px; }
.seats-9 .hm-seat:nth-child(9) { right: 26%; bottom: 12px; }

.tc-stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.tc-stats div { display: flex; flex-direction: column; gap: 2px; padding: 7px 8px; border-radius: 8px; background: #0a1710; border: 1px solid #152a1e; }
.tc-stats b { font-size: 12.5px; }
.tiny { font-size: 10.5px; }
@media (max-width: 720px) { .table-grid { grid-template-columns: 1fr; } }
</style>
