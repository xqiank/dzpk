<template>
  <div class="page lobby-page private-page">
    <header class="pv-head">
      <div class="col">
        <b class="pv-title">私人局</b>
        <span class="pv-sub">自己开房 · 房号加入 · 只消耗房卡</span>
      </div>
      <button class="card-pill" @click="buyCards">
        <Icon name="card" :size="15" color="#efcf8e" />
        <b>{{ auth.profile?.cards ?? 0 }}</b>
        <span>张房卡</span>
      </button>
    </header>

    <div class="pv-actions">
      <button class="pv-action primary" @click="openCreate">
        <span class="pv-action-ico"><Icon name="plus" :size="20" color="#04220f" /></span>
        <div class="col">
          <b>创建房间</b>
          <span>设置玩法 / 盲注 / 人数，1 张房卡开一桌</span>
        </div>
        <Icon name="chevron" :size="16" color="#8fe6bd" />
      </button>
      <button class="pv-action" @click="tab = 'join'">
        <span class="pv-action-ico ghost"><Icon name="key" :size="19" color="#efcf8e" /></span>
        <div class="col">
          <b>加入房间</b>
          <span>输入好友发来的 4 位房号</span>
        </div>
        <Icon name="chevron" :size="16" color="#7fae97" />
      </button>
    </div>

    <div class="segmented dark pv-tabs">
      <button class="seg" :class="{ active: tab === 'mine' }" @click="tab = 'mine'">我的房间 ({{ rooms.length }})</button>
      <button class="seg" :class="{ active: tab === 'join' }" @click="tab = 'join'">加入房间</button>
    </div>

    <!-- 我的房间 -->
    <template v-if="tab === 'mine'">
      <div v-if="loading && !rooms.length" class="pv-empty">
        <Icon name="refresh" :size="26" color="#7fd8ae" />
        <b>正在读取我的房间…</b>
      </div>
      <div v-else-if="!rooms.length" class="pv-empty">
        <div class="doc-box"><Icon name="card" :size="28" color="#9fe6c4" /></div>
        <b>还没有私人局</b>
        <span class="sm muted">开一桌发给好友房号，或输入房号加入别人的桌</span>
        <div class="row gap-8" style="margin-top:14px">
          <button class="btn btn-sm btn-primary" @click="openCreate">创建房间</button>
          <button class="btn btn-sm btn-ghost" @click="tab = 'join'">加入房间</button>
        </div>
      </div>
      <div v-else class="pv-list">
        <div v-for="r in rooms" :key="r.roomNo" class="pv-room-wrap">
        <article
          class="poker-table-card pv-card"
          :class="[gameModeClass(r), occClass(r), { mine: r.hostUid === auth.profile?.uid }]"
          @click="enterRoom(r)"
        >
          <div class="card-occupancy-bar" :class="occClass(r)" aria-hidden="true">
            <span class="card-occupancy-fill" :style="{ height: Math.round((r.online / r.seatCount) * 100) + '%' }"></span>
          </div>
          <div class="card-watermark-overlay"></div>
          <div class="card-tether-watermark" aria-hidden="true">
            <svg viewBox="0 0 100 100" fill="currentColor"><path d="M50 0C22.4 0 0 22.4 0 50s22.4 50 50 50 50-22.4 50-50S77.6 0 50 0zm21.6 28.5H54.2v6.2h17.4v6.2H54.2v5.3c15.2.8 26.6 4.3 26.6 8.5 0 4.2-11.4 7.7-26.6 8.5v16.1H45.8V63.2c-15.2-.8-26.6-4.3-26.6-8.5 0-4.2 11.4-7.7 26.6-8.5v-5.3H28.4v-6.2h17.4v-6.2H28.4v-6.2h43.2v6.2zm-21.6 30c12.3 0 22.5-2.2 24.8-5.3-2.3-3.1-12.5-5.3-24.8-5.3s-22.5 2.2-24.8 5.3c2.3 3.1 12.5 5.3 24.8 5.3z"></path></svg>
          </div>
          <div class="card-donut-col">
            <div class="table-topdown-box">
              <svg class="seat-arc-svg" viewBox="-6 -6 67 51" aria-hidden="true">
                <path class="seat-arc-glow" d="M 16.26 -0.49 A 20.25 20.25 0 0 1 19.50 -0.75 L 35.50 -0.75 A 20.25 20.25 0 0 1 38.74 39.49"></path>
                <path class="seat-arc-line" d="M 16.26 -0.49 A 20.25 20.25 0 0 1 19.50 -0.75 L 35.50 -0.75 A 20.25 20.25 0 0 1 38.74 39.49" stroke="url(#lobbySeatArcGrad)"></path>
              </svg>
              <span
                v-for="(d, i) in seatDots(r)"
                :key="i"
                class="s-dot"
                :class="{ on: i < r.online }"
                :style="{ left: d.x + 'px', top: d.y + 'px' }"
              ></span>
              <b class="table-seat-num">{{ r.online }}/{{ r.seatCount }}</b>
            </div>
          </div>
          <div class="card-info-col">
            <div class="room-title-line">
              <span class="room-title">{{ r.name }}</span>
              <span v-if="r.needsCode" class="pv-lock" title="需要密码">🔒</span>
              <span v-if="r.hostUid === auth.profile?.uid" class="pv-mine-tag">我的</span>
            </div>
            <div class="room-blinds-line">
              <strong class="blinds-num">{{ chips(r.sb) }}/{{ chips(r.bb) }}</strong>
              <span class="room-safety-badge" aria-hidden="true">
                <svg viewBox="0 0 24 26" width="16" height="18"><path d="M12 1 2.5 4.7v7.1c0 5.8 3.8 10.9 9.5 13.2 5.7-2.3 9.5-7.4 9.5-13.2V4.7L12 1Zm0 3.1 6.8 2.6v5.1c0 4.4-2.6 8.5-6.8 10.4-4.2-1.9-6.8-6-6.8-10.4V6.7L12 4.1Z" fill="currentColor" fill-rule="evenodd"></path><path d="M12 4.1v8H5.2V6.7L12 4.1Zm0 8h6.8c-.1 4.3-2.7 8.3-6.8 10.2V12.1Z" fill="currentColor"></path></svg>
              </span>
              <span class="room-badge" :class="'room-badge--' + r.gameType"><span>{{ r.gameName }}</span></span>
            </div>
            <div class="room-meta-line">
              <span class="meta-time">
                <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                <span>{{ remain(r.expireAt) }}</span>
              </span>
              <span class="meta-id">ID: {{ r.roomNo }}</span>
            </div>
          </div>
          <div class="card-buyin-col">
            <span class="buyin-title">最小买入</span>
            <div class="buyin-amount-wrap">
              <span class="usdt-green-icon">₮</span><b class="buyin-value">{{ chips(r.minBuyin) }}</b>
            </div>
            <span class="vpip-text">每手 {{ r.cardCost }} 张房卡</span>
          </div>
        </article>
        <div class="pv-row-actions">
          <button class="pv-mini primary" :disabled="busy" @click="enterRoom(r)">进入房间</button>
          <button class="pv-mini" @click="copyNo(r.roomNo)">复制房号</button>
          <button v-if="r.hostUid === auth.profile?.uid" class="pv-mini danger" :disabled="busy" @click="closeRoom(r)">解散</button>
        </div>
        </div>
      </div>
      <p class="pv-tip">私人局只消耗房卡，不计入现金场；房间到期后会自动解散。</p>
    </template>

    <!-- 加入房间 -->
    <template v-else>
      <div class="pv-join">
        <b class="pv-join-title">输入房号加入</b>
        <span class="tiny muted">房号由房主分享，4 位数字</span>
        <div class="pv-code-row">
          <input
            v-model="code"
            class="pv-code"
            inputmode="numeric"
            maxlength="4"
            placeholder="0000"
            @keyup.enter="findRoom"
          />
          <button class="btn btn-primary pv-join-btn" :disabled="busy || code.length < 4" @click="findRoom">查找</button>
        </div>
        <div v-if="findError && /密码/.test(findError)" class="pv-code-row" style="margin-top: 10px">
          <input v-model="pwd" class="pv-code" maxlength="8" placeholder="房间密码" style="letter-spacing: 4px; font-size: 17px" />
          <button class="btn btn-primary pv-join-btn" :disabled="busy" @click="findRoom">解锁</button>
        </div>
        <div v-if="found" class="pv-found">
          <div class="col" style="min-width:0">
            <b>{{ found.name }}</b>
            <span class="tiny muted">{{ found.gameName }} · {{ chips(found.sb) }}/{{ chips(found.bb) }} · {{ found.online }}/{{ found.seatCount }} 人在座</span>
          </div>
          <button class="btn btn-sm btn-primary" :disabled="busy" @click="enterRoom(found)">进入</button>
        </div>
        <p v-if="findError" class="pv-error">{{ findError }}</p>
      </div>

      <div v-if="recent.length" class="pv-recent">
        <b class="pv-recent-title">最近加入</b>
        <button v-for="r in recent" :key="r.roomNo" class="pv-recent-chip" @click="code = String(r.roomNo); findRoom()">
          <b>#{{ r.roomNo }}</b><span>{{ r.name }}</span>
        </button>
      </div>
    </template>

  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Icon from '@/components/Icon.vue'
import { useAuth } from '@/stores/auth'
import { useUI } from '@/stores/ui'
import { useTable } from '@/stores/table'
import { socket } from '@/lib/ws'
import { chips, remain } from '@/lib/format'
import { sfx } from '@/lib/sfx'
import '@/styles/lobby-official.css'

const router = useRouter()
const auth = useAuth()
const ui = useUI()
const table = useTable()

const gameTypes = [
  { key: 'nlh', label: '长牌' },
  { key: 'short6', label: '短牌' },
  { key: 'omaha', label: '奥马哈' },
  { key: 'niuniu', label: '牛牛' }
]

const tab = ref<'mine' | 'join'>('mine')
const rooms = ref<any[]>([])
const loading = ref(false)
const busy = ref(false)
const code = ref('')
const pwd = ref('')
const found = ref<any>(null)
const findError = ref('')
const recent = ref<any[]>(JSON.parse(localStorage.getItem('u5_private_recent') || '[]'))

const OFFICIAL_SEAT_DOTS = [
  { x: 16.26, y: -0.49 }, { x: 38.74, y: -0.49 }, { x: 55.75, y: 19.5 },
  { x: 38.74, y: 39.49 }, { x: 16.26, y: 39.49 }, { x: -0.75, y: 19.5 }
]
function seatDots(r: any) {
  const n = Math.max(2, Math.min(9, Number(r?.seatCount) || 6))
  if (n === 6) return OFFICIAL_SEAT_DOTS
  const cx = 27.5, cy = 19.5, rx = 28.25, ry = 20
  return Array.from({ length: n }, (_, i) => {
    const a = -Math.PI / 2 + (i * Math.PI * 2) / n + Math.PI / n
    return { x: +(cx + rx * Math.cos(a)).toFixed(2), y: +(cy + ry * Math.sin(a)).toFixed(2) }
  })
}
function occClass(r: any) {
  const p = (Number(r?.online) || 0) / Math.max(1, Number(r?.seatCount) || 6)
  if (p >= 0.8) return 'occ-hot'
  if (p >= 0.4) return 'occ-mid'
  return 'occ-low'
}
function gameModeClass(r: any) {
  return 'mode-' + (r?.gameType || 'nlh')
}


let poll = 0

async function loadRooms() {
  loading.value = true
  try {
    const list = await socket.call<any[]>('lobby.privateRooms')
    rooms.value = Array.isArray(list) ? list : []
  } catch {
    rooms.value = []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  sfx.tap()
  router.push('/private/create')
}


async function findRoom() {
  const no = Number(code.value)
  if (!no || String(no).length < 4) {
    findError.value = '请输入 4 位房号'
    return
  }
  findError.value = ''
  found.value = null
  busy.value = true
  try {
    found.value = await socket.call('lobby.findRoom', { roomNo: no, code: pwd.value.trim() })
    sfx.tap()
  } catch (e: any) {
    found.value = null
    findError.value = e.message
  } finally {
    busy.value = false
  }
}

async function enterRoom(r: any) {
  if (busy.value) return
  busy.value = true
  try {
    const seated = r.online < r.seatCount
    await table.join(r.roomNo, seated ? r.minBuyin : 0, -1, !seated, pwd.value.trim())
    remember(r)
    sfx.chip()
    router.push(`/table/${r.roomNo}${seated ? '' : '?observe=1'}`)
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    busy.value = false
  }
}

function remember(r: any) {
  const next = [{ roomNo: r.roomNo, name: r.name }, ...recent.value.filter((x: any) => x.roomNo !== r.roomNo)].slice(0, 6)
  recent.value = next
  localStorage.setItem('u5_private_recent', JSON.stringify(next))
}

async function closeRoom(r: any) {
  if (busy.value) return
  busy.value = true
  try {
    await socket.call('lobby.closeRoom', { roomNo: r.roomNo })
    ui.ok('房间已解散')
    await loadRooms()
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    busy.value = false
  }
}

function copyNo(no: number) {
  navigator.clipboard?.writeText(String(no)).then(() => ui.ok('房号已复制')).catch(() => ui.err('复制失败'))
  sfx.tap()
}

async function buyCards() {
  try {
    const d = await socket.call<{ cards: number }>('user.buyCards', { count: 10 })
    auth.applyBalance({ cards: d.cards })
    ui.ok('已购买 10 张房卡')
    sfx.chip()
  } catch (e: any) {
    ui.err(e.message)
  }
}

onMounted(async () => {
  await auth.refresh()
  loadRooms()
  poll = window.setInterval(loadRooms, 5000)
})

onUnmounted(() => {
  if (poll) window.clearInterval(poll)
})
</script>

<style scoped>
.private-page { padding-top: 12px; }
.pv-head { display: flex; align-items: center; justify-content: space-between; gap: 10px; margin-bottom: 12px; }
.pv-title { font-size: 21px; font-weight: 900; letter-spacing: .5px; color: #f1f8f5; }
.pv-sub { font-size: 12px; color: var(--muted); }
.card-pill { display: inline-flex; align-items: center; gap: 4px; height: 32px; padding: 0 10px; border-radius: 16px;
  background: rgba(10, 32, 22, .9); border: 1px solid rgba(214, 176, 92, .3); font-size: 11px; color: #cfd8c6; }
.card-pill b { font-size: 13px; color: #efcf8e; }
.pv-actions { display: flex; flex-direction: column; gap: 8px; margin-bottom: 12px; }
.pv-action { display: flex; align-items: center; gap: 10px; padding: 12px 14px; border-radius: 16px; text-align: left;
  background: linear-gradient(150deg, rgba(19, 48, 33, .9), rgba(8, 26, 17, .94)); border: 1px solid var(--line); }
.pv-action.primary { border-color: rgba(46, 224, 122, .38); background: linear-gradient(150deg, rgba(24, 66, 45, .95), rgba(9, 30, 20, .95)); }
.pv-action .col { flex: 1; min-width: 0; gap: 2px; }
.pv-action b { font-size: 14px; color: #eafaf1; }
.pv-action span { font-size: 11.5px; color: var(--muted); }
.pv-action-ico { width: 38px; height: 38px; border-radius: 12px; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(180deg, #55ea9c, #12b767); box-shadow: 0 6px 14px rgba(18, 183, 103, .28); }
.pv-action-ico.ghost { background: rgba(214, 176, 92, .16); border: 1px solid rgba(214, 176, 92, .34); box-shadow: none; }
.pv-tabs { margin-bottom: 12px; }
.pv-list { display: flex; flex-direction: column; gap: 10px; }
.pv-room { padding: 13px 14px; border-radius: 16px; background: linear-gradient(150deg, rgba(17, 44, 31, .92), rgba(7, 22, 15, .95));
  border: 1px solid rgba(243, 204, 122, .22); }
.pv-room.mine { border-color: rgba(46, 224, 122, .34); }
.pv-room-top { display: flex; align-items: flex-start; justify-content: space-between; gap: 8px; }
.pv-room-name { font-size: 14.5px; font-weight: 800; color: #f0f9f4; }
.pv-room-meta { font-size: 11.5px; color: var(--muted); margin-top: 2px; }
.pv-no { flex: none; font-size: 12px; font-weight: 800; color: #efcf8e; padding: 2px 8px; border-radius: 8px;
  background: rgba(214, 176, 92, .14); border: 1px solid rgba(214, 176, 92, .3); }
.pv-room-stats { display: flex; flex-wrap: wrap; gap: 6px 12px; margin: 9px 0 10px; font-size: 11px; color: #9fb8a9; }
.pv-room-stats .dot { width: 6px; height: 6px; border-radius: 50%; display: inline-block; margin-right: 4px; background: #2ee07a; }
.pv-room-btns { display: flex; gap: 8px; }
.pv-room-btns .btn { height: 34px; }
.btn.danger { color: #ffb3a7; }
.pv-tip { margin-top: 12px; font-size: 11px; color: var(--muted-2); text-align: center; }
.pv-empty { display: flex; flex-direction: column; align-items: center; gap: 6px; padding: 34px 16px; text-align: center;
  border-radius: 18px; background: rgba(9, 28, 19, .7); border: 1px dashed rgba(120, 200, 160, .22); }
.pv-empty b { font-size: 14px; color: #e6f4ee; }
.pv-join { padding: 16px; border-radius: 18px; background: linear-gradient(150deg, rgba(19, 48, 33, .9), rgba(8, 26, 17, .94));
  border: 1px solid var(--line); display: flex; flex-direction: column; gap: 6px; }
.pv-join-title { font-size: 15px; color: #eafaf1; }
.pv-code-row { display: flex; gap: 8px; margin-top: 8px; }
.pv-code { flex: 1; min-width: 0; height: 52px; border-radius: 14px; text-align: center; letter-spacing: 10px;
  font-size: 22px; font-weight: 800; color: #f2ffe8; background: rgba(4, 16, 10, .7); border: 1px solid rgba(122, 196, 156, .28); }
.pv-join-btn { width: 92px; height: 52px; }
.pv-found { display: flex; align-items: center; justify-content: space-between; gap: 10px; margin-top: 12px; padding: 10px 12px;
  border-radius: 12px; background: rgba(12, 44, 30, .8); border: 1px solid rgba(46, 224, 122, .3); }
.pv-error { font-size: 12px; color: #ffb3a7; margin-top: 8px; }
.pv-recent { margin-top: 16px; display: flex; flex-wrap: wrap; gap: 8px; }
.pv-recent-title { width: 100%; font-size: 12px; color: var(--muted); margin-bottom: 2px; }
.pv-recent-chip { display: inline-flex; align-items: center; gap: 6px; height: 32px; padding: 0 12px; border-radius: 16px;
  background: rgba(10, 32, 22, .9); border: 1px solid var(--line); font-size: 11.5px; color: #cfe0d4; }
.pv-recent-chip b { color: #efcf8e; }

/* 私人局：官方房间卡之上的操作行 / 标记 */
.pv-room-wrap { display: flex; flex-direction: column; }
.pv-row-actions { display: flex; gap: 6px; padding: 8px 6px 2px; }
.pv-mini { height: 28px; padding: 0 10px; border-radius: 9px; font-size: 11.5px; font-weight: 700;
  color: #cfe6d8; background: rgba(10, 32, 22, .9); border: 1px solid rgba(120, 200, 160, .24); }
.pv-mini.primary { color: #04220f; background: linear-gradient(180deg, #55ea9c, #12b767); border-color: transparent; }
.pv-mini.danger { color: #ffb3a7; }
.pv-lock { font-size: 11px; }
.pv-mine-tag { flex: none; font-size: 9.5px; font-weight: 800; color: #2b1400; padding: 1px 6px; border-radius: 6px;
  background: linear-gradient(180deg, #ffe08a, #e0a83a); }
.pv-card { cursor: pointer; }
.pv-card .card-info-col { gap: 3px; }
</style>
