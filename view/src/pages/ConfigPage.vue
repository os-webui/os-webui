<script setup lang="ts">
import { computed } from 'vue';
import ConfigPageInit from './ConfigPageInit.vue'
import { useRoute } from 'vue-router';
const route = useRoute()
const id = computed(() => `${route.params.id}`)
</script>
<template>
  <ConfigPageInit :key="id" :plugin="id" />
</template>
<style scoped></style>
