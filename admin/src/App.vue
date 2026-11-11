<template>
  <router-view />
  <div class="toasts">
    <div v-for="t in toasts" :key="t.id" class="toast" :class="t.kind">{{ t.text }}</div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue'

const toasts = ref<{ id: number; text: string; kind: string }[]>([])
let seq = 1
function push(text: string, kind = 'ok') {
  const id = seq++
  toasts.value.push({ id, text, kind })
  setTimeout(() => { toasts.value = toasts.value.filter(t => t.id !== id) }, 2600)
}
onMounted(() => window.addEventListener('admin-toast', ((e: any) => push(e.detail.text, e.detail.kind)) as any))
onBeforeUnmount(() => window.removeEventListener('admin-toast', ((e: any) => push(e.detail.text, e.detail.kind)) as any))
</script>
