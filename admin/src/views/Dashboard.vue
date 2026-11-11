<template>
  <div class="dash">
    <div class="kpis">
      <div v-for="k in kpis" :key="k.key" class="card kpi">
        <div class="kpi-top">
          <span class="muted tiny">{{ k.label }}</span>
          <span class="kpi-ico" v-html="k.icon"></span>
        </div>
        <b class="kpi-val">{{ k.value }}</b>
        <div class="kpi-foot">
          <em>{{ k.unit }}</em>
          <span v-if="k.delta !== null" class="delta" :class="k.delta >= 0 ? 'up' : 'down'">
            {{ k.delta >= 0 ? '▲' : '▼' }} {{ Math.abs(k.delta) }}{{ k.deltaUnit }}
          </span>
        </div>
      </div>
    </div>

    <div class="chart-grid">
      <div class="card chart-card">
        <div class="between" style="margin-bottom:6px">
          <div class="col">
            <b>近 7 日牌局走势</b>
            <span class="muted tiny">每日完成手数 · 鼠标悬停查看明细</span>
          </div>
          <div class="row">
            <span class="legend-dot hands"></span><span class="muted tiny">牌局</span>
            <span class="legend-dot users"></span><span class="muted tiny">新增用户</span>
          </div>
        </div>
        <svg class="trend-svg" viewBox="0 0 720 220" preserveAspectRatio="none">
          <defs>
            <linearGradient id="areaHands" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="rgba(46,224,122,.45)" />
              <stop offset="100%" stop-color="rgba(46,224,122,0)" />
            </linearGradient>
            <linearGradient id="lineHands" x1="0" y1="0" x2="1" y2="0">
              <stop offset="0%" stop-color="#2ee07a" />
              <stop offset="100%" stop-color="#9affc6" />
            </linearGradient>
          </defs>
          <g class="grid">
            <line v-for="g in gridLines" :key="'g' + g" x1="0" :y1="g" x2="720" :y2="g" />
          </g>
          <polygon :points="areaPoints" fill="url(#areaHands)" />
          <polyline :points="linePoints" fill="none" stroke="url(#lineHands)" stroke-width="2.4" stroke-linejoin="round" stroke-linecap="round" />
          <polyline :points="userPoints" fill="none" stroke="#ffd875" stroke-width="1.8" stroke-dasharray="5 4" stroke-linejoin="round" />
          <g v-for="(p, i) in pointList" :key="'p' + i">
            <circle :cx="p.x" :cy="p.y" r="3.6" fill="#06150e" stroke="#2ee07a" stroke-width="2" />
            <text :x="p.x" y="214" class="axis-label" text-anchor="middle">{{ p.label }}</text>
            <text :x="p.x" :y="p.y - 10" class="value-label" text-anchor="middle">{{ p.value }}</text>
          </g>
        </svg>
        <div class="chart-foot">
          <span>合计 <b class="mono gold">{{ trendUsers }} 人</b> 新增 · 充值 <b class="mono gold">{{ money(trendRecharge) }}</b></span>
          <span class="muted tiny">数据每 2 秒推送</span>
        </div>
      </div>

      <div class="card dist-card">
        <b>玩法分布</b>
        <span class="muted tiny">今日各玩法手数占比</span>
        <div class="dist-bar">
          <span v-for="g in byGame" :key="g.game" class="dist-seg" :style="{ width: gamePct(g) + '%', background: gameColor(g.game) }"></span>
        </div>
        <div class="dist-legend">
          <div v-for="g in byGame" :key="'l' + g.game" class="dist-row">
            <span class="legend-dot" :style="{ background: gameColor(g.game) }"></span>
            <span class="dist-name">{{ gameName(g.game) }}</span>
            <span class="mono muted">{{ g.hands }}</span>
            <b class="mono">{{ gamePct(g) }}%</b>
          </div>
          <div v-if="!byGame.length" class="muted tiny">暂无对局数据</div>
        </div>
      </div>
    </div>

    <div class="panels">
      <div class="card">
        <div class="between" style="margin-bottom:10px">
          <b>实时牌桌</b>
          <span class="muted tiny">共 {{ tables.length }} 桌 · 每 2 秒刷新</span>
        </div>
        <table class="table">
          <thead><tr><th>房间</th><th>玩法</th><th>状态</th><th>手数</th><th>底池</th><th>座位</th></tr></thead>
          <tbody>
            <tr v-for="t in tables" :key="t.roomNo">
              <td>{{ t.roomNo }} · {{ t.name }}</td>
              <td><span class="tag">{{ gameName(t.gameType) }}</span></td>
              <td><span class="state-pill" :class="t.state">{{ stateName(t.state) }}</span></td>
              <td class="mono">#{{ t.handNo }}</td>
              <td class="mono">{{ money(t.pot) }}</td>
              <td>{{ t.seated }}/{{ t.seats }} <span class="muted">旁观 {{ t.observers }}</span></td>
            </tr>
            <tr v-if="!tables.length"><td colspan="6" class="muted" style="text-align:center;padding:18px">暂无运行中的牌桌</td></tr>
          </tbody>
        </table>
      </div>
      <div class="card">
        <div class="between" style="margin-bottom:10px"><b>最近牌局</b><router-link to="/hands" class="muted tiny">全部记录 →</router-link></div>
        <div v-for="h in hands" :key="h.id" class="hand-row">
          <span class="tag gray mono">#{{ h.handNo }}</span>
          <span class="mono muted">{{ h.roomNo }}</span>
          <span class="mono">{{ h.board || '翻前结束' }}</span>
          <b class="mono gold" style="margin-left:auto">{{ money(h.pot) }}</b>
        </div>
        <div v-if="!hands.length" class="muted tiny" style="padding:10px 0">暂无牌局</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { api, money, fmtTime } from '@/lib/api'
import { adminLive } from '@/lib/live'

const d = ref<any>({})
const tables = ref<any[]>([])
const hands = ref<any[]>([])
let timer = 0

const peakTime = computed(() => (d.value.peakAt ? fmtTime(d.value.peakAt) : '—'))
const gameName = (g: string) => ({ nlh: '长牌', short6: '短牌', omaha: '奥马哈', niuniu: '牛牛' } as any)[g] || g
const stateName = (s: string) => ({ playing: '进行中', showdown: '结算中', waiting: '等待中', insurance: '保险中', '-': '空闲' } as any)[s] || s
const byGame = computed(() => (d.value.byGame || []) as { game: string; hands: number }[])
const trendList = computed(() => (d.value.trend || []) as { day: string; hands: number; users?: number }[])
const trendUsers = computed(() => Number(d.value.trendUsers || 0))
const trendRecharge = computed(() => Number(d.value.trendRecharge || 0))

const GAME_COLORS: Record<string, string> = { nlh: '#2ee07a', short6: '#7dd3fc', omaha: '#ffd875', niuniu: '#c084fc' }
const gameColor = (g: string) => GAME_COLORS[g] || '#6b7280'
const gamePct = (g: { hands: number }) => {
  const total = byGame.value.reduce((a, b) => a + Number(b.hands || 0), 0) || 1
  return Math.round((Number(g.hands || 0) / total) * 100)
}

const gridLines = [30, 75, 120, 165]
const chartMax = computed(() => {
  const vals = trendList.value.map(t => Number(t.hands || 0))
  return Math.max(4, ...vals)
})
function chartXY(i: number, v: number) {
  const n = Math.max(1, trendList.value.length - 1)
  const x = trendList.value.length <= 1 ? 360 : 30 + (i / n) * 660
  const y = 180 - (Number(v || 0) / chartMax.value) * 140
  return { x, y }
}
const pointList = computed(() => trendList.value.map((t, i) => {
  const { x, y } = chartXY(i, Number(t.hands || 0))
  return { x: +x.toFixed(1), y: +y.toFixed(1), value: Number(t.hands || 0), label: String(t.day || '').slice(5, 10) }
}))
const linePoints = computed(() => pointList.value.map(p => `${p.x},${p.y}`).join(' '))
const areaPoints = computed(() => {
  if (!pointList.value.length) return ''
  const first = pointList.value[0]
  const last = pointList.value[pointList.value.length - 1]
  return `${linePoints.value} ${last.x},180 ${first.x},180`
})
const userPoints = computed(() => {
  const vals = trendList.value.map(t => Number(t.users || 0))
  const max = Math.max(2, ...vals)
  const n = Math.max(1, trendList.value.length - 1)
  return vals.map((v, i) => {
    const x = trendList.value.length <= 1 ? 360 : 30 + (i / n) * 660
    const y = 180 - (v / max) * 140
    return `${x.toFixed(1)},${y.toFixed(1)}`
  }).join(' ')
})

const ICON = {
  online: '<circle cx="12" cy="12" r="9"/><path d="M3 12h18"/>',
  users: '<circle cx="9" cy="8" r="3.4"/><path d="M2.5 19c1-3 3.4-4.6 6.5-4.6s5.5 1.6 6.5 4.6"/>',
  robot: '<rect x="4" y="7" width="16" height="12" rx="3"/><path d="M12 3v4M9 12h.01M15 12h.01"/>',
  table: '<circle cx="12" cy="12" r="8.5"/>',
  hands: '<path d="M6 3h8l4 4v14H6z"/>',
  coin: '<ellipse cx="12" cy="7" rx="7.5" ry="3"/><path d="M4.5 7v10c0 1.7 3.4 3 7.5 3s7.5-1.3 7.5-3V7"/>',
  net: '<path d="M3 17l5-6 4 4 6-8"/><path d="M14 6h6v6"/>'
}
function svgIcon(path: string) {
  return `<svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round">${path}</svg>`
}
function deltaOf(idx: number) {
  const n = trendList.value.length
  if (n < 2 || idx >= n) return null
  const cur = Number((trendList.value[n - 1] as any)[idx === 0 ? 'hands' : 'users'] || 0)
  const prev = Number((trendList.value[n - 2] as any)[idx === 0 ? 'hands' : 'users'] || 0)
  if (prev === 0 && cur === 0) return null
  return { delta: cur - prev }
}
const kpis = computed(() => {
  const dh = deltaOf(0)?.delta ?? null
  const du = deltaOf(1)?.delta ?? null
  return [
    { key: 'online', label: '在线连接', value: d.value.online || 0, unit: 'WebSocket 实时', icon: svgIcon(ICON.online), delta: null, deltaUnit: '' },
    { key: 'users', label: '注册用户', value: d.value.users || 0, unit: '非机器人账号', icon: svgIcon(ICON.users), delta: du, deltaUnit: ' 人' },
    { key: 'robots', label: '机器人', value: d.value.robots || 0, unit: '官方房填充', icon: svgIcon(ICON.robot), delta: null, deltaUnit: '' },
    { key: 'tables', label: '实时牌桌', value: d.value.liveTables || 0, unit: '运行中', icon: svgIcon(ICON.table), delta: null, deltaUnit: '' },
    { key: 'hands', label: '累计牌局', value: d.value.hands || 0, unit: '手', icon: svgIcon(ICON.hands), delta: null, deltaUnit: '' },
    { key: 'todayHands', label: '今日牌局', value: d.value.todayHands || 0, unit: '手', icon: svgIcon(ICON.hands), delta: dh, deltaUnit: ' 手' },
    { key: 'todayUsers', label: '今日新增', value: d.value.todayUsers || 0, unit: '人', icon: svgIcon(ICON.users), delta: null, deltaUnit: '' },
    { key: 'todayRecharge', label: '今日充值', value: money(d.value.todayRecharge), unit: 'USDT', icon: svgIcon(ICON.coin), delta: null, deltaUnit: '' },
    { key: 'totalBalance', label: '用户总资产', value: money(d.value.totalBalance), unit: 'USDT', icon: svgIcon(ICON.coin), delta: null, deltaUnit: '' },
    { key: 'qps', label: '实时下发流量', value: Math.round((d.value.qps || 0) / 1024), unit: 'KB / 秒', icon: svgIcon(ICON.net), delta: null, deltaUnit: '' },
    { key: 'peak', label: '在线峰值', value: d.value.peak || 0, unit: peakTime.value, icon: svgIcon(ICON.online), delta: null, deltaUnit: '' },
    { key: 'sentMB', label: '累计下发', value: (d.value.sentMB || 0).toFixed(1), unit: 'MB', icon: svgIcon(ICON.net), delta: null, deltaUnit: '' }
  ]
})

async function load() {
  try { d.value = await api.get('dashboard') } catch {}
  try { tables.value = await api.get('tables') } catch {}
  try { hands.value = await api.get('hands?limit=8') } catch {}
}

// 实时数据走 WS 推送；最近牌局列表低频兜底刷新
let offStats: (() => void) | null = null
let offTables: (() => void) | null = null
onMounted(() => {
  load()
  offStats = adminLive.on('admin.stats', v => { if (v) d.value = v })
  offTables = adminLive.on('admin.tables', v => { if (Array.isArray(v)) tables.value = v })
  timer = window.setInterval(async () => {
    try { hands.value = await api.get('hands?limit=8') } catch {}
  }, 15000)
})
onBeforeUnmount(() => {
  clearInterval(timer)
  offStats?.()
  offTables?.()
})
</script>

<style scoped>
.col { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.row { display: flex; align-items: center; gap: 8px; }
.kpis { display: grid; grid-template-columns: repeat(6, 1fr); gap: 11px; }
.kpi { display: flex; flex-direction: column; gap: 6px; padding: 13px 14px; position: relative; overflow: hidden; }
.kpi::after { content: ""; position: absolute; right: -30px; top: -30px; width: 92px; height: 92px; border-radius: 50%;
  background: radial-gradient(circle, rgba(46, 224, 122, .16), transparent 68%); pointer-events: none; }
.kpi-top { display: flex; align-items: center; justify-content: space-between; gap: 6px; }
.kpi-ico { width: 26px; height: 26px; border-radius: 8px; display: flex; align-items: center; justify-content: center;
  background: rgba(46, 224, 122, .12); border: 1px solid rgba(46, 224, 122, .24); color: #7dffb0; }
.kpi-val { font-size: 24px; font-weight: 800; color: #eaf8f0; letter-spacing: .2px; font-variant-numeric: tabular-nums; }
.kpi-foot { display: flex; align-items: center; justify-content: space-between; gap: 6px; }
.kpi-foot em { font-style: normal; font-size: 10.5px; color: #6f8d7c; }
.delta { font-size: 10px; font-weight: 700; padding: 1px 6px; border-radius: 7px; }
.delta.up { background: rgba(46, 224, 122, .16); color: #7dffb0; }
.delta.down { background: rgba(255, 122, 122, .16); color: #ffa8a8; }

.chart-grid { display: grid; grid-template-columns: 1.55fr 1fr; gap: 12px; margin-top: 12px; }
.chart-card { display: flex; flex-direction: column; gap: 6px; }
.trend-svg { width: 100%; height: 230px; overflow: visible; }
.trend-svg .grid line { stroke: rgba(60, 110, 84, .35); stroke-width: 1; stroke-dasharray: 3 5; }
.axis-label { fill: #6f8d7c; font-size: 10px; }
.value-label { fill: #b7ffd4; font-size: 10.5px; font-weight: 700; }
.legend-dot { width: 8px; height: 8px; border-radius: 3px; display: inline-block; }
.legend-dot.hands { background: #2ee07a; }
.legend-dot.users { background: #ffd875; }
.chart-foot { display: flex; align-items: center; justify-content: space-between; gap: 8px; margin-top: 2px; }

.dist-card { display: flex; flex-direction: column; gap: 8px; }
.dist-bar { display: flex; height: 12px; border-radius: 7px; overflow: hidden; background: rgba(255, 255, 255, .05); }
.dist-seg { height: 100%; transition: width .3s ease; }
.dist-legend { display: flex; flex-direction: column; gap: 7px; margin-top: 2px; }
.dist-row { display: flex; align-items: center; gap: 8px; font-size: 12px; }
.dist-name { flex: 1; color: #cfe7da; }

.panels { display: grid; grid-template-columns: 1.6fr 1fr; gap: 12px; margin-top: 12px; }
.hand-row { display: flex; align-items: center; gap: 8px; padding: 7px 0; border-bottom: 1px solid #14261b; font-size: 12px; }
.state-pill { display: inline-flex; align-items: center; height: 21px; padding: 0 8px; border-radius: 7px; font-size: 10.5px; font-weight: 700;
  background: rgba(255, 255, 255, .06); color: #9fb8ab; border: 1px solid rgba(120, 180, 150, .2); }
.state-pill.playing { background: rgba(46, 224, 122, .16); color: #7dffb0; border-color: rgba(46, 224, 122, .35); }
.state-pill.showdown { background: rgba(255, 216, 117, .16); color: #ffd875; border-color: rgba(255, 216, 117, .35); }
.state-pill.insurance { background: rgba(125, 211, 252, .14); color: #7dd3fc; border-color: rgba(125, 211, 252, .32); }
.tiny { font-size: 10.5px; }
@media (max-width: 1500px) { .kpis { grid-template-columns: repeat(4, 1fr); } }
@media (max-width: 1200px) { .kpis { grid-template-columns: repeat(2, 1fr); } .panels, .chart-grid { grid-template-columns: 1fr; } }
</style>
