import { ref } from 'vue'
import { type PluginInfo } from './HomeView'
import { useRouter } from 'vue-router'
export interface FeatureInfo {
  id: string
  name: string
  description: string
}
export interface PluginFeatures {
  info: PluginInfo
  features: FeatureInfo[]
}
export interface Props {
  signal: AbortSignal
  plugin: string
  data: PluginFeatures
}
export function createPluginView(props: Props) {
  const disabled = ref(false)
  const router = useRouter()
  return {
    disabled,
    actions: {
      config() {
        if (disabled.value) {
          return
        }
        router.push(`/config/${encodeURIComponent(props.plugin)}`)
      },
      reload() {
        if (disabled.value) {
          return
        }
        console.log('reload')
      },
      restart() {
        if (disabled.value) {
          return
        }
        console.log('restart')
      },
      start() {
        if (disabled.value) {
          return
        }
        console.log('start')
      },
      stop() {
        if (disabled.value) {
          return
        }
        console.log('stop')
      },
    },
  }
}
