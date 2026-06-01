<script setup lang="ts">
import { type IPageLoader, PageStatus } from '@/internal/core/loader';
import {
  NSpin, NResult, NButton,
} from 'naive-ui'
import { computed } from 'vue';
const props = defineProps<{
  loader?: IPageLoader
}>()

const err = computed(() => {
  const error = props.loader?.error
  return error instanceof Error
    ? error.message
    : String(error);
})
</script>
<template>
  <div class="view">
    <slot v-if="loader?.status === PageStatus.Ok" name="ok">ok</slot>
    <slot v-else-if="loader?.status === PageStatus.Error" name="error">
      <n-result status="error" :title="$t('main.error')" :description="err">
        <template #footer>
          <n-button @click="loader?.fetch()">{{ $t('main.retry') }}</n-button>
        </template>
      </n-result>
    </slot>
    <slot v-else-if="loader?.status === PageStatus.Loading" name="loading">
      <div class="flex align-items-center justify-content-center">
        <n-spin size="large" />
      </div>
    </slot>
  </div>
</template>

<style scoped></style>
