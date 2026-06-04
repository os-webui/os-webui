<script setup lang="ts">
import { onUnmounted, onMounted, watchEffect } from 'vue'
import NavBar from '@/ui/NavBar.vue'
import { RouterView, useRouter } from 'vue-router'
import {
  NConfigProvider, NGlobalStyle,
  NButton, NIcon, NTooltip, NMessageProvider,
} from 'naive-ui'
import {
  HomeOutlined, InfoOutlined, StorefrontOutlined,
} from '@vicons/material'
import { Github } from '@vicons/fa'

import { useThemeStore } from '@/stores/theme'
import { useBreakpointStore } from '@/stores/breakpoint'
import { useI18n } from 'vue-i18n'

import ThemeMenu from '@/ui/ThemeMenu.vue'
import LangMenu from '@/ui/LangMenu.vue'
// import DevMenu from '@/ui/DevMenu.vue'

import { useLocaleStore } from './stores/locale'
import { useTitle } from './stores/title'

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


const title = useTitle()
useRouter().beforeEach((to) => {
  const pageTitle = to.meta.title
  title.set(typeof pageTitle === "string" || Array.isArray(pageTitle) ? pageTitle : [])
})
watchEffect(() => {
  title.updateTitle()
  const main = title.title
  let get = false
  const strs: string[] = title.get.map((v) => {
    if (v.startsWith('main.')) {
      get = true
      return i18n.t(v)
    }
    return v
  })
  if (!get) {
    i18n.t('main.home')
  }
  strs.push(main)
  document.title = strs.join(' - ')
})
</script>

<template>
  <n-config-provider inline-theme-disabled :theme="theme.theme" :locale="locale.locale.message"
    :date-locale="locale.locale.date">
    <n-config-provider inline-theme-disabled :theme="theme.theme">
      <n-message-provider :keep-alive-on-hover="true" :duration="5000">
        <n-global-style />

        <header class="sticky top-0 z-5">
          <NavBar>
            <!-- brand 無論手機還是桌面都會顯示到 導航欄左側 -->
            <template v-slot:brand>
              <n-tooltip trigger="hover" placement="bottom">
                <template #trigger>
                  <RouterLink to="/" class="flex align-items-center justify-content-center">
                    <n-button :text="true">
                      <n-icon size="1.3rem">
                        <HomeOutlined />
                      </n-icon>
                    </n-button>
                  </RouterLink>
                </template>
                {{ $t('main.home') }}
              </n-tooltip>
            </template>

            <!-- menu 在桌面系統顯示到 導航欄左側，手機顯示到導航欄 摺疊部分 -->
            <!-- <template v-slot:menu>
            <DevMenu :placement="breakpoint.md ? 'bottom-start' : 'left-start'" />
          </template> -->

            <!-- right-menu 在桌面系統顯示到 導航欄右側，手機顯示到導航欄 摺疊部分 -->
            <template v-slot:right-menu>
              <n-tooltip trigger="hover" :placement="breakpoint.md ? 'bottom' : 'left'">
                <template #trigger>
                  <RouterLink to="/store" class="flex align-items-center justify-content-center">
                    <n-button :text="true">
                      <n-icon size="1.3rem">
                        <StorefrontOutlined />
                      </n-icon>
                    </n-button>
                  </RouterLink>
                </template>
                {{ $t('main.store') }}
              </n-tooltip>
              <n-tooltip trigger="hover" placement="bottom">
                <template #trigger>
                  <RouterLink to="/about" class="flex align-items-center justify-content-center">
                    <n-button :text="true">
                      <n-icon size="1.3rem">
                        <InfoOutlined />
                      </n-icon>
                    </n-button>
                  </RouterLink>
                </template>
                {{ $t('main.about') }}
              </n-tooltip>
            </template>
            <!-- right-brand 無論手機還是桌面都會顯示到 導航欄右側 -->
            <template v-slot:right-brand>
              <LangMenu placement="bottom-end" />
              <ThemeMenu placement="bottom-end" />
            </template>
          </NavBar>
        </header>

        <main class="flex justify-content-center pt-3 pb-6">
          <div class="container">
            <RouterView />
          </div>
        </main>

        <footer :class="theme.name == 'dark' ? 'footer-dark' : 'footer-light'"
          class="flex justify-content-center flex-wrap">
          <div class="container">
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
      </n-message-provider>
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
