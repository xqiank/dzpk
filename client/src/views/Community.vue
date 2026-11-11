<template>
  <div class="page community">
    <header class="c-head">
      <button class="head-btn" @click="router.push('/lobby')"><Icon name="back" :size="16" /> 返回</button>
      <div class="col" style="align-items:center">
        <b class="h3">U5德州官方交流社区 <span class="tag tag-gold" style="height:17px">官方认证</span></b>
        <span class="tiny green"><i class="live"></i> 全员开放 · {{ online }} 人在线 · {{ seated }} 人在座</span>
      </div>
      <span style="width:60px"></span>
    </header>

    <div class="notice-bar">
      <span class="tag tag">官方公告</span>
      <span class="tiny">欢迎加入 U5德州官方交流社区，禁止发布广告与违规内容</span>
    </div>

    <div class="msgs" ref="listEl">
      <div v-for="(m, i) in messages" :key="i" class="msg">
        <Avatar :src="m.avatar" :size="34" />
        <div class="col" style="flex:1;min-width:0">
          <span class="who">{{ m.nick }} <em class="tiny">{{ fmtTime(m.time) }}</em></span>
          <div class="bubble">{{ m.content }}</div>
        </div>
      </div>
    </div>

    <div class="quick-row scroll-x">
      <button v-for="q in quickPhrases" :key="q" class="qchip" @click="send(q)">{{ q }}</button>
    </div>

    <div class="composer">
      <input v-model="text" class="field" placeholder="在官方群说点什么…" @keyup.enter="send(text)" />
      <button class="btn btn-primary send" :disabled="!text.trim()" @click="send(text)">发送</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import Avatar from '@/components/Avatar.vue'
import Icon from '@/components/Icon.vue'
import { socket } from '@/lib/ws'
import { fmtTime } from '@/lib/format'
import { useAuth } from '@/stores/auth'
import { useUI } from '@/stores/ui'
import { sfx } from '@/lib/sfx'

const router = useRouter()
const auth = useAuth()
const ui = useUI()
const messages = ref<any[]>([])
const text = ref('')
const listEl = ref<HTMLElement | null>(null)
const online = ref(0)
const seated = ref(0)
const quickPhrases = ['赶紧上桌，就差你了！', '这把必须跟！', '抓个同花顺！', '老板大气，再发一车！']

function scrollBottom() {
  nextTick(() => {
    if (listEl.value) listEl.value.scrollTop = listEl.value.scrollHeight
  })
}

async function load() {
  try {
    messages.value = await socket.call('chat.world.history')
    if (!messages.value.length) {
      messages.value = [
        { uid: 0, nick: '官方', avatar: '/avatars/av1.svg', content: '欢迎来到 U5 德州官方交流社区，文明交流、纯粹竞技。', time: Date.now() - 3600000 },
        { uid: 0, nick: '孤岛沉思', avatar: '/avatars/av16.svg', content: '大厅哪桌有鱼？来个带路的', time: Date.now() - 600000 }
      ]
    }
    scrollBottom()
  } catch {}
  refreshOnline()
}

async function send(t: string) {
  const v = (t || '').trim()
  if (!v) return
  try {
    await socket.call('chat.world.send', { text: v })
    text.value = ''
    sfx.tap()
  } catch (e: any) { ui.err(e.message) }
}

function onWorld(e: Event) {
  const m = (e as CustomEvent).detail
  if (!m) return
  messages.value.push(m)
  if (messages.value.length > 120) messages.value.shift()
  scrollBottom()
}

async function refreshOnline() {
  try {
    const d = await socket.call<{ online: number; seated: number }>('chat.world.online')
    online.value = d.online || 0
    seated.value = d.seated || 0
  } catch { /* 忽略 */ }
}

let onlineTimer = 0
onMounted(() => {
  load()
  onlineTimer = window.setInterval(refreshOnline, 5000)
  window.addEventListener('u5-world', onWorld as EventListener)
})
onBeforeUnmount(() => {
  window.removeEventListener('u5-world', onWorld as EventListener)
  clearInterval(onlineTimer)
})
</script>

<style scoped>
.community { display: flex; flex-direction: column; height: 100vh; padding-bottom: calc(84px + var(--safe-b)); }
.c-head { display: flex; align-items: center; gap: 8px; padding-top: 14px; }
.head-btn { height: 34px; padding: 0 12px; border-radius: 10px; background: rgba(255,255,255,.06); border: 1px solid var(--line); font-size: 13px; display: flex; align-items: center; gap: 4px; }
.live { display: inline-block; width: 6px; height: 6px; border-radius: 50%; background: #2ee07a; margin-right: 2px; box-shadow: 0 0 8px rgba(46,224,122,.9); }
.notice-bar { display: flex; align-items: center; gap: 8px; margin: 10px 0; padding: 8px 10px; border-radius: 10px; background: linear-gradient(90deg, rgba(20,64,42,.9), rgba(8,26,17,.5)); border: 1px solid var(--line); white-space: nowrap; overflow: hidden; }
.msgs { flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 10px; padding: 6px 0 10px; }
.msg { display: flex; gap: 8px; }
.who { font-size: 11px; color: #9ec7b3; margin-bottom: 3px; }
.who em { font-style: normal; color: var(--muted-2); margin-left: 4px; }
.bubble { display: inline-block; max-width: 100%; padding: 7px 11px; border-radius: 12px; font-size: 13px; line-height: 1.45;
  background: linear-gradient(160deg, rgba(26,58,42,.9), rgba(14,36,25,.92)); border: 1px solid var(--line); box-shadow: 0 4px 12px rgba(0,0,0,.25); }
.quick-row { display: flex; gap: 6px; padding: 6px 0; }
.qchip { height: 26px; padding: 0 9px; border-radius: 8px; font-size: 11px; color: #cfe8d9; background: rgba(14,40,27,.8); border: 1px solid var(--line); white-space: nowrap; }
.composer { position: fixed; bottom: 0; left: 50%; transform: translateX(-50%); width: 100%; max-width: 460px; display: flex; gap: 8px; padding: 10px 14px calc(12px + var(--safe-b)); background: linear-gradient(180deg, rgba(4,18,11,.6), rgba(3,14,9,.98)); border-top: 1px solid var(--line); }
.composer .field { height: 42px; }
.send { width: 76px; height: 42px; padding: 0; }
</style>
