<script setup lang="ts">
import InitPage from '@/components/InitPage.vue'
import HomeView from '@/views/HomeView.vue'
import { createHomePage } from './HomePage';
const page = createHomePage()
</script>
<template>
  <InitPage :loader="page.loader.value">
    <template #ok>
      <HomeView :abort="page.abort.signal" :items="page.plugins.value" />
      <!-- <HomeView :abort="page.abort.signal" :items="[]" /> -->

    </template>
  </InitPage>
</template>
<style scoped></style>
