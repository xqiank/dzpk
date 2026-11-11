<template>
  <div class="page rooms-page">
    <div class="toolbar">
      <h2 class="page-title">房间管理</h2>
      <span class="sum-chip">官方房 <b>{{ officialCount }}</b></span>
      <span class="sum-chip">私人房 <b>{{ privateCount }}</b></span>
      <span class="sum-chip">在线牌桌 <b>{{ liveCount }}</b></span>
      <span class="sum-chip" :class="{ warn: officialDown > 0 }">已下架官方房 <b>{{ officialDown }}</b></span>
      <div class="toolbar-right">
        <select v-model="gameFilter" class="input" style="width: 118px">
          <option value="">全部玩法</option><option value="nlh">长牌</option>
          <option value="short6">短牌</option><option value="omaha">奥马哈</option><option value="niuniu">牛牛</option>
        </select>
        <button class="btn btn-sm" @click="toggleClubs">{{ showClubs ? '隐藏俱乐部' : '俱乐部' }}</button>
        <button class="btn btn-sm" @click="clearPrivate">清理私人房</button>
        <button class="btn btn-sm btn-danger" :disabled="!selected.length" @click="batchOff">批量下架 ({{ selected.length }})</button>
        <button class="btn btn-sm" @click="load">刷新</button>
        <button class="btn btn-primary btn-sm" @click="edit(null)">+ 新建房间</button>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <b>房间列表</b>
        <span class="muted">共 {{ shownRooms.length }} 个房间 · 上架 / 下架即时生效</span>
      </div>
      <div class="panel-body">
    <table v-if="showClubs" class="data-table" style="margin-bottom:14px">
      <thead><tr><th>ID</th><th>编号</th><th>名称</th><th>房主</th><th>卡池</th><th>成员</th><th>状态</th><th>操作</th></tr></thead>
      <tbody>
        <tr v-for="c in clubs" :key="c.id" @click="loadMembers(c.id)" style="cursor:pointer">
          <td class="mono">{{ c.id }}</td><td class="mono">{{ c.clubNo }}</td><td>{{ c.name }}</td>
          <td class="mono">{{ c.owner }}</td><td class="mono">{{ c.cardPool }}</td><td class="mono">{{ c.members }}</td>
          <td>
            <span class="state-pill" :class="c.status ? 'green' : 'red'">{{ c.status ? '正常' : '停用' }}</span>
            <button class="btn btn-sm" style="margin-left:6px" @click.stop="toggleClub(c)">{{ c.status ? '停用' : '启用' }}</button>
            <button class="btn btn-sm" :class="c.certified ? 'btn-primary' : ''" style="margin-left:6px" @click.stop="toggleCert(c)">
              {{ c.certified ? '已认证' : (c.certifyRequest ? '待审·通过' : '认证') }}
            </button>
            <span v-if="c.certifyRequest && !c.certified" class="tag warn" style="margin-left:6px">申请中</span>
          </td>
          <td>
            <button class="btn btn-sm" @click.stop="addCards(c)">补卡</button>
            <button class="btn btn-sm" style="margin-left:6px" @click.stop="editClub(c)">改名</button>
            <button class="btn btn-sm btn-danger" style="margin-left:6px" @click.stop="kick(c)">踢人</button>
          </td>
        </tr>
        <tr v-if="membersOf">
          <td colspan="8" class="muted tiny">成员：{{ membersOf }}</td>
        </tr>
      </tbody>
    </table>
    <table class="data-table">
      <thead><tr><th><input type="checkbox" :checked="allChecked" @change="toggleAll" /></th><th>房间号</th><th>名称</th><th>玩法</th><th>人数</th><th>盲注</th><th>买入</th><th>房卡</th><th>保险</th><th>类型</th><th>在线</th><th>归属</th><th>操作</th></tr></thead>
      <tbody>
        <tr v-for="r in shownRooms" :key="r.RoomNo">
          <td><input type="checkbox" :value="r.RoomNo" v-model="selected" /></td>
          <td class="mono">{{ r.RoomNo }}</td>
          <td>{{ r.Name }}</td>
          <td><span class="tag">{{ gameName(r.GameType) }}</span></td>
          <td>{{ r.SeatCount }}</td>
          <td class="mono">{{ money(r.SB) }}/{{ money(r.BB) }}</td>
          <td class="mono">{{ money(r.MinBuyin) }}~{{ money(r.MaxBuyin) }}</td>
          <td>{{ r.CardCost }}</td>
          <td>{{ r.Insurance ? '开' : '关' }}</td>
          <td><span class="state-pill" :class="r.IsOfficial ? 'gold' : ''">{{ r.IsOfficial ? '官方' : '私人' }}</span></td>
          <td class="mono"><router-link :to="'/tables'" class="tag gray" style="text-decoration:none">{{ live[r.RoomNo] || '0' }}</router-link></td>
          <td class="mono">{{ r.ClubID ? '俱乐部 ' + r.ClubID : (r.HostUID ? '玩家 ' + r.HostUID : '系统') }}</td>
          <td>
            <button class="btn btn-sm" @click="edit(r)">编辑</button>
            <button class="btn btn-sm" style="margin-left:6px" @click="clone(r)">复制</button>
            <button class="btn btn-sm" style="margin-left:6px" @click="relist(r)">上架</button>
            <button class="btn btn-sm btn-danger" style="margin-left:6px" @click="del(r)">下架</button>
          </td>
        </tr>
      </tbody>
    </table>
      </div>
    </div>

  <div v-if="form" class="modal-mask" @click.self="form = null">
    <div class="modal">
      <b>{{ form.RoomNo ? '编辑房间 ' + form.RoomNo : '新建房间' }}</b>
      <div class="form-grid">
        <label><span>房间号</span><input v-model.number="form.RoomNo" class="input" type="number" /></label>
        <label><span>名称</span><input v-model="form.Name" class="input" /></label>
        <label><span>玩法</span>
          <select v-model="form.GameType" class="input">
            <option value="nlh">长牌</option><option value="short6">短牌</option>
            <option value="omaha">奥马哈</option><option value="niuniu">牛牛</option>
          </select>
        </label>
        <label><span>人数</span><input v-model.number="form.SeatCount" class="input" type="number" /></label>
        <label><span>小盲(分)</span><input v-model.number="form.SB" class="input" type="number" /></label>
        <label><span>大盲(分)</span><input v-model.number="form.BB" class="input" type="number" /></label>
        <label><span>最小买入(分)</span><input v-model.number="form.MinBuyin" class="input" type="number" /></label>
        <label><span>最大买入(分)</span><input v-model.number="form.MaxBuyin" class="input" type="number" /></label>
        <label><span>房卡消耗</span><input v-model.number="form.CardCost" class="input" type="number" /></label>
        <label><span>保险</span><select v-model.number="form.Insurance" class="input"><option :value="1">开启</option><option :value="0">关闭</option></select></label>
        <label><span>官方房</span><select v-model.number="form.IsOfficial" class="input"><option :value="1">是</option><option :value="0">否</option></select></label>
        <label><span>排序</span><input v-model.number="form.Sort" class="input" type="number" /></label>
      </div>
      <p class="tiny" :style="{ color: blindHint.ok ? '#7ff0b0' : '#ff9f9f', marginTop: '10px' }">{{ blindHint.text }}</p>
      <div class="row" style="justify-content:flex-end;margin-top:14px">
        <button class="btn" @click="form = null">取消</button>
        <button class="btn btn-primary" @click="save">保存</button>
      </div>
    </div>
  </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, money } from '@/lib/api'

const rooms = ref<any[]>([])
const form = ref<any>(null)
const showClubs = ref(false)
const clubs = ref<any[]>([])
const membersOf = ref('')
const officialCount = computed(() => rooms.value.filter(r => r.IsOfficial).length)
const privateCount = computed(() => rooms.value.filter(r => !r.IsOfficial).length)
const liveCount = computed(() => Object.values(live.value || {}).filter(v => Number(String(v).split('/')[0]) > 0).length)
const officialDown = computed(() => rooms.value.filter(r => r.IsOfficial && Number(r.Sort || 0) < 0).length)
async function toggleCert(c: any) {
  try {
    const r = await api.post<{ certified: boolean }>('club/certify', { id: c.id, on: !c.certified })
    c.certified = r.certified
    toast(r.certified ? '已设为官方认证' : '已取消认证')
  } catch (e: any) {
    toast(e.message)
  }
}

async function toggleClub(c: any) {
  try {
    await api.post('club/status', { id: c.id, status: c.status ? 0 : 1 })
    toast('已更新')
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

async function editClub(c: any) {
  const name = prompt('俱乐部名称', c.name)
  if (name === null) return
  const notice = prompt('俱乐部公告', '') || ''
  try {
    await api.post('club/update', { id: c.id, name, notice })
    toast('已保存')
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

async function kick(c: any) {
  const raw = prompt('要移出哪个 UID？')
  if (!raw) return
  try {
    await api.post('club/kick', { id: c.id, uid: Number(raw) })
    toast('已移出')
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

async function addCards(c: any) {
  const n = Number(prompt('补多少张房卡给「' + c.name + '」？', '10'))
  if (!n) return
  try {
    const d = await api.post<{ cardPool: number }>('club/cards', { id: c.id, delta: n })
    c.cardPool = d.cardPool
    toast('卡池已更新为 ' + d.cardPool)
  } catch (e: any) { toast(e.message, 'err') }
}

async function loadMembers(id: number) {
  try {
    const list = await api.get<any[]>('club/members?id=' + id)
    membersOf.value = list.map((m: any) => `${m.nick}(${m.uid})`).join('、') || '暂无成员'
  } catch { membersOf.value = '' }
}
async function toggleClubs() {
  showClubs.value = !showClubs.value
  if (showClubs.value && !clubs.value.length) { try { clubs.value = await api.get('clubs') } catch {} }
}
const selected = ref<number[]>([])
const gameFilter = ref('')
const blindHint = computed(() => {
  const f = form.value
  if (!f) return { ok: true, text: '' }
  if (!(f.BB > 0) || !(f.SB > 0) || f.SB > f.BB) return { ok: false, text: '盲注需满足 0 < 小盲 ≤ 大盲' }
  if (!(f.MinBuyin > 0) || f.MaxBuyin < f.MinBuyin) return { ok: false, text: '买入需满足 0 < 最小 ≤ 最大' }
  return { ok: true, text: '配置合法：' + (f.SB / 100).toFixed(2) + '/' + (f.BB / 100).toFixed(2) + '，买入 ' + (f.MinBuyin / 100).toFixed(2) + '~' + (f.MaxBuyin / 100).toFixed(2) }
})
const shownRooms = computed(() => (gameFilter.value ? rooms.value.filter(r => r.GameType === gameFilter.value) : rooms.value))
const allChecked = computed(() => shownRooms.value.length > 0 && selected.value.length === shownRooms.value.length)
function toggleAll(e: any) {
  selected.value = e.target.checked ? shownRooms.value.map(r => r.RoomNo) : []
}
async function clearPrivate() {
  const targets = rooms.value.filter(r => !r.IsOfficial)
  if (!targets.length) { toast('暂无私人房'); return }
  for (const r of targets) { try { await api.post('room/delete', { roomNo: r.RoomNo }) } catch {} }
  toast('已清理 ' + targets.length + ' 间私人房')
  load()
}

async function batchOff() {
  if (!selected.value.length) return
  for (const no of selected.value) {
    try { await api.post('room/delete', { roomNo: no }) } catch {}
  }
  toast('已下架 ' + selected.value.length + ' 间房')
  selected.value = []
  load()
}

const gameName = (g: string) => ({ nlh: '长牌', short6: '短牌', omaha: '奥马哈', niuniu: '牛牛' } as any)[g] || g

const live = ref<Record<number, string>>({})
async function load() {
  try { rooms.value = await api.get<any[]>('rooms') } catch {}
  try {
    const t = await api.get<any[]>('tables')
    const m: Record<number, string> = {}
    for (const x of t) m[x.roomNo] = `${x.seated}/${x.seats}`
    live.value = m
  } catch {}
}

function clone(r: any) {
  const copy = { ...r, RoomNo: 3000 + Math.floor(Math.random() * 6000), Name: r.Name + ' · 副本' }
  form.value = { ...copy, Insurance: r.Insurance ? 1 : 0, IsOfficial: r.IsOfficial ? 1 : 0 }
}

function edit(r: any) {
  form.value = r ? { ...r, Insurance: r.Insurance ? 1 : 0, IsOfficial: r.IsOfficial ? 1 : 0 } : {
    RoomNo: Math.floor(Math.random() * 9000) + 1000, Name: '新房间', GameType: 'nlh', SeatCount: 6,
    SB: 20, BB: 40, Ante: 0, MinBuyin: 2000, MaxBuyin: 8000, BetLimit: 'nl', IsOfficial: 1, CardCost: 0,
    Insurance: 1, Duration: 120, Sort: 0
  }
}

async function save() {
  const f = form.value
  if (!f) return
  try {
    await api.post('room/save', { ...f, Insurance: !!f.Insurance, IsOfficial: !!f.IsOfficial })
    toast('已保存')
    form.value = null
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

async function relist(r: any) {
  try {
    await api.post('room/save', { ...r, Status: 1, Insurance: !!r.Insurance, IsOfficial: !!r.IsOfficial })
    toast('已重新上架')
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

async function del(r: any) {
  try {
    await api.post('room/delete', { roomNo: r.RoomNo })
    toast('已下架')
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

function toast(text: string, kind = 'ok') {
  window.dispatchEvent(new CustomEvent('admin-toast', { detail: { text, kind } }))
}

onMounted(load)
</script>

<style scoped>
.form-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-top: 12px; }
label { display: flex; flex-direction: column; gap: 5px; font-size: 11px; color: var(--muted); }
label .input { height: 32px; }
</style>
