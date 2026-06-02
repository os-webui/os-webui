<script setup lang="ts">
import InitPage from '@/components/InitPage.vue'
import PluginView from '@/views/PluginView.vue'
import { createPluginPage } from './PluginPage';
interface Props {
  plugin: string
}
const props = defineProps<Props>()
const page = createPluginPage(props.plugin)

</script>
<template>
  <InitPage :loader="page.loader.value">
    <template #ok>
      <PluginView :abort="page.abort.signal" :plugin="page.plugin.value!" />
    </template>
  </InitPage>
</template>
<style scoped></style>
