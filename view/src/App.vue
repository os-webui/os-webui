<script setup lang="ts">
import { onUnmounted, onMounted, watchEffect } from 'vue'
import NavBar from '@/components/NavBar.vue'
import { RouterView } from 'vue-router'
import {
  NConfigProvider, NGlobalStyle,
  NButton, NIcon,
} from 'naive-ui'
import {
  HomeOutlined, InfoOutlined,
} from '@vicons/material'
import { Github } from '@vicons/fa'

import { useThemeStore } from '@/stores/theme'
import { useBreakpointStore } from '@/stores/breakpoint'
import { useI18n } from 'vue-i18n'

import ThemeMenu from '@/components/ThemeMenu.vue'
import LangMenu from '@/components/LangMenu.vue'
import DevMenu from '@/components/DevMenu.vue'

import { useLocaleStore } from './stores/locale'
const theme = useThemeStore()
const breakpoint = useBreakpointStore()
const i18n = useI18n()
const locale = useLocaleStore()
watchEffect(() => {
  i18n.locale.value = locale.locale.id
})
let cleanupLocale: (() => void) | undefined
let cleanupBreakpoint: (() => void) | undefined
onMounted(() => {
  cleanupLocale = locale.start()
  cleanupBreakpoint = breakpoint.start()
})
onUnmounted(() => {
  if (cleanupLocale) {
    cleanupLocale()
  }
  if (cleanupBreakpoint) {
    cleanupBreakpoint()
  }
})
i18n.locale
</script>

<template>
  <n-config-provider inline-theme-disabled :theme="theme.theme" :locale="locale.locale.message"
    :date-locale="locale.locale.date">
    <n-config-provider inline-theme-disabled :theme="theme.theme">
      <n-global-style />

      <header class="sticky top-0 z-5">
        <NavBar>
          <!-- brand 無論手機還是桌面都會顯示到 導航欄左側 -->
          <template v-slot:brand>
            <RouterLink to="/" class="flex align-items-center justify-content-center">
              <n-button :text="true">
                <n-icon size="1.3rem">
                  <HomeOutlined />
                </n-icon>
              </n-button>
            </RouterLink>
          </template>

          <!-- menu 在桌面系統顯示到 導航欄左側，手機顯示到導航欄 摺疊部分 -->
          <template v-slot:menu>
            <DevMenu :placement="breakpoint.md ? 'bottom-start' : 'left-start'" />
          </template>

          <!-- right-menu 在桌面系統顯示到 導航欄右側，手機顯示到導航欄 摺疊部分 -->
          <template v-slot:right-menu>
            <RouterLink to="/about" class="flex align-items-center justify-content-center">
              <n-button :text="true">
                <n-icon size="1.3rem">
                  <InfoOutlined />
                </n-icon>
              </n-button>
            </RouterLink>
          </template>
          <!-- right-brand 無論手機還是桌面都會顯示到 導航欄右側 -->
          <template v-slot:right-brand>
            <LangMenu placement="bottom-end" />
            <ThemeMenu placement="bottom-end" />
          </template>
        </NavBar>
      </header>

      <main class="flex justify-content-center pt-3 pb-6">
        <div class="container ">
          <RouterView />
        </div>
      </main>

      <footer :class="theme.name == 'dark' ? 'footer-dark' : 'footer-light'"
        class="flex justify-content-center flex-wrap">
        <div class="container ">
          <a href="https://github.com/os-webui/os-webui" target="_blank">
            <n-button :text="true">
              <template #icon>
                <n-icon>
                  <Github />
                </n-icon>
              </template>
              https://github.com/os-webui/os-webui
            </n-button>
          </a>
        </div>
      </footer>
    </n-config-provider>
  </n-config-provider>
</template>

<style scoped>
.footer-dark {
  background-color: hsl(221, 14%, 11%);
  padding: 3rem 1.5rem 6rem;
}

.footer-light {
  background-color: hsl(221, 14%, 98%);
  padding: 3rem 1.5rem 6rem;
}
</style>
