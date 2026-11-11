<template>
  <div class="hands-page">
    <div class="toolbar card">
      <div class="row">
        <input v-model="room" class="input" placeholder="房间号" style="width:100px" />
        <input v-model="uid" class="input" placeholder="玩家 UID" style="width:115px" />
        <select v-model="game" class="input" style="width:110px">
          <option value="">全部玩法</option><option value="nlh">长牌</option>
          <option value="short6">短牌</option><option value="omaha">奥马哈</option><option value="niuniu">牛牛</option>
        </select>
        <input v-model="winner" class="input" placeholder="赢家 / 牌型 / UID" style="width:135px" />
        <select v-model.number="range" class="input" style="width:110px">
          <option :value="0">全部时间</option><option :value="1">近 1 小时</option>
          <option :value="24">近 24 小时</option><option :value="168">近 7 天</option>
        </select>
        <button class="btn btn-sm" @click="exportCsv">导出 CSV</button>
      </div>
      <div class="row">
        <span class="sum-chip">显示 <b>{{ filtered.length }}</b> / {{ hands.length }}</span>
        <span class="sum-chip gold">底池合计 <b class="mono">{{ money(totalPot) }}</b></span>
        <span class="sum-chip green">摊牌 <b>{{ showdownCount }}</b></span>
      </div>
    </div>

    <div class="card panel">
      <div class="panel-head">
        <b>牌局记录</b>
        <span class="muted tiny">点击任意行展开行动线</span>
      </div>
      <div class="panel-body">
        <table class="table">
          <thead><tr><th>ID</th><th>手数</th><th>房间</th><th>玩法</th><th>公共牌</th><th>底池</th><th>赢家 / 牌型</th><th>时间</th><th style="width:150px">操作</th></tr></thead>
          <tbody>
            <template v-for="h in filtered" :key="h.id">
            <tr @click="toggleRow(h)" style="cursor:pointer">
              <td class="mono">{{ h.id }}</td>
              <td class="mono">#{{ h.handNo }}</td>
              <td class="mono">{{ h.roomNo }}</td>
              <td><span class="tag" :class="gameTag(h.gameType)">{{ gameName(h.gameType) }}</span></td>
              <td class="mono"><span v-for="(c, i) in boardCards(h.board)" :key="i" class="bcard" :class="c.red ? 'red' : ''">{{ c.txt }}</span><span v-if="!h.board" class="muted">—</span></td>
              <td class="mono gold">{{ money(h.pot) }}</td>
              <td>
                <span :class="(h.winner || '').includes(' / ') ? '' : 'green'">{{ h.winner || '—' }}</span>
                <span v-if="h.handName" class="tag" style="margin-left:6px">{{ h.handName }}</span>
              </td>
              <td class="muted">{{ fmtTime(h.time) }}</td>
              <td>
                <button class="btn btn-sm" @click.stop="toggleRow(h)">{{ expanded === h.id ? '收起' : '展开' }}</button>
                <button class="btn btn-sm" style="margin-left:6px" @click.stop="open(h)">详情</button>
              </td>
            </tr>
            <tr v-if="expanded === h.id" class="expand-row">
              <td colspan="9">
                <div class="exp-wrap">
                  <div class="exp-line">
                    <b class="tiny">行动线</b>
                    <span class="muted tiny">{{ (h.actions || []).length }} 条</span>
                  </div>
                  <div class="acts">
                    <span v-for="a in h.actions || []" :key="a.seq" class="act">
                      <b class="mono">S{{ a.street }}</b> 座{{ a.seat }} {{ actName(a.action) }} {{ a.amount ? money(a.amount) : '' }}
                    </span>
                    <span v-if="!(h.actions || []).length" class="muted tiny">本手无行动记录</span>
                  </div>
                </div>
              </td>
            </tr>
            </template>
            <tr v-if="!filtered.length"><td colspan="9" class="muted" style="text-align:center;padding:20px">没有符合条件的牌局</td></tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, money, fmtTime } from '@/lib/api'

const hands = ref<any[]>([])
const room = ref('')
const uid = ref('')
const game = ref('')
const winner = ref('')
const range = ref(0)
const detail = ref<any>(null)
const expanded = ref<number>(0)
// "KhTs6dJdQh" → [{txt:'K♥',red:true},...]
function boardCards(board?: string) {
  if (!board) return []
  const map: Record<string, string> = { h: '♥', d: '♦', c: '♣', s: '♠' }
  const out: { txt: string; red: boolean }[] = []
  const str = String(board)
  for (let i = 0; i + 1 < str.length; i += 2) {
    const rank = str[i]
    const suit = map[str[i + 1]] || str[i + 1]
    out.push({ txt: rank + suit, red: str[i + 1] === 'h' || str[i + 1] === 'd' })
  }
  return out
}
const totalPot = computed(() => filtered.value.reduce((a, b) => a + Number(b.pot || 0), 0))
const showdownCount = computed(() => filtered.value.filter((h: any) => !!h.handName).length)
const gameName = (g: string) => ({ nlh: '长牌', short6: '短牌', omaha: '奥马哈', niuniu: '牛牛' } as any)[g] || g
const gameTag = (g: string) => ({ nlh: '', short6: 'purple', omaha: 'blue', niuniu: 'orange' } as any)[g] || ''
const actName = (a: string) => ({ fold: '弃牌', check: '过牌', call: '跟注', bet: '下注', raise: '加注', allin: '全下', 'all-in': '全下', post: '盲注' } as any)[a] || a

// 行内展开需要行动线：按需拉取并回填到行对象
const actCache = new Map<number, any>()
async function toggleRow(h: any) {
  if (expanded.value === h.id) { expanded.value = 0; return }
  expanded.value = h.id
  if (!h.actions && !actCache.has(h.id)) {
    try {
      const d = await api.get<any>('hand?id=' + h.id)
      h.actions = d?.actions || []
      h.players = d?.players || h.players || []
      actCache.set(h.id, d)
    } catch { h.actions = [] }
  } else if (actCache.has(h.id)) {
    h.actions = actCache.get(h.id)?.actions || []
  }
}

const filtered = computed(() => {
  const k = room.value.trim()
  const u = uid.value.trim()
  let rows = hands.value
  if (k) rows = rows.filter(h => String(h.roomNo).includes(k))
  if (game.value) rows = rows.filter(h => h.gameType === game.value)
  const w = winner.value.trim()
  if (w) {
    rows = rows.filter(h => String(h.winner || '').includes(w) || String(h.handName || '').includes(w) ||
      (h.winnerUids || []).some((u: number) => String(u).includes(w)))
  }
  if (u) rows = rows.filter(h => (h.players || []).some((p: any) => String(p.uid).includes(u)))
  if (range.value > 0) {
    const cut = Date.now() - range.value * 3600 * 1000
    rows = rows.filter(h => (h.time || 0) >= cut)
  }
  return rows
})

function exportCsv() {
  const rows = [['id', 'handNo', 'roomNo', 'gameType', 'board', 'pot', 'winner', 'handName', 'time']]
  for (const h of filtered.value) {
    rows.push([h.id, h.handNo, h.roomNo, h.gameType, h.board || '', (h.pot / 100).toFixed(2), h.winner || '', h.handName || '', new Date(h.time).toISOString()])
  }
  const csv = rows.map(r => r.join(',')).join('\n')
  const blob = new Blob([new Uint8Array([0xEF, 0xBB, 0xBF]), csv], { type: 'text/csv;charset=utf-8' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `hands-${Date.now()}.csv`
  a.click()
  URL.revokeObjectURL(a.href)
}

async function load() {
  try { hands.value = await api.get<any[]>('hands?limit=100') } catch {}
}

async function open(h: any) {
  try { detail.value = await api.get('hand?id=' + h.id) } catch (e: any) {
    window.dispatchEvent(new CustomEvent('admin-toast', { detail: { text: e.message, kind: 'err' } }))
  }
}

onMounted(load)
</script>

<style scoped>
.hands-page { display: flex; flex-direction: column; gap: 12px; }
.exp-wrap { padding: 6px 2px 10px; display: flex; flex-direction: column; gap: 6px; }
.exp-line { display: flex; align-items: center; gap: 8px; }
.acts { display: flex; flex-wrap: wrap; gap: 6px; }
.act { font-size: 11px; color: #cfe7da; padding: 3px 8px; border-radius: 7px; background: rgba(255, 255, 255, .04); border: 1px solid rgba(60, 120, 88, .28); }
.act b { color: #7dd3fc; margin-right: 4px; }
.expand-row td { background: rgba(8, 20, 14, .72); }
.bcard { display: inline-block; font-size: 11px; padding: 1px 4px; margin-right: 3px; border-radius: 4px; background: rgba(255, 255, 255, .07); }
.bcard.red { color: #ff9f9f; }
.green { color: #7ff0b0; }
.red { color: #ff9f9f; }
.acts { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 8px; }
.act { font-size: 11px; padding: 3px 7px; border-radius: 6px; background: #0a1710; border: 1px solid #152a1e; }

.expand-row td { background: #081410; }
.exp-wrap { padding: 6px 2px 10px; }
.exp-line { display: flex; align-items: baseline; gap: 8px; }
.exp-players { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 6px; }
.exp-player {
  display: inline-flex; align-items: center; gap: 5px; font-size: 11px; padding: 4px 8px;
  border-radius: 7px; background: #0a1710; border: 1px solid #152a1e; color: #cfe6d9;
}
.exp-player.win { border-color: rgba(127, 240, 176, .4); }
.exp-player i { font-style: normal; color: #ffd875; }

.bcard {
  display: inline-block; min-width: 22px; text-align: center; margin-right: 3px; padding: 1px 3px;
  border-radius: 4px; background: #0d1f16; border: 1px solid #1d3a2a; font-size: 11px; color: #dfeee6;
}
.bcard.red { color: #ff9f9f; border-color: #3a2026; }
</style>
