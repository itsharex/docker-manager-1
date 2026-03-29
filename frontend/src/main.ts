import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { initSettings } from './ui/settings'
import { initAuth } from './ui/auth'
import router from './router'
import { i18n } from './i18n'

initSettings()
initAuth()
createApp(App).use(router).use(i18n).mount('#app')
