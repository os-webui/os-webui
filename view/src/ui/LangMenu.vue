<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">

import {
  CircleOutlined, CheckCircleOutlined,
  TranslateOutlined,
} from '@vicons/material'
import {
  NButton,
  NDropdown,
  NIcon,
} from 'naive-ui';
import type {
  DropdownOption, DropdownRenderOption, DropdownDividerOption,
} from 'naive-ui';
import { useLocaleStore, Locales } from '@/stores/locale'
import { useI18n } from 'vue-i18n'
import { renderLabel } from '@/internal/render/label'

const i18n = useI18n()
const locale = useLocaleStore()

const menus: Array<DropdownOption | DropdownRenderOption | DropdownDividerOption> =
  Locales.map((v) => {
    return {
      key: v.id,
      label: renderLabel({
        label: v.name,
        prefixRender: () => locale.choose === v.id ? CheckCircleOutlined : CircleOutlined,
      }),
    }
  });
menus.push(
  {
    key: 'divider',
    type: 'divider',
  },
  {
    key: 'auto',
    label: renderLabel({
      label: () => i18n.t('lang.auto'),
      prefixRender: () => locale.choose === 'auto' ? CheckCircleOutlined : CircleOutlined,
    }),
  }
)
function handleSelect(key: string,
  // options: DropdownOption | DropdownRenderOption | DropdownDividerOption
) {
  locale.choose = key
}
const props = defineProps<{
  placement?: any
}>()
</script>
<template>
  <n-dropdown trigger="hover" :placement="props.placement" :show-arrow="true" :options="menus" @select="handleSelect">
    <!-- <n-dropdown :show="true" :placement="props.placement" :show-arrow="true" :options="menus" @select="handleSelect"> -->
    <n-button :text="true">
      <template #icon>
        <n-icon size="1.3rem">
          <TranslateOutlined />
        </n-icon>
      </template>
    </n-button>
  </n-dropdown>
</template>
