import { DefaultHttpClient } from "@/internal/api"
import { PageLoader, PageValue } from "@/internal/core/loader"
import type { PluginFeatures } from "@/views/PluginView";
import { ref, onUnmounted, onMounted } from 'vue';


export function createPluginPage(id: string) {
  const abort = new AbortController()
  const plugin = new PageValue(() => DefaultHttpClient.get<PluginFeatures>(`/api/v1/plugins/${id}`, {
    signal: abort.signal,
  }).then((resp) => resp.data))
  const loader = ref(new PageLoader([
    plugin,
  ]))

  onMounted(() => {
    loader.value.fetch()
  })
  onUnmounted(() => {
    abort.abort('page close')
  })
  return {
    abort: abort,
    loader: loader,
    plugin: plugin,
  }
}
