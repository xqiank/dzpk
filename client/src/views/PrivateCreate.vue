<template>
  <div class="create-page">
  <div class="create-room-page">
    <header class="create-page-header">
      <button class="back-nav-btn" @click="goBack">
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 19l-7-7 7-7"></path></svg>
        <span>返回</span>
      </button>
      <div class="header-center-info">
        <h1 class="header-page-title">创建牌桌</h1>
        <div class="pull-down-indicator"><span class="indicator-bar"></span></div>
      </div>
      <div class="header-card-badge" @click="buyCards">
        <span class="card-icon">
          <svg viewBox="0 0 24 24" width="15" height="15" fill="currentColor"><path d="M22 10V6c0-1.1-.9-2-2-2H4c-1.1 0-2 .9-2 2v4c1.1 0 2 .9 2 2s-.9 2-2 2v4c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2v-4c-1.1 0-2-.9-2-2s.9-2 2-2zm-2-1.46c-1.19.69-2 1.99-2 3.46s.81 2.77 2 3.46V18H4v-1.46c1.19-.69 2-1.99 2-3.46s-.81-2.77-2-3.46V6h16v2.54z"></path></svg>
        </span>
        <span class="card-count">{{ auth.profile?.cards ?? 0 }} 张房卡</span>
      </div>
    </header>

    <main class="create-page-scroll">
      <form class="create-form-content" @submit.prevent>
        <!-- 1 玩法 -->
        <section class="form-card game-type-section">
          <div class="section-title"><span class="section-badge">1</span><strong>选择玩法</strong></div>
          <div class="game-type-tabs">
            <button
              v-for="g in gameTypes"
              :key="g.key"
              type="button"
              class="game-tab-btn"
              :class="{ active: form.gameType === g.key }"
              @click="form.gameType = g.key"
            >
              <span class="tab-title">{{ g.title }}</span>
              <small class="tab-sub">{{ g.sub }}</small>
            </button>
          </div>
        </section>

        <!-- 2 基础信息 -->
        <section class="form-card">
          <div class="section-title"><span class="section-badge">2</span><strong>基础信息</strong></div>

          <div class="field-item">
            <label class="field-label">房间名称</label>
            <input v-model="form.name" class="u5-input" type="text" maxlength="20" placeholder="U5德州 VIP" />
          </div>

          <div class="field-item club-setting-item">
            <div class="toggle-row">
              <div class="toggle-text">
                <span class="field-label" style="margin: 0; display: inline-flex; align-items: center; gap: 6px">
                  <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor" style="color: #ffd700; flex-shrink: 0"><path d="M12 2L4 5.5v5.5c0 5 3.5 9.5 8 11 4.5-1.5 8-6 8-11V5.5L12 2z"></path></svg>
                  俱乐部专属牌桌
                </span>
                <span class="toggle-sub">开启后牌桌属于指定俱乐部，仅俱乐部成员可入座，带有金色专属勋章</span>
              </div>
              <label class="switch-control">
                <input v-model="form.clubOnly" type="checkbox" :disabled="!myClubId" />
                <span class="switch-slider"></span>
              </label>
            </div>
          </div>

          <div class="field-item">
            <label class="field-label">人数上限</label>
            <div class="pill-select-group">
              <button
                v-for="s in seatOptions"
                :key="s.value"
                type="button"
                class="pill-btn"
                :class="{ active: form.seatCount === s.value }"
                @click="form.seatCount = s.value"
              >{{ s.label }} <span v-if="s.badge" class="pill-badge">{{ s.badge }}</span></button>
            </div>
          </div>

          <div class="field-item">
            <label class="field-label">牌局时长</label>
            <div class="pill-select-group">
              <button
                v-for="d in durationOptions"
                :key="d.value"
                type="button"
                class="pill-btn"
                :class="{ active: form.duration === d.value }"
                @click="form.duration = d.value"
              >{{ d.label }}</button>
            </div>
            <div class="duration-cost-tip">
              <svg viewBox="0 0 24 24" width="13" height="13" fill="#ffd700" style="vertical-align: -2px; margin-right: 4px"><path d="M12 2C8.13 2 5 5.13 5 9c0 2.38 1.19 4.47 3 5.74V17c0 .55.45 1 1 1h6c.55 0 1-.45 1-1v-2.26c1.81-1.27 3-3.36 3-5.74 0-3.87-3.13-7-7-7zm-3 18c0 .55.45 1 1 1h4c.55 0 1-.45 1-1v-1H9v1z"></path></svg>
              到期后牌局自动封盘解散并生成详细战绩。预计消耗 <b>{{ cardCost }}</b> 张房卡。
            </div>
          </div>

          <div class="field-item">
            <label class="field-label">房间密码（选填）</label>
            <input v-model="form.password" class="u5-input" type="text" maxlength="8" placeholder="留空则公开无密码，任何人可加入" />
          </div>
        </section>

        <!-- 3 盲注与筹码买入 -->
        <section class="form-card">
          <div class="section-title"><span class="section-badge">3</span><strong>盲注与筹码买入</strong></div>

          <div class="field-item">
            <label class="field-label">大小盲级别</label>
            <div class="blind-pills-grid">
              <button
                v-for="b in blindOptions"
                :key="b.key"
                type="button"
                class="blind-card-btn"
                :class="{ active: form.blindKey === b.key }"
                @click="pickBlind(b.key)"
              >
                <strong class="blind-val">{{ b.label }}</strong>
                <span class="blind-tag">{{ b.tag }}</span>
              </button>
            </div>
            <div v-if="form.blindKey === 'custom'" class="form-row-2" style="margin-top: 10px">
              <div class="field-item">
                <label class="field-label">小盲（USDT）</label>
                <input v-model.number="form.sbUsdt" type="number" min="0.01" step="0.01" class="u5-input" />
              </div>
              <div class="field-item">
                <label class="field-label">大盲（USDT）</label>
                <input v-model.number="form.bbUsdt" type="number" min="0.01" step="0.01" class="u5-input" />
              </div>
            </div>
            <div class="field-sub-note">当前配置：小盲 {{ showSb }} · 大盲 {{ showBb }}</div>
          </div>

          <div class="form-row-2">
            <div class="field-item">
              <label class="field-label">最小买入筹码</label>
              <input v-model.number="form.minBuyin" type="number" min="0.01" step="0.01" class="u5-input" />
            </div>
            <div class="field-item">
              <label class="field-label">最大买入筹码</label>
              <input v-model.number="form.maxBuyin" type="number" min="0.01" step="0.01" class="u5-input" />
            </div>
          </div>

          <div class="field-item">
            <label class="field-label">前注 Ante（可选）</label>
            <input v-model.number="form.ante" type="number" min="0" step="0.01" class="u5-input" placeholder="0 为无前注" />
          </div>
        </section>

        <!-- 4 特殊模式 -->
        <section class="form-card special-mode-section">
          <div class="section-title"><span class="section-badge hot-badge">热门特色玩法</span></div>
          <div class="mode-toggle-block">
            <label class="mode-header-label">
              <input v-model="form.squid" type="checkbox" class="mode-checkbox" />
              <div class="mode-label-info">
                <div class="mode-name"><span class="mode-icon">🦑</span>鱿鱼模式（烧鱿鱼）<span class="mode-tag-pill">俱乐部高热度</span></div>
                <p class="mode-summary">全员开局被标记，赢池脱险，最后未脱险者向全桌赔付鱿鱼赏金</p>
              </div>
            </label>
          </div>
          <div class="mode-toggle-block">
            <label class="mode-header-label">
              <input v-model="form.bomb" type="checkbox" class="mode-checkbox" />
              <div class="mode-label-info">
                <div class="mode-name"><span class="mode-icon">💣</span>暴击玩法（炸弹底池 Bomb Pot）<span class="mode-tag-pill bomb-tag">翻牌大乱斗</span></div>
                <p class="mode-summary">全员强制投入 5BB 前注，跳过翻牌前直接发翻牌抢大池</p>
              </div>
            </label>
          </div>
        </section>

        <!-- 5 功能开关 -->
        <section class="form-card">
          <div class="section-title"><span class="section-badge">4</span><strong>牌桌功能</strong></div>

          <div class="switches-list">
            <label class="switch-row">
              <div class="switch-info"><strong>自动带入模式</strong><p>筹码不足最小买入时，自动从钱包扣除筹码补充到桌上</p></div>
              <input v-model="form.autoBuyin" type="checkbox" class="toggle-checkbox" />
            </label>
            <label class="switch-row">
              <div class="switch-info"><strong>延迟看牌</strong><p>轮到自己行动时才翻开底牌，沉浸式线下扑克体验</p></div>
              <input v-model="form.delayPeek" type="checkbox" class="toggle-checkbox" />
            </label>
            <label class="switch-row">
              <div class="switch-info"><strong>允许 Straddle（大盲 2 倍）</strong><p>翻牌前枪口位（UTG）可自愿投入 2 倍大盲强制盲注抢占后行动优势</p></div>
              <input v-model="form.straddle" type="checkbox" class="toggle-checkbox" />
            </label>
            <label class="switch-row">
              <div class="switch-info"><strong>全下跑两次（Run it Twice）</strong><p>无后续下注时发两组公共牌分占底池；与保险互斥</p></div>
              <input v-model="form.runTwice" type="checkbox" class="toggle-checkbox" />
            </label>
            <label class="switch-row">
              <div class="switch-info"><strong>全下保险（Insurance）</strong><p>领先方逐街投保，保费从钱包扣取；与发两次互斥</p></div>
              <input v-model="form.insurance" type="checkbox" class="toggle-checkbox" />
            </label>
          </div>
        </section>

        <div class="security-cert-box">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="#6ee7b1" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path></svg>
          <div class="cert-copy">
            <strong>区块链公平洗发牌凭证系统</strong>
            <p>每手牌发牌前自动由服务端生成加密哈希承诺，结算后公开种子，杜绝透视与伙牌。</p>
          </div>
        </div>
        <div class="scroll-bottom-spacer"></div>
      </form>
    </main>

    <footer class="create-page-footer">
      <div class="footer-cost-info">
        <span class="cost-title">创建消耗</span>
        <strong class="cost-value">{{ cardCost }} 张房卡</strong>
      </div>
      <button type="button" class="create-submit-btn" :disabled="busy" @click="submit">
        <span>{{ busy ? '创建中…' : '立即创建牌桌' }}</span>
      </button>
    </footer>
  </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { useUI } from '@/stores/ui'
import { useTable } from '@/stores/table'
import { socket } from '@/lib/ws'
import { sfx } from '@/lib/sfx'
import '@/styles/create-room-official.css'

const router = useRouter()
const auth = useAuth()
const ui = useUI()
const table = useTable()

const gameTypes = [
  { key: 'nlh', title: '标准德州', sub: "NL Hold'em" },
  { key: 'short6', title: '短牌 6+', sub: 'Short Deck' },
  { key: 'niuniu', title: '牛牛', sub: '固定庄 / 抢庄' },
  { key: 'omaha', title: '奥马哈', sub: 'Omaha High · 四张底牌' }
]
const seatOptions = [
  { value: 2, label: '2 人桌', badge: '单挑' },
  { value: 6, label: '6 人桌', badge: '' },
  { value: 8, label: '8 人桌', badge: '推荐' },
  { value: 9, label: '9 人桌', badge: '' }
]
const durationOptions = [
  { value: 60, label: '60 分钟' },
  { value: 120, label: '120 分钟 (标准)' },
  { value: 240, label: '240 分钟' },
  { value: 360, label: '360 分钟' }
]
const blindOptions = [
  { key: 'b1', label: '1 / 2', tag: '新手常用', sb: 1, bb: 2 },
  { key: 'b2', label: '2 / 5', tag: '常规桌', sb: 2, bb: 5 },
  { key: 'b3', label: '5 / 10', tag: '中额桌', sb: 5, bb: 10 },
  { key: 'b4', label: '10 / 20', tag: '高额桌', sb: 10, bb: 20 },
  { key: 'custom', label: '自定义', tag: '自由设定', sb: 1, bb: 2 }
]

const form = ref({
  gameType: 'nlh',
  name: '',
  seatCount: 8,
  duration: 120,
  password: '',
  blindKey: 'b1',
  sbUsdt: 1,
  bbUsdt: 2,
  minBuyin: 100,
  maxBuyin: 400,
  ante: 0,
  clubOnly: false,
  autoBuyin: false,
  delayPeek: false,
  straddle: false,
  runTwice: false,
  insurance: false,
  squid: false,
  bomb: false
})
const busy = ref(false)
const myClubId = ref(0)

const cardCost = computed(() => Math.max(1, Math.min(6, Math.round(form.value.duration / 60))))
const showSb = computed(() => (form.value.blindKey === 'custom' ? form.value.sbUsdt : blindOf().sb))
const showBb = computed(() => (form.value.blindKey === 'custom' ? form.value.bbUsdt : blindOf().bb))

function blindOf() {
  return blindOptions.find(b => b.key === form.value.blindKey) || blindOptions[0]
}
function pickBlind(key: string) {
  form.value.blindKey = key
  const b = blindOf()
  if (key !== 'custom') {
    form.value.sbUsdt = b.sb
    form.value.bbUsdt = b.bb
    const bb = b.bb * 100
    form.value.minBuyin = (bb * 50) / 100
    form.value.maxBuyin = (bb * 200) / 100
  }
  sfx.tap()
}

async function submit() {
  if (busy.value) return
  const sb = Math.round(Number(showSb.value) * 100)
  const bb = Math.round(Number(showBb.value) * 100)
  if (sb <= 0 || bb <= 0) {
    ui.err('盲注不合法')
    return
  }
  const minBuyin = Math.round(Number(form.value.minBuyin) * 100)
  const maxBuyin = Math.round(Number(form.value.maxBuyin) * 100)
  if (minBuyin < bb * 10) {
    ui.err('最小买入至少 10BB')
    return
  }
  if (maxBuyin < minBuyin) {
    ui.err('最大买入不能小于最小买入')
    return
  }
  busy.value = true
  try {
    const d = await socket.call<{ roomNo: number; name: string; cardCost: number }>('lobby.createRoom', {
      name: form.value.name,
      gameType: form.value.gameType,
      seatCount: form.value.seatCount,
      sb,
      bb,
      ante: Math.round(Number(form.value.ante || 0) * 100),
      minBuyin,
      maxBuyin,
      insurance: form.value.insurance,
      duration: form.value.duration,
      password: form.value.password,
      squid: form.value.squid,
      bomb: form.value.bomb,
      autoBuyin: form.value.autoBuyin,
      delayPeek: form.value.delayPeek,
      clubId: form.value.clubOnly ? myClubId.value : 0
    })
    ui.ok(`房间已创建，消耗 ${d.cardCost} 张房卡`)
    sfx.chip()
    await auth.refresh()
    await table.join(d.roomNo, minBuyin)
    router.push(`/table/${d.roomNo}`)
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    busy.value = false
  }
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

function goBack() {
  sfx.tap()
  router.push('/private')
}

onMounted(async () => {
  await auth.refresh()
  try {
    const clubs = await socket.call<any[]>('club.list')
    if (Array.isArray(clubs) && clubs.length) myClubId.value = Number(clubs[0].id) || 0
  } catch {
    myClubId.value = 0
  }
})
</script>

<style scoped>
/* 官方那套 iOS 风格开关（该组件样式不在已抓取的 CSS chunk 里，按官方截图复刻） */
.create-page .switch-control { position: relative; display: inline-flex; align-items: center; width: 20px; height: 18px; flex-shrink: 0; }
.create-page .switch-control input { width: 20px; height: 18px; margin: 0; accent-color: #00e575; cursor: pointer; }
.create-page .switch-control input:disabled { opacity: .4; }
.create-page .switch-slider { display: none; }
/* 底栏已固定，滚动区预留空间 */
.create-page .create-page-scroll { padding-bottom: 24px; }
.create-page .toggle-checkbox { width: 20px; height: 20px; accent-color: #00e575; flex-shrink: 0; }
.create-page .pill-btn { height: 32px; min-height: 32px; padding: 6px 12px; }
.create-page .header-page-title { line-height: 24px; }
.create-page .header-card-badge { height: 30px; line-height: normal; }
.create-page .back-nav-btn { padding-left: 15px; padding-right: 15px; }
.create-page .mode-toggle-block { margin-top: 10px; }
.create-page .toggle-text .field-label { line-height: 17px; }
.create-page .toggle-text { gap: 3px; }
.create-page .mode-checkbox { width: 20px; height: 20px; accent-color: #00e575; flex-shrink: 0; }
.create-page .switch-row { cursor: pointer; }
</style>
