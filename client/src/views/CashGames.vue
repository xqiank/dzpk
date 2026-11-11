<template>
  <section class="lobby-page lobby-view">
    <div class="mobile-lobby-shell">
      <div class="lobby-fixed-header">
        <!-- 轮播 + 公告跑马灯 -->
        <div class="lobby-carousel-wrap has-banners">
          <div class="carousel-inner">
            <div
              v-for="(b, i) in banners"
              :key="b.key"
              class="carousel-slide-item"
              :class="{ active: carouselIndex === i }"
              role="button"
              tabindex="0"
              :aria-label="b.title"
              :aria-hidden="carouselIndex !== i"
              :style="{ background: b.bg }"
              @click="onBanner(b)"
            >
              <img :src="b.img" :alt="b.title" class="banner-full-img" />
              <div class="banner-image-copy">
                <span class="banner-tag">{{ b.tag }}</span>
                <strong class="banner-heading">{{ b.title }}</strong>
                <p class="banner-subcopy">{{ b.sub }}</p>
              </div>
            </div>
            <div class="carousel-dots-bar">
              <button
                v-for="(b, i) in banners"
                :key="'dot' + b.key"
                type="button"
                class="carousel-dot"
                :class="{ active: carouselIndex === i }"
                :aria-label="`轮播第${i + 1}张：${b.title}`"
                :aria-pressed="carouselIndex === i"
                @click.stop="setCarousel(i)"
              ></button>
            </div>
          </div>
          <div class="lobby-broadcast-ticker">
            <div class="ticker-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24" width="13" height="13" fill="currentColor"><path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"></path></svg>
            </div>
            <div class="ticker-track">
              <div class="ticker-marquee-inner">
                <span class="ticker-text">{{ announce }}</span>
                <span class="ticker-text" aria-hidden="true">{{ announce }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 用户 / 俱乐部 Hero 卡 -->
        <header class="lobby-club-hero-card">
          <div class="hero-top-row">
            <div class="hero-left-col">
              <div class="club-avatar-box" @click="go('/profile')">
                <img class="club-avatar-img" :src="avatarUrl(auth.profile?.avatar)" :alt="`${auth.profile?.nick || '游客'}头像`" />
                <span class="vip-avatar-badge">VIP</span>
              </div>
              <div class="club-info-wrap">
                <div class="club-identity-row">
                  <div class="club-name-row" title="个人中心" @click="go('/profile')">
                    <span class="crown-icon">
                      <svg viewBox="0 0 24 24" width="14" height="14" fill="#ffd700"><path d="M5 16L3 5l5.5 5L12 4l3.5 6L21 5l-2 11H5zm14 3c0 .6-.4 1-1 1H6c-.6 0-1-.4-1-1v-1h14v1z"></path></svg>
                    </span>
                    <strong class="club-name-text">{{ auth.profile?.nick || '游客' }}</strong>
                    <span class="club-switch-arrow">›</span>
                  </div>
                  <div class="club-id-row">
                    <span class="club-id-text">ID: {{ auth.profile?.uid || '----' }}</span>
                    <button class="club-id-copy" title="复制ID" @click="copyId">
                      <svg viewBox="0 0 24 24" width="11" height="11" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
                    </button>
                  </div>
                </div>
                <div class="club-asset-row">
                  <div class="compact-asset-badge" title="点击充值 / 管理资产" @click="openRecharge">
                    <span class="compact-tether-icon">₮</span>
                    <span class="compact-asset-val">{{ hideBal ? '****' : money(auth.profile?.balance) }}</span>
                    <button class="compact-asset-eye" :title="hideBal ? '显示余额' : '隐藏余额'" @click.stop="hideBal = !hideBal">
                      <svg viewBox="0 0 24 24" width="13" height="13" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                    </button>
                    <button class="compact-asset-plus" title="快速充币" @click.stop="openRecharge">
                      <svg viewBox="0 0 24 24" width="11" height="11" fill="none" stroke="currentColor" stroke-width="3.2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <div class="hero-right-col">
              <button class="green-square-btn" title="官方圈交流社区" @click="showChat = true">
                <div class="square-icon-wrap">
                  <svg viewBox="0 0 24 24" width="19" height="19" aria-hidden="true">
                    <defs><linearGradient id="green-action-v-0" gradientUnits="userSpaceOnUse" x1="2" y1="2" x2="22" y2="22"><stop offset="0%" stop-color="#fff6d6"></stop><stop offset="50%" stop-color="#fde68a"></stop><stop offset="100%" stop-color="#d97706"></stop></linearGradient></defs>
                    <g stroke="url(#green-action-v-0)" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></g>
                  </svg>
                </div>
                <span class="square-btn-text">官方圈</span>
              </button>
              <button class="green-square-btn" title="会员" @click="showVip = true">
                <div class="square-icon-wrap">
                  <svg viewBox="0 0 24 24" width="19" height="19" aria-hidden="true">
                    <defs><linearGradient id="green-action-v-1" gradientUnits="userSpaceOnUse" x1="2" y1="2" x2="22" y2="22"><stop offset="0%" stop-color="#fff6d6"></stop><stop offset="50%" stop-color="#fde68a"></stop><stop offset="100%" stop-color="#d97706"></stop></linearGradient></defs>
                    <g stroke="url(#green-action-v-1)" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></g>
                  </svg>
                </div>
                <span class="square-btn-text">会员</span>
              </button>
              <button class="green-square-btn" title="查看战绩记录" @click="go('/records')">
                <div class="square-icon-wrap">
                  <svg viewBox="0 0 24 24" width="19" height="19" aria-hidden="true">
                    <defs><linearGradient id="green-action-v-2" gradientUnits="userSpaceOnUse" x1="2" y1="2" x2="22" y2="22"><stop offset="0%" stop-color="#fff6d6"></stop><stop offset="50%" stop-color="#fde68a"></stop><stop offset="100%" stop-color="#d97706"></stop></linearGradient></defs>
                    <g stroke="url(#green-action-v-2)" fill="none" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3.5" y="3.5" width="17" height="17" rx="3.5"></rect><path d="M8 16V12M12 16V8M16 16V10.5" stroke-width="2.2"></path></g>
                  </svg>
                </div>
                <span class="square-btn-text">牌局记录</span>
              </button>
            </div>
          </div>
        </header>

        <!-- 玩法分类 + 房卡 / 开局 -->
        <div class="game-category-bar">
          <div class="category-tabs">
            <button
              v-for="t in filters"
              :key="t.key"
              class="cat-tab"
              :class="{ active: lobby.filter === t.key }"
              @click="lobby.setFilter(t.key)"
            >
              <span>{{ t.label }}</span>
              <i class="active-bar"></i>
            </button>
          </div>
          <div class="category-actions">
            <button class="join-room-key-btn" title="输入房间号加入" @click="go('/private')">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor"><path d="M12.65 10C11.83 7.67 9.61 6 7 6c-3.31 0-6 2.69-6 6s2.69 6 6 6c2.61 0 4.83-1.67 5.65-4H17v4h4v-4h2v-4H12.65zM7 14c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2z"></path></svg>
            </button>
            <button class="create-room-plus-btn" title="创建牌桌" @click="goCreate">
              <span class="plus-sym">＋</span><span>开局</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 桌子列表 -->
      <main class="lobby-scrollable-body">
        <div class="poker-tables-container">
          <!-- 限时积分冲榜赛 -->
          <article v-if="mtt" class="mtt-lobby-banner-card" @click="openMtt">
            <div class="mtt-banner-left">
              <div class="mtt-banner-badges">
                <span class="mtt-badge-gold">
                  <svg viewBox="0 0 24 24" width="12" height="12" fill="currentColor" style="vertical-align: -1.5px; margin-right: 3px"><path d="M19 5h-2V3H7v2H5c-1.1 0-2 .9-2 2v1c0 2.55 1.92 4.63 4.39 4.94.63 1.5 1.98 2.63 3.61 2.96V19H7v2h10v-2h-4v-3.1c1.63-.33 2.98-1.46 3.61-2.96C19.08 12.63 21 10.55 21 8V7c0-1.1-.9-2-2-2zM5 8V7h2v3.82C5.84 10.4 5 9.3 5 8zm14 0c0 1.3-.84 2.4-2 2.82V7h2v1z"></path></svg>
                  限时积分冲榜赛
                </span>
                <span class="mtt-badge-free">免费报名</span>
              </div>
              <strong class="mtt-banner-title">{{ mtt.title }}</strong>
              <div class="mtt-banner-meta">
                <span>{{ mtt.schedule }}</span>
                <span>·</span>
                <span>初始 {{ mtt.initialPoints.toLocaleString('en-US') }} 赛事积分</span>
                <span>·</span>
                <span>已报 {{ mtt.registered }} 人</span>
              </div>
            </div>
            <div class="mtt-banner-right">
              <div class="mtt-prize-pill">
                <small>当前总奖池</small>
                <strong>₮ {{ chips(mtt.prizePool) }}</strong>
              </div>
              <button type="button" class="mtt-enter-btn" :class="{ signed: mtt.signed }" @click.stop="openMtt">
                {{ mtt.signed ? '已报名' : '免费报名' }}
              </button>
            </div>
          </article>

          <!-- 房间卡 -->
          <article
            v-for="r in shownRooms"
            :key="r.roomNo"
            class="poker-table-card"
            :class="[gameModeClass(r), occClass(r)]"
            :data-room-no="r.roomNo"
            @click="openRoom(r)"
          >
            <div class="card-occupancy-bar" :class="occClass(r)" aria-hidden="true">
              <span class="card-occupancy-fill" :style="{ height: Math.round((r.online / r.seatCount) * 100) + '%' }"></span>
            </div>
            <div class="card-watermark-overlay"></div>
            <div class="card-tether-watermark" aria-hidden="true">
              <svg viewBox="0 0 100 100" fill="currentColor"><path d="M50 0C22.4 0 0 22.4 0 50s22.4 50 50 50 50-22.4 50-50S77.6 0 50 0zm21.6 28.5H54.2v6.2h17.4v6.2H54.2v5.3c15.2.8 26.6 4.3 26.6 8.5 0 4.2-11.4 7.7-26.6 8.5v16.1H45.8V63.2c-15.2-.8-26.6-4.3-26.6-8.5 0-4.2 11.4-7.7 26.6-8.5v-5.3H28.4v-6.2h17.4v-6.2H28.4v-6.2h43.2v6.2zm-21.6 30c12.3 0 22.5-2.2 24.8-5.3-2.3-3.1-12.5-5.3-24.8-5.3s-22.5 2.2-24.8 5.3c2.3 3.1 12.5 5.3 24.8 5.3z"></path></svg>
            </div>
            <div class="card-donut-col">
              <div class="table-topdown-box" :class="{ 'online-bump': lobby.bumpOf(r.roomNo).online, 'full-bump': lobby.bumpOf(r.roomNo).full }">
                <svg class="seat-arc-svg" viewBox="-6 -6 67 51" aria-hidden="true">
                  <path class="seat-arc-glow" d="M 16.26 -0.49 A 20.25 20.25 0 0 1 19.50 -0.75 L 35.50 -0.75 A 20.25 20.25 0 0 1 38.74 39.49"></path>
                  <path class="seat-arc-line" d="M 16.26 -0.49 A 20.25 20.25 0 0 1 19.50 -0.75 L 35.50 -0.75 A 20.25 20.25 0 0 1 38.74 39.49" stroke="url(#lobbySeatArcGrad)"></path>
                </svg>
                <span
                  v-for="(d, i) in seatDots(r)"
                  :key="i"
                  class="s-dot"
                  :class="{ on: i < r.online, 'just-on': i === r.online - 1 && lobby.bumpOf(r.roomNo).online }"
                  :style="{ left: d.x + 'px', top: d.y + 'px' }"
                ></span>
                <b class="table-seat-num">{{ r.online }}/{{ r.seatCount }}</b>
              </div>
            </div>
            <div class="card-info-col">
              <div class="room-title-line">
                <span class="room-title">{{ roomTitle(r) }}</span>
                <span v-if="extraTags(r).length" class="room-flag-mini">{{ flagText(r) }}</span>
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
              <span class="vpip-text" :class="{ 'vpip-bump': lobby.bumpOf(r.roomNo).vpip }">入池率 {{ r.vpip }}%</span>
            </div>
          </article>

          <div v-if="!shownRooms.length && lobby.loaded" class="lobby-empty-box">
            <b>当前筛选下没有可坐的牌桌</b>
            <span>换个玩法，或去私人局自己开一桌</span>
            <div class="row gap-8" style="margin-top:12px">
              <button class="btn btn-sm btn-ghost" @click="lobby.setFilter('all')">看全部</button>
              <button class="btn btn-sm btn-primary" @click="goCreate">自己开一桌</button>
            </div>
          </div>
        </div>
        <div v-if="shownRooms.length" class="lobby-list-end-hint"><span>已展示全部 {{ shownRooms.length }} 张牌桌</span></div>
      </main>

      <!-- 好运牌型奖池浮标 -->
      <div class="lobby-bottom-right-floats">
        <div class="jp-float-entrance-wrap">
          <button type="button" class="lobby-jackpot-float-widget" title="Jackpot 好运牌型大奖池" @click="showJackpot = true">
            <div class="jp-chip-shower-stage" aria-hidden="true">
              <span v-for="i in 6" :key="i" :class="'jp-chip-drop-x c' + i">
                <span class="jp-chip-drop-y"><img src="/assets/table/chip.png" alt="" class="jp-room-chip-img" draggable="false" /></span>
              </span>
            </div>
            <div class="jp-widget-shake-body">
              <div class="jp-widget-emblem-stage">
                <span class="jp-widget-aura-ring" aria-hidden="true"></span>
                <span class="jp-widget-sparkle s1" aria-hidden="true">✦</span>
                <span class="jp-widget-sparkle s2" aria-hidden="true">✦</span>
                <svg class="jp-widget-logo-svg" viewBox="0 0 88 68" aria-hidden="true">
                  <defs>
                    <linearGradient id="jpGoldRimA" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" stop-color="#fffbeb"></stop><stop offset="45%" stop-color="#fde68a"></stop><stop offset="100%" stop-color="#b4831e"></stop></linearGradient>
                    <linearGradient id="jpObsidianCardA" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" stop-color="#194a36"></stop><stop offset="45%" stop-color="#0a2218"></stop><stop offset="100%" stop-color="#040d09"></stop></linearGradient>
                    <linearGradient id="jpIvoryCardA" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" stop-color="#ffffff"></stop><stop offset="70%" stop-color="#fef3c7"></stop><stop offset="100%" stop-color="#fde68a"></stop></linearGradient>
                    <linearGradient id="jpSpadeGoldA" x1="0%" y1="0%" x2="0%" y2="100%"><stop offset="0%" stop-color="#ffffff"></stop><stop offset="40%" stop-color="#fde68a"></stop><stop offset="100%" stop-color="#d4af37"></stop></linearGradient>
                  </defs>
                  <circle cx="44" cy="36" r="26" fill="none" stroke="url(#jpGoldRimA)" stroke-width="1.2" stroke-dasharray="3 2" opacity="0.68"></circle>
                  <g transform="translate(29, 35) rotate(-16) translate(-18, -25)">
                    <rect x="0" y="0" width="35" height="48" rx="5.5" fill="url(#jpIvoryCardA)" stroke="url(#jpGoldRimA)" stroke-width="2"></rect>
                    <text x="5" y="13" fill="#e11d48" font-size="11" font-weight="900" font-family="Arial, sans-serif">A</text>
                    <path d="M 8.5 15.5 C 7 13.8, 4.5 14.2, 4.5 16.5 C 4.5 18.5, 8.5 21, 8.5 22.5 C 8.5 21, 12.5 18.5, 12.5 16.5 C 12.5 14.2, 10 13.8, 8.5 15.5 Z" fill="#e11d48"></path>
                    <path d="M 18 22 C 15.2 18.5, 10.5 19.5, 10.5 24 C 10.5 28, 18 32.5, 18 35.5 C 18 32.5, 25.5 28, 25.5 24 C 25.5 19.5, 20.8 18.5, 18 22 Z" fill="#e11d48"></path>
                  </g>
                  <g transform="translate(53, 34) rotate(11) translate(-19, -26)">
                    <rect x="0" y="0" width="37" height="50" rx="6" fill="url(#jpObsidianCardA)" stroke="url(#jpGoldRimA)" stroke-width="2.2"></rect>
                    <rect x="2.5" y="2.5" width="32" height="45" rx="4" fill="none" stroke="rgba(253,230,138,0.32)" stroke-width="0.8"></rect>
                    <text x="5" y="13.5" fill="#fde68a" font-size="11.5" font-weight="900" font-family="Arial, sans-serif">A</text>
                    <path d="M 18.5 15 C 18.5 15, 8.5 23.5, 8.5 29 C 8.5 32.5, 11.2 34.8, 14.2 34.8 C 16 34.8, 17.2 33.8, 17.8 32.8 L 16.2 38.5 L 20.8 38.5 L 19.2 32.8 C 19.8 33.8, 21 34.8, 22.8 34.8 C 25.8 34.8, 28.5 32.5, 28.5 29 C 28.5 23.5, 18.5 15, 18.5 15 Z" fill="url(#jpSpadeGoldA)"></path>
                  </g>
                  <path d="M 44 1 L 46.2 6.8 L 52 9 L 46.2 11.2 L 44 17 L 41.8 11.2 L 36 9 L 41.8 6.8 Z" fill="#fffbeb"></path>
                  <circle cx="44" cy="9" r="2" fill="#00ff88"></circle>
                </svg>
              </div>
              <div class="jp-widget-prize-plaque">
                <span class="jp-plaque-top-ribbon">好运牌型奖池</span>
                <span class="jp-plaque-shimmer-clip" aria-hidden="true"><span class="jp-plaque-shimmer"></span></span>
                <div class="jp-plaque-amount-row">
                  <span class="jp-plaque-currency">₮</span><span class="jp-plaque-digits">{{ jackpot ? jackpot.poolText : '10,777.36' }}</span>
                </div>
                <div class="jp-plaque-sub-tag"><span class="jp-sub-dot"></span><span>立即参与大奖</span></div>
              </div>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- 买入弹层 -->
    <div v-if="pick" class="mask" @click.self="pick = null">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:14px">
          <div class="col">
            <b class="h2">{{ pick.name }}</b>
            <span class="sm muted">盲注 {{ money(pick.sb, { bare: true }) }}/{{ money(pick.bb, { bare: true }) }} · {{ pick.seatCount }}人桌</span>
          </div>
          <button class="close-x" @click="pick = null"><Icon name="close" :size="18" /></button>
        </div>
        <div v-if="pick.online >= pick.seatCount" class="full-hint">
          当前牌桌已满员，可先旁观等待空位
        </div>
        <div v-else class="buyin-boxes">
          <button v-for="amt in buyinOptions" :key="amt" class="buyin-box" :class="{ active: buyin === amt }"
                  :disabled="amt > (auth.profile?.balance || 0)" @click="buyin = amt">
            <span class="tiny muted">买入</span>
            <b>{{ money(amt) }}</b>
            <span class="tiny" :style="{ color: 'var(--muted-2)' }">{{ Math.round(amt / pick.bb) }} BB</span>
          </button>
        </div>
        <div class="row-between" style="margin:14px 2px 8px">
          <span class="sm muted">我的余额</span>
          <b class="gold num">{{ money(auth.profile?.balance) }}</b>
        </div>
        <button v-if="pick.online < pick.seatCount" class="btn btn-primary btn-block" style="height:50px"
                :disabled="busy" @click="enterRoom(true)">
          带入并入座
        </button>
        <button class="btn btn-primary btn-block" style="height:50px" :disabled="busy" @click="enterRoom(false)">
          旁观牌桌
        </button>
      </div>
    </div>

    <!-- 公平竞技说明 -->
    <div v-if="showFair" class="mask" @click.self="showFair = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:12px">
          <b class="h2">公平竞技 · 区块发牌</b>
          <button class="close-x" @click="showFair = false"><Icon name="close" :size="18" /></button>
        </div>
        <div class="fair-list">
          <div class="fair-row"><b>可验证发牌</b><span>每手牌的牌堆哈希在开局前生成，结算后可复核。</span></div>
          <div class="fair-row"><b>安全卫士</b><span>同环境 IP 两人进房会实时提醒全房用户。</span></div>
          <div class="fair-row"><b>禁止作弊</b><span>AI、伙牌、虫洞、辅助软件一律封禁处理。</span></div>
        </div>
        <button class="btn btn-primary btn-block" style="height:46px;margin-top:14px" @click="showFair = false">知道了</button>
      </div>
    </div>

    <!-- 官方交流社区（官方 OfficialGroupChatDrawer） -->
    <OfficialChatDrawer
      :visible="showChat"
      :member-count="mtt ? mtt.registered : 0"
      @close="showChat = false"
      @openLogin="router.push('/login')"
    />

    <!-- 锦标赛详情（官方 MttCampaignModal） -->
    <MttCampaignModal
      :show="showMtt"
      :tournament="mttTournament"
      :user-id="auth.profile?.uid"
      @close="showMtt = false"
      @updated="onMttUpdated"
      @enterRoom="onMttEnterRoom"
    />

    <!-- 好运牌型大奖池（官方 JackpotModal） -->
    <JackpotModal
      :visible="showJackpot"
      :status="jackpot"
      :profile="jackpotProfile"
      @close="showJackpot = false"
      @refresh="loadJackpot"
      @joined="onJackpotJoined"
      @openLogin="router.push('/login')"
    />

    <!-- 新人一次性房卡红包（官方 NewbieRedPacketModal） -->
    <NewbieRedPacketModal
      :open="showNewbie"
      @close="showNewbie = false"
      @claimed="onNewbieClaimed"
    />

    <!-- 充值 -->
    <div v-if="showRecharge" class="mask" @click.self="showRecharge = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:14px">
          <b class="h2">立即充值</b>
          <button class="close-x" @click="showRecharge = false"><Icon name="close" :size="18" /></button>
        </div>
        <div class="grid2">
          <button v-for="a in [1000, 5000, 10000, 50000, 100000, 500000]" :key="a" class="buyin-box"
                  :class="{ active: rechargeAmt === a }" @click="rechargeAmt = a">
            <b>{{ money(a) }}</b><span class="tiny muted">USDT</span>
          </button>
        </div>
        <button class="btn btn-primary btn-block" style="height:50px;margin-top:14px" @click="doRecharge">确认充值</button>
      </div>
    </div>

    <!-- VIP -->
    <div v-if="showVip" class="mask" @click.self="showVip = false">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:14px">
          <b class="h2">U5 VIP CLUB</b>
          <button class="close-x" @click="showVip = false"><Icon name="close" :size="18" /></button>
        </div>
        <p class="sm muted">VIP 享受专属牌桌、优先客服、房卡折扣与专属头像框。</p>
        <div class="vip-list">
          <div v-for="v in vipTiers" :key="v.lv" class="vip-row">
            <Icon name="crown" :size="18" :color="v.lv <= (auth.profile?.vip || 0) ? '#efcf8e' : '#5b7a6a'" />
            <div class="col" style="flex:1">
              <b>VIP {{ v.lv }}</b>
              <span class="tiny muted">{{ v.desc }}</span>
            </div>
            <span class="tag" :class="v.lv <= (auth.profile?.vip || 0) ? 'tag-gold' : 'tag-dim'">
              {{ v.lv <= (auth.profile?.vip || 0) ? '已解锁' : '未解锁' }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Icon from '@/components/Icon.vue'
import JackpotModal from '@/components/lobby/JackpotModal.vue'
import MttCampaignModal from '@/components/lobby/MttCampaignModal.vue'
import OfficialChatDrawer from '@/components/lobby/OfficialChatDrawer.vue'
import NewbieRedPacketModal from '@/components/lobby/NewbieRedPacketModal.vue'
import { useAuth } from '@/stores/auth'
import { useLobby, type LobbyRoom } from '@/stores/lobby'
import { useUI } from '@/stores/ui'
import { useTable } from '@/stores/table'
import { socket } from '@/lib/ws'
import { money, chips, remain, gameTagOf, avatarUrl } from '@/lib/format'
import { sfx } from '@/lib/sfx'
import '@/styles/lobby-official.css'
import '@/styles/newbie-redpacket-official.css'

const router = useRouter()
const auth = useAuth()
const lobby = useLobby()
const ui = useUI()
const table = useTable()

const filters = [
  { key: 'all', label: '全部' },
  { key: 'nlh', label: '长牌' },
  { key: 'short6', label: '短牌' },
  { key: 'omaha', label: '奥马哈' },
  { key: 'niuniu', label: '牛牛' }
] as const

const vipTiers = [
  { lv: 1, desc: '每日赠送 1 张房卡' },
  { lv: 2, desc: '专属桌布 + 房卡 9 折' },
  { lv: 3, desc: '专属客服 + 房卡 8 折' },
  { lv: 4, desc: '高额桌准入 + 专属头像框' }
]

const pick = ref<LobbyRoom | null>(null)
const showJackpot = ref(false)
const showNewbie = ref(false)
// 领取成功：本地同步房卡数并收起弹层
function onNewbieClaimed(_cards: number, total: number) {
  if (auth.profile) (auth.profile as any).cards = total
  sfx.play('chips')
}
const showMtt = ref(false)
const mttTournament = ref<any>(null)
const showChat = ref(false)
const mtt = ref<any>(null)
const jackpot = ref<any>(null)

const OFFICIAL_SEAT_DOTS = [
  { x: 16.26, y: -0.49 }, { x: 38.74, y: -0.49 }, { x: 55.75, y: 19.5 },
  { x: 38.74, y: 39.49 }, { x: 16.26, y: 39.49 }, { x: -0.75, y: 19.5 }
]
function seatDots(r: any) {
  const n = Math.max(2, Math.min(9, Number(r?.seatCount) || 6))
  if (n === 6) return OFFICIAL_SEAT_DOTS
  const cx = 27.5
  const cy = 19.5
  const rx = 28.25
  const ry = 20
  return Array.from({ length: n }, (_, i) => {
    const a = (-Math.PI / 2 + (i * Math.PI * 2) / n + Math.PI / n) as number
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
function roomTitle(r: any) {
  return String(r?.name || '')
}

async function loadMTT() {
  try {
    const d = await socket.call<any>('mtt.banner')
    mtt.value = d
    mttTournament.value = d?.tournament || null
  } catch { mtt.value = null; mttTournament.value = null }
}
// 锦标赛详情（官方 MttCampaignModal）
function openMtt() {
  if (!mtt.value) return
  showMtt.value = true
  sfx.play('click')
}
// 赛事桌开启后从弹层直接入场
function onMttEnterRoom(room: any) {
  if (room?.roomNo) {
    showMtt.value = false
    router.push(`/table/${room.roomNo}`)
    return
  }
  ui.toast('赛事桌将在 21:00 准时开启，届时点这里直接入场', 'info', 2600)
}
function onMttUpdated(t: any) {
  if (!t) return
  mttTournament.value = t
  mtt.value = {
    ...(mtt.value || {}),
    signed: !!t.my_registered,
    registered: t.players,
    prizePool: Math.round(Number(t.prize_pool || 0) * 100),
    initialPoints: t.start_chips,
    tournament: t
  }
}
async function loadJackpot() {
  try {
    const d = await socket.call<any>('jackpot.status')
    jackpot.value = { ...d, poolText: money(d.pool) }
  } catch { jackpot.value = null }
}

// Jackpot 参与成功：后端已扣门票并写入倒计时，这里同步奖池数字
function onJackpotJoined(d: any) {
  if (!d) return
  jackpot.value = { ...(jackpot.value || {}), ...d, poolText: money(d.pool) }
  sfx.play('chips')
  auth.refresh()
}
const onlyAvailable = ref(false)
const hideBal = ref(false)
// 官方盲注/买入显示为去尾零形式：0.2/0.4、20
// 官方房特有条目：鱿鱼 / 暴击 / 极速（按房号映射，前端展示层）
const EXTRA_TAGS: Record<number, { text: string; kind: string }[]> = {
  216: [{ text: '鱿鱼', kind: 'squid' }, { text: '暴击', kind: 'boom' }],
  224: [{ text: '暴击', kind: 'boom' }],
  217: [{ text: '极速', kind: 'fast' }]
}
function flagText(r: any) {
  const t = extraTags(r)
  if (!t.length) return ''
  if (r.roomNo === 216) return '鱿鱼暴击'
  return t.map(x => x.text).join(' ')
}

function ringTone(r: any) {
  const gt = r?.gameType
  if (gt === 'omaha') return 'tone-purple'
  if (gt === 'short6') return 'tone-red'
  if (gt === 'niuniu') return 'tone-gold'
  return ''
}

function extraTags(r: any) {
  return EXTRA_TAGS[r.roomNo] || []
}

function blindsOf(r: any) {
  return `${chips(r.sb)}/${chips(r.bb)}`
}
const isGuest = computed(() => !auth.token || auth.profile?.guest === true)
const jackpotProfile = computed(() => ({
  id: auth.profile?.uid,
  nick: auth.profile?.nick,
  wallet: auth.profile?.balance
}))
const gamePill = (g: string) => ({ nlh: 'pill-blue', short6: 'pill-pink', omaha: 'pill-purple', niuniu: 'pill-gold' } as any)[g] || 'pill-blue'
const sortMode = ref(localStorage.getItem('u5_lobby_sort') || 'order')
// 公告：首屏从 /api/config 拉，之后由 WS announcement 事件实时替换（后台改完立刻生效）
const banners = [
  {
    key: 'blockchain', title: '区块发牌', tag: '公平竞技', sub: '发牌过程公开可验证',
    img: '/assets/branding/carousel-blockchain.webp', bg: 'linear-gradient(135deg, #032b1b 0%, #0d6e46 100%)'
  },
  {
    key: 'u5', title: 'U5 德州', tag: '全天开局', sub: '长牌 · 短牌 · 奥马哈 · 牛牛',
    img: '/assets/branding/carousel-club.webp', bg: 'linear-gradient(135deg, #062a1e 0%, #0a5c3c 100%)'
  },
  {
    key: 'club', title: '俱乐部共创', tag: '组队开局', sub: '与您携手共创公平俱乐部',
    img: '/img/carousel-learning.webp', bg: 'linear-gradient(135deg, #12261a 0%, #0f6b45 100%)'
  }
]
const carouselIndex = ref(0)
let carouselTimer = 0
function setCarousel(i: number) {
  carouselIndex.value = ((i % banners.length) + banners.length) % banners.length
  restartCarousel()
}
function restartCarousel() {
  if (carouselTimer) window.clearInterval(carouselTimer)
  carouselTimer = window.setInterval(() => {
    if (document.hidden) return
    carouselIndex.value = (carouselIndex.value + 1) % banners.length
  }, 4200)
}
function onBanner(b: { key: string }) {
  if (b.key === 'blockchain') showFair.value = true
  else go('/club')
}
onMounted(() => restartCarousel())
onBeforeUnmount(() => { if (carouselTimer) window.clearInterval(carouselTimer) })

const announce = ref('IP 两人在房间安全卫士实时追踪提醒全房间用户，禁止 AI、伙牌、虫洞、串通作弊。')
let offAnnounce: (() => void) | null = null
const shownRooms = computed(() => (onlyAvailable.value ? lobby.visible.filter(r => r.online < r.seatCount) : lobby.visible))
const recentIds = ref<number[]>(JSON.parse(localStorage.getItem('u5_recent') || '[]'))
const recentRooms = computed(() =>
  recentIds.value.map(id => lobby.rooms.find(r => r.roomNo === id)).filter(Boolean).slice(0, 4) as LobbyRoom[]
)
const hits = ref<Record<string, number>>(JSON.parse(localStorage.getItem('u5_recent_hits') || '{}'))
const hitOf = (no: number) => hits.value[String(no)] || 0
const SORTS = ['order', 'fill', 'hands', 'pot']
function toggleSort() {
  const i = (SORTS.indexOf(sortMode.value) + 1) % SORTS.length
  sortMode.value = SORTS[i]
  localStorage.setItem('u5_lobby_sort', sortMode.value)
  ui.toast(sortMode.value === 'hands' ? '已按手数排序' : sortMode.value === 'pot' ? '已按底池排序' : '已按人气排序', 'info')
}

function toggleAvailable() {
  onlyAvailable.value = !onlyAvailable.value
  localStorage.setItem('u5_lobby_avail', onlyAvailable.value ? '1' : '0')
}

function rememberRoom(no: number) {
  const next = [no, ...recentIds.value.filter(x => x !== no)].slice(0, 6)
  recentIds.value = next
  localStorage.setItem('u5_recent', JSON.stringify(next))
  hits.value[String(no)] = hitOf(no) + 1
  localStorage.setItem('u5_recent_hits', JSON.stringify(hits.value))
}
const buyin = ref(0)
const busy = ref(false)
const showFair = ref(false)
const showRecharge = ref(false)
const showVip = ref(false)
const rechargeAmt = ref(1000)

const buyinOptions = computed(() => {
  if (!pick.value) return []
  const r = pick.value
  const opts = [r.minBuyin, Math.round((r.minBuyin + r.maxBuyin) / 2), r.maxBuyin]
  return Array.from(new Set(opts)).sort((a, b) => a - b)
})

function sinceText(ts: number) {
  const d = Math.max(0, Math.floor((Date.now() - ts) / 1000))
  if (d < 60) return '刚刚有牌局'
  if (d < 3600) return Math.floor(d / 60) + ' 分钟前'
  return Math.floor(d / 3600) + ' 小时前'
}

function stripeOf(r: LobbyRoom) {
  if (r.gameType === 'short6') return 'pink'
  if (r.gameType === 'omaha') return 'purple'
  if (r.gameType === 'niuniu') return 'orange'
  return 'blue'
}
function ringColor(r: LobbyRoom) {
  const p = r.online / r.seatCount
  if (p >= 0.9) return '#ff9f6b'
  if (p >= 0.5) return '#2ee07a'
  return '#4a8f6b'
}
function ringStyle(_r: LobbyRoom) {
  return {}
}
function dotStyle(n: number, total: number) {
  const angle = ((n - 1) / total) * Math.PI * 2 - Math.PI / 2
  return {
    left: `${42 + 32 * Math.cos(angle) - 3}px`,
    top: `${30 + 20 * Math.sin(angle) - 3}px`
  }
}

function go(p: string) {
  sfx.tap()
  router.push(p)
}

function copyId() {
  const id = String(auth.profile?.uid || '')
  navigator.clipboard?.writeText(id).then(() => ui.ok('已复制 ID')).catch(() => ui.err('复制失败'))
  sfx.tap()
}

function openRoom(r: LobbyRoom) {
  rememberRoom(r.roomNo)
  sfx.tap()
  pick.value = r
  buyin.value = r.minBuyin
}

async function enterRoom(takeSeat: boolean) {
  if (!pick.value || busy.value) return
  busy.value = true
  try {
    await table.join(pick.value.roomNo, takeSeat ? buyin.value : 0, -1, !takeSeat)
    sfx.chip()
    const roomNo = pick.value.roomNo
    pick.value = null
    router.push(`/table/${roomNo}${takeSeat ? '' : '?observe=1'}`)
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    busy.value = false
  }
}

function goCreate() {
  sfx.tap()
  router.push('/private/create')
}

function openRecharge() {
  sfx.tap()
  showRecharge.value = true
}

async function doRecharge() {
  try {
    const d = await socket.call<{ balance: number }>('user.recharge', { amount: rechargeAmt.value })
    auth.applyBalance({ balance: d.balance })
    ui.ok('充值成功')
    showRecharge.value = false
    sfx.win()
  } catch (e: any) {
    ui.err(e.message)
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

let jackpotTimer = 0
onMounted(() => {
  onlyAvailable.value = localStorage.getItem('u5_lobby_avail') === '1'
  lobby.fetch()
  auth.refresh()
  loadMTT()
  loadJackpot()
  // 奖池浮标 30s 轮询一次：别人参与后大厅数字实时涨，切后台不空转
  jackpotTimer = window.setInterval(() => { if (!document.hidden) loadJackpot() }, 30000)
  // 新人房卡红包：未领过就自动弹一次（官方进入大厅即弹）
  if (!sessionStorage.getItem('u5_newbie_rp_seen')) {
    socket.call('redpacket.newbieStatus').then((st: any) => {
      if (st?.can_claim) {
        showNewbie.value = true
        sessionStorage.setItem('u5_newbie_rp_seen', '1')
      }
    }).catch(() => {})
  }
  // 公告首屏拉一次
  fetch('/api/config').then(r => r.json()).then(d => {
    const text = d?.d?.announcement
    if (text) announce.value = text
  }).catch(() => {})
  // 后台改公告 → WS 实时替换
  offAnnounce = socket.on('announcement', (m: any) => {
    if (m?.text) announce.value = String(m.text)
  })
})
onUnmounted(() => { offAnnounce?.(); if (jackpotTimer) window.clearInterval(jackpotTimer) })
</script>

<style scoped>.top { padding: 2px 0 6px; }
.user-main { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.id-line { font-size: 11.5px; color: var(--muted); }
.bal-row { display: flex; align-items: center; gap: 4px; margin-top: 2px; min-width: 0; }
.eye-btn { width: 24px; height: 24px; border-radius: 50%; background: rgba(8,40,30,.9);
  border: 1px solid rgba(120,220,180,.3); display: flex; align-items: center; justify-content: center; color: #bff0d8; }
.plus-btn { width: 26px; height: 26px; border-radius: 50%; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(180deg, #ffe08a, #d8a63c); color: #2a1a02; border: 1px solid rgba(255,246,205,.7);
  box-shadow: 0 4px 10px rgba(216,166,60,.35); }
.user { display: flex; align-items: center; gap: 10px; flex: 1; min-width: 0; }
.user .col { min-width: 0; }
.user .row { min-width: 0; }
.ava-wrap { position: relative; flex: none; }
.vip { position: absolute; bottom: -4px; left: 50%; transform: translateX(-50%);
  font-size: 8px; font-weight: 800; letter-spacing: .3px; color: #2b1d05;
  background: linear-gradient(180deg, #ffe6a8, #d8ab5b); padding: 0 6px; border-radius: 6px;
  border: 1px solid rgba(255,255,255,.5); }
.idbox { display: inline-flex; align-items: center; gap: 3px; white-space: nowrap; flex: none; font-size: 9.5px; color: var(--muted); background: rgba(255,255,255,.05); padding: 1px 6px; border-radius: 6px; border: 1px solid rgba(255,255,255,.06); }
.bal { font-size: 12.5px; font-weight: 800; white-space: nowrap; color: #7ff0b0; display: inline-flex; align-items: center; gap: 3px; }
.coin { position: relative; font-style: normal; display: inline-flex; width: 16px; height: 16px; border-radius: 50%; align-items: center; justify-content: center; font-size: 10px; font-weight: 900; color: #04220f; background: linear-gradient(180deg, #7ff0b0, #12b767); }
.coin::before, .coin::after { content: ""; position: absolute; left: 50%; transform: translateX(-50%); width: 8px; height: 1px; background: currentColor; opacity: .85; }
.coin::before { top: 44%; }
.coin::after { top: 60%; }
.coin.sm { width: 18px; height: 18px; font-size: 11px; color: #04220f; border: 1px solid rgba(255,255,255,.75); background: linear-gradient(180deg, #8dffc4, #17b877); box-shadow: 0 0 10px rgba(23,255,168,.5); }
.mini-btn { display: inline-flex; align-items: center; gap: 2px; height: 20px; padding: 0 6px; white-space: nowrap; flex: none; border-radius: 6px; font-size: 10px; font-weight: 700; color: #04220f; background: linear-gradient(180deg, #55ea9c, #12b767); }
.top-tools { display: flex; gap: 5px; flex: none; }
.horn { font-size: 13px; flex: none; }
.marquee-track { flex: 1; overflow: hidden; }
.marquee-track span { display: block; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  font-size: 10.5px; color: #a9d8bd; }
from { transform: translateX(0); }
to { transform: translateX(-100%); }
.tabs { display: flex; gap: 4px; flex: 1; min-width: 0; overflow-x: auto; scrollbar-width: none; }
.tab.active { color: #241a06; background: linear-gradient(180deg, #f2d493, #cfa151); box-shadow: 0 6px 16px rgba(201,163,92,.24); }
.avail-btn { width: 38px; height: 36px; border-radius: 11px; display: flex; align-items: center; justify-content: center; color: var(--muted); background: rgba(10, 30, 20, .6); border: 1px solid transparent; }
.avail-btn.on { color: #04220f; background: linear-gradient(180deg, #55ea9c, #12b767); }
.fair-banner { padding: 14px 16px; min-height: 139px; display: flex; flex-direction: column; justify-content: flex-end; }
.banner-img { position: absolute; inset: 0; width: 100%; height: 100%; object-fit: cover; }
.fair-banner .banner-txt { position: relative; z-index: 2; }
.fair-banner::after { content: ""; position: absolute; inset: 0; background: linear-gradient(90deg, rgba(3, 20, 13, .92) 0%, rgba(3, 20, 13, .55) 45%, rgba(3, 20, 13, 0) 78%); }
.fair-list { display: flex; flex-direction: column; gap: 10px; }
.fair-row { display: flex; flex-direction: column; gap: 3px; padding: 10px 12px; border-radius: 12px; background: rgba(10, 32, 22, .78); border: 1px solid var(--line); }
.fair-row b { font-size: 13px; color: #d8f6e4; }
.fair-row span { font-size: 12px; color: var(--muted); line-height: 1.6; }
.banner { position: relative; overflow: hidden; border-radius: 18px; padding: 16px; margin-bottom: 14px;
  background: linear-gradient(115deg, rgba(12, 52, 34, .96), rgba(30, 92, 58, .5)), radial-gradient(circle at 88% 20%, rgba(60, 200, 130, .4), transparent 55%);
  border: 1px solid rgba(122, 196, 156, .2); }
.banner-txt { position: relative; display: flex; flex-direction: column; gap: 4px; }
.banner-title { font-size: 22px; font-weight: 900; letter-spacing: .5px; }
.banner-sub { font-size: 12.5px; color: #a9cdb9; }
.banner-cards { position: absolute; right: 18px; top: 10px; display: flex; }
.bc { width: 34px; height: 48px; border-radius: 5px; background: #f6f7f2; color: #12241a; display: flex; align-items: center; justify-content: center; font-size: 22px; box-shadow: 0 6px 14px rgba(0,0,0,.4); }
.bc.red { color: #d8292f; }
.banner-bg { position: absolute; inset: 0; background-image: radial-gradient(rgba(255,255,255,.05) 1px, transparent 1px); background-size: 4px 4px; opacity: .6; }
.rooms { display: flex; flex-direction: column; gap: 12px; }
.list-tools { display: flex; align-items: center; gap: 6px; margin-bottom: 7px; }
.lt-btn { height: 26px; padding: 0 10px; border-radius: 13px; font-size: 11px; color: #a9c9b7;
  background: rgba(10,34,24,.8); border: 1px solid rgba(150,200,172,.18); }
.lt-btn.on { color: #2a1a02; background: linear-gradient(180deg, #ffe08a, #d8a63c); border-color: transparent; font-weight: 700; }
.recent { display: flex; align-items: center; gap: 6px; margin-bottom: 9px; overflow-x: auto; }
.recent-label { font-size: 11px; color: var(--muted-2); flex: none; }
.recent-chip { flex: none; display: inline-flex; align-items: center; gap: 4px; height: 26px; padding: 0 9px;
  border-radius: 13px; font-size: 11px; color: #d9efe2; background: rgba(12, 40, 28, .82);
  border: 1px solid rgba(157, 194, 192, .18); }
.recent-chip b { color: #ffe688; font-weight: 800; }
.recent-chip i { font-style: normal; color: var(--muted); }
.recent-chip em { font-style: normal; margin-left: 2px; font-size: 10px; font-weight: 800; color: #2a1a02; background: linear-gradient(180deg, #ffe08a, #d8a63c); border-radius: 8px; padding: 0 5px; line-height: 15px; }
.room-stripe { position: absolute; left: 0; top: 14px; bottom: 14px; width: 3px; border-radius: 0 3px 3px 0;
  background: linear-gradient(180deg, #34d399, #059669); box-shadow: 0 0 8px rgba(52, 211, 153, .55); }
.room-stripe.pink { background: linear-gradient(180deg, #f472b6, #db2777); box-shadow: 0 0 8px rgba(244, 114, 182, .55); }
.room-stripe.purple { background: linear-gradient(180deg, #60a5fa, #2563eb); box-shadow: 0 0 8px rgba(96, 165, 250, .55); }
.room-stripe.orange { background: linear-gradient(180deg, #fbbf24, #d97706); box-shadow: 0 0 8px rgba(251, 191, 36, .55); }
.pill-blue { background: rgba(38,86,140,.55); color: #a9d0ff; border: 1px solid rgba(90,150,220,.35); }
.pill-pink { background: rgba(130,40,80,.5); color: #ffb0cf; border: 1px solid rgba(220,90,140,.35); }
.pill-purple { background: rgba(90,60,150,.5); color: #d3c1ff; border: 1px solid rgba(150,110,230,.35); }
.pill-gold { background: rgba(120,90,30,.5); color: #ffe6a8; border: 1px solid rgba(220,180,90,.35); }
.close-x { width: 32px; height: 32px; border-radius: 10px; background: rgba(255,255,255,.06); display: flex; align-items: center; justify-content: center; color: var(--muted); }
.buyin-boxes, .grid2 { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.buyin-box { padding: 10px 6px; border-radius: 12px; display: flex; flex-direction: column; align-items: center; gap: 2px;
  background: rgba(8, 26, 17, .8); border: 1px solid var(--line); font-size: 15px; font-weight: 800; }
.buyin-box.active { border-color: rgba(46, 224, 122, .6); background: rgba(16, 60, 38, .8); box-shadow: 0 0 0 2px rgba(46,224,122,.14); }
.fld { display: flex; flex-direction: column; gap: 5px; }
.fld > span { font-size: 11px; color: var(--muted); }
.fld .field { height: 40px; font-size: 14px; }
.switch { width: 46px; height: 26px; border-radius: 20px; background: rgba(255,255,255,.1); position: relative; transition: background .2s; }
.switch i { position: absolute; top: 3px; left: 3px; width: 20px; height: 20px; border-radius: 50%; background: #d6e8dd; transition: transform .2s; }
.switch.on { background: linear-gradient(180deg, #2ee07a, #10a75c); }
.switch.on i { transform: translateX(20px); background: #fff; }
.cost-line { display: flex; align-items: center; gap: 6px; margin-top: 12px; font-size: 12px; color: #cfe8d9; }
.vip-list { display: flex; flex-direction: column; gap: 8px; margin-top: 12px; }
.vip-row { display: flex; align-items: center; gap: 10px; padding: 10px 12px; border-radius: 12px; background: rgba(8, 26, 17, .7); border: 1px solid var(--line); }
.empty { display: flex; flex-direction: column; gap: 10px; }
.empty-card { display: flex; flex-direction: column; align-items: center; gap: 6px; padding: 26px 18px; border-radius: 16px;
  background: linear-gradient(160deg, rgba(17,44,34,.82), rgba(8,24,17,.9)); border: 1px solid var(--line); text-align: center; }
.empty-card b { font-size: 14px; }
.full-hint { padding: 14px; border-radius: 12px; font-size: 12px; color: #ffd98a; background: rgba(60,44,12,.5); border: 1px solid rgba(255,210,130,.24); text-align: center; }
/* ===== 官方实测值（DOM 计算样式 dump，2026-11） ===== */
.banner { padding: 14px 15px; border-radius: 16px; }
.banner-title { font-size: 18px; font-weight: 900; }
.banner-sub { font-size: 13px; font-weight: 600; }
/* ============ 按官方截图实测（430×932 @2x）重排 ============ */
.user-card { padding: 8px 12px; border-radius: 16px; min-height: 73px; }
.id-line { font-size: 11.5px; }
.banner-dots { position: absolute; right: 14px; bottom: 10px; display: flex; gap: 5px; z-index: 3; }
.banner-dots i { width: 7px; height: 4px; border-radius: 3px; background: rgba(255,255,255,.4); }
.banner-dots i.on { width: 16px; background: #ffd875; }
.marquee-track span { font-size: 11.5px; }
.tabs { gap: 2px; }
.tab.active { font-weight: 900; color: #23180a; }
.rooms { display: flex; flex-direction: column; gap: 12px; padding-top: 6px; }
.room-stripe { left: 0; top: 14px; bottom: 14px; width: 4px; border-radius: 0 3px 3px 0; }
/* 官方座位指示器：60×44 椭圆 + 青铜双环 + 6 个点 */
.seat-ring { position: relative; width: 60px; height: 44px; flex: none; border-radius: 999px;
  background: radial-gradient(ellipse at center, #147043, #07331e 66%, #03170d);
  border: 2px solid #a16207;
  box-shadow: 0 0 0 1.2px rgba(253, 230, 138, .68), 0 4px 10px rgba(0, 0, 0, .72), inset 0 0 10px rgba(0, 0, 0, .65); }
.ring-dot.lit { background: #eafff7; border-color: #fff;
  box-shadow: 0 0 7px rgba(160, 255, 225, .95), 0 0 14px rgba(60, 240, 180, .55); }
.ring-dot.dark { background: #4a3218; border-color: rgba(253, 230, 138, .35); }
.ring-dot.d1 { left: 24%; top: 16%; }
.ring-dot.d2 { left: 76%; top: 16%; }
.ring-dot.d3 { left: 12%; top: 50%; }
.ring-dot.d4 { left: 88%; top: 50%; }
.ring-dot.d5 { left: 24%; top: 84%; }
.ring-dot.d6 { left: 76%; top: 84%; }
/* 卡片左上外挂玩法标签（官方「鱿鱼暴击」） */
.room-flag { position: absolute; left: 0; top: 0; height: 18px; padding: 0 8px 0 6px;
  display: inline-flex; align-items: center; gap: 3px; border-radius: 0 0 8px 0;
  font-size: 10.5px; font-weight: 800; color: #fff;
  background: linear-gradient(180deg, #a06bff, #6f3ad0); }
.room-flag.boom { background: linear-gradient(180deg, #ffa06b, #d8562f); }
.room-flag.fast { background: linear-gradient(180deg, #6bd7ff, #2f8fd8); }
.seat-ring.tone-purple { box-shadow: inset 0 0 0 2px rgba(168, 132, 255, .92), inset 0 0 12px rgba(0,0,0,.5), 0 0 14px rgba(168,132,255,.2); }
.seat-ring.tone-purple::after { border-color: rgba(168, 132, 255, .45); }
.seat-ring.tone-red { box-shadow: inset 0 0 0 2px rgba(255, 140, 150, .9), inset 0 0 12px rgba(0,0,0,.5), 0 0 14px rgba(255,140,150,.18); }
.seat-ring.tone-red::after { border-color: rgba(255, 140, 150, .42); }
.seat-ring.tone-gold { box-shadow: inset 0 0 0 2px rgba(255, 208, 120, .95), inset 0 0 12px rgba(0,0,0,.5), 0 0 14px rgba(255,208,120,.22); }
.seat-ring.tone-gold::after { border-color: rgba(255, 208, 120, .48); }
/* ===== 官方数值（DOM 计算样式实测）===== */
.lobby { padding: 0 15px calc(96px + var(--safe-b)); }
.room {
  position: relative; display: flex; align-items: center; gap: 12px; text-align: left;
  padding: 13px 14px 13px 15px; border-radius: 16px; min-height: 90px;
  background:
    radial-gradient(circle at 88% 50%, rgba(212,175,55,.08) 0%, transparent 58%),
    linear-gradient(135deg, rgba(17,43,32,.93), rgba(10,25,18,.96) 55%, rgba(6,16,12,.98)),
    url(/img/room-card-bg.webp) right center / cover no-repeat;
  border: 1px solid rgba(243,204,122,.24); overflow: hidden;
  box-shadow: 0 6px 18px rgba(0,0,0,.52), inset 0 1px rgba(255,248,220,.14);
  backdrop-filter: blur(12px);
}
.room:active { transform: scale(.994); }
.room-stripe { position: absolute; left: 0; top: 14px; bottom: 14px; width: 3px; border-radius: 0 3px 3px 0;
  background: linear-gradient(180deg,#34d399,#059669); box-shadow: 0 0 8px rgba(52,211,153,.55); }
.room-stripe.pink { background: linear-gradient(180deg,#f472b6,#db2777); box-shadow: 0 0 8px rgba(244,114,182,.55); }
.room-stripe.purple { background: linear-gradient(180deg,#60a5fa,#2563eb); box-shadow: 0 0 8px rgba(96,165,250,.55); }
.room-stripe.orange { background: linear-gradient(180deg,#fbbf24,#d97706); box-shadow: 0 0 8px rgba(251,191,36,.55); }
.seat-ring {
  position: relative; width: 60px; height: 44px; flex: none; border-radius: 999px;
  background: radial-gradient(ellipse at center, #147043, #07331e 66%, #03170d);
  border: 2px solid #a16207;
  box-shadow: 0 0 0 1.2px rgba(253,230,138,.68), 0 4px 10px rgba(0,0,0,.72), inset 0 0 10px rgba(0,0,0,.65);
}
.seat-ring::after { content: ""; position: absolute; inset: 3px; border-radius: 999px; border: 1px solid rgba(253,230,138,.26); }
.ring-dot { position: absolute; width: 8px; height: 8px; border-radius: 50%; margin: -4px 0 0 -4px;
  background: #2b4a3a; border: 1px solid rgba(255,255,255,.22); }
.ring-dot.lit { background: #eafff7; border-color: #fff; box-shadow: 0 0 7px rgba(160,255,225,.95), 0 0 14px rgba(60,240,180,.55); }
.ring-dot.dark { background: #4a3218; border-color: rgba(253,230,138,.35); }
.ring-dot.d1 { left: 24%; top: 16%; } .ring-dot.d2 { left: 76%; top: 16%; }
.ring-dot.d3 { left: 12%; top: 50%; } .ring-dot.d4 { left: 88%; top: 50%; }
.ring-dot.d5 { left: 24%; top: 84%; } .ring-dot.d6 { left: 76%; top: 84%; }
.ring-num { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center;
  font-size: 12.5px; font-weight: 900; letter-spacing: -.4px; color: #fef08a;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  text-shadow: 0 1px 4px rgba(0,0,0,.98), 0 0 8px rgba(250,204,21,.35); }
.room-main { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 1px; }
.room-name { font-size: 13.5px; font-weight: 800; color: #f0f9f4; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.blinds-row { flex-wrap: nowrap; overflow: hidden; gap: 6px !important; margin-top: 0; }
.blinds { font-size: 20px; font-weight: 900; letter-spacing: -.2px; line-height: 1.08; white-space: nowrap;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums;
  background: linear-gradient(180deg, #fffbeb, #fde68a 52%, #d4af37);
  -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent;
  filter: drop-shadow(0 2px 6px rgba(0,0,0,.6)); }
.shield-badge { width: 30px; height: 30px; border-radius: 9px; display: inline-flex; align-items: center; justify-content: center;
  background: linear-gradient(180deg,#3fd98f,#159e5f); border: 1px solid rgba(190,255,220,.55); flex: none; }
.game-pill { font-size: 11px; font-weight: 700; padding: 3px 10px; border-radius: 11px; white-space: nowrap; flex: none; }
.room-meta { display: flex; align-items: center; gap: 5px; margin-top: 0; line-height: 1.2; font-size: 10.5px; color: #82a394; letter-spacing: .02em; }
.room-side { flex: none; display: flex; flex-direction: column; align-items: flex-end; justify-content: center; gap: 2.5px;
  padding-left: 10px; min-width: 81px; border-left: 1px dashed rgba(243,204,122,.24); }
.room-side .tiny { font-size: 9.5px; font-weight: 600; line-height: 1.1; color: #93b0a2; letter-spacing: .02em; }
.buyin { font-size: 20px; font-weight: 900; color: #fff; display: flex; align-items: center; gap: 3.5px; letter-spacing: -.5px;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-variant-numeric: tabular-nums; }
.vpip { font-size: 9.5px; font-weight: 700; letter-spacing: .02em; line-height: 1.1; color: #6ee7b7;
  background: rgba(16,185,129,.14); border: 1px solid rgba(52,211,153,.3); padding: 1.5px 6px; border-radius: 999px; }
.rooms { display: flex; flex-direction: column; gap: 10px; padding-top: 6px; }

.room-flag { position: absolute; left: 0; top: 0; height: 18px; padding: 0 8px 0 6px;
  display: inline-flex; align-items: center; gap: 3px; border-radius: 0 0 8px 0;
  font-size: 10.5px; font-weight: 800; color: #fff;
  background: linear-gradient(180deg,#a06bff,#6f3ad0); z-index: 4; }
.room-flag.boom { background: linear-gradient(180deg,#ffa06b,#d8562f); }
.room-flag.fast { background: linear-gradient(180deg,#6bd7ff,#2f8fd8); }

/* ===== 基础布局（清理时被删，重新补齐）===== */
.marquee { display: flex; align-items: center; gap: 8px; margin: 6px 0 8px; padding: 6px 10px; border-radius: 10px;
  background: linear-gradient(90deg, rgba(10,40,28,.9), rgba(8,28,20,.9)); border: 1px solid rgba(120,220,180,.16); }
.marquee .horn { flex: none; font-size: 13px; }
.marquee-track { flex: 1; overflow: hidden; }
.marquee-track span { display: block; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; font-size: 11.5px; color: #cfe6d9; }

.user-card { display: flex; align-items: center; gap: 11px; padding: 6px 12px; border-radius: 16px; min-height: 73px;
  background: linear-gradient(150deg, rgba(14,42,30,.94), rgba(8,24,17,.96)); border: 1px solid rgba(212,183,104,.28);
  box-shadow: 0 8px 22px rgba(0,0,0,.35), inset 0 1px 0 rgba(255,255,255,.05); }
.ava-wrap { position: relative; flex: none; }
.vip { position: absolute; bottom: -4px; left: 50%; transform: translateX(-50%); font-size: 8px; font-weight: 800;
  letter-spacing: .3px; color: #2b1d05; background: linear-gradient(180deg,#ffe6a8,#d8ab5b); padding: 0 6px;
  border-radius: 6px; border: 1px solid rgba(255,246,205,.7); }
.nick { font-size: 16px; font-weight: 800; letter-spacing: .1px; }
.login-pill { height: 30px; padding: 0 14px; border-radius: 15px; font-size: 13px; font-weight: 800; min-width: 116px;
  display: inline-flex; align-items: center; justify-content: center;
  color: #ffe6a8; background: rgba(20,50,34,.7); border: 1px solid rgba(239,207,142,.75); }
.bal-row { display: flex; align-items: center; gap: 4px; margin-top: 2px; min-width: 0; }
.bal-pill { height: 30px; display: inline-flex; align-items: center; gap: 4px; padding: 0 12px 0 4px; border-radius: 15px;
  font-size: 19px; font-weight: 800; color: #fff; background: rgba(10,34,24,.9); border: 1px solid rgba(120,220,180,.24); }
.bal-pill .coin { width: 19px; height: 19px; font-size: 11.5px; }
.eye-btn { width: 24px; height: 24px; border-radius: 50%; background: rgba(8,40,30,.9);
  border: 1px solid rgba(120,220,180,.3); display: flex; align-items: center; justify-content: center; color: #bff0d8; }
.plus-btn { width: 26px; height: 26px; border-radius: 50%; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(180deg,#ffe08a,#d8a63c); color: #2a1a02; border: 1px solid rgba(255,246,205,.7);
  box-shadow: 0 4px 10px rgba(216,166,60,.35); }
.top-tools { display: flex; align-items: center; gap: 6px; flex: none; }
.tool { width: 52px; height: 50px; border-radius: 13px; display: flex; flex-direction: column; align-items: center;
  justify-content: center; gap: 3px; font-size: 9.5px; font-weight: 700; color: #d9ecdf;
  background: rgba(10,34,24,.8); border: 1px solid rgba(157,194,192,.18); }

.filters { display: flex; align-items: center; gap: 6px; margin-bottom: 10px; }
.filters .tabs { flex: 1 1 auto; min-width: 0; display: flex; gap: 2px; overflow-x: auto; scrollbar-width: none; }
.filters .tabs::-webkit-scrollbar { display: none; }
.tab { height: 35px; font-size: 13px; font-weight: 700; padding: 0 10px; border-radius: 12px; white-space: nowrap;
  color: #d9ecdf; background: transparent; border: 1px solid transparent; flex: none; }
.tab.active { color: #23180a; font-weight: 900; background: linear-gradient(180deg,#ffe9a8,#e7bf5c); border-color: transparent; }
.key-btn { width: 34px; height: 35px; border-radius: 12px; flex: none; display: flex; flex-direction: column;
  align-items: center; justify-content: center; gap: 1px; font-size: 8.5px; color: #efcf8e;
  background: rgba(10,34,24,.85); border: 1px solid rgba(239,207,142,.28); }
.open-btn { height: 35px; padding: 0 12px; border-radius: 12px; flex: none; display: inline-flex; align-items: center; gap: 4px;
  font-size: 13.5px; font-weight: 900; color: #04230f; background: linear-gradient(180deg,#7ef0b4,#2bbd7a);
  box-shadow: 0 6px 16px rgba(35,190,125,.28); }

/* 官方结构之外的补充：房间卡小角标 / 奖池弹层 / 空态 */
.lobby-page {
  --header-safe-pt: 12px;
  font-size: 16px;
  line-height: normal;
}
.lobby-page .club-name-row { min-width: 46px; }
.lobby-page .club-name-text { line-height: 17px; min-width: 44px !important; }
.lobby-page .club-identity-row { gap: 4px; }
.lobby-page .club-id-row { flex-shrink: 1 !important; }
.lobby-page .club-id-row { min-height: 17px; }
.lobby-page .room-title-line, .lobby-page .room-title { line-height: 19px; }
.lobby-page .room-blinds-line { min-height: 24px; }
.lobby-page .jp-plaque-digits { line-height: 14px; }
.lobby-page .club-id-row { min-width: 0; overflow: hidden; }
.lobby-page .club-id-text { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.room-flag-mini { flex: none; font-size: 9.5px; font-weight: 800; color: #2b1400; padding: 1px 6px; border-radius: 6px;
  background: linear-gradient(180deg, #ffe08a, #e0a83a); }
.lobby-empty-box { display: flex; flex-direction: column; align-items: center; gap: 6px; padding: 40px 16px; text-align: center;
  border-radius: 18px; background: rgba(9, 28, 19, .7); border: 1px dashed rgba(120, 200, 160, .22); }
.lobby-empty-box b { font-size: 14px; color: #e6f4ee; }
.lobby-empty-box span { font-size: 12px; color: var(--muted); }
.jackpot-hero { display: flex; flex-direction: column; align-items: center; gap: 4px; padding: 16px; border-radius: 16px;
  background: radial-gradient(circle at 50% 0%, rgba(212,175,55,.18), rgba(8,26,17,.9)); border: 1px solid rgba(214,176,92,.3); }
.jackpot-amount { font-size: 28px; font-weight: 900; color: #ffd875; }
.jackpot-tiers { display: flex; flex-direction: column; gap: 8px; margin-top: 12px; }
.jackpot-tier { display: flex; align-items: center; justify-content: space-between; padding: 10px 12px; border-radius: 12px;
  background: rgba(10, 32, 22, .8); border: 1px solid var(--line); font-size: 13px; color: #d8ecdf; }
.jackpot-tier b { color: #efcf8e; }
.mtt-enter-btn.signed { opacity: .7; }
</style>
