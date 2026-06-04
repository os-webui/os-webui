<script setup lang="ts">
import InitPage from '@/ui/InitPage.vue'
import ConfigView from '@/views/ConfigView.vue'
import { createConfigPageInit, type Props } from './ConfigPageInit';
const props = defineProps<Props>()
const { signal, loader, data } = createConfigPageInit(props)
</script>
<template>
  <InitPage :loader="loader">
    <template #ok>
      <ConfigView :signal="signal" :data="data.value!" :plugin="plugin"/>
    </template>
  </InitPage>
</template>
<style scoped></style>
