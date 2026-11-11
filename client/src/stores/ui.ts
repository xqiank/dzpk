import { defineStore } from 'pinia'
import { ref } from 'vue'
import { sfx } from '@/lib/sfx'

export const useUI = defineStore('ui', () => {
  const toasts = ref<{ id: number; text: string; kind: 'ok' | 'err' | 'info' }[]>([])
  const sound = ref(localStorage.getItem('u5_sound') !== '0')
  const voice = ref(localStorage.getItem('u5_voice') !== '0')
  const voicePack = ref(localStorage.getItem('u5_voice_pack') || 'sweet')
  let seq = 1

  function toast(text: string, kind: 'ok' | 'err' | 'info' = 'info', ms = 2200) {
    const id = seq++
    toasts.value.push({ id, text, kind })
    setTimeout(() => { toasts.value = toasts.value.filter(t => t.id !== id) }, ms)
  }
  function ok(text: string) { toast(text, 'ok') }
  function err(text: string) { toast(text, 'err') }

  function toggleSound() {
    sound.value = !sound.value
    sfx.setEnabled(sound.value)
  }
  function toggleVoice() {
    voice.value = !voice.value
    sfx.setVoice(voice.value)
  }
  function setVoicePack(pack: string) {
    voicePack.value = pack
    sfx.setPack(pack as any)
  }

  return { toasts, sound, voice, voicePack, toast, ok, err, toggleSound, toggleVoice, setVoicePack }
})
