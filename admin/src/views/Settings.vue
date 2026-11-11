<template>
  <div class="grid" style="grid-template-columns: 1fr 1fr">
    <div class="card">
      <b>运营参数</b>
      <div class="settings">
        <label v-for="k in keys" :key="k">
          <span>{{ label(k) }}</span>
          <input v-model="form[k]" class="input" />
        </label>
      </div>
      <button class="btn btn-primary" style="margin-top:12px" @click="save">保存设置</button>
    </div>
    <div class="card">
      <b>全服广播</b>
      <p class="muted tiny">通过 WebSocket 推送给所有在线客户端</p>
      <textarea v-model="broadcast" class="input" style="width:100%;height:90px;padding:10px;resize:vertical" placeholder="输入广播内容…"></textarea>
      <button class="btn btn-primary" style="margin-top:10px" @click="sendBroadcast">立即广播</button>
      <div style="margin-top:18px;border-top:1px solid #14261b;padding-top:14px">
        <b>关于</b>
        <p class="muted tiny" style="line-height:1.7">
          U5 德扑 · 1:1 复刻版<br />
          客户端 Vue3 / 服务端 Go / 数据库 MySQL 5.7<br />
          全量消息通过 WebSocket 推送，服务端权威发牌与结算
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api } from '@/lib/api'

const keys = ['site_name', 'announcement', 'withdraw_fee', 'robot_fill', 'min_buyin_rate']
const form = ref<Record<string, string>>({})
const broadcast = ref('')

const label = (k: string) => ({
  site_name: '站点名称', announcement: '大厅公告', withdraw_fee: '提现费率(%)',
  robot_fill: '机器人自动填充(1/0)', min_buyin_rate: '最小买入倍数'
} as any)[k] || k

async function load() {
  try { form.value = await api.get<Record<string, string>>('settings') } catch {}
}

async function save() {
  try {
    for (const k of keys) await api.post('setting', { k, v: String(form.value[k] ?? '') })
    toast('设置已保存')
  } catch (e: any) { toast(e.message, 'err') }
}

async function sendBroadcast() {
  if (!broadcast.value.trim()) return
  try {
    await api.post('broadcast', { text: broadcast.value })
    toast('广播已推送')
    broadcast.value = ''
  } catch (e: any) { toast(e.message, 'err') }
}

function toast(text: string, kind = 'ok') {
  window.dispatchEvent(new CustomEvent('admin-toast', { detail: { text, kind } }))
}

onMounted(load)
</script>

<style scoped>
.settings { display: flex; flex-direction: column; gap: 10px; margin-top: 12px; }
.settings label { display: flex; flex-direction: column; gap: 5px; font-size: 11px; color: var(--muted); }
.settings .input { height: 34px; }
</style>
