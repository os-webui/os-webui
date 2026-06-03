<script setup lang="ts">
import { useTitle } from '@/stores/title';
import { type Props } from './PluginView'
import {
  NCard, NTag, NSpace, NBadge, NDivider
} from 'naive-ui';
const offset = [10, -5] as const
const props = defineProps<Props>()
const title = useTitle()
title.set(props.plugin.info.name)
</script>
<template>
  <n-card :segmented="{
    content: true,
    footer: 'soft',
  }">
    <template #header>
      <n-badge :value="plugin.info.version" type="success" :offset="offset"
        v-if="plugin.info.version && plugin.info.version !== ''">
        {{ plugin.info.name }}
      </n-badge>
      <div v-else>
        {{ plugin.info.name }}
      </div>
    </template>
    <template #header-extra>
      {{ plugin.info.author }}
    </template>

    {{ plugin.info.description }}

    <template #footer>
      <n-space>
        <n-tag v-for="platform in plugin.info.platform" :key="platform">{{ platform }}</n-tag>
      </n-space>
    </template>
  </n-card>


  <div v-if="plugin.features && plugin.features.length">
    <n-divider />
    <div class="grid">
      <div v-for="item in plugin.features" :key="item.id" class="xl:col-3 lg:col-4 md:col-6 col-12">
        <RouterLink :to="`/plugin/${encodeURIComponent(plugin.info.id)}/${encodeURIComponent(item.id)}`"
          class="plain-link">
          <n-card :segmented="{
            content: true,
            footer: 'soft',
          }" hoverable>
            <template #header>
              <div>
                {{ item.name }}
              </div>
            </template>
            {{ item.description }}
          </n-card>
        </RouterLink>
      </div>
    </div>
  </div>
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
