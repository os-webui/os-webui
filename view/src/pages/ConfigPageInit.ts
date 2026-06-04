import { DefaultHttpClient } from "@/internal/api"
import { PageLoader, PageValue } from "@/internal/core/loader"
import type { PluginInfo } from "@/views/HomeView";

import { ref, onUnmounted, onMounted } from 'vue';
export interface Props {
  plugin: string
}

export function createConfigPageInit(props: Props) {
  const abort = new AbortController()
  const data = new PageValue(() => DefaultHttpClient.get<{
    info: PluginInfo
    data: string
  }>(`/api/v1/plugins/${encodeURIComponent(props.plugin)}/config`, {
    signal: abort.signal,
  }).then((resp) => resp.data))
  const loader = ref(new PageLoader([
    data,
  ]))

  onMounted(() => {
    loader.value.fetch()
  })
  onUnmounted(() => {
    abort.abort('page close')
  })
  return {
    signal: abort.signal,
    loader: loader,
    data: data,
  }
}
