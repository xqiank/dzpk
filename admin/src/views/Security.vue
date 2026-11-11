<template>
  <div class="security-page">
    <div class="toolbar card">
      <div class="row">
        <b>玩家举报队列</b>
        <span class="sum-chip">当前 <b>{{ reports.length }}</b> 条</span>
        <span class="sum-chip" :class="pendingCount ? 'red' : 'green'">待处理 <b>{{ pendingCount }}</b></span>
      </div>
      <div class="row">
        <select v-model.number="rpStatus" class="input" style="width:120px" @change="loadReports">
          <option :value="0">待处理</option><option :value="1">已处理</option><option :value="-1">全部</option>
        </select>
        <button class="btn btn-sm" @click="loadReports">刷新</button>
      </div>
    </div>

    <div class="card panel"><div class="panel-head"><b>举报明细</b><span class="muted tiny">处理后会写入安全审计日志</span></div><div class="panel-body">
<table class="table" style="margin-top:10px">
      <thead><tr><th>ID</th><th>举报人</th><th>被举报人</th><th>房间</th><th>原因</th><th>说明</th><th>时间</th><th>操作</th></tr></thead>
      <tbody>
        <tr v-for="r in reports" :key="r.id">
          <td class="mono">{{ r.id }}</td>
          <td>{{ r.nick || r.uid }} <span class="muted mono">#{{ r.uid }}</span></td>
          <td>{{ r.targetNick || r.targetUid }} <span class="muted mono">#{{ r.targetUid }}</span></td>
          <td class="mono">{{ r.roomNo }}</td>
          <td><span class="tag" :class="reasonTag(r.reason)">{{ reasonName(r.reason) }}</span></td>
          <td class="muted">{{ r.detail || '-' }}</td>
          <td class="muted">{{ fmtTime(r.createdAt) }}</td>
          <td>
            <span v-if="r.status" class="tag gray">已处理</span>
            <button v-else class="btn btn-sm btn-primary" @click="resolve(r)">标记处理</button>
          </td>
        </tr>
        <tr v-if="!reports.length"><td colspan="8" class="muted" style="text-align:center;padding:14px">暂无记录</td></tr>
      </tbody>
    </table>
    </div></div>
  <div class="grid" style="grid-template-columns: 1fr 1fr">
    <div class="card">
      <div class="between"><b>安全事件</b>
        <select v-model="kind" class="input" style="width:120px">
          <option value="">全部类型</option><option value="rate_limit">限流</option>
          <option value="sql_inject">注入探测</option><option value="admin_try">后台爆破</option>
          <option value="ws_abuse">连接滥用</option><option value="bad_token">令牌异常</option>
        </select>
        <select v-model.number="hours" class="input" style="width:112px">
          <option :value="0">全部</option><option :value="1">近 1 小时</option>
          <option :value="24">近 24 小时</option><option :value="168">近 7 天</option>
        </select>
      </div>
      <table class="table" style="margin-top:10px">
        <thead><tr><th>类型</th><th>UID</th><th>IP</th><th>详情</th><th>时间</th></tr></thead>
        <tbody>
          <tr v-for="(e, i) in shownEvents" :key="i">
            <td><span class="tag red">{{ kindName(e.kind) }}</span></td>
            <td class="mono">{{ e.uid || '-' }}</td>
            <td class="mono">{{ e.ip }}</td>
            <td class="muted">{{ e.detail }}</td>
            <td class="muted">{{ fmtTime(e.time) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="card">
      <div class="between"><b>后台操作日志</b>
        <input v-model="adminQ" class="input" placeholder="管理员ID" style="width:100px" />
        <input v-model="actionQ" class="input" placeholder="按动作筛选" style="width:120px" />
        <select v-model.number="logHours" class="input" style="width:112px">
          <option :value="0">全部时间</option><option :value="1">近 1 小时</option>
          <option :value="24">近 24 小时</option><option :value="168">近 7 天</option>
        </select>
        <button class="btn btn-sm" @click="exportLogs">导出</button>
      </div>
      <table class="table" style="margin-top:10px">
        <thead><tr><th>管理员</th><th>动作</th><th>对象</th><th>明细</th><th>IP</th><th>时间</th></tr></thead>
        <tbody>
          <tr v-for="(l, i) in shownLogs" :key="i">
            <td class="mono">{{ l.adminId }}</td>
            <td>{{ l.action }}</td>
            <td class="mono">{{ l.target }}</td>
            <td class="muted">{{ l.detail }}</td>
            <td class="mono muted">{{ l.ip }}</td>
            <td class="muted">{{ fmtTime(l.time) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, fmtTime } from '@/lib/api'

const reports = ref<any[]>([])
const pendingCount = computed(() => reports.value.filter(r => !r.status).length)
const rpStatus = ref(0)
const events = ref<any[]>([])
const hours = ref(0)
const kind = ref('')
const actionQ = ref('')
const adminQ = ref('')
const logHours = ref(0)
const logs = ref<any[]>([])
const shownEvents = computed(() => {
  let rows = events.value
  if (kind.value) rows = rows.filter((e: any) => e.kind === kind.value)
  if (hours.value) {
    const cut = Date.now() - hours.value * 3600 * 1000
    rows = rows.filter((e: any) => (e.time || 0) >= cut)
  }
  return rows
})
const shownLogs = computed(() => {
  let rows = logs.value
  const k = actionQ.value.trim()
  const a = adminQ.value.trim()
  if (k) rows = rows.filter((l: any) => String(l.action || '').includes(k))
  if (a) rows = rows.filter((l: any) => String(l.adminId || '').includes(a))
  if (logHours.value) {
    const cut = Date.now() - logHours.value * 3600 * 1000
    rows = rows.filter((l: any) => (l.time || 0) >= cut)
  }
  return rows
})
const reasonName = (k: string) => ({ abuse: '辱骂骚扰', idle: '消极游戏', cheat: '疑似作弊', other: '其他' } as any)[k] || k
const reasonTag = (k: string) => (k === 'cheat' ? 'red' : k === 'abuse' ? 'warn' : 'gold')
async function loadReports() {
  try { reports.value = await api.get<any[]>(`reports?status=${rpStatus.value}`) } catch { reports.value = [] }
}
async function resolve(r: any) {
  try {
    await api.post('report/resolve', { id: r.id, note: '后台标记已处理' })
    await loadReports()
  } catch { /* 失败保持原状，下次刷新可见 */ }
}
const kindName = (k: string) => ({ rate_limit: '限流', bad_token: '令牌异常', sql_inject: '注入探测', ws_abuse: '连接滥用', admin_try: '后台爆破', xss: 'XSS', brute_force: '撞库锁定', login_lock: '登录锁定', player_report: '玩家举报' } as any)[k] || k

function exportLogs() {
  const rows = [['adminId', 'action', 'target', 'detail', 'ip', 'time']]
  for (const l of shownLogs.value) rows.push([l.adminId, l.action, l.target, l.detail, l.ip, new Date(l.time).toISOString()])
  const csv = rows.map(r => r.map(v => String(v ?? '').replace(/,/g, ' ')).join(',')).join('\n')
  const blob = new Blob([new Uint8Array([0xEF, 0xBB, 0xBF]), csv], { type: 'text/csv;charset=utf-8' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `admin-logs-${Date.now()}.csv`
  a.click()
  URL.revokeObjectURL(a.href)
}

async function load() {
  loadReports()
  try { events.value = await api.get<any[]>('security') } catch {}
  try { logs.value = await api.get<any[]>('logs') } catch {}
}
onMounted(load)
</script>
