<script setup lang="ts">
import InitPage from '@/components/InitPage.vue'
import HomeView from '@/views/HomeView.vue'
import { createHomePage } from './HomePage';
const { signal, loader, plugins } = createHomePage()
</script>
<template>
  <InitPage :loader="loader">
    <template #ok>
      <HomeView :signal="signal" :items="plugins.value" />
      <!-- <HomeView :signal="page.abort.signal" :items="[]" /> -->

    </template>
  </InitPage>
</template>
<style scoped></style>
