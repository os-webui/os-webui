<script setup lang="ts">
import { useTitle } from '@/stores/title';
import { createConfigView, type Props } from './ConfigView'
import {
  NCard, NButton, NSpace, NFormItem, NPopconfirm, NIcon
} from 'naive-ui';
import { DeleteSweepOutlined, RestartAltOutlined, SaveOutlined } from '@vicons/material'
import InputTextarea from '@/ui/input/InputTextarea.vue'
import { useNav } from '../stores/plugin';
const props = defineProps<Props>()
const title = useTitle()
title.set(['main.config', props.data.info.name])
const nav = useNav()
nav.setPlugin(props.plugin, props.data.info.name)


const { textValue,
  disabled, disabledClear, disabledReset,
  actions } = createConfigView(props)
const autosize = { minRows: 3, maxRows: 5 }
const config = `$ConfigDir/${props.plugin}/plugin.txt`
</script>
<template>
  <n-card :title="$t('main.config') + ' - ' + props.data.info.name" :bordered="false">
    <n-form-item :label="config">
      <InputTextarea class="textarea" v-model:value="textValue" :disabled="disabled" :autosize="autosize" />
    </n-form-item>

    <template #footer>
      <n-space justify="end">
        <n-popconfirm @positive-click="actions.clear()">
          <template #trigger>
            <n-button :disabled="disabledClear">
              <template #icon>
                <NIcon>
                  <DeleteSweepOutlined />
                </NIcon>
              </template>
              {{ $t('ui.clear') }}</n-button>
          </template>
          {{ $t('ui.clearConfirm') }}
        </n-popconfirm>
        <n-popconfirm @positive-click="actions.reset()">
          <template #trigger>
            <n-button :disabled="disabledReset">
              <template #icon>
                <NIcon>
                  <RestartAltOutlined />
                </NIcon>
              </template>
              {{ $t('ui.reset') }}</n-button>
          </template>
          {{ $t('ui.resetConfirm') }}
        </n-popconfirm>

        <n-popconfirm @positive-click="actions.save()">
          <template #trigger>
            <n-button :disabled="disabledReset">
              <template #icon>
                <NIcon>
                  <SaveOutlined />
                </NIcon>
              </template>
              {{ $t('ui.save') }}</n-button>
          </template>
          {{ $t('ui.saveConfirm') }}
        </n-popconfirm>
      </n-space>
    </template>
  </n-card>

</template>

<style scoped>
.n-form-item {
  height: 100%;
}

.textarea {
  height: calc(100dvh - 17rem);
  min-height: 13rem;
  width: 100%;
}
</style>
