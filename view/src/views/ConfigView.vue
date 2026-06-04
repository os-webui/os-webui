<script setup lang="ts">
import { useTitle } from '@/stores/title';
import { createConfigView, type Props } from './ConfigView'
import {
  NInput, NCard, NButton, NSpace
} from 'naive-ui';
const props = defineProps<Props>()
const title = useTitle()
title.set(['main.config', props.data.info.name])
const { textValue, handleKeyDown,
  disabled, disabledClear, disabledReset,
  actions } = createConfigView(props)
const autosize = { minRows: 3, maxRows: 5 }
</script>
<template>
  <n-card :title="$t('main.config') + ' - ' + props.data.info.name" :bordered="false">
    <n-input type="textarea" v-model:value="textValue" @keydown="handleKeyDown" placeholder="" :disabled="disabled"
      :autosize="autosize" autofocus show-count autocapitalize="off" autocomplete="off" autocorrect="off"
      spellcheck="false" />

    <template #footer>
      <n-space justify="end">
        <n-button :disabled="disabledClear" @click="actions.clear()">{{ $t('ui.clear') }}</n-button>
        <n-button :disabled="disabledReset" @click="actions.reset()">{{ $t('ui.reset') }}</n-button>
        <n-button :disabled="disabledReset" @click="actions.save()">{{ $t('ui.save') }}</n-button>
      </n-space>
    </template>
  </n-card>

</template>

<style scoped>
.n-input {
  height: 100%;
}

.n-card {
  height: calc(100dvh - 18rem);
  min-height: 20rem;
}
</style>
