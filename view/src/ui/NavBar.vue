<script setup lang="ts">
import { ref } from 'vue'
import {
  MenuOutlined,
  MenuOpenOutlined,
} from '@vicons/material'

import {
  NEl,
  NButton, NIcon,
} from 'naive-ui'
import { useThemeStore } from '@/stores/theme'
import { useBreakpointStore } from '@/stores/breakpoint'

const theme = useThemeStore()
const breakpoint = useBreakpointStore()

const open = ref(false)
</script>

<template>
  <n-el>
    <!-- 爲桌面系統佈局導航 -->
    <div :class="theme.name == 'dark' ? 'navbar-dark' : 'navbar-light'"
      class="flex justify-content-center align-items-center" v-if="breakpoint.md">
      <div class="container flex">
        <div class="flex flex-grow-1 gap-2">
          <slot name="brand"></slot>
          <slot name="menu"></slot>
        </div>

        <div class="flex flex-grow-0 gap-2">
          <slot name="right-menu"></slot>
          <slot name="right-brand"></slot>
        </div>
      </div>
    </div>
    <!-- 爲手機佈局導航 -->
    <div :class="theme.name == 'dark' ? 'navbar-dark' : 'navbar-light'"
      class="flex justify-content-center align-items-center" v-else>
      <div class="container">
        <div class="flex flex-column align-items-start gap-2 w-full">
          <div class="flex flex-row align-items-center w-full">
            <div class="flex flex-grow-1 align-items-center gap-2">
              <slot name="brand"></slot>
            </div>
            <div class="flex flex-grow-0 align-items-center gap-2">
              <slot name="right-brand"></slot>
              <div class="flex pr-2">
                <n-button :text="true" @click="open = !open">
                  <template #icon>
                    <n-icon v-if="open" size="1.3rem">
                      <MenuOpenOutlined />
                    </n-icon>
                    <n-icon v-if="!open" size="1.3rem">
                      <MenuOutlined />
                    </n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </div>
          <div class="flex grid w-full" v-if="open">
            <div class="col-2 col-offset-10">
              <div class="flex flex-column align-items-center gap-2 w-full">
                <slot name="menu"></slot>
                <slot name="right-menu"></slot>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </n-el>
</template>


<style scoped>
/* .navbar {
    padding: 0.75rem 0;
    background-color: var(--card-color);
    box-shadow: var(--box-shadow-1);
} */
.navbar-dark {
  padding: 0.75rem 0;
  background-color: hsla(221deg, 14%, 21%, 1);
  box-shadow: var(--box-shadow-1);
}

.navbar-light {
  padding: 0.75rem 0;
  background-color: hsla(221deg, 14%, 96%, 1);
  box-shadow: var(--box-shadow-1);
}
</style>
