<script setup lang="ts">
import { computed } from 'vue';
import PluginPageLoader from './PluginPageLoader.vue'
import { useRoute } from 'vue-router';
const route = useRoute()
const id = computed(() => `${route.params.id}`)
</script>
<template>
  <PluginPageLoader :key="id" :plugin="id" />
</template>
<style scoped></style>
