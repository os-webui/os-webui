<!-- eslint-disable @typescript-eslint/no-explicit-any -->
<script setup lang="ts">

import {
  BackupTableFilled, BugReportOutlined,
} from '@vicons/material'

import {
  NButton,
  NDropdown,
  NIcon,
} from 'naive-ui';
import { renderLabel } from '@/internal/render/label'


const menus = [
  {
    key: 'table',
    label: renderLabel({
      label: 'table',
      url: '/example/table',
      prefix: BackupTableFilled,
    }),
  },

]

const props = defineProps<{
  placement?: any
}>()
</script>
<template>
  <n-dropdown trigger="hover" :placement="props.placement" :show-arrow="true" :options="menus">
    <!-- <n-dropdown :show="true" :placement="props.placement" :show-arrow="true" :options="menus"> -->
    <n-button :text="true" class="flex align-items-center justify-content-center">
      <n-icon size="1.3rem">
        <BugReportOutlined />
      </n-icon>
    </n-button>
  </n-dropdown>
</template>
