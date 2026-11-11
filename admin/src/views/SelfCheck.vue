<template>
  <div class="selfcheck-page">
    <div class="toolbar card">
      <div class="row">
        <span class="sum-chip" :class="data && data.dbOk ? 'green' : 'red'">DB <b>{{ data ? (data.dbOk ? '正常' : '异常') : '—' }}</b> · {{ data ? data.dbMs : 0 }}ms</span>
        <span class="sum-chip">运行 <b>{{ uptime }}</b></span>
        <span class="sum-chip">Goroutine <b>{{ data ? data.goroutines : 0 }}</b></span>
        <span class="sum-chip">堆内存 <b>{{ data ? data.heapMB : 0 }} MB</b></span>
      </div>
      <div class="row">
        <button class="btn btn-sm" @click="copyCmd">{{ copied ? '已复制' : '复制验收命令' }}</button>
        <button class="btn btn-sm" @click="load">刷新</button>
      </div>
    </div>

    <div v-if="!data" class="card muted tiny" style="padding:16px">读取中…</div>
    <div v-else class="sc-kpis">
      <div v-for="k in kpis" :key="k.label" class="card sc-kpi">
        <span class="muted tiny">{{ k.label }}</span>
        <b :class="k.warn ? 'red' : ''">{{ k.value }}</b>
        <em v-if="k.sub">{{ k.sub }}</em>
      </div>
    </div>

    <div class="card panel">
      <div class="panel-head">
        <b>最近验收报告</b>
        <div class="row gap-8">
          <select v-if="history.length" v-model="picked" class="input" style="width:190px" @change="loadReport">
            <option value="">最新一次</option>
            <option v-for="h in history" :key="h.file" :value="h.file">{{ h.file.replace('.md','') }} · {{ (h.size/1024).toFixed(1) }}KB</option>
          </select>
          <span class="muted tiny" v-if="report.exists">更新于 {{ fmtTime(report.mtime) }}</span>
        </div>
      </div>
      <div style="padding: 12px 14px">
        <div v-if="!report.exists" class="muted tiny">
          还没有报告，终端执行 <code>make verify-suite</code> 生成
        </div>
        <template v-else>
          <div class="summary-line" v-html="headline"></div>
          <table class="table" style="margin-top:8px">
            <thead><tr><th>#</th><th>套件</th><th>断言</th><th>结果</th><th>耗时</th></tr></thead>
            <tbody>
              <tr v-for="r in rows" :key="r.idx">
                <td class="mono">{{ r.idx }}</td>
                <td>{{ r.name }}</td>
                <td class="mono">{{ r.claims }}</td>
                <td><span class="tag" :class="r.ok ? '' : 'red'">{{ r.ok ? '通过' : '失败' }}</span></td>
                <td class="muted">{{ r.cost }}</td>
              </tr>
            </tbody>
          </table>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, fmtTime } from '@/lib/api'

const data = ref<any>(null)
const report = ref<{ exists: boolean; content: string; mtime: number }>({ exists: false, content: '', mtime: 0 })
const history = ref<any[]>([])
const picked = ref('')
const copied = ref(false)
async function copyCmd() {
  try {
    await navigator.clipboard.writeText('make verify-suite')
    copied.value = true
    setTimeout(() => { copied.value = false }, 1800)
  } catch { copied.value = false }
}
async function loadReport() {
  try {
    const q = picked.value ? `?file=${encodeURIComponent(picked.value)}` : ''
    report.value = await api.get('verify-report' + q)
  } catch {
    report.value = { exists: false, content: '', mtime: 0 }
  }
}

const kpis = computed(() => {
  const d = data.value || {}
  return [
    { label: '在线 / 峰值', value: `${d.online || 0} / ${d.peak || 0}` },
    { label: '牌桌', value: d.tables || 0 },
    { label: '俱乐部', value: d.clubs || 0 },
    { label: '真人 / 机器人', value: `${d.users || 0} / ${d.robots || 0}` },
    { label: '近 24h 手数', value: d.hands24h || 0 },
    { label: '近 24h 流水', value: d.bills24h || 0 },
    { label: '待处理举报', value: d.pendingReports || 0, warn: !!d.pendingReports },
    { label: '未读通知', value: d.unreadNotices || 0 },
    { label: '堆内存', value: `${d.heapMB || 0} MB` },
    { label: 'GC 次数', value: d.gcRuns || 0, sub: `运行 ${uptime.value}` }
  ]
})

const uptime = computed(() => {
  const s = Number(data.value?.uptimeSec || 0)
  const h = Math.floor(s / 3600), m = Math.floor((s % 3600) / 60)
  return h > 0 ? `${h}h ${m}m` : `${m}m ${s % 60}s`
})

// 从 markdown 里抽总断言数与表格行
const headline = computed(() => {
  const m = report.value.content.match(/套件：(.+?)　断言：(.+?)　耗时：(.+)/)
  if (!m) return ''
  const claims = m[2].replace(/\*/g, '')
  return `<b>${m[1].replace(/\*/g, '')}</b> · 断言 <b class="green">${claims}</b> · 总耗时 ${m[3]}`
})
const rows = computed(() => {
  const out: any[] = []
  for (const line of report.value.content.split('\n')) {
    const m = line.match(/^\|\s*(\d+)\s*\|\s*(.+?)（`(.+?)`）\s*\|\s*([\d?\/]+)\s*\|\s*(✅|❌)\s*\|\s*(.+?)\s*\|/)
    if (m) out.push({ idx: m[1], name: m[2], file: m[3], claims: m[4], ok: m[5] === '✅', cost: m[6] })
  }
  return out
})

async function load() {
  try { data.value = await api.get('selfcheck') } catch { data.value = null }
  try { history.value = await api.get('verify-report-list') } catch { history.value = [] }
  await loadReport()
}
onMounted(load)
</script>

<style scoped>
.selfcheck-page { display: flex; flex-direction: column; gap: 12px; }
.sc-kpis { display: grid; grid-template-columns: repeat(5, 1fr); gap: 11px; }
.sc-kpi { display: flex; flex-direction: column; gap: 4px; padding: 13px 14px; }
.sc-kpi b { font-size: 20px; font-weight: 800; color: #eaf8f0; font-variant-numeric: tabular-nums; }
.sc-kpi b.red { color: #ffa8a8; }
.sc-kpi em { font-style: normal; font-size: 10.5px; color: #6f8d7c; }
@media (max-width: 1400px) { .sc-kpis { grid-template-columns: repeat(3, 1fr); } }
.kpis { display: grid; grid-template-columns: repeat(6, 1fr); gap: 8px; margin-top: 10px; }
.kpi { display: flex; flex-direction: column; gap: 2px; padding: 9px 10px; border-radius: 8px; background: #0a1710; border: 1px solid #152a1e; }
.kpi b { font-size: 17px; font-weight: 800; color: #eaf8f0; }
.kpi i { font-style: normal; font-size: 11px; font-weight: 500; }
.green { color: #7ff0b0; }
.red { color: #ff9f9f; }
.summary-line { font-size: 12.5px; color: #a9c4b6; }
</style>
