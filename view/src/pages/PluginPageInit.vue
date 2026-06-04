<script setup lang="ts">
import InitPage from '@/components/InitPage.vue'
import PluginView from '@/views/PluginView.vue'
import { createPluginPageInit, type Props } from './PluginPageInit';
const props = defineProps<Props>()
const { signal, loader, data } = createPluginPageInit(props)
</script>
<template>
  <InitPage :loader="loader">
    <template #ok>
      <PluginView :signal="signal" :data="data.value!" :plugin="plugin"/>
    </template>
  </InitPage>
</template>
<style scoped></style>
