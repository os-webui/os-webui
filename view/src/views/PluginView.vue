<script setup lang="ts">
import { useTitle } from '@/stores/title';
import { createPluginView, type Props } from './PluginView'
import {
  NCard, NButton, NSpace, NBadge, NDivider
} from 'naive-ui';
const offset = [10, -5] as const
const props = defineProps<Props>()
const title = useTitle()
title.set(props.data.info.name)
const { disabled, actions } = createPluginView(props)

</script>
<template>
  <n-card :segmented="{
    content: true,
    footer: 'soft',
  }">
    <template #header>
      <n-badge :value="data.info.version" type="success" :offset="offset"
        v-if="data.info.version && data.info.version !== ''">
        {{ data.info.name }}
      </n-badge>
      <div v-else>
        {{ data.info.name }}
      </div>
    </template>
    <template #header-extra>
      {{ data.info.author }}
    </template>

    {{ data.info.description }}

    <template #footer>
      <n-space justify="end">
        <n-button :disabled="disabled" @click="actions.config()">{{ $t('main.config') }}</n-button>
        <n-button :disabled="disabled" @click="actions.reload()">{{ $t('main.reload') }}</n-button>
        <n-button :disabled="disabled" @click="actions.restart()">{{ $t('main.restart') }}</n-button>
        <n-button :disabled="disabled" @click="actions.start()">{{ $t('main.start') }}</n-button>
        <n-button :disabled="disabled" @click="actions.stop()">{{ $t('main.stop') }}</n-button>
      </n-space>
    </template>
  </n-card>


  <div v-if="data.features && data.features.length">
    <n-divider />
    <div class="grid">
      <div v-for="item in data.features" :key="item.id" class="xl:col-3 lg:col-4 md:col-6 col-12">
        <RouterLink :to="`/plugin/${encodeURIComponent(data.info.id)}/${encodeURIComponent(item.id)}`"
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
