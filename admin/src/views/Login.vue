<template>
  <div class="login-page">
    <div class="glow"></div>
    <div class="login-card">
      <div class="brand">
        <div class="mark">♠</div>
        <div class="col">
          <b class="title">U5 德扑 · 运营后台</b>
          <span class="muted tiny">Operations Console</span>
        </div>
      </div>
      <label class="fld"><span>管理员账号</span>
        <input v-model.trim="username" class="input" placeholder="admin" @keyup.enter="submit" />
      </label>
      <label class="fld"><span>密码</span>
        <input v-model="password" class="input" type="password" placeholder="••••••" @keyup.enter="submit" />
      </label>
      <button class="btn btn-primary submit" :disabled="busy" @click="submit">{{ busy ? '登录中…' : '登录后台' }}</button>
      <p class="muted tiny hint">默认账号 admin / admin888 · 首次登录请及时修改</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '@/lib/api'

const router = useRouter()
const username = ref('admin')
const password = ref('admin888')
const busy = ref(false)

async function submit() {
  if (busy.value) return
  busy.value = true
  try {
    await login(username.value, password.value)
    router.replace('/dashboard')
  } catch (e: any) {
    window.dispatchEvent(new CustomEvent('admin-toast', { detail: { text: e.message, kind: 'err' } }))
  } finally {
    busy.value = false
  }
}
</script>

<style scoped>
.login-page { height: 100%; display: flex; align-items: center; justify-content: center; position: relative; overflow: hidden; }
.glow { position: absolute; width: 900px; height: 900px; border-radius: 50%; background: radial-gradient(circle, rgba(46, 224, 122, .14), transparent 62%); top: -380px; left: 50%; transform: translateX(-50%); }
.login-card { position: relative; width: 360px; padding: 26px; border-radius: 16px; background: linear-gradient(160deg, #0e2117, #08130d); border: 1px solid var(--line); box-shadow: 0 30px 80px rgba(0, 0, 0, .6); display: flex; flex-direction: column; gap: 14px; }
.brand { display: flex; align-items: center; gap: 12px; margin-bottom: 6px; }
.mark { width: 44px; height: 44px; border-radius: 12px; background: rgba(46, 224, 122, .12); border: 1px solid #255c3d; color: var(--gold); display: flex; align-items: center; justify-content: center; font-size: 22px; }
.title { font-size: 16px; }
.tiny { font-size: 11px; }
.fld { display: flex; flex-direction: column; gap: 6px; }
.fld span { font-size: 11px; color: var(--muted); }
.fld .input { height: 38px; }
.submit { height: 40px; margin-top: 4px; font-size: 13px; border-radius: 10px; }
.hint { text-align: center; margin: 2px 0 0; }
.col { display: flex; flex-direction: column; }
</style>
