<template>
  <div class="page robots-page">
    <div class="toolbar">
      <h2 class="page-title">机器人</h2>
      <span class="sum-chip">机器人池 <b>{{ bots.length }}</b></span>
      <span class="sum-chip">启用 <b>{{ enabledCount }}</b></span>
      <span class="sum-chip">已绑定房间 <b>{{ boundCount }}</b></span>
      <div class="toolbar-right">
        <select v-model.number="fillTarget" class="input" style="width: 118px">
          <option :value="2">每桌补到 2 人</option>
          <option :value="3">每桌补到 3 人</option>
          <option :value="4">每桌补到 4 人</option>
          <option :value="6">每桌补到 6 人</option>
        </select>
        <button class="btn btn-sm" :disabled="busy" @click="fillAll">一键补位</button>
        <button class="btn btn-sm btn-danger" :disabled="busy" @click="clearAll">一键清场</button>
        <button class="btn btn-sm" @click="load">刷新</button>
      </div>
    </div>

    <div class="panel">
      <div class="panel-head">
        <b>机器人池</b>
        <span class="muted">风格与激进度在线热调，下一手决策即生效</span>
      </div>
      <div class="panel-body">
        <div class="toolbar" style="padding: 0 0 10px; border: none; background: none">
          <select v-model.number="newRoom" class="input" style="width: 190px">
            <option :value="0">自动分配房间</option>
            <option v-for="r in rooms" :key="r.RoomNo" :value="r.RoomNo">{{ r.RoomNo }} · {{ r.Name }}</option>
          </select>
          <select v-model="newStyle" class="input" style="width: 110px">
            <option value="tag">紧凶</option><option value="lag">松凶</option>
            <option value="nit">岩石</option><option value="station">跟注站</option>
            <option value="maniac">疯子</option><option value="fish">娱乐</option>
          </select>
          <button class="btn btn-primary btn-sm" @click="add">+ 添加机器人</button>
        </div>
        <div v-if="!bots.length" class="robots-empty">机器人池为空，点「添加机器人」先造几个</div>
        <table v-else class="data-table">
          <thead><tr><th>UID</th><th>昵称</th><th>风格</th><th>激进度</th><th>VPIP</th><th>诈唬</th><th>绑定房间</th><th>状态</th><th>操作</th></tr></thead>
          <tbody>
            <tr v-for="b in bots" :key="b.UID">
              <td class="mono">#{{ b.UID }}</td>
              <td>{{ b.Nick }}</td>
              <td>
                <select class="input mini" :value="b.Style" @change="tune(b, 'Style', ($event.target as HTMLSelectElement).value)">
                  <option value="tag">紧凶</option><option value="lag">松凶</option>
                  <option value="nit">岩石</option><option value="station">跟注站</option>
                  <option value="maniac">疯子</option><option value="fish">娱乐</option>
                </select>
              </td>
              <td><input class="input mini" type="number" min="0" max="100" :value="b.Aggr" @change="tune(b, 'Aggr', ($event.target as HTMLInputElement).value)" /></td>
              <td><input class="input mini" type="number" min="5" max="90" :value="b.VPIP" @change="tune(b, 'VPIP', ($event.target as HTMLInputElement).value)" /></td>
              <td><input class="input mini" type="number" min="0" max="100" :value="b.Bluff" @change="tune(b, 'Bluff', ($event.target as HTMLInputElement).value)" /></td>
              <td class="mono">{{ b.RoomNo || '自动' }}</td>
              <td><span class="state-pill" :class="b.Enabled ? 'green' : 'red'">{{ b.Enabled ? '启用' : '停用' }}</span></td>
              <td>
                <div class="row-actions">
                  <button class="btn btn-xs primary" :disabled="busy" @click="seatNow(b)">坐进牌桌</button>
                  <button class="btn btn-xs" @click="assign(b)">{{ b.RoomNo ? '改绑定' : '绑定房间' }}</button>
                  <button class="btn btn-xs" :class="b.Enabled ? 'danger' : 'primary'" @click="toggle(b)">{{ b.Enabled ? '停用' : '启用' }}</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="assigning" class="modal-mask" @click.self="assigning = null">
      <div class="modal">
        <b>绑定机器人 · {{ assigning.Nick }}</b>
        <p class="muted tiny">绑定后该机器人只会在指定房间入座；设为自动则由系统按需填充官方房。</p>
        <select v-model.number="assignRoom" class="input" style="width:100%;margin:10px 0">
          <option :value="0">自动分配</option>
          <option v-for="r in rooms" :key="r.RoomNo" :value="r.RoomNo">{{ r.RoomNo }} · {{ r.Name }}</option>
        </select>
        <div class="row" style="justify-content:flex-end">
          <button class="btn" @click="assigning = null">取消</button>
          <button class="btn btn-primary" @click="saveAssign">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from '@/lib/api'

const bots = ref<any[]>([])
const rooms = ref<any[]>([])
const newRoom = ref(0)
const newStyle = ref('tag')
const assigning = ref<any>(null)
const assignRoom = ref(0)
const fillTarget = ref(3)
const busy = ref(false)
const enabledCount = computed(() => bots.value.filter(b => b.Enabled).length)
const boundCount = computed(() => bots.value.filter(b => Number(b.RoomNo) > 0).length)

async function fillAll() {
  if (busy.value) return
  busy.value = true
  try {
    const d = await api.post<any>('robots/fill', { target: fillTarget.value })
    toast(`已补位 ${d.seated} 个座位（机器人池 ${d.pool}）`)
    await load()
  } catch (e: any) {
    toast(e?.message || '补位失败', 'err')
  } finally {
    busy.value = false
  }
}

async function clearAll() {
  if (busy.value) return
  busy.value = true
  try {
    const d = await api.post<any>('robots/clear', {})
    toast(`已请下 ${d.removed} 个机器人`)
    await load()
  } catch (e: any) {
    toast(e?.message || '清场失败', 'err')
  } finally {
    busy.value = false
  }
}



// 参数在线热调：改完立刻写库 + ReloadBots，下一手决策即生效
async function tune(b: any, field: 'Style' | 'Aggr' | 'VPIP' | 'Bluff', v: string) {
  const next = { Style: b.Style, Aggr: b.Aggr, VPIP: b.VPIP, Bluff: b.Bluff }
  if (field === 'Style') next.Style = String(v)
  else next[field] = Math.round(Number(v) || 0)
  try {
    await api.post('robot/profile', { uid: b.UID, style: next.Style, aggr: next.Aggr, bluff: next.Bluff, vpip: next.VPIP })
    b.Style = next.Style
    b.Aggr = next.Aggr
    b.VPIP = next.VPIP
    b.Bluff = next.Bluff
  } catch (e: any) {
    toast(e.message, 'err')
    load()
  }
}

const styleName = (s: string) => ({ tag: '紧凶', lag: '松凶', nit: '岩石', station: '跟注站', maniac: '疯子', fish: '娱乐' } as any)[s] || s

async function load() {
  try { bots.value = await api.get<any[]>('robots') } catch {}
  try { rooms.value = await api.get<any[]>('rooms') } catch {}
}

async function add() {
  try {
    const r = await api.post<{ nick: string; seat: number }>('robot/add', { roomNo: newRoom.value, style: newStyle.value, name: '' })
    toast(r?.seat >= 0 ? `${r.nick} 已入座 ${newRoom.value} 桌（座位 ${r.seat + 1}）` : '机器人已添加')
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

function assign(b: any) {
  assigning.value = b
  assignRoom.value = b.RoomNo || 0
}

async function seatNow(b: any) {
  const roomNo = b.RoomNo || newRoom.value || Number(prompt('输入房间号') || 0)
  if (!roomNo) { toast('请先选择房间', 'err'); return }
  try {
    const r = await api.post<{ seat: number }>('robot/seat', { uid: b.UID, roomNo })
    toast(`已入座 ${roomNo} 桌（座位 ${r.seat + 1}）`)
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

async function saveAssign() {
  try {
    await api.post('robot/assign', { uid: assigning.value.UID, roomNo: assignRoom.value, enabled: true })
    toast('已绑定')
    assigning.value = null
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

async function toggle(b: any) {
  try {
    await api.post('robot/assign', { uid: b.UID, roomNo: b.RoomNo, enabled: !b.Enabled })
    toast('已更新')
    load()
  } catch (e: any) { toast(e.message, 'err') }
}

function toast(text: string, kind = 'ok') {
  window.dispatchEvent(new CustomEvent('admin-toast', { detail: { text, kind } }))
}

onMounted(load)
</script>

<style scoped>
.input.mini { width: 74px; height: 28px; padding: 0 6px; font-size: 12px; }
</style>
