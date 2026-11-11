<template>
  <div class="profile-page user-profile-view">
    <div class="mobile-lobby-shell mobile-profile-shell">
      <header class="mobile-lobby-toolbar profile-toolbar">
        <div class="header-action-icons">
          <button class="icon-touch-btn" title="消息通知" @click="openNotices">
            <Icon name="bell" :size="17" />
            <em v-if="unread" class="bell-badge-pill">{{ unread > 99 ? '99+' : unread }}</em>
          </button>
          <button class="icon-touch-btn" title="设置" @click="deskOpen = true">
            <Icon name="settings" :size="17" />
          </button>
          <button class="icon-touch-btn" title="平台说明与会员准则" @click="helpOpen = true">
            <Icon name="help" :size="17" />
          </button>
        </div>
      </header>

      <main class="profile-content">
        <section class="mine-user-card">
          <div class="user-avatar-wrap" @click="editOpen = true">
            <img class="user-main-avatar" :src="avatarUrl(p?.avatar)" :alt="`${p?.nick || '玩家'}头像`" />
          </div>
          <div class="user-info-meta">
            <div class="user-nickname-row">
              <span class="user-nickname-text">{{ p?.nick || '玩家' }}</span>
              <span class="user-tag-icon">
                <Icon name="star" :size="11" color="#ffd700" />
                <span>VIP{{ p?.vip ? ' ' + p.vip : '' }}</span>
              </span>
            </div>
            <div class="user-id-row">
              <span class="user-id-text">ID: {{ p?.uid || '----' }}</span>
              <button class="copy-id-btn" title="复制ID" @click="copyUid">
                <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
              </button>
            </div>
            <div class="user-signature-row">
              <span class="user-signature-quote">“</span>
              <span class="user-signature-text">{{ p?.signature || '点击设置个性签名…' }}</span>
            </div>
          </div>
          <button class="profile-link-btn" title="编辑个人资料" @click="editOpen = true">
            <span>个人资料</span><span class="chevron-right">›</span>
          </button>
        </section>

        <div class="profile-scroll-body">
          <section class="mine-wallet-hero-card vip-metal-card">
            <div class="card-chip-hologram">
              <div class="smart-emv-chip">
                <span class="chip-line horizontal"></span>
                <span class="chip-line vertical"></span>
              </div>
              <div class="card-brand-tag">
                <span class="brand-crown-icon" aria-hidden="true"><Icon name="crown" :size="13" color="#ffd700" /></span>
                <span class="brand-text">U5 VIP CLUB</span>
              </div>
            </div>
            <div class="card-main-content">
              <div class="card-top-row">
                <div class="asset-title-group">
                  <span class="usdt-gold-badge">₮</span>
                  <span class="asset-label-text">账户总资产 (USDT)</span>
                </div>
                <button class="room-card-badge-btn" title="查看房卡与明细" @click="buyCards">
                  <span class="card-badge-icon" aria-hidden="true"><Icon name="card" :size="13" /></span>
                  <span>房卡 {{ p?.cards ?? 0 }} 张</span>
                  <span class="badge-arrow">›</span>
                </button>
              </div>
              <div class="card-balance-row" title="点击查看资金明细" @click="router.push('/bills')">
                <div class="balance-display-main">
                  <span class="currency-symbol">₮</span>
                  <strong class="balance-amount">{{ money(p?.balance) }}</strong>
                </div>
                <div class="balance-sub-hint">
                  <span class="live-dot"></span><span>1:1 链上即时承兑</span>
                </div>
              </div>
              <div class="card-action-pills">
                <button class="action-pill-btn btn-primary-deposit" @click="rechargeOpen = true">
                  <Icon name="plus" :size="14" /><span>立即充值</span>
                </button>
                <button class="action-pill-btn btn-secondary-withdraw" @click="withdrawOpen = true">
                  <Icon name="arrowDown" :size="14" /><span>极速提现</span>
                </button>
                <button class="action-pill-btn btn-tertiary-records" @click="router.push('/bills')">
                  <Icon name="doc" :size="13" /><span>资产明细</span>
                </button>
              </div>
            </div>
          </section>

          <section class="mine-quick-grid">
            <button class="quick-grid-card" @click="router.push('/bills')">
              <div class="grid-card-icon"><Icon name="doc" :size="22" /></div>
              <span class="grid-card-label">账单</span>
            </button>
            <button class="quick-grid-card" @click="router.push('/records')">
              <div class="grid-card-icon"><Icon name="star" :size="22" /></div>
              <span class="grid-card-label">收藏</span>
            </button>
            <button class="quick-grid-card customer-service-btn" @click="supportOpen = true">
              <div class="grid-card-icon cs-icon-wrapper">
                <Icon name="help" :size="22" />
                <span class="live-cs-pulse"></span>
              </div>
              <span class="grid-card-label">客服</span>
            </button>
            <button class="quick-grid-card" @click="openInvite">
              <div class="grid-card-icon"><Icon name="users" :size="22" /></div>
              <span class="grid-card-label">我的邀请</span>
            </button>
          </section>

          <section class="mine-stats-section">
            <h2 class="stats-section-title">数据</h2>
            <div class="stats-card-main">
              <div class="stats-game-tabs">
                <button
                  v-for="g in games"
                  :key="g.key"
                  class="game-tab-btn"
                  :class="{ active: gameTab === g.key }"
                  @click="gameTab = g.key"
                >{{ g.label }}</button>
              </div>
              <div class="stats-data-row">
                <div class="stats-col stats-col-left">
                  <div class="stat-unit">
                    <b class="stat-num">{{ money(stat.net) }}</b>
                    <span class="stat-label">总输赢</span>
                  </div>
                  <div class="stat-unit">
                    <b class="stat-num">{{ stat.hands || 0 }}</b>
                    <span class="stat-label">手牌数</span>
                  </div>
                </div>
                <div class="stats-col stats-col-center">
                  <div class="donut-chart-container">
                    <svg class="donut-svg" viewBox="0 0 80 80">
                      <circle class="donut-bg" cx="40" cy="40" r="32"></circle>
                      <circle class="donut-fg" cx="40" cy="40" r="32" :stroke-dasharray="donutLen" :stroke-dashoffset="donutOffset"></circle>
                    </svg>
                    <div class="donut-text-center"><b class="donut-percent-num">{{ raiseRate }}%</b></div>
                  </div>
                  <span class="stat-label donut-label">加注率</span>
                </div>
                <div class="stats-col stats-col-right">
                  <div class="stat-unit">
                    <b class="stat-num">{{ vpipRate }}%</b>
                    <span class="stat-label">入池率</span>
                  </div>
                  <div class="stat-unit">
                    <b class="stat-num">{{ money(avgWin, { sign: false }) }}</b>
                    <span class="stat-label">手均盈利</span>
                  </div>
                </div>
              </div>
            </div>
          </section>

          <section class="mine-menu-group">
            <div class="menu-list-item" @click="payOpen = true">
              <div class="item-left-label">
                <span class="menu-item-icon icon-password"><Icon name="lock" :size="16" /></span>
                <span class="item-title">支付密码</span>
                <span v-if="!p?.hasPayPwd" class="info-alert-circle">!</span>
              </div>
              <span class="item-chevron">›</span>
            </div>
            <div class="menu-list-item" @click="deskOpen = true">
              <div class="item-left-label">
                <span class="menu-item-icon icon-settings"><Icon name="settings" :size="16" /></span>
                <span class="item-title">桌面设置</span>
                <span class="pay-pwd-set-tag">{{ deskSub }}</span>
              </div>
              <span class="item-chevron">›</span>
            </div>
            <div class="menu-list-item" @click="addToDesktop">
              <div class="item-left-label">
                <span class="menu-item-icon icon-desktop"><Icon name="grid" :size="16" /></span>
                <span class="item-title">添加到桌面</span>
              </div>
              <span class="item-chevron">›</span>
            </div>
            <div class="menu-list-item logout-item" @click="logout">
              <div class="item-left-label">
                <span class="menu-item-icon icon-logout"><Icon name="logout" :size="16" /></span>
                <span class="item-title">退出当前账号</span>
              </div>
              <span class="item-chevron">›</span>
            </div>
          </section>
        </div>
      </main>
    </div>

    <!-- 充币 / 提币 / 充提记录（官方 UsdtModal） -->
    <UsdtModal :open="walletOpen" :initial-tab="walletTab" @close="walletOpen = false"
               @balance-updated="auth.refresh()" />

    <!-- 资料 -->
    <div v-if="editOpen" class="mask" @click.self="editOpen = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:14px"><b class="h2">个人资料</b>
          <button class="close-x" @click="editOpen = false"><Icon name="close" :size="18" /></button></div>
        <div class="avatar-row">
          <Avatar :src="form.avatar" :size="64" gold />
          <div class="row gap-6" style="flex-wrap:wrap">
            <button v-for="i in 8" :key="i" class="av-pick" :class="{ on: form.avatar === `/avatars/av${i}.svg` }"
                    @click="form.avatar = `/avatars/av${i}.svg`">
              <img :src="`/avatars/av${i}.svg`" alt="" />
            </button>
          </div>
        </div>
        <input v-model="form.nick" class="field" style="margin-top:12px" placeholder="昵称" maxlength="16" />
        <input v-model="form.signature" class="field" style="margin-top:10px" placeholder="个性签名" maxlength="32" />
        <button class="btn btn-primary btn-block" style="height:48px;margin-top:14px" @click="saveProfile">保存</button>
      </div>
    </div>

    <!-- 支付密码 -->
    <div v-if="payOpen" class="mask" @click.self="payOpen = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:14px"><b class="h2">设置支付密码</b>
          <button class="close-x" @click="payOpen = false"><Icon name="close" :size="18" /></button></div>
        <input v-model="payPwd" class="field" type="password" maxlength="6" placeholder="6 位数字支付密码" />
        <button class="btn btn-primary btn-block" style="height:48px;margin-top:14px" @click="savePay">保存</button>
      </div>
    </div>

    <!-- 帮助 -->
    <div v-if="inviteOpen" class="mask" @click.self="inviteOpen = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:10px">
          <b class="h2">我的邀请</b>
          <button class="close-x" @click="inviteOpen = false"><Icon name="close" :size="18" /></button>
        </div>
        <div class="invite-card">
          <span class="tiny muted">我的邀请码</span>
          <b class="invite-code">{{ p?.uid }}</b>
          <button class="btn btn-primary btn-sm" @click="copyInvite">复制邀请链接</button>
        </div>
        <div class="row-between" style="margin:14px 0 6px">
          <b class="h3">我的下级</b>
          <span class="tiny muted">{{ agentSubs.length }} 人</span>
        </div>
        <div v-if="!agentSubs.length" class="tiny muted" style="padding:10px 0">还没有下级，把邀请链接发给牌友即可绑定。</div>
        <div v-for="a in agentSubs" :key="a.uid" class="list-item" style="padding:8px 0">
          <div class="row gap-10">
            <Avatar :src="a.avatar" :size="32" />
            <div class="col"><b>{{ a.nick }}</b><span class="tiny muted">ID: {{ a.uid }}</span></div>
          </div>
          <span class="tag tag-gold">L{{ a.level }} · 返佣 {{ a.rebate }}‰</span>
        </div>
      </div>
    </div>

    <!-- 桌面设置 -->
    <div v-if="deskOpen" class="mask" @click.self="deskOpen = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:12px">
          <b class="h2">桌面设置</b>
          <button class="close-x" @click="deskOpen = false"><Icon name="close" :size="18" /></button>
        </div>
        <div class="pref-title">牌桌呢布</div>
        <div class="cloth-grid">
          <button v-for="c in CLOTHS" :key="c.key" class="cloth-item" :class="{ on: cloth === c.key }" @click="setCloth(c.key)">
            <img :src="c.url" alt="" />
            <span>{{ c.name }}</span>
          </button>
        </div>
        <div class="pref-title" style="margin-top:12px">牌背</div>
        <div class="back-grid">
          <button v-for="b in CARD_BACKS" :key="b.key" class="back-item" :class="{ on: cardBack === b.key }" @click="setCardBack(b.key)">
            <img :src="b.url" alt="" />
            <span>{{ b.name }}</span>
          </button>
        </div>
        <div class="pref-row">
          <div class="col"><b>音效</b><span class="row-sub">发牌、筹码、行动、胜利</span></div>
          <button class="switch" :class="{ on: ui.sound }" @click="ui.toggleSound()"><i></i></button>
        </div>
        <div class="pref-row">
          <div class="col"><b>语音播报</b><span class="row-sub">跟注 / 加注 / 过牌 / 轮到你了</span></div>
          <button class="switch" :class="{ on: ui.voice }" @click="ui.toggleVoice()"><i></i></button>
        </div>
        <div class="pref-row">
          <div class="col"><b>筹码显示</b><span class="row-sub">实际筹码 / 大盲倍数</span></div>
          <button class="switch" :class="{ on: chipUnit === 'bb' }" @click="toggleChipUnit"><i></i></button>
        </div>
      </div>
    </div>

    <!-- 通知中心 -->
    <div v-if="noticeOpen" class="mask" @click.self="noticeOpen = false">
      <div class="sheet notice-sheet">
        <div class="row-between" style="margin-bottom:12px">
          <b class="h2">消息通知</b>
          <div class="row gap-8">
            <button class="mini-btn" :disabled="!unread" @click="markRead">全部已读</button>
            <button class="close-x" @click="noticeOpen = false"><Icon name="close" :size="18" /></button>
          </div>
        </div>
        <div v-if="!notices.length" class="notice-empty">
          <Icon name="bell" :size="30" color="#4d6b5c" />
          <span class="sm muted">暂时没有新消息</span>
        </div>
        <div v-for="n in notices" :key="n.id" class="notice-item" :class="{ unread: !n.read }">
          <span class="notice-ico" :class="n.type">
            <Icon :name="iconOf(n.type)" :size="15" color="#0b1e15" />
          </span>
          <div class="col" style="flex:1;min-width:0">
            <div class="row-between">
              <b class="notice-title">{{ n.title }}</b>
              <span class="tiny muted">{{ ago(n.createdAt) }}</span>
            </div>
            <span class="notice-body">{{ n.body }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 平台说明与会员准则（官方 ProfileHelpModal） -->
    <HelpGuideModal :open="helpOpen" @close="helpOpen = false" />
    <!-- 在线客服（官方 SupportChatModal） -->
    <SupportChatModal :visible="supportOpen" @close="supportOpen = false" @openLogin="router.push('/login')" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
import Avatar from '@/components/Avatar.vue'
import Icon from '@/components/Icon.vue'
import HelpGuideModal from '@/components/HelpGuideModal.vue'
import SupportChatModal from '@/components/SupportChatModal.vue'
import UsdtModal from '@/components/UsdtModal.vue'
import { useAuth } from '@/stores/auth'
import { useUI } from '@/stores/ui'
import { socket } from '@/lib/ws'
import { money } from '@/lib/format'
import { CLOTHS, CARD_BACKS, clothOf, cardBackOf } from '@/lib/tableSkins'
import { sfx } from '@/lib/sfx'
import { avatarUrl } from '@/lib/format'
import '@/styles/profile-official.css'

const router = useRouter()
const auth = useAuth()
const ui = useUI()

const p = computed(() => auth.profile)
const games = [
  { key: 'nlh', label: '长牌' },
  { key: 'short6', label: '短牌' },
  { key: 'omaha', label: '奥马哈' },
  { key: 'niuniu', label: '牛牛' }
]
const gameTab = ref('nlh')
const stat = computed(() => (p.value?.stats || {})[gameTab.value] || {})
const raiseRate = computed(() => {
  const s: any = stat.value
  const hands = Number(s.hands || 0)
  if (!hands) return 0
  return Math.min(100, Math.round((Number(s.raises || 0) / hands) * 100))
})
const vpipRate = computed(() => {
  const s: any = stat.value
  const hands = Number(s.hands || 0)
  if (!hands) return 0
  return Math.min(100, Math.round((Number(s.vpip || 0) / hands) * 100))
})
const avgWin = computed(() => {
  const s: any = stat.value
  const hands = Number(s.hands || 0)
  if (!hands) return 0
  return Math.round(Number(s.net || 0) / hands)
})
const donutLen = 201
const donutOffset = computed(() => Math.round(donutLen - (donutLen * raiseRate.value) / 100))
function addToDesktop() {
  const isStandalone = window.matchMedia('(display-mode: standalone)').matches
  if (isStandalone) {
    ui.ok('已经在桌面模式运行')
    return
  }
  ui.toast('在浏览器菜单里选择「添加到主屏幕」即可', 'info', 3600)
}
const quicks = [
  { icon: 'doc', label: '账单', path: '/bills' },
  { icon: 'chart', label: '战绩', path: '/records' },
  { icon: 'chat', label: '客服', dot: true },
  { icon: 'users', label: '我的邀请' }
]
const walletOpen = ref(false)
const walletTab = ref<'deposit' | 'withdraw' | 'orders'>('deposit')
const rechargeOpen = computed({ get: () => walletOpen.value && walletTab.value === 'deposit', set: (v: boolean) => { if (v) { walletTab.value = 'deposit'; walletOpen.value = true } else if (walletTab.value === 'deposit') walletOpen.value = false } })
const withdrawOpen = computed({ get: () => walletOpen.value && walletTab.value === 'withdraw', set: (v: boolean) => { if (v) { walletTab.value = 'withdraw'; walletOpen.value = true } else if (walletTab.value === 'withdraw') walletOpen.value = false } })
const editOpen = ref(false)
const payOpen = ref(false)
const helpOpen = ref(false)
const supportOpen = ref(false)
const inviteOpen = ref(false)
const agentSubs = ref<any[]>([])
const deskOpen = ref(false)
const cloth = ref(localStorage.getItem('u5_cloth') || 'green')
const cardBack = ref(localStorage.getItem('u5_cardback') || 'emerald')
const chipUnit = ref(localStorage.getItem('u5_chips_unit') === 'bb' ? 'bb' : 'chips')
const deskSub = computed(() => {
  const c = clothOf(cloth.value).name
  const b = cardBackOf(cardBack.value).name
  return `${c} · ${b} · ${chipUnit.value === 'bb' ? 'BB' : '实际'}筹码`
})
function setCloth(k: string) {
  cloth.value = k
  localStorage.setItem('u5_cloth', k)
  ui.toast('牌桌呢布已切换', 'info', 1500)
}
function setCardBack(k: string) {
  cardBack.value = k
  localStorage.setItem('u5_cardback', k)
  ui.toast('牌背已切换', 'info', 1500)
}
function toggleChipUnit() {
  chipUnit.value = chipUnit.value === 'bb' ? 'chips' : 'bb'
  localStorage.setItem('u5_chips_unit', chipUnit.value)
}
const noticeOpen = ref(false)
const notices = ref<any[]>([])
const unread = ref(0)
const amount = ref(10000)
const payPwd = ref('')
const form = reactive({ nick: '', avatar: '', signature: '' })

function iconOf(t: string) {
  return t === 'wallet' ? 'doc' : t === 'club' ? 'shield' : t === 'agent' ? 'users' : t === 'game' ? 'cards' : 'bell'
}
function ago(ms: number) {
  const d = Math.max(0, Date.now() - (ms || 0))
  if (d < 60e3) return '刚刚'
  if (d < 3600e3) return Math.floor(d / 60e3) + ' 分钟前'
  if (d < 86400e3) return Math.floor(d / 3600e3) + ' 小时前'
  if (d < 7 * 86400e3) return Math.floor(d / 86400e3) + ' 天前'
  const t = new Date(ms)
  return `${t.getMonth() + 1}/${t.getDate()}`
}
async function loadNotices(silent = true) {
  try {
    const d = await socket.call<{ list: any[]; unread: number }>('user.notices', {})
    notices.value = d.list || []
    unread.value = d.unread || 0
  } catch {
    if (!silent) ui.err('通知加载失败')
  }
}
async function openNotices() {
  noticeOpen.value = true
  sfx.tap()
  await loadNotices(false)
}
async function markRead() {
  try {
    await socket.call('user.noticesRead', {})
    notices.value = notices.value.map(n => ({ ...n, read: true }))
    unread.value = 0
    ui.ok('已全部标记为已读')
  } catch (e: any) { ui.err(e.message) }
}

async function openInvite() {
  inviteOpen.value = true
  try {
    const r = await socket.call<{ list: any[] }>('club.agentSubs', {})
    agentSubs.value = r.list || []
  } catch {
    agentSubs.value = []
  }
}
function copyInvite() {
  const link = `${location.origin}/login?invite=${p.value?.uid || ''}`
  navigator.clipboard?.writeText(link).then(() => ui.ok('邀请链接已复制')).catch(() => ui.err('复制失败'))
}

function copyUid() {
  navigator.clipboard?.writeText(String(p.value?.uid || '')).then(() => ui.ok('已复制')).catch(() => {})
}

function quick(q: any) {
  if (q.path) return router.push(q.path)
  if (q.label === '客服') return (supportOpen.value = true)
  if (q.label === '我的邀请') return openInvite()
  ui.toast('邀请链接已复制', 'ok')
}

async function doRecharge() {
  try {
    const d = await socket.call<{ balance: number }>('user.recharge', { amount: amount.value })
    auth.applyBalance({ balance: d.balance })
    ui.ok('充值成功')
    rechargeOpen.value = false
    sfx.win()
  } catch (e: any) { ui.err(e.message) }
}

async function doWithdraw() {
  try {
    const d = await socket.call<{ balance: number }>('user.withdraw', { amount: amount.value, payPwd: '' })
    auth.applyBalance({ balance: d.balance })
    ui.ok('提现申请已提交')
    withdrawOpen.value = false
  } catch (e: any) { ui.err(e.message) }
}

async function saveProfile() {
  try {
    await socket.call('user.updateProfile', { nick: form.nick, avatar: form.avatar, signature: form.signature })
    await auth.refresh()
    editOpen.value = false
    ui.ok('已保存')
  } catch (e: any) { ui.err(e.message) }
}

async function savePay() {
  if (!/^\d{6}$/.test(payPwd.value)) return ui.err('请输入 6 位数字')
  try {
    await socket.call('user.setPayPwd', { password: payPwd.value })
    payOpen.value = false
    ui.ok('支付密码已设置')
  } catch (e: any) { ui.err(e.message) }
}

async function buyCards() {
  try {
    const d = await socket.call<{ cards: number }>('user.buyCards', { count: 10 })
    auth.applyBalance({ cards: d.cards })
    ui.ok('已购买 10 张房卡')
  } catch (e: any) { ui.err(e.message) }
}

function logout() {
  auth.logout()
  router.replace('/login')
}

onMounted(async () => {
  if (route?.query?.recharge === '1') rechargeOpen.value = true
  await auth.refresh()
  form.nick = p.value?.nick || ''
  form.avatar = p.value?.avatar || '/avatars/av1.svg'
  form.signature = p.value?.signature || ''
  loadNotices()
})
</script>

<style scoped>
.invite-card {
  display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 14px;
  border-radius: 14px; background: rgba(255, 216, 117, .06);
  border: 1px dashed rgba(255, 216, 117, .34);
}
.invite-code { font-size: 24px; font-weight: 900; letter-spacing: 4px; color: #ffe08a; font-variant-numeric: tabular-nums; }

.vip-badge {
  display: inline-flex; align-items: center; gap: 3px; height: 20px; padding: 0 8px;
  border-radius: 10px; font-size: 10px; font-weight: 800; color: #04220f; letter-spacing: .4px;
  background: linear-gradient(180deg, #2ee08a, #12a862);
  box-shadow: 0 2px 8px rgba(20, 200, 120, .35);
}
.club-label { font-size: 11px; font-weight: 800; letter-spacing: 1.6px; color: #e8cf8e; }
.tether { position: relative; font-size: 22px; font-weight: 800; color: #ffd875; line-height: 1; }
.tether::before, .tether::after { content: ""; position: absolute; left: 50%; transform: translateX(-50%); width: 12px; height: 1.5px; background: currentColor; opacity: .9; }
.tether::before { top: 40%; }
.tether::after { top: 58%; }

.profile { padding-top: 14px; }
.p-head { display: flex; align-items: center; gap: 8px; }
.p-head .round { width: 42px; height: 42px; }
.head-btn { height: 34px; padding: 0 12px; border-radius: 10px; background: rgba(255,255,255,.06); border: 1px solid var(--line); font-size: 13px; display: flex; align-items: center; gap: 4px; }
.round { position: relative; width: 42px; height: 42px; border-radius: 50%; background: rgba(16, 44, 30, .8); border: 1px solid var(--line); display: flex; align-items: center; justify-content: center; }
.user-card { display: flex; align-items: center; gap: 11px; margin-top: 10px; padding: 9px 12px; }
.copy { color: #7ff0b0; margin-left: 4px; }
.sig { display: block; margin-top: 3px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.asset-card { margin-top: 10px; display: flex; flex-direction: column; gap: 6px; background: linear-gradient(150deg, rgba(26,60,42,.9), rgba(10,28,19,.94)); }
.cards-pill { display: inline-flex; align-items: center; gap: 4px; height: 24px; padding: 0 10px; border-radius: 8px; font-size: 11px; color: #ffe3a3; background: rgba(239,207,142,.1); border: 1px solid rgba(239,207,142,.22); }
.bal-big { font-size: 29px; font-weight: 800; letter-spacing: .5px; color: #fff; text-shadow: 0 2px 16px rgba(255,255,255,.14); }
.dot-green { display: inline-block; width: 6px; height: 6px; border-radius: 50%; background: #2ee07a; box-shadow: 0 0 8px rgba(46,224,122,.8); margin-right: 4px; }
.asset-actions { display: flex; gap: 8px; margin-top: 6px; }
.asset-actions .btn { padding: 0 8px; font-size: 13px; white-space: nowrap; }
.quick-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 7px; margin-top: 10px; }
.quick { position: relative; height: 80px; border-radius: 14px; background: rgba(14, 40, 27, .7); border: 1px solid var(--line); display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 4px; font-size: 12px; color: #d9ecdf; }
.qdot { position: absolute; top: 10px; right: 14px; width: 6px; height: 6px; border-radius: 50%; background: #2ee07a; }
.stats { display: grid; grid-template-columns: 1fr 1.15fr 1fr; grid-template-rows: auto auto; gap: 14px 8px; align-items: center; padding: 16px 12px; }
.stats .radial { grid-column: 2; grid-row: 1 / span 2; }
.stats .stat-cell:nth-child(1) { grid-column: 1; grid-row: 1; }
.stats .stat-cell:nth-child(3) { grid-column: 3; grid-row: 1; }
.stats .stat-cell:nth-child(4) { grid-column: 1; grid-row: 2; }
.stats .stat-cell:nth-child(5) { grid-column: 3; grid-row: 2; }
.sec-h { font-size: 14px; font-weight: 800; }
.link-btn {
  display: inline-flex; align-items: center; gap: 4px; height: 30px; padding: 0 12px; border-radius: 9px;
  font-size: 12px; font-weight: 700; color: #cfe9dc;
  background: rgba(255,255,255,.05); border: 1px solid rgba(122,196,156,.28);
}
.stat-cell { display: flex; flex-direction: column; align-items: center; }
.stat-cell b { font-size: 20px; font-weight: 800; }
.stat-cell span { font-size: 11px; color: var(--muted); }
.radial { position: relative; display: flex; align-items: center; justify-content: center; }
.radial svg { width: 110px; }
.radial-txt { position: absolute; display: flex; flex-direction: column; align-items: center; }
.radial-txt b { font-size: 15px; color: #7ff0b0; }
.radial-txt span { font-size: 11px; color: var(--muted); }
.foot { grid-column: 1 / 4; display: flex; justify-content: space-around; padding-top: 10px; border-top: 1px solid var(--line); }
.foot div { display: flex; flex-direction: column; align-items: center; }
.foot b { font-size: 16px; }
.foot span { font-size: 11px; color: var(--muted); }
.close-x { width: 32px; height: 32px; border-radius: 10px; background: rgba(255,255,255,.06); display: flex; align-items: center; justify-content: center; color: var(--muted); }
.grid3 { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.amt { padding: 10px 6px; border-radius: 12px; background: rgba(8,26,17,.8); border: 1px solid var(--line); display: flex; flex-direction: column; align-items: center; gap: 2px; font-weight: 800; }
.amt.on { border-color: rgba(46,224,122,.6); background: rgba(16,60,38,.8); }
.avatar-row { display: flex; gap: 12px; align-items: center; }
.av-pick { width: 36px; height: 36px; border-radius: 50%; overflow: hidden; border: 2px solid transparent; padding: 0; }
.av-pick.on { border-color: #2ee07a; }
.av-pick img { width: 100%; height: 100%; object-fit: cover; }
.dot-badge {
  position: absolute; top: -2px; right: -3px; min-width: 15px; height: 15px; padding: 0 4px;
  border-radius: 8px; background: linear-gradient(180deg, #ff6b6b, #e5312f); color: #fff;
  font-size: 9px; font-weight: 800; font-style: normal; line-height: 15px; text-align: center;
  box-shadow: 0 2px 6px rgba(229, 49, 47, .5);
}
.notice-sheet { max-height: 74vh; overflow-y: auto; }
.notice-empty { display: flex; flex-direction: column; align-items: center; gap: 10px; padding: 34px 0; }
.notice-item {
  display: flex; gap: 10px; padding: 11px 0; border-bottom: 1px solid rgba(255, 255, 255, .05);
}
.notice-item:last-child { border-bottom: 0; }
.notice-item.unread .notice-title { color: #ffe6a8; }
.notice-ico {
  flex: none; width: 30px; height: 30px; border-radius: 10px; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(180deg, #8ff0bc, #33b878);
}
.notice-ico.wallet { background: linear-gradient(180deg, #ffe08a, #e0b455); }
.notice-ico.club { background: linear-gradient(180deg, #9be3ff, #45a9d8); }
.notice-ico.agent { background: linear-gradient(180deg, #d3b8ff, #8a6cd8); }
.notice-title { font-size: 13px; font-weight: 800; }
.notice-body { display: block; margin-top: 3px; font-size: 11.5px; line-height: 1.5; color: #9fb8ac; }
.mini-btn {
  height: 26px; padding: 0 10px; border-radius: 8px; font-size: 11px; font-weight: 700;
  color: #cfeedc; background: rgba(255, 255, 255, .06); border: 1px solid var(--line);
}
.mini-btn:disabled { opacity: .4; }

/* ---- 官方资料卡对齐 ---- */
.nick { font-size: 16px; font-weight: 800; letter-spacing: .2px; color: #fff; }
.uid-line { font-size: 12px; color: #9fb8ac; }
.sig { display: block; margin-top: 3px; font-size: 11px; color: #7f9a8d; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.asset-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 10px; }
.cash-stack { position: relative; width: 42px; height: 30px; flex: none; }
.cash-sheet {
  position: absolute; left: 0; top: 0; width: 38px; height: 26px; border-radius: 5px;
  background: linear-gradient(160deg, #eef4f8, #b9c8d2);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, .55), 0 3px 8px rgba(0, 0, 0, .35);
}
.cash-sheet::after {
  content: ""; position: absolute; inset: 5px 6px auto; height: 1px; background: rgba(90, 120, 140, .45);
  box-shadow: 0 5px 0 rgba(90, 120, 140, .35), 0 10px 0 rgba(90, 120, 140, .25);
}
.cash-sheet.b { left: 5px; top: 4px; opacity: .5; }
.asset-side { display: flex; flex-direction: column; align-items: flex-end; gap: 4px; }
.asset-label { display: flex; align-items: center; gap: 5px; margin-top: 0; font-size: 12px; color: #b9d2c4; }
.tether.sm { font-size: 11px; }
.tether.sm::before, .tether.sm::after { width: 7px; height: 1px; }
.asset-mid { display: flex; align-items: flex-end; justify-content: space-between; gap: 10px; }
.settle { font-size: 11px; color: #8fb3a1; white-space: nowrap; }
.settle .dot-green { margin-right: 5px; }

.quick-ico {
  width: 40px; height: 40px; border-radius: 12px; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(180deg, #7ef0b4, #2bbd7a);
  box-shadow: 0 4px 12px rgba(35, 190, 125, .28), inset 0 1px 0 rgba(255, 255, 255, .45);
}
.quick span:last-child { font-size: 12px; }
.row-ico {
  width: 34px; height: 34px; border-radius: 10px; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(180deg, rgba(126, 240, 180, .22), rgba(35, 160, 106, .16));
  border: 1px solid rgba(109, 243, 176, .3);
}
.row-sub { font-size: 10px; color: #7f9a8d; margin-top: 1px; }
.row-ico.danger { background: rgba(255, 120, 120, .12); border-color: rgba(255, 120, 120, .32); }
.warn-dot {
  width: 15px; height: 15px; border-radius: 50%; font-size: 10px; font-weight: 900; font-style: normal;
  line-height: 15px; text-align: center; color: #2a0a0a; background: linear-gradient(180deg, #ff9f9f, #e05a5a);
}

/* 桌面设置弹层 */
.pref-title { font-size: 11px; font-weight: 800; letter-spacing: .8px; color: #8fb3a1; margin-bottom: 8px; }
.cloth-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.cloth-item { border-radius: 10px; overflow: hidden; border: 1.5px solid transparent; background: rgba(255,255,255,.04); }
.cloth-item img { width: 100%; height: 54px; object-fit: cover; display: block; }
.cloth-item span { display: block; font-size: 10px; padding: 4px 0 5px; color: #bcd5c8; }
.cloth-item.on { border-color: rgba(126,240,180,.85); box-shadow: 0 4px 14px rgba(35,190,125,.25); }
.back-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px; }
.back-item { display: flex; align-items: center; gap: 8px; padding: 6px 8px; border-radius: 10px; background: rgba(255,255,255,.04); border: 1.5px solid transparent; }
.back-item img { width: 26px; height: 36px; border-radius: 4px; object-fit: cover; }
.back-item span { font-size: 11px; color: #bcd5c8; }
.back-item.on { border-color: rgba(126,240,180,.85); }
.pref-row { display: flex; align-items: center; justify-content: space-between; padding: 10px 0; border-top: 1px solid rgba(255,255,255,.06); margin-top: 6px; }
.pref-row b { font-size: 13px; }
.switch { position: relative; width: 44px; height: 26px; border-radius: 13px; background: rgba(255,255,255,.12); }
.switch i { position: absolute; left: 2px; top: 2px; width: 22px; height: 22px; border-radius: 50%; background: #cfe6d9; transition: transform .18s ease; }
.switch.on { background: linear-gradient(180deg, #2ee08a, #12a862); }
.switch.on i { transform: translateX(18px); background: #fff; }

/* 官方：图标为细线描，无实心底 */
.quick { gap: 6px; }
.quick span:last-child { font-size: 12.5px; color: #d9ecdf; }
.copy-ico { vertical-align: -2px; margin-left: 3px; }

/* 个人中心：与官方同基准（整页固定、16px 基准、normal 行高） */
.profile-page { height: 100dvh; overflow: hidden; font-size: 16px; line-height: normal; }
.profile-page .user-nickname-text { line-height: 18px; }
.profile-page .card-brand-tag { line-height: 15px; }
.profile-page .balance-amount { font-size: 30px; line-height: 34px; }
</style>
