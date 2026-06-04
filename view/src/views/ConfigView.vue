<script setup lang="ts">
import { useTitle } from '@/stores/title';
import { createConfigView, type Props } from './ConfigView'
import {
  NCard, NButton, NSpace, NFormItem
} from 'naive-ui';
import InputTextarea from '@/ui/input/InputTextarea.vue'
const props = defineProps<Props>()
const title = useTitle()
title.set(['main.config', props.data.info.name])
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
        <n-button :disabled="disabledClear" @click="actions.clear()">{{ $t('ui.clear') }}</n-button>
        <n-button :disabled="disabledReset" @click="actions.reset()">{{ $t('ui.reset') }}</n-button>
        <n-button :disabled="disabledReset" @click="actions.save()">{{ $t('ui.save') }}</n-button>
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
