<template>
  <div class="page records records-official">
    <header class="rec-head">
      <div class="col" style="align-items:center;width:100%">
        <b class="page-title">战绩与荣耀榜 <span class="tag tag-gold pro-tag">PRO</span></b>
        <span class="tiny muted">实时手牌复盘 · 密码学存证验牌 · 全网王者排行</span>
      </div>
    </header>

    <div class="rec-tabs">
      <button class="rec-tab" :class="{ active: tab === 'hands' }" @click="tab = 'hands'"><Icon name="doc" :size="15" /> 牌局明细</button>
      <button class="rec-tab" :class="{ active: tab === 'stats' }" @click="tab = 'stats'"><Icon name="chart" :size="15" /> 生涯数据</button>
      <button class="rec-tab" :class="{ active: tab === 'rank' }" @click="tab = 'rank'"><Icon name="trophy" :size="15" /> 全网荣耀榜</button>
    </div>

    <div v-if="filterOpen" class="mask" @click.self="filterOpen = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:12px">
          <b class="h2">筛选与排序</b>
          <button class="close-x" @click="filterOpen = false"><Icon name="close" :size="18" /></button>
        </div>
        <div class="pref-title">盈亏</div>
        <div class="seg-row">
          <button v-for="o in [{ k: 'all', t: '全部' }, { k: 'win', t: '盈利局' }, { k: 'loss', t: '亏损局' }, { k: 'star', t: '收藏' }]" :key="o.k"
                  class="seg-cell" :class="{ on: filter === o.k }" @click="filter = o.k as any">{{ o.t }}</button>
        </div>
        <div class="pref-title" style="margin-top:14px">牌型</div>
        <div class="seg-row wrap">
          <button class="seg-cell" :class="{ on: handType === '' }" @click="handType = ''">全部</button>
          <button v-for="t in handTypes" :key="t.name" class="seg-cell" :class="{ on: handType === t.name }" @click="handType = t.name">{{ t.name }} ({{ t.count }})</button>
        </div>
        <div class="pref-title" style="margin-top:14px">排序</div>
        <div class="seg-row wrap">
          <button v-for="o in sortOptions" :key="o.k" class="seg-cell" :class="{ on: sortKey === o.k }" @click="sortKey = o.k as any">{{ o.t }}</button>
        </div>
        <button class="btn btn-primary btn-block" style="height:46px;margin-top:16px" @click="filterOpen = false">完成</button>
      </div>
    </div>

    <template v-if="tab === 'hands'">
      <div class="game-filter-capsules">
        <button v-for="g in games" :key="g.key" class="capsule-btn" :class="{ active: gameFilter === g.key }"
                @click="gameFilter = g.key">{{ g.capsule }}</button>
        <button class="capsule-btn filter-capsule" :class="{ active: filter === 'loss' || !!handType || sortKey !== 'time' }"
                @click="filterOpen = true">筛选</button>
      </div>

      <div class="row gap-8 filters-row">
        <button class="chip" :class="{ on: filter === 'all' }" @click="filter = 'all'">全部 ({{ hands.length }})</button>
        <button class="chip" :class="{ on: filter === 'win' }" @click="filter = 'win'"><Icon name="trophy" :size="13" /> 盈利局</button>
        <button class="chip" :class="{ on: filter === 'loss' }" @click="filter = 'loss'">亏损局</button>
        <button class="chip" :class="{ on: filter === 'star' }" @click="filter = 'star'"><Icon name="star" :size="13" /> 收藏 ({{ favorites.size }})</button>
        <button class="chip refresh-chip" @click="load"><Icon name="refresh" :size="13" /></button>
      </div>

      <div v-if="!visibleHands.length" class="empty-hands-state">
        <div class="empty-icon-box"><Icon name="doc" :size="30" color="#9fe6c4" /></div>
        <b>暂无牌局对战记录</b>
        <p>进入德州牌桌参与实战对局后，每一手牌均会自动生成可追溯存证与复盘记录。</p>
        <button class="capsule-btn" style="margin-top:6px" @click="router.push('/cash')">前往现金场开局</button>
      </div>

      <div v-else class="hand-list">
        <article v-for="h in visibleHands" :key="h.id" class="hand-card-wrapper">
          <header class="hand-card-header">
            <div class="header-left">
              <span class="room-badge">{{ h.roomName || gameNameOf(h.gameType) }}</span>
              <span class="hand-no-badge">#{{ h.handNo }}局</span>
            </div>
            <div class="header-right">
              <time class="time-label">{{ fmtTime(h.time, true) }}</time>
              <span class="pot-num">底池: <b>{{ money(h.pot) }}</b></span>
            </div>
          </header>
          <div class="hand-card-body">
            <div v-if="holeCardsOf(h).length" class="cards-column my-hole-cards">
              <span class="column-title">我的手牌</span>
              <div class="cards-row">
                <PlayingCard v-for="(c, i) in holeCardsOf(h)" :key="'hole' + i" :card="c" size="sm" />
              </div>
            </div>
            <div v-if="parseBoard(h.board).length" class="cards-column community-board">
              <span class="column-title">公共牌 ({{ parseBoard(h.board).length }})</span>
              <div class="cards-row board-row">
                <PlayingCard v-for="(c, i) in parseBoard(h.board)" :key="'board' + i" :card="c" size="sm" />
              </div>
            </div>
            <div class="cards-column result-column">
              <span v-if="h.handName" class="hand-type-tag">{{ h.handName }}</span>
              <div class="profit-badge" :class="h.net > 0 ? 'is-win' : h.net < 0 ? 'is-loss' : 'is-unknown'">
                {{ h.net > 0 ? '+' : '' }}{{ netOf(h) }}
              </div>
            </div>
          </div>
          <footer class="hand-card-footer">
            <button class="footer-btn replay-btn" title="回放这局牌" @click="open(h)">
              <Icon name="play" :size="13" /><span>牌谱回放</span>
            </button>
            <button class="footer-btn fairness-btn" title="查看随机种子与哈希承诺" @click="router.push(`/records/hand/${h.id}`)">
              <Icon name="shield" :size="13" /><span>{{ h.seedRevealed ? '种子已揭示' : '密码学存证' }}</span>
            </button>
            <button class="footer-btn favorite-btn" :class="{ active: favorites.has(h.id) }" title="收藏该手牌" @click="toggleFav(h.id)">
              <Icon name="star" :size="13" /><span>{{ favorites.has(h.id) ? '已收藏' : '收藏' }}</span>
            </button>
          </footer>
        </article>
        <button v-if="hasMore" class="btn btn-ghost btn-block" style="margin-top:10px" @click="loadMore">加载更多</button>
        <p v-else class="tiny muted" style="text-align:center;margin-top:10px">已显示全部 {{ hands.length }} 手</p>
      </div>
    </template>

    <template v-else-if="tab === 'stats'">
      <div class="career-stats-container">
        <div class="game-filter-capsules">
          <button v-for="g in games" :key="g.key" class="capsule-btn" :class="{ active: statGame === g.key }"
                  @click="statGame = g.key">{{ g.capsule }}</button>
        </div>

        <section class="stats-hero-card">
          <div class="ring-gauge-wrap">
            <svg class="ring-gauge-svg" viewBox="0 0 140 140">
              <circle cx="70" cy="70" r="60" fill="none" stroke="rgba(255,255,255,.07)" stroke-width="10" />
              <circle cx="70" cy="70" r="60" fill="none" stroke="#2ee07a" stroke-width="10" stroke-linecap="round"
                      :stroke-dasharray="377" :stroke-dashoffset="377 - Math.min(377, (Number(stat.vpip || 0) / 100) * 377)"
                      transform="rotate(-90 70 70)" />
              <circle cx="70" cy="70" r="46" fill="none" stroke="rgba(255,255,255,.06)" stroke-width="8" />
              <circle cx="70" cy="70" r="46" fill="none" stroke="#ffd875" stroke-width="8" stroke-linecap="round"
                      :stroke-dasharray="289" :stroke-dashoffset="289 - Math.min(289, (Number(stat.aggr || 0) / 100) * 289)"
                      transform="rotate(-90 70 70)" />
            </svg>
            <div class="ring-center-content">
              <b class="vpip-num">{{ stat.vpip || 0 }}%</b>
              <span>入池率</span>
            </div>
          </div>
          <div class="style-assessment-col">
            <span class="style-badge-tag">{{ styleTag }}</span>
            <p class="style-desc-text">{{ styleDesc }}</p>
            <div class="gauge-legend-row">
              <span class="legend-item outer"><i class="dot-outer"></i>入池率 {{ stat.vpip || 0 }}%</span>
              <span class="legend-item inner"><i class="dot-inner"></i>加注率 {{ stat.aggr || 0 }}%</span>
            </div>
          </div>
        </section>

        <section class="technical-metrics-grid">
          <div class="metric-card"><b class="metric-val">{{ stat.hands || 0 }}</b><span class="metric-sub">总手数</span></div>
          <div class="metric-card"><b class="metric-val text-green">{{ stat.wins || 0 }}</b><span class="metric-sub">获胜</span></div>
          <div class="metric-card"><b class="metric-val text-gold">{{ stat.winRate || 0 }}%</b><span class="metric-sub">总胜率</span></div>
          <div class="metric-card"><b class="metric-val" :class="(stat.avgProfit || 0) >= 0 ? 'text-green' : 'text-red'">{{ money(stat.avgProfit || 0) }}</b><span class="metric-sub">手均盈亏</span></div>
          <div class="metric-card"><b class="metric-val">{{ stat.showdownRate || 0 }}%</b><span class="metric-sub">摊牌胜率 W$SD</span></div>
          <div class="metric-card"><b class="metric-val text-gold">{{ (Number(stat.af || 0)).toFixed(2) }}</b><span class="metric-sub">激进系数 AF</span></div>
        </section>

        <section class="profit-trend-card">
          <div class="trend-card-header">
            <span>近期筹码盈亏走势 (Trend)</span>
            <b class="total-profit-num" :class="(stat.net || 0) >= 0 ? 'green' : 'red'">{{ money(stat.net || 0) }}</b>
          </div>
          <div class="chart-canvas-box">
            <svg class="trend-svg" viewBox="0 0 340 100" preserveAspectRatio="none">
              <polyline :points="trendPoints" fill="none" stroke="#39e3a6" stroke-width="2" stroke-linejoin="round" stroke-linecap="round" />
              <polyline :points="trendArea" fill="rgba(57,227,166,.14)" stroke="none" />
            </svg>
          </div>
          <p class="tiny muted" style="margin:6px 0 0">仅展示服务端完整记录的最近 {{ (stat.trend || []).length }} 手</p>
        </section>
      </div>
    </template>

    <template v-else>
      <div class="leaderboard-view">
        <div class="category-tabs">
          <button v-for="c in rankKinds" :key="c.key" class="category-btn" :class="{ active: rankKind === c.key }"
                  @click="setRankKind(c.key)">
            <div class="cat-icon-box"><Icon :name="c.icon" :size="15" :color="c.color" /></div>
            <span class="cat-label">{{ c.label }}</span>
            <span class="cat-badge">{{ c.badge }}</span>
          </button>
        </div>
        <div class="period-tabs-wrap">
          <div class="period-tabs">
            <button v-for="p in rankPeriods" :key="p.key" class="period-tab" :class="{ active: rankPeriod === p.key }"
                    @click="setRankPeriod(p.key)">{{ p.label }}</button>
          </div>
        </div>

        <div v-if="rankLoading" class="lb-loading-state">榜单加载中…</div>

        <template v-else>
          <section v-if="rankTop.length" class="podium-section">
            <div v-if="rankTop[1]" class="podium-col rank-2">
              <div class="player-avatar-wrap"><Avatar :src="rankTop[1].avatar" :size="54" /></div>
              <div class="podium-info">
                <span class="player-name">{{ rankTop[1].nick }}</span>
                <span class="player-tier">{{ rankValueText(rankTop[1]) }}</span>
              </div>
            </div>
            <div v-if="rankTop[0]" class="podium-col rank-1">
              <div class="crown-icon-wrap gold-crown-wrap"><Icon name="crown" :size="22" color="#ffd875" /></div>
              <div class="player-avatar-wrap champion-wrap"><Avatar :src="rankTop[0].avatar" :size="66" /></div>
              <div class="podium-info">
                <span class="player-name champion-name">{{ rankTop[0].nick }}</span>
                <span class="player-tier champion-tier">{{ rankValueText(rankTop[0]) }}</span>
              </div>
            </div>
            <div v-if="rankTop[2]" class="podium-col rank-3">
              <div class="player-avatar-wrap"><Avatar :src="rankTop[2].avatar" :size="54" /></div>
              <div class="podium-info">
                <span class="player-name">{{ rankTop[2].nick }}</span>
                <span class="player-tier">{{ rankValueText(rankTop[2]) }}</span>
              </div>
            </div>
          </section>

          <div class="rank-list-container">
            <div class="rank-list-header">
              <span class="col-rank">名次</span>
              <span class="col-player">玩家</span>
              <span class="col-tier">战绩</span>
              <span class="col-stat">{{ rankStatLabel }}</span>
            </div>
            <div v-if="!rankRest.length && !rankTop.length" class="empty-list">
              <div class="empty-icon-box"><Icon name="trophy" :size="30" color="#9fc9b6" /></div>
              <b>榜单虚位以待</b>
              <p>打满 3 手即进入排名，赢下的每一分都算数。</p>
            </div>
            <div v-for="r in rankRest" :key="'rk' + r.uid" class="rank-row">
              <span class="col-rank"><span class="rank-badge">{{ r.rank }}</span></span>
              <span class="col-player">
                <Avatar :src="r.avatar" :size="30" />
                <span class="cell-names"><b class="cell-nickname">{{ r.nick }}</b></span>
              </span>
              <span class="col-tier"><span class="tier-pill">{{ r.hands ? r.winRate + '%' : '—' }}</span></span>
              <span class="col-stat" :class="rankKind === 'profit' ? (r.net >= 0 ? 'val-profit' : 'val-loss') : rankKind === 'wealth' ? 'val-wealth' : 'val-winrate'">
                {{ rankValueText(r) }}
              </span>
            </div>
          </div>

          <div v-if="myRank > 0" class="my-rank-dock">
            <div class="my-dock-inner">
              <span class="my-rank-badge">我的名次 {{ myRank }}</span>
              <span class="my-dock-text">与榜首相差 {{ rankGapText }}</span>
            </div>
          </div>
        </template>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Avatar from '@/components/Avatar.vue'
import Icon from '@/components/Icon.vue'
import PlayingCard from '@/components/PlayingCard.vue'
import { socket } from '@/lib/ws'
import { money, fmtTime, gameNameOf } from '@/lib/format'
import { useUI } from '@/stores/ui'
import '@/styles/records-official.css'

const router = useRouter()
const ui = useUI()

const tab = ref<'hands' | 'stats' | 'rank'>('hands')
const games = [
  { key: 'all', label: '全部玩法', capsule: '全部牌局' },
  { key: 'nlh', label: '长牌', capsule: '长牌常规' },
  { key: 'short6', label: '短牌(6+)', capsule: '短牌6+' },
  { key: 'omaha', label: '奥马哈', capsule: '奥马哈' }
]
const gameFilter = ref('all')
const statGame = ref('nlh')
const filter = ref<'all' | 'win' | 'loss' | 'star'>('all')
const sortKey = ref<'time' | 'netDesc' | 'netAsc' | 'potDesc'>('time')
const handType = ref('')
const filterOpen = ref(false)
const sortOptions = [
  { k: 'time', t: '最新优先' },
  { k: 'netDesc', t: '输赢从高到低' },
  { k: 'netAsc', t: '输赢从低到高' },
  { k: 'potDesc', t: '底池从大到小' }
]
const hands = ref<any[]>([])
const offset = ref(0)
const hasMore = ref(false)
const PAGE = 40
const stats = ref<Record<string, any>>({})
const favorites = ref(new Set<number>())
async function syncFavorites() {
  // 服务端为准：列表数据里带 favorited 标记
  const set = new Set<number>()
  for (const h of hands.value) if (h.favorited) set.add(h.id)
  favorites.value = set
}
async function toggleFav(id: number) {
  const on = !favorites.value.has(id)
  try {
    await socket.call('records.favorite', { handId: id, on })
    if (on) favorites.value.add(id)
    else favorites.value.delete(id)
    favorites.value = new Set(favorites.value)
    const row = hands.value.find(h => h.id === id)
    if (row) row.favorited = on
    ui.toast(on ? '已收藏' : '已取消收藏', 'ok')
  } catch (e: any) {
    ui.err(e?.message || '收藏失败')
  }
}
const ranking = ref<any[]>([])
const rankDays = ref(1)
const rankKind = ref<'profit' | 'winrate' | 'wealth'>('profit')
const rankPeriod = ref<'today' | 'week' | 'total'>('today')
const rankLoading = ref(false)
const myRank = ref(0)
const rankKinds = [
  { key: 'profit' as const, label: '盈利榜', badge: '热门', icon: 'chart' as const, color: '#ffb066' },
  { key: 'winrate' as const, label: '胜率大师', badge: '精湛', icon: 'target' as const, color: '#7fd8ff' },
  { key: 'wealth' as const, label: '财富榜', badge: '巨擘', icon: 'trophy' as const, color: '#ffd875' }
]
const rankPeriods = [
  { key: 'today' as const, label: '今日榜' },
  { key: 'week' as const, label: '本周榜' },
  { key: 'total' as const, label: '全服总榜' }
]
const rankTop = computed(() => ranking.value.slice(0, 3))
const rankRest = computed(() => ranking.value.slice(3))
const rankStatLabel = computed(() => (rankKind.value === 'profit' ? '累计盈利' : rankKind.value === 'winrate' ? '真实胜率' : '总筹码资产'))
function rankValueText(r: any) {
  if (!r) return ''
  if (rankKind.value === 'profit') return money(r.net || 0)
  if (rankKind.value === 'winrate') return `${r.winRate || 0}%`
  return money(r.value || r.net || 0)
}
const rankGapText = computed(() => {
  const top = ranking.value[0]
  if (!top) return '—'
  if (rankKind.value === 'winrate') return `${Math.max(0, (top.winRate || 0) - 50)}%`
  return money(Math.max(0, Number(top.net || 0)))
})

const stat = computed(() => stats.value[statGame.value] || {})

function holeCardsOf(h: any): string[] {
  const raw = h?.holeCards
  if (!raw) return []
  const s = String(raw).trim()
  if (!s) return []
  if (s.length % 2 === 0 && /^([2-9tjqkaTJQKA][shdc])+$/.test(s)) return s.match(/../g) || []
  return s.split(/[ ,|]+/).filter(Boolean)
}
function netOf(h: any) {
  const v = Number(h?.net || 0) / 100
  return Math.abs(v).toFixed(2)
}

const styleTag = computed(() => {
  const vpip = Number(stat.value.vpip || 0)
  const aggr = Number(stat.value.aggr || 0)
  const winRate = Number(stat.value.winRate || 0)
  if (!Number(stat.value.hands || 0)) return '暂无实战'
  if (vpip > 30) return aggr > 18 ? '凶悍激进型 (TAG)' : '松凶跟注型 (Loose)'
  if (winRate >= 50) return 'TAG 紧凶技术流'
  return '稳健防守型 (Solid)'
})
const styleDesc = computed(() => {
  const handsN = Number(stat.value.hands || 0)
  if (!handsN) return '尚未在该玩法完成对局，开局后将实时生成专业牌风诊断与走势图'
  const label = games.find(g => g.key === statGame.value)?.label || '长牌'
  return `在${label}牌局中已完成 ${handsN} 手对局，胜率 ${stat.value.winRate || 0}%。`
})
const trendPoints = computed(() => {
  const pts = (stat.value.trend || []) as { cum: number }[]
  if (pts.length < 2) return '0,50 340,50'
  const vals = pts.map(p => Number(p.cum || 0))
  const min = Math.min(...vals, 0)
  const max = Math.max(...vals, 0)
  const span = max - min || 1
  return vals.map((v, i) => {
    const x = (i / (vals.length - 1)) * 340
    const y = 92 - ((v - min) / span) * 84
    return `${x.toFixed(1)},${y.toFixed(1)}`
  }).join(' ')
})
const trendArea = computed(() => `${trendPoints.value} 340,100 0,100`)

const visibleHands = computed(() => {
  let list = hands.value
  if (gameFilter.value !== 'all') list = list.filter(h => h.gameType === gameFilter.value)
  if (filter.value === 'win') list = list.filter(h => h.net > 0)
  if (filter.value === 'loss') list = list.filter(h => h.net < 0)
  if (filter.value === 'star') list = list.filter(h => favorites.value.has(h.id))
  if (handType.value) list = list.filter(h => (h.handName || '') === handType.value)
  const key = sortKey.value
  if (key === 'netDesc') list = [...list].sort((a, b) => (b.net || 0) - (a.net || 0))
  else if (key === 'netAsc') list = [...list].sort((a, b) => (a.net || 0) - (b.net || 0))
  else if (key === 'potDesc') list = [...list].sort((a, b) => (b.pot || 0) - (a.pot || 0))
  return list
})
function trimMoney(cents: number) {
  const v = Number(cents || 0) / 100
  const n = Math.round(v * 100) / 100
  return Number.isInteger(n) ? String(n) : n.toFixed(2)
}
// 按牌型聚合（只列真实出现过的），供下拉筛选
const handTypes = computed(() => {
  const map = new Map<string, number>()
  for (const h of hands.value) {
    const n = h.handName || ''
    if (!n) continue
    map.set(n, (map.get(n) || 0) + 1)
  }
  const order = ['皇家同花顺', '同花顺', '四条', '葫芦', '同花', '顺子', '三条', '两对', '一对', '高牌', '五小牛', '炸弹', '五花牛', '牛牛', '牛九', '牛八']
  return [...map.entries()]
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => {
      const ia = order.indexOf(a.name)
      const ib = order.indexOf(b.name)
      return (ia < 0 ? 99 : ia) - (ib < 0 ? 99 : ib)
    })
})
const summary = computed(() => {
  const st = stats.value || {}
  if (gameFilter.value !== 'all') {
    const g = st[gameFilter.value] || {}
    return { hands: g.hands || 0, net: g.net || 0, winRate: g.winRate || 0 }
  }
  let handsSum = 0, netSum = 0, winW = 0
  for (const g of Object.values(st) as any[]) {
    handsSum += g.hands || 0
    netSum += g.net || 0
    winW += Math.round(((g.winRate || 0) * (g.hands || 0)) / 100)
  }
  return { hands: handsSum, net: netSum, winRate: handsSum ? Math.round((winW / handsSum) * 100) : 0 }
})
const total = computed(() => summary.value.hands)
const netTotal = computed(() => summary.value.net)
const winRate = computed(() => summary.value.winRate)

function tagOf(gt: string) {
  return gt === 'short6' ? 'tag-purple' : gt === 'omaha' ? 'tag-blue' : 'tag-gold'
}
function parseBoard(b: string) {
  const out: string[] = []
  for (let i = 0; i + 1 < (b || '').length; i += 2) out.push(b.slice(i, i + 2))
  return out
}

async function load() {
  try {
    const [h, s] = await Promise.all([
      socket.call<any[]>('records.hands', { limit: PAGE, offset: 0 }),
      socket.call<{ stats: any; recent: any[] }>('records.stats')
    ])
    hands.value = h
    offset.value = h.length
    hasMore.value = h.length >= PAGE
    stats.value = s.stats || {}
    syncFavorites()
    await loadRank()
  } catch (e: any) {
    ui.err(e.message)
  }
}

async function loadMore() {
  try {
    const more = await socket.call<any[]>('records.hands', { limit: PAGE, offset: offset.value })
    if (!more.length) { hasMore.value = false; return }
    hands.value = [...hands.value, ...more]
    offset.value += more.length
    hasMore.value = more.length >= PAGE
    syncFavorites()
  } catch (e: any) {
    ui.err(e.message)
  }
}

async function loadRank() {
  rankLoading.value = true
  try {
    const r = await socket.call<any>('records.rank', { kind: rankKind.value, period: rankPeriod.value, limit: 30 })
    ranking.value = r?.list || []
    myRank.value = Number(r?.myRank || 0)
  } catch {
    ranking.value = []
    myRank.value = 0
  } finally {
    rankLoading.value = false
  }
}

function setRankKind(k: 'profit' | 'winrate' | 'wealth') {
  rankKind.value = k
  loadRank()
}
function setRankPeriod(p: 'today' | 'week' | 'total') {
  rankPeriod.value = p
  loadRank()
}

function open(h: any) {
  router.push(`/records/hand/${h.id}`)
}

onMounted(load)
</script>

<style scoped>
.fav-btn {
  width: 26px; height: 26px; border-radius: 8px; display: inline-flex; align-items: center; justify-content: center;
  color: rgba(210, 235, 220, .4); background: rgba(255, 255, 255, .05); border: 1px solid rgba(122, 196, 156, .16);
}
.fav-btn.on { color: #ffd875; border-color: rgba(255, 216, 117, .5); background: rgba(255, 216, 117, .12); }
.records { padding-top: 14px; }
.rec-head { display: flex; align-items: center; gap: 8px; }
.head-btn { height: 34px; padding: 0 12px; border-radius: 10px; background: rgba(255,255,255,.06); border: 1px solid var(--line); font-size: 13px; display: flex; align-items: center; gap: 4px; }
.summary { display: flex; align-items: center; padding: 12px 10px; border-radius: 14px; }
.sum-item { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 4px; }
.sum-item b { font-size: 18px; font-weight: 800; }
.divider { width: 1px; height: 26px; background: var(--line); }
.chip { height: 30px; padding: 0 11px; border-radius: 16px; font-size: 12px; color: var(--muted); background: rgba(10, 30, 20, .7); border: 1px solid var(--line); display: inline-flex; align-items: center; gap: 4px; }
.chip.on { color: #241a06; background: linear-gradient(180deg, #f2d493, #cfa151); border-color: transparent; font-weight: 700; }
.empty { padding: 30px 16px; display: flex; flex-direction: column; align-items: center; }
.page-title { font-size: 22px; font-weight: 800; }
.pro-tag { font-size: 10px; height: 18px; }
.doc-box { width: 104px; height: 104px; border-radius: 28px; display: flex; align-items: center; justify-content: center; background: rgba(20, 50, 36, .7); border: 1px solid var(--line); margin-bottom: 10px; }
.hand-list { display: flex; flex-direction: column; gap: 7px; }
.hand-item { padding: 10px 12px; border-radius: 13px; }
.hand-item:active { transform: scale(.99); }
.stats-card { display: grid; grid-template-columns: 1fr 1.2fr 1fr; grid-template-rows: auto auto; gap: 10px; align-items: center; }
.stat-cell { display: flex; flex-direction: column; align-items: center; gap: 2px; }
.stat-cell b { font-size: 20px; font-weight: 800; }
.stat-cell span { font-size: 11px; color: var(--muted); }
.radial { position: relative; display: flex; align-items: center; justify-content: center; }
.radial svg { width: 100%; max-width: 130px; }
.radial-txt { position: absolute; display: flex; flex-direction: column; align-items: center; }
.radial-txt b { font-size: 19px; color: #7ff0b0; }
.radial-txt span { font-size: 10px; color: var(--muted); }
.stat-foot { grid-column: 1 / 4; display: flex; justify-content: space-around; padding-top: 10px; border-top: 1px solid var(--line); }
.stat-foot div { display: flex; flex-direction: column; align-items: center; }
.stat-foot b { font-size: 16px; }
.stat-foot span { font-size: 11px; color: var(--muted); }
.rank-hero { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.rank-list { display: flex; flex-direction: column; gap: 8px; }
.rank-row { display: flex; align-items: center; gap: 10px; padding: 10px 12px; border-radius: 12px; background: rgba(12, 34, 22, .7); border: 1px solid var(--line); }
.rank-no { width: 22px; text-align: center; font-weight: 900; color: var(--muted); }
.rank-no.r1 { color: #ffd98a; } .rank-no.r2 { color: #d8e6ff; } .rank-no.r3 { color: #ffb07a; }

.chip-sel {
  height: 28px; padding: 0 8px; border-radius: 9px; font-size: 11.5px; color: #cfe6d9;
  background: rgba(255, 255, 255, .06); border: 1px solid var(--line);
}

.filters-row { flex-wrap: wrap; row-gap: 6px; margin-bottom: 12px; }
.filters-row .chip { flex: none; white-space: nowrap; }
.filters-row .chip-sel { flex: none; }
.refresh-chip { margin-left: auto; }

/* ===== 官方战绩页样式 ===== */
.rec-tabs {
  display: flex; gap: 4px; padding: 4px; margin: 14px 0 12px; border-radius: 16px;
  background: rgba(8, 26, 18, .9); border: 1px solid rgba(214, 176, 92, .22);
}
.rec-tab {
  flex: 1; height: 44px; border-radius: 13px; display: inline-flex; align-items: center; justify-content: center; gap: 5px;
  font-size: 14px; font-weight: 700; color: #d9ecdf; white-space: nowrap;
}
.rec-tab.active {
  color: #ffe08a; font-weight: 900;
  background: rgba(30, 62, 44, .9);
  box-shadow: inset 0 0 0 1.5px rgba(255, 216, 117, .85);
}
.game-tabs {
  display: flex; gap: 4px; padding: 5px; margin: 12px 0; border-radius: 14px;
  background: rgba(8, 26, 18, .9); border: 1px solid rgba(255, 255, 255, .06);
}
.game-tab {
  flex: 1; height: 36px; border-radius: 11px; font-size: 13.5px; font-weight: 700; color: #cfe6d9; white-space: nowrap;
}
.game-tab.active {
  color: #23180a; font-weight: 900;
  background: linear-gradient(180deg, #ffe9a8, #e7bf5c);
  box-shadow: 0 4px 12px rgba(231, 191, 92, .28);
}
.filters-row .chip {
  height: 32px; padding: 0 12px; border-radius: 16px; font-size: 12.5px; font-weight: 700;
  color: #cfe6d9; background: rgba(255, 255, 255, .04); border: 1px solid rgba(255, 255, 255, .1);
  display: inline-flex; align-items: center; gap: 4px;
}
.filters-row .chip.on {
  color: #ffe08a; border-color: rgba(255, 216, 117, .8); background: rgba(48, 38, 12, .55);
}
.more-chip { margin-left: auto; width: 34px; padding: 0; justify-content: center; }
.refresh-chip { margin-left: 0; }
.filters-row { flex-wrap: nowrap; overflow: hidden; gap: 6px !important; }
.filters-row .chip { flex: none; padding: 0 10px; font-size: 12px; }
.pref-title { font-size: 11px; font-weight: 800; letter-spacing: .8px; color: #8fb3a1; margin-bottom: 8px; }
.seg-row { display: flex; gap: 8px; }
.seg-row.wrap { flex-wrap: wrap; }
.seg-cell {
  flex: 1; min-width: 88px; height: 40px; border-radius: 11px; font-size: 13px; font-weight: 700;
  color: #cfe6d9; background: rgba(255, 255, 255, .04); border: 1px solid rgba(255, 255, 255, .1);
}
.seg-cell.on { color: #23180a; background: linear-gradient(180deg, #ffe9a8, #e7bf5c); border-color: transparent; }

.doc-box {
  width: 64px; height: 64px; border-radius: 16px; display: flex; align-items: center; justify-content: center;
  background: rgba(20, 50, 36, .55); border: 1px solid rgba(126, 240, 180, .35);
}
</style>
