<template>
  <div class="table-root">
    <div class="scene" :class="'cloth-' + cloth" :style="sceneStyle" @touchstart="onTwoFinger" @pointerup="onSceneTap">
      <div class="caustics"></div>
      <div class="bubbles" aria-hidden="true">
        <i v-for="b in bubbles" :key="b.id" class="bubble" :style="b.style"></i>
      </div>
      <div class="reef"></div>
      <div class="felt"></div>

      <!-- 顶部栏（官方：菜单 / 纸飞机 / 安全卫士；牌池在右上） -->
      <div class="t-top">
        <div class="t-left">
          <div class="t-rounds">
            <button class="t-round" @click="menu = true"><Icon name="menu" :size="18" /></button>
            <button class="t-round" @click="invite"><Icon name="send" :size="16" /></button>
          </div>
        </div>
        <div class="t-center">
          <span class="t-badge t-badge-btn" @click="safetyOpen = true"><Icon name="shield" :size="13" /> 安全卫士</span>
        </div>
      </div>


      <!-- 座位 -->
      <div v-for="s in positionedSeats" :key="s.seat" class="seat" :data-seat="s.seat"
           :class="{
             acting: st?.toAct === s.seat && st?.state === 'playing',
             'seat-me': s.isMe,
            'hero-docked': s.isMe && showCluster, 'seat-lift': s.isMe && showDock, folded: s.folded, away: s.sittingOut,
             winner: winnersSet.has(s.seat),
             'has-cards': s.isMe && (s.cards?.length || 0) > 0,

             'just-sat': s.isMe && justSat,
             'hero-face-down': s.isMe && (s.cards?.length || 0) > 0 && !faceHero,
             shake: shakeSeat === s.seat
           }"
           :style="{ left: s.pos.left, top: s.pos.top }"
           >
        <div class="avatar-wrapper" :class="{ 'avatar-winner-active': winnersSet.has(s.seat) }">
          <div v-if="winnersSet.has(s.seat)" class="winner-sunburst-rays" aria-hidden="true">
            <div class="winner-ray-layer"></div>
          </div>
          <div v-if="winnersSet.has(s.seat)" class="winner-sparkles" aria-hidden="true">
            <span v-for="n in 14" :key="n" class="sparkle" :class="'sp-' + n"></span>
          </div>
          <span v-if="winnersSet.has(s.seat)" class="showdown-winner-crown" aria-hidden="true">
            <svg viewBox="0 0 24 20" fill="none">
              <path d="M2 16h20l-1.6-9-4.9 4L12 3 8.5 11 3.6 7 2 16z" fill="url(#crownGold)" stroke="#8a5a00" stroke-width="1.2" stroke-linejoin="round" />
              <defs>
                <linearGradient id="crownGold" x1="0" y1="0" x2="1" y2="1">
                  <stop offset="0%" stop-color="#fff6c4" /><stop offset="55%" stop-color="#ffd700" /><stop offset="100%" stop-color="#ff9500" />
                </linearGradient>
              </defs>
            </svg>
          </span>
          <div class="avatar-box" @click.stop="openPlayerCard(s)">
            <Avatar :src="s.avatar" :size="62" :gold="s.isMe" :class="'avatar-img seat-avatar'" />
            <TimerRing :deadline="s.deadline" :active="st?.toAct === s.seat && st?.state === 'playing'"
                       :suppress="s.isMe && ((s.cards?.length || 0) > 0 || (s.cardCount || 0) > 0)" />
          </div>
          <span v-if="isNiu && st?.banker === s.seat" class="dealer-coin-badge bank" :class="s.dealerPos">庄</span>
          <img v-else-if="st?.dealer === s.seat" class="dealer-coin-img" :class="s.dealerPos"
               src="/img/table/dealer.png" alt="D" />
          <div class="opponent-cards-fan" :class="[s.fanSide, { 'fan-muck': s.folded && !s.muckShown }]" v-if="!s.isMe && ((s.cards?.length || 0) > 0 || s.cardCount > 0)">
            <template v-if="s.show && !suspense && (s.cards?.length || 0) > 0">
              <div v-for="(c, i) in s.cards.slice(0, 4)" :key="'sc' + i" class="show-card" :class="'sc-' + i"
                   :style="{ visibility: revealHidden.has(s.seat) ? 'hidden' : 'visible' }">
                <span class="reveal-in" :style="{ animationDelay: (i * 85) + 'ms' }">
                  <PlayingCard :card="c" :face="true" size="xs" :four-color="tablePrefs.fourColor" />
                </span>
              </div>
            </template>
            <template v-else-if="s.folded && s.muckShown && (s.cards?.length || 0) > 0">
              <div v-for="(c, i) in s.cards!.slice(0, 2)" :key="'mk' + i" class="show-card" :class="'sc-' + i">
                <PlayingCard :card="c" :face="true" size="xs" :four-color="tablePrefs.fourColor" />
              </div>
            </template>
            <template v-else>
              <span v-for="(c, i) in displayCards(s)" :key="i" class="mini-card-back"
                    :class="['card-back-' + (i + 1), { 'muck-peekable': s.folded }]"
                    :title="s.folded ? '点击查看盖牌' : ''"
                    :style="{ backgroundImage: `url(${cardBackImg})` }"
                    @click.stop="s.folded && openMuck(s)"></span>
              <span v-if="s.folded && !s.muckShown" class="muck-lock-badge" title="加密牌背 · 点击查看"
                    @click.stop="openMuck(s)"><Icon name="lock" :size="9" /></span>
            </template>
          </div>
          <span v-if="s.insured" class="insured-tag">保险</span>
          <span v-if="!s.robot && s.connected === false" class="seat-offline">断线</span>
        </div>
        <div class="hole-cards-wrapper" v-if="s.isMe && ((s.cards?.length || 0) > 0 || s.cardCount > 0)"
             :class="{ deciding: showCluster, 'hero-winner-cards': heroIsWinner, 'hero-loser-cards': heroIsLoser }"
             @pointerdown="startPeek" @pointerup="endPeek" @pointercancel="endPeek" @pointerleave="endPeek">
          <div v-for="(c, i) in displayCards(s)" :key="i" class="hero-card-slot" :class="['hero-card-' + i, 'slot-' + i, { 'land-pop': heroLand === i }]">
            <PlayingCard :card="c" :face="faceHero" :size="'sm'" :four-color="tablePrefs.fourColor" />
          </div>
          <span v-if="!faceHero && !peeking && showPeekHint" class="peek-hint">按住看牌</span>
        </div>
        <div class="seat-tag" :class="{ 'allin-anim': s.allin && !s.folded, 'win-pill-glow': winnersSet.has(s.seat) }">
          <span class="seat-name">{{ s.nick }}</span>
          <span class="seat-chips" :class="{ 'winner-chips-flash': winnerFlashSeat === s.seat }">{{ chipsx(s.stack) }}</span>
        </div>
        <span v-if="s.bet > 0" class="seat-bet-chips" :class="s.betPos">
          <i class="chip-mini-icon"></i><span class="bet-chip-val">{{ chipsx(s.bet) }}</span>
        </span>
        <transition name="bubble">
          <span v-if="s.lastAction && actVisible(s)" class="action-bubble" :class="'bubble-' + s.lastAction" :style="bubblePos(s)">
            {{ actionLabel(s) }}
          </span>
        </transition>
        <span v-if="s.show && s.handName" class="hand-name-badge">
          {{ s.handName }}<em v-if="isNiu && niuSeatMult(s) > 1" class="mult">×{{ niuSeatMult(s) }}</em>
        </span>
      </div>

      <!-- 空位（官方虚线圆 · 点击入座） -->
      <div v-for="v in vacantSeats" :key="'vac' + v.seat" class="seat seat-vacant vacant-anim"
           :class="{ 'is-preselected': sitOpen && sitSeat === v.seat, 'is-full-house': !vacantSeats.length }"
           :style="{ left: v.pos.left, top: v.pos.top, '--vac-delay': (v.seat * 160) + 'ms' }" @click="openSit(v.seat)">
        <div class="avatar-wrapper">
          <div class="avatar-box vacant-box">
            <span class="vacant-plus">▼</span>
            <span class="vacant-label">入座</span>
          </div>
        </div>
      </div>

      <!-- 牌池（官方位置：右上，牌背 + 剩余张数） -->
      <div class="deck-pool">
        <div class="deck-net-status">
          <i class="net-signal-bars"></i>
          <span class="net-ms-text">{{ netState.latency }}ms</span>
        </div>
        <div class="deck-shuffle">
          <img :src="cardBackImg" alt="" />
          <img :src="cardBackImg" alt="" />
          <img :src="cardBackImg" alt="" />
          <img :src="cardBackImg" alt="" />
        </div>
        <div class="deck-pool-info">
          <span class="deck-pool-title">牌池</span>
          <span class="deck-pool-count">剩余 {{ st?.deckLeft ?? 0 }} 张</span>
        </div>
      </div>

      <!-- 特效层 -->
      <div class="fx-layer" ref="fxLayer"></div>

      <!-- 底池（官方位置：毛毡区 21%） -->
      <div class="center-pot-area" :class="{ bump: potBump }">
        <div class="pot-felt-tray">
          <div class="tray-watermark">T</div>
          <div class="pot-usdt-chip-stack has-chips">
            <i class="chip-stack-disc disc-layer-1"></i>
            <div class="chip-stack-top-wrapper">
              <img class="chip-stack-top" src="/img/table/usdt_chip_3d.svg" alt="" />
            </div>
          </div>
        </div>
        <div class="pot-pill-wrapper">
          <div class="pot-chip-pill">
            <span class="pot-pill-label">底池</span>
            <span class="pot-chip-amount">{{ chipsx(st?.pot || 0) }}</span>
          </div>
        </div>
      </div>

      <!-- 中央桌面信息 -->
      <div class="center-room-title">
        <span class="room-title-text">{{ st?.name || '牌桌' }}</span>
        <span class="blinds-badge">{{ money(st?.sb || 0, { bare: true }) }}/{{ money(st?.bb || 0, { bare: true }) }}</span>
      </div>

      <!-- 公共牌区（官方：毛毡区 47%，牌 45×64 / 间距 3） -->
      <div class="community-board-area" :class="{ 'has-showdown-winner': !!st?.winners?.length }">
        <div class="board-slots" :class="{ flipping: boardFlip }" v-if="!isNiu">
          <template v-for="i in 5" :key="i">
            <div class="card-slot-placeholder" :class="{ 'has-card': !!st?.board?.[i - 1] }">
              <div v-if="st?.board?.[i - 1]" class="card-item card-dealt"
                   :class="{ 'land-pop': boardLand === i - 1, 'win-board': winBoard.has(st!.board[i - 1]), 'win-card-golden': winBoard.has(st!.board[i - 1]) && !!st?.winners?.length, 'is-winning-board-card': winBoard.has(st!.board[i - 1]) && !!st?.winners?.length }"
                   :style="{ animationDelay: (i * 40) + 'ms' }">
                <PlayingCard :card="st!.board[i - 1]" size="md" :four-color="tablePrefs.fourColor" />
              </div>
            </div>
          </template>
        </div>
        <div class="table-meta-watermark">
          <div class="meta-line blinds-line">
            现金场-{{ gameLabel }} <span class="meta-dot">·</span>
            {{ money(st?.sb || 0, { bare: true }) }}/{{ money(st?.bb || 0, { bare: true }) }}
            <span class="meta-dot">·</span> ID:{{ st?.roomNo }}
          </div>
          <div class="meta-line niu-meta-line" v-if="isNiu">
            牛牛牌桌 · {{ st?.seatCount }}人桌 · {{ niuBankerMode === 'fixed' ? '固定庄' : '抢庄' }}
            <span class="meta-dot">·</span> 庄家 {{ bankerName || '待定' }}
            <span class="meta-dot">·</span> 底注 {{ chips(st?.bb || 0) }}
          </div>
          <div class="meta-line" v-else>按行动时间看牌 · {{ st?.insuranceEnabled ? '保险已开启' : '保险未开启' }}</div>
          <div class="meta-line room-countdown-line">房间剩余 {{ remainText }}</div>
        </div>
      </div>

      <!-- 观战互动（默认收成小圆钮，点开才铺开，避免遮挡牌桌） -->
      <template v-if="(st?.mySeat ?? -1) < 0">
        <button v-if="!specOpen" class="spec-fab" @click="specOpen = true" title="互动">
          <img src="/img/props/prop-flower.svg" alt="" />
          <span>互动</span>
        </button>
        <transition name="fade">
          <div v-if="specOpen" class="spec-bar">

            <button v-for="g in specGifts" :key="g.key" class="spec-gift" :title="g.label" @click="sendGift(g.key)">
              <img :src="propSrc(g.key)" alt="" />
            </button>
            <button class="spec-urge" @click="urgeHand">催牌</button>
            <button class="spec-more" @click="specOpen = false"><Icon name="close" :size="12" /></button>
          </div>
        </transition>
      </template>

      <!-- 旁观入座 -->

      <!-- 结算（官方样式：飘在赢家座位上方的胶囊） -->
      <transition name="fade">
        <div v-if="showResult" class="result-banner" :class="resultTier" :style="resultStyle" @click="openLastHand">
          <div class="who">{{ resultTitle }}</div>
          <div class="amt">{{ resultAmount }}</div>
          <div v-if="topHandName" class="result-hand">{{ topHandName }}</div>
        </div>
      </transition>

      <!-- 金币粒子 -->
      <div v-if="coins.length" class="coins">
        <i v-for="c in coins" :key="c.id" class="fx-coin" :style="c.style"></i>
      </div>

      <!-- 摊牌悬念 -->
      <transition name="fade">
        <div v-if="suspense" class="suspense-tag">
          <span class="ss-dot"></span> 摊牌…
        </div>
      </transition>

      <!-- 断线重连遮罩 -->
      <transition name="fade">
        <div v-if="reconnecting" class="reconnect-veil">
          <div class="rc-card">
            <span class="rc-spinner"></span>
            <b>正在重连牌桌</b>
            <span class="tiny muted">第 {{ netState.attempt || 1 }} 次尝试 · 恢复后自动同步你的座位与筹码</span>
          </div>
        </div>
      </transition>

      <!-- 全下冲击波 -->
      <div v-for="w in allinWaves" :key="w.id" class="allin-wave" :style="w.style">
        <span class="ring"></span><span class="ring r2"></span>
        <b class="tag-text">ALL IN</b>
      </div>

      <!-- 大牌型结算影院 -->
      <transition name="cinema">
        <div v-if="bigWin.show" class="win-cinema" :class="bigWin.tier" @click="bigWin.show = false">
          <div class="rays"></div>
          <div class="cinema-inner">
            <span class="cinema-kicker">{{ bigWin.kicker }}</span>
            <b class="cinema-hand">{{ bigWin.name }}</b>
            <span class="cinema-amount">{{ bigWin.amount }}</span>
            <span class="cinema-who">{{ bigWin.who }}</span>
          </div>
          <i v-for="n in 14" :key="n" class="spark" :style="`--i:${n}`"></i>
        </div>
      </transition>

      <!-- 聊天浮层（最新消息气泡） -->
      <div v-if="chatOpen" class="chat-panel">
        <div v-for="(m, i) in (table.messages || []).slice(-6)" :key="i" class="chat-line">
          <b>{{ m.nick }}</b>：<img v-if="m.type === 'emoji'" class="chat-inline"  :src="emojiSrc(m.asset || '')" alt="" />
          <img v-else-if="m.type === 'prop'" class="chat-inline" :src="propSrc(m.asset || '')" alt="" />
          <span v-else>{{ m.content }}</span>
        </div>
      </div>

      <!-- 按住说话浮层（官方 ptt-hud-overlay） -->
      <div class="tc-scope">
      <div v-if="recording" class="ptt-hud-overlay">
        <div class="ptt-hud-card">
          <div class="ptt-mic-wave-box">
            <div class="ptt-pulse-ring ring-1"></div>
            <div class="ptt-pulse-ring ring-2"></div>
            <div class="ptt-mic-icon-center">
              <svg viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="#25e5b6" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 2a3 3 0 0 0-3 3v7a3 3 0 0 0 6 0V5a3 3 0 0 0-3-3z" />
                <path d="M19 10v2a7 7 0 0 1-14 0v-2" />
                <line x1="12" y1="19" x2="12" y2="23" /><line x1="8" y1="23" x2="16" y2="23" />
              </svg>
            </div>
          </div>
          <div class="ptt-soundwave-bars">
            <span class="wave-bar b1"></span><span class="wave-bar b2"></span><span class="wave-bar b3"></span>
            <span class="wave-bar b4"></span><span class="wave-bar b5"></span>
          </div>
          <div class="ptt-time-counter">00:{{ String(Math.min(10, Math.floor(recMs / 1000))).padStart(2, '0') }} / 00:10</div>
          <p class="ptt-instruction">正在对讲广播 · 松开发送</p>
        </div>
      </div>
      </div>

      <!-- 语音消息气泡 -->
      <div v-for="v in voiceBubbles" :key="v.id" class="voice-bubble" :style="chatFxPos(v, 56)">
        <button class="vb-play" :class="{ playing: v.playing, played: v.played }" @click="playVoice(v)">
          <template v-if="v.playing">
            <i class="vb-wave" v-for="n in 3" :key="n" :style="{ animationDelay: (n * 110) + 'ms' }"></i>
          </template>
          <Icon v-else name="play" :size="12" color="#0a2418" />
          <span>{{ (v.dur / 1000).toFixed(1) }}s</span>
        </button>
        <span class="vb-progress"><i :style="{ width: (v.progress || 0) + '%' }"></i></span>
        <span class="vb-nick">{{ v.nick }}<template v-if="v.listeners && v.listeners.length"> · 已听 {{ v.listeners.length }}</template><template v-else-if="v.played && !v.playing"> · 已播放</template></span>
      </div>

      <!-- 底池飞向赢家（官方 flying-pot-chip / flyPotToWinner） -->
      <div class="ts-scope">
        <div class="fx-pot-layer" ref="potFxLayer">
          <div
            v-for="p in potFlyers"
            :key="p.id"
            class="flying-pot-chip"
            :style="{ '--pot-start-x': p.sx + 'px', '--pot-start-y': p.sy + 'px', '--pot-target-x': p.tx + 'px', '--pot-target-y': p.ty + 'px', animationDelay: p.delay + 'ms' }"
          >
            <img src="/assets/table/chip.png" alt="" class="flying-pot-chip-img" />
          </div>
          <div
            v-if="potAmountFly"
            class="flying-pot-amount"
            :style="{ '--pot-start-x': potAmountFly.sx + 'px', '--pot-start-y': potAmountFly.sy + 'px', '--pot-target-x': potAmountFly.tx + 'px', '--pot-target-y': potAmountFly.ty + 'px' }"
          >
            <div class="flying-pot-pill-box">
              <span class="flying-pot-val">+{{ potAmountFly.text }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 发牌飞行（官方 flying-card-item / flyCardKeyframe） -->
      <div class="td-scope">
        <div class="deal-fly-layer">
          <div
            v-for="f in dealFlyers"
            :key="f.id"
            class="flying-card-item"
            :class="{ 'is-board': f.board }"
            :style="{
              '--start-x': f.sx + 'px',
              '--start-y': f.sy + 'px',
              '--flight-dx': f.dx + 'px',
              '--flight-dy': f.dy + 'px',
              '--start-rot': f.board ? '-10deg' : '-22deg',
              '--scale-start': f.board ? '.68' : '.5',
              '--target-rot': '0deg',
              animationDelay: f.delay + 'ms'
            }"
          >
            <img src="/img/table/pokerback_emerald.png" alt="" class="flying-card-img" />
          </div>
        </div>
      </div>

      <!-- 摊牌亮牌扇飞（官方 flying-opponent-fan / flyOpponentFanAnim） -->
      <div class="tr-scope">
        <div class="reveal-fan-layer">
          <div
            v-for="f in revealFans"
            :key="f.id"
            class="flying-opponent-fan"
            :style="{ '--start-x': f.sx + 'px', '--start-y': f.sy + 'px', '--flight-dx': f.dx + 'px', '--flight-dy': f.dy + 'px' }"
          >
            <span v-for="(c, i) in f.cards" :key="i" class="fan-card">
              <PlayingCard :card="c" :face="true" size="xs" :four-color="tablePrefs.fourColor" />
            </span>
          </div>
        </div>
      </div>

      <!-- 社交道具飘浮（官方 flying-props-layer + prop-fly-trajectory） -->
      <div class="tc-scope">
      <div class="flying-props-layer">
        <div
          v-for="f in flyingProps"
          :key="f.id"
          class="flying-prop-item"
          :style="{ left: f.x + 'px', top: f.y + 'px', '--target-x': f.tx + 'px', '--target-y': f.ty + 'px' }"
        >
          <img v-if="f.type === 'prop'" class="prop-real-fly-img" :src="propSrc(f.asset)" alt="道具" />
          <img v-else class="prop-real-fly-img" :src="emojiSrc(f.asset)" alt="表情" />
        </div>
      </div>
      </div>

      <!-- 底部工具栏（官方四圆钮） -->
      <div class="table-toolbar" :class="{ dim: st?.state === 'showdown' }">
        <div class="tb-group">
          <button class="tb-btn tb-voice-toggle" :class="{ off: !ui.voice }" :title="ui.voice ? '关闭语音播报' : '开启语音播报'" @click="toggleVoiceAnnounce">
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M11 5 6 9H3v6h3l5 4V5z" />
              <path v-if="ui.voice" d="M15.5 8.5a5 5 0 0 1 0 7M18.5 6a9 9 0 0 1 0 12" />
              <path v-else d="M16 9l5 6M21 9l-5 6" />
            </svg>
          </button>
          <button class="tb-btn" @click="menu = true"><Icon name="doc" :size="18" /></button>
          <button class="tb-btn" @click="chatSheet = true"><Icon name="chat" :size="18" /></button>
        </div>
        <div class="tb-group">
          <button class="tb-btn tb-voice-btn" :class="{ 'recording': recording }" title="按住说话" aria-label="按住说话"
                  @pointerdown.prevent="startVoice" @pointerup.prevent="stopVoice" @pointercancel="cancelVoice"
                  @pointerleave="cancelVoice"><Icon name="mic" :size="18" /></button>
          <button class="tb-btn" @click="openRebuy"><Icon name="doc" :size="18" /></button>
        </div>
      </div>

      <!-- 官方垂直加注器（overlay + 竖向滑轨） -->
      <transition name="raise-picker">
        <div v-if="raiseOpen" class="vertical-bet-overlay" @click.self="raiseOpen = false">
          <div class="vertical-bet-picker">
            <div class="vertical-bet-value">{{ chipsx(raiseTo) }}</div>
            <div class="vertical-bet-ratio" :style="{ bottom: ratioBottom + 'px' }">
              <b>{{ ratioPct }}%</b>
              <span>{{ chipsx(raiseTo) }}</span>
            </div>
            <div class="vertical-bet-rail" :class="{ glowing: railDrag }" ref="railEl"
                 @pointerdown="railDown" @pointermove="railMove" @pointerup="railUp" @pointercancel="railUp">
              <div class="vertical-bet-dots">
                <i v-for="(d, i) in railDots" :key="i" :class="{ active: d }"></i>
              </div>
              <div class="vertical-bet-handle" :style="{ bottom: handleBottom + 'px' }">{{ raiseStepIcon }}</div>
            </div>
            <div class="vertical-bet-quick-presets">
              <button class="vb-preset" @click="raisePreset(0.5)">1/2</button>
              <button class="vb-preset" @click="raisePreset(0.75)">3/4</button>
              <button class="vb-preset" @click="raisePreset(1)">满池</button>
              <button class="vb-preset" @click="raisePreset(9)">全下</button>
            </div>
            <button class="vertical-bet-confirm" @click="confirmRaise">
              {{ st?.actions?.isRaiseOfBet ? '下注' : '加注' }}<br />{{ chipsx(raiseTo) }}
            </button>
          </div>
        </div>
      </transition>

      <!-- 官方预操作条（轮到自己前的预选） -->
      <div class="pre-action-bar" v-if="canPreAction">
        <button class="pre-action-pill" :class="{ active: preAction === 'check' }" @click="togglePre('check')">
          ✓ 过牌/跟注
        </button>
        <button class="pre-action-pill" :class="{ active: preAction === 'fold' }" @click="togglePre('fold')">
          ✕ 弃牌
        </button>
      </div>

      <!-- 官方圆形动作按钮簇（弃牌 / 跟注·过牌 / 加注） -->
      <div class="action-cluster" v-if="showCluster">
        <button class="wepoker-action-btn btn-fold" :disabled="!st?.actions?.canFold" @click="doAct('fold')">
          <span class="btn-text">弃牌</span>
        </button>
        <button v-if="st?.actions?.canCheck" class="wepoker-action-btn btn-call is-check" @click="doAct('check')">
          <span class="btn-text">过牌</span>
        </button>
        <button v-else class="wepoker-action-btn btn-call" @click="doAct('call')">
          <span class="call-k">跟注</span>
          <span class="call-chip-amount">{{ chipsx(st?.actions?.callAmount || 0) }}</span>
        </button>
        <button class="wepoker-action-btn btn-raise" :disabled="!st?.actions?.canRaise" @click="openRaise">
          <span class="btn-text">{{ st?.actions?.isRaiseOfBet ? '下注' : '加注' }}</span>
          <span class="raise-sub-hint">{{ chipsx(raiseTo) }}</span>
        </button>
      </div>

      <!-- 入座（旁观时，官方右侧虚线圆） -->

      <!-- 底部操作 -->
      <!-- 牛牛抢庄（官方 NiuniuTableHost 抢庄流程） -->
      <div v-if="isNiu && st?.state === 'bidding'" class="niu-bid-layer">
        <div class="niu-bid-panel">
          <div class="niu-bid-head">
            <span class="niu-bid-title">抢庄加倍</span>
            <span class="niu-bid-timer">{{ niuBidLeft }}s</span>
          </div>
          <p class="niu-bid-desc">五张定牛型 · 喊得最高当庄 · 庄家倍数全场结算</p>
          <div class="niu-bid-options">
            <button v-for="opt in NIU_BID_OPTIONS" :key="opt.mult"
                    class="niu-bid-btn" :class="{ active: myBid === opt.mult, pass: opt.mult === 0 }"
                    :disabled="!canBid || niuBidBusy" @click="doNiuBid(opt.mult)">{{ opt.label }}</button>
          </div>
          <div class="niu-bid-status">
            <span v-if="myBid > 0" class="niu-bid-mine">已喊 {{ myBid }} 倍，等其他人…</span>
            <span v-else-if="mySeatObj?.cards?.length" class="niu-bid-mine muted">选一个倍数开始抢庄</span>
            <span v-else class="niu-bid-mine muted">本手未参与，等待定庄</span>
          </div>
          <div class="niu-bid-board">
            <span v-for="s in st.seats.filter(x => (x.bid || 0) > 0)" :key="'nb' + s.seat" class="niu-bid-chip">
              {{ s.seat + 1 }}号位 ×{{ s.bid }}
            </span>
          </div>
        </div>
      </div>

      <!-- 牛牛闲家押注（官方抢庄后押注档位） -->
      <div v-if="isNiu && st?.state === 'betting'" class="niu-bid-layer">
        <div class="niu-bid-panel niu-stake-panel">
          <div class="niu-bid-head">
            <span class="niu-bid-title">{{ isBankerSeat ? '已坐庄 · 等待闲家押注' : '押注加倍' }}</span>
            <span class="niu-bid-timer">{{ niuStakeLeft }}s</span>
          </div>
          <p class="niu-bid-desc">庄家 {{ bankerName || '待定' }} · 抢庄 ×{{ st?.bankerBid || 1 }} · 底注 {{ chips(st?.bb || 0) }}</p>
          <div v-if="!isBankerSeat" class="niu-bid-options niu-stake-options">
            <button v-for="opt in NIU_STAKE_OPTIONS" :key="opt"
                    class="niu-bid-btn" :class="{ active: myStake === opt }"
                    :disabled="!canStake || niuBidBusy" @click="doNiuStake(opt)">
              {{ opt }} 倍
            </button>
          </div>
          <div class="niu-bid-status">
            <span v-if="isBankerSeat" class="niu-bid-mine muted">庄家不押注，等闲家下注后统一开牌</span>
            <span v-else-if="myStake > 0" class="niu-bid-mine">已押 {{ myStake }} 倍底注</span>
            <span v-else class="niu-bid-mine muted">选一个注码开始押注</span>
          </div>
          <div class="niu-bid-board">
            <span v-for="s in st.seats.filter(x => (x.stake || 0) > 0)" :key="'ns' + s.seat" class="niu-bid-chip">
              {{ s.seat + 1 }}号位 ×{{ s.stake }}
            </span>
          </div>
        </div>
      </div>

      <!-- 牛牛开牌：逐张翻开 -->
      <div v-if="isNiu && niuRevealStage > 0" class="niu-reveal-mask" @click="skipNiuReveal">
        <div class="niu-reveal-card">
          <div class="niu-reveal-title">{{ niuRevealStage >= 2 ? '开牌！' : '正在开牌…' }}</div>
          <div class="niu-reveal-row">
            <div v-for="(c, i) in myNiuCards" :key="'nr' + i" class="niu-reveal-slot" :class="{ open: i < niuRevealCount }">
              <PlayingCard :card="c" :face="i < niuRevealCount" :four-color="tablePrefs.fourColor" size="sm" />
            </div>
          </div>
          <div v-if="niuRevealStage >= 2 && myNiuName" class="niu-reveal-name">{{ myNiuName }}<em v-if="niuSeatMult(mySeatObj) > 1"> ×{{ niuSeatMult(mySeatObj) }}</em></div>
        </div>
      </div>

      <div class="actions" v-if="st && showDock && (st.mySeat ?? -1) >= 0">
        <template v-if="raiseOpen">
          <div class="raise-placeholder"></div>
        </template>
        <template v-else-if="isNiu">
          <div class="action-row">
            <button class="act-btn" style="flex:2" @click="menu = true">
              <span v-if="st.state === 'bidding'">抢庄中 · 选择倍数</span>
              <span v-else-if="st.state === 'betting'">{{ isBankerSeat ? '已坐庄 · 等闲家押注' : '押注中 · 选择注码' }}</span>
              <span v-else>{{ st.state === 'playing' ? '发牌中 · 等待开牌' : '本局已开牌' }}</span>
              <span class="sub" v-if="st.state === 'bidding'">喊得最高当庄 · 庄家倍数全场结算</span>
              <span class="sub" v-else-if="st.state === 'betting'">闲家押注 × 牛型 × 庄家倍数</span>
              <span class="sub" v-else>{{ st.state === 'playing' ? '五张定牛型 · 庄闲对赌' : '下一局即将开始' }}</span>
            </button>
            <button class="act-btn" style="flex:1" @click="openRebuy">补码</button>
          </div>
        </template>
        <template v-else>
          <div class="action-row">
            <button class="act-btn" style="flex:2" @click="menu = true">
              <span v-if="(st?.mySeat ?? -1) >= 0">{{ st?.state === 'playing' ? '等待中…' : '本手已结束' }}</span>
              <span v-else class="spec-line">旁观中 · {{ st?.name }}</span>
              <span class="sub" v-if="(st?.mySeat ?? -1) >= 0">点击打开牌桌菜单</span>
              <span class="sub" v-else>点桌面空位，或这里直接入座</span>
            </button>
            <button v-if="(st?.mySeat ?? -1) >= 0" class="act-btn" style="flex:1" @click="toggleSitOut">
              {{ mySeatObj?.sittingOut ? '回座' : '暂离' }}
            </button>
            <button v-else class="act-btn act-call" style="flex:1" @click="quickSit">带入入座</button>
            <button class="act-btn" style="flex:1" @click="openRebuy">补码</button>
          </div>
        </template>
      </div>

      <!-- 牌桌聊天（官方 TableChat side-drawer） -->
      <ChatDrawer
        :open="chatSheet"
        :messages="table.messages || []"
        :can-speak="!!st && st.mySeat >= 0"
        :self-uid="auth.profile?.uid"
        :room-no="st?.roomNo"
        :emojis="emojis"
        :prop-items="props"
        :emoji-src="emojiSrc"
        :prop-src="propSrc"
        :seat-of="seatOfUid"
        :hand-log="handLog"
        @close="chatSheet = false"
        @send="onChatSend"
        @play-voice="playVoiceMsg"
      />

      <!-- 道具目标选择 -->
      <div v-if="targetPick" class="mask" @click.self="targetPick = null">
        <div class="sheet">
          <div class="row-between" style="margin-bottom:10px">
            <b class="h2">把{{ propLabel(targetPick) }}投给谁？</b>
            <button class="close-x" @click="targetPick = null"><Icon name="close" :size="18" /></button>
          </div>
          <div class="target-grid">
            <button class="target-item" @click="sendProp(0)">
              <span class="target-ava desk"><Icon name="grid" :size="20" color="#9fd8b8" /></span>
              <span>桌面中央</span>
            </button>
            <button v-for="s in (st?.seats || []).filter(x => x.seat !== st?.mySeat)" :key="s.seat" class="target-item"
                    @click="sendProp(s.uid)">
              <Avatar :src="s.avatar" :size="40" />
              <span>{{ s.nick }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 菜单 -->
      <!-- 牌桌快捷菜单（官方 TableRoomMenu：左侧抽屉，图标/文案/底栏与官方一致） -->
      <div v-if="menu" class="tm-official">
        <div class="drawer-modal-backdrop" @click.self="menu = false">
          <aside class="side-drawer-container menu-drawer">
            <header class="drawer-top-nav">
              <button class="drawer-back-btn" @click="menu = false">‹ 返回</button>
              <div class="drawer-title-group">
                <strong>牌桌快捷菜单</strong>
                <small class="drawer-sub-id">{{ st?.name }} · 盲注 {{ blindsLabel }}</small>
              </div>
              <button class="dialog-close-btn drawer-corner-close" aria-label="关闭" @click="menu = false"><Icon name="close" :size="14" /></button>
            </header>
            <div class="drawer-scroll-content menu-scroll-area">
              <div class="menu-section-group">
                <small class="menu-section-title">牌局交互与复盘</small>
                <button class="menu-nav-item" @click="menu = false; invite()">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/icons/telegram-plane.svg" alt="飞机分享" /></span>
                    <div class="menu-item-texts"><b>邀请好友分享</b><small>飞机机器人分享 / 专属链接</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
                <button class="menu-nav-item" @click="menu = false; openRebuy()">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/icons/chip.png" alt="补码/退码" /></span>
                    <div class="menu-item-texts"><b>筹码管理 (补码 / 退码)</b><small>补充带入筹码或将多余筹码退回钱包</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
                <button class="menu-nav-item" @click="menu = false; showReport()">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/icons/data.png" alt="房间报表" /></span>
                    <div class="menu-item-texts"><b>房间报表统计</b><small>总带入、总手数与战绩排行</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
                <button class="menu-nav-item" @click="menu = false; chatSheet = true">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/mine/message.svg" alt="牌桌聊天" /></span>
                    <div class="menu-item-texts"><b>牌桌聊天互动</b><small>实时频道、快捷语与动态表情</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
                <button v-if="st?.handId" class="menu-nav-item" @click="menu = false; openLastHand()">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><Icon name="replay" :size="20" color="#7dd3fc" /></span>
                    <div class="menu-item-texts"><b>本局手牌回顾</b><small>逐街行动线与公平哈希</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
              </div>

              <div class="menu-section-group">
                <small class="menu-section-title">牌桌外观与偏好</small>
                <button class="menu-nav-item" @click="menu = false; prefsOpen = true">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/mine/settings.svg" alt="牌桌设置" /></span>
                    <div class="menu-item-texts"><b>牌桌偏好设置</b><small>音效、触感震动与语音播报</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
                <button class="menu-nav-item" @click="switchCloth">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/icons/theme.svg" alt="牌桌主题" /></span>
                    <div class="menu-item-texts"><b>切换牌桌呢布</b><small>经典绿呢 / 财神送宝 / 鱿鱼游戏 / 七彩祥云 / 金色大漠 / 深海幽蓝</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
              </div>

              <div class="menu-section-group">
                <small class="menu-section-title">安全存证与规则</small>
                <button class="menu-nav-item" @click="menu = false; hashOpen = true">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/icons/shield.png" alt="公平存证" /></span>
                    <div class="menu-item-texts"><b>公平哈希存证</b><small>SHA-256 牌局加密哈希与随机种子</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
                <button class="menu-nav-item" @click="menu = false; rulesOpen = true">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/icons/rule-card.svg" alt="德州规则" /></span>
                    <div class="menu-item-texts"><b>{{ rulesTitle }}</b><small>{{ rulesSub }}</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
                <button class="menu-nav-item" @click="menu = false; openReportPlayer()">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><Icon name="warning" :size="20" color="#fca5a5" /></span>
                    <div class="menu-item-texts"><b>举报违规玩家</b><small>辱骂 / 消极 / 疑似作弊，运营实时受理</small></div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
              </div>

              <div class="menu-section-group">
                <small class="menu-section-title">席位与筹码状态</small>
                <button v-if="(st?.mySeat ?? -1) >= 0" class="menu-nav-item situp-item" @click="menu = false; toggleSitOut()">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/icons/avatar.png" alt="围观模式" /></span>
                    <div class="menu-item-texts">
                      <b>{{ mySeatObj?.sittingOut ? '回到牌桌继续游戏' : '站起为围观模式' }}</b>
                      <small>离席保留房间观战视角</small>
                    </div>
                  </div>
                  <span class="menu-nav-arrow">›</span>
                </button>
                <div v-else class="menu-nav-item spectator-info-item">
                  <div class="menu-nav-item-left">
                    <span class="menu-item-icon"><img class="menu-real-img" src="/assets/icons/avatar.png" alt="围观模式" /></span>
                    <div class="menu-item-texts"><b>当前为围观状态</b><small>点击牌桌任意空座即可带入入座</small></div>
                  </div>
                  <span class="status-chip-spectator">围观中</span>
                </div>
              </div>
            </div>
            <footer class="drawer-bottom-action menu-bottom-action">
              <div class="drawer-wallet-bar">
                <div class="drawer-wallet-left">
                  <div class="wallet-icon-wrap">
                    <svg class="wallet-svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M20 12V8H6a2 2 0 0 1-2-2c0-1.1.9-2 2-2h12v4" />
                      <path d="M4 6v12c0 1.1.9 2 2 2h14v-4" />
                      <path d="M18 12a2 2 0 0 0-2 2c0 1.1.9 2 2 2h4v-4h-4z" />
                    </svg>
                    <span class="wallet-usdt-badge">₮</span>
                  </div>
                  <div class="drawer-wallet-info">
                    <span class="drawer-wallet-label">钱包可用余额</span>
                    <strong class="drawer-wallet-val">${{ money(auth.profile?.balance) }}</strong>
                  </div>
                </div>
                <button class="drawer-wallet-recharge-btn" title="充值筹码" @click="menu = false; goRecharge()">
                  <span>充值</span><span class="plus-icon">＋</span>
                </button>
              </div>
              <button class="drawer-leave-room-btn leave-btn" @click="menu = false; exitRoom()">
                <img class="menu-real-btn-icon" src="/assets/icons/exit.svg" alt="退出房间" /><span>退出房间</span>
              </button>
              <button v-if="isHost" class="drawer-leave-room-btn dissolve-btn" @click="menu = false; dissolveRoom()">
                <img class="menu-real-btn-icon" src="/assets/icons/dissolve.svg" alt="解散牌桌" /><span>解散牌桌</span>
              </button>
            </footer>
          </aside>
        </div>
      </div>

      <!-- 安全卫士 -->
      <div v-if="safetyOpen" class="mask" @click.self="safetyOpen = false">
        <div class="sheet">
          <div class="row-between" style="margin-bottom:10px">
            <b class="h2">安全卫士</b>
            <button class="close-x" @click="safetyOpen = false"><Icon name="close" :size="18" /></button>
          </div>
          <div class="safety-hero">
            <Icon name="shield" :size="30" color="#6df3b0" />
            <div class="col">
              <b>本桌防护运行中</b>
              <span class="tiny muted">发牌哈希可验证 · 异常行为实时拦截 · 同 IP 同桌监测</span>
            </div>
          </div>
          <div class="list-item" style="padding-left:0" @click="safetyOpen = false; hashOpen = true">
            <span class="sf-ico"><Icon name="lock" :size="17" color="#6df3b0" /></span><div class="col" style="flex:1"><b>公平哈希存证</b><span class="tiny muted">本手牌 SHA-256 与种子</span></div><Icon name="chevron" :size="14" />
          </div>
          <div class="list-item" style="padding-left:0" @click="safetyOpen = false; openReportPlayer()">
            <span class="sf-ico danger"><Icon name="flag" :size="17" color="#ff9f9f" /></span><div class="col" style="flex:1"><b>举报违规玩家</b><span class="tiny muted">辱骂 / 消极 / 疑似作弊</span></div><Icon name="chevron" :size="14" />
          </div>
          <div class="list-item" style="padding-left:0" @click="safetyOpen = false; router.push('/profile')">
            <span class="sf-ico info"><Icon name="help" :size="17" color="#7fd8ff" /></span><div class="col" style="flex:1"><b>客服与申诉</b><span class="tiny muted">7×24 在线，平均 30 秒响应</span></div><Icon name="chevron" :size="14" />
          </div>
          <p class="tiny muted" style="margin-top:10px;line-height:1.6">
            发牌种子由服务端提前承诺哈希、落库存证，牌局结束后可回溯核对。
          </p>
        </div>
      </div>

      <!-- 离桌结算（传送到 body，避免被牌桌缩放容器影响定位） -->
      <Teleport to="body">
      <div v-if="settle" class="mask center settle-mask" @click.self="closeSettle">
        <div class="settle-card">
          <Icon :name="settle.net >= 0 ? 'trophy' : 'coins'" :size="26" :color="settle.net >= 0 ? '#7ff0b0' : '#ff9f9f'" />
          <b class="h2" style="margin-top:8px">{{ settle.net >= 0 ? '这一桌赢了' : '这一桌输了' }}</b>
          <div class="settle-net" :class="settle.net >= 0 ? 'green' : 'red'">{{ settle.net >= 0 ? '+' : '' }}{{ chips(settle.net) }}</div>
          <div class="settle-rows">
            <div class="settle-row"><span>带入合计</span><b>{{ chips(settle.buyin) }}</b></div>
            <div class="settle-row"><span>离桌筹码</span><b>{{ chips(settle.stack) }}</b></div>
            <div class="settle-row"><span>打了几手</span><b>{{ settle.hands }} 手</b></div>
            <div class="settle-row"><span>保险收支</span><b :class="settle.ins >= 0 ? 'green' : 'red'">{{ settle.ins >= 0 ? '+' : '' }}{{ chips(settle.ins) }}</b></div>
            <div class="settle-row"><span>房间</span><b>{{ settle.roomName }}</b></div>
          </div>
          <div v-if="settle.session && settle.session.length" class="settle-hands">
            <span class="tiny muted">近 {{ settle.session.length }} 手</span>
            <div class="sh-list">
              <span v-for="h in settle.session" :key="h.handNo" class="sh-chip" :class="h.delta >= 0 ? 'win' : 'lose'">
                #{{ h.handNo }} {{ h.delta >= 0 ? '+' : '' }}{{ chips(h.delta) }}
              </span>
            </div>
          </div>
          <p class="tiny muted" style="margin:8px 2px 0;text-align:center">筹码已退回钱包，可在「战绩」查看每一手复盘</p>
          <div class="settle-acts">
            <button class="btn btn-ghost" @click="closeSettle">留在大厅</button>
            <button class="btn btn-primary" @click="goRecords">看战绩</button>
          </div>
        </div>
      </div>
      </Teleport>

      <!-- 举报玩家 -->
      <div v-if="reportPlayerOpen" class="mask" @click.self="reportPlayerOpen = false">
        <div class="sheet">
          <div class="row-between" style="margin-bottom:10px">
            <b class="h2">举报玩家</b>
            <button class="close-x" @click="reportPlayerOpen = false"><Icon name="close" :size="18" /></button>
          </div>
          <p class="tiny muted" style="margin:0 0 10px">举报内容会同步到运营后台，恶意举报同样会被记录。</p>
          <div class="rp-list">
            <button v-for="s in reportTargets" :key="s.uid" class="rp-item" :class="{ on: rpTarget === s.uid }" @click="rpTarget = s.uid">
              <Avatar :src="s.avatar" :size="26" />
              <span class="ellipsis">{{ s.nick }}</span>
              <i v-if="s.robot" class="tiny">BOT</i>
            </button>
          </div>
          <div class="rp-reasons">
            <button v-for="r in reportReasons" :key="r.key" class="rp-chip" :class="{ on: rpReason === r.key }" @click="rpReason = r.key">{{ r.label }}</button>
          </div>
          <textarea v-model="rpDetail" class="field rp-detail" maxlength="120" placeholder="补充说明（选填，最多 120 字）"></textarea>
          <button class="btn btn-primary btn-block" style="height:46px;margin-top:12px" :disabled="!rpTarget || rpBusy" @click="submitReportPlayer">
            {{ rpBusy ? '提交中…' : '提交举报' }}
          </button>
        </div>
      </div>

      <!-- 报表 -->
      <!-- 房间报表（官方 TableReportDrawer） -->
      <ReportDrawer :open="reportOpen" :report="table.report" :room="st"
                    @close="reportOpen = false" />

      <!-- 玩家卡（官方 PlayerCardDialog：点头像弹出） -->
      <PlayerCardDialog :open="cardOpen" :seat="cardSeat" :muted="cardMuted" :remark="cardRemark"
                        :props-list="props" @close="cardOpen = false" @copy-id="copyPlayerId"
                        @smiley="smileyTo" @toggle-mute="toggleSeatMute" @mail="mailToSeat"
                        @save-remark="saveSeatRemark" @throw-prop="throwPropToSeat" />

    <!-- 带入筹码入座（官方 JoinDialog） -->
    <div v-if="sitOpen" class="join-official">
      <div class="modal-backdrop" @click.self="sitOpen = false">
        <section class="cocos-dialog-box buyin-dialog">
          <header class="dialog-header">
            <div class="dialog-title-col">
              <strong class="dialog-main-title">带入筹码 · 座位 {{ sitSeat + 1 }}</strong>
              <small class="dialog-subtitle">钱包可用余额 {{ money(auth.profile?.balance) }} USDT</small>
            </div>
            <button class="dialog-close-btn" aria-label="关闭" @click="sitOpen = false"><Icon name="close" :size="15" /></button>
          </header>
          <div class="buyin-seat-picker" role="group" aria-label="选择座位">
            <span class="seat-picker-label">选择座位</span>
            <div class="seat-picker-row">
              <button v-for="v in vacantSeats" :key="'pick' + v.seat" type="button" class="seat-pick-chip"
                      :class="{ active: sitSeat === v.seat }" @click="sitSeat = v.seat; sfx.tap()">
                {{ v.seat + 1 }} 号位
              </button>
              <span v-if="!vacantSeats.length" class="seat-picker-empty">本桌已满，等下一手或换桌</span>
            </div>
          </div>
          <div class="buyin-amount-hero">
            <span class="currency-glyph">₮</span>
            <strong class="buyin-num">{{ money(sitAmount) }}</strong>
            <span v-if="sitBB > 0" class="buyin-bb-hint">({{ sitBB }}BB)</span>
          </div>
          <div v-if="sitInsufficient" class="buyin-insufficient-banner">
            <div class="insufficient-text">
              <svg class="alert-mini-icon" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor"
                   stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" />
              </svg>
              <span>钱包余额不足最低带入要求</span>
            </div>
            <button class="buyin-recharge-btn" @click="sitOpen = false; goRecharge()">去充值 ＋</button>
          </div>
          <div class="buyin-slider-box">
            <div class="slider-track-wrap">
              <input class="cocos-range-slider" type="range" :min="sitMin" :max="sitMax" step="1"
                     :value="sitAmount" :style="{ '--slider-pct': sitPct + '%' }" @input="onSitSlide" />
            </div>
            <div class="slider-limits">
              <small class="limit-tag">最低 {{ money(sitMin) }}</small>
              <small class="limit-tag">最高 {{ money(sitMax) }}</small>
            </div>
          </div>
          <div class="dialog-ratio-row">
            <button class="ratio-pill" :class="{ active: sitRatio === 'min' }" @click="setSitRatio('min')">最低带入</button>
            <button class="ratio-pill" :class="{ active: sitRatio === '100bb' }" @click="setSitRatio('100bb')">100BB标准</button>
            <button class="ratio-pill" :class="{ active: sitRatio === 'half' }" @click="setSitRatio('half')">标准中码</button>
            <button class="ratio-pill" :class="{ active: sitRatio === 'max' }" @click="setSitRatio('max')">满额带入</button>
          </div>
          <div class="dialog-actions-row">
            <button class="dlg-btn-cancel" @click="sitOpen = false">取消</button>
            <button class="dlg-btn-confirm" :disabled="sitBusy || sitInsufficient || sitOver" @click="doSit">
              {{ sitBusy ? '确认入座中…' : sitInsufficient ? '余额不足' : sitOver ? '超出可用余额' : '确认带入坐下' }}
            </button>
          </div>
        </section>
      </div>
    </div>

      <!-- 桌面筹码管理（官方 RebuyDialog：补码 / 退码） -->
      <RebuyDialog
        :open="rebuyOpen"
        :user-balance="Number(auth.profile?.balance || 0)"
        :self-chips="Number(mySeatObj?.stack || 0)"
        :amount="rebuyAmt"
        :min-amount="Number(st?.minBuyin || 0)"
        :max-amount="Number(st?.maxBuyin || 0)"
        :available-max="Number(auth.profile?.balance || 0)"
        :submitting="rebuyBusy"
        :big-blind="Number(st?.bb || 0)"
        @close="rebuyOpen = false"
        @update-amount="rebuyAmt = $event"
        @submit="doRebuy"
        @submit-reduce="doReduceChips"
        @open-recharge="rebuyOpen = false; goRecharge()"
      />

      <!-- 退桌结算（官方 SquidExitDialog） -->
      <SquidExitDialog :open="exitOpen" :preview="exitPreview" :busy="exitBusy" :leaving="exitLeaving"
                       @close="exitOpen = false" @confirm="doExitConfirm" />

      <!-- 退码确认 -->
      <div v-if="cashoutOpen" class="mask center" @click.self="cashoutOpen = false">
        <div class="sheet sheet-center">
          <b class="h2">确认退码</b>
          <p class="sm muted" style="margin:8px 0 0">离桌后剩余筹码 <b class="gold">{{ chipsx(mySeatObj?.stack || 0) }}</b> 立即退回钱包</p>
          <div class="row gap-10" style="margin-top:14px">
            <button class="btn btn-ghost" style="flex:1;height:44px" @click="cashoutOpen = false">再玩一会</button>
            <button class="btn btn-danger" style="flex:1.3;height:44px" @click="doCashout">确认退码</button>
          </div>
        </div>
      </div>

      <!-- 保险（官方 InsuranceDialog：全下摊牌前的限时报价 + 三档投保方案） -->
      <div v-if="ins && ins.seat === (st?.mySeat ?? -1)" class="ins-official">
        <div class="modal-backdrop insurance-backdrop">
          <section class="cocos-dialog-box insurance-dialog" :class="{ 'is-urgent-countdown': insLeft <= 5 && !insBusy }">
            <header class="dialog-header insurance-header">
              <div class="dialog-title-col">
                <strong class="insurance-title">
                  <span class="insurance-shield-icon"><Icon name="shield" :size="20" /></span>德州扑克底池保险
                </strong>
                <small class="dialog-subtitle">{{ (ins.boardCount || 0) >= 4 ? '转牌圈保险 · 仅承保下一张河牌' : '翻牌圈保险 · 仅承保下一张转牌' }}</small>
                <div class="insurance-status-line">
                  <span class="insurance-countdown-pill" :class="{ urgent: insLeft <= 5 }">
                    <Icon name="clock" :size="12" />
                    <span>决断剩余 {{ insLeft }}s</span>
                  </span>
                  <small class="wallet-rem-text">钱包可用: {{ money(auth.profile?.balance) }} USDT</small>
                </div>
              </div>
            </header>
            <div class="insurance-body">
              <div class="insurance-stat-grid">
                <div class="stat-card">
                  <div class="stat-label">待保底池 / 承保上限</div>
                  <div class="stat-val highlight">{{ money(ins.cover) }}</div>
                </div>
                <div class="stat-card">
                  <div class="stat-label">危险 Outs 数</div>
                  <div class="stat-val danger">{{ ins.outs }} 张</div>
                </div>
                <div class="stat-card">
                  <div class="stat-label">参考赔率</div>
                  <div class="stat-val odds-val">1 : {{ (ins.odds || 0).toFixed(1) }}</div>
                </div>
              </div>
              <div class="insurance-outs-section" v-if="ins.outCards && ins.outCards.length">
                <div class="outs-heading">对手可能逆转的反超牌 (共 {{ ins.outCards.length }} 张)</div>
                <div class="outs-cards-container">
                  <div v-for="c in ins.outCards" :key="'out-' + c" class="out-card-mini">
                    <PlayingCard :card="c" :face="true" size="xs" :four-color="tablePrefs.fourColor" />
                  </div>
                </div>
              </div>
              <div class="insurance-plans-section">
                <div class="plans-heading">选择投保方案</div>
                <div class="plans-list">
                  <div v-for="p in insPlanList" :key="p.key" class="plan-card" :class="{ active: insPlan === p.key }"
                       @click="insPlan = p.key">
                    <div class="plan-header">
                      <span class="plan-label">{{ p.label }}</span>
                      <span class="plan-radio-dot"></span>
                    </div>
                    <div class="plan-details">
                      <div class="detail-row">
                        <span>钱包扣取保费</span>
                        <strong class="premium-cost">{{ money(p.premium) }} 筹码</strong>
                      </div>
                      <div class="detail-row">
                        <span>最大理赔额</span>
                        <strong class="payout-amount">+{{ money(p.payout) }} 筹码</strong>
                      </div>
                      <div class="detail-row highlight-net">
                        <span>净理赔金额</span>
                        <strong class="payout-amount net-claim">+{{ money(p.payout) }} 筹码</strong>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="insurance-tips">
                • 保费购买时从钱包扣取；本保单仅保下一张公牌，出牌后无论输赢均即刻生效。<br />
                • 仅赔付承保底池损失；净赔付筹码直接进入桌内。
              </div>
            </div>
            <div class="dialog-actions-row insurance-actions">
              <button class="dlg-btn-cancel insurance-btn-pass" :disabled="insBusy" @click="table.decideInsurance(false)"> 放弃保险 </button>
              <button class="dlg-btn-confirm insurance-btn-buy" :disabled="insBusy || insLeft <= 0" @click="buyInsurance">
                <Icon name="shield" :size="16" />
                <span>{{ insBusy ? '确认中…' : '购买所选保险' }}</span>
              </button>
            </div>
          </section>
        </div>
      </div>

      <!-- 保险等待条（官方 insurance-waiting-banner：别人在评估时旁观） -->
      <div v-if="ins && ins.seat !== (st?.mySeat ?? -1)" class="insurance-waiting-banner">
        <span class="insurance-spin-shield"><Icon name="shield" :size="20" /></span>
        <span class="waiting-text">全下领先玩家正在评估底池保险方案…</span>
      </div>

      <!-- 等待下一局（官方 WaitingNextHandDialog：入座后本局进行中） -->
      <div v-if="waitNextOpen" class="modal-backdrop" @click.self="closeWaitNext">
        <section class="cocos-dialog-box waiting-next-hand-dialog" role="dialog" aria-labelledby="waitingDialogTitle">
          <header class="dialog-header">
            <div class="dialog-title-col">
              <strong id="waitingDialogTitle">等待下一局</strong>
              <small class="dialog-subtitle">当前对局进行中 · 自动排队发牌</small>
            </div>
            <button class="dialog-close-btn" aria-label="关闭" @click="closeWaitNext"><Icon name="close" :size="14" /></button>
          </header>
          <div class="waiting-dialog-body">
            <div class="waiting-pulse-box">
              <div class="waiting-glow-ring ring-1"></div>
              <div class="waiting-glow-ring ring-2"></div>
              <div class="waiting-progress-ring"
                   :class="{ 'is-live': waitPhase === 'play', 'is-deal': waitPhase === 'deal' }"
                   :style="{ '--wp': waitRingPct }"></div>
              <div class="waiting-chip-core">
                <svg viewBox="0 0 24 24" width="36" height="36" fill="none" stroke="#25e5b6"
                     stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10" stroke-dasharray="4 2" />
                  <polyline points="12 6 12 12 16 14" />
                </svg>
              </div>
            </div>
            <div class="waiting-status-pill"><span class="live-dot"></span><span>已成功入座 · 下局自动进入</span></div>
            <div v-if="waitPhase === 'deal'" class="waiting-countdown-row">
              <strong class="waiting-countdown-num">{{ waitNextSec }}</strong>
              <span class="waiting-countdown-unit">秒后自动为您发牌</span>
            </div>
            <p class="waiting-dialog-desc">当前牌桌手牌正在进行中，<strong>系统已为您锁定席位</strong>，将在本局结算完毕后立即自动为您派发手牌。</p>
            <p class="waiting-dialog-subdesc">在此期间您可以安心观战，无需任何额外操作。</p>
          </div>
          <div class="dialog-actions-row">
            <button class="dlg-btn-confirm dlg-btn-teal" @click="closeWaitNext">好的，我知道了</button>
          </div>
        </section>
      </div>

      <!-- 牌桌偏好（音效/语音） -->
      <!-- 牌桌偏好设置（官方 SettingsDrawer） -->
      <SettingsDrawer :open="prefsOpen" :prefs="tablePrefs" @close="prefsOpen = false"
                      @update="patchPrefs" @preview-voice="previewVoicePack"
                      @standup="standUpToWatch" @leave-table="leaveFromSettings" />

      <!-- 查看对手盖牌（官方 MuckRevealDialog） -->
      <MuckRevealDialog :open="muckOpen" :seat="muckSeat" :hand-no="st?.handNo" :cost="muckCost"
                        :card-back="cardBackKey" :vip="muckVipFree" :busy="muckBusy"
                        @close="muckOpen = false" @confirm="confirmMuck" />

      <!-- 主题 / 牌背（官方 ThemeDialog） -->
      <ThemeDialog :open="clothSheet" :cloth="cloth" :card-back="cardBackKey"
                   @close="clothSheet = false" @pick-cloth="setCloth" @pick-card-back="setCardBack" />

      <!-- 哈希 -->
      <!-- 公平存证（官方 FairnessProofDrawer） -->
      <FairnessDrawer :open="hashOpen" :room-no="st?.roomNo" @close="hashOpen = false" />

      <!-- 规则 -->
      <!-- 规则与牌型强弱（官方 RulesDialog） -->
      <RulesDialog :open="rulesOpen" :room="rulesRoom" @close="rulesOpen = false" />

      <!-- 加入提示 -->
      <div v-if="!st" class="loading-layer">
        <div class="spinner"></div>
        <span class="sm muted">正在进入牌桌…</span>
        <button class="btn btn-ghost btn-sm" style="margin-top:10px" @click="leave">返回大厅</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Avatar from '@/components/Avatar.vue'
import Icon from '@/components/Icon.vue'
import PlayingCard from '@/components/PlayingCard.vue'
import TimerRing from '@/components/TimerRing.vue'
import RulesDialog from '@/components/table/RulesDialog.vue'
import ChatDrawer from '@/components/table/ChatDrawer.vue'
import RebuyDialog from '@/components/table/RebuyDialog.vue'
import ReportDrawer from '@/components/table/ReportDrawer.vue'
import PlayerCardDialog from '@/components/table/PlayerCardDialog.vue'
import FairnessDrawer from '@/components/table/FairnessDrawer.vue'
import ThemeDialog from '@/components/table/ThemeDialog.vue'
import MuckRevealDialog from '@/components/table/MuckRevealDialog.vue'
import SettingsDrawer, { type TablePrefs } from '@/components/table/SettingsDrawer.vue'
import SquidExitDialog from '@/components/table/SquidExitDialog.vue'
import { useTable, type SeatState } from '@/stores/table'
import { useAuth } from '@/stores/auth'
import { useUI } from '@/stores/ui'
import { socket, netState } from '@/lib/ws'
import { CLOTHS, CARD_BACKS, clothOf, cardBackOf } from '@/lib/tableSkins'
import { money, chips, remain, actionText, gameNameOf } from '@/lib/format'
import { sfx } from '@/lib/sfx'

const route = useRoute()
const router = useRouter()
const table = useTable()
// 诊断句柄：无头验收脚本可读取当前牌桌快照并驱动行动
;(window as any).__u5table = table
const auth = useAuth()
const ui = useUI()

const st = computed(() => table.state)
const mySeatObj = computed(() => table.mySeat)
const menu = ref(false)
const prefsOpen = ref(false)

/* ---- 牌桌偏好（官方 SettingsDrawer 全量开关，localStorage 持久化） ---- */
const PREF_KEY = 'u5_table_prefs'
function loadPrefs(): TablePrefs {
  const raw = JSON.parse(localStorage.getItem(PREF_KEY) || '{}')
  return {
    sound: localStorage.getItem('u5_sound') !== '0',
    vibration: localStorage.getItem('u5_haptic') !== '0',
    voice: localStorage.getItem('u5_voice') !== '0',
    voicePack: localStorage.getItem('u5_voice_pack') || 'sweet',
    chipsInBB: (localStorage.getItem('u5_chips_unit') || 'chips') === 'bb',
    confirmAllin: raw.confirmAllin !== false,
    autoMuck: raw.autoMuck !== false,
    muteProps: raw.muteProps === true,
    fourColor: raw.fourColor === true,
    deckStyle: raw.deckStyle || 'simple',
    effects: raw.effects || 'auto',
    cloth: localStorage.getItem('u5_cloth') || 'green',
    cardBack: localStorage.getItem('u5_cardback') || 'emerald'
  }
}
const tablePrefs = ref<TablePrefs>(loadPrefs())
function persistPrefs() {
  localStorage.setItem(PREF_KEY, JSON.stringify({
    confirmAllin: tablePrefs.value.confirmAllin,
    autoMuck: tablePrefs.value.autoMuck,
    muteProps: tablePrefs.value.muteProps,
    fourColor: tablePrefs.value.fourColor,
    deckStyle: tablePrefs.value.deckStyle,
    effects: tablePrefs.value.effects
  }))
}
function patchPrefs(patch: Partial<TablePrefs>) {
  tablePrefs.value = { ...tablePrefs.value, ...patch }
  if ('sound' in patch) { ui.sound = !!patch.sound; sfx.setEnabled(!!patch.sound) }
  if ('voice' in patch) { ui.voice = !!patch.voice; sfx.setVoice(!!patch.voice) }
  if ('voicePack' in patch && patch.voicePack) { ui.setVoicePack(patch.voicePack); tablePrefs.value.voicePack = patch.voicePack }
  if ('vibration' in patch) localStorage.setItem('u5_haptic', patch.vibration ? '1' : '0')
  if ('chipsInBB' in patch) {
    chipUnit.value = patch.chipsInBB ? 'bb' : 'chips'
    localStorage.setItem('u5_chips_unit', chipUnit.value)
  }
  if ('cloth' in patch && patch.cloth) {
    cloth.value = patch.cloth
    localStorage.setItem('u5_cloth', patch.cloth)
  }
  if ('cardBack' in patch && patch.cardBack) {
    cardBackKey.value = patch.cardBack
  }
  persistPrefs()
}
function previewVoicePack(key: string) {
  ui.setVoicePack(key as any)
  tablePrefs.value.voicePack = key
  try { sfx.voice('call', { minGap: 0 }) } catch { /* 试听失败忽略 */ }
}
function applyFxLevel() {
  const root = document.querySelector('.table-root') as HTMLElement | null
  if (!root) return
  root.classList.toggle('fx-light', tablePrefs.value.effects === 'light')
  root.classList.toggle('fx-full', tablePrefs.value.effects === 'full')
}
watch(() => tablePrefs.value.effects, applyFxLevel, { immediate: false })
onMounted(() => {
  applyFxLevel()
  window.addEventListener('u5:preview-cinema', previewBigWin)
})
onUnmounted(() => window.removeEventListener('u5:preview-cinema', previewBigWin))

const gestureOn = ref(localStorage.getItem('u5_gesture') !== '0')

function toggleGesture() {
  gestureOn.value = !gestureOn.value
  localStorage.setItem('u5_gesture', gestureOn.value ? '1' : '0')
}
const clothSheet = ref(false)
const packs: Record<string, string> = { sweet: '甜美女声', gentle: '温柔女声', male: '沉稳男声', en: '英文语音', classic: '经典播报' }
function choosePack(key: string) {
  ui.setVoicePack(key)
  sfx.voice('your-turn', { minGap: 0 })
}
function testVoice() {
  sfx.voice('raise', { minGap: 0 })
  sfx.play('chips')
}
const raiseOpen = ref(false)
const raiseTo = ref(0)
const chatOpen = ref(true)
const chatSheet = ref(false)
const chatTab = ref<'quick' | 'emoji' | 'prop'>('quick')
const chatText = ref('')
const targetPick = ref<string | null>(null)
const flyingProps = ref<{ id: number; type: string; asset: string; x: number; y: number; tx: number; ty: number }[]>([])
const potFlyers = ref<{ id: number; sx: number; sy: number; tx: number; ty: number; delay: number }[]>([])
const dealFlyers = ref<{ id: number; sx: number; sy: number; dx: number; dy: number; delay: number; board: boolean }[]>([])
const revealFans = ref<{ id: number; sx: number; sy: number; dx: number; dy: number; cards: string[] }[]>([])
type HandLogEntry = { id: number; kind: 'hand' | 'street' | 'action' | 'result'; street: number; seat: number; nick: string; text: string; amount: number; ts: number }
const handLog = ref<HandLogEntry[]>([])
const STREET_NAMES = ['翻牌前', '翻牌', '转牌', '河牌', '摊牌']
let logSeq = 1
let logActionKey = ''
let logWinKey = ''
function pushLog(e: Omit<HandLogEntry, 'id'>) {
  handLog.value.push({ id: logSeq++, ...e })
  if (handLog.value.length > 140) handLog.value.splice(0, handLog.value.length - 140)
}
const revealHidden = reactive(new Set<number>())
const revealedSeats = new Set<string>()
let revealSeq = 1
const boardLand = ref(-1)
const heroLand = ref(-1)
let dealSeq = 1
const potAmountFly = ref<{ sx: number; sy: number; tx: number; ty: number; text: string } | null>(null)
const potFxLayer = ref<HTMLElement | null>(null)
let potFlySeq = 1
const winnerFlashSeat = ref(-1)
let flashTimer = 0
let flySeq = 1
const shakeSeat = ref<number>(-1)
let chatFxSeq = 1
const quickPhrases = ['赶紧上桌，就差你了！', '这把必须跟！', '抓个同花顺！', '老板大气，再发一车！', '好牌慢慢来', '溜了溜了']
const emojis = ['e1', 'e2', 'e3', 'e4', 'e5', 'e6', 'e7', 'e8']
const props = [
  { key: 'prop-beer', label: '啤酒' }, { key: 'prop-bomb', label: '炸弹' },
  { key: 'prop-chicken', label: '小鸡' }, { key: 'prop-clap', label: '鼓掌' },
  { key: 'prop-flower', label: '鲜花' }, { key: 'prop-gun', label: '手枪' },
  { key: 'prop-kiss', label: '飞吻' }, { key: 'prop-shark', label: '鲨鱼' },
  { key: 'prop-slipper', label: '拖鞋' }, { key: 'prop-smoke', label: '烟雾' },
  { key: 'prop-thumbsup', label: '点赞' }, { key: 'prop-tomato', label: '番茄' },
  { key: 'prop-water', label: '泼水' }
]
const emojiSrc = (a: string) => `/img/emoji/${a}.png`
const propSrc = (a: string) => `/img/props/${a}.svg`

async function sendText() {
  const v = chatText.value.trim()
  if (!v) return
  try { await table.send(v); chatText.value = '' } catch (e: any) { ui.err(e.message) }
}

// ---- 官方聊天抽屉的接线 ----
function seatOfUid(uid: number) {
  const t = st.value
  if (!t) return -1
  const s = (t.seats || []).find((x: any) => Number(x.uid) === Number(uid))
  return s ? Number(s.seat) : -1
}
async function onChatSend(p: { text: string; type?: string; asset?: string; target?: number }) {
  try {
    if (p.type === 'emoji' || p.type === 'prop') {
      await table.send('', p.type, p.asset || '', p.target || 0)
      sfx.chip()
      return
    }
    await table.send(p.text, p.type === 'quick' ? 'quick' : undefined)
    if (p.type === 'quick') sfx.tap()
  } catch (e: any) {
    ui.err(e.message)
  }
}
function playVoiceMsg(m: any) {
  if (!m?.data) { sfx.voice('call', { minGap: 0 }); return }
  playVoice({ data: m.data, dur: m.dur || 0, mid: m.vid, uid: m.uid })
  const me = auth.profile?.uid
  if (m.vid && me && Number(m.uid) !== Number(me)) {
    const list = m.listeners || (m.listeners = [])
    if (!list.includes(Number(me))) list.push(Number(me))
    socket.call('table.voicePlayed', { id: m.vid }).catch(() => {})
  }
}
async function sendQuick(q: string) {
  try { await table.send(q); sfx.tap() } catch (e: any) { ui.err(e.message) }
}
async function sendAsset(type: 'emoji' | 'prop', asset: string, target = 0) {
  try { await table.send('', type, asset, target); sfx.chip() } catch (e: any) { ui.err(e.message) }
}
function propLabel(key: string | null) {
  if (!key) return '道具'
  return (props.find(p => p.key === key)?.label) || '道具'
}
function pickTarget(key: string) {
  targetPick.value = key
}
async function sendProp(uid: number) {
  const key = targetPick.value
  targetPick.value = null
  if (!key) return
  chatSheet.value = false
  await sendAsset('prop', key, uid)
}
function chatFxPos(c: { seat: number }, dx = 0) {
  const t = st.value
  const isMe = !!t && t.mySeat >= 0 && c.seat === t.mySeat
  if (isMe) {
    // docked 状态：贴在自己头像正上方
    const el = document.querySelector('.seat.hero-docked') as HTMLElement | null
    if (el) {
      const r = el.getBoundingClientRect()
      const scene = document.querySelector('.scene') as HTMLElement | null
      const sr = scene?.getBoundingClientRect()
      if (sr) {
        const scale = sr.width / 430 || 1
        return {
          left: ((r.left - sr.left) / scale + r.width / scale + 40 + dx) + 'px',
          top: ((r.top - sr.top) / scale + 6) + 'px'
        }
      }
    }
    return { left: '148px', top: '712px' }
  }
  const ps = positionedSeats.value.find(x => x.seat === c.seat)
  if (!ps) return { left: '50%', top: '40%' }
  return { left: ps.pos.left, top: `calc(${ps.pos.top} - 78px - ${dx}px)` }
}
// 官方 flying-props-layer：从发送者座位起飞 → 目标座位，trajectory 由 CSS keyframes 驱动
function showChatFx(m: any) {
  if (!m || (m.type !== 'emoji' && m.type !== 'prop')) return
  if (tablePrefs.value.muteProps && m.type === 'prop') return
  const s = st.value
  if (!s) return
  const fromSeat = s.seats.find(x => x.uid === m.uid)?.seat
  const toSeat = m.target ? s.seats.find(x => x.uid === m.target)?.seat : undefined
  const from = fromSeat !== undefined ? seatPx(fromSeat) : { x: SCENE_W / 2, y: SCENE_H - 40 }
  const to = toSeat !== undefined ? seatPx(toSeat) : { x: from.x, y: from.y - 80 }
  const id = flySeq++
  flyingProps.value.push({
    id, type: m.type, asset: m.asset || '',
    x: from.x, y: from.y, tx: to.x - from.x, ty: to.y - from.y
  })
  setTimeout(() => { flyingProps.value = flyingProps.value.filter(x => x.id !== id) }, 900)
  if (toSeat !== undefined) {
    setTimeout(() => {
      shakeSeat.value = toSeat
      setTimeout(() => { shakeSeat.value = -1 }, 260)
      burstCoins(to)
    }, 800)
  }
}
const reportOpen = ref(false)
// 观战互动：礼物与催牌
const specGifts = [
  { key: 'prop-flower', label: '鲜花' },
  { key: 'prop-beer', label: '啤酒' },
  { key: 'prop-clap', label: '鼓掌' },
  { key: 'prop-kiss', label: '飞吻' },
  { key: 'prop-thumbsup', label: '点赞' }
]
// 筹码显示单位：实际筹码 / 大盲倍数
const chipUnit = ref<'chips' | 'bb'>((localStorage.getItem('u5_chips_unit') as any) === 'bb' ? 'bb' : 'chips')
function chipsx(cents: number | null | undefined) {
  const v = Number(cents || 0)
  if (chipUnit.value === 'chips') return chips(v)
  const bb = Number(st.value?.bb || 0)
  if (!bb) return chips(v)
  const n = v / bb
  const s2 = Math.abs(n - Math.round(n)) < 0.05 ? String(Math.round(n)) : n.toFixed(1)
  return s2 + 'BB'
}
function toggleChipUnit() {
  chipUnit.value = chipUnit.value === 'chips' ? 'bb' : 'chips'
  localStorage.setItem('u5_chips_unit', chipUnit.value)
  ui.toast(chipUnit.value === 'bb' ? '筹码按大盲显示' : '筹码按实际显示', 'info', 1600)
}

const specOpen = ref(false)
const urgeAt = ref(0)
async function sendGift(key: string) {
  const s = st.value
  if (!s) return
  // 目标是玩家 uid（服务端按 uid 校验是否同桌），优先给正在行动的人
  const bySeat = (seat: number) => s.seats.find(x => x.seat === seat)
  const target = (s.toAct >= 0 ? bySeat(s.toAct) : null) ||
    [...s.seats].sort((a, b) => (b.stack || 0) - (a.stack || 0))[0]
  if (!target?.uid) return
  try {
    await socket.call('table.chat', { type: 'prop', asset: key, target: target.uid })
    sfx.chip()
  } catch (e: any) {
    ui.err(e.message)
  }
}
async function urgeHand() {
  if (Date.now() - urgeAt.value < 8000) {
    ui.toast('催牌冷却中，再等几秒', 'info', 1600)
    return
  }
  urgeAt.value = Date.now()
  try {
    await socket.call('table.chat', { type: 'quick', text: '催一拍，等你出手 👀' })
    sfx.tap()
  } catch (e: any) {
    ui.err(e.message)
  }
}

// 座位断线/回线提示（对手侧）
const connMap = new Map<number, boolean>()
watch(() => st.value?.seats, (seats) => {
  if (!seats) return
  for (const s of seats) {
    if (s.robot || !s.uid) continue
    const prev = connMap.get(s.uid)
    if (prev === true && s.connected === false) ui.toast(`${s.nick} 掉线，本手自动弃牌`, 'info', 2600)
    if (prev === false && s.connected === true) ui.toast(`${s.nick} 已回桌`, 'info', 2000)
    connMap.set(s.uid, !!s.connected)
  }
}, { deep: true })

const suspense = ref(false)
let suspenseTimer = 0
watch(() => st.value?.state, (nv, ov) => {
  if (nv === 'showdown' && ov && ov !== 'showdown') {
    suspense.value = true
    clearTimeout(suspenseTimer)
    suspenseTimer = window.setTimeout(() => { suspense.value = false }, 560)
  }
  if (nv !== 'showdown') suspense.value = false
})

const reconnecting = computed(() => netState.everOpen && netState.status !== 'open')
const safetyOpen = ref(false)
const reportPlayerOpen = ref(false)
const rpTarget = ref(0)
const rpReason = ref('abuse')
const rpDetail = ref('')
const rpBusy = ref(false)
const reportReasons = [
  { key: 'abuse', label: '辱骂骚扰' },
  { key: 'idle', label: '消极游戏' },
  { key: 'cheat', label: '疑似作弊' },
  { key: 'other', label: '其他问题' }
]
const reportTargets = computed(() => {
  const s = st.value
  if (!s) return []
  return s.seats.filter(x => x.uid && !x.isMe)
})
function openReportPlayer() {
  menu.value = false
  rpTarget.value = 0
  rpReason.value = 'abuse'
  rpDetail.value = ''
  reportPlayerOpen.value = true
}
async function submitReportPlayer() {
  if (!rpTarget.value || rpBusy.value) return
  rpBusy.value = true
  try {
    await socket.call('table.reportPlayer', { uid: rpTarget.value, reason: rpReason.value, detail: rpDetail.value })
    ui.ok('举报已提交，运营会尽快处理')
    reportPlayerOpen.value = false
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    rpBusy.value = false
  }
}
const rebuyOpen = ref(false)
const cashoutOpen = ref(false)
const settle = ref<any>(null)
// 本桌逐手战绩（离桌结算用）：记录每手开始筹码与结算后筹码
const sessionHands = ref<{ handNo: number; delta: number }[]>([])
const handStart = new Map<number, number>()
watch(() => [st.value?.version, st.value?.handNo, st.value?.state] as const, () => {
  const s = st.value
  const me = s?.seats.find(x => x.isMe)
  if (!s || !me) return
  const no = Number(s.handNo || 0)
  if (!no) return
  if (s.state === 'playing' && !handStart.has(no)) handStart.set(no, Number(me.stack || 0) + Number(me.bet || 0))
  if (s.state === 'showdown') {
    const start = handStart.get(no)
    if (start === undefined) return
    if (sessionHands.value.some(x => x.handNo === no)) return
    const end = Number(me.stack || 0)
    sessionHands.value = [...sessionHands.value, { handNo: no, delta: end - start }].slice(-30)
  }
})
function closeSettle() {
  settle.value = null
  router.replace('/lobby')
}
function goRecords() {
  const id = st.value?.handId
  settle.value = null
  router.push(id ? '/records/hand/' + id : '/records')
}
const hashOpen = ref(false)
const rulesOpen = ref(false)
const rulesRoom = computed(() => st.value)

/* ---- 查看对手盖牌（官方 MuckRevealDialog） ---- */
const muckOpen = ref(false)
const muckSeatNo = ref(-1)
const muckBusy = ref(false)
const muckSeat = computed(() => positionedSeats.value.find(x => x.seat === muckSeatNo.value) || null)
const muckCost = computed(() => Number(muckSeat.value?.muckFee || (st.value?.bb || 0) * 2 || 0))
const muckVipFree = computed(() => Number(auth.profile?.vip || 0) >= 3)
function openMuck(s: SeatState) {
  if (!s?.folded) return
  muckSeatNo.value = s.seat
  muckOpen.value = true
  sfx.tap()
}
async function confirmMuck() {
  const s = muckSeat.value
  if (!s || muckBusy.value) return
  muckBusy.value = true
  try {
    const r: any = await table.revealMuck(s.seat)
    muckOpen.value = false
    ui.ok(`已查看 ${s.nick} 的盖牌${r?.fee ? ` · 消耗 ${(Number(r.fee) / 100).toFixed(2)} 筹码` : ''}`)
  } catch (e: any) {
    ui.err(e?.message || '查看失败')
  } finally {
    muckBusy.value = false
  }
}

/* ---- 玩家卡（官方 PlayerCardDialog） ---- */
const cardOpen = ref(false)
const cardSeatNo = ref(-1)
const cardSeat = computed(() => positionedSeats.value.find(x => x.seat === cardSeatNo.value) || null)
const remarks = ref<Record<string, string>>(JSON.parse(localStorage.getItem('u5_remarks') || '{}'))
const mutedSeats = ref<Record<string, boolean>>(JSON.parse(localStorage.getItem('u5_muted_seats') || '{}'))
const cardRemark = computed(() => (cardSeat.value ? remarks.value[String(cardSeat.value.uid)] || '' : ''))
const cardMuted = computed(() => (cardSeat.value ? !!mutedSeats.value[String(cardSeat.value.seat)] : false))
function openPlayerCard(s: SeatState) {
  if (!s || s.isMe) return
  cardSeatNo.value = s.seat
  cardOpen.value = true
  sfx.tap()
}
async function copyPlayerId(uid: number) {
  try {
    await navigator.clipboard.writeText(String(uid))
    ui.ok('玩家 ID 已复制')
  } catch {
    ui.err('复制失败')
  }
}
function smileyTo(seatNo: number) {
  const seatObj = positionedSeats.value.find(x => x.seat === seatNo)
  const emoji = emojis[Math.floor(Math.random() * emojis.length)]
  if (seatObj) sendAsset('emoji', emoji, seatObj.uid)
}
function toggleSeatMute(seatNo: number) {
  const key = String(seatNo)
  mutedSeats.value = { ...mutedSeats.value, [key]: !mutedSeats.value[key] }
  localStorage.setItem('u5_muted_seats', JSON.stringify(mutedSeats.value))
  ui.ok(mutedSeats.value[key] ? '已静音该玩家语音' : '已取消静音')
}
function mailToSeat(seatNo: number) {
  const seatObj = positionedSeats.value.find(x => x.seat === seatNo)
  cardOpen.value = false
  chatSheet.value = true
  if (seatObj?.nick) chatText.value = `@${seatObj.nick} `
}
function saveSeatRemark(text: string) {
  if (!cardSeat.value) return
  const uid = String(cardSeat.value.uid)
  remarks.value = { ...remarks.value, [uid]: text }
  localStorage.setItem('u5_remarks', JSON.stringify(remarks.value))
  ui.ok(text ? '备注已保存' : '备注已清空')
}
function throwPropToSeat(key: string) {
  const seatObj = cardSeat.value
  if (!seatObj) return
  sendAsset('prop', key, seatObj.uid)
  cardOpen.value = false
}


const rebuyAmt = ref(0)
const showResult = ref(false)
// ===== 全下冲击波 =====
const allinWaves = ref<{ id: number; style: string }[]>([])
let waveSeq = 1
function spawnAllinWave(seat: number) {
  const pos = seat === st.value?.mySeat ? heroPos() : seatPx(seat)
  const id = waveSeq++
  allinWaves.value = [...allinWaves.value, { id, style: `left:${pos.x}px;top:${pos.y}px` }]
  sfx.allin()
  if (tablePrefs.value.vibration) (navigator as any).vibrate?.([12, 30, 18])
  setTimeout(() => { allinWaves.value = allinWaves.value.filter(w => w.id !== id) }, 1400)
}
// ===== 大牌型结算影院 =====
const BIG_TIERS: Record<string, string> = {
  // NLH / 短牌 / 奥马哈
  皇家同花顺: 'legend', 同花顺: 'legend', 四条: 'epic', 葫芦: 'epic', 同花: 'rare', 顺子: 'rare',
  // 牛牛
  五小牛: 'legend', 炸弹: 'legend', 五花牛: 'epic', 牛牛: 'rare'
}
const bigWin = ref({ show: false, name: '', amount: '', who: '', tier: 'rare', kicker: '' })
const previews = [
  { name: '同花顺', tier: 'legend' }, { name: '四条', tier: 'epic' },
  { name: '葫芦', tier: 'epic' }, { name: '同花', tier: 'rare' },
  { name: '牛牛', tier: 'rare' }, { name: '五小牛', tier: 'legend' }
]
let previewIdx = 0
function previewBigWin() {
  prefsOpen.value = false
  setTimeout(() => previewNow(), 220)
}
function previewNow() {
  const p = previews[previewIdx++ % previews.length]
  bigWin.value = {
    show: true, name: p.name, amount: '+1,280.00', who: '特效预览 · 无实际牌局影响',
    tier: p.tier, kicker: p.tier === 'legend' ? 'LEGENDARY HAND' : p.tier === 'epic' ? 'BIG HAND' : 'NICE HAND'
  }
  sfx.win()
  setTimeout(() => { bigWin.value.show = false }, p.tier === 'legend' ? 3200 : 2400)
}

function triggerBigWin() {
  const s = st.value
  const w = (s?.winners || [])[0]
  if (!s || !w) return
  const tier = BIG_TIERS[w.handName || '']
  if (!tier || tier !== 'legend') return
  bigWin.value = {
    show: true,
    name: w.handName || '',
    amount: '+' + money(w.amount || 0),
    who: w.seat === s.mySeat ? '你赢下了这一手' : `${w.nick} 赢下底池`,
    tier,
    kicker: tier === 'legend' ? 'LEGENDARY HAND' : tier === 'epic' ? 'BIG HAND' : 'NICE HAND'
  }
  setTimeout(() => { bigWin.value.show = false }, tier === 'legend' ? 3200 : 2400)
}
const coins = ref<{ id: number; style: any }[]>([])
const portrait = ref(false)
const winSize = ref({ w: window.innerWidth, h: window.innerHeight })
// 竖屏桌面：430x932 设计稿等比缩放到当前视口
const sceneStyle = computed(() => {
  const scale = Math.min(winSize.value.w / 430, winSize.value.h / 932)
  return { transform: `scale(${scale})`, '--felt-img': `url(${clothUrl.value})` } as any
})
const now = ref(Date.now())
const cloth = ref(localStorage.getItem('u5_cloth') || 'green')
const clothUrl = computed(() => clothOf(cloth.value).url)

const clothName = computed(() => clothOf(cloth.value).name)
const blindsLabel = computed(() => {
  const t = st.value
  if (!t) return ''
  return `${chips(t.sb)}/${chips(t.bb)}`
})

const bubbles = Array.from({ length: 16 }, (_, i) => {
  const size = 4 + (i % 5) * 3
  return {
    id: i,
    style: {
      left: (i * 6.3 + 3) + '%',
      width: size + 'px',
      height: size + 'px',
      animationDuration: (11 + (i % 7) * 2) + 's',
      animationDelay: (-i * 1.7) + 's'
    }
  }
})

const fxLayer = ref<HTMLElement | null>(null)
const potBump = ref(false)
const boardFlip = ref(false)
const SCENE_W = 430
const SCENE_H = 932
let lastDealKey = ''
let lastWinKey = 0

function seatPx(seat: number): { x: number; y: number } {
  const ps = positionedSeats.value.find(x => x.seat === seat)
  if (!ps) return { x: SCENE_W / 2, y: SCENE_H / 2 }
  return { x: (parseFloat(ps.pos.left) / 100) * SCENE_W, y: (parseFloat(ps.pos.top) / 100) * SCENE_H }
}
// 官方牌池：右侧 14 + 宽 52 → 中心 x≈390；牌背中心 y≈62
const deckPos = () => ({ x: 390, y: 62 })
// 官方底池：毛毡区 21% → y≈256（筹码堆中心）
const potPos = () => ({ x: SCENE_W / 2, y: 258 })
// 官方玩家位：毛毡区 84% → y≈791
const heroPos = () => ({ x: SCENE_W / 2, y: 791 })

function fxEl(cls: string, x: number, y: number, w: number, h: number) {
  const layer = fxLayer.value
  if (!layer) return null
  const d = document.createElement('div')
  d.className = cls
  d.style.width = w + 'px'
  d.style.height = h + 'px'
  d.style.left = x - w / 2 + 'px'
  d.style.top = y - h / 2 + 'px'
  layer.appendChild(d)
  // 保险清理：动画被中断时也能回收，避免残留
  setTimeout(() => d.remove(), 2600)
  if (layer.childElementCount > 60) {
    const first = layer.firstElementChild as HTMLElement | null
    if (first) first.remove()
  }
  return d
}

// 官方发牌：牌背从牌池按 flyCardKeyframe 飞到落点，落点触发 cardCrackBurst 弹入
function flyCard(card: { seat: number; index: number }, delay: number) {
  const from = deckPos()
  const s = st.value
  const isBoard = card.seat < 0
  const isMine = !!s && card.seat === s.mySeat
  const to = isBoard
    ? { x: SCENE_W / 2 + ((card.index || 0) - 2) * 48, y: SCENE_H * 0.45 - 10 }
    : isMine
      ? heroPos()
      : { x: seatPx(card.seat).x + 26, y: seatPx(card.seat).y - 16 }
  const id = dealSeq++
  dealFlyers.value.push({
    id, sx: from.x, sy: from.y, dx: to.x - from.x, dy: to.y - from.y, delay, board: isBoard
  })
  const dur = 420
  window.setTimeout(() => {
    dealFlyers.value = dealFlyers.value.filter(f => f.id !== id)
    if (isBoard) {
      boardLand.value = card.index ?? -1
      window.setTimeout(() => { if (boardLand.value === card.index) boardLand.value = -1 }, 220)
    } else if (isMine) {
      heroLand.value = card.index ?? 0
      window.setTimeout(() => { if (heroLand.value === card.index) heroLand.value = -1 }, 220)
    }
  }, delay + dur)
}

function flyChip(from: { x: number; y: number }, to: { x: number; y: number }, delay = 0) {
  const c = fxEl('fx-chip', from.x, from.y, 20, 20)
  if (!c) return
  const dx = to.x - from.x
  const dy = to.y - from.y
  const anim = c.animate(
    [
      { transform: 'translate(0,0) scale(.7)', opacity: 0.2 },
      { transform: `translate(${dx * 0.5}px,${dy * 0.5 - 26}px) scale(1.1)`, opacity: 1, offset: 0.6 },
      { transform: `translate(${dx}px,${dy}px) scale(1)`, opacity: 1 }
    ],
    { duration: 300, delay, easing: 'cubic-bezier(.2,.85,.3,1.1)', fill: 'both' }
  )
  anim.onfinish = () => c.remove()
}

// 官方底池飞币：6 枚错峰 + 1 枚金额浮签，落点触发赢家筹码闪动
function flyPotToWinner(to: { x: number; y: number }, amount: number, seat: number) {
  const from = potPos()
  const base = potFlySeq++
  const list = Array.from({ length: 6 }, (_, i) => ({
    id: base * 10 + i,
    sx: from.x + (i % 3) * 6 - 6,
    sy: from.y + Math.floor(i / 3) * 4 - 2,
    tx: to.x + (i % 3) * 10 - 10,
    ty: to.y + Math.floor(i / 3) * 6 - 3,
    delay: i * 55
  }))
  potFlyers.value = [...potFlyers.value, ...list]
  potAmountFly.value = { sx: from.x, sy: from.y - 26, tx: to.x, ty: to.y - 34, text: money(amount) }
  window.setTimeout(() => {
    potFlyers.value = potFlyers.value.filter(x => !list.some(l => l.id === x.id))
    potAmountFly.value = null
  }, 1600)
  window.setTimeout(() => {
    burstCoins(to)
    winnerFlashSeat.value = seat
    window.clearTimeout(flashTimer)
    flashTimer = window.setTimeout(() => { winnerFlashSeat.value = -1 }, 950)
  }, 1050)
}

function burstCoins(at: { x: number; y: number }) {
  const layer = fxLayer.value
  if (!layer) return
  for (let i = 0; i < 14; i++) {
    const c = fxEl('fx-coin', at.x, at.y, 9, 9)
    if (!c) continue
    const ang = Math.random() * Math.PI * 2
    const dist = 40 + Math.random() * 110
    const anim = c.animate(
      [
        { transform: 'translate(0,0) scale(.5)', opacity: 1 },
        { transform: `translate(${Math.cos(ang) * dist}px,${Math.sin(ang) * dist}px) scale(1.1)`, opacity: 0 }
      ],
      { duration: 700 + Math.random() * 300, delay: i * 22, easing: 'cubic-bezier(.2,.7,.4,1)', fill: 'both' }
    )
    anim.onfinish = () => c.remove()
  }
}

// 摊牌亮牌：对手翻牌瞬间从座位飞出两张开牌扇（官方 flyOpponentFanAnim）
function runRevealFx() {
  const s = st.value
  if (!s) return
  for (const seat of s.seats || []) {
    if (!seat.show || !seat.cards?.length || seat.isMe) continue
    const key = `${s.handNo}:${seat.seat}`
    if (revealedSeats.has(key)) continue
    revealedSeats.add(key)
    const base = seatPx(seat.seat)
    const id = revealSeq++
    revealFans.value.push({
      id, sx: base.x, sy: base.y - 44, dx: 18, dy: 72,
      cards: seat.cards.slice(0, 2)
    })
    revealHidden.add(seat.seat)
    window.setTimeout(() => {
      revealFans.value = revealFans.value.filter(f => f.id !== id)
      revealHidden.delete(seat.seat)
    }, 460)
  }
}

function runFx() {
  const s = st.value
  const p = table.prev
  if (!s) return
  try { runRevealFx() } catch { /* 亮牌动画失败不影响牌局 */ }
  // 发牌
  if (s.newCards?.length) {
    const key = `${s.handNo}:${s.street}:${s.newCards.length}`
    if (key !== lastDealKey) {
      lastDealKey = key
      const cards = [...s.newCards].sort((a, b) => a.index - b.index)
      cards.forEach((c, i) => flyCard(c, i * 55))
    }
  }
  // 下注筹码
  if (p && p.roomNo === s.roomNo && p.handNo === s.handNo) {
    for (const seat of s.seats) {
      const old = p.seats.find(x => x.seat === seat.seat)
      const oldBet = old?.bet ?? 0
      if (seat.bet > oldBet) flyChip(seatPx(seat.seat), potPos(), 0)
    }
  }
  // 底池飞向赢家（官方 flying-pot-chip 曲线 + 金额浮签 + 赢家筹码闪动）
  if (s.winEv && s.winEv.pot !== lastWinKey && s.winEv.seat >= 0) {
    lastWinKey = s.winEv.pot
    const to = s.winEv.seat === s.mySeat ? heroPos() : seatPx(s.winEv.seat)
    flyPotToWinner(to, s.winEv.amount || s.winEv.pot, s.winEv.seat)
  }
}

watch(() => st.value?.pot, (n, o) => {
  if (typeof n === 'number' && typeof o === 'number' && n > o) {
    potBump.value = true
    setTimeout(() => { potBump.value = false }, 380)
  }
})
watch(() => st.value?.street, (n, o) => {
  if (typeof n === 'number' && typeof o === 'number' && n > o) {
    boardFlip.value = true
    setTimeout(() => { boardFlip.value = false }, 520)
  }
})
const allinSeen = new Set<string>()
watch(() => st.value?.version, () => {
  const s = st.value
  if (s) {
    for (const seat of s.seats) {
      const key = `${s.handNo}:${seat.seat}`
      if (seat.allin && !seat.folded && !allinSeen.has(key)) {
        allinSeen.add(key)
        spawnAllinWave(seat.seat)
      }
    }
  }
  try { runFx() } catch { /* FX 失败不影响牌局 */ }
})

// ===== 语音播报：喇叭开关 + 剩余 5 秒催促 =====
function toggleVoiceAnnounce() {
  ui.toggleVoice()
  if (ui.voice) sfx.voice('your-turn', { minGap: 0 })
}

// ===== 牌局战报时间线（官方「牌局战报」页签） =====
watch(() => st.value?.handNo, (n, o) => {
  if (typeof n !== 'number' || n === o) return
  handLog.value = []
  logActionKey = ''
  logWinKey = ''
  pushLog({ kind: 'hand', street: 0, seat: -1, nick: '', text: `第 ${n} 手`, amount: 0, ts: Date.now() })
  pushLog({ kind: 'street', street: 0, seat: -1, nick: '', text: STREET_NAMES[0], amount: 0, ts: Date.now() })
})
watch(() => st.value?.street, (n, o) => {
  if (typeof n !== 'number' || typeof o !== 'number' || n <= o) return
  pushLog({ kind: 'street', street: n, seat: -1, nick: '', text: STREET_NAMES[Math.min(n, 4)], amount: 0, ts: Date.now() })
})
watch(() => st.value?.actionEv, (ev) => {
  const t = st.value
  if (!ev || !t) return
  const key = `${t.handNo}:${ev.ts}:${ev.seat}:${ev.action}:${ev.amount}`
  if (key === logActionKey) return
  logActionKey = key
  const seat = (t.seats || []).find(x => x.seat === ev.seat)
  pushLog({
    kind: 'action', street: ev.street, seat: ev.seat,
    nick: seat?.nick || `${ev.seat + 1}号位`,
    text: actionText(ev.action), amount: ev.amount,
    ts: (() => { const v = Number(ev.ts || 0); if (!v) return Date.now(); return v > 1e12 ? v : v * 1000 })()
  })
})
watch(() => st.value?.winners, (ws) => {
  const t = st.value
  if (!t || !ws?.length) return
  const key = `${t.handNo}:${ws.length}:${ws[0]?.seat}:${ws[0]?.amount}`
  if (key === logWinKey) return
  logWinKey = key
  for (const w of ws) {
    pushLog({
      kind: 'result', street: 4, seat: w.seat || 0,
      nick: w.nick || `${(w.seat || 0) + 1}号位`,
      text: w.handName ? `摊牌 ${w.handName}` : '赢得底池',
      amount: w.amount || 0, ts: Date.now()
    })
  }
})

// ===== 官方语音播报：轮到自己 + 各家行动 =====
const actSig = new Map<number, string>()
let lastTurnKey = ''
watch(
  () => [st.value?.version, st.value?.toAct, st.value?.handNo, st.value?.street] as const,
  () => {
    const t = st.value
    if (!t) return
    // 轮到自己
    const turnKey = `${t.handNo}:${t.street}:${t.toAct}`
    if (t.state === 'playing' && t.mySeat >= 0 && t.toAct === t.mySeat && turnKey !== lastTurnKey) {
      lastTurnKey = turnKey
      sfx.voice('your-turn')
      sfx.play('turn')
    }
    // 各家行动播报
    for (const seat of t.seats || []) {
      const act = seat.lastAction || ''
      const sig = `${t.handNo}:${act}:${seat.lastAmount}`
      const prev = actSig.get(seat.seat)
      if (prev === sig) continue
      actSig.set(seat.seat, sig)
      if (!act || act === 'sb' || act === 'bb' || act === 'post') continue
      if (seat.isMe) continue
      const map: Record<string, any> = { call: 'call', check: 'check', fold: 'fold', raise: 'raise', bet: 'raise', allin: 'all-in' }
      const v = map[act]
      if (v) sfx.voice(v, { minGap: 120 })
      if (act === 'call' || act === 'raise' || act === 'bet' || act === 'allin') sfx.play('chips')
    }
  }
)
let offVoice: (() => void) | null = null
let offVoicePlayed: (() => void) | null = null
onMounted(() => {
  offVoice = socket.on('table.voice', showVoice)
  offVoicePlayed = socket.on('table.voicePlayed', (m: any) => {
    const b = voiceBubbles.value.find(x => x.mid && x.mid === m.id)
    if (b) {
      const list = b.listeners || (b.listeners = [])
      if (!list.includes(m.uid)) list.push(m.uid)
    }
    const msg = (table.messages || []).find(x => x.vid && x.vid === m.id)
    if (msg) {
      const list = msg.listeners || (msg.listeners = [])
      if (!list.includes(m.uid)) list.push(m.uid)
    }
  })
})
onUnmounted(() => { offVoice?.(); offVoicePlayed?.() })

watch(() => (table.messages || []).length, () => {
  const m = table.messages[table.messages.length - 1]
  if (m) showChatFx(m)
})

const isNiu = computed(() => st.value?.gameType === 'niuniu')

/* ---- 牛牛抢庄 / 开牌（官方 NiuniuTableHost 流程） ---- */
const NIU_BID_OPTIONS = [
  { mult: 0, label: '不抢' },
  { mult: 1, label: '1 倍' },
  { mult: 2, label: '2 倍' },
  { mult: 3, label: '3 倍' },
  { mult: 4, label: '4 倍' }
]
const niuBankerMode = ref(localStorage.getItem('u5_niu_mode') || 'rob')
const niuBidBusy = ref(false)
const niuBidTick = ref(Date.now())
const niuBidLeft = computed(() => {
  const dl = Number(st.value?.deadline || 0)
  if (!dl) return 6
  return Math.max(0, Math.ceil((dl - niuBidTick.value) / 1000))
})
const myBid = computed(() => Number(mySeatObj.value?.bid || 0))
const NIU_STAKE_OPTIONS = [1, 2, 5, 10]
const myStake = computed(() => Number(mySeatObj.value?.stake || 0))
const isBankerSeat = computed(() => !!st.value && Number(st.value.dealer) === Number(st.value.mySeat) && (st.value.mySeat ?? -1) >= 0)
const niuStakeLeft = computed(() => {
  const dl = Number(st.value?.deadline || 0)
  if (!dl) return 6
  return Math.max(0, Math.ceil((dl - niuBidTick.value) / 1000))
})
const canStake = computed(() => {
  const t = st.value
  const me = mySeatObj.value
  return !!t && t.state === 'betting' && !!me && (me.cards?.length || 0) >= 5 && Number(t.dealer) !== Number(t.mySeat)
})
async function doNiuStake(mult: number) {
  if (niuBidBusy.value) return
  niuBidBusy.value = true
  try {
    await socket.call('table.niuStake', { mult })
    sfx.chip()
  } catch (e: any) {
    ui.err(e?.message || '押注失败')
  } finally {
    niuBidBusy.value = false
  }
}
const canBid = computed(() => {
  const t = st.value
  const me = mySeatObj.value
  return !!t && t.state === 'bidding' && !!me && (me.cards?.length || 0) >= 5 && myBid.value === 0
})
function niuSeatMult(s: any) {
  if (!s) return 0
  const base = Number(s.mult || 1)
  const banker = Number(st.value?.bankerBid || 1)
  const m = base * (banker > 1 ? banker : 1)
  return m > 0 ? m : 0
}
async function doNiuBid(mult: number) {
  if (niuBidBusy.value) return
  niuBidBusy.value = true
  try {
    await socket.call('table.niuBid', { mult })
    sfx.tap()
  } catch (e: any) {
    ui.err(e?.message || '抢庄失败')
  } finally {
    niuBidBusy.value = false
  }
}

// 开牌：自己 5 张逐张翻开，再显示牛型
const niuRevealStage = ref(0)
const niuRevealCount = ref(0)
const myNiuCards = computed(() => (mySeatObj.value?.cards || []).slice(0, 5))
const myNiuName = computed(() => String(mySeatObj.value?.handName || ''))
let niuRevealTimer = 0
function runNiuReveal() {
  clearInterval(niuRevealTimer)
  if (!isNiu.value || myNiuCards.value.length < 5) { niuRevealStage.value = 0; return }
  niuRevealStage.value = 1
  niuRevealCount.value = 0
  niuRevealTimer = window.setInterval(() => {
    niuRevealCount.value++
    if (niuRevealCount.value >= 5) {
      clearInterval(niuRevealTimer)
      niuRevealStage.value = 2
      sfx.chip()
      setTimeout(() => { niuRevealStage.value = 0 }, 2200)
    }
  }, 260)
}
function skipNiuReveal() {
  clearInterval(niuRevealTimer)
  niuRevealStage.value = 0
}
watch(() => st.value?.state, nv => {
  if (!isNiu.value) return
  if (nv === 'playing' || nv === 'showdown') runNiuReveal()
  if (nv !== 'bidding' && nv !== 'betting') { clearInterval(niuBidTimer) } else { startNiuBidTick() }
})
let niuBidTimer = 0
function startNiuBidTick() {
  clearInterval(niuBidTimer)
  niuBidTick.value = Date.now()
  niuBidTimer = window.setInterval(() => { niuBidTick.value = Date.now() }, 300)
}
onBeforeUnmount(() => { clearInterval(niuBidTimer); clearInterval(niuRevealTimer) })
const myCards = computed(() => st.value?.myCards || [])
const myHandName = computed(() => {
  const me = st.value?.seats.find(x => x.isMe)
  return me?.show ? (me.handName || '') : ''
})
const heroDocked = computed(() => !!(st.value && st.value.mySeat >= 0))
// 入座入场动画：从旁观切到在座的瞬间点亮 1.6 秒
const justSat = ref(false)
let lastSeatVal = -2
watch(() => st.value?.mySeat, (v) => {
  const cur = v ?? -1
  if (lastSeatVal === -1 && cur >= 0) {
    justSat.value = true
    setTimeout(() => { justSat.value = false }, 1600)
    sfx.win()
  }
  lastSeatVal = cur
}, { immediate: true })
const chipIcon = '/img/chip.png'
const cardBackKey = ref(localStorage.getItem('u5_cardback') || 'emerald')
const cardBackImg = computed(() => cardBackOf(cardBackKey.value).url)
watch(cardBackKey, v => localStorage.setItem('u5_cardback', v))
const bankerName = computed(() => {
  const s = st.value
  if (!s) return '-'
  return s.seats.find(x => x.seat === s.banker)?.nick || '-'
})

const gameLabel = computed(() => {
  const g = st.value?.gameType || 'nlh'
  return { nlh: 'NLH 长牌', short6: '6+ 短牌', omaha: '奥马哈', niuniu: '牛牛' }[g] || 'NLH'
})
// 规则抽屉标题：官方 NLH 桌写「德州牌桌规则 / 标准德州规则与 10 级牌型图解」，其它玩法跟随玩法名
const rulesTitle = computed(() => {
  const g = st.value?.gameType || 'nlh'
  return { nlh: '德州牌桌规则', short6: '短牌牌桌规则', omaha: '奥马哈牌桌规则', niuniu: '牛牛牌桌规则' }[g] || '德州牌桌规则'
})
const rulesSub = computed(() => {
  const g = st.value?.gameType || 'nlh'
  return { nlh: '标准德州规则与 10 级牌型图解', short6: '短牌专属牌型大小与 10 级牌型图解', omaha: '四张底牌取两张与 10 级牌型图解', niuniu: '牛牛比牌规则与倍数说明' }[g] || '标准德州规则与 10 级牌型图解'
})


const rebuyOptions = computed(() => {
  const s = st.value
  if (!s) return []
  const min = s.minBuyin
  const max = s.maxBuyin
  return Array.from(new Set([min, Math.round((min + max) / 2), max])).sort((a, b) => a - b)
})

const callAmt = computed(() => st.value?.actions?.callAmount || 0)

const winnersSet = computed(() => new Set((st.value?.winners || []).map(w => w.seat)))

/* 官方 hero-winner-cards / hero-loser-cards：摊牌时英雄底牌金边或抖动 */
const heroIsWinner = computed(() => {
  const t = st.value
  return !!(t && t.state === 'showdown' && t.mySeat >= 0 && winnersSet.value.has(t.mySeat))
})
const heroIsLoser = computed(() => {
  const t = st.value
  if (!t || t.state !== 'showdown' || t.mySeat < 0) return false
  const me = (t.seats || []).find(x => x.isMe)
  return !!me && !me.folded && !winnersSet.value.has(t.mySeat)
})

const resultTitle = computed(() => {
  const w = st.value?.winners || []
  if (!w.length) return '本手结束'
  if (w.length === 1) return `${w[0].nick} 赢得底池`
  return `${w.map(x => x.nick).join(' / ')} 平分底池`
})
const resultAmount = computed(() => {
  const w = st.value?.winners || []
  const total = w.reduce((a, b) => a + b.amount, 0)
  return '+' + money(total)
})
const winBoard = computed(() => {
  const s = st.value
  if (!s || s.state !== 'showdown') return new Set<string>()
  const board = new Set(s.board || [])
  const out = new Set<string>()
  for (const w of s.winners || []) {
    for (const c of w.best || []) if (board.has(c)) out.add(c)
  }
  return out
})
const resultStyle = computed(() => {
  const s = st.value
  const w = s?.winners?.[0]
  const ps = w ? positionedSeats.value.find(x => x.seat === w.seat) : null
  if (!ps) return { left: '50%', top: '30%', transform: 'translateX(-50%)' }
  const top = parseFloat(ps.pos.top)
  const left = Math.min(76, Math.max(24, parseFloat(ps.pos.left)))
  // 上方座位把结算胶囊放到名牌下方，避免压住头像与亮牌
  const below = top < 36
  // 英雄自己赢时底牌上提了 40px，胶囊再抬高一点，留出底牌空间
  const up = ps.isMe ? 168 : 128
  return { left: left + '%', top: `calc(${top}% ${below ? '+' : '-'} ${below ? 112 : up}px)`, transform: 'translateX(-50%)' }
})
const resultTier = computed(() => {
  const name = (st.value?.winners || [])[0]?.handName || ''
  const tier = BIG_TIERS[name]
  return tier ? 'tier-' + tier : ''
})
const topHandName = computed(() => {
  const w = st.value?.winners || []
  const name = w[0]?.handName
  return name || ''
})
const resultSub = computed(() => {
  const s = st.value
  if (!s) return ''
  const me = s.seats.find(x => x.isMe)
  if (!me) return s.winners?.[0]?.handName || ''
  const mine = (s.winners || []).find(w => w.seat === me.seat)
  if (mine) return `我赢得 ${money(mine.amount)} · ${mine.handName || '对手弃牌'}`
  if ((me.committed || 0) === 0) return '本手未入池 · 等待下一局'
  return me.folded ? '本手已弃牌' : `本手 -${money(Math.max(0, me.committed))}`
})

const insLeft = computed(() => {
  const d = st.value?.insurance?.deadline || 0
  if (!d) return 0
  return Math.max(0, Math.ceil((d - now.value) / 1000))
})
const insDash = computed(() => {
  const total = 12
  return Math.round((insLeft.value / total) * 120)
})

/* ---- 官方 WaitingNextHandDialog：已入座但本局进行中（等待下一手） ---- */
const waitNextDismissed = ref(-1)
const heroWaitingNext = computed(() => {
  const t = st.value
  if (!t || t.mySeat < 0) return false
  if (t.state !== 'playing' && t.state !== 'showdown') return false
  const me = (t.seats || []).find(x => x.isMe)
  if (!me) return false
  // 摊牌阶段自己仍无牌 => 本局已结算，正在发下一手
  return (me.cards?.length || 0) === 0 && (me.cardCount || 0) === 0
})
const waitNextOpen = computed(() => heroWaitingNext.value && waitNextDismissed.value !== (st.value?.handNo || 0))
function closeWaitNext() { waitNextDismissed.value = st.value?.handNo || 0 }
watch(() => st.value?.handNo, () => { waitNextDismissed.value = -1 })

/* 等待下一局：结算倒计时进度环（在官方 DOM 之上叠加的动效层） */
const waitTick = ref(Date.now())
const waitPhase = computed(() => (st.value?.state === 'showdown' ? 'deal' : 'play'))
const waitLeftMs = computed(() => {
  const t = st.value
  if (!t?.nextHandAt) return 0
  return Math.max(0, t.nextHandAt - waitTick.value)
})
const waitNextSec = computed(() => Math.max(0, Math.ceil(waitLeftMs.value / 1000)))
// 结算阶段第一次量到的剩余时间作为分母，进度环从 100% 走到底，不多不少正好对上发牌时刻
const waitHoldStart = ref(0)
watch([waitPhase, waitNextOpen], () => {
  if (waitPhase.value === 'deal' && waitNextOpen.value) {
    if (waitLeftMs.value > (waitHoldStart.value || 0)) waitHoldStart.value = waitLeftMs.value
  } else if (waitPhase.value !== 'deal') {
    waitHoldStart.value = 0
  }
}, { immediate: true })
const waitRingPct = computed(() => {
  if (waitPhase.value !== 'deal') return '100%'
  const hold = waitHoldStart.value || st.value?.resultHoldMs || 5200
  const p = hold > 0 ? Math.max(0, Math.min(1, waitLeftMs.value / hold)) : 0
  return Math.round(p * 100) + '%'
})

/* ---- 官方保险弹层：报价 / 三档方案 / 提交 ---- */
const ins = computed(() => st.value?.insurance || null)
const insPlan = ref('')
const insBusy = ref(false)
const insPlanList = computed(() => {
  const o = ins.value
  if (!o) return []
  if (o.plans && o.plans.length) return o.plans
  return [{ key: 'full', label: '满池保险 · 全额', premium: o.premium, payout: o.payout }]
})
watch(() => [ins.value?.deadline, ins.value?.planKey], () => {
  insPlan.value = ins.value?.planKey || insPlanList.value[0]?.key || 'full'
  insBusy.value = false
}, { immediate: true })
async function buyInsurance() {
  if (insBusy.value) return
  insBusy.value = true
  try {
    await table.decideInsurance(true, insPlan.value)
  } finally {
    setTimeout(() => { insBusy.value = false }, 600)
  }
}

const remainText = computed(() => {
  const s = st.value
  if (!s?.expireAt) return '--:--'
  const left = Math.max(0, Math.floor((s.expireAt - now.value) / 1000))
  const h = Math.floor(left / 3600)
  const m = Math.floor((left % 3600) / 60)
  const sec = left % 60
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(h)}:${p(m)}:${p(sec)}`
})

interface PositionedSeat extends SeatState {
  pos: { left: string; top: string }
  timerPct: number
  rel: number
  betPos: string
  dealerPos: string
  fanSide: string
}

// 官方座位百分位（相对毛毡区 880px：y 起点 52，见 docs/official-table-baseline.md）
const SEAT_MAP: Record<number, [number, number][]> = {
  2: [[50, 84], [50, 9]],
  3: [[50, 84], [13, 47], [87, 44]],
  4: [[50, 84], [13, 47], [50, 9], [87, 44]],
  5: [[50, 84], [13, 60], [13, 28], [50, 9], [87, 42]],
  6: [[50, 84], [13, 60], [13, 28], [50, 9], [87, 26], [87, 60]],
  7: [[50, 77], [13.5, 66], [12, 46], [13, 25], [50, 11], [87, 25], [88, 46]],
  8: [[50, 77], [13.5, 66], [12, 46], [13, 25], [50, 11], [87, 25], [88, 46], [86.5, 66]],
  9: [[50, 77], [13, 72], [11, 55], [11.5, 37], [15, 22], [85, 22], [88.5, 37], [89, 55], [87, 72]]
}
const SEAT_TABLE_TOP = 52
const SEAT_TABLE_H = 880
const SCENE_H_PX = 932

// 毛毡区百分比 → 场景坐标百分比
function feltTop(pct: number) {
  return ((SEAT_TABLE_TOP + (SEAT_TABLE_H * pct) / 100) / SCENE_H_PX) * 100
}

function slotOf(n: number, rel: number): [number, number] {
  const map = SEAT_MAP[n] || SEAT_MAP[6]
  return map[Math.min(rel, map.length - 1)]
}

const positionedSeats = computed<PositionedSeat[]>(() => {
  const s = st.value
  if (!s) return []
  const n = s.seatCount || 6
  const myIdx = s.mySeat >= 0 ? s.mySeat : 0
  return s.seats.map(seat => {
    const rel = ((seat.seat - myIdx) % n + n) % n
    const [x, y] = slotOf(n, rel)
    const atTop = y < 20
    const onLeft = x < 40
    const onRight = x > 60
    // 下注筹码与庄家按钮一律朝桌面内侧摆放（官方左右列都是内侧）
    const betPos = rel === 0
      ? 'bet-pos-hero'
      : atTop
        ? 'bet-pos-top-center'
        : onLeft
          ? 'bet-pos-right'
          : 'bet-pos-left'
    const dealerPos = rel === 0 ? 'dealer-pos-hero' : atTop ? 'dealer-pos-top' : onLeft ? 'dealer-pos-right' : 'dealer-pos-left'
    const fanSide = onLeft ? 'fan-left' : atTop ? 'fan-top' : 'fan-right'
    return {
      ...seat,
      rel,
      betPos,
      dealerPos,
      fanSide,
      pos: {
        left: x.toFixed(2) + '%',
        top: feltTop(y).toFixed(2) + '%'
      },
      timerPct: 0
    }
  })
})

const vacantSeats = computed(() => {
  const t = st.value
  if (!t) return []
  const n = t.seatCount || 6
  const myIdx = t.mySeat >= 0 ? t.mySeat : 0
  const occupied = new Set((t.seats || []).map(x => x.seat))
  const out: { seat: number; pos: { left: string; top: string } }[] = []
  for (let i = 0; i < n; i++) {
    const seatNo = (myIdx + i) % n
    if (occupied.has(seatNo)) continue
    const [x, y] = slotOf(n, i)
    out.push({ seat: seatNo, pos: { left: x.toFixed(2) + '%', top: feltTop(y).toFixed(2) + '%' } })
  }
  return out
})

const spectating = computed(() => {
  const t = st.value
  return !!(t && t.mySeat < 0)
})
function quickSit() {
  const t = st.value
  if (!t) return
  const first = vacantSeats.value[0]
  openSit(first ? first.seat : -1)
}

const sitOpen = ref(false)
const sitSeat = ref(-1)
const sitAmount = ref(0)
const sitBusy = ref(false)
const sitMin = computed(() => st.value?.minBuyin || 0)
const sitMax = computed(() => {
  const t = st.value
  if (!t) return 0
  return t.maxBuyin && t.maxBuyin > t.minBuyin ? t.maxBuyin : t.minBuyin || 0
})
const sitPct = computed(() => {
  const span = sitMax.value - sitMin.value
  if (span <= 0) return 0
  return Math.max(0, Math.min(100, ((sitAmount.value - sitMin.value) / span) * 100))
})
const sitBB = computed(() => {
  const bb = st.value?.bb || 0
  return bb > 0 ? Math.round(sitAmount.value / bb) : 0
})
const sitInsufficient = computed(() => (auth.profile?.balance || 0) < sitMin.value)
const sitOver = computed(() => sitAmount.value > (auth.profile?.balance || 0))
const sitRatio = computed(() => {
  const bb100 = Math.max(sitMin.value, Math.min(sitMax.value, (st.value?.bb || 1) * 100))
  if (Math.abs(sitAmount.value - sitMin.value) < 1) return 'min'
  if (Math.abs(sitAmount.value - sitMax.value) < 1) return 'max'
  if (Math.abs(sitAmount.value - bb100) < 1) return '100bb'
  if (Math.abs(sitAmount.value - (sitMin.value + sitMax.value) / 2) < 1) return 'half'
  return ''
})
function onSitSlide(e: Event) {
  const v = Number((e.target as HTMLInputElement).value)
  if (!Number.isNaN(v)) sitAmount.value = Math.max(sitMin.value, Math.min(sitMax.value, v))
  sfx.tap()
}
function setSitRatio(k: string) {
  const t = st.value
  const bb = t?.bb || 1
  const val = k === 'min'
    ? sitMin.value
    : k === 'max'
      ? sitMax.value
      : k === '100bb'
        ? Math.max(sitMin.value, Math.min(sitMax.value, bb * 100))
        : Math.round((sitMin.value + sitMax.value) / 2)
  sitAmount.value = val
  sfx.tap()
}
function openSit(seat: number) {
  if ((st.value?.mySeat ?? -1) >= 0) { ui.toast('你已经在座', 'ok'); return }
  sitSeat.value = seat
  sitAmount.value = sitMin.value || 0
  sitOpen.value = true
  sfx.tap()
}
async function doSit() {
  const t = st.value
  if (!t || !sitAmount.value) return
  sitBusy.value = true
  try {
    await table.join(t.roomNo, sitAmount.value, sitSeat.value, false)
    sitOpen.value = false
    ui.ok('已入座')
    sfx.chip()
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    sitBusy.value = false
  }
}

function displayCards(s: SeatState): string[] {
  if (s.cards && s.cards.length) return s.cards
  // 弃牌也要留牌背：官方在盖牌上挂「加密牌背」角标，可付费查看
  if (s.cardCount > 0) {
    return Array.from({ length: s.cardCount }, () => '')
  }
  return []
}

function bubblePos(s: PositionedSeat) {
  return s.isMe
    ? { top: '-60px', left: '50%', transform: 'translateX(-50%)' }
    : { top: '-30px', left: '50%', transform: 'translateX(-50%)' }
}

const preAction = ref<'' | 'check' | 'fold'>('')
const canPreAction = computed(() => {
  const t = st.value
  return !!(t && t.mySeat >= 0 && t.state === 'playing' && t.toAct !== t.mySeat && !raiseOpen.value)
})
function togglePre(v: 'check' | 'fold') {
  preAction.value = preAction.value === v ? '' : v
  sfx.tap()
}
watch(
  () => [st.value?.toAct, st.value?.state, st.value?.handNo] as const,
  () => {
    const t = st.value
    if (!t || t.mySeat < 0 || t.state !== 'playing' || t.toAct !== t.mySeat) return
    const want = preAction.value
    if (!want) return
    preAction.value = ''
    const act = want === 'check' ? (t.actions?.canCheck ? 'check' : 'call') : 'fold'
    doAct(act).catch(() => {})
  }
)
watch(() => st.value?.handNo, () => { preAction.value = '' })

const myTurnNow = computed(() => {
  const t = st.value
  return !!(t && t.mySeat >= 0 && t.state === 'playing' && t.toAct === t.mySeat && t.actions)
})
const showCluster = computed(() => myTurnNow.value && !isNiu.value && !raiseOpen.value)
const showDock = computed(() => !!(raiseOpen.value || isNiu.value || !myTurnNow.value))
function actVisible(s: SeatState) {
  const a = s.lastAction
  if (!a) return false
  if (a === 'sb' || a === 'bb' || a === 'post') return false
  // 抢庄阶段结束、进入开牌后不再挂「抢庄」气泡，避免整桌一直顶着倍数
  if ((a === 'bid' || a === 'stake' || a === 'banker') && st.value?.state !== 'bidding' && st.value?.state !== 'betting') return false
  return true
}
function actionLabel(s: SeatState) {
  // 牛牛：抢庄 / 押注 倍数气泡
  if (s.lastAction === 'bid') return (Number(s.bid || 0) > 0 ? `抢庄 ×${s.bid}` : '不抢')
  if (s.lastAction === 'stake') return `押注 ×${s.stake || 1}`
  if (s.lastAction === 'banker') return `坐庄 ×${st.value?.bankerBid || 1}`
  if (s.lastAction === 'allin') return '全下 ' + chips(s.lastAmount)
  if (s.lastAction === 'raise' || s.lastAction === 'bet') return actionText(s.lastAction) + ' ' + chips(s.lastAmount)
  if (s.lastAction === 'call') return '跟注 ' + chips(s.lastAmount)
  return actionText(s.lastAction)
}

// ===== 官方垂直加注器：滑轨 200 / 手柄 44 / 行程 156 =====
const railEl = ref<HTMLElement | null>(null)
const railDrag = ref(false)
const RAIL_TRAVEL = 156
const raiseMin = computed(() => st.value?.actions?.minRaiseTo || 0)
const raiseMax = computed(() => st.value?.actions?.maxRaiseTo || 1)
const raiseRatio = computed(() => {
  const span = Math.max(1, raiseMax.value - raiseMin.value)
  return Math.min(1, Math.max(0, (raiseTo.value - raiseMin.value) / span))
})
const handleBottom = computed(() => Math.round(raiseRatio.value * RAIL_TRAVEL))
const ratioBottom = computed(() => Math.round(32 + raiseRatio.value * RAIL_TRAVEL))
const ratioPct = computed(() => {
  const pot = Math.max(1, (st.value?.pot || 0) + callAmt.value)
  return Math.round((raiseTo.value / pot) * 100)
})
const raiseStepIcon = computed(() => {
  const units = Math.max(1, Math.round((raiseTo.value - raiseMin.value) / Math.max(1, st.value?.bb || 1)))
  return units >= 10 ? 'All' : String(units)
})
const railDots = computed(() => {
  const n = 11
  const fill = Math.round(raiseRatio.value * (n - 1))
  return Array.from({ length: n }, (_, i) => i <= fill)
})
function setRaiseByRatio(r: number) {
  const v = raiseMin.value + Math.min(1, Math.max(0, r)) * (raiseMax.value - raiseMin.value)
  raiseTo.value = clampRaise(Math.round(v))
}
function railPos(ev: PointerEvent) {
  const el = railEl.value
  if (!el) return 0
  const rect = el.getBoundingClientRect()
  const y = ev.clientY - rect.top
  return 1 - Math.min(1, Math.max(0, (y - 22) / Math.max(1, rect.height - 44)))
}
function railDown(ev: PointerEvent) {
  railDrag.value = true
  ;(ev.currentTarget as HTMLElement).setPointerCapture?.(ev.pointerId)
  setRaiseByRatio(railPos(ev))
  sfx.tap()
}
function railMove(ev: PointerEvent) {
  if (!railDrag.value) return
  setRaiseByRatio(railPos(ev))
}
function railUp() {
  railDrag.value = false
}
function raisePreset(k: number) {
  const pot = (st.value?.pot || 0) + callAmt.value
  const target = k >= 9 ? raiseMax.value : pot * k + callAmt.value
  raiseTo.value = clampRaise(Math.max(raiseMin.value, Math.round(target)))
  sfx.tap()
}

function clampRaise(v: number) {
  const a = st.value?.actions
  if (!a) return v
  return Math.max(a.minRaiseTo, Math.min(a.maxRaiseTo, Math.round(v)))
}

function openRaise() {
  const a = st.value?.actions
  if (!a) return
  raiseTo.value = clampRaise(Math.round((st.value!.pot + callAmt.value) * 0.66))
  raiseOpen.value = true
  sfx.tap()
}

async function doAct(action: string) {
  try {
    await table.act(action)
    raiseOpen.value = false
    if (action === 'fold') sfx.fold()
  } catch (e: any) {
    ui.err(e.message)
  }
}

async function confirmRaise() {
  try {
    await table.act('raise', raiseTo.value)
    raiseOpen.value = false
  } catch (e: any) {
    ui.err(e.message)
  }
}

async function allIn() {
  if (tablePrefs.value.confirmAllin) {
    const ok = window.confirm('确认全下？全下后将推入全部筹码')
    if (!ok) return
  }
  try {
    await table.act('allin')
    raiseOpen.value = false
    sfx.allin()
  } catch (e: any) {
    ui.err(e.message)
  }
}

const isHost = computed(() => {
  const t = st.value
  const uid = auth.profile?.uid
  return !!(t && t.hostUid && uid && t.hostUid === uid)
})
function goRecharge() {
  router.push('/profile?recharge=1')
}
async function dissolveRoom() {
  const t = st.value
  if (!t) return
  try {
    await socket.call('lobby.closeRoom', { roomNo: t.roomNo })
    ui.ok('牌桌已解散，筹码已退回钱包')
    router.push('/private')
  } catch (e: any) {
    ui.err(e?.message || '解散失败')
  }
}
async function toggleSitOut() {
  try {
    await table.sitOut(!mySeatObj.value?.sittingOut)
    ui.ok(mySeatObj.value?.sittingOut ? '已暂离牌桌' : '已回座')
  } catch (e: any) {
    ui.err(e.message)
  }
}

function openRebuy() {
  menu.value = false
  rebuyAmt.value = rebuyOptions.value[0] || 0
  rebuyOpen.value = true
}

function askCashout() {
  cashoutOpen.value = true
  rebuyOpen.value = false
  sfx.tap()
}
async function doCashout() {
  cashoutOpen.value = false
  const s = st.value
  const me = mySeatObj.value
  const ins = Number((me as any)?.insPayout || 0)
  // 先离桌再展示结算：留在此页，用户关掉结算卡后再回大厅
  await table.leave()
  const buyinTotal = Number(me?.buyin || 0)
  const stackLeft = Number(me?.stack || 0)
  settle.value = s && me ? {
    // 真实盈亏按「离桌筹码 − 累计带入」算，避免只统计摊牌赢池造成的偏差
    net: stackLeft - buyinTotal,
    buyin: buyinTotal,
    stack: stackLeft,
    hands: Number(me.hands || 0),
    ins,
    session: sessionHands.value.slice(-8),
    roomName: s.name
  } : { net: 0, buyin: 0, stack: 0, hands: 0, ins: 0, session: [], roomName: '牌桌' }
  sfx.win()
}

const rebuyBusy = ref(false)
async function doRebuy() {
  if (rebuyBusy.value) return
  rebuyBusy.value = true
  try {
    await table.rebuy(rebuyAmt.value)
    ui.ok('带入成功')
    rebuyOpen.value = false
    sfx.chip()
    auth.refresh()
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    rebuyBusy.value = false
  }
}
// 官方 RebuyDialog「退码」：把桌上多余筹码提前退回钱包（服务端保留至少 1BB）
async function doReduceChips(amount: number) {
  if (rebuyBusy.value || !(amount > 0)) return
  rebuyBusy.value = true
  try {
    const d = await socket.call<{ amount: number; stack: number; balance: number }>('table.reduceChips', { amount })
    ui.ok(`已退回 ${money(d.amount)} 到钱包`)
    rebuyOpen.value = false
    sfx.chip()
    auth.applyBalance({ balance: d.balance })
  } catch (e: any) {
    ui.err(e?.message || '退码失败')
  } finally {
    rebuyBusy.value = false
  }
}

async function showReport() {
  menu.value = false
  try {
    await table.loadReport()
    reportOpen.value = true
  } catch (e: any) {
    ui.err(e.message)
  }
}

let swipeY = 0
function swipeStart(e: PointerEvent) { swipeY = e.clientY }
function swipeEnd(e: PointerEvent) {
  // 上滑收起结算横幅；下滑/轻点保持原有「查看回顾」行为
  if (swipeY - e.clientY > 36) showResult.value = false
}

// ===== 语音消息（按住说话） =====
const recording = ref(false)
const recMs = ref(0)
let mediaRec: MediaRecorder | null = null
let recChunks: Blob[] = []
let recTimer: number | null = null
let recStarted = 0
const voiceBubbles = ref<{ id: number; mid: string; uid: number; seat: number; nick: string; data: string; dur: number; playing?: boolean; played?: boolean; progress?: number; listeners?: number[] }[]>([])
let voiceSeq = 1
const voiceCache = new Map<string, HTMLAudioElement>()

async function startVoice() {
  if (recording.value) return
  if (!navigator.mediaDevices?.getUserMedia || typeof MediaRecorder === 'undefined') {
    ui.err('当前环境不支持录音')
    return
  }
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    recChunks = []
    mediaRec = new MediaRecorder(stream, { audioBitsPerSecond: 16000 })
    mediaRec.ondataavailable = e => { if (e.data.size) recChunks.push(e.data) }
    mediaRec.start()
    recStarted = Date.now()
    recMs.value = 0
    recording.value = true
    recTimer = window.setInterval(() => {
      recMs.value = Date.now() - recStarted
      if (recMs.value >= 10000) stopVoice()
    }, 100)
  } catch {
    ui.err('需要麦克风权限')
  }
}

function cancelVoice() {
  if (!recording.value) return
  finishVoice(true)
}

async function stopVoice() {
  if (!recording.value) return
  finishVoice(false)
}

async function finishVoice(cancel: boolean) {
  recording.value = false
  if (recTimer) { clearInterval(recTimer); recTimer = null }
  const dur = Math.min(10000, Date.now() - recStarted)
  const rec = mediaRec
  mediaRec = null
  if (!rec) return
  if (cancel || dur < 400) {
    rec.stream.getTracks().forEach(t => t.stop())
    if (!cancel) ui.toast('说话时间太短', 'info')
    return
  }
  const stopped = new Promise<void>(res => { rec.onstop = () => res() })
  rec.stop()
  await stopped
  rec.stream.getTracks().forEach(t => t.stop())
  const blob = new Blob(recChunks, { type: rec.mimeType || 'audio/webm' })
  if (blob.size > 150 * 1024) { ui.err('语音过长'); return }
  const buf = await blob.arrayBuffer()
  let bin = ''
  const bytes = new Uint8Array(buf)
  for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i])
  const b64 = btoa(bin)
  try {
    await socket.call('table.voice', { data: b64, dur })
  } catch (e: any) {
    ui.err(e.message)
  }
}

function showVoice(m: any) {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const t = st.value
  if (!t) return
  const seat = t.seats.find(x => x.uid === m.uid)?.seat
  if (seat === undefined) return
  const id = voiceSeq++
  const bubble = { id, mid: String(m.id || ''), uid: m.uid, seat, nick: m.nick, data: m.data, dur: m.dur, listeners: [] as number[] }
  voiceBubbles.value.push(bubble)
  // 同步进聊天消息流：聊天抽屉里用官方 .voice-bubble-body 展示，可点播
  if (!(table.messages || []).some(x => x.vid && x.vid === bubble.mid)) {
    table.messages.push({
      uid: m.uid, nick: m.nick, avatar: m.avatar || '', content: '', type: 'voice',
      time: Number(m.time || Date.now()), data: m.data, dur: m.dur, vid: bubble.mid, listeners: []
    } as any)
    if (table.messages.length > 80) table.messages.shift()
  }
  setTimeout(() => { voiceBubbles.value = voiceBubbles.value.filter(x => x.id !== id) }, 8000)
  playVoice(bubble)
}

function playVoice(v: any) {
  try {
    let audio = voiceCache.get(v.data)
    if (!audio) {
      audio = new Audio('data:audio/webm;base64,' + v.data)
      voiceCache.set(v.data, audio)
      if (voiceCache.size > 20) {
        const first = voiceCache.keys().next().value
        if (first) voiceCache.delete(first)
      }
    }
    const target = voiceBubbles.value.find(x => x.data === v.data) || (v.mid ? v : null)
    if (target) { target.playing = true; target.played = false; target.progress = 0 }
    audio.onplay = () => { if (target) { target.playing = true; target.progress = 0 } }
    audio.ontimeupdate = () => {
      if (!target || !audio.duration) return
      target.progress = Math.min(100, (audio.currentTime / audio.duration) * 100)
    }
    audio.onended = () => {
      if (target) { target.playing = false; target.played = true; target.progress = 100 }
      // 别人发的语音播完 → 回报「已听」
      const mid = (v as any).mid || target?.mid
      if (mid && (v as any).uid !== mySeatObj.value?.uid) socket.call('table.voicePlayed', { id: mid }).catch(() => {})
    }
    audio.currentTime = 0
    audio.play().catch(() => { if (target) target.playing = false })
  } catch { /* 忽略播放异常 */ }
}

// ===== 长按看牌 =====
const peeking = ref(false)
const faceHero = computed(() => {
  const t = st.value
  if (!t) return false
  if (peeking.value) return true
  // 摊牌或本手结束直接明牌
  if (t.state === 'showdown') return true
  if (t.mySeat >= 0 && t.seats.find(x => x.isMe)?.show) return true
  return false
})
const showPeekHint = ref(localStorage.getItem('u5_peek_hint') !== 'done')
function startPeek() {
  peeking.value = true
  if (showPeekHint.value) {
    showPeekHint.value = false
    localStorage.setItem('u5_peek_hint', 'done')
  }
}
function endPeek() { peeking.value = false }

// ===== 双击桌面 = 过牌/跟注 =====
let lastTap = 0
function onSceneTap(e: PointerEvent) {
  const target = e.target as HTMLElement
  if (target.closest('button, input, .sheet, .mask, .action-cluster, .vertical-bet-overlay, .seat')) return
  const now = Date.now()
  if (now - lastTap < 320 && myTurnNow.value) {
    lastTap = 0
    const t = st.value
    if (!t?.actions) return
    doAct(t.actions.canCheck ? 'check' : 'call').catch(() => {})
    sfx.chip()
    return
  }
  lastTap = now
}

function onTwoFinger(e: TouchEvent) {
  // 双指轻触桌面唤起菜单（移动端快捷手势）
  if (!gestureOn.value) return
  if (e.touches && e.touches.length === 2) {
    menu.value = true
    sfx.tap()
  }
}

function openLastHand() {
  const id = st.value?.handId
  menu.value = false
  if (!id) {
    // 牌谱落库有毫秒级延迟，先提示并跳到战绩页，避免点了没反应
    ui.toast('牌谱生成中，已在战绩页等你', 'info', 2200)
    router.push('/records')
    return
  }
  router.push('/records/hand/' + id)
}

function copyTableLink() {
  const no = st.value?.roomNo
  if (!no) return
  const url = `${location.origin}/table/${no}`
  navigator.clipboard?.writeText(url).then(() => ui.ok('牌桌链接已复制')).catch(() => ui.err('复制失败'))
  menu.value = false
  sfx.tap()
}

function invite() {
  const url = `${location.origin}/table/${st.value?.roomNo}`
  navigator.clipboard?.writeText(url).then(() => ui.ok('专属牌桌链接已复制')).catch(() => ui.err('复制失败'))
  menu.value = false
}

function switchCloth() {
  menu.value = false
  clothSheet.value = true
  sfx.tap()
}
function setCloth(key: string) {
  cloth.value = key
  localStorage.setItem('u5_cloth', key)
  ui.ok('已切换桌布：' + (CLOTHS.find(c => c.key === key)?.name || ''))
  sfx.tap()
}
function setCardBack(key: string) {
  cardBackKey.value = key
  ui.ok('已切换牌背：' + cardBackOf(key).name)
  sfx.tap()
}

async function leave() {
  await table.leave()
  router.replace('/lobby')
}

// 官方设置抽屉底部：站起围观（离座但留在本桌） / 离开牌桌（在座先走退码结算）
async function standUpToWatch(skipPreview = false) {
  const no = st.value?.roomNo
  prefsOpen.value = false
  if (!no) return
  if (st.value && st.value.mySeat >= 0 && !skipPreview) { openExitDialog(false); return }
  try {
    await table.leave()
    await table.join(no, 0, -1, true)
    const t = st.value
    if (t && t.mySeat >= 0) {
      // 手中还有牌（已下注），服务端会在本手结束后自动离座
      ui.toast('本手结束后自动站起，先以围观身份看牌', 'info', 2800)
    } else {
      ui.ok('已站起，切换为围观模式')
    }
  } catch (e: any) {
    ui.err(e?.message || '站起失败')
  }
}
// 菜单「退出房间」：在座先走退码结算，围观状态直接离开
const exitOpen = ref(false)
const exitBusy = ref(false)
const exitLeaving = ref(false)
const exitPreview = ref<any>(null)

// 打开官方退桌结算弹层（SquidExitDialog），拿服务端实时预览
async function openExitDialog(leaving: boolean) {
  if (st.value && st.value.mySeat < 0) { leave(); return }
  exitLeaving.value = leaving
  try {
    exitPreview.value = await socket.call('table.exitPreview', { leaving })
    exitOpen.value = true
    sfx.tap()
  } catch (e: any) {
    ui.err(e?.message || '结算预览失败')
    if (leaving) leave()
  }
}
async function doExitConfirm() {
  if (exitBusy.value) return
  exitBusy.value = true
  try {
    if (exitLeaving.value) {
      exitOpen.value = false
      await doCashout()
    } else {
      exitOpen.value = false
      await standUpToWatch(true)
    }
  } finally {
    exitBusy.value = false
  }
}
function exitRoom() {
  if (st.value && st.value.mySeat >= 0) { openExitDialog(true); return }
  leave()
}
function leaveFromSettings() {
  prefsOpen.value = false
  if (st.value && st.value.mySeat >= 0) { openExitDialog(true); return }
  leave()
}

function spawnCoins() {
  const layer = { id: Date.now(), style: {} }
  coins.value = []
  for (let i = 0; i < 16; i++) {
    const angle = Math.random() * Math.PI * 2
    const dist = 60 + Math.random() * 120
    coins.value.push({
      id: i,
      style: {
        left: '50%',
        top: '42%',
        '--dx': Math.cos(angle) * dist + 'px',
        '--dy': Math.sin(angle) * dist + 'px',
        animationDelay: (i * 18) + 'ms'
      }
    })
  }
  setTimeout(() => { coins.value = [] }, 1400)
}

let timer = 0
let waitTimer = 0
let hurryKey = ''
function tickClock() {
  now.value = Date.now()
  // 官方 hurry 语音：自己剩余时间不足 6 秒时提醒一次
  const t = st.value
  const me = t?.seats?.find(x => x.isMe)
  if (!t || !me?.deadline || t.state !== 'playing' || t.toAct !== t.mySeat) return
  const left = (me.deadline - now.value) / 1000
  const key = `${t.handNo}:${me.deadline}`
  if (left > 0 && left <= 6 && hurryKey !== key) {
    hurryKey = key
    sfx.voice('hurry', { minGap: 0 })
  }
}

watch(() => st.value?.state, (nv, ov) => {
  if (nv === 'showdown' && ov !== 'showdown') {
    // 结算时收起聊天抽屉，避免挡住摊牌与赢家光效
    chatSheet.value = false
    clothSheet.value = false
    showResult.value = true
    setTimeout(() => { showResult.value = false }, 4200)
    setTimeout(() => { try { triggerBigWin() } catch { /* 特效失败不影响牌局 */ } }, 620)
  }
})

watch(() => st.value?.winEv, (w) => {
  if (w && w.amount > 0) {
    spawnCoins()
    if (tablePrefs.value.vibration) (navigator as any).vibrate?.([18, 40, 24])
  }
})

watch(() => st.value?.handNo, () => {
  raiseOpen.value = false
  showResult.value = false
})

onMounted(async () => {
  const roomNo = Number(route.params.roomNo)
  if (!roomNo) { router.replace('/lobby'); return }
  if (!table.state || table.state.roomNo !== roomNo) {
    let matched = false
    try {
      await socket.call('table.resume')
      matched = table.state?.roomNo === roomNo
    } catch { /* 没有可恢复的牌局 */ }
    if (!matched) {
      // 目标房间与当前对局不一致：以旁观身份进入目标桌，避免串桌
      try {
        await table.join(roomNo, 0, -1, true)
      } catch (e: any) {
        ui.err(e.message)
        router.replace('/lobby')
        return
      }
    }
  }
  timer = window.setInterval(tickClock, 1000)
  waitTimer = window.setInterval(() => { waitTick.value = Date.now() }, 200)
  const check = () => {
    winSize.value = { w: window.innerWidth, h: window.innerHeight }
    portrait.value = window.innerHeight > window.innerWidth * 1.05
  }
  check()
  window.addEventListener('resize', check)
  ;(window as any).__u5check = check
})

onBeforeUnmount(() => {
  clearInterval(timer)
  clearInterval(waitTimer)
  const check = (window as any).__u5check
  if (check) window.removeEventListener('resize', check)
})
</script>

<style scoped>
.bubble-enter-active, .bubble-leave-active { transition: opacity .2s ease, transform .2s ease; }
.bubble-enter-from, .bubble-leave-to { opacity: 0; transform: translateY(4px) scale(.9); }
.loading-layer { position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 10px; background: rgba(2, 12, 20, .9); z-index: 60; }
.spinner { width: 34px; height: 34px; border-radius: 50%; border: 3px solid rgba(150, 225, 255, .2); border-top-color: #6df3b0; animation: spin .8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
.close-x { width: 32px; height: 32px; border-radius: 10px; background: rgba(255, 255, 255, .06); display: flex; align-items: center; justify-content: center; color: var(--muted); flex: none; }
.menu-group { margin-top: 12px; }
.menu-title { font-size: 12px; font-weight: 700; color: var(--gold); margin: 0 2px 6px; letter-spacing: .3px; }
.menu-item { display: flex; align-items: center; gap: 10px; padding: 11px 12px; border-radius: 12px; background: rgba(8, 26, 17, .66); border: 1px solid var(--line); margin-bottom: 7px; }
.menu-item:active { background: rgba(16, 44, 30, .8); }
.menu-item b { font-size: 14px; }
.report-head, .report-row { display: grid; grid-template-columns: 1.6fr .8fr .6fr .9fr; gap: 6px; align-items: center; font-size: 12px; }
.report-head { color: var(--muted); padding: 8px 4px; border-bottom: 1px solid var(--line); }
.report-row { padding: 9px 4px; border-bottom: 1px solid rgba(122, 196, 156, .07); }
.ellipsis { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ins-row { display: flex; justify-content: space-between; padding: 7px 0; border-bottom: 1px dashed rgba(122, 196, 156, .12); font-size: 13px; }
.hash-box { display: flex; flex-direction: column; gap: 6px; padding: 10px 12px; border-radius: 12px; background: rgba(4, 16, 10, .7); border: 1px solid var(--line); margin-bottom: 10px; }
.hash-box code { font-size: 11px; color: #9ef7c4; word-break: break-all; line-height: 1.5; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
.rule-list { display: flex; flex-direction: column; gap: 8px; }
.rule-row { display: flex; gap: 9px; font-size: 13px; color: #d6ecdf; line-height: 1.5; }
.rule-idx { width: 20px; height: 20px; flex: none; border-radius: 50%; background: rgba(46, 224, 122, .16); color: #7ff0b0; font-size: 11px; display: flex; align-items: center; justify-content: center; font-weight: 700; }
.buyin-box { padding: 10px 6px; border-radius: 12px; display: flex; flex-direction: column; align-items: center; gap: 2px;
  background: rgba(8, 26, 17, .8); border: 1px solid var(--line); font-size: 15px; font-weight: 800; }
.buyin-box.active { border-color: rgba(46, 224, 122, .6); background: rgba(16, 60, 38, .8); }
.grid2 { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }

/* ===== 全下冲击波 ===== */
.allin-wave {
  position: absolute; width: 0; height: 0; pointer-events: none; z-index: 42;
  transform: translate(-50%, -50%);
}
.allin-wave .ring {
  position: absolute; left: 50%; top: 50%; width: 90px; height: 90px; margin: -45px 0 0 -45px;
  border-radius: 50%; border: 2px solid rgba(255, 214, 102, .9);
  box-shadow: 0 0 24px rgba(255, 190, 60, .7), inset 0 0 20px rgba(255, 190, 60, .5);
  animation: allinRing .9s cubic-bezier(.2,.7,.3,1) forwards;
}
.allin-wave .ring.r2 { animation-delay: .16s; border-color: rgba(255, 120, 90, .85); }
.allin-wave .tag-text {
  position: absolute; left: 50%; top: 50%; transform: translate(-50%, -50%);
  font-size: 12px; font-weight: 900; letter-spacing: 1.6px; color: #2a1200; white-space: nowrap;
  padding: 3px 9px; border-radius: 10px;
  background: linear-gradient(180deg, #ffe9a8, #f0a53c);
  box-shadow: 0 4px 16px rgba(255, 170, 40, .55);
  animation: allinTag 1.15s cubic-bezier(.2,.8,.3,1) forwards;
}
@keyframes allinRing {
  0% { transform: scale(.35); opacity: .95; }
  70% { opacity: .5; }
  100% { transform: scale(2.5); opacity: 0; }
}
@keyframes allinTag {
  0% { transform: translate(-50%, -120%) scale(.5); opacity: 0; }
  22% { transform: translate(-50%, -150%) scale(1.08); opacity: 1; }
  72% { transform: translate(-50%, -205%) scale(1); opacity: 1; }
  100% { transform: translate(-50%, -250%) scale(.95); opacity: 0; }
}

/* ===== 大牌型结算影院 ===== */
.win-cinema {
  position: absolute; inset: 0; z-index: 260; display: flex; align-items: center; justify-content: center;
  background: radial-gradient(circle at 50% 42%, rgba(0, 40, 22, .58), rgba(0, 0, 0, .78) 72%);
  overflow: hidden;
}
.win-cinema .rays {
  position: absolute; width: 220vmax; height: 220vmax; left: 50%; top: 42%; transform: translate(-50%, -50%);
  background: conic-gradient(from 0deg, rgba(255, 226, 150, .13) 0 2.5deg, transparent 2.5deg 13deg);
  animation: raysSpin 22s linear infinite; opacity: .75;
  filter: blur(.3px);
  mask-image: radial-gradient(circle at 50% 50%, #000 12%, transparent 62%);
  -webkit-mask-image: radial-gradient(circle at 50% 50%, #000 12%, transparent 62%);
}
@keyframes raysSpin { to { transform: translate(-50%, -50%) rotate(360deg); } }
.cinema-inner::before {
  content: ""; position: absolute; left: 50%; top: 50%; width: 420px; height: 420px;
  transform: translate(-50%, -50%); border-radius: 50%; pointer-events: none; z-index: -1;
  background: radial-gradient(circle, rgba(255, 206, 110, .3), rgba(255, 206, 110, 0) 65%);
  animation: haloPulse 2.6s ease-in-out infinite;
}
@keyframes haloPulse {
  0%, 100% { opacity: .6; transform: translate(-50%, -50%) scale(.92); }
  50% { opacity: 1; transform: translate(-50%, -50%) scale(1.06); }
}
.cinema-inner {
  position: relative; z-index: 2; display: flex; flex-direction: column; align-items: center; gap: 6px;
  padding: 22px 26px; border-radius: 20px; text-align: center;
  background: linear-gradient(170deg, rgba(20, 56, 38, .92), rgba(6, 22, 14, .95));
  border: 1px solid rgba(255, 216, 117, .5);
  box-shadow: 0 20px 60px rgba(0, 0, 0, .6), inset 0 1px 0 rgba(255, 255, 255, .12);
  animation: cinemaPop .5s cubic-bezier(.2,1.2,.3,1) both;
}
@keyframes cinemaPop {
  0% { transform: scale(.72); opacity: 0; }
  100% { transform: scale(1); opacity: 1; }
}
.cinema-kicker { font-size: 10px; font-weight: 800; letter-spacing: 2.6px; color: rgba(255, 216, 117, .8); }
.cinema-hand {
  font-size: 30px; font-weight: 900; letter-spacing: 2px; line-height: 1.15;
  background: linear-gradient(180deg, #fff6d8, #f0b64a);
  -webkit-background-clip: text; background-clip: text; color: transparent;
  text-shadow: 0 6px 24px rgba(240, 182, 74, .38);
}
.cinema-amount { font-size: 20px; font-weight: 900; color: #7ff0b0; letter-spacing: .6px; }
.cinema-who { font-size: 11.5px; color: #b9d2c4; }
.win-cinema.legend .cinema-inner { border-color: rgba(255, 236, 170, .8); box-shadow: 0 20px 70px rgba(255, 190, 60, .32); }
.win-cinema.legend .cinema-hand { font-size: 34px; }
.win-cinema.epic .cinema-inner { border-color: rgba(150, 220, 255, .55); }
.win-cinema .spark {
  position: absolute; left: 50%; top: 42%; width: 5px; height: 5px; border-radius: 50%;
  background: #ffe9a8; opacity: 0;
  transform: rotate(calc(var(--i) * 25.7deg)) translateY(-40px);
  animation: sparkFly 1.6s ease-out forwards; animation-delay: calc(var(--i) * 22ms);
}
@keyframes sparkFly {
  0% { opacity: 0; transform: rotate(calc(var(--i) * 25.7deg)) translateY(-30px) scale(.6); }
  30% { opacity: 1; }
  100% { opacity: 0; transform: rotate(calc(var(--i) * 25.7deg)) translateY(-190px) scale(1.1); }
}
.cinema-enter-active, .cinema-leave-active { transition: opacity .25s ease; }
.cinema-enter-from, .cinema-leave-to { opacity: 0; }

/* 低功耗/减少动效偏好：保留信息，去掉大位移动画 */
@media (prefers-reduced-motion: reduce) {
  .win-cinema .rays, .win-cinema .spark, .cinema-inner::before { animation: none; }
  .allin-wave .ring, .allin-wave .tag-text { animation-duration: .42s; }
  .cinema-inner { animation: none; }
}

/* 举报玩家 */
.rp-list { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 10px; max-height: 132px; overflow-y: auto; }
.rp-item {
  display: flex; align-items: center; gap: 6px; height: 36px; padding: 0 10px 0 4px;
  border-radius: 10px; background: rgba(255, 255, 255, .045); border: 1px solid var(--line);
  font-size: 12px; color: #d8ecdf; max-width: 48%;
}
.rp-item.on { border-color: rgba(255, 150, 120, .8); background: rgba(255, 120, 90, .14); }
.rp-item i { font-size: 9px; font-style: normal; color: #ffd875; }
.rp-reasons { display: grid; grid-template-columns: repeat(2, 1fr); gap: 6px; margin-bottom: 10px; }
.rp-chip {
  height: 38px; border-radius: 10px; font-size: 12.5px; font-weight: 700; color: #cfe6d9;
  background: rgba(255, 255, 255, .05); border: 1px solid var(--line);
}
.rp-chip.on { color: #2a1000; background: linear-gradient(180deg, #ffd0a8, #f09060); border-color: transparent; }
.rp-detail { height: 78px; padding: 10px 12px; resize: none; line-height: 1.5; }
</style>
