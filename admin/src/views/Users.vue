<template>
  <div class="users-page">
    <div class="toolbar card">
      <div class="row">
        <input v-model="q" class="input" placeholder="搜索 UID / 昵称" style="width:210px" />
        <select v-model="kind" class="input" style="width:120px">
          <option value="all">全部用户</option><option value="human">仅真人</option>
          <option value="robot">仅机器人</option><option value="muted">禁言中</option>
        </select>
        <button class="btn btn-sm" @click="load">查询</button>
      </div>
      <div class="row">
        <span class="sum-chip">共 <b>{{ filtered.length }}</b> 个用户</span>
        <span class="sum-chip">真人 <b>{{ humanCount }}</b></span>
        <span class="sum-chip">机器人 <b>{{ robotCount }}</b></span>
        <span class="sum-chip gold">总余额 <b class="mono">{{ money(totalBalance) }}</b></span>
      </div>
    </div>

    <div class="card list-card">
      <table class="table">
        <thead><tr><th>用户</th><th>UID</th><th>VIP</th><th>余额</th><th>房卡</th><th>类型</th><th>状态</th><th>注册时间</th><th style="width:120px">操作</th></tr></thead>
        <tbody>
          <tr v-for="u in filtered" :key="u.uid">
            <td>
              <div class="user-cell">
                <span class="u-avatar">{{ (u.nick || 'U').slice(0, 1) }}</span>
                <b class="u-nick">{{ u.nick }}</b>
              </div>
            </td>
            <td class="mono">{{ u.uid }}</td>
            <td><span class="tag gold">VIP {{ u.vip }}</span></td>
            <td class="mono">{{ money(u.balance) }}</td>
            <td class="mono">{{ u.cards }}</td>
            <td><span class="tag" :class="u.robot ? 'gray' : ''">{{ u.robot ? '机器人' : '真人' }}</span></td>
            <td>
              <span class="state-pill" :class="u.status === 1 ? 'ok' : (u.status === 2 ? 'warn' : 'bad')">
                {{ u.status === 1 ? '正常' : u.status === 2 ? '禁言' : '封禁' }}
              </span>
            </td>
            <td class="muted">{{ fmtTime(u.createdAt) }}</td>
            <td>
              <div class="row" style="gap:6px">
                <button class="btn btn-sm" @click="openDetail(u)">详情</button>
                <div class="menu-wrap">
                  <button class="btn btn-sm" @click.stop="menuFor = menuFor === u.uid ? 0 : u.uid">⋯</button>
                  <div v-if="menuFor === u.uid" class="mini-menu" @click.stop>
                    <button @click="adjust(u, 'balance'); menuFor = 0">改余额</button>
                    <button @click="adjust(u, 'cards'); menuFor = 0">改房卡</button>
                    <button @click="mute(u); menuFor = 0">{{ u.status === 2 ? '解除禁言' : '禁言' }}</button>
                    <button @click="kick(u); menuFor = 0">强制下线</button>
                    <button @click="showLogs(u); menuFor = 0">操作时间轴</button>
                    <button class="danger" @click="toggle(u); menuFor = 0">{{ u.status === 1 ? '封禁账号' : '解除封禁' }}</button>
                  </div>
                </div>
              </div>
            </td>
          </tr>
          <tr v-if="!filtered.length"><td colspan="9" class="muted" style="text-align:center;padding:20px">没有匹配的用户</td></tr>
        </tbody>
      </table>
    </div>
  </div>

  <!-- 用户详情抽屉 -->
  <div v-if="detail" class="drawer-mask" @click.self="detail = null">
    <aside class="drawer">
      <header class="drawer-head">
        <div class="col">
          <b>用户详情</b>
          <span class="muted tiny">UID {{ detail.uid }} · 注册于 {{ fmtTime(detail.createdAt) }}</span>
        </div>
        <button class="btn btn-sm" @click="detail = null">关闭</button>
      </header>
      <div class="drawer-body">
        <div class="profile-card">
          <span class="u-avatar lg">{{ (detail.nick || 'U').slice(0, 1) }}</span>
          <div class="col" style="flex:1">
            <b class="dv-nick">{{ detail.nick }}</b>
            <span class="muted tiny">VIP {{ detail.vip }} · {{ detail.status === 1 ? '正常' : detail.status === 2 ? '禁言中' : '已封禁' }}</span>
          </div>
          <div class="col" style="align-items:flex-end">
            <b class="mono gold">{{ money(detail.balance) }}</b>
            <span class="muted tiny">房卡 {{ detail.cards }} 张</span>
          </div>
        </div>

        <div class="metric-grid">
          <div class="metric"><span class="muted tiny">总手数</span><b class="mono">{{ detail.hands }}</b></div>
          <div class="metric"><span class="muted tiny">胜率</span><b class="mono">{{ detail.winRate }}%</b></div>
          <div class="metric"><span class="muted tiny">累计盈亏</span><b class="mono" :class="detail.net >= 0 ? 'green' : 'red'">{{ money(detail.net) }}</b></div>
          <div class="metric"><span class="muted tiny">冻结</span><b class="mono">{{ money(detail.frozen) }}</b></div>
        </div>

        <div class="drawer-actions">
          <button class="btn btn-sm" @click="adjust(detail, 'balance')">改余额</button>
          <button class="btn btn-sm" @click="adjust(detail, 'cards')">改房卡</button>
          <button class="btn btn-sm" @click="mute(detail)">{{ detail.status === 2 ? '解除禁言' : '禁言' }}</button>
          <button class="btn btn-sm" @click="kick(detail)">强制下线</button>
          <button class="btn btn-sm btn-danger" @click="toggle(detail)">{{ detail.status === 1 ? '封禁' : '解封' }}</button>
        </div>

        <div class="section-title">最近资金流水</div>
        <div class="mini-table">
          <div v-for="b in detail.bills || []" :key="b.id" class="mini-row">
            <span class="tag gray">{{ b.type }}</span>
            <span class="mono" :class="b.amount >= 0 ? 'green' : 'red'">{{ b.amount >= 0 ? '+' : '' }}{{ (b.amount / 100).toFixed(2) }}</span>
            <span class="muted tiny">{{ b.remark }}</span>
            <span class="muted tiny" style="margin-left:auto">{{ fmtTime(b.time) }}</span>
          </div>
          <div v-if="!(detail.bills || []).length" class="muted tiny" style="padding:8px 0">暂无流水</div>
        </div>

        <div class="section-title">后台操作记录</div>
        <div class="mini-table">
          <div v-for="(l, i) in detail.logs || []" :key="'l' + i" class="mini-row">
            <span class="tag">{{ l.action }}</span>
            <span class="muted tiny" style="flex:1">{{ l.detail }}</span>
            <span class="muted tiny">{{ fmtTime(l.createdAt) }}</span>
          </div>
          <div v-if="!(detail.logs || []).length" class="muted tiny" style="padding:8px 0">暂无操作记录</div>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, money, fmtTime } from '@/lib/api'

const list = ref<any[]>([])
const q = ref('')
const kind = ref('all')
const editing = ref<any>(null)
const banning = ref<any>(null)
const logUser = ref<any>(null)
const userLogs = ref<any[]>([])
const menuFor = ref(0)
const detail = ref<any>(null)

const humanCount = computed(() => filtered.value.filter(u => !u.robot).length)
const robotCount = computed(() => filtered.value.filter(u => u.robot).length)
const totalBalance = computed(() => filtered.value.reduce((a, b) => a + Number(b.balance || 0), 0))

async function openDetail(u: any) {
  menuFor.value = 0
  try {
    detail.value = await api.get<any>(`user/detail?uid=${u.uid}`)
  } catch {
    detail.value = null
    toast('用户详情读取失败', 'err')
  }
}

const filtered = computed(() => {
  const k = q.value.trim()
  let rows = list.value
  if (kind.value === 'human') rows = rows.filter(u => !u.robot)
  else if (kind.value === 'robot') rows = rows.filter(u => u.robot)
  else if (kind.value === 'muted') rows = rows.filter(u => u.status === 2)
  if (!k) return rows
  return rows.filter(u => String(u.uid).includes(k) || (u.nick || '').includes(k))
})

async function load() {
  try { list.value = await api.get<any[]>('users?limit=200') } catch {}
}

function adjust(u: any, kind: 'balance' | 'cards') {
  editing.value = { u, kind, delta: kind === 'balance' ? 10000 : 10, remark: '' }
}

async function submitAdjust() {
  const e = editing.value
  if (!e || !e.delta) return
  try {
    if (e.kind === 'balance') {
      await api.post('user/balance', { uid: e.u.uid, delta: e.delta, remark: e.remark })
    } else {
      await api.post('user/cards', { uid: e.u.uid, delta: e.delta })
    }
    toast('调整成功')
    editing.value = null
    load()
  } catch (err: any) { toast(err.message, 'err') }
}

function mute(u: any) {
  if (u.status === 2) return setStatus(u, 1, '已解禁')
  banning.value = { u, status: 2, reason: '' }
}

function toggle(u: any) {
  if (u.status !== 1) return setStatus(u, 1, '已解封')
  banning.value = { u, status: 0, reason: '' }
}

async function confirmBan() {
  const b = banning.value
  if (!b) return
  await setStatus(b.u, b.status, b.status === 0 ? '已封禁并踢下线' : '已禁言', b.reason)
  banning.value = null
}

async function setStatus(u: any, status: number, okText: string, reason = '') {
  try {
    await api.post('user/status', { uid: u.uid, status, reason })
    toast(okText)
    load()
  } catch (err: any) { toast(err.message, 'err') }
}

async function kick(u: any) {
  try {
    await api.post('user/kick', { uid: u.uid })
    toast('已强制下线')
  } catch (err: any) { toast(err.message, 'err') }
}

async function showLogs(u: any) {
  logUser.value = u
  userLogs.value = []
  try { userLogs.value = await api.get<any[]>(`user/logs?uid=${u.uid}`) } catch { userLogs.value = [] }
}

function toast(text: string, kind = 'ok') {
  window.dispatchEvent(new CustomEvent('admin-toast', { detail: { text, kind } }))
}

onMounted(load)
</script>

<style scoped>
.timeline-row { display: flex; align-items: center; gap: 8px; padding: 7px 0; border-bottom: 1px solid #152a1e; }
.timeline-row:last-child { border-bottom: 0; }
.tl-dot { width: 7px; height: 7px; border-radius: 50%; background: #7ff0b0; box-shadow: 0 0 8px rgba(127, 240, 176, .6); flex: none; }
</style>

<style scoped>
.users-page { display: flex; flex-direction: column; gap: 12px; }
.toolbar { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 12px 14px; flex-wrap: wrap; }
.sum-chip { font-size: 11.5px; color: #9fb8ab; padding: 3px 9px; border-radius: 8px; background: rgba(255, 255, 255, .04); border: 1px solid rgba(60, 120, 88, .3); }
.sum-chip b { color: #d7ffe9; }
.sum-chip.gold b { color: var(--gold); }
.list-card { padding: 6px 0 0; }
.user-cell { display: flex; align-items: center; gap: 8px; }
.u-avatar { width: 26px; height: 26px; border-radius: 8px; display: flex; align-items: center; justify-content: center;
  background: radial-gradient(circle at 34% 30%, rgba(46, 224, 122, .3), rgba(9, 30, 20, .95)); border: 1px solid rgba(46, 224, 122, .35);
  color: #d7ffe9; font-size: 12px; font-weight: 800; flex: none; }
.u-avatar.lg { width: 44px; height: 44px; border-radius: 12px; font-size: 18px; }
.u-nick { font-size: 12.5px; }
.state-pill { display: inline-flex; align-items: center; height: 21px; padding: 0 8px; border-radius: 7px; font-size: 10.5px; font-weight: 700; }
.state-pill.ok { background: rgba(46, 224, 122, .16); color: #7dffb0; border: 1px solid rgba(46, 224, 122, .35); }
.state-pill.warn { background: rgba(255, 209, 102, .16); color: #ffd166; border: 1px solid rgba(255, 209, 102, .35); }
.state-pill.bad { background: rgba(255, 122, 122, .16); color: #ffa8a8; border: 1px solid rgba(255, 122, 122, .35); }
.menu-wrap { position: relative; }
.mini-menu { position: absolute; right: 0; top: 30px; z-index: 30; min-width: 132px; padding: 6px; border-radius: 10px;
  background: #0c1b13; border: 1px solid var(--line); box-shadow: 0 12px 28px rgba(0, 0, 0, .5); display: flex; flex-direction: column; }
.mini-menu button { text-align: left; padding: 7px 9px; border-radius: 7px; background: none; border: none; color: var(--text); font-size: 12px; cursor: pointer; }
.mini-menu button:hover { background: rgba(46, 224, 122, .12); }
.mini-menu button.danger { color: #ffa8a8; }
.col { display: flex; flex-direction: column; gap: 2px; }
.drawer-mask { position: fixed; inset: 0; background: rgba(2, 8, 5, .6); z-index: 300; display: flex; justify-content: flex-end; }
.drawer { width: min(480px, 96vw); height: 100%; background: linear-gradient(180deg, #0b1c14, #06110c); border-left: 1px solid var(--line);
  display: flex; flex-direction: column; animation: drawerIn .22s ease; }
@keyframes drawerIn { from { transform: translateX(24px); opacity: .4 } to { transform: none; opacity: 1 } }
.drawer-head { display: flex; align-items: center; justify-content: space-between; padding: 14px 16px; border-bottom: 1px solid #14261b; }
.drawer-body { flex: 1; overflow-y: auto; padding: 14px 16px 26px; display: flex; flex-direction: column; gap: 12px; }
.profile-card { display: flex; align-items: center; gap: 10px; padding: 12px; border-radius: 12px;
  background: rgba(18, 42, 30, .8); border: 1px solid rgba(60, 120, 88, .3); }
.dv-nick { font-size: 15px; }
.metric-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; }
.metric { display: flex; flex-direction: column; gap: 2px; padding: 8px 9px; border-radius: 9px; background: #0a1710; border: 1px solid #152a1e; }
.drawer-actions { display: flex; flex-wrap: wrap; gap: 8px; }
.section-title { font-size: 12px; font-weight: 700; color: var(--gold); margin-top: 2px; }
.mini-table { display: flex; flex-direction: column; gap: 5px; }
.mini-row { display: flex; align-items: center; gap: 8px; font-size: 12px; padding: 6px 8px; border-radius: 8px; background: rgba(255, 255, 255, .03); }
.green { color: #7dffb0; }
.red { color: #ffa8a8; }
</style>
