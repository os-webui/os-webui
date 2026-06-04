<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">
import {
  DarkModeOutlined, LightModeOutlined,
  CircleOutlined, CheckCircleOutlined,
} from '@vicons/material'
import { DarkTheme20Filled } from '@vicons/fluent'
import {
  NButton,
  NDropdown,
  NIcon,
} from 'naive-ui'
import type { DropdownOption, DropdownRenderOption, DropdownDividerOption } from 'naive-ui'
import { useThemeStore } from '@/stores/theme'
import { useI18n } from 'vue-i18n'
import { renderLabel } from '@/internal/render/label'
const i18n = useI18n()
const theme = useThemeStore()

const menus: Array<DropdownOption | DropdownRenderOption | DropdownDividerOption> = [
  {
    key: 'dark',
    label: renderLabel({
      label: () => i18n.t('theme.dark'),
      prefixRender: () => theme.choose === 'dark' ? CheckCircleOutlined : CircleOutlined,
      suffix: DarkModeOutlined,
    }),
  },
  {
    key: 'light',
    label: renderLabel({
      label: () => i18n.t('theme.light'),
      prefixRender: () => theme.choose === 'light' ? CheckCircleOutlined : CircleOutlined,
      suffix: LightModeOutlined,
    }),
  },
  {
    key: 'divider',
    type: 'divider',
  },
  {
    key: 'auto',
    label: renderLabel({
      label: () => i18n.t('theme.auto'),
      prefixRender: () => theme.choose === 'auto' ? CheckCircleOutlined : CircleOutlined,
      suffix: DarkTheme20Filled,
    }),
  },
]
function handleSelect(key: string,
  // options: DropdownOption | DropdownRenderOption | DropdownDividerOption
) {
  theme.choose = key
}
defineProps<{
  placement?: any
}>()
</script>
<template>
  <n-dropdown trigger="hover" :show-arrow="true" :placement="placement" :options="menus" @select="handleSelect">
    <!-- <n-dropdown :show="true" :show-arrow="true" :placement="placement" :options="menus" @select="handleSelect"> -->
    <n-button :text="true">
      <template #icon>
        <n-icon v-if="theme.name == 'dark'" size="1.3rem">
          <DarkModeOutlined />
        </n-icon>
        <n-icon v-if="theme.name != 'dark'" size="1.3rem">
          <LightModeOutlined />
        </n-icon>
      </template>
    </n-button>
  </n-dropdown>
</template>
