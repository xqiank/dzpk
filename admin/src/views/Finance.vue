<template>
  <div class="page finance-page">
    <div class="toolbar">
      <h2 class="page-title">充提审核</h2>
      <span class="sum-chip">待核实充币 <b>{{ stats.pendingDeposits ?? 0 }}</b> 笔 · ₮ {{ Number(stats.pendingDepositAmount || 0).toFixed(2) }}</span>
      <span class="sum-chip" :class="{ warn: Number(stats.pendingWithdraws || 0) > 0 }">待打款提币 <b>{{ stats.pendingWithdraws ?? 0 }}</b> 笔</span>
      <span class="sum-chip">累计已打款 ₮ {{ Number(stats.paidWithdrawAmount || 0).toFixed(2) }}</span>
      <div class="toolbar-right">
        <select v-model="status" class="input" style="width: 130px" @change="load">
          <option value="all">全部状态</option>
          <option value="pending">待处理</option>
          <option value="completed">已完成</option>
          <option value="rejected">已驳回</option>
          <option value="cancelled">已取消</option>
        </select>
        <select v-model="type" class="input" style="width: 120px" @change="load">
          <option value="all">全部类型</option>
          <option value="deposit">充币</option>
          <option value="withdraw">提币</option>
        </select>
        <button class="btn btn-sm" @click="load">刷新</button>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <b>USDT 订单</b>
        <span class="muted">共 {{ list.length }} 笔 · 充币核实入账 / 提币审核打款 / 驳回自动退冻结</span>
      </div>
      <div class="panel-body">
        <div v-if="!list.length" class="empty">暂无符合条件的订单</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th style="width: 70px">订单号</th>
              <th style="width: 140px">用户</th>
              <th style="width: 76px">类型</th>
              <th style="width: 90px">金额</th>
              <th style="width: 70px">网络</th>
              <th>地址 / TxID</th>
              <th style="width: 90px">状态</th>
              <th style="width: 150px">时间</th>
              <th style="width: 210px">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="o in list" :key="o.id">
              <td class="mono">#{{ o.id }}</td>
              <td>
                <div class="user-cell">
                  <b>{{ o.nick || '用户' }}</b>
                  <small class="mono">UID {{ o.uid }}</small>
                </div>
              </td>
              <td><span class="tag" :class="o.type === 'deposit' ? 'blue' : 'orange'">{{ o.type === 'deposit' ? '充币' : '提币' }}</span></td>
              <td class="mono">
                <b>{{ o.type === 'deposit' ? '+' : '-' }}{{ Number(o.amount).toFixed(2) }}</b>
                <small v-if="Number(o.fee) > 0" class="muted"> 费 {{ Number(o.fee).toFixed(2) }}</small>
              </td>
              <td class="mono">{{ o.network }}</td>
              <td class="mono ellipsis" :title="o.address">{{ o.address || '—' }}</td>
              <td><span class="state-pill" :class="stateClass(o.status)">{{ stateLabel(o.status) }}</span></td>
              <td class="mono muted">{{ fmt(o.createdAt) }}</td>
              <td>
                <div v-if="o.status === 'pending'" class="row-actions">
                  <button v-if="o.type === 'deposit'" class="btn btn-xs primary" :disabled="busy === o.id" @click="approve(o)">核实入账</button>
                  <button v-else class="btn btn-xs primary" :disabled="busy === o.id" @click="pay(o)">审核打款</button>
                  <button class="btn btn-xs danger" :disabled="busy === o.id" @click="reject(o)">驳回</button>
                </div>
                <span v-else class="muted">{{ o.remark || '—' }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { api } from '@/lib/api'

const list = ref<any[]>([])
const stats = reactive<any>({})
const status = ref('all')
const type = ref('all')
const busy = ref(0)
const errMsg = ref('')

function stateLabel(s: string) {
  return s === 'pending' ? '待处理' : s === 'completed' ? '已完成' : s === 'rejected' ? '已驳回' : s === 'cancelled' ? '已取消' : s
}
function stateClass(s: string) {
  return s === 'pending' ? 'gold' : s === 'completed' ? 'green' : 'red'
}
function fmt(sec: number) {
  if (!sec) return '—'
  const d = new Date(sec * 1000)
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

async function load() {
  errMsg.value = ''
  try {
    const q = `?status=${status.value}&type=${type.value}`
    const res = await api.get<any>('usdt/orders' + q)
    list.value = res.list || []
    Object.assign(stats, res.stats || {})
  } catch (e: any) {
    errMsg.value = e?.message || '加载失败'
  }
}

async function approve(o: any) {
  busy.value = o.id
  try {
    const res = await api.post<any>('usdt/approve', { id: o.id })
    Object.assign(stats, res.stats || {})
    await load()
  } catch (e: any) {
    errMsg.value = e?.message || '入账失败'
  } finally {
    busy.value = 0
  }
}

async function pay(o: any) {
  busy.value = o.id
  try {
    const res = await api.post<any>('usdt/pay', { id: o.id, txid: `manual-${Date.now()}` })
    Object.assign(stats, res.stats || {})
    await load()
  } catch (e: any) {
    errMsg.value = e?.message || '打款失败'
  } finally {
    busy.value = 0
  }
}

async function reject(o: any) {
  busy.value = o.id
  try {
    const res = await api.post<any>('usdt/reject', { id: o.id, remark: '运营审核驳回' })
    Object.assign(stats, res.stats || {})
    await load()
  } catch (e: any) {
    errMsg.value = e?.message || '驳回失败'
  } finally {
    busy.value = 0
  }
}

onMounted(load)
</script>

<style scoped>
.finance-page { display: flex; flex-direction: column; gap: 14px; }
.user-cell { display: flex; flex-direction: column; gap: 2px; }
.user-cell small { font-size: 11px; color: var(--muted); }
.ellipsis { max-width: 240px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.row-actions { display: flex; gap: 6px; }
.empty { padding: 36px; text-align: center; color: var(--muted); }
</style>
