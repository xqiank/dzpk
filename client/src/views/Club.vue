<template>
  <div class="club-page">
  <div class="club-center-container">
    <!-- ============ 俱乐部大厅（进入俱乐部后） ============ -->
    <div v-if="hubClub" class="club-hub-view">
      <div class="club-hub-fixed-header">
        <div class="hub-nav-bar">
          <button class="hub-back-btn" @click="hubClub = null">
            <Icon name="back" :size="18" />
          </button>
          <div class="hub-nav-title">{{ hubClub.name }}</div>
          <div class="hub-right-tools">
            <button class="hub-members-icon-btn" title="成员" @click="openMembers(hubClub)"><Icon name="users" :size="16" /></button>
            <button class="hub-tool-btn" title="公告" @click="openNotice(hubClub)"><Icon name="bell" :size="18" /></button>
          </div>
        </div>

        <div class="hub-hero-card">
          <div class="hero-left-info">
            <div class="hero-avatar-wrap">
              <img :src="hubClub.avatar || '/avatars/av5.svg'" alt="俱乐部头像" />
            </div>
            <div class="hero-meta-wrap">
              <div class="hero-name-row">
                <Icon name="crown" :size="15" color="#ffd700" />
                <span class="hero-club-name">{{ hubClub.name }}</span>
                <span v-if="hubClub.certified" class="club-role-badge badge-public">认证</span>
              </div>
              <div class="hero-id-row" @click="copyClubNo(hubClub.clubNo)">
                <span class="hero-id-text">ID: {{ hubClub.clubNo }}</span>
                <Icon name="copy" :size="12" />
              </div>
            </div>
          </div>
          <div class="hero-right-buttons">
            <button class="sq-icon" title="房卡池" @click="openRecharge">
              <Icon name="card" :size="18" color="#efcf8e" />
              <span>{{ hubClub.cardPool }}</span>
            </button>
            <button class="sq-icon" title="卡池明细" @click="openCardLogs(hubClub)">
              <Icon name="history" :size="18" color="#9fe0c0" />
            </button>
            <button class="sq-icon" title="成员列表" @click="openMembers(hubClub)">
              <Icon name="users" :size="18" color="#efcf8e" />
              <span>{{ hubClub.members }}</span>
            </button>
          </div>
        </div>

        <div class="hub-tabs-row">
          <div class="game-type-tabs">
            <button
              v-for="t in hubTabs"
              :key="t.key"
              class="hub-tab"
              :class="{ active: hubTab === t.key }"
              @click="hubTab = t.key"
            >{{ t.label }}</button>
          </div>
          <div class="hub-action-tools">
            <button v-if="hubClub.owner" class="hub-tool primary" title="开一张俱乐部桌" @click="openRoomForm">
              <Icon name="plus" :size="16" color="#04220f" />
            </button>
            <button v-if="hubClub.owner" class="hub-tool" title="入会申请" @click="openApplications(hubClub)">
              <Icon name="doc" :size="16" color="#efcf8e" />
            </button>
          </div>
        </div>
      </div>

      <main class="club-scrollable-body">
        <div class="hub-tables-list">
          <article
            v-for="r in hubRooms"
            :key="r.roomNo"
            class="poker-table-card"
            :class="[gameModeClass(r), occClass(r)]"
            @click="enter(r)"
          >
            <div class="card-occupancy-bar" :class="occClass(r)" aria-hidden="true">
              <span class="card-occupancy-fill" :style="{ height: Math.round((r.online / r.seats) * 100) + '%' }"></span>
            </div>
            <div class="card-watermark-overlay"></div>
            <div class="card-tether-watermark" aria-hidden="true">
              <svg viewBox="0 0 100 100" fill="currentColor"><path d="M50 0C22.4 0 0 22.4 0 50s22.4 50 50 50 50-22.4 50-50S77.6 0 50 0zm21.6 28.5H54.2v6.2h17.4v6.2H54.2v5.3c15.2.8 26.6 4.3 26.6 8.5 0 4.2-11.4 7.7-26.6 8.5v16.1H45.8V63.2c-15.2-.8-26.6-4.3-26.6-8.5 0-4.2 11.4-7.7 26.6-8.5v-5.3H28.4v-6.2h17.4v-6.2H28.4v-6.2h43.2v6.2zm-21.6 30c12.3 0 22.5-2.2 24.8-5.3-2.3-3.1-12.5-5.3-24.8-5.3s-22.5 2.2-24.8 5.3c2.3 3.1 12.5 5.3 24.8 5.3z"></path></svg>
            </div>
            <div class="card-donut-col">
              <div class="table-topdown-box">
                <svg class="seat-arc-svg" viewBox="-6 -6 67 51" aria-hidden="true">
                  <path class="seat-arc-glow" d="M 16.26 -0.49 A 20.25 20.25 0 0 1 19.50 -0.75 L 35.50 -0.75 A 20.25 20.25 0 0 1 38.74 39.49"></path>
                  <path class="seat-arc-line" d="M 16.26 -0.49 A 20.25 20.25 0 0 1 19.50 -0.75 L 35.50 -0.75 A 20.25 20.25 0 0 1 38.74 39.49" stroke="url(#clubSeatArcGrad)"></path>
                </svg>
                <span
                  v-for="(d, i) in seatDots(r)"
                  :key="i"
                  class="s-dot"
                  :class="{ on: i < r.online }"
                  :style="{ left: d.x + 'px', top: d.y + 'px' }"
                ></span>
                <b class="table-seat-num">{{ r.online }}/{{ r.seats }}</b>
              </div>
            </div>
            <div class="card-info-col">
              <div class="room-title-line"><span class="room-title">{{ r.name }}</span></div>
              <div class="room-blinds-line">
                <strong class="blinds-num">{{ trim(r.sb) }}/{{ trim(r.bb) }}</strong>
                <span class="room-safety-badge" aria-hidden="true">
                  <svg viewBox="0 0 24 26" width="16" height="18"><path d="M12 1 2.5 4.7v7.1c0 5.8 3.8 10.9 9.5 13.2 5.7-2.3 9.5-7.4 9.5-13.2V4.7L12 1Zm0 3.1 6.8 2.6v5.1c0 4.4-2.6 8.5-6.8 10.4-4.2-1.9-6.8-6-6.8-10.4V6.7L12 4.1Z" fill="currentColor" fill-rule="evenodd"></path><path d="M12 4.1v8H5.2V6.7L12 4.1Zm0 8h6.8c-.1 4.3-2.7 8.3-6.8 10.2V12.1Z" fill="currentColor"></path></svg>
                </span>
                <span class="room-badge" :class="'room-badge--' + r.gameType"><span>{{ gameName(r.gameType) }}</span></span>
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
                <span class="usdt-green-icon">₮</span><b class="buyin-value">{{ trim(r.minBuyin) }}</b>
              </div>
              <span class="vpip-text">每手 {{ r.cardCost || 1 }} 张房卡</span>
            </div>
          </article>

          <div v-if="!hubRooms.length" class="empty-hub-rooms">
            <p>当前俱乐部暂无开放对局</p>
            <button class="empty-hub-create-btn" :disabled="!hubClub.owner" @click="openRoomForm">立即开一张俱乐部桌</button>
          </div>
        </div>
      </main>
    </div>

    <!-- ============ 我的俱乐部（目录） ============ -->
    <div v-else class="club-directory-view">
      <div class="club-directory-fixed-header">
        <div class="club-actions-standalone-row">
          <button class="club-action-standalone-join" title="输入俱乐部ID申请加入俱乐部" @click="joinOpen = true">
            <Icon name="plus" :size="16" />
            <span>加入俱乐部</span>
          </button>
          <button class="club-action-standalone-create" title="创建属于自己的俱乐部" @click="createOpen = true">
            <Icon name="crown" :size="16" color="#04220f" />
            <span>创建俱乐部</span>
          </button>
        </div>
      </div>

      <main class="club-scrollable-body">
        <div class="club-cards-stack">
          <article
            v-for="c in clubs"
            :key="c.id"
            class="club-item-card"
            :class="{ 'has-corner-medal': true, 'is-owner': c.owner, 'is-joined': !c.owner }"
            @click="openHub(c)"
          >
            <div class="club-corner-medal" :class="c.owner ? 'medal-owner' : 'medal-joined'" :title="c.owner ? '我创建的俱乐部' : '已加入该俱乐部'">
              <svg class="medal-icon" width="11" height="11" viewBox="0 0 24 24" fill="currentColor"><path d="M5 19h14v1.5H5V19zm1-7.5 3 2.5 3-6 3 6 3-2.5-1.5 6H7.5L6 11.5z"></path></svg>
              <span>{{ c.owner ? '创建' : '已加入' }}</span>
            </div>
            <div class="club-avatar-column" @click.stop="openMembers(c)">
              <div class="club-avatar-frame">
                <img :src="c.avatar || '/avatars/av5.svg'" alt="俱乐部头像" class="club-img" />
                <span class="club-chat-badge" title="成员列表">群聊</span>
              </div>
            </div>
            <div class="club-meta-column">
              <div class="club-meta-name-row">
                <svg class="crown-svg" width="15" height="15" viewBox="0 0 24 24" fill="currentColor"><path d="M5 19h14v1.5H5V19zm1-7.5 3 2.5 3-6 3 6 3-2.5-1.5 6H7.5L6 11.5z"></path></svg>
                <strong class="club-name-text">{{ c.name }}</strong>
                <span v-if="c.owner" class="club-role-badge badge-admin">会长</span>
                <span v-else class="club-role-badge badge-public">成员</span>
              </div>
              <div class="club-meta-id-row">
                <span class="club-id-tag">ID: {{ c.clubNo }}</span>
                <button class="club-copy-tag" title="复制ID" @click.stop="copyClubNo(c.clubNo)">
                  <svg class="copy-svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
                </button>
                <span class="club-member-count">
                  <Icon name="users" :size="12" /> {{ c.members }} 成员
                </span>
                <span class="club-active-tables-tag"><span class="pulse-live-dot"></span> 牌局进行中</span>
              </div>
              <p class="club-announcement-snippet">{{ c.notice || 'U5德州专属俱乐部，区块链公平存证，保障每一局公平竞技。' }}</p>
            </div>
            <div class="club-action-column">
              <button class="club-chat-quick-btn" title="成员列表" @click.stop="openMembers(c)">
                <Icon name="users" :size="13" /><span>成员</span>
              </button>
              <button class="club-enter-hub-btn" title="进入俱乐部大厅" @click.stop="openHub(c)">进入</button>
            </div>
          </article>

          <div v-if="!clubs.length" class="empty-clubs-state">
            <div class="empty-icon-shield">
              <div class="shield-halo-box">
                <Icon name="crown" :size="34" color="#ffd700" />
              </div>
            </div>
            <b>暂无俱乐部</b>
            <p>点击上方「创建俱乐部」，开启您的专属德州竞技圈，或点击「加入俱乐部」输入 ID 申请加入</p>
          </div>
        </div>
      </main>
    </div>

    <div v-if="roomForm" class="mask" @click.self="roomForm = null">
      <div class="sheet">
        <div class="row-between" style="margin-bottom:12px">
          <b class="h2">俱乐部开局</b>
          <button class="close-x" @click="roomForm = null"><Icon name="close" :size="18" /></button>
        </div>
        <div class="segmented dark" style="margin-bottom:10px">
          <button v-for="g in gameTypes" :key="g.key" class="seg" :class="{ active: roomForm.gameType === g.key }"
                  @click="roomForm.gameType = g.key">{{ g.label }}</button>
        </div>
        <div class="grid2">
          <label class="fld"><span>人数</span>
            <select v-model.number="roomForm.seatCount" class="field">
              <option v-for="n in [2,3,4,5,6,7,8,9]" :key="n" :value="n">{{ n }} 人</option>
            </select>
          </label>
          <label class="fld"><span>时长</span>
            <select v-model.number="roomForm.duration" class="field">
              <option :value="60">60 分钟</option><option :value="90">90 分钟</option>
              <option :value="120">120 分钟</option><option :value="240">240 分钟</option>
            </select>
          </label>
          <label class="fld"><span>小盲</span>
            <select v-model.number="roomForm.sb" class="field" @change="roomForm.bb = roomForm.sb * 2">
              <option :value="10">0.10</option><option :value="20">0.20</option>
              <option :value="50">0.50</option><option :value="100">1.00</option>
            </select>
          </label>
          <label class="fld"><span>大盲</span>
            <select v-model.number="roomForm.bb" class="field">
              <option :value="20">0.20</option><option :value="40">0.40</option>
              <option :value="100">1.00</option><option :value="200">2.00</option>
            </select>
          </label>
          <label class="fld"><span>最小买入</span>
            <select v-model.number="roomForm.minBuyin" class="field">
              <option :value="1000">10.00</option><option :value="2000">20.00</option>
              <option :value="5000">50.00</option><option :value="10000">100.00</option>
            </select>
          </label>
          <label class="fld"><span>最大买入</span>
            <select v-model.number="roomForm.maxBuyin" class="field">
              <option :value="8000">80.00</option><option :value="20000">200.00</option>
              <option :value="40000">400.00</option>
            </select>
          </label>
        </div>
        <div class="cost-line" style="margin-top:12px">
          <Icon name="card" :size="14" color="#efcf8e" />
          <span>俱乐部卡池 {{ activeClub?.cardPool ?? 0 }} 张 · 建房扣 1 张，每手再扣 1 张</span>
        </div>
        <button class="btn btn-primary btn-block" style="height:48px;margin-top:14px" @click="createClubRoom">创建并进入</button>
      </div>
    </div>

    <div v-if="rechargeOpen" class="club-dialog-mask" @click.self="rechargeOpen = false">
      <div class="club-dialog-card">
        <div class="dialog-header">
          <span class="modal-dialog-title"><Icon name="card" :size="18" class="modal-title-svg" />转入房卡到俱乐部</span>
          <button class="dialog-close-btn" @click="rechargeOpen = false">×</button>
        </div>
        <div class="dialog-body">
          <p class="dialog-tip">从我的房卡转给俱乐部卡池，所有成员开局共用。</p>
          <div class="form-input-group">
            <div class="form-label-with-hint"><label>转入数量</label><span class="input-limit-hint">我的房卡 {{ auth.profile?.cards ?? 0 }} 张</span></div>
            <div class="row gap-8">
              <button v-for="n in [1,5,10,20]" :key="n" class="chip" :class="{ on: rechargeCount === n }" @click="rechargeCount = n">{{ n }} 张</button>
            </div>
          </div>
          <button class="dialog-submit-btn" @click="doRecharge">确认转入</button>
        </div>
      </div>
    </div>

    <div v-if="joinOpen" class="club-dialog-mask" @click.self="joinOpen = false">
      <div class="club-dialog-card">
        <div class="dialog-header">
          <span class="modal-dialog-title"><Icon name="users" :size="18" class="modal-title-svg" />加入俱乐部</span>
          <button class="dialog-close-btn" @click="joinOpen = false">×</button>
        </div>
        <div class="dialog-body">
          <p class="dialog-tip">输入俱乐部 ID 提交申请，会长审核后即可进入。</p>
          <div class="form-input-group">
            <label>俱乐部 ID</label>
            <input v-model="clubNo" class="field" placeholder="输入俱乐部 ID（8 位数字）" maxlength="8" />
          </div>
          <button class="dialog-submit-btn" @click="joinClub">申请加入</button>
        </div>
      </div>
    </div>

    <div v-if="createOpen" class="club-dialog-mask" @click.self="createOpen = false">
      <div class="club-dialog-card">
        <div class="dialog-header">
          <span class="modal-dialog-title"><Icon name="crown" :size="18" class="modal-title-svg" />创建俱乐部</span>
          <button class="dialog-close-btn" @click="createOpen = false">×</button>
        </div>
        <div class="dialog-body">
          <p class="dialog-tip">创建后你就是会长，可邀请成员、发起俱乐部牌桌。</p>
          <div class="form-input-group">
            <div class="form-label-with-hint"><label>俱乐部名称</label><span class="input-limit-hint" :class="{ error: newName.length > 16 }">{{ newName.length }}/16</span></div>
            <input v-model="newName" class="field" placeholder="俱乐部名称" maxlength="16" />
          </div>
          <div class="form-input-group">
            <div class="form-label-with-hint"><label>俱乐部公告</label><span class="input-limit-hint">可选</span></div>
            <input v-model="newNotice" class="field" placeholder="例如：每晚 20:00 开桌" maxlength="80" />
          </div>
          <button class="dialog-submit-btn" :disabled="!newName.trim()" @click="createClub">创建</button>
        </div>
      </div>
    </div>

    <div v-if="cardLogOpen" class="club-dialog-mask" @click.self="cardLogOpen = false">
      <div class="club-dialog-card">
        <div class="dialog-header">
          <span class="modal-dialog-title"><Icon name="card" :size="18" class="modal-title-svg" />俱乐部卡池明细</span>
          <button class="dialog-close-btn" @click="cardLogOpen = false">×</button>
        </div>
        <div class="dialog-body">
        <div v-if="!cardLogs.length" class="tiny muted" style="padding:16px 0;text-align:center">暂无卡池流水</div>
        <div v-for="l in cardLogs" :key="l.id" class="cardlog-row">
          <span class="cl-badge" :class="l.delta > 0 ? 'in' : 'out'">{{ l.delta > 0 ? '+' : '' }}{{ l.delta }}</span>
          <div class="col" style="flex:1;min-width:0">
            <b class="cl-reason">{{ reasonName(l.reason) }}</b>
            <span class="tiny muted ellipsis">{{ l.nick }} · {{ l.remark || '—' }}</span>
          </div>
          <div class="col" style="align-items:flex-end">
            <b class="tiny">余 {{ l.balance }} 张</b>
            <span class="tiny muted">{{ cardAgo(l.time) }}</span>
          </div>
        </div>
        </div>
      </div>
    </div>

    <div v-if="noticeOpen" class="club-dialog-mask" @click.self="noticeOpen = false">
      <div class="club-dialog-card">
        <div class="dialog-header">
          <span class="modal-dialog-title"><Icon name="bell" :size="18" class="modal-title-svg" />俱乐部公告</span>
          <button class="dialog-close-btn" @click="noticeOpen = false">×</button>
        </div>
        <div class="dialog-body">
          <p class="dialog-tip">公告会推送给所有在线成员，最多 80 字。</p>
          <div class="form-input-group">
            <div class="form-label-with-hint"><label>公告内容</label><span class="input-limit-hint" :class="{ error: noticeDraft.length > 80 }">{{ noticeDraft.length }}/80</span></div>
            <textarea v-model="noticeDraft" class="field notice-input" maxlength="80" placeholder="例如：每晚 20:00 开桌 · 新手友好"></textarea>
          </div>
          <button class="dialog-submit-btn" :disabled="noticeBusy" @click="saveNotice">{{ noticeBusy ? '保存中…' : '发布公告' }}</button>
        </div>
      </div>
    </div>

    <!-- 成员 / 入会申请（官方 club-members-fullscreen） -->
    <div v-if="membersOpen" class="club-members-fullscreen">
      <header class="members-header">
        <div class="mh-left">
          <button class="hub-back-btn" @click="membersOpen = false">‹ 返回</button>
          <b class="mh-title">俱乐部成员</b>
          <span class="club-member-count">共 {{ members.length }} 人</span>
        </div>
        <button class="hub-members-icon-btn" @click="loadMembers" title="刷新"><Icon name="refresh" :size="15" /></button>
      </header>

      <div class="members-sub-tabs">
        <button :class="{ active: memberTab === 'all' }" @click="memberTab = 'all'">全部成员</button>
        <button :class="{ active: memberTab === 'apply' }" @click="memberTab = 'apply'">
          入会申请{{ isOwner && applications.length ? ` (${applications.length})` : '' }}
        </button>
      </div>

      <div class="members-content-scroll">
        <template v-if="memberTab === 'all'">
          <div v-for="m in members" :key="m.uid" class="member-card-row">
            <div class="member-avatar-box"><Avatar :src="m.avatar" :size="44" class="m-avatar" /></div>
            <div class="member-info-col">
              <span class="m-name">{{ m.nick }}</span>
              <span class="m-id">ID: {{ m.uid }}</span>
            </div>
            <span class="member-role-tag" :class="m.role === 2 ? 'owner' : m.role === 1 ? 'admin' : ''">
              {{ m.role === 2 ? '创建者' : m.role === 1 ? '管理员' : '成员' }}
            </span>
            <div v-if="isOwner && m.role < 2 && m.uid !== auth.profile?.uid" class="member-actions">
              <button class="btn btn-sm btn-ghost" @click="setRole(m, m.role === 1 ? 0 : 1)">{{ m.role === 1 ? '取消管理' : '设为管理' }}</button>
              <button class="btn btn-sm btn-ghost" @click="transfer(m)">转让</button>
              <button class="btn btn-sm btn-danger" @click="kick(m)">移出</button>
            </div>
            <div v-else-if="m.uid !== auth.profile?.uid" class="member-actions">
              <button class="btn btn-sm btn-ghost" @click="kick(m)">移出</button>
            </div>
          </div>
          <div v-if="!members.length" class="members-empty">暂无成员</div>
        </template>

        <template v-else>
          <div v-for="a in applications" :key="'app' + a.uid" class="app-card-row">
            <div class="member-avatar-box"><Avatar :src="a.avatar" :size="44" class="m-avatar" /></div>
            <div class="member-info-col">
              <span class="m-name">{{ a.nick }}</span>
              <span class="m-id">ID: {{ a.uid }} · 申请加入</span>
            </div>
            <div class="app-action-btns">
              <button class="btn-app-approve" @click="review(a, true)">同意</button>
              <button class="btn-app-reject" @click="review(a, false)">拒绝</button>
            </div>
          </div>
          <div v-if="!applications.length" class="empty-app-text">
            {{ isOwner ? '暂无待处理的入会申请' : '只有会长可以审核入会申请' }}
          </div>
        </template>
      </div>
    </div>
  </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import Avatar from '@/components/Avatar.vue'
import Icon from '@/components/Icon.vue'
import { socket } from '@/lib/ws'
import { useUI } from '@/stores/ui'
import { useTable } from '@/stores/table'
import { useAuth } from '@/stores/auth'
import { sfx } from '@/lib/sfx'
import { remain } from '@/lib/format'
import '@/styles/club-official.css'
import '@/styles/lobby-official.css'

const router = useRouter()
const ui = useUI()
const table = useTable()
const auth = useAuth()

const clubs = ref<any[]>([])
const rooms = ref<any[]>([])
const joinOpen = ref(false)
const createOpen = ref(false)
const membersOpen = ref(false)
const memberTab = ref<'all' | 'apply'>('all')
const members = ref<any[]>([])
const applications = ref<any[]>([])
const appliedClubId = ref(0)
const membersClub = ref<any>(null)
const isOwner = computed(() => !!(membersClub.value?.owner))
const roomForm = ref<any>(null)
const rechargeOpen = ref(false)
const rechargeCount = ref(5)
const gameTypes = [
  { key: 'nlh', label: '长牌' }, { key: 'short6', label: '短牌' },
  { key: 'omaha', label: '奥马哈' }, { key: 'niuniu', label: '牛牛' }
]
const activeClub = computed(() => clubs.value[0] || null)
const clubNo = ref('')
const newName = ref('')
const newNotice = ref('')

async function load() {
  try {
    const d = await socket.call<{ clubs: any[]; rooms: any[] }>('club.list')
    clubs.value = d.clubs
    rooms.value = d.rooms
    d.clubs.forEach((c: any) => {
      if (!c.owner) return
      certifyMap.value = { ...certifyMap.value, [c.id]: { certified: !!c.certified, pending: !!c.certifyRequest } }
      loadCertify(c)
    })
  } catch (e: any) {
    ui.err(e.message)
  }
}

const hubClub = ref<any>(null)
const hubTab = ref<'all' | 'nlh' | 'short6' | 'omaha' | 'niuniu'>('all')
const hubTabs = [
  { key: 'all', label: '全部' },
  { key: 'nlh', label: '长牌' },
  { key: 'short6', label: '短牌' },
  { key: 'omaha', label: '奥马哈' },
  { key: 'niuniu', label: '牛牛' }
] as const
const hubRooms = computed(() => {
  if (!hubClub.value) return []
  const list = rooms.value.filter((r: any) => r.clubId === hubClub.value.id)
  const filtered = hubTab.value === 'all' ? list : list.filter((r: any) => r.gameType === hubTab.value)
  return [...filtered].sort((a: any, b: any) => (b.online || 0) - (a.online || 0) || (a.roomNo || 0) - (b.roomNo || 0))
})
const gameName = (g: string) => ({ nlh: '长牌', short6: '短牌', omaha: '奥马哈', niuniu: '牛牛' } as any)[g] || '长牌'
const OFFICIAL_SEAT_DOTS = [
  { x: 16.26, y: -0.49 }, { x: 38.74, y: -0.49 }, { x: 55.75, y: 19.5 },
  { x: 38.74, y: 39.49 }, { x: 16.26, y: 39.49 }, { x: -0.75, y: 19.5 }
]
function seatDots(r: any) {
  const n = Math.max(2, Math.min(9, Number(r?.seats) || 6))
  if (n === 6) return OFFICIAL_SEAT_DOTS
  const cx = 27.5, cy = 19.5, rx = 28.25, ry = 20
  return Array.from({ length: n }, (_, i) => {
    const a = -Math.PI / 2 + (i * Math.PI * 2) / n + Math.PI / n
    return { x: +(cx + rx * Math.cos(a)).toFixed(2), y: +(cy + ry * Math.sin(a)).toFixed(2) }
  })
}
function occClass(r: any) {
  const p = (Number(r?.online) || 0) / Math.max(1, Number(r?.seats) || 6)
  if (p >= 0.8) return 'occ-hot'
  if (p >= 0.4) return 'occ-mid'
  return 'occ-low'
}
function gameModeClass(r: any) {
  return 'mode-' + (r?.gameType || 'nlh')
}
function openHub(c: any) {
  hubClub.value = c
  hubTab.value = 'all'
  sfx.tap()
}
function copyClubNo(no: any) {
  navigator.clipboard?.writeText(String(no || '')).then(() => ui.ok('俱乐部 ID 已复制')).catch(() => ui.err('复制失败'))
  sfx.tap()
}
async function openApplications(c: any) {
  membersClub.value = c
  membersOpen.value = true
  try {
    applications.value = await socket.call('club.applications', { clubId: c.id })
  } catch { applications.value = [] }
}

const roomSort = ref<'online' | 'blind' | 'new'>('online')
const noticeOfClub = computed(() => (clubs.value[0]?.notice || '').trim())
const roomSorts = [
  { key: 'online', label: '人气' },
  { key: 'blind', label: '盲注' },
  { key: 'new', label: '最新' }
] as const
function trim(cents: number) {
  const v = Number(cents || 0) / 100
  const n = Math.round(v * 100) / 100
  return Number.isInteger(n) ? String(n) : n.toFixed(2)
}
const sortedRooms = computed(() => {
  const list = [...rooms.value]
  if (roomSort.value === 'blind') return list.sort((a, b) => (a.bb || 0) - (b.bb || 0))
  if (roomSort.value === 'new') return list.sort((a, b) => (b.roomNo || 0) - (a.roomNo || 0))
  return list.sort((a, b) => (b.online || 0) - (a.online || 0) || (a.roomNo || 0) - (b.roomNo || 0))
})
const cardLogOpen = ref(false)
const cardLogs = ref<any[]>([])
const reasonName = (r: string) => ({ recharge: '成员转入', room: '创建房间', hand: '牌局消耗', adjust: '后台调整' } as any)[r] || r
function cardAgo(ms: number) {
  const d = Math.max(0, Date.now() - (ms || 0))
  if (d < 60e3) return '刚刚'
  if (d < 3600e3) return Math.floor(d / 60e3) + ' 分钟前'
  if (d < 86400e3) return Math.floor(d / 3600e3) + ' 小时前'
  const t = new Date(ms)
  return `${t.getMonth() + 1}/${t.getDate()}`
}
async function openCardLogs(c: any) {
  cardLogOpen.value = true
  cardLogs.value = []
  try {
    const d = await socket.call<{ list: any[] }>('club.cardLogs', { clubId: c.id })
    cardLogs.value = d.list || []
  } catch (e: any) {
    ui.err(e.message)
  }
}
const noticeOpen = ref(false)
const noticeDraft = ref('')
const noticeBusy = ref(false)
const noticeClub = ref<any>(null)
function openNotice(c: any) {
  noticeClub.value = c
  noticeDraft.value = c.notice || ''
  noticeOpen.value = true
}
async function saveNotice() {
  if (!noticeClub.value || noticeBusy.value) return
  noticeBusy.value = true
  try {
    const r = await socket.call<{ notice: string }>('club.updateNotice', { clubId: noticeClub.value.id, notice: noticeDraft.value })
    noticeClub.value.notice = r.notice
    noticeOpen.value = false
    ui.ok('公告已推送给在线成员')
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    noticeBusy.value = false
  }
}
const certifyMap = ref<Record<number, { certified: boolean; pending: boolean }>>({})
function pendingOf(c: any) { return !!certifyMap.value[c.id]?.pending }
async function loadCertify(c: any) {
  if (!c?.owner) return
  try {
    const d = await socket.call<{ certified: boolean; pending: boolean }>('club.officialStatus', { clubId: c.id })
    certifyMap.value = { ...certifyMap.value, [c.id]: { certified: !!d.certified, pending: !!d.pending } }
    if (d.certified) c.certified = 1
  } catch {}
}
async function requestCertify(c: any) {
  if (pendingOf(c) || c.certified) return
  try {
    const d = await socket.call<{ pending: boolean; certified: boolean }>('club.requestCertify', { clubId: c.id })
    certifyMap.value = { ...certifyMap.value, [c.id]: { certified: !!d.certified, pending: !!d.pending } }
    ui.ok(d.pending ? '认证申请已提交，等待官方审核' : '俱乐部已通过认证')
    if (d.certified) c.certified = 1
  } catch (e: any) { ui.err(e.message) }
}

async function createClub() {
  if (!newName.value.trim()) return ui.err('请填写名称')
  try {
    await socket.call('club.create', { name: newName.value, notice: newNotice.value, avatar: '/avatars/av7.svg' })
    ui.ok('俱乐部已创建')
    createOpen.value = false
    load()
    sfx.win()
  } catch (e: any) {
    ui.err(e.message)
  }
}

async function joinClub() {
  if (!clubNo.value.trim()) return ui.err('请输入俱乐部 ID')
  try {
    const r = await socket.call('club.join', { clubNo: clubNo.value.trim() })
    ui.ok(r?.joined ? '已加入俱乐部' : '申请已提交，等待会长审批')
    if (!r?.joined) appliedClubId.value = r?.id || 0
    joinOpen.value = false
    load()
  } catch (e: any) {
    ui.err(e.message)
  }
}

async function kick(m: any) {
  const cl = clubs.value[0]
  if (!cl) return
  try {
    await socket.call('club.kick', { clubId: cl.id, uid: m.uid })
    ui.ok('已移出成员')
    members.value = members.value.filter((x: any) => x.uid !== m.uid)
    load()
  } catch (e: any) { ui.err(e.message) }
}

async function review(a: any, approve: boolean) {
  if (!membersClub.value) return
  try {
    await socket.call('club.review', { clubId: membersClub.value.id, uid: a.uid, approve })
    ui.ok(approve ? '已同意入会' : '已拒绝申请')
    applications.value = applications.value.filter(x => x.uid !== a.uid)
    if (approve) {
      members.value = await socket.call('club.members', { clubId: membersClub.value.id })
      load()
    }
  } catch (e: any) {
    ui.err(e.message)
  }
}
async function setRole(m: any, role: number) {
  if (!membersClub.value) return
  try {
    await socket.call('club.setRole', { clubId: membersClub.value.id, uid: m.uid, role })
    m.role = role
    ui.ok(role === 1 ? m.nick + ' 已成为管理员' : m.nick + ' 已取消管理员')
  } catch (e: any) {
    ui.err(e.message)
  }
}
async function transfer(m: any) {
  if (!membersClub.value) return
  try {
    await socket.call('club.transfer', { clubId: membersClub.value.id, uid: m.uid })
    ui.ok('会长已转让给 ' + m.nick)
    membersOpen.value = false
    load()
  } catch (e: any) {
    ui.err(e.message)
  }
}
async function loadMembers() {
  const c = membersClub.value
  if (!c) return
  try {
    members.value = await socket.call('club.members', { clubId: c.id })
    if (c.owner) {
      try {
        const r: any = await socket.call('club.applications', { clubId: c.id })
        applications.value = r.list || []
      } catch { /* 非会长忽略 */ }
    }
  } catch (e: any) {
    ui.err(e.message)
  }
}
async function openMembers(c: any) {
  try {
    membersClub.value = c
    members.value = await socket.call('club.members', { clubId: c.id })
    applications.value = []
    if (c.owner) {
      try {
        const r = await socket.call('club.applications', { clubId: c.id })
        applications.value = r.list || []
      } catch { /* 非会长忽略 */ }
    }
    membersOpen.value = true
  } catch (e: any) {
    ui.err(e.message)
  }
}

async function enter(r: any) {
  try {
    await table.join(r.roomNo, 0, -1, true)
    router.push(`/table/${r.roomNo}?observe=1`)
  } catch (e: any) {
    ui.err(e.message)
  }
}

function go(p: string) { router.push(p) }

function openRoomForm() {
  if (!clubs.value.length) return ui.err('先创建或加入俱乐部')
  roomForm.value = { gameType: 'nlh', seatCount: 6, sb: 20, bb: 40, minBuyin: 2000, maxBuyin: 8000, insurance: true, duration: 120 }
}

function openRecharge() {
  if (!clubs.value.length) return ui.err('先创建或加入俱乐部')
  rechargeOpen.value = true
}

async function doRecharge() {
  const cl = clubs.value[0]
  try {
    const d = await socket.call<{ cards: number; cardPool: number }>('club.recharge', { clubId: cl.id, count: rechargeCount.value })
    cl.cardPool = d.cardPool
    if (auth.profile) auth.profile.cards = d.cards
    ui.ok('已转入 ' + rechargeCount.value + ' 张')
    rechargeOpen.value = false
    load()
  } catch (e: any) { ui.err(e.message) }
}

async function createClubRoom() {
  const cl = clubs.value[0]
  if (!cl || !roomForm.value) return
  try {
    const d = await socket.call<{ roomNo: number; name: string; cardPool: number }>('club.createRoom', {
      clubId: cl.id, ...roomForm.value
    })
    cl.cardPool = d.cardPool
    roomForm.value = null
    ui.ok('房间已创建')
    await table.join(d.roomNo, 2000)
    router.push('/table/' + d.roomNo)
  } catch (e: any) { ui.err(e.message) }
}

// 会长发布公告后，在线成员实时收到推送
let offNotice: (() => void) | null = null
onMounted(() => {
  load()
  offNotice = socket.on('club.notice', (m: any) => {
    const hit = clubs.value.find((c: any) => c.id === m.clubId)
    if (hit) hit.notice = m.notice
    if (m.from !== auth.profile?.uid) ui.toast('俱乐部公告：' + m.notice, 'info', 3200)
  })
})
onUnmounted(() => { offNotice?.() })
</script>

<style scoped>
.club-tools .btn { white-space: nowrap; }
.mini-btn {
  width: 34px; height: 30px; border-radius: 9px; font-size: 15px; line-height: 1;
  display: inline-flex; align-items: center; justify-content: center; flex: none;
  background: rgba(30, 96, 72, .5); border: 1px solid rgba(167, 243, 208, .3); color: #d8f6e4;
}
.mini-btn:active { transform: scale(.94); }

.cert-tag {
  display: inline-flex; align-items: center; gap: 3px; height: 20px; padding: 0 8px;
  border-radius: 10px; font-size: 10px; font-weight: 800; color: #04220f;
  background: linear-gradient(180deg, #7ef0b4, #23c07a);
  box-shadow: 0 2px 8px rgba(35, 192, 122, .35);
}
.cert-tag.pending {
  color: #3a2a00; cursor: pointer;
  background: linear-gradient(180deg, #ffe9a8, #e7bf5c);
  box-shadow: 0 2px 8px rgba(212, 183, 104, .3);
}

.apply-box {
  margin-bottom: 10px; padding: 10px; border-radius: 12px;
  background: rgba(255, 216, 117, .06); border: 1px solid rgba(255, 216, 117, .24);
}

.club { padding-top: 12px; }
.page-head { display: flex; align-items: center; gap: 8px; }
.head-btn { height: 34px; padding: 0 12px; border-radius: 10px; background: rgba(255,255,255,.06); border: 1px solid var(--line); font-size: 13px; display: flex; align-items: center; gap: 4px; }
.club-actions { display: flex; gap: 12px; margin: 6px 0 0; }
.club-btn {
  flex: 1; height: 54px; border-radius: 14px; display: flex; align-items: center; justify-content: center;
  gap: 7px; font-size: 14px; font-weight: 800;
}
.club-btn.ghost { background: rgba(12,40,30,.9); border: 1px solid rgba(212,183,104,.55); color: #ffe6a8; }
.club-btn.primary { color: #06230f; }
.club-btn.primary {
  background: linear-gradient(180deg, #2be183, #10a75c); color: #04220f;
  border: 1px solid rgba(190,255,219,.45); box-shadow: 0 8px 20px rgba(22,190,108,.3);
}
.empty-state { display: flex; flex-direction: column; align-items: center; gap: 4px; padding: 78px 0 0; }
.crown-box { width: 96px; height: 96px; border-radius: 22px; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(160deg, rgba(24,60,42,.95), rgba(10,30,20,.98)); border: 1px solid rgba(239,207,142,.35);
  box-shadow: 0 10px 26px rgba(0,0,0,.4), inset 0 1px 0 rgba(255,255,255,.06); margin-bottom: 20px; }
.empty-title { font-size: 17px; font-weight: 800; color: #f2fbf6; }
.empty-desc { font-size: 14px; line-height: 1.8; color: #cfe6d9; text-align: center; max-width: 86%; margin: 10px 0 0; }
.club-list { display: flex; flex-direction: column; gap: 10px; }
.club-item { border-radius: 14px; padding: 12px; }
.club-rooms { margin-top: 6px; }
.club-tools { display: flex; gap: 8px; margin-bottom: 12px; }
.club-tools .btn { flex: 1; }
.grid2 { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; }
.fld { display: flex; flex-direction: column; gap: 5px; }
.fld > span { font-size: 11px; color: var(--muted); }
.fld .field { height: 40px; }
.chip { height: 34px; padding: 0 14px; border-radius: 17px; font-size: 12px; color: var(--muted); background: rgba(10,30,20,.7); border: 1px solid var(--line); }
.chip.on { color: #241a06; background: linear-gradient(180deg, #f2d493, #cfa151); font-weight: 700; border-color: transparent; }
.cost-line { display: flex; align-items: center; gap: 6px; font-size: 12px; color: #cfe8d9; }
.room { display: flex; align-items: center; gap: 9px; padding: 10px 12px; border-radius: 14px; background: rgba(12, 34, 22, .7); border: 1px solid var(--line); margin-bottom: 8px; }
.close-x { width: 32px; height: 32px; border-radius: 10px; background: rgba(255,255,255,.06); display: flex; align-items: center; justify-content: center; color: var(--muted); }

.notice-row { display: flex; align-items: center; gap: 6px; margin-top: 10px; }
.notice-text { flex: 1; min-width: 0; margin: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.notice-edit {
  flex: none; height: 24px; padding: 0 9px; border-radius: 8px; font-size: 10.5px; font-weight: 800;
  color: #ffe3a3; background: rgba(239, 207, 142, .12); border: 1px solid rgba(239, 207, 142, .3);
}
.notice-edit:active { transform: scale(.95); }
.notice-input { height: 84px; padding: 10px 12px; resize: none; line-height: 1.5; }

.sort-chip {
  height: 24px; padding: 0 9px; border-radius: 8px; font-size: 10.5px; font-weight: 700;
  color: #a9c4b6; background: rgba(255, 255, 255, .04); border: 1px solid var(--line);
}
.sort-chip.on { color: #06210f; background: linear-gradient(180deg, #7ef0b4, #2bbd7a); border-color: transparent; }

.cardlog-row { display: flex; align-items: center; gap: 10px; padding: 9px 0; border-bottom: 1px solid rgba(255, 255, 255, .05); }
.cardlog-row:last-child { border-bottom: 0; }
.cl-badge {
  flex: none; min-width: 42px; height: 26px; padding: 0 8px; border-radius: 8px; font-size: 12px; font-weight: 800;
  display: inline-flex; align-items: center; justify-content: center;
}
.cl-badge.in { color: #06210f; background: linear-gradient(180deg, #7ef0b4, #2bbd7a); }
.cl-badge.out { color: #2a0a0a; background: linear-gradient(180deg, #ffb3b3, #e0655a); }
.cl-reason { font-size: 12.5px; }

.horn-badge {
  font-size: 11px; line-height: 1; padding: 3px 6px; border-radius: 7px;
  background: rgba(255, 216, 117, .14); border: 1px solid rgba(255, 216, 117, .34);
}

/* 官方空态：96×96 圆角方块 + 金色细边 + 皇冠 */
.empty-icon {
  width: 96px; height: 96px; border-radius: 22px; display: flex; align-items: center; justify-content: center;
  background: rgba(14, 40, 28, .8); border: 1px solid rgba(214, 176, 92, .5);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, .05);
}
.empty-desc { font-size: 14px; line-height: 1.8; color: #cfe6d9; }

/* 俱乐部中心：与其它官方页同基准（安全区 12px、16px 基准字号、normal 行高） */
.club-page { --header-safe-pt: 12px; font-size: 16px; line-height: normal; }
.club-page .club-hub-view { flex: 1; min-height: 0; display: flex; flex-direction: column; }
.club-page .empty-hub-rooms { padding: 28px 16px; text-align: center; color: var(--muted); font-size: 13px; }
.club-page .empty-hub-create-btn { margin-top: 12px; height: 38px; padding: 0 18px; border-radius: 12px; font-size: 13px; font-weight: 700;
  color: #04220f; background: linear-gradient(180deg, #55ea9c, #12b767); }
.club-page .empty-hub-create-btn:disabled { opacity: .45; }
.club-page .hub-back-btn { width: 34px; height: 34px; border-radius: 12px; display: flex; align-items: center; justify-content: center;
  background: rgba(10, 36, 28, .8); border: 1px solid rgba(74, 133, 107, .5); color: #f4d47e; }
.club-page .hub-tool-btn { width: 34px; height: 34px; border-radius: 12px; display: flex; align-items: center; justify-content: center;
  background: rgba(10, 36, 28, .8); border: 1px solid rgba(74, 133, 107, .5); color: #d8ecdf; }
.club-page .hub-tab { height: 30px; padding: 0 12px; border-radius: 9px; font-size: 12.5px; font-weight: 700; color: #9bb6a4;
  background: transparent; border: 1px solid transparent; }
.club-page .hub-tab.active { color: #04220f; background: linear-gradient(180deg, #fde68a, #d4af37); }
.club-page .hub-tool { width: 32px; height: 32px; border-radius: 10px; display: flex; align-items: center; justify-content: center;
  background: rgba(10, 32, 22, .9); border: 1px solid rgba(120, 200, 160, .24); }

/* 空态文案与官方一致（官方全局样式未随组件导出，这里补回数值） */
.club-page .empty-clubs-state b { font-size: 16px; font-weight: 700; color: #fff; }
.club-page .empty-clubs-state p { max-width: 260px; margin: 0; font-size: 13px; line-height: 19.5px; color: #8da699; }
.club-page .shield-halo-box { width: 68px; height: 68px; display: flex; align-items: center; justify-content: center; border-radius: 20px; }

.club-page .hub-tool.primary { background: linear-gradient(180deg, #55ea9c, #12b767); border-color: transparent; }
</style>
