import { DefaultHttpClient } from "@/internal/api"
import { PageLoader, PageValue } from "@/internal/core/loader"
import type { PluginFeatures } from "@/views/PluginView";
import { ref, onUnmounted, onMounted } from 'vue';
export interface Props {
  plugin: string
}

export function createPluginPageInit(props: Props) {
  const abort = new AbortController()
  const data = new PageValue(() => DefaultHttpClient.get<PluginFeatures>(`/api/v1/plugins/${props.plugin}`, {
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
