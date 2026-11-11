<template>
  <div class="ledger-page">
    <div class="toolbar card">
      <div class="row">
        <select v-model="type" class="input" style="width:130px">
          <option value="">全部类型</option>
          <option value="recharge">充值</option><option value="withdraw">提现</option>
          <option value="buyin">带入</option><option value="cashout">带出</option>
          <option value="card">房卡</option><option value="adjust">后台调整</option>
          <option value="insurance">保险</option>
        </select>
        <select v-model.number="range" class="input" style="width:120px">
          <option :value="0">全部时间</option><option :value="1">近 1 小时</option>
          <option :value="24">近 24 小时</option><option :value="168">近 7 天</option>
        </select>
        <input v-model="q" class="input" placeholder="昵称 / UID" style="width:150px" />
        <button class="btn btn-sm" @click="exportCsv">导出 CSV</button>
      </div>
      <div class="row">
        <span class="sum-chip">记录 <b>{{ filtered.length }}</b></span>
        <span class="sum-chip green">充值 <b class="mono">{{ money(sumIn) }}</b></span>
        <span class="sum-chip red">提现 <b class="mono">{{ money(sumOut) }}</b></span>
        <span class="sum-chip gold">净流入 <b class="mono">{{ money(net) }}</b></span>
        <span class="sum-chip">房卡 <b>{{ cardCount }} 张</b></span>
      </div>
    </div>

    <div class="card panel">
      <div class="panel-head">
        <b>财务流水</b>
        <span class="muted tiny">按类型汇总（当前筛选结果）</span>
      </div>
      <div style="padding: 12px 14px 4px">
        <div class="dist-bar">
          <span v-for="d in typeDist" :key="d.key" class="dist-seg" :style="{ width: d.pct + '%', background: d.color }"></span>
        </div>
        <div class="dist-legend">
          <span v-for="d in typeDist" :key="'l' + d.key" class="dist-item">
            <i :style="{ background: d.color }"></i>{{ d.label }} <b>{{ d.count }}</b>
          </span>
          <span v-if="!typeDist.length" class="muted tiny">暂无流水</span>
        </div>
      </div>
      <div class="panel-body">
        <table class="table">
          <thead><tr><th>ID</th><th>UID</th><th>昵称</th><th>类型</th><th>金额</th><th>变动后</th><th>备注</th><th>时间</th></tr></thead>
          <tbody>
            <tr v-for="b in filtered" :key="b.id">
              <td class="mono">{{ b.id }}</td>
              <td class="mono">{{ b.uid }}</td>
              <td>{{ b.nick }}</td>
              <td><span class="tag" :class="typeTag(b.type)">{{ typeName(b.type) }}</span></td>
              <td class="mono" :class="b.amount >= 0 ? 'green' : 'red'">{{ b.type === 'card' ? Math.abs(b.amount) + ' 张' : money(b.amount) }}</td>
              <td class="mono">{{ b.type === 'card' ? '—' : money(b.after) }}</td>
              <td class="muted">{{ b.remark }}</td>
              <td class="muted">{{ fmtTime(b.time) }}</td>
            </tr>
            <tr v-if="!filtered.length"><td colspan="8" class="muted" style="text-align:center;padding:20px">没有符合条件的流水</td></tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, money, fmtTime } from '@/lib/api'

const list = ref<any[]>([])
const type = ref('')
const range = ref(0)
const q = ref('')
const filtered = computed(() => {
  let rows = list.value
  if (type.value) rows = rows.filter(x => x.type === type.value)
  if (range.value) {
    const cut = Date.now() - range.value * 3600 * 1000
    rows = rows.filter(x => (x.time || 0) >= cut)
  }
  const k = q.value.trim()
  if (k) rows = rows.filter(x => String(x.nick || '').includes(k) || String(x.uid || '').includes(k))
  return rows
})
const sumIn = computed(() => filtered.value.filter(x => x.type === 'recharge').reduce((a, b) => a + (b.amount || 0), 0))
const sumOut = computed(() => Math.abs(filtered.value.filter(x => x.type === 'withdraw').reduce((a, b) => a + (b.amount || 0), 0)))
const net = computed(() => sumIn.value - sumOut.value)
const cardCount = computed(() => Math.abs(filtered.value.filter(x => x.type === 'card').reduce((a, b) => a + (b.amount || 0), 0)))
const TYPE_COLORS: Record<string, string> = {
  recharge: '#2ee07a', withdraw: '#ff7a7a', buyin: '#7dd3fc', cashout: '#ffd875',
  card: '#c084fc', adjust: '#94a3b8', insurance: '#34d399', muck_peek: '#f472b6', club: '#60a5fa'
}
const typeTag = (t: string) => ({ recharge: '', withdraw: 'red', buyin: 'blue', cashout: 'orange', card: 'purple', adjust: 'gray', insurance: '', muck_peek: 'purple', club: 'blue' } as any)[t] || ''
const typeDist = computed(() => {
  const map = new Map<string, number>()
  for (const b of filtered.value) map.set(b.type, (map.get(b.type) || 0) + 1)
  const total = filtered.value.length || 1
  return [...map.entries()]
    .sort((a, b) => b[1] - a[1])
    .map(([key, count]) => ({
      key, count, pct: Math.max(2, Math.round((count / total) * 100)),
      label: typeName(key), color: TYPE_COLORS[key] || '#6b7280'
    }))
})
const typeName = (t: string) => ({ recharge: '充值', withdraw: '提现', buyin: '带入', cashout: '带出', card: '房卡', adjust: '后台调整', insurance: '保险', muck_peek: '盖牌窥探', club: '俱乐部房卡' } as any)[t] || t

function exportCsv() {
  const rows = [['id', 'uid', 'nick', 'type', 'amount', 'after', 'remark', 'time']]
  for (const b of filtered.value) {
    rows.push([b.id, b.uid, b.nick, b.type, (b.amount / 100).toFixed(2), (b.after / 100).toFixed(2), b.remark || '', new Date(b.time).toISOString()])
  }
  const csv = rows.map(r => r.map(v => String(v ?? '').replace(/,/g, ' ')).join(',')).join('\n')
  const blob = new Blob([new Uint8Array([0xEF, 0xBB, 0xBF]), csv], { type: 'text/csv;charset=utf-8' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `ledger-${Date.now()}.csv`
  a.click()
  URL.revokeObjectURL(a.href)
}

async function load() {
  try { list.value = await api.get<any[]>('ledger?limit=500') } catch {}
}
onMounted(load)
</script>
<style scoped>
.ledger-page { display: flex; flex-direction: column; gap: 12px; }
.green { color: #7ff0b0; } .red { color: #ff9f9f; } .gold { color: #ffd875; }
.sum-row { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; }
.sum-cell { display: flex; flex-direction: column; gap: 2px; padding: 9px 10px; border-radius: 8px; background: #0a1710; border: 1px solid #152a1e; }
.sum-cell b { font-size: 17px; font-weight: 800; }
</style>
