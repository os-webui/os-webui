import { DefaultHttpClient } from "@/internal/api"
import { PageLoader, PageValue } from "@/internal/core/loader"
import type { PluginInfo } from "@/views/HomeView";
import { ref, onUnmounted, onMounted } from 'vue';


export function createHomePage() {
  const abort = new AbortController()
  const plugins = new PageValue(() => DefaultHttpClient.get<PluginInfo[]>('/api/v1/plugins', {
    signal: abort.signal,
  }).then((resp) => {
    return resp.data.sort((l, r) => l.id.localeCompare(r.id))
  }))
  const loader = ref(new PageLoader([
    plugins,
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
    plugins: plugins,
  }
}
