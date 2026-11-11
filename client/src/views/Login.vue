<template>
  <div class="login">
    <div class="glow"></div>
    <div class="brand">
      <div class="logo logo-img">
        <img :src="logo" alt="U5" />
      </div>
      <h1 class="h1">U5 德扑</h1>
      <p class="muted sm">德州扑克 · 短牌 · 奥马哈 · 牛牛</p>
    </div>

    <div class="card card-pad form">
      <div class="segmented" style="margin-bottom:14px">
        <button class="seg" :class="{ active: mode === 'login' }" @click="mode = 'login'">登录</button>
        <button class="seg" :class="{ active: mode === 'register' }" @click="mode = 'register'">注册</button>
      </div>
      <div class="col gap-10">
        <input v-model.trim="username" class="field" placeholder="账号（3-20位字母数字）" autocomplete="username" />
        <input v-if="mode === 'register'" v-model.trim="nick" class="field" placeholder="昵称" />
        <input v-model="password" class="field" type="password" placeholder="密码（6位以上）" autocomplete="current-password" @keyup.enter="submit" />
      </div>
      <button class="btn btn-primary btn-block" style="margin-top:16px;height:50px" :disabled="busy" @click="submit">
        {{ busy ? '处理中…' : mode === 'login' ? '进入牌桌' : '立即注册' }}
      </button>
      <button class="btn btn-ghost btn-block" style="margin-top:10px" :disabled="busy" @click="guest">
        游客体验 · 送 $50 + 2 房卡
      </button>
      <p class="tiny muted" style="text-align:center;margin:14px 0 0">
        安全卫士已启动 · 全链路加密 · 公平哈希存证
      </p>
    </div>

    <div class="perks">
      <div class="perk"><b>12</b><span>官方牌桌</span></div>
      <div class="perk"><b>1:1</b><span>链上承兑</span></div>
      <div class="perk"><b>30</b><span>在线机器人</span></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { useUI } from '@/stores/ui'
import { sfx } from '@/lib/sfx'

const router = useRouter()
const route = useRoute()
const auth = useAuth()
const ui = useUI()

const logo = '/img/u5-qipai-logo-login.webp'
const mode = ref<'login' | 'register'>('login')
const username = ref('')
const nick = ref('')
const password = ref('')
const busy = ref(false)

async function api(path: string, body: any) {
  const r = await fetch('/api/auth/' + path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  const data = await r.json()
  if (!data.ok) throw new Error(data.e || '请求失败')
  return data.d
}

async function submit() {
  if (busy.value) return
  if (username.value.length < 3) return ui.err('账号至少 3 位')
  if (password.value.length < 6) return ui.err('密码至少 6 位')
  busy.value = true
  try {
    const d = await api(mode.value, { username: username.value, password: password.value, nick: nick.value })
    sfx.win()
    auth.setSession(d.token, d)
    router.replace(afterLogin())
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    busy.value = false
  }
}

function afterLogin() {
  const r = route.query.redirect
  if (typeof r === 'string' && r.startsWith('/') && !r.startsWith('//')) return r
  return '/lobby'
}

async function guest() {
  if (busy.value) return
  busy.value = true
  try {
    const d = await api('guest', {})
    auth.setSession(d.token, d)
    router.replace(afterLogin())
  } catch (e: any) {
    ui.err(e.message)
  } finally {
    busy.value = false
  }
}
</script>

<style scoped>
.login { position: relative; min-height: 100vh; padding: 52px 22px 34px; display: flex; flex-direction: column; justify-content: center; gap: 26px; }
.glow { position: absolute; inset: -30% -20% auto -20%; height: 70%; background: radial-gradient(ellipse at 50% 40%, rgba(46, 224, 122, .28), transparent 62%); filter: blur(30px); pointer-events: none; }
.brand { text-align: center; position: relative; }
.logo-img { width: 128px; height: 88px; background: transparent; border: none; box-shadow: none; }
.logo-img img { width: 100%; height: 100%; object-fit: contain; filter: drop-shadow(0 10px 26px rgba(0, 0, 0, .5)); }
.logo { width: 78px; height: 78px; margin: 0 auto 14px; border-radius: 24px; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(160deg, rgba(46, 224, 122, .18), rgba(10, 40, 26, .9)); border: 1px solid rgba(46, 224, 122, .3);
  box-shadow: 0 12px 40px rgba(20, 180, 100, .25), inset 0 1px 0 rgba(255,255,255,.1); }
.logo-mark { font-size: 40px; color: var(--gold); text-shadow: 0 0 20px rgba(255, 220, 140, .6); }
.form { position: relative; }
.perks { display: flex; justify-content: space-around; }
.perk { display: flex; flex-direction: column; align-items: center; gap: 2px; }
.perk b { font-size: 18px; color: var(--gold); font-weight: 800; }
.perk span { font-size: 11px; color: var(--muted); }
</style>
