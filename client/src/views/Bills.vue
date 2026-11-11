<template>
  <div class="page bills">
    <header class="b-head">
      <button class="head-btn" @click="router.back()"><Icon name="back" :size="16" /> 返回</button>
      <b class="h3">资产明细</b>
      <span style="width:60px"></span>
    </header>

    <div class="card card-pad balance">
      <div class="row-between">
        <div class="col">
          <span class="tiny muted">账户总资产 (USDT)</span>
          <b class="big">{{ money(auth.profile?.balance) }}</b>
        </div>
        <button class="btn btn-primary btn-sm" @click="router.push('/profile')">充值</button>
      </div>
    </div>

    <div class="row gap-8" style="margin:12px 0">
      <button v-for="t in types" :key="t.key" class="chip" :class="{ on: filter === t.key }" @click="filter = t.key">{{ t.label }}</button>
    </div>

    <div class="card list">
      <div v-for="(b, i) in visible" :key="i" class="bill-item">
        <div class="ico" :class="b.amount >= 0 ? 'in' : 'out'">
          <Icon :name="b.amount >= 0 ? 'arrowDown' : 'arrowUp'" :size="15" />
        </div>
        <div class="col" style="flex:1">
          <b class="sm">{{ labelOf(b.type) }}</b>
          <span class="tiny muted">{{ fmtTime(b.time, true) }}{{ b.remark ? ' · ' + b.remark : '' }}</span>
        </div>
        <div class="col" style="align-items:flex-end">
          <b :class="b.amount >= 0 ? 'green' : 'red'">{{ b.amount >= 0 ? '+' : '' }}{{ money(b.amount) }}</b>
          <span class="tiny muted">余额 {{ money(b.after) }}</span>
        </div>
      </div>
      <div v-if="!visible.length" class="empty tiny muted">暂无记录</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Icon from '@/components/Icon.vue'
import { socket } from '@/lib/ws'
import { money, fmtTime } from '@/lib/format'
import { useAuth } from '@/stores/auth'

const router = useRouter()
const auth = useAuth()
const bills = ref<any[]>([])
const filter = ref('all')
const types = [
  { key: 'all', label: '全部' },
  { key: 'recharge', label: '充值' },
  { key: 'withdraw', label: '提现' },
  { key: 'buyin', label: '带入' },
  { key: 'cashout', label: '带出' },
  { key: 'card', label: '房卡' }
]
const labels: Record<string, string> = {
  recharge: '充值', withdraw: '提现', buyin: '带入牌桌', cashout: '离开牌桌',
  card: '房卡购买/调整', rake: '服务费', reward: '奖励', adjust: '后台调整', insurance: '保险结算'
}
function labelOf(t: string) { return labels[t] || t }
const visible = computed(() => (filter.value === 'all' ? bills.value : bills.value.filter(b => b.type === filter.value)))

onMounted(async () => {
  try { bills.value = await socket.call('user.bills', { limit: 100 }) } catch {}
})
</script>

<style scoped>
.bills { padding-top: 14px; }
.b-head { display: flex; align-items: center; gap: 10px; }
.head-btn { height: 34px; padding: 0 12px; border-radius: 10px; background: rgba(255,255,255,.06); border: 1px solid var(--line); font-size: 13px; display: flex; align-items: center; gap: 4px; }
.balance { margin-top: 14px; }
.big { font-size: 28px; font-weight: 900; }
.chip { height: 30px; padding: 0 12px; border-radius: 15px; font-size: 12px; color: var(--muted); background: rgba(10,30,20,.7); border: 1px solid var(--line); }
.chip.on { color: #241a06; background: linear-gradient(180deg, #f2d493, #cfa151); border-color: transparent; font-weight: 700; }
.list { overflow: hidden; }
.bill-item { display: flex; align-items: center; gap: 10px; padding: 12px 14px; border-bottom: 1px solid rgba(122,196,156,.08); }
.ico { width: 32px; height: 32px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex: none; }
.ico.in { background: rgba(46,224,122,.14); color: #7ff0b0; }
.ico.out { background: rgba(255,120,120,.12); color: #ff9f9f; }
.empty { padding: 24px; text-align: center; }
</style>
