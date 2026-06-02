<script setup lang="ts">
import { type Props } from './HomeView'
import {
  NCard, NTag, NSpace, NBadge, NEmpty, NButton
} from 'naive-ui';
const offset = [10, -5] as const

defineProps<Props>()

</script>
<template>
  <div class="grid" v-if="items && items.length">
    <div v-for="item in items" :key="item.id" class="xl:col-3 lg:col-4 md:col-6 col-12">
      <RouterLink :to="`/plugin/${encodeURIComponent(item.id)}`" class="plain-link">
        <n-card :segmented="{
          content: true,
          footer: 'soft',
        }" hoverable>
          <template #header>
            <n-badge :value="item.version" type="success" :offset="offset" v-if="item.version && item.version !== ''">
              {{ item.name }}
            </n-badge>
            <div v-else>
              {{ item.name }}
            </div>
          </template>
          <template #header-extra>
            {{ item.author }}
          </template>

          {{ item.description }}

          <template #footer>
            <n-space>
              <n-tag v-for="platform in item.platform" :key="platform">{{ platform }}</n-tag>
            </n-space>
          </template>
        </n-card>
      </RouterLink>
    </div>
  </div>
  <n-empty size="large" :description="$t('plugins.emby')" v-else>
    <template #extra>
      <RouterLink to="/store">
        <n-button size="small">
          {{ $t('plugins.fromStore') }}
        </n-button>
      </RouterLink>
    </template>
  </n-empty>

</template>
<style scoped>
.n-card {
  height: 100%;
}

.plain-link {
  color: inherit;
  text-decoration: none;
  cursor: pointer;
}
</style>
