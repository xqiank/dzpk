import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { socket, netState } from '@/lib/ws'
import { sfx } from '@/lib/sfx'
import router from './router'
import './styles/main.css'
import './styles/table.css'
import './styles/insurance-official.css'
import './styles/table-menu-official.css'
import './styles/buyin-official.css'
import './styles/table-settle-official.css'
import './styles/table-deal-official.css'
import './styles/table-reveal-official.css'
import './styles/table-winner-aura-official.css'
import './styles/rebuy-official.css'
import './styles/guide-official.css'
import './styles/usdt-official.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
// 诊断句柄：客服/QA 可在控制台查看连接状态并手动触发重连
;(window as any).__u5socket = socket
;(window as any).__u5sfx = sfx
;(window as any).__u5net = netState

app.mount('#app')
