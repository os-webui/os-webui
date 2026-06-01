import './assets/main.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { createI18n } from 'vue-i18n'
import zhTW from './locale/zhTW.json'
import zhCN from './locale/zhCN.json'
import enUS from './locale/enUS.json'

const app = createApp(App)

app.use(createPinia())
app.use(createI18n({
  legacy: false,
  locale: 'en-US',
  fallbackLocale: 'en-US',
  messages: {
    'zh-TW': zhTW,
    'zh-CN': zhCN,
    'en-US': enUS,
  }
}))
app.use(router)

app.mount('#app')
